package task_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
)

// TestRequireAttachmentsApprove 测试附件要求 - 同意操作
func TestRequireAttachmentsApprove(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateWithRequireAttachments()
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

	// 尝试无附件的同意操作(应该失败)
	err = taskMgr.ApproveWithAttachments(tsk.ID, "approval-node", "user-1", "approved", []string{})
	if err == nil {
		t.Error("ApproveWithAttachments() should fail when attachments are required but empty")
	}

	// 尝试有附件的同意操作(应该成功)
	err = taskMgr.ApproveWithAttachments(tsk.ID, "approval-node", "user-1", "approved", []string{"file1.pdf", "file2.jpg"})
	if err != nil {
		t.Fatalf("ApproveWithAttachments() with attachments should succeed, got error: %v", err)
	}
}

// TestRequireAttachmentsReject 测试附件要求 - 拒绝操作
func TestRequireAttachmentsReject(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateWithRequireAttachments()
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

	// 尝试无附件的拒绝操作(应该失败)
	err = taskMgr.RejectWithAttachments(tsk.ID, "approval-node", "user-1", "rejected", []string{})
	if err == nil {
		t.Error("RejectWithAttachments() should fail when attachments are required but empty")
	}

	// 尝试有附件的拒绝操作(应该成功)
	err = taskMgr.RejectWithAttachments(tsk.ID, "approval-node", "user-1", "rejected", []string{"file1.pdf"})
	if err != nil {
		t.Fatalf("RejectWithAttachments() with attachments should succeed, got error: %v", err)
	}
}

// TestNoRequireAttachments 测试不要求附件的情况
func TestNoRequireAttachments(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateWithoutRequireAttachments()
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

	// 尝试无附件的同意操作(应该成功,因为不要求附件)
	err = taskMgr.ApproveWithAttachments(tsk.ID, "approval-node", "user-1", "approved", []string{})
	if err != nil {
		t.Fatalf("ApproveWithAttachments() without attachments should succeed when not required, got error: %v", err)
	}
}

// createTestTemplateWithRequireAttachments 创建要求附件的测试模板
func createTestTemplateWithRequireAttachments() *template.Template {
	return &template.Template{
		ID:          "tpl-001",
		Name:        "Test Template With Require Attachments",
		Description: "Test template with require attachments",
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
					RequireAttachmentsField: true,
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

// createTestTemplateWithoutRequireAttachments 创建不要求附件的测试模板
func createTestTemplateWithoutRequireAttachments() *template.Template {
	return &template.Template{
		ID:          "tpl-001",
		Name:        "Test Template Without Require Attachments",
		Description: "Test template without require attachments",
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
					RequireAttachmentsField: false,
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

