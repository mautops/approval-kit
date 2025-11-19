package statemachine_test

import (
	"testing"

	"github.com/mautops/approval-kit/internal/statemachine"
	"github.com/mautops/approval-kit/internal/task"
)

// TestStateTransitionsDefined 验证状态转换规则已定义
func TestStateTransitionsDefined(t *testing.T) {
	// 验证所有状态都有转换规则定义(即使是空规则)
	allStates := []task.TaskState{
		task.TaskStatePending,
		task.TaskStateSubmitted,
		task.TaskStateApproving,
		task.TaskStateApproved,
		task.TaskStateRejected,
		task.TaskStateCancelled,
		task.TaskStateTimeout,
	}

	transitions := statemachine.GetStateTransitions()

	for _, state := range allStates {
		_, exists := transitions[state]
		if !exists {
			t.Errorf("状态 %s 没有定义转换规则", state)
		}
	}
}

// TestStateTransitionsRules 验证状态转换规则的正确性
func TestStateTransitionsRules(t *testing.T) {
	transitions := statemachine.GetStateTransitions()

	tests := []struct {
		name      string
		from      task.TaskState
		validTos  []task.TaskState
		invalidTos []task.TaskState
	}{
		{
			name: "pending state transitions",
			from: task.TaskStatePending,
			validTos: []task.TaskState{
				task.TaskStateSubmitted,
				task.TaskStateCancelled,
			},
			invalidTos: []task.TaskState{
				task.TaskStateApproving,
				task.TaskStateApproved,
				task.TaskStateRejected,
				task.TaskStateTimeout,
			},
		},
		{
			name: "submitted state transitions",
			from: task.TaskStateSubmitted,
			validTos: []task.TaskState{
				task.TaskStateApproving,
				task.TaskStateCancelled,
				task.TaskStatePending, // 撤回功能允许
			},
			invalidTos: []task.TaskState{
				task.TaskStateApproved,
				task.TaskStateRejected,
				task.TaskStateTimeout,
			},
		},
		{
			name: "approving state transitions",
			from: task.TaskStateApproving,
			validTos: []task.TaskState{
				task.TaskStateApproved,
				task.TaskStateRejected,
				task.TaskStateCancelled,
				task.TaskStateTimeout,
				task.TaskStatePending, // 撤回功能允许
			},
			invalidTos: []task.TaskState{
				task.TaskStateSubmitted,
			},
		},
		{
			name: "approved state (final state)",
			from: task.TaskStateApproved,
			validTos: []task.TaskState{}, // 终态,不能转换
			invalidTos: []task.TaskState{
				task.TaskStatePending,
				task.TaskStateSubmitted,
				task.TaskStateApproving,
				task.TaskStateRejected,
				task.TaskStateCancelled,
				task.TaskStateTimeout,
			},
		},
		{
			name: "rejected state (final state)",
			from: task.TaskStateRejected,
			validTos: []task.TaskState{}, // 终态,不能转换
			invalidTos: []task.TaskState{
				task.TaskStatePending,
				task.TaskStateSubmitted,
				task.TaskStateApproving,
				task.TaskStateApproved,
				task.TaskStateCancelled,
				task.TaskStateTimeout,
			},
		},
		{
			name: "cancelled state (final state)",
			from: task.TaskStateCancelled,
			validTos: []task.TaskState{}, // 终态,不能转换
			invalidTos: []task.TaskState{
				task.TaskStatePending,
				task.TaskStateSubmitted,
				task.TaskStateApproving,
				task.TaskStateApproved,
				task.TaskStateRejected,
				task.TaskStateTimeout,
			},
		},
		{
			name: "timeout state (final state)",
			from: task.TaskStateTimeout,
			validTos: []task.TaskState{}, // 终态,不能转换
			invalidTos: []task.TaskState{
				task.TaskStatePending,
				task.TaskStateSubmitted,
				task.TaskStateApproving,
				task.TaskStateApproved,
				task.TaskStateRejected,
				task.TaskStateCancelled,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validTos, exists := transitions[tt.from]
			if !exists {
				t.Fatalf("状态 %s 没有定义转换规则", tt.from)
			}

			// 验证有效转换
			for _, validTo := range tt.validTos {
				found := false
				for _, to := range validTos {
					if to == validTo {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("状态 %s 应该允许转换到 %s,但规则中未定义", tt.from, validTo)
				}
			}

			// 验证无效转换
			for _, invalidTo := range tt.invalidTos {
				for _, to := range validTos {
					if to == invalidTo {
						t.Errorf("状态 %s 不应该允许转换到 %s", tt.from, invalidTo)
					}
				}
			}
		})
	}
}

// TestFinalStates 验证终态不能转换到其他状态
func TestFinalStates(t *testing.T) {
	finalStates := []task.TaskState{
		task.TaskStateApproved,
		task.TaskStateRejected,
		task.TaskStateCancelled,
		task.TaskStateTimeout,
	}

	transitions := statemachine.GetStateTransitions()

	for _, state := range finalStates {
		validTos := transitions[state]
		if len(validTos) > 0 {
			t.Errorf("终态 %s 不应该允许转换到其他状态,但定义了 %d 个转换", state, len(validTos))
		}
	}
}

