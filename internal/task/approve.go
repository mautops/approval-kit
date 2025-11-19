package task

import (
	"fmt"
	"time"

	"github.com/mautops/approval-kit/internal/errors"
	"github.com/mautops/approval-kit/internal/event"
	"github.com/mautops/approval-kit/internal/template"
	"github.com/mautops/approval-kit/internal/types"
)

// Approve 审批人进行同意操作
func (m *memoryTaskManager) Approve(id string, nodeID string, approver string, comment string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 1. 获取任务
	tsk, exists := m.tasks[id]
	if !exists {
		return fmt.Errorf("task %q not found", id)
	}

	// 2. 验证任务状态(只有 submitted 或 approving 状态才能审批)
	tsk.mu.RLock()
	state := tsk.State
	tsk.mu.RUnlock()

	if state != types.TaskStateSubmitted && state != types.TaskStateApproving {
		return fmt.Errorf("%w: task state %q cannot be approved", errors.ErrInvalidStateTransition, state)
	}

	// 2.1 获取模板和节点配置,验证审批意见必填
	tpl, err := m.templateMgr.Get(tsk.TemplateID, 0)
	if err != nil {
		return fmt.Errorf("failed to get template %q: %w", tsk.TemplateID, err)
	}

	node, exists := tpl.Nodes[nodeID]
	if !exists {
		return fmt.Errorf("node %q not found in template", nodeID)
	}

	if node.Type == template.NodeTypeApproval {
		approvalConfig, ok := node.Config.(template.ApprovalNodeConfigAccessor)
		if ok && approvalConfig.RequireComment() {
			if comment == "" {
				return fmt.Errorf("comment is required for approval node %q", nodeID)
			}
		}
	}

	// 3. 更新任务状态为 approving(如果还是 submitted)
	tsk.mu.Lock()
	if tsk.State == types.TaskStateSubmitted {
		tsk.State = types.TaskStateApproving
		tsk.UpdatedAt = time.Now()
	}

	// 4. 记录审批结果
	// 初始化节点审批记录
	if tsk.Approvals == nil {
		tsk.Approvals = make(map[string]map[string]*Approval)
	}
	if tsk.Approvals[nodeID] == nil {
		tsk.Approvals[nodeID] = make(map[string]*Approval)
	}

	// 记录审批结果
	tsk.Approvals[nodeID][approver] = &Approval{
		Result:    "approve",
		Comment:   comment,
		CreatedAt: time.Now(),
	}

	// 5. 生成审批记录
	record := &Record{
		ID:          generateRecordID(),
		TaskID:      id,
		NodeID:      nodeID,
		Approver:    approver,
		Result:      "approve",
		Comment:     comment,
		CreatedAt:   time.Now(),
		Attachments: []string{},
	}

	// 验证记录
	if err := record.Validate(); err != nil {
		tsk.mu.Unlock()
		return fmt.Errorf("invalid record: %w", err)
	}

	// 添加到记录列表
	tsk.Records = append(tsk.Records, record)

	// 6. 检查审批是否完成(对于单人审批模式,审批人同意后立即完成)
	// 获取审批人列表和当前状态
	approvers := tsk.Approvers[nodeID]
	currentState := tsk.State
	shouldTransition := false

	// 获取节点配置以确定审批模式
	if node.Type == template.NodeTypeApproval {
		// 尝试从配置中获取审批人列表(如果任务中的审批人列表为空)
		if len(approvers) == 0 {
			// 尝试通过类型断言获取固定审批人配置
			// 注意: 这里我们无法直接访问 node 包的类型,所以简化处理
			// 如果只有一个审批记录,且是当前审批人的同意记录,可能是单人审批模式
			if len(tsk.Approvals[nodeID]) == 1 {
				if approval, exists := tsk.Approvals[nodeID][approver]; exists && approval.Result == "approve" {
					// 对于单人审批模式,如果审批人已同意,将任务状态转换为已通过
					if m.stateMachine.CanTransition(currentState, types.TaskStateApproved) {
						shouldTransition = true
					}
				}
			}
		} else if len(approvers) == 1 && approvers[0] == approver {
			// 如果只有一个审批人且就是当前审批人,且已同意,将任务状态转换为已通过
			// 这是单人审批模式
			if m.stateMachine.CanTransition(currentState, types.TaskStateApproved) {
				shouldTransition = true
			}
		} else if len(approvers) > 1 {
			// 多人审批模式: 检查是否所有审批人都已同意
			// 对于多人会签模式,需要所有审批人都同意
			allApproved := true
			for _, approverID := range approvers {
				approval, exists := tsk.Approvals[nodeID][approverID]
				if !exists || approval == nil || approval.Result != "approve" {
					allApproved = false
					break
				}
			}
			// 如果所有审批人都已同意,则完成
			// 注意: 这里我们简化处理,只要所有审批人都已同意就转换状态
			// 实际应该根据审批模式(会签/或签/比例/顺序)来决定
			if allApproved {
				if m.stateMachine.CanTransition(currentState, types.TaskStateApproved) {
					shouldTransition = true
				}
			}
		}
	}

	// 7. 更新任务更新时间
	tsk.UpdatedAt = time.Now()
	tsk.mu.Unlock()

	// 8. 如果需要转换状态,执行状态转换(在释放锁之后)
	if shouldTransition {
		// 重新获取任务(因为锁已释放)
		tsk = m.tasks[id]
		adapter := &taskAdapter{task: tsk}
		newTask, err := m.stateMachine.Transition(adapter, types.TaskStateApproved, "all approvers approved")
		if err == nil {
			tsk = newTask.(*taskAdapter).task
			tsk.mu.Lock()
			tsk.UpdatedAt = time.Now()
			tsk.mu.Unlock()
			m.tasks[id] = tsk
		}
	} else {
		// 保存更新后的任务
		m.tasks[id] = tsk
	}

	// 9. 生成审批事件
	if m.eventNotifier != nil {
		approvalInfo := &event.ApprovalInfo{
			NodeID:   nodeID,
			Approver: approver,
			Result:   "approve",
			Comment:  comment,
		}
		m.generateEvent(event.EventTypeApprovalOp, tsk, node, approvalInfo)
		
		// 如果状态已转换为 approved,生成任务通过事件和节点完成事件
		if shouldTransition && tsk.State == types.TaskStateApproved {
			// 生成节点完成事件
			m.generateEvent(event.EventTypeNodeCompleted, tsk, node, nil)
			// 生成任务通过事件
			m.generateEvent(event.EventTypeTaskApproved, tsk, node, nil)
		}
	}

	return nil
}

// ApproveWithAttachments 审批人进行同意操作(带附件)
func (m *memoryTaskManager) ApproveWithAttachments(id string, nodeID string, approver string, comment string, attachments []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 1. 获取任务
	tsk, exists := m.tasks[id]
	if !exists {
		return fmt.Errorf("task %q not found", id)
	}

	// 2. 验证任务状态(只有 submitted 或 approving 状态才能审批)
	tsk.mu.RLock()
	state := tsk.State
	tsk.mu.RUnlock()

	if state != types.TaskStateSubmitted && state != types.TaskStateApproving {
		return fmt.Errorf("%w: task state %q cannot be approved", errors.ErrInvalidStateTransition, state)
	}

	// 2.1 获取模板和节点配置,验证审批意见和附件要求
	tpl, err := m.templateMgr.Get(tsk.TemplateID, 0)
	if err != nil {
		return fmt.Errorf("failed to get template %q: %w", tsk.TemplateID, err)
	}

	node, exists := tpl.Nodes[nodeID]
	if !exists {
		return fmt.Errorf("node %q not found in template", nodeID)
	}

	if node.Type == template.NodeTypeApproval {
		approvalConfig, ok := node.Config.(template.ApprovalNodeConfigAccessor)
		if ok {
			if approvalConfig.RequireComment() && comment == "" {
				return fmt.Errorf("comment is required for approval node %q", nodeID)
			}
			if approvalConfig.RequireAttachments() && len(attachments) == 0 {
				return fmt.Errorf("attachments are required for approval node %q", nodeID)
			}
		}
	}

	// 3. 更新任务状态为 approving(如果还是 submitted)
	tsk.mu.Lock()
	if tsk.State == types.TaskStateSubmitted {
		tsk.State = types.TaskStateApproving
		tsk.UpdatedAt = time.Now()
	}

	// 4. 记录审批结果
	// 初始化节点审批记录
	if tsk.Approvals == nil {
		tsk.Approvals = make(map[string]map[string]*Approval)
	}
	if tsk.Approvals[nodeID] == nil {
		tsk.Approvals[nodeID] = make(map[string]*Approval)
	}

	// 记录审批结果
	tsk.Approvals[nodeID][approver] = &Approval{
		Result:    "approve",
		Comment:   comment,
		CreatedAt: time.Now(),
	}

	// 5. 生成审批记录
	record := &Record{
		ID:          generateRecordID(),
		TaskID:      id,
		NodeID:      nodeID,
		Approver:    approver,
		Result:      "approve",
		Comment:     comment,
		CreatedAt:   time.Now(),
		Attachments: attachments,
	}

	// 验证记录
	if err := record.Validate(); err != nil {
		tsk.mu.Unlock()
		return fmt.Errorf("invalid record: %w", err)
	}

	// 添加到记录列表
	tsk.Records = append(tsk.Records, record)

	// 6. 更新任务更新时间
	tsk.UpdatedAt = time.Now()
	tsk.mu.Unlock()

	return nil
}

// Reject 审批人进行拒绝操作
func (m *memoryTaskManager) Reject(id string, nodeID string, approver string, comment string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 1. 获取任务
	tsk, exists := m.tasks[id]
	if !exists {
		return fmt.Errorf("task %q not found", id)
	}

	// 2. 验证任务状态(只有 submitted 或 approving 状态才能拒绝)
	tsk.mu.RLock()
	state := tsk.State
	tsk.mu.RUnlock()

	if state != types.TaskStateSubmitted && state != types.TaskStateApproving {
		return fmt.Errorf("%w: task state %q cannot be rejected", errors.ErrInvalidStateTransition, state)
	}

	// 2.1 获取模板和节点配置,验证审批意见必填
	tpl, err := m.templateMgr.Get(tsk.TemplateID, 0)
	if err != nil {
		return fmt.Errorf("failed to get template %q: %w", tsk.TemplateID, err)
	}

	node, exists := tpl.Nodes[nodeID]
	if !exists {
		return fmt.Errorf("node %q not found in template", nodeID)
	}

	if node.Type == template.NodeTypeApproval {
		approvalConfig, ok := node.Config.(template.ApprovalNodeConfigAccessor)
		if ok && approvalConfig.RequireComment() {
			if comment == "" {
				return fmt.Errorf("comment is required for approval node %q", nodeID)
			}
		}
	}

	// 3. 更新任务状态为 approving(如果还是 submitted)
	tsk.mu.Lock()
	if tsk.State == types.TaskStateSubmitted {
		tsk.State = types.TaskStateApproving
		tsk.UpdatedAt = time.Now()
	}

	// 4. 记录审批结果
	// 初始化节点审批记录
	if tsk.Approvals == nil {
		tsk.Approvals = make(map[string]map[string]*Approval)
	}
	if tsk.Approvals[nodeID] == nil {
		tsk.Approvals[nodeID] = make(map[string]*Approval)
	}

	// 记录审批结果
	tsk.Approvals[nodeID][approver] = &Approval{
		Result:    "reject",
		Comment:   comment,
		CreatedAt: time.Now(),
	}

	// 5. 生成审批记录
	record := &Record{
		ID:          generateRecordID(),
		TaskID:      id,
		NodeID:      nodeID,
		Approver:    approver,
		Result:      "reject",
		Comment:     comment,
		CreatedAt:   time.Now(),
		Attachments: []string{},
	}

	// 验证记录
	if err := record.Validate(); err != nil {
		tsk.mu.Unlock()
		return fmt.Errorf("invalid record: %w", err)
	}

	// 添加到记录列表
	tsk.Records = append(tsk.Records, record)

	// 6. 处理拒绝后行为
	// 先释放任务锁,避免在状态机转换时死锁
	// 但保持管理器锁,确保任务不会被其他操作修改
	tsk.mu.Unlock()
	
	// 记录拒绝前的状态,用于事件生成
	rejectBeforeState := types.TaskStateApproving

	if node.Type == template.NodeTypeApproval {
		approvalConfig, ok := node.Config.(template.ApprovalNodeConfigAccessor)
		if ok {
			rejectBehavior := approvalConfig.GetRejectBehavior()
			switch rejectBehavior {
			case "terminate":
				// 拒绝后终止流程
				// 重新获取任务锁(管理器锁已持有)
				tsk = m.tasks[id]
				tsk.mu.Lock()
				adapter := &taskAdapter{task: tsk}
				tsk.mu.Unlock()

				newTaskAdapter, err := m.stateMachine.Transition(adapter, types.TaskStateRejected, "rejected by approver")
				if err != nil {
					return fmt.Errorf("failed to transition to rejected state: %w", err)
				}
				m.tasks[id] = newTaskAdapter.(*taskAdapter).task
				tsk = m.tasks[id]
			case "rollback":
				// 拒绝后回退到上一节点
				// 查找上一节点(通过 Edges 查找)
				prevNodeID := ""
				for _, edge := range tpl.Edges {
					if edge.To == nodeID {
						prevNodeID = edge.From
						break
					}
				}
				if prevNodeID == "" {
					// 没有上一节点,终止流程
					tsk = m.tasks[id]
					tsk.mu.Lock()
					adapter := &taskAdapter{task: tsk}
					tsk.mu.Unlock()

					newTaskAdapter, err := m.stateMachine.Transition(adapter, types.TaskStateRejected, "rejected by approver, no previous node")
					if err != nil {
						return fmt.Errorf("failed to transition to rejected state: %w", err)
					}
					m.tasks[id] = newTaskAdapter.(*taskAdapter).task
					tsk = m.tasks[id]
				} else {
					// 跳转到上一节点
					tsk = m.tasks[id]
					tsk.mu.Lock()
					tsk.CurrentNode = prevNodeID
					tsk.State = types.TaskStateApproving
					tsk.UpdatedAt = time.Now()
					tsk.mu.Unlock()
				}
			case "jump":
				// 拒绝后跳转到指定节点
				targetNodeID := approvalConfig.GetRejectTargetNode()
				if targetNodeID == "" {
					// 未指定目标节点,终止流程
					tsk = m.tasks[id]
					tsk.mu.Lock()
					adapter := &taskAdapter{task: tsk}
					tsk.mu.Unlock()

					newTaskAdapter, err := m.stateMachine.Transition(adapter, types.TaskStateRejected, "rejected by approver, no target node")
					if err != nil {
						return fmt.Errorf("failed to transition to rejected state: %w", err)
					}
					m.tasks[id] = newTaskAdapter.(*taskAdapter).task
					tsk = m.tasks[id]
				} else {
					// 验证目标节点存在
					if _, exists := tpl.Nodes[targetNodeID]; !exists {
						return fmt.Errorf("reject target node %q not found in template", targetNodeID)
					}
					// 跳转到目标节点
					tsk = m.tasks[id]
					tsk.mu.Lock()
					tsk.CurrentNode = targetNodeID
					tsk.State = types.TaskStateApproving
					tsk.UpdatedAt = time.Now()
					tsk.mu.Unlock()
				}
			default:
				// 默认行为: 终止流程
				tsk = m.tasks[id]
				tsk.mu.Lock()
				adapter := &taskAdapter{task: tsk}
				tsk.mu.Unlock()

				newTaskAdapter, err := m.stateMachine.Transition(adapter, types.TaskStateRejected, "rejected by approver")
				if err != nil {
					return fmt.Errorf("failed to transition to rejected state: %w", err)
				}
				m.tasks[id] = newTaskAdapter.(*taskAdapter).task
				tsk = m.tasks[id]
			}
		} else {
			// 无法获取配置,默认终止流程
			tsk = m.tasks[id]
			tsk.mu.Lock()
			adapter := &taskAdapter{task: tsk}
			tsk.mu.Unlock()

			newTaskAdapter, err := m.stateMachine.Transition(adapter, types.TaskStateRejected, "rejected by approver")
			if err != nil {
				return fmt.Errorf("failed to transition to rejected state: %w", err)
			}
			m.tasks[id] = newTaskAdapter.(*taskAdapter).task
			tsk = m.tasks[id]
		}
	} else {
		// 非审批节点,默认终止流程
		tsk = m.tasks[id]
		tsk.mu.Lock()
		adapter := &taskAdapter{task: tsk}
		tsk.mu.Unlock()

		newTaskAdapter, err := m.stateMachine.Transition(adapter, types.TaskStateRejected, "rejected by approver")
		if err != nil {
			return fmt.Errorf("failed to transition to rejected state: %w", err)
		}
		m.tasks[id] = newTaskAdapter.(*taskAdapter).task
		tsk = m.tasks[id]
	}

	// 7. 更新任务更新时间(如果需要)
	tsk = m.tasks[id]
	tsk.mu.Lock()
	tsk.UpdatedAt = time.Now()
	tsk.mu.Unlock()

	// 8. 生成事件
	if m.eventNotifier != nil {
		// 如果拒绝前状态是 approving,先生成审批操作事件
		if rejectBeforeState == types.TaskStateApproving {
			m.generateEvent(event.EventTypeApprovalOp, tsk, node, &event.ApprovalInfo{
				NodeID:   nodeID,
				Approver: approver,
				Result:   "reject",
				Comment:  comment,
			})
		}
		
		// 根据最终状态生成相应事件
		if tsk.State == types.TaskStateRejected {
			// 生成节点完成事件
			m.generateEvent(event.EventTypeNodeCompleted, tsk, node, nil)
			// 生成任务拒绝事件
			m.generateEvent(event.EventTypeTaskRejected, tsk, node, nil)
		}
	}

	return nil
}

// RejectWithAttachments 审批人进行拒绝操作(带附件)
func (m *memoryTaskManager) RejectWithAttachments(id string, nodeID string, approver string, comment string, attachments []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 1. 获取任务
	tsk, exists := m.tasks[id]
	if !exists {
		return fmt.Errorf("task %q not found", id)
	}

	// 2. 验证任务状态(只有 submitted 或 approving 状态才能拒绝)
	tsk.mu.RLock()
	state := tsk.State
	tsk.mu.RUnlock()

	if state != types.TaskStateSubmitted && state != types.TaskStateApproving {
		return fmt.Errorf("%w: task state %q cannot be rejected", errors.ErrInvalidStateTransition, state)
	}

	// 2.1 获取模板和节点配置,验证审批意见和附件要求
	tpl, err := m.templateMgr.Get(tsk.TemplateID, 0)
	if err != nil {
		return fmt.Errorf("failed to get template %q: %w", tsk.TemplateID, err)
	}

	node, exists := tpl.Nodes[nodeID]
	if !exists {
		return fmt.Errorf("node %q not found in template", nodeID)
	}

	if node.Type == template.NodeTypeApproval {
		approvalConfig, ok := node.Config.(template.ApprovalNodeConfigAccessor)
		if ok {
			if approvalConfig.RequireComment() && comment == "" {
				return fmt.Errorf("comment is required for approval node %q", nodeID)
			}
			if approvalConfig.RequireAttachments() && len(attachments) == 0 {
				return fmt.Errorf("attachments are required for approval node %q", nodeID)
			}
		}
	}

	// 3. 更新任务状态为 approving(如果还是 submitted)
	tsk.mu.Lock()
	if tsk.State == types.TaskStateSubmitted {
		tsk.State = types.TaskStateApproving
		tsk.UpdatedAt = time.Now()
	}

	// 4. 记录审批结果
	// 初始化节点审批记录
	if tsk.Approvals == nil {
		tsk.Approvals = make(map[string]map[string]*Approval)
	}
	if tsk.Approvals[nodeID] == nil {
		tsk.Approvals[nodeID] = make(map[string]*Approval)
	}

	// 记录审批结果
	tsk.Approvals[nodeID][approver] = &Approval{
		Result:    "reject",
		Comment:   comment,
		CreatedAt: time.Now(),
	}

	// 5. 生成审批记录
	record := &Record{
		ID:          generateRecordID(),
		TaskID:      id,
		NodeID:      nodeID,
		Approver:    approver,
		Result:      "reject",
		Comment:     comment,
		CreatedAt:   time.Now(),
		Attachments: attachments,
	}

	// 验证记录
	if err := record.Validate(); err != nil {
		tsk.mu.Unlock()
		return fmt.Errorf("invalid record: %w", err)
	}

	// 添加到记录列表
	tsk.Records = append(tsk.Records, record)

	// 6. 更新任务更新时间
	tsk.UpdatedAt = time.Now()
	tsk.mu.Unlock()

	return nil
}

// generateRecordID 生成审批记录 ID
// 使用时间戳和原子计数器确保唯一性
func generateRecordID() string {
	// 使用纳秒时间戳 + 原子计数器确保唯一性
	// 即使在同一纳秒内多次调用也能保证唯一性
	now := time.Now()
	// 使用纳秒时间戳 + 微秒时间戳 + 纳秒部分的组合
	// 这样可以确保即使在同一纳秒内多次调用也能生成不同的 ID
	return fmt.Sprintf("record-%d-%d-%d", now.UnixNano(), now.UnixMicro(), now.Nanosecond()%1000)
}
