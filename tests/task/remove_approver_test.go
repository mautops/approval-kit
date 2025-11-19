package task_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
)

// TestRemoveApprover 测试减签功能
func TestRemoveApprover(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateWithRemoveApproverPermission()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	// 创建任务管理器
	taskMgr := task.NewTaskManager(templateMgr, nil)

	// 创建任务
	params := json.RawMessage(`{"amount": 1000}`)
	tsk, err := taskMgr.Create("tpl-003", "biz-001", params)
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

	err = taskMgr.AddApprover(tsk.ID, "approval-node", "user-2", "setup approver")
	if err != nil {
		t.Fatalf("AddApprover() failed: %v", err)
	}

	// 减签
	err = taskMgr.RemoveApprover(tsk.ID, "approval-node", "user-1", "remove approver reason")
	if err != nil {
		t.Fatalf("RemoveApprover() failed: %v", err)
	}

	// 验证审批人列表已更新
	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	approvers := tsk.Approvers["approval-node"]
	foundUser1 := false
	foundUser2 := false
	for _, approver := range approvers {
		if approver == "user-1" {
			foundUser1 = true
		}
		if approver == "user-2" {
			foundUser2 = true
		}
	}

	if foundUser1 {
		t.Error("RemoveApprover() should remove user-1 from approvers list")
	}
	if !foundUser2 {
		t.Error("RemoveApprover() should keep user-2 in approvers list")
	}

	// 验证减签记录已生成
	if len(tsk.Records) < 2 {
		t.Fatal("RemoveApprover() should generate a record")
	}
	lastRecord := tsk.Records[len(tsk.Records)-1]
	if lastRecord.Result != "remove_approver" {
		t.Errorf("Record.Result = %q, want %q", lastRecord.Result, "remove_approver")
	}
}

// createTestTemplateWithRemoveApproverPermission 创建允许减签的测试模板
func createTestTemplateWithRemoveApproverPermission() *template.Template {
	return &template.Template{
		ID:          "tpl-003",
		Name:        "Test Template With Remove Approver Permission",
		Description: "Test template with remove approver permission",
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
						AllowRemoveApprover: true,
						AllowAddApprover:    true, // 需要允许加签才能设置审批人列表
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

// TestRemoveApproverWithoutPermission 测试无权限减签(应该失败)
func TestRemoveApproverWithoutPermission(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateWithoutRemoveApproverPermission()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	// 创建任务管理器
	taskMgr := task.NewTaskManager(templateMgr, nil)

	// 创建任务
	params := json.RawMessage(`{"amount": 1000}`)
	tsk, err := taskMgr.Create("tpl-003", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 提交任务
	err = taskMgr.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Submit() failed: %v", err)
	}

	// 尝试减签(应该失败,因为没有权限)
	err = taskMgr.RemoveApprover(tsk.ID, "approval-node", "user-1", "remove approver reason")
	if err == nil {
		t.Error("RemoveApprover() should fail when permission is not allowed")
	}
}

// TestRemoveApproverNotFound 测试减签不存在的任务
func TestRemoveApproverNotFound(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	taskMgr := task.NewTaskManager(templateMgr, nil)

	// 尝试减签不存在的任务
	err := taskMgr.RemoveApprover("non-existent", "approval-node", "user-1", "remove approver reason")
	if err == nil {
		t.Error("RemoveApprover() should fail when task does not exist")
	}
}

// TestRemoveApproverApproversNotFound 测试减签时审批人列表不存在
func TestRemoveApproverApproversNotFound(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateWithRemoveApproverPermission()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	// 创建任务管理器
	taskMgr := task.NewTaskManager(templateMgr, nil)

	// 创建任务
	params := json.RawMessage(`{"amount": 1000}`)
	tsk, err := taskMgr.Create("tpl-003", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 提交任务
	err = taskMgr.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Submit() failed: %v", err)
	}

	// 尝试减签(审批人列表为空)
	err = taskMgr.RemoveApprover(tsk.ID, "approval-node", "user-1", "remove approver reason")
	if err == nil {
		t.Error("RemoveApprover() should fail when approvers list is empty")
	}
}

// TestRemoveApproverNotInList 测试减签不在列表中的审批人
func TestRemoveApproverNotInList(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateWithRemoveApproverPermission()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	// 创建任务管理器
	taskMgr := task.NewTaskManager(templateMgr, nil)

	// 创建任务
	params := json.RawMessage(`{"amount": 1000}`)
	tsk, err := taskMgr.Create("tpl-003", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 提交任务
	err = taskMgr.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Submit() failed: %v", err)
	}

	// 添加审批人
	err = taskMgr.AddApprover(tsk.ID, "approval-node", "user-1", "setup approver")
	if err != nil {
		t.Fatalf("AddApprover() failed: %v", err)
	}

	// 尝试减签不在列表中的审批人
	err = taskMgr.RemoveApprover(tsk.ID, "approval-node", "user-999", "remove approver reason")
	if err == nil {
		t.Error("RemoveApprover() should fail when approver is not in list")
	}
}

// createTestTemplateWithoutRemoveApproverPermission 在 permissions_test.go 中定义

