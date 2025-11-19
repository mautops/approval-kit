package task_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
)

// TestTransfer 测试转交审批功能
func TestTransfer(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateWithTransferPermission()
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

	// 手动设置审批人列表(因为固定审批人需要在节点激活时设置)
	err = taskMgr.AddApprover(tsk.ID, "approval-node", "user-1", "setup approver")
	if err != nil {
		t.Fatalf("AddApprover() failed: %v", err)
	}

	// 转交审批
	err = taskMgr.Transfer(tsk.ID, "approval-node", "user-1", "user-2", "transfer reason")
	if err != nil {
		t.Fatalf("Transfer() failed: %v", err)
	}

	// 验证审批人列表已更新
	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	approvers := tsk.Approvers["approval-node"]
	foundUser2 := false
	foundUser1 := false
	for _, approver := range approvers {
		if approver == "user-2" {
			foundUser2 = true
		}
		if approver == "user-1" {
			foundUser1 = true
		}
	}

	if !foundUser2 {
		t.Error("Transfer() should add user-2 to approvers list")
	}
	if foundUser1 {
		t.Error("Transfer() should remove user-1 from approvers list")
	}

	// 验证转交记录已生成
	if len(tsk.Records) == 0 {
		t.Fatal("Transfer() should generate a record")
	}
	lastRecord := tsk.Records[len(tsk.Records)-1]
	if lastRecord.Result != "transfer" {
		t.Errorf("Record.Result = %q, want %q", lastRecord.Result, "transfer")
	}
}

// TestTransferWithoutPermission 测试无权限转交(应该失败)
func TestTransferWithoutPermission(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateWithoutTransferPermission()
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

	// 尝试转交审批(应该失败,因为没有权限)
	err = taskMgr.Transfer(tsk.ID, "approval-node", "user-1", "user-2", "transfer reason")
	if err == nil {
		t.Error("Transfer() should fail when permission is not allowed")
	}
}

// TestTransferNotApprover 测试非审批人转交(应该失败)
func TestTransferNotApprover(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateWithTransferPermission()
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

	// 尝试非审批人转交(应该失败)
	err = taskMgr.Transfer(tsk.ID, "approval-node", "user-3", "user-2", "transfer reason")
	if err == nil {
		t.Error("Transfer() should fail when user is not an approver")
	}
}

// TestTransferNotFound 测试转交不存在的任务
func TestTransferNotFound(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	taskMgr := task.NewTaskManager(templateMgr, nil)

	// 尝试转交不存在的任务
	err := taskMgr.Transfer("non-existent", "approval-node", "user-1", "user-2", "transfer reason")
	if err == nil {
		t.Error("Transfer() should fail when task does not exist")
	}
}

// TestTransferToExistingApprover 测试转交给已在列表中的审批人
func TestTransferToExistingApprover(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateWithTransferPermission()
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

	// 添加多个审批人
	err = taskMgr.AddApprover(tsk.ID, "approval-node", "user-1", "setup approver")
	if err != nil {
		t.Fatalf("AddApprover() failed: %v", err)
	}

	err = taskMgr.AddApprover(tsk.ID, "approval-node", "user-2", "setup approver")
	if err != nil {
		t.Fatalf("AddApprover() failed: %v", err)
	}

	// 转交给已在列表中的审批人(user-2)
	err = taskMgr.Transfer(tsk.ID, "approval-node", "user-1", "user-2", "transfer reason")
	if err != nil {
		t.Fatalf("Transfer() failed: %v", err)
	}

	// 验证审批人列表: user-1 应被移除, user-2 应保留(不重复)
	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	approvers := tsk.Approvers["approval-node"]
	user1Count := 0
	user2Count := 0
	for _, approver := range approvers {
		if approver == "user-1" {
			user1Count++
		}
		if approver == "user-2" {
			user2Count++
		}
	}

	if user1Count > 0 {
		t.Error("Transfer() should remove user-1 from approvers list")
	}
	if user2Count != 1 {
		t.Errorf("Transfer() should keep user-2 in approvers list exactly once, got %d", user2Count)
	}
}

// TestTransferApproversNotFound 测试转交时审批人列表不存在
func TestTransferApproversNotFound(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateWithTransferPermission()
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

	// 尝试转交(审批人列表为空)
	err = taskMgr.Transfer(tsk.ID, "approval-node", "user-1", "user-2", "transfer reason")
	if err == nil {
		t.Error("Transfer() should fail when approvers list is empty")
	}
}

// createTestTemplateWithTransferPermission 创建允许转交的测试模板
func createTestTemplateWithTransferPermission() *template.Template {
	return &template.Template{
		ID:          "tpl-001",
		Name:        "Test Template With Transfer",
		Description: "Test template with transfer permission",
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
					Permissions: node.OperationPermissions{
						AllowTransfer:      true,
						AllowAddApprover:   true, // 需要允许加签才能设置审批人列表
						AllowRemoveApprover: true,
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
			{From: "start", To: "approval-node"},
			{From: "approval-node", To: "end"},
		},
	}
}

// createTestTemplateWithoutTransferPermission 在 permissions_test.go 中定义

