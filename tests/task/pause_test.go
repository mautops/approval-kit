package task_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
	"github.com/mautops/approval-kit/internal/types"
)

// TestPause 测试任务暂停功能
func TestPause(t *testing.T) {
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

	// 暂停任务
	err = taskMgr.Pause(tsk.ID, "user paused")
	if err != nil {
		t.Fatalf("Pause() failed: %v", err)
	}

	// 验证任务状态已变为已暂停
	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	if tsk.GetState() != types.TaskStatePaused {
		t.Errorf("Task state = %q, want %q", tsk.GetState(), types.TaskStatePaused)
	}

	// 验证暂停字段已设置
	if tsk.PausedAt == nil {
		t.Error("Task.PausedAt should not be nil")
	}
	if tsk.PausedState != types.TaskStatePending {
		t.Errorf("Task.PausedState = %q, want %q", tsk.PausedState, types.TaskStatePending)
	}
}

// TestPauseSubmittedTask 测试暂停已提交的任务
func TestPauseSubmittedTask(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplate()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	taskMgr := task.NewTaskManager(templateMgr, nil)

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

	// 暂停任务
	err = taskMgr.Pause(tsk.ID, "user paused")
	if err != nil {
		t.Fatalf("Pause() failed: %v", err)
	}

	// 验证任务状态
	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	if tsk.GetState() != types.TaskStatePaused {
		t.Errorf("Task state = %q, want %q", tsk.GetState(), types.TaskStatePaused)
	}
	if tsk.PausedState != types.TaskStateSubmitted {
		t.Errorf("Task.PausedState = %q, want %q", tsk.PausedState, types.TaskStateSubmitted)
	}
}

// TestPauseApprovingTask 测试暂停审批中的任务
func TestPauseApprovingTask(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplate()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	taskMgr := task.NewTaskManager(templateMgr, nil)

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

	// 暂停任务
	err = taskMgr.Pause(tsk.ID, "user paused")
	if err != nil {
		t.Fatalf("Pause() failed: %v", err)
	}

	// 验证任务状态
	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	if tsk.GetState() != types.TaskStatePaused {
		t.Errorf("Task state = %q, want %q", tsk.GetState(), types.TaskStatePaused)
	}
}

// TestPauseInvalidState 测试在无效状态下暂停任务
func TestPauseInvalidState(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplate()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	taskMgr := task.NewTaskManager(templateMgr, nil)

	params := json.RawMessage(`{"amount": 1000}`)
	tsk, err := taskMgr.Create("tpl-001", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 取消任务
	err = taskMgr.Cancel(tsk.ID, "cancelled")
	if err != nil {
		t.Fatalf("Cancel() failed: %v", err)
	}

	// 尝试暂停已取消的任务(应该失败)
	err = taskMgr.Pause(tsk.ID, "user paused")
	if err == nil {
		t.Error("Pause() should fail for cancelled task")
	}
}

// TestPauseNotFound 测试暂停不存在的任务
func TestPauseNotFound(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	taskMgr := task.NewTaskManager(templateMgr, nil)

	err := taskMgr.Pause("non-existent", "user paused")
	if err == nil {
		t.Error("Pause() should fail for non-existent task")
	}
}

// TestPauseStateHistory 测试暂停时记录状态变更历史
func TestPauseStateHistory(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplate()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	taskMgr := task.NewTaskManager(templateMgr, nil)

	params := json.RawMessage(`{"amount": 1000}`)
	tsk, err := taskMgr.Create("tpl-001", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 暂停任务
	err = taskMgr.Pause(tsk.ID, "user paused")
	if err != nil {
		t.Fatalf("Pause() failed: %v", err)
	}

	// 验证状态变更历史
	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	history := tsk.GetStateHistory()
	if len(history) < 1 {
		t.Fatal("State history should have at least 1 record")
	}

	lastChange := history[len(history)-1]
	if lastChange.From != types.TaskStatePending {
		t.Errorf("Last state change From = %q, want %q", lastChange.From, types.TaskStatePending)
	}
	if lastChange.To != types.TaskStatePaused {
		t.Errorf("Last state change To = %q, want %q", lastChange.To, types.TaskStatePaused)
	}
	if lastChange.Reason != "user paused" {
		t.Errorf("Last state change Reason = %q, want %q", lastChange.Reason, "user paused")
	}
}

// TestPauseTimestamp 测试暂停时间戳
func TestPauseTimestamp(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplate()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	taskMgr := task.NewTaskManager(templateMgr, nil)

	params := json.RawMessage(`{"amount": 1000}`)
	tsk, err := taskMgr.Create("tpl-001", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	beforePause := time.Now()

	// 暂停任务
	err = taskMgr.Pause(tsk.ID, "user paused")
	if err != nil {
		t.Fatalf("Pause() failed: %v", err)
	}

	afterPause := time.Now()

	// 验证暂停时间戳
	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	if tsk.PausedAt == nil {
		t.Fatal("Task.PausedAt should not be nil")
	}

	if tsk.PausedAt.Before(beforePause) || tsk.PausedAt.After(afterPause) {
		t.Errorf("Task.PausedAt = %v, should be between %v and %v", tsk.PausedAt, beforePause, afterPause)
	}
}

