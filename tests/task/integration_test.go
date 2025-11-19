package task_test

import (
	"encoding/json"
	"testing"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
	"github.com/mautops/approval-kit/internal/types"
)

// TestApprovalFlowEndToEnd 测试端到端审批流程
// 完整流程: 创建模板 -> 创建任务 -> 提交 -> 审批 -> 完成
func TestApprovalFlowEndToEnd(t *testing.T) {
	// 1. 创建管理器和模板
	_, taskManager, tmpl := createTestManagers(t)

	// 4. 创建审批任务
	tsk, err := taskManager.Create(tmpl.ID, "business-001", json.RawMessage(`{"amount": 1000}`))
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// 验证初始状态
	if tsk.State != types.TaskStatePending {
		t.Errorf("Task initial state = %q, want %q", tsk.State, types.TaskStatePending)
	}
	if tsk.CurrentNode != "start" {
		t.Errorf("Task initial node = %q, want %q", tsk.CurrentNode, "start")
	}

	// 5. 提交任务
	err = taskManager.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Failed to submit task: %v", err)
	}

	// 验证提交后状态
	tsk, err = taskManager.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}
	if tsk.State != types.TaskStateSubmitted {
		t.Errorf("Task state after submit = %q, want %q", tsk.State, types.TaskStateSubmitted)
	}
	if tsk.SubmittedAt == nil {
		t.Error("Task.SubmittedAt should be set after submit")
	}

	// 6. 执行审批操作
	err = taskManager.Approve(tsk.ID, "approval-001", "user-001", "approved")
	if err != nil {
		t.Fatalf("Failed to approve task: %v", err)
	}

	// 验证审批后状态(单人审批模式,审批后立即完成)
	tsk, err = taskManager.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}
	if tsk.State != types.TaskStateApproved {
		t.Errorf("Task state after approve = %q, want %q", tsk.State, types.TaskStateApproved)
	}

	// 验证审批记录
	if len(tsk.Records) == 0 {
		t.Fatal("Task should have approval records")
	}
	record := tsk.Records[0]
	if record.Result != "approve" {
		t.Errorf("Record.Result = %q, want %q", record.Result, "approve")
	}
	if record.Approver != "user-001" {
		t.Errorf("Record.Approver = %q, want %q", record.Approver, "user-001")
	}

	// 验证审批结果已记录
	approvals, exists := tsk.Approvals["approval-001"]
	if !exists {
		t.Fatal("Task should have approvals for node approval-001")
	}
	approval, exists := approvals["user-001"]
	if !exists {
		t.Fatal("Task should have approval from user-001")
	}
	if approval.Result != "approve" {
		t.Errorf("Approval.Result = %q, want %q", approval.Result, "approve")
	}
}

// TestApprovalFlowWithRejection 测试包含拒绝的审批流程
func TestApprovalFlowWithRejection(t *testing.T) {
	// 创建管理器和模板
	_, taskManager, tmpl := createTestManagers(t)

	// 创建任务
	tsk, err := taskManager.Create(tmpl.ID, "business-002", nil)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// 提交任务
	err = taskManager.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Failed to submit task: %v", err)
	}

	// 执行拒绝操作
	err = taskManager.Reject(tsk.ID, "approval-001", "user-001", "rejected")
	if err != nil {
		t.Fatalf("Failed to reject task: %v", err)
	}

	// 验证拒绝记录
	tsk, err = taskManager.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}

	// 验证拒绝记录
	if len(tsk.Records) == 0 {
		t.Fatal("Task should have rejection records")
	}
	record := tsk.Records[0]
	if record.Result != "reject" {
		t.Errorf("Record.Result = %q, want %q", record.Result, "reject")
	}
	if record.Comment != "rejected" {
		t.Errorf("Record.Comment = %q, want %q", record.Comment, "rejected")
	}

	// 验证审批结果已记录
	approvals, exists := tsk.Approvals["approval-001"]
	if !exists {
		t.Fatal("Task should have approvals for node approval-001")
	}
	approval, exists := approvals["user-001"]
	if !exists {
		t.Fatal("Task should have approval from user-001")
	}
	if approval.Result != "reject" {
		t.Errorf("Approval.Result = %q, want %q", approval.Result, "reject")
	}
}

// TestApprovalFlowMultipleTasks 测试多个任务的审批流程
func TestApprovalFlowMultipleTasks(t *testing.T) {
	// 创建管理器和模板
	_, taskManager, tmpl := createTestManagers(t)

	// 创建多个任务
	taskIDs := make([]string, 3)
	for i := 0; i < 3; i++ {
		tsk, err := taskManager.Create(tmpl.ID, "business-003", nil)
		if err != nil {
			t.Fatalf("Failed to create task %d: %v", i, err)
		}
		taskIDs[i] = tsk.ID

		// 提交任务
		err = taskManager.Submit(tsk.ID)
		if err != nil {
			t.Fatalf("Failed to submit task %d: %v", i, err)
		}
	}

	// 验证所有任务都已提交
	for i, taskID := range taskIDs {
		tsk, err := taskManager.Get(taskID)
		if err != nil {
			t.Fatalf("Failed to get task %d: %v", i, err)
		}
		if tsk.State != types.TaskStateSubmitted {
			t.Errorf("Task %d state = %q, want %q", i, tsk.State, types.TaskStateSubmitted)
		}
	}
}

// createSimpleApprovalTemplate 创建简单的审批模板
// 包含: 开始节点 -> 审批节点 -> 结束节点
// 这是一个通用的测试辅助函数,可以在多个测试中复用
func createSimpleApprovalTemplate() *template.Template {
	return &template.Template{
		ID:   "simple-template",
		Name: "Simple Approval Template",
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "Start Node",
				Type: template.NodeTypeStart,
			},
			"approval-001": {
				ID:   "approval-001",
				Name: "Approval Node",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeSingle,
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{"user-001"},
					},
				},
			},
			"end": {
				ID:   "end",
				Name: "End Node",
				Type: template.NodeTypeEnd,
			},
		},
		Edges: []*template.Edge{
			{From: "start", To: "approval-001"},
			{From: "approval-001", To: "end"},
		},
		Version: 1,
	}
}

// createTestManagers 创建测试用的管理器和模板
// 返回: 模板管理器、任务管理器和已创建的模板
func createTestManagers(t *testing.T) (template.TemplateManager, task.TaskManager, *template.Template) {
	templateManager := template.NewTemplateManager()
	taskManager := task.NewTaskManager(templateManager, nil)

	tmpl := createSimpleApprovalTemplate()
	err := templateManager.Create(tmpl)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	return templateManager, taskManager, tmpl
}

// TestApprovalFlowWithWithdraw 测试包含撤回的审批流程
func TestApprovalFlowWithWithdraw(t *testing.T) {
	// 创建管理器和模板
	_, taskManager, tmpl := createTestManagers(t)

	// 创建任务
	tsk, err := taskManager.Create(tmpl.ID, "business-004", nil)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// 提交任务
	err = taskManager.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Failed to submit task: %v", err)
	}

	// 验证提交后状态
	tsk, err = taskManager.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}
	if tsk.State != types.TaskStateSubmitted {
		t.Errorf("Task state after submit = %q, want %q", tsk.State, types.TaskStateSubmitted)
	}

	// 撤回任务
	err = taskManager.Withdraw(tsk.ID, "withdraw reason")
	if err != nil {
		t.Fatalf("Failed to withdraw task: %v", err)
	}

	// 验证撤回后状态
	tsk, err = taskManager.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}
	if tsk.State != types.TaskStatePending {
		t.Errorf("Task state after withdraw = %q, want %q", tsk.State, types.TaskStatePending)
	}
}

// TestApprovalFlowWithRejectJump 测试拒绝后跳转的审批流程
func TestApprovalFlowWithRejectJump(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	
	// 创建包含拒绝后跳转配置的模板
	tpl := &template.Template{
		ID:   "jump-template",
		Name: "Reject Jump Template",
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "Start Node",
				Type: template.NodeTypeStart,
			},
			"approval-001": {
				ID:   "approval-001",
				Name: "Approval Node 1",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode:            node.ApprovalModeSingle,
					RejectBehavior:  node.RejectBehaviorJump,
					RejectTargetNode: "approval-002",
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{"user-001"},
					},
				},
			},
			"approval-002": {
				ID:   "approval-002",
				Name: "Approval Node 2",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeSingle,
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{"user-002"},
					},
				},
			},
			"end": {
				ID:   "end",
				Name: "End Node",
				Type: template.NodeTypeEnd,
			},
		},
		Edges: []*template.Edge{
			{From: "start", To: "approval-001"},
			{From: "approval-001", To: "approval-002"},
			{From: "approval-002", To: "end"},
		},
		Version: 1,
	}
	
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	// 创建任务管理器
	taskMgr := task.NewTaskManager(templateMgr, nil)

	// 创建任务
	tsk, err := taskMgr.Create("jump-template", "business-005", nil)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// 提交任务
	err = taskMgr.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Failed to submit task: %v", err)
	}

	// 拒绝审批(应该跳转到 approval-002)
	err = taskMgr.Reject(tsk.ID, "approval-001", "user-001", "rejected")
	if err != nil {
		t.Fatalf("Failed to reject task: %v", err)
	}

	// 验证任务已跳转到目标节点
	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}
	if tsk.CurrentNode != "approval-002" {
		t.Errorf("Task.CurrentNode = %q, want %q", tsk.CurrentNode, "approval-002")
	}
	if tsk.State != types.TaskStateApproving {
		t.Errorf("Task.State = %q, want %q", tsk.State, types.TaskStateApproving)
	}
}

