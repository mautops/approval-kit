package task_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
)

// TestRequireCommentApprove 测试审批意见必填 - 同意操作
func TestRequireCommentApprove(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateWithRequireComment()
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

	// 尝试无审批意见的同意操作(应该失败)
	err = taskMgr.Approve(tsk.ID, "approval-node", "user-1", "")
	if err == nil {
		t.Error("Approve() should fail when comment is required but empty")
	}

	// 尝试有审批意见的同意操作(应该成功)
	err = taskMgr.Approve(tsk.ID, "approval-node", "user-1", "approved")
	if err != nil {
		t.Fatalf("Approve() with comment should succeed, got error: %v", err)
	}
}

// TestRequireCommentReject 测试审批意见必填 - 拒绝操作
func TestRequireCommentReject(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateWithRequireComment()
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

	// 尝试无审批意见的拒绝操作(应该失败)
	err = taskMgr.Reject(tsk.ID, "approval-node", "user-1", "")
	if err == nil {
		t.Error("Reject() should fail when comment is required but empty")
	}

	// 尝试有审批意见的拒绝操作(应该成功)
	err = taskMgr.Reject(tsk.ID, "approval-node", "user-1", "rejected")
	if err != nil {
		t.Fatalf("Reject() with comment should succeed, got error: %v", err)
	}
}

// TestNoRequireComment 测试不要求审批意见的情况
func TestNoRequireComment(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateWithoutRequireComment()
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

	// 尝试无审批意见的同意操作(应该成功,因为不要求审批意见)
	err = taskMgr.Approve(tsk.ID, "approval-node", "user-1", "")
	if err != nil {
		t.Fatalf("Approve() without comment should succeed when not required, got error: %v", err)
	}
}

// createTestTemplateWithRequireComment 创建要求审批意见的测试模板
func createTestTemplateWithRequireComment() *template.Template {
	return &template.Template{
		ID:          "tpl-001",
		Name:        "Test Template With Require Comment",
		Description: "Test template with require comment",
		Version:     1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "Start Node",
				Type: template.NodeTypeStart,
			},
			"approval-node": {
				ID:   "approval-node",
				Name: "Approval Node",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeSingle,
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{"user-1"},
					},
					RequireCommentField: true,
				},
			},
			"end": {
				ID:   "end",
				Name: "End Node",
				Type: template.NodeTypeEnd,
			},
		},
		Edges: []*template.Edge{
			{From: "start", To: "approval-node"},
			{From: "approval-node", To: "end"},
		},
	}
}

// createTestTemplateWithoutRequireComment 创建不要求审批意见的测试模板
func createTestTemplateWithoutRequireComment() *template.Template {
	return &template.Template{
		ID:          "tpl-001",
		Name:        "Test Template Without Require Comment",
		Description: "Test template without require comment",
		Version:     1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "Start Node",
				Type: template.NodeTypeStart,
			},
			"approval-node": {
				ID:   "approval-node",
				Name: "Approval Node",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeSingle,
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{"user-1"},
					},
					RequireCommentField: false,
				},
			},
			"end": {
				ID:   "end",
				Name: "End Node",
				Type: template.NodeTypeEnd,
			},
		},
		Edges: []*template.Edge{
			{From: "start", To: "approval-node"},
			{From: "approval-node", To: "end"},
		},
	}
}

