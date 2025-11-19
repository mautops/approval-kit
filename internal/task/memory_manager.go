package task

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mautops/approval-kit/internal/errors"
	"github.com/mautops/approval-kit/internal/event"
	"github.com/mautops/approval-kit/internal/statemachine"
	"github.com/mautops/approval-kit/internal/template"
	"github.com/mautops/approval-kit/internal/types"
)

// memoryTaskManager 内存实现的任务管理器
type memoryTaskManager struct {
	mu                sync.RWMutex
	tasks             map[string]*Task // taskID -> Task
	templateMgr       template.TemplateManager
	stateMachine      statemachine.StateMachine
	approverFetcherFunc func(*template.Template, *Task) error // 审批人获取函数(可选,用于任务创建时获取动态审批人)
	eventNotifier     *event.EventNotifier // 事件通知器(可选)
}

// NewTaskManager 创建新的任务管理器实例(内存实现)
// templateMgr: 模板管理器,用于获取模板信息
// approverFetcherFunc: 审批人获取函数(可选,用于任务创建时获取动态审批人)
func NewTaskManager(templateMgr template.TemplateManager, approverFetcherFunc func(*template.Template, *Task) error) TaskManager {
	return &memoryTaskManager{
		tasks:              make(map[string]*Task),
		templateMgr:        templateMgr,
		stateMachine:       statemachine.NewStateMachine(),
		approverFetcherFunc: approverFetcherFunc,
		eventNotifier:      nil,
	}
}

// NewTaskManagerWithNotifier 创建带事件通知器的任务管理器实例
// templateMgr: 模板管理器,用于获取模板信息
// approverFetcherFunc: 审批人获取函数(可选,用于任务创建时获取动态审批人)
// notifier: 事件通知器(可选)
func NewTaskManagerWithNotifier(templateMgr template.TemplateManager, approverFetcherFunc func(*template.Template, *Task) error, notifier *event.EventNotifier) TaskManager {
	return &memoryTaskManager{
		tasks:              make(map[string]*Task),
		templateMgr:        templateMgr,
		stateMachine:       statemachine.NewStateMachine(),
		approverFetcherFunc: approverFetcherFunc,
		eventNotifier:      notifier,
	}
}

// Create 基于模板创建审批任务实例
func (m *memoryTaskManager) Create(templateID string, businessID string, params json.RawMessage) (*Task, error) {
	// 获取模板(使用最新版本)
	tpl, err := m.templateMgr.Get(templateID, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get template %q: %w", templateID, err)
	}

	// 生成任务 ID
	taskID := generateTaskID()

	// 创建任务对象
	now := time.Now()
	tsk := &Task{
		ID:             taskID,
		TemplateID:     templateID,
		TemplateVersion: tpl.Version,
		BusinessID:     businessID,
		Params:         params,
		State:          TaskStatePending,
		CurrentNode:    findStartNode(tpl),
		CreatedAt:      now,
		UpdatedAt:      now,
		SubmittedAt:    nil,
		NodeOutputs:    make(map[string]json.RawMessage),
		Approvers:      make(map[string][]string),
		Approvals:      make(map[string]map[string]*Approval),
		Records:        []*Record{},
		StateHistory:   []*StateChange{},
	}

	// 如果参数为 nil,初始化为空 JSON 对象
	if tsk.Params == nil {
		tsk.Params = json.RawMessage("{}")
	}

	// 处理动态审批人(任务创建时获取)
	// 如果提供了审批人获取函数,调用它来获取配置为 on_create 时机的动态审批人
	if m.approverFetcherFunc != nil {
		if err := m.approverFetcherFunc(tpl, tsk); err != nil {
			// 如果获取失败,记录错误但不阻止任务创建
			// 审批人可以在节点激活时重新获取
		}
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// 存储任务
	m.tasks[taskID] = tsk

	// 生成任务创建事件
	if m.eventNotifier != nil {
		m.generateEvent(event.EventTypeTaskCreated, tsk, nil, nil)
	}

	// 返回任务的副本,确保隔离性
	return tsk.Clone(), nil
}

// Get 获取审批任务详情
// 返回任务的快照副本,确保隔离性
func (m *memoryTaskManager) Get(id string) (*Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tsk, exists := m.tasks[id]
	if !exists {
		return nil, fmt.Errorf("task %q not found", id)
	}

	// 返回任务的深拷贝,确保隔离性
	return tsk.Clone(), nil
}

// Submit 提交任务进入审批流程
// 使用状态机进行状态转换,从 pending 转换为 submitted
func (m *memoryTaskManager) Submit(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 获取任务
	tsk, exists := m.tasks[id]
	if !exists {
		return fmt.Errorf("task %q not found", id)
	}

	// 验证当前状态允许提交
	if !m.stateMachine.CanTransition(tsk.GetState(), types.TaskStateSubmitted) {
		return errors.ErrInvalidStateTransition
	}

	// 使用状态机执行状态转换
	// 创建一个适配器,让 Task 实现 TransitionableTask 接口
	adapter := &taskAdapter{task: tsk}
	newTask, err := m.stateMachine.Transition(adapter, types.TaskStateSubmitted, "task submitted")
	if err != nil {
		return fmt.Errorf("state transition failed: %w", err)
	}

	// 更新任务对象
	tsk = newTask.(*taskAdapter).task

	// 设置提交时间
	now := time.Now()
	tsk.SubmittedAt = &now
	tsk.UpdatedAt = now

	// 保存更新后的任务
	m.tasks[id] = tsk

	// 生成任务提交事件
	if m.eventNotifier != nil {
		m.generateEvent(event.EventTypeTaskSubmitted, tsk, nil, nil)
		
		// 生成节点激活事件(提交后当前节点被激活)
		// 如果当前节点是 start,需要找到下一个节点并激活
		if tsk.CurrentNode != "" {
			if tpl, err := m.templateMgr.Get(tsk.TemplateID, 0); err == nil {
				currentNode := tsk.CurrentNode
				// 如果当前节点是 start,找到下一个节点
				if node, exists := tpl.Nodes[currentNode]; exists && node.Type == template.NodeTypeStart {
					// 查找从 start 节点出发的下一个节点
					nextNodeID := findNextNode(tpl, currentNode)
					if nextNodeID != "" {
						if nextNode, exists := tpl.Nodes[nextNodeID]; exists {
							// 更新当前节点为下一个节点
							tsk.CurrentNode = nextNodeID
							// 生成下一个节点的激活事件
							m.generateEvent(event.EventTypeNodeActivated, tsk, nextNode, nil)
						}
					}
				} else if node, exists := tpl.Nodes[currentNode]; exists {
					// 当前节点不是 start,直接生成激活事件
					m.generateEvent(event.EventTypeNodeActivated, tsk, node, nil)
				}
			}
		}
	}

	return nil
}

// Query 查询任务(支持多条件组合)
func (m *memoryTaskManager) Query(filter *TaskFilter) ([]*Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []*Task

	for _, tsk := range m.tasks {
		tsk.mu.RLock()
		matches := true

		// 按状态过滤
		if filter.State != types.TaskState("") && tsk.State != filter.State {
			matches = false
		}

		// 按模板 ID 过滤
		if filter.TemplateID != "" && tsk.TemplateID != filter.TemplateID {
			matches = false
		}

		// 按业务 ID 过滤
		if filter.BusinessID != "" && tsk.BusinessID != filter.BusinessID {
			matches = false
		}

		// 按审批人过滤(查询待审批任务)
		if filter.Approver != "" {
			// 检查该审批人是否在任一节点的审批人列表中
			found := false
			if tsk.Approvers != nil {
				for _, approvers := range tsk.Approvers {
					for _, approver := range approvers {
						if approver == filter.Approver {
							found = true
							break
						}
					}
					if found {
						break
					}
				}
			}
			if !found {
				matches = false
			}
		}

		// 按时间范围过滤
		if !filter.StartTime.IsZero() && tsk.CreatedAt.Before(filter.StartTime) {
			matches = false
		}
		if !filter.EndTime.IsZero() && tsk.CreatedAt.After(filter.EndTime) {
			matches = false
		}

		tsk.mu.RUnlock()

		if matches {
			// 返回任务的副本
			results = append(results, tsk.Clone())
		}
	}

	return results, nil
}

// taskAdapter 适配器,让 Task 实现 TransitionableTask 接口
// 用于解决循环依赖问题
type taskAdapter struct {
	task *Task
}

func (a *taskAdapter) GetState() types.TaskState {
	return a.task.GetState()
}

func (a *taskAdapter) SetState(state types.TaskState) {
	a.task.SetState(state)
}

func (a *taskAdapter) GetUpdatedAt() time.Time {
	return a.task.GetUpdatedAt()
}

func (a *taskAdapter) SetUpdatedAt(t time.Time) {
	a.task.SetUpdatedAt(t)
}

func (a *taskAdapter) GetStateHistory() []*statemachine.StateChange {
	history := a.task.GetStateHistory()
	result := make([]*statemachine.StateChange, len(history))
	for i, sc := range history {
		result[i] = &statemachine.StateChange{
			From:   sc.From,
			To:     sc.To,
			Reason: sc.Reason,
			Time:   sc.Time,
		}
	}
	return result
}

func (a *taskAdapter) AddStateChange(change *statemachine.StateChange) {
	a.task.AddStateChangeRecord(change.From, change.To, change.Reason, change.Time)
}

func (a *taskAdapter) Clone() statemachine.TransitionableTask {
	return &taskAdapter{task: a.task.Clone()}
}

// Approve 和 Reject 方法实现在 approve.go 文件中

// Cancel 取消任务
// 将任务从 pending、submitted 或 approving 状态转换为 cancelled 状态
// 已通过、已拒绝、已取消、已超时的任务不能取消
func (m *memoryTaskManager) Cancel(id string, reason string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 获取任务
	tsk, exists := m.tasks[id]
	if !exists {
		return fmt.Errorf("task %q not found", id)
	}

	// 检查当前状态是否允许取消
	currentState := tsk.GetState()
	if !m.stateMachine.CanTransition(currentState, types.TaskStateCancelled) {
		return fmt.Errorf("cannot cancel task in state %q", currentState)
	}

	// 使用状态机执行状态转换
	adapter := &taskAdapter{task: tsk}
	newTask, err := m.stateMachine.Transition(adapter, types.TaskStateCancelled, reason)
	if err != nil {
		return fmt.Errorf("state transition failed: %w", err)
	}

	// 更新任务对象
	tsk = newTask.(*taskAdapter).task
	tsk.UpdatedAt = time.Now()

	// 保存更新后的任务
	m.tasks[id] = tsk

	// 生成取消事件
	if m.eventNotifier != nil {
		m.generateEvent(event.EventTypeTaskCancelled, tsk, nil, nil)
	}

	return nil
}

// Withdraw 撤回任务
// 将任务从 submitted 或 approving 状态撤回回 pending 状态
// 如果任务已有审批记录,不允许撤回
func (m *memoryTaskManager) Withdraw(id string, reason string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 获取任务
	tsk, exists := m.tasks[id]
	if !exists {
		return fmt.Errorf("task %q not found", id)
	}

	// 检查当前状态是否允许撤回
	currentState := tsk.GetState()
	if currentState != types.TaskStateSubmitted && currentState != types.TaskStateApproving {
		return fmt.Errorf("cannot withdraw task in state %q, only submitted or approving tasks can be withdrawn", currentState)
	}

	// 检查是否有审批记录(如果有,不允许撤回)
	records := tsk.GetRecords()
	if len(records) > 0 {
		return fmt.Errorf("cannot withdraw task with approval records")
	}

	// 验证当前状态允许转换为 pending
	if !m.stateMachine.CanTransition(currentState, types.TaskStatePending) {
		return errors.ErrInvalidStateTransition
	}

	// 使用状态机执行状态转换
	adapter := &taskAdapter{task: tsk}
	newTask, err := m.stateMachine.Transition(adapter, types.TaskStatePending, reason)
	if err != nil {
		return fmt.Errorf("state transition failed: %w", err)
	}

	// 更新任务对象
	tsk = newTask.(*taskAdapter).task

	// 清空提交时间
	tsk.SubmittedAt = nil
	tsk.UpdatedAt = time.Now()

	// 保存更新后的任务
	m.tasks[id] = tsk

	// 生成撤回事件
	if m.eventNotifier != nil {
		m.generateEvent(event.EventTypeTaskWithdrawn, tsk, nil, nil)
	}

	return nil
}

// Transfer 转交审批
// 将审批任务从原审批人转交给新审批人
func (m *memoryTaskManager) Transfer(id string, nodeID string, fromApprover string, toApprover string, reason string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 1. 获取任务
	tsk, exists := m.tasks[id]
	if !exists {
		return fmt.Errorf("task %q not found", id)
	}

	// 2. 获取模板
	tpl, err := m.templateMgr.Get(tsk.TemplateID, 0)
	if err != nil {
		return fmt.Errorf("failed to get template %q: %w", tsk.TemplateID, err)
	}

	// 3. 获取节点配置
	node, exists := tpl.Nodes[nodeID]
	if !exists {
		return fmt.Errorf("node %q not found in template", nodeID)
	}

	// 4. 检查节点类型是否为审批节点
	if node.Type != template.NodeTypeApproval {
		return fmt.Errorf("node %q is not an approval node", nodeID)
	}

	// 5. 获取审批节点配置
	approvalConfig, ok := node.Config.(template.ApprovalNodeConfigAccessor)
	if !ok {
		return fmt.Errorf("node %q config is not ApprovalNodeConfig", nodeID)
	}

	// 6. 检查是否允许转交
	perms, ok := approvalConfig.GetPermissions().(template.OperationPermissionsAccessor)
	if !ok || !perms.AllowTransfer() {
		return fmt.Errorf("transfer is not allowed for node %q", nodeID)
	}

	// 7. 检查原审批人是否是任务的审批人
	tsk.mu.Lock()
	approvers, exists := tsk.Approvers[nodeID]
	if !exists {
		tsk.mu.Unlock()
		return fmt.Errorf("approvers not found for node %q", nodeID)
	}

	// 检查原审批人是否在审批人列表中
	found := false
	for _, approver := range approvers {
		if approver == fromApprover {
			found = true
			break
		}
	}
	if !found {
		tsk.mu.Unlock()
		return fmt.Errorf("user %q is not an approver for node %q", fromApprover, nodeID)
	}

	// 8. 更新审批人列表(移除原审批人,添加新审批人)
	newApprovers := make([]string, 0, len(approvers))
	for _, approver := range approvers {
		if approver != fromApprover {
			newApprovers = append(newApprovers, approver)
		}
	}
	// 检查新审批人是否已在列表中
	toApproverExists := false
	for _, approver := range newApprovers {
		if approver == toApprover {
			toApproverExists = true
			break
		}
	}
	if !toApproverExists {
		newApprovers = append(newApprovers, toApprover)
	}
	tsk.Approvers[nodeID] = newApprovers

	// 9. 更新审批记录(如果有原审批人的审批记录,需要更新)
	if tsk.Approvals != nil && tsk.Approvals[nodeID] != nil {
		// 如果原审批人已有审批记录,将其移除或标记为转交
		if approval, exists := tsk.Approvals[nodeID][fromApprover]; exists {
			// 创建转交记录,保留原审批记录但标记为转交
			// 这里我们保留原审批记录,但会在 Records 中创建转交记录
			_ = approval // 保留原审批记录
		}
	}

	// 10. 生成转交记录
	record := &Record{
		ID:         generateRecordID(),
		TaskID:     id,
		NodeID:     nodeID,
		Approver:   fromApprover,
		Result:     "transfer",
		Comment:    reason,
		CreatedAt:  time.Now(),
		Attachments: []string{},
	}

	// 验证记录
	if err := record.Validate(); err != nil {
		tsk.mu.Unlock()
		return fmt.Errorf("invalid record: %w", err)
	}

	// 添加到记录列表
	tsk.Records = append(tsk.Records, record)

	// 11. 更新任务更新时间
	tsk.UpdatedAt = time.Now()
	tsk.mu.Unlock()

	// 12. 保存更新后的任务
	m.tasks[id] = tsk

	// 13. 生成转交事件
	if m.eventNotifier != nil {
		// 获取节点信息
		var node *template.Node
		if tpl, err := m.templateMgr.Get(tsk.TemplateID, 0); err == nil {
			if n, exists := tpl.Nodes[nodeID]; exists {
				node = n
			}
		}
		m.generateEvent(event.EventTypeApprovalOp, tsk, node, &event.ApprovalInfo{
			NodeID:   nodeID,
			Approver: fromApprover,
			Result:   "transfer",
			Comment:  reason,
		})
	}

	return nil
}

// AddApprover 加签
// 在审批人列表中添加新的审批人
func (m *memoryTaskManager) AddApprover(id string, nodeID string, approver string, reason string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 1. 获取任务
	tsk, exists := m.tasks[id]
	if !exists {
		return fmt.Errorf("task %q not found", id)
	}

	// 2. 获取模板
	tpl, err := m.templateMgr.Get(tsk.TemplateID, 0)
	if err != nil {
		return fmt.Errorf("failed to get template %q: %w", tsk.TemplateID, err)
	}

	// 3. 获取节点配置
	node, exists := tpl.Nodes[nodeID]
	if !exists {
		return fmt.Errorf("node %q not found in template", nodeID)
	}

	// 4. 检查节点类型是否为审批节点
	if node.Type != template.NodeTypeApproval {
		return fmt.Errorf("node %q is not an approval node", nodeID)
	}

	// 5. 获取审批节点配置
	approvalConfig, ok := node.Config.(template.ApprovalNodeConfigAccessor)
	if !ok {
		return fmt.Errorf("node %q config is not ApprovalNodeConfig", nodeID)
	}

	// 6. 检查是否允许加签
	perms, ok := approvalConfig.GetPermissions().(template.OperationPermissionsAccessor)
	if !ok || !perms.AllowAddApprover() {
		return fmt.Errorf("add approver is not allowed for node %q", nodeID)
	}

	// 7. 更新审批人列表
	tsk.mu.Lock()
	if tsk.Approvers == nil {
		tsk.Approvers = make(map[string][]string)
	}

	// 获取当前审批人列表
	approvers, exists := tsk.Approvers[nodeID]
	if !exists {
		approvers = []string{}
	}

	// 检查新审批人是否已在列表中
	for _, existingApprover := range approvers {
		if existingApprover == approver {
			tsk.mu.Unlock()
			return fmt.Errorf("approver %q already exists in node %q", approver, nodeID)
		}
	}

	// 添加新审批人
	approvers = append(approvers, approver)
	tsk.Approvers[nodeID] = approvers

	// 8. 生成加签记录
	record := &Record{
		ID:         generateRecordID(),
		TaskID:     id,
		NodeID:     nodeID,
		Approver:   approver,
		Result:     "add_approver",
		Comment:    reason,
		CreatedAt:  time.Now(),
		Attachments: []string{},
	}

	// 验证记录
	if err := record.Validate(); err != nil {
		tsk.mu.Unlock()
		return fmt.Errorf("invalid record: %w", err)
	}

	// 添加到记录列表
	tsk.Records = append(tsk.Records, record)

	// 9. 更新任务更新时间
	tsk.UpdatedAt = time.Now()
	tsk.mu.Unlock()

	// 10. 保存更新后的任务
	m.tasks[id] = tsk

	// 11. 生成加签事件
	if m.eventNotifier != nil {
		// 获取节点信息
		var node *template.Node
		if tpl, err := m.templateMgr.Get(tsk.TemplateID, 0); err == nil {
			if n, exists := tpl.Nodes[nodeID]; exists {
				node = n
			}
		}
		m.generateEvent(event.EventTypeApprovalOp, tsk, node, &event.ApprovalInfo{
			NodeID:   nodeID,
			Approver: approver,
			Result:   "add_approver",
			Comment:  reason,
		})
	}

	return nil
}

// RemoveApprover 减签
// 从审批人列表中移除指定的审批人
func (m *memoryTaskManager) RemoveApprover(id string, nodeID string, approver string, reason string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 1. 获取任务
	tsk, exists := m.tasks[id]
	if !exists {
		return fmt.Errorf("task %q not found", id)
	}

	// 2. 获取模板
	tpl, err := m.templateMgr.Get(tsk.TemplateID, 0)
	if err != nil {
		return fmt.Errorf("failed to get template %q: %w", tsk.TemplateID, err)
	}

	// 3. 获取节点配置
	node, exists := tpl.Nodes[nodeID]
	if !exists {
		return fmt.Errorf("node %q not found in template", nodeID)
	}

	// 4. 检查节点类型是否为审批节点
	if node.Type != template.NodeTypeApproval {
		return fmt.Errorf("node %q is not an approval node", nodeID)
	}

	// 5. 获取审批节点配置
	approvalConfig, ok := node.Config.(template.ApprovalNodeConfigAccessor)
	if !ok {
		return fmt.Errorf("node %q config is not ApprovalNodeConfig", nodeID)
	}

	// 6. 检查是否允许减签
	perms, ok := approvalConfig.GetPermissions().(template.OperationPermissionsAccessor)
	if !ok || !perms.AllowRemoveApprover() {
		return fmt.Errorf("remove approver is not allowed for node %q", nodeID)
	}

	// 7. 更新审批人列表
	tsk.mu.Lock()
	if tsk.Approvers == nil {
		tsk.mu.Unlock()
		return fmt.Errorf("approvers not found for node %q", nodeID)
	}

	// 获取当前审批人列表
	approvers, exists := tsk.Approvers[nodeID]
	if !exists {
		tsk.mu.Unlock()
		return fmt.Errorf("approvers not found for node %q", nodeID)
	}

	// 检查要移除的审批人是否在列表中
	found := false
	newApprovers := make([]string, 0, len(approvers))
	for _, existingApprover := range approvers {
		if existingApprover == approver {
			found = true
			// 跳过这个审批人,不添加到新列表
			continue
		}
		newApprovers = append(newApprovers, existingApprover)
	}

	if !found {
		tsk.mu.Unlock()
		return fmt.Errorf("approver %q not found in node %q", approver, nodeID)
	}

	// 更新审批人列表
	tsk.Approvers[nodeID] = newApprovers

	// 8. 生成减签记录
	record := &Record{
		ID:         generateRecordID(),
		TaskID:     id,
		NodeID:     nodeID,
		Approver:   approver,
		Result:     "remove_approver",
		Comment:    reason,
		CreatedAt:  time.Now(),
		Attachments: []string{},
	}

	// 验证记录
	if err := record.Validate(); err != nil {
		tsk.mu.Unlock()
		return fmt.Errorf("invalid record: %w", err)
	}

	// 添加到记录列表
	tsk.Records = append(tsk.Records, record)

	// 9. 更新任务更新时间
	tsk.UpdatedAt = time.Now()
	tsk.mu.Unlock()

	// 10. 保存更新后的任务
	m.tasks[id] = tsk

	// 11. 生成减签事件
	if m.eventNotifier != nil {
		// 获取节点信息
		var node *template.Node
		if tpl, err := m.templateMgr.Get(tsk.TemplateID, 0); err == nil {
			if n, exists := tpl.Nodes[nodeID]; exists {
				node = n
			}
		}
		m.generateEvent(event.EventTypeApprovalOp, tsk, node, &event.ApprovalInfo{
			NodeID:   nodeID,
			Approver: approver,
			Result:   "remove_approver",
			Comment:  reason,
		})
	}

	return nil
}

// HandleTimeout 处理任务超时
func (m *memoryTaskManager) HandleTimeout(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 获取任务
	tsk, exists := m.tasks[id]
	if !exists {
		return fmt.Errorf("task %q not found", id)
	}

	// 检查是否超时
	timeout, _ := m.CheckTimeout(tsk)
	if !timeout {
		// 未超时,直接返回
		return nil
	}

	// 验证当前状态允许转换为超时状态
	if !m.stateMachine.CanTransition(tsk.GetState(), types.TaskStateTimeout) {
		return errors.ErrInvalidStateTransition
	}

	// 使用状态机执行状态转换
	adapter := &taskAdapter{task: tsk}
	newTask, err := m.stateMachine.Transition(adapter, types.TaskStateTimeout, "task timeout")
	if err != nil {
		return fmt.Errorf("state transition failed: %w", err)
	}

	// 更新任务对象
	tsk = newTask.(*taskAdapter).task
	tsk.UpdatedAt = time.Now()

	// 保存更新后的任务
	m.tasks[id] = tsk

	// 生成超时事件
	if m.eventNotifier != nil {
		m.generateEvent(event.EventTypeTaskTimeout, tsk, nil, nil)
	}

	return nil
}

// taskIDCounter 任务 ID 计数器,用于确保唯一性
var taskIDCounter int64

// generateTaskID 生成任务 ID
// 使用时间戳 + 原子计数器确保唯一性,即使在同一纳秒内多次调用也能保证不同
func generateTaskID() string {
	// 使用纳秒时间戳 + 原子计数器确保唯一性
	now := time.Now()
	counter := atomic.AddInt64(&taskIDCounter, 1)
	return fmt.Sprintf("task-%d-%d", now.UnixNano(), counter)
}

// findStartNode 查找模板中的开始节点
func findStartNode(tpl *template.Template) string {
	for _, node := range tpl.Nodes {
		if node.Type == template.NodeTypeStart {
			return node.ID
		}
	}
	// 如果找不到开始节点,返回空字符串
	// 这应该不会发生,因为模板验证已经确保有开始节点
	return ""
}

// findNextNode 查找指定节点的下一个节点
func findNextNode(tpl *template.Template, nodeID string) string {
	for _, edge := range tpl.Edges {
		if edge.From == nodeID {
			return edge.To
		}
	}
	return ""
}

