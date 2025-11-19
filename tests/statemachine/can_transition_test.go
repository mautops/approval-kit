package statemachine_test

import (
	"testing"

	"github.com/mautops/approval-kit/internal/statemachine"
	"github.com/mautops/approval-kit/internal/task"
)

// TestCanTransition 测试 CanTransition 方法
// 使用表驱动测试,覆盖所有状态转换场景
func TestCanTransition(t *testing.T) {
	sm := statemachine.NewStateMachine()

	tests := []struct {
		name    string
		from    task.TaskState
		to      task.TaskState
		want    bool
		wantErr bool
	}{
		// 合法转换
		{
			name:    "pending to submitted",
			from:    task.TaskStatePending,
			to:      task.TaskStateSubmitted,
			want:    true,
			wantErr: false,
		},
		{
			name:    "pending to cancelled",
			from:    task.TaskStatePending,
			to:      task.TaskStateCancelled,
			want:    true,
			wantErr: false,
		},
		{
			name:    "submitted to approving",
			from:    task.TaskStateSubmitted,
			to:      task.TaskStateApproving,
			want:    true,
			wantErr: false,
		},
		{
			name:    "submitted to cancelled",
			from:    task.TaskStateSubmitted,
			to:      task.TaskStateCancelled,
			want:    true,
			wantErr: false,
		},
		{
			name:    "approving to approved",
			from:    task.TaskStateApproving,
			to:      task.TaskStateApproved,
			want:    true,
			wantErr: false,
		},
		{
			name:    "approving to rejected",
			from:    task.TaskStateApproving,
			to:      task.TaskStateRejected,
			want:    true,
			wantErr: false,
		},
		{
			name:    "approving to cancelled",
			from:    task.TaskStateApproving,
			to:      task.TaskStateCancelled,
			want:    true,
			wantErr: false,
		},
		{
			name:    "approving to timeout",
			from:    task.TaskStateApproving,
			to:      task.TaskStateTimeout,
			want:    true,
			wantErr: false,
		},

		// 非法转换: 终态不能转换
		{
			name:    "approved to pending (final state)",
			from:    task.TaskStateApproved,
			to:      task.TaskStatePending,
			want:    false,
			wantErr: false,
		},
		{
			name:    "rejected to pending (final state)",
			from:    task.TaskStateRejected,
			to:      task.TaskStatePending,
			want:    false,
			wantErr: false,
		},
		{
			name:    "cancelled to pending (final state)",
			from:    task.TaskStateCancelled,
			to:      task.TaskStatePending,
			want:    false,
			wantErr: false,
		},
		{
			name:    "timeout to pending (final state)",
			from:    task.TaskStateTimeout,
			to:      task.TaskStatePending,
			want:    false,
			wantErr: false,
		},

		// 非法转换: 不允许的转换路径
		{
			name:    "pending to approving (invalid path)",
			from:    task.TaskStatePending,
			to:      task.TaskStateApproving,
			want:    false,
			wantErr: false,
		},
		{
			name:    "pending to approved (invalid path)",
			from:    task.TaskStatePending,
			to:      task.TaskStateApproved,
			want:    false,
			wantErr: false,
		},
		{
			name:    "submitted to approved (invalid path)",
			from:    task.TaskStateSubmitted,
			to:      task.TaskStateApproved,
			want:    false,
			wantErr: false,
		},
		{
			name:    "approving to submitted (invalid path)",
			from:    task.TaskStateApproving,
			to:      task.TaskStateSubmitted,
			want:    false,
			wantErr: false,
		},

		// 边界情况: 相同状态
		{
			name:    "pending to pending (same state)",
			from:    task.TaskStatePending,
			to:      task.TaskStatePending,
			want:    false,
			wantErr: false,
		},
		{
			name:    "approved to approved (same state)",
			from:    task.TaskStateApproved,
			to:      task.TaskStateApproved,
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sm.CanTransition(tt.from, tt.to)
			if got != tt.want {
				t.Errorf("CanTransition(%q, %q) = %v, want %v", tt.from, tt.to, got, tt.want)
			}
		})
	}
}

// TestCanTransitionEdgeCases 测试边界情况
func TestCanTransitionEdgeCases(t *testing.T) {
	sm := statemachine.NewStateMachine()

	// 测试所有状态组合
	allStates := []task.TaskState{
		task.TaskStatePending,
		task.TaskStateSubmitted,
		task.TaskStateApproving,
		task.TaskStateApproved,
		task.TaskStateRejected,
		task.TaskStateCancelled,
		task.TaskStateTimeout,
	}

	// 验证每个状态都有明确的转换规则
	for _, from := range allStates {
		for _, to := range allStates {
			result := sm.CanTransition(from, to)
			// 结果应该是确定的(不应该是随机或未定义的)
			_ = result // 确保方法不会 panic
		}
	}
}

// TestCanTransitionWithdraw 测试撤回功能的状态转换
func TestCanTransitionWithdraw(t *testing.T) {
	sm := statemachine.NewStateMachine()

	// 测试撤回: submitted -> pending
	if !sm.CanTransition(task.TaskStateSubmitted, task.TaskStatePending) {
		t.Error("CanTransition(submitted, pending) should return true for withdraw")
	}

	// 测试撤回: approving -> pending
	if !sm.CanTransition(task.TaskStateApproving, task.TaskStatePending) {
		t.Error("CanTransition(approving, pending) should return true for withdraw")
	}
}

// TestCanTransitionUndefinedState 测试未定义状态的处理
func TestCanTransitionUndefinedState(t *testing.T) {
	sm := statemachine.NewStateMachine()

	// 测试从未定义的状态转换(应该返回 false)
	undefinedState := task.TaskState("undefined")
	if sm.CanTransition(undefinedState, task.TaskStatePending) {
		t.Error("CanTransition(undefined, pending) should return false")
	}
}
