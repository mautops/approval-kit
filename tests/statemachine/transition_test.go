package statemachine_test

import (
	stderrors "errors"
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/errors"
	"github.com/mautops/approval-kit/internal/statemachine"
	"github.com/mautops/approval-kit/internal/task"
)

// TestTransition 测试状态转换执行
// 使用纯函数测试,验证不可变实现
func TestTransition(t *testing.T) {
	sm := statemachine.NewStateMachine()

	tests := []struct {
		name      string
		from      task.TaskState
		to        task.TaskState
		reason    string
		wantErr   bool
		wantErrType error
	}{
		// 合法转换
		{
			name:      "pending to submitted",
			from:      task.TaskStatePending,
			to:        task.TaskStateSubmitted,
			reason:    "user submitted",
			wantErr:   false,
			wantErrType: nil,
		},
		{
			name:      "submitted to approving",
			from:      task.TaskStateSubmitted,
			to:        task.TaskStateApproving,
			reason:    "entered approval flow",
			wantErr:   false,
			wantErrType: nil,
		},
		{
			name:      "approving to approved",
			from:      task.TaskStateApproving,
			to:        task.TaskStateApproved,
			reason:    "all approvers approved",
			wantErr:   false,
			wantErrType: nil,
		},
		{
			name:      "approving to rejected",
			from:      task.TaskStateApproving,
			to:        task.TaskStateRejected,
			reason:    "approver rejected",
			wantErr:   false,
			wantErrType: nil,
		},

		// 非法转换
		{
			name:      "pending to approved (invalid)",
			from:      task.TaskStatePending,
			to:        task.TaskStateApproved,
			reason:    "invalid transition",
			wantErr:   true,
			wantErrType: errors.ErrInvalidStateTransition,
		},
		{
			name:      "approved to pending (final state)",
			from:      task.TaskStateApproved,
			to:        task.TaskStatePending,
			reason:    "cannot transition from final state",
			wantErr:   true,
			wantErrType: errors.ErrInvalidStateTransition,
		},
		{
			name:      "same state transition",
			from:      task.TaskStatePending,
			to:        task.TaskStatePending,
			reason:    "same state",
			wantErr:   true,
			wantErrType: errors.ErrInvalidStateTransition,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试任务
			testTask := createTestTask(tt.from)

			// 执行状态转换
			newTask, err := sm.Transition(testTask, tt.to, tt.reason)

			// 验证错误
			if (err != nil) != tt.wantErr {
				t.Errorf("Transition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				// 验证错误类型
			if tt.wantErrType != nil && !stderrors.Is(err, tt.wantErrType) {
				t.Errorf("Transition() error type = %v, want %v", err, tt.wantErrType)
			}
				// 转换失败时,不应该返回新任务
				if newTask != nil {
					t.Errorf("Transition() should not return task on error, got %v", newTask)
				}
				return
			}

			// 验证转换成功
			if newTask == nil {
				t.Fatal("Transition() should return new task on success")
			}

			// 验证新任务状态
			if newTask.GetState() != tt.to {
				t.Errorf("Transition() new state = %v, want %v", newTask.GetState(), tt.to)
			}

			// 验证原任务未被修改(不可变性)
			if testTask.GetState() != tt.from {
				t.Errorf("Transition() original task was modified, state = %v, want %v", 
					testTask.GetState(), tt.from)
			}

			// 验证状态变更历史
			history := newTask.GetStateHistory()
			if len(history) == 0 {
				t.Error("Transition() should record state change history")
			} else {
				lastChange := history[len(history)-1]
				if lastChange.From != tt.from {
					t.Errorf("StateHistory.From = %v, want %v", lastChange.From, tt.from)
				}
				if lastChange.To != tt.to {
					t.Errorf("StateHistory.To = %v, want %v", lastChange.To, tt.to)
				}
				if lastChange.Reason != tt.reason {
					t.Errorf("StateHistory.Reason = %v, want %v", lastChange.Reason, tt.reason)
				}
			}
		})
	}
}

// TestTransitionImmutability 测试状态转换的不可变性
func TestTransitionImmutability(t *testing.T) {
	sm := statemachine.NewStateMachine()

	originalTask := createTestTask(task.TaskStatePending)
	originalState := originalTask.GetState()
	originalUpdatedAt := originalTask.GetUpdatedAt()

	// 执行转换
	newTask, err := sm.Transition(originalTask, task.TaskStateSubmitted, "test")
	if err != nil {
		t.Fatalf("Transition() failed: %v", err)
	}

	// 等待一小段时间,确保时间戳不同
	time.Sleep(10 * time.Millisecond)

	// 验证原任务未被修改
	if originalTask.GetState() != originalState {
		t.Error("Original task state was modified")
	}
	if originalTask.GetUpdatedAt() != originalUpdatedAt {
		t.Error("Original task UpdatedAt was modified")
	}

	// 验证新任务有新的时间戳
	if newTask.GetUpdatedAt().Before(originalTask.GetUpdatedAt()) || 
		newTask.GetUpdatedAt().Equal(originalTask.GetUpdatedAt()) {
		t.Error("New task should have updated timestamp")
	}
}

// TestTransitionHistory 测试状态变更历史记录
func TestTransitionHistory(t *testing.T) {
	sm := statemachine.NewStateMachine()

	testTask := createTestTask(task.TaskStatePending)

	// 执行多次转换
	transitions := []struct {
		to     task.TaskState
		reason string
	}{
		{task.TaskStateSubmitted, "submitted"},
		{task.TaskStateApproving, "entered approval"},
		{task.TaskStateApproved, "approved"},
	}

	var err error
	var currentTask statemachine.TransitionableTask = testTask
	for _, trans := range transitions {
		currentTask, err = sm.Transition(currentTask, trans.to, trans.reason)
		if err != nil {
			t.Fatalf("Transition() failed: %v", err)
		}
	}

	// 验证历史记录
	history := currentTask.GetStateHistory()
	if len(history) != len(transitions) {
		t.Errorf("StateHistory length = %d, want %d", len(history), len(transitions))
	}

	// 验证历史记录的顺序和内容
	expectedFrom := task.TaskStatePending
	for i, trans := range transitions {
		change := history[i]
		if change.From != expectedFrom {
			t.Errorf("History[%d].From = %v, want %v", i, change.From, expectedFrom)
		}
		if change.To != trans.to {
			t.Errorf("History[%d].To = %v, want %v", i, change.To, trans.to)
		}
		if change.Reason != trans.reason {
			t.Errorf("History[%d].Reason = %v, want %v", i, change.Reason, trans.reason)
		}
		expectedFrom = trans.to
	}
}

// 辅助函数: 创建测试任务
// 为了支持测试,我们需要一个最小化的 Task 接口
func createTestTask(state task.TaskState) *TestTask {
	return &TestTask{
		state:     state,
		updatedAt: time.Now(),
		history:   []*statemachine.StateChange{},
	}
}

// TestTask 测试用的任务结构
// 实现 TransitionableTask 接口以支持状态转换测试
type TestTask struct {
	state     task.TaskState
	updatedAt time.Time
	history   []*statemachine.StateChange
}

func (t *TestTask) GetState() task.TaskState {
	return t.state
}

func (t *TestTask) SetState(state task.TaskState) {
	t.state = state
}

func (t *TestTask) GetUpdatedAt() time.Time {
	return t.updatedAt
}

func (t *TestTask) SetUpdatedAt(tm time.Time) {
	t.updatedAt = tm
}

func (t *TestTask) GetStateHistory() []*statemachine.StateChange {
	return t.history
}

func (t *TestTask) AddStateChange(change *statemachine.StateChange) {
	t.history = append(t.history, change)
}

func (t *TestTask) Clone() statemachine.TransitionableTask {
	historyCopy := make([]*statemachine.StateChange, len(t.history))
	copy(historyCopy, t.history)
	return &TestTask{
		state:     t.state,
		updatedAt: t.updatedAt,
		history:   historyCopy,
	}
}

