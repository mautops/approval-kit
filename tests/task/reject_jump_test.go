package task_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
	"github.com/mautops/approval-kit/internal/types"
)

// TestRejectJumpToTargetNode 测试拒绝后跳转到指定节点
func TestRejectJumpToTargetNode(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateWithRejectJump()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	// 创建任务管理器
	taskMgr := task.NewTaskManager(templateMgr, nil)

	// 创建任务
	params := json.RawMessage(`{"amount": 1000}`)
	tsk, err := taskMgr.Create("tpl-001", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 提交任务
	err = taskMgr.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Submit() failed: %v", err)
	}

	// 拒绝审批
	err = taskMgr.Reject(tsk.ID, "approval-001", "user-001", "rejected")
	if err != nil {
		t.Fatalf("Reject() failed: %v", err)
	}

	// 获取任务
	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	// 验证任务已跳转到目标节点
	if tsk.CurrentNode != "approval-002" {
		t.Errorf("CurrentNode = %q, want %q", tsk.CurrentNode, "approval-002")
	}

	// 验证任务状态
	if tsk.State != types.TaskStateApproving {
		t.Errorf("Task.State = %q, want %q", tsk.State, types.TaskStateApproving)
	}
}

// TestRejectRollback 测试拒绝后回滚到上一节点
func TestRejectRollback(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateWithRejectRollback()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	// 创建任务管理器
	taskMgr := task.NewTaskManager(templateMgr, nil)

	// 创建任务
	params := json.RawMessage(`{"amount": 1000}`)
	tsk, err := taskMgr.Create("tpl-001", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 提交任务
	err = taskMgr.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Submit() failed: %v", err)
	}

	// 拒绝审批
	err = taskMgr.Reject(tsk.ID, "approval-002", "user-002", "rejected")
	if err != nil {
		t.Fatalf("Reject() failed: %v", err)
	}

	// 获取任务
	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	// 验证任务已回滚到上一节点
	if tsk.CurrentNode != "approval-001" {
		t.Errorf("CurrentNode = %q, want %q", tsk.CurrentNode, "approval-001")
	}
}

// TestRejectTerminate 测试拒绝后终止流程
func TestRejectTerminate(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateWithRejectTerminate()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	// 创建任务管理器
	taskMgr := task.NewTaskManager(templateMgr, nil)

	// 创建任务
	params := json.RawMessage(`{"amount": 1000}`)
	tsk, err := taskMgr.Create("tpl-001", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 提交任务
	err = taskMgr.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Submit() failed: %v", err)
	}

	// 拒绝审批
	err = taskMgr.Reject(tsk.ID, "approval-001", "user-001", "rejected")
	if err != nil {
		t.Fatalf("Reject() failed: %v", err)
	}

	// 获取任务
	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	// 验证任务状态已变为已拒绝
	if tsk.State != types.TaskStateRejected {
		t.Errorf("Task.State = %q, want %q", tsk.State, types.TaskStateRejected)
	}
}

// createTestTemplateWithRejectJump 创建包含拒绝后跳转配置的测试模板
func createTestTemplateWithRejectJump() *template.Template {
	return &template.Template{
		ID:          "tpl-001",
		Name:        "Test Template",
		Description: "Test template description",
		Version:     1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
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
					Mode:    node.ApprovalModeSingle,
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
	}
}

// createTestTemplateWithRejectRollback 创建包含拒绝后回滚配置的测试模板
func createTestTemplateWithRejectRollback() *template.Template {
	return &template.Template{
		ID:          "tpl-001",
		Name:        "Test Template",
		Description: "Test template description",
		Version:     1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
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
					Mode:    node.ApprovalModeSingle,
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
					Mode:           node.ApprovalModeSingle,
					RejectBehavior: node.RejectBehaviorRollback,
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
	}
}

// createTestTemplateWithRejectTerminate 创建包含拒绝后终止配置的测试模板
func createTestTemplateWithRejectTerminate() *template.Template {
	return &template.Template{
		ID:          "tpl-001",
		Name:        "Test Template",
		Description: "Test template description",
		Version:     1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
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
					Mode:           node.ApprovalModeSingle,
					RejectBehavior: node.RejectBehaviorTerminate,
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
	}
}
