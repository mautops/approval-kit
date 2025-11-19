package task_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
)

// TestTaskManagerSubmit 测试任务提交功能
func TestTaskManagerSubmit(t *testing.T) {
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

	// 验证初始状态
	if tsk.State != task.TaskStatePending {
		t.Errorf("Initial task state = %v, want %v", tsk.State, task.TaskStatePending)
	}

	// 提交任务
	err = taskMgr.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Submit() failed: %v", err)
	}

	// 验证状态已转换
	submittedTask, err := taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}
	if submittedTask.State != task.TaskStateSubmitted {
		t.Errorf("Task state after Submit = %v, want %v", submittedTask.State, task.TaskStateSubmitted)
	}
	if submittedTask.SubmittedAt == nil {
		t.Error("Task.SubmittedAt should be set after Submit")
	}
}

// TestTaskManagerSubmitStateTransition 测试状态转换
func TestTaskManagerSubmitStateTransition(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplate()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	taskMgr := task.NewTaskManager(templateMgr, nil)

	tsk, err := taskMgr.Create("tpl-001", "biz-001", nil)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 提交任务
	err = taskMgr.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Submit() failed: %v", err)
	}

	// 验证状态变更历史
	submittedTask, err := taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	if len(submittedTask.StateHistory) == 0 {
		t.Error("Task.StateHistory should contain state change record")
	}

	// 验证状态变更记录
	lastChange := submittedTask.StateHistory[len(submittedTask.StateHistory)-1]
	if lastChange.From != task.TaskStatePending {
		t.Errorf("StateChange.From = %v, want %v", lastChange.From, task.TaskStatePending)
	}
	if lastChange.To != task.TaskStateSubmitted {
		t.Errorf("StateChange.To = %v, want %v", lastChange.To, task.TaskStateSubmitted)
	}
	if lastChange.Reason == "" {
		t.Error("StateChange.Reason should not be empty")
	}
	if lastChange.Time.IsZero() {
		t.Error("StateChange.Time should be set")
	}
}

// TestTaskManagerSubmitInvalidState 测试无效状态提交
func TestTaskManagerSubmitInvalidState(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplate()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	taskMgr := task.NewTaskManager(templateMgr, nil)

	tsk, err := taskMgr.Create("tpl-001", "biz-001", nil)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 第一次提交应该成功
	err = taskMgr.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Submit() failed: %v", err)
	}

	// 第二次提交应该失败(状态已经是 submitted)
	err = taskMgr.Submit(tsk.ID)
	if err == nil {
		t.Error("Submit() should fail when task is already submitted")
	}
}

// TestTaskManagerSubmitNotFound 测试提交不存在的任务
func TestTaskManagerSubmitNotFound(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	taskMgr := task.NewTaskManager(templateMgr, nil)

	// 尝试提交不存在的任务
	err := taskMgr.Submit("non-existent")
	if err == nil {
		t.Error("Submit() should fail when task does not exist")
	}
}

// TestTaskManagerSubmitUpdatesTimestamp 测试提交更新时间戳
func TestTaskManagerSubmitUpdatesTimestamp(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplate()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	taskMgr := task.NewTaskManager(templateMgr, nil)

	tsk, err := taskMgr.Create("tpl-001", "biz-001", nil)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 等待一小段时间
	time.Sleep(10 * time.Millisecond)

	originalUpdatedAt := tsk.UpdatedAt

	// 提交任务
	err = taskMgr.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Submit() failed: %v", err)
	}

	// 验证更新时间已更新
	submittedTask, err := taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}
	if !submittedTask.UpdatedAt.After(originalUpdatedAt) {
		t.Error("Task.UpdatedAt should be updated after Submit")
	}
}

