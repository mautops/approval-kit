package statemachine_test

import (
	"testing"

	"github.com/mautops/approval-kit/internal/statemachine"
	"github.com/mautops/approval-kit/internal/task"
)

// TestStateMachineInterface 验证 StateMachine 接口的所有方法
func TestStateMachineInterface(t *testing.T) {
	sm := statemachine.NewStateMachine()

	// 验证接口类型
	var _ statemachine.StateMachine = sm

	// 验证所有方法都存在且可调用
	_ = sm.CanTransition(task.TaskStatePending, task.TaskStateSubmitted)
	_ = sm.GetValidTransitions(task.TaskStatePending)
	_, _ = sm.Transition(createTestTask(task.TaskStatePending), task.TaskStateSubmitted, "test")
}

// TestStateMachineCanTransition 验证 CanTransition 方法
func TestStateMachineCanTransition(t *testing.T) {
	sm := statemachine.NewStateMachine()

	tests := []struct {
		name string
		from task.TaskState
		to   task.TaskState
		want bool
	}{
		{
			name: "valid transition",
			from: task.TaskStatePending,
			to:   task.TaskStateSubmitted,
			want: true,
		},
		{
			name: "invalid transition",
			from: task.TaskStatePending,
			to:   task.TaskStateApproved,
			want: false,
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

// TestStateMachineGetValidTransitions 验证 GetValidTransitions 方法
func TestStateMachineGetValidTransitions(t *testing.T) {
	sm := statemachine.NewStateMachine()

	tests := []struct {
		name           string
		state          task.TaskState
		wantCount      int
		wantContains   []task.TaskState
		wantNotContains []task.TaskState
	}{
		{
			name:      "pending state",
			state:     task.TaskStatePending,
			wantCount: 2,
			wantContains: []task.TaskState{
				task.TaskStateSubmitted,
				task.TaskStateCancelled,
			},
			wantNotContains: []task.TaskState{
				task.TaskStateApproving,
				task.TaskStateApproved,
			},
		},
		{
			name:      "submitted state",
			state:     task.TaskStateSubmitted,
			wantCount: 3, // approving, cancelled, pending (withdraw)
			wantContains: []task.TaskState{
				task.TaskStateApproving,
				task.TaskStateCancelled,
				task.TaskStatePending, // 撤回功能允许
			},
			wantNotContains: []task.TaskState{
				task.TaskStateApproved,
			},
		},
		{
			name:      "approving state",
			state:     task.TaskStateApproving,
			wantCount: 5, // approved, rejected, cancelled, timeout, pending (withdraw)
			wantContains: []task.TaskState{
				task.TaskStateApproved,
				task.TaskStateRejected,
				task.TaskStateCancelled,
				task.TaskStateTimeout,
				task.TaskStatePending, // 撤回功能允许
			},
			wantNotContains: []task.TaskState{
				task.TaskStateSubmitted,
			},
		},
		{
			name:      "approved state (final)",
			state:     task.TaskStateApproved,
			wantCount: 0,
			wantContains: []task.TaskState{},
			wantNotContains: []task.TaskState{
				task.TaskStatePending,
				task.TaskStateSubmitted,
				task.TaskStateApproving,
			},
		},
		{
			name:      "rejected state (final)",
			state:     task.TaskStateRejected,
			wantCount: 0,
			wantContains: []task.TaskState{},
			wantNotContains: []task.TaskState{
				task.TaskStatePending,
				task.TaskStateSubmitted,
				task.TaskStateApproving,
			},
		},
		{
			name:      "cancelled state (final)",
			state:     task.TaskStateCancelled,
			wantCount: 0,
			wantContains: []task.TaskState{},
			wantNotContains: []task.TaskState{
				task.TaskStatePending,
				task.TaskStateSubmitted,
				task.TaskStateApproving,
			},
		},
		{
			name:      "timeout state (final)",
			state:     task.TaskStateTimeout,
			wantCount: 0,
			wantContains: []task.TaskState{},
			wantNotContains: []task.TaskState{
				task.TaskStatePending,
				task.TaskStateSubmitted,
				task.TaskStateApproving,
			},
		},
		{
			name:      "undefined state",
			state:     task.TaskState("undefined"),
			wantCount: 0,
			wantContains: []task.TaskState{},
			wantNotContains: []task.TaskState{
				task.TaskStatePending,
				task.TaskStateSubmitted,
				task.TaskStateApproving,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transitions := sm.GetValidTransitions(tt.state)

			// 验证数量
			if len(transitions) != tt.wantCount {
				t.Errorf("GetValidTransitions(%q) count = %d, want %d", 
					tt.state, len(transitions), tt.wantCount)
			}

			// 验证包含的状态
			for _, wantState := range tt.wantContains {
				found := false
				for _, state := range transitions {
					if state == wantState {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("GetValidTransitions(%q) should contain %q", tt.state, wantState)
				}
			}

			// 验证不包含的状态
			for _, notWantState := range tt.wantNotContains {
				for _, state := range transitions {
					if state == notWantState {
						t.Errorf("GetValidTransitions(%q) should not contain %q", tt.state, notWantState)
					}
				}
			}
		})
	}
}

// TestStateMachineTransition 验证 Transition 方法
func TestStateMachineTransition(t *testing.T) {
	sm := statemachine.NewStateMachine()

	t.Run("valid transition", func(t *testing.T) {
		testTask := createTestTask(task.TaskStatePending)
		newTask, err := sm.Transition(testTask, task.TaskStateSubmitted, "test")

		if err != nil {
			t.Errorf("Transition() error = %v, want nil", err)
		}
		if newTask == nil {
			t.Error("Transition() should return new task on success")
		}
		if newTask.GetState() != task.TaskStateSubmitted {
			t.Errorf("Transition() new state = %v, want %v", 
				newTask.GetState(), task.TaskStateSubmitted)
		}
	})

	t.Run("invalid transition", func(t *testing.T) {
		testTask := createTestTask(task.TaskStatePending)
		newTask, err := sm.Transition(testTask, task.TaskStateApproved, "invalid")

		if err == nil {
			t.Error("Transition() should return error for invalid transition")
		}
		if newTask != nil {
			t.Error("Transition() should not return task on error")
		}
	})
}

// TestStateMachineInterfaceConsistency 验证接口方法的一致性
func TestStateMachineInterfaceConsistency(t *testing.T) {
	sm := statemachine.NewStateMachine()

	// 验证 CanTransition 和 GetValidTransitions 的一致性
	allStates := []task.TaskState{
		task.TaskStatePending,
		task.TaskStateSubmitted,
		task.TaskStateApproving,
		task.TaskStateApproved,
		task.TaskStateRejected,
		task.TaskStateCancelled,
		task.TaskStateTimeout,
	}

	for _, from := range allStates {
		validTransitions := sm.GetValidTransitions(from)

		for _, to := range allStates {
			canTransition := sm.CanTransition(from, to)

			// 如果 CanTransition 返回 true,应该在有效转换列表中
			if canTransition {
				found := false
				for _, validTo := range validTransitions {
					if validTo == to {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Inconsistency: CanTransition(%q, %q) = true, but %q not in GetValidTransitions(%q)",
						from, to, to, from)
				}
			}

			// 如果在有效转换列表中,CanTransition 应该返回 true
			for _, validTo := range validTransitions {
				if validTo == to && !canTransition {
					t.Errorf("Inconsistency: %q in GetValidTransitions(%q), but CanTransition(%q, %q) = false",
						to, from, from, to)
				}
			}
		}
	}
}

// TestStateMachineReturnValueImmutability 验证返回值的不可变性
func TestStateMachineReturnValueImmutability(t *testing.T) {
	sm := statemachine.NewStateMachine()

	// GetValidTransitions 应该返回副本,修改返回值不应该影响内部状态
	transitions1 := sm.GetValidTransitions(task.TaskStatePending)
	transitions2 := sm.GetValidTransitions(task.TaskStatePending)

	// 修改第一个返回值
	if len(transitions1) > 0 {
		transitions1[0] = task.TaskStateApproved
	}

	// 第二个返回值不应该受影响
	if len(transitions2) > 0 && transitions2[0] == task.TaskStateApproved {
		t.Error("GetValidTransitions() should return immutable copy")
	}
}

