package task_test

import (
	"encoding/json"
	"testing"

	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
	"github.com/mautops/approval-kit/internal/types"
)

// TestCancel 测试任务取消功能
func TestCancel(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplate()
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

	// 取消任务
	err = taskMgr.Cancel(tsk.ID, "user cancelled")
	if err != nil {
		t.Fatalf("Cancel() failed: %v", err)
	}

	// 验证任务状态已变为已取消
	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	if tsk.GetState() != types.TaskStateCancelled {
		t.Errorf("Task state = %q, want %q", tsk.GetState(), types.TaskStateCancelled)
	}
}

// TestCancelSubmittedTask 测试取消已提交的任务
func TestCancelSubmittedTask(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplate()
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

	// 取消任务
	err = taskMgr.Cancel(tsk.ID, "user cancelled")
	if err != nil {
		t.Fatalf("Cancel() failed: %v", err)
	}

	// 验证任务状态已变为已取消
	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	if tsk.GetState() != types.TaskStateCancelled {
		t.Errorf("Task state = %q, want %q", tsk.GetState(), types.TaskStateCancelled)
	}
}

// TestCancelApprovedTask 测试取消已通过的任务(应该失败)
func TestCancelApprovedTask(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateWithApprovalNode()
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

	// 执行审批操作
	err = taskMgr.Approve(tsk.ID, "approval-001", "user-001", "approved")
	if err != nil {
		t.Fatalf("Approve() failed: %v", err)
	}

	// 尝试取消已通过的任务(应该失败)
	err = taskMgr.Cancel(tsk.ID, "user cancelled")
	if err == nil {
		t.Error("Cancel() should fail for approved task")
	}
}

