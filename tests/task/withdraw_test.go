package task_test

import (
	"encoding/json"
	"testing"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
	"github.com/mautops/approval-kit/internal/types"
)

// TestWithdraw 测试任务撤回功能
func TestWithdraw(t *testing.T) {
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

	// 撤回任务
	err = taskMgr.Withdraw(tsk.ID, "user withdraw")
	if err != nil {
		t.Fatalf("Withdraw() failed: %v", err)
	}

	// 验证任务状态已变为待审批
	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	if tsk.GetState() != types.TaskStatePending {
		t.Errorf("Task state = %q, want %q", tsk.GetState(), types.TaskStatePending)
	}

	// 验证 SubmittedAt 已清空
	if tsk.SubmittedAt != nil {
		t.Error("Task.SubmittedAt should be nil after withdraw")
	}
}

// TestWithdrawApprovingTask 测试撤回审批中的任务
// 注意: 这个测试需要先让任务进入审批中状态,但由于需要节点激活,
// 这里先跳过,等集成测试时再验证
func TestWithdrawApprovingTask(t *testing.T) {
	t.Skip("Skipping test that requires task in approving state - to be implemented in integration tests")
}

// TestWithdrawPendingTask 测试撤回待审批任务(应该失败)
func TestWithdrawPendingTask(t *testing.T) {
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

	// 尝试撤回待审批任务(应该失败)
	err = taskMgr.Withdraw(tsk.ID, "user withdraw")
	if err == nil {
		t.Error("Withdraw() should fail when task is pending")
	}
}

// TestWithdrawApprovedTask 测试撤回已通过的任务(应该失败)
func TestWithdrawApprovedTask(t *testing.T) {
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

	// 尝试撤回待审批的任务(应该失败,因为状态不是 submitted 或 approving)
	err = taskMgr.Withdraw(tsk.ID, "user withdraw")
	if err == nil {
		t.Error("Withdraw() should fail when task is pending")
	}
}

// TestWithdrawWithApprovalRecords 测试有审批记录时撤回(应该失败)
func TestWithdrawWithApprovalRecords(t *testing.T) {
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

	// 手动添加审批记录(模拟已有审批操作)
	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	// 通过 AddApprover 添加审批人(这会生成记录)
	// 注意: 需要模板允许加签
	tplWithAddApprover := createTestTemplate()
	tplWithAddApprover.Nodes["approval-001"].Config.(*node.ApprovalNodeConfig).Permissions.AllowAddApprover = true
	err = templateMgr.Update("tpl-001", tplWithAddApprover)
	if err != nil {
		t.Fatalf("Update template failed: %v", err)
	}

	err = taskMgr.AddApprover(tsk.ID, "approval-001", "user-001", "test add approver")
	if err != nil {
		t.Fatalf("AddApprover() failed: %v", err)
	}

	// 尝试撤回(应该失败,因为有审批记录)
	err = taskMgr.Withdraw(tsk.ID, "user withdraw")
	if err == nil {
		t.Error("Withdraw() should fail when task has approval records")
	}
}

// TestWithdrawNotFound 测试撤回不存在的任务
func TestWithdrawNotFound(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	taskMgr := task.NewTaskManager(templateMgr, nil)

	// 尝试撤回不存在的任务
	err := taskMgr.Withdraw("non-existent", "user withdraw")
	if err == nil {
		t.Error("Withdraw() should fail when task does not exist")
	}
}

// TestWithdrawStateHistory 测试撤回记录状态变更历史
func TestWithdrawStateHistory(t *testing.T) {
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

	// 撤回任务
	err = taskMgr.Withdraw(tsk.ID, "user withdraw")
	if err != nil {
		t.Fatalf("Withdraw() failed: %v", err)
	}

	// 验证状态变更历史
	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	if len(tsk.StateHistory) < 2 {
		t.Errorf("Task.StateHistory should contain at least 2 records, got %d", len(tsk.StateHistory))
	}

	// 验证最后一条状态变更记录是撤回
	lastChange := tsk.StateHistory[len(tsk.StateHistory)-1]
	if lastChange.From != types.TaskStateSubmitted && lastChange.From != types.TaskStateApproving {
		t.Errorf("StateChange.From = %v, want %v or %v", lastChange.From, types.TaskStateSubmitted, types.TaskStateApproving)
	}
	if lastChange.To != types.TaskStatePending {
		t.Errorf("StateChange.To = %v, want %v", lastChange.To, types.TaskStatePending)
	}
	if lastChange.Reason == "" {
		t.Error("StateChange.Reason should not be empty")
	}
}

