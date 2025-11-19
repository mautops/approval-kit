package task_test

import (
	stderrors "errors"
	"testing"

	"github.com/mautops/approval-kit/internal/errors"
	"github.com/mautops/approval-kit/internal/task"
)

// TestStateTransitionError 验证 StateTransitionError 错误类型
func TestStateTransitionError(t *testing.T) {
	from := task.TaskStatePending
	to := task.TaskStateApproved
	baseErr := errors.ErrInvalidStateTransition

	err := &task.StateTransitionError{
		From: from,
		To:   to,
		Err:  baseErr,
	}

	// 验证错误消息格式
	expectedMsg := "state transition failed: pending -> approved: invalid state transition"
	if err.Error() != expectedMsg {
		t.Errorf("StateTransitionError.Error() = %q, want %q", err.Error(), expectedMsg)
	}

	// 验证 Unwrap 方法
	if stderrors.Unwrap(err) != baseErr {
		t.Errorf("StateTransitionError.Unwrap() = %v, want %v", stderrors.Unwrap(err), baseErr)
	}

	// 验证 Is 方法
	if !stderrors.Is(err, baseErr) {
		t.Errorf("errors.Is() should return true for wrapped error")
	}
}

// TestStateTransitionErrorAllStates 验证所有状态组合的错误消息
func TestStateTransitionErrorAllStates(t *testing.T) {
	states := []task.TaskState{
		task.TaskStatePending,
		task.TaskStateSubmitted,
		task.TaskStateApproving,
		task.TaskStateApproved,
		task.TaskStateRejected,
		task.TaskStateCancelled,
		task.TaskStateTimeout,
	}

	for _, from := range states {
		for _, to := range states {
			if from == to {
				continue
			}
			err := &task.StateTransitionError{
				From: from,
				To:   to,
				Err:  errors.ErrInvalidStateTransition,
			}

			// 验证错误消息包含状态信息
			msg := err.Error()
			if msg == "" {
				t.Errorf("StateTransitionError.Error() should not be empty")
			}
			// 验证错误可以被展开
			if !stderrors.Is(err, errors.ErrInvalidStateTransition) {
				t.Errorf("StateTransitionError should be unwrappable to ErrInvalidStateTransition")
			}
		}
	}
}

