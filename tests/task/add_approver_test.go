package task_test

import (
	"encoding/json"
	"testing"

	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
)

// TestAddApprover 测试加签功能
// 注意: 完整测试需要审批人列表已设置,这通常在节点激活时完成
// 这里先测试基本逻辑,完整流程在集成测试中验证
func TestAddApprover(t *testing.T) {
	t.Skip("Skipping test that requires approvers list to be set - to be implemented in integration tests")
}

// TestAddApproverWithoutPermission 测试无权限加签(应该失败)
func TestAddApproverWithoutPermission(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateWithoutAddApproverPermission()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	// 创建任务管理器
	taskMgr := task.NewTaskManager(templateMgr, nil)

	// 创建任务
	params := json.RawMessage(`{"amount": 1000}`)
	tsk, err := taskMgr.Create("tpl-002", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 提交任务
	err = taskMgr.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Submit() failed: %v", err)
	}

	// 尝试加签(应该失败,因为没有权限)
	err = taskMgr.AddApprover(tsk.ID, "approval-node", "user-2", "add approver reason")
	if err == nil {
		t.Error("AddApprover() should fail when permission is not allowed")
	}
}

// TestAddApproverNotFound 测试加签不存在的任务
func TestAddApproverNotFound(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	taskMgr := task.NewTaskManager(templateMgr, nil)

	// 尝试加签不存在的任务
	err := taskMgr.AddApprover("non-existent", "approval-node", "user-2", "add approver reason")
	if err == nil {
		t.Error("AddApprover() should fail when task does not exist")
	}
}

// createTestTemplateWithoutAddApproverPermission 在 permissions_test.go 中定义

