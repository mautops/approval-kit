package task_test

import (
	"encoding/json"
	"testing"

	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
	"github.com/mautops/approval-kit/internal/types"
)

// TestRollbackToNode 测试回退到指定节点功能
func TestRollbackToNode(t *testing.T) {
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

	// 执行审批操作,使任务完成
	err = taskMgr.Approve(tsk.ID, "approval-001", "user-001", "approved")
	if err != nil {
		t.Fatalf("Approve() failed: %v", err)
	}

	// 验证任务已完成
	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}
	if tsk.GetState() != types.TaskStateApproved {
		t.Fatalf("Task state before rollback = %q, want %q", tsk.GetState(), types.TaskStateApproved)
	}

	// 回退到审批节点
	err = taskMgr.RollbackToNode(tsk.ID, "approval-001", "user rollback")
	if err != nil {
		t.Fatalf("RollbackToNode() failed: %v", err)
	}

	// 验证任务状态已回退
	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	// 回退后,任务应该回到审批中状态
	if tsk.GetState() != types.TaskStateApproving {
		t.Errorf("Task state = %q, want %q", tsk.GetState(), types.TaskStateApproving)
	}
	if tsk.CurrentNode != "approval-001" {
		t.Errorf("Task.CurrentNode = %q, want %q", tsk.CurrentNode, "approval-001")
	}
}

// TestRollbackToNodeNotFound 测试回退到不存在的节点
func TestRollbackToNodeNotFound(t *testing.T) {
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

	// 尝试回退到不存在的节点(应该失败)
	err = taskMgr.RollbackToNode(tsk.ID, "non-existent", "user rollback")
	if err == nil {
		t.Error("RollbackToNode() should fail for non-existent node")
	}
}

// TestRollbackToNodeNotCompleted 测试回退到未完成的节点
func TestRollbackToNodeNotCompleted(t *testing.T) {
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

	// 尝试回退到未完成的节点(应该失败)
	err = taskMgr.RollbackToNode(tsk.ID, "approval-001", "user rollback")
	if err == nil {
		t.Error("RollbackToNode() should fail for non-completed node")
	}
}

// TestRollbackToNodeNotFound 测试回退不存在的任务
func TestRollbackToNodeTaskNotFound(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	taskMgr := task.NewTaskManager(templateMgr, nil)

	err := taskMgr.RollbackToNode("non-existent", "approval-001", "user rollback")
	if err == nil {
		t.Error("RollbackToNode() should fail for non-existent task")
	}
}

// TestRollbackToNodeStateHistory 测试回退时记录状态变更历史
func TestRollbackToNodeStateHistory(t *testing.T) {
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

	// 执行审批操作
	err = taskMgr.Approve(tsk.ID, "approval-001", "user-001", "approved")
	if err != nil {
		t.Fatalf("Approve() failed: %v", err)
	}

	// 回退到审批节点
	err = taskMgr.RollbackToNode(tsk.ID, "approval-001", "user rollback")
	if err != nil {
		t.Fatalf("RollbackToNode() failed: %v", err)
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

	// 检查回退的状态变更记录
	lastChange := history[len(history)-1]
	if lastChange.To != types.TaskStateApproving {
		t.Errorf("Last state change To = %q, want %q", lastChange.To, types.TaskStateApproving)
	}
	if lastChange.Reason != "user rollback" {
		t.Errorf("Last state change Reason = %q, want %q", lastChange.Reason, "user rollback")
	}
}

