package task_test

import (
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/task"
)

// TestStateChangeStruct 验证 StateChange 结构体定义
func TestStateChangeStruct(t *testing.T) {
	// 验证 StateChange 类型存在
	var change *task.StateChange
	if change != nil {
		_ = change
	}
}

// TestStateChangeFields 验证 StateChange 结构体的所有字段
func TestStateChangeFields(t *testing.T) {
	now := time.Now()
	change := &task.StateChange{
		From:   task.TaskStatePending,
		To:     task.TaskStateSubmitted,
		Reason: "user submitted",
		Time:   now,
	}

	// 验证字段值
	if change.From != task.TaskStatePending {
		t.Errorf("StateChange.From = %v, want %v", change.From, task.TaskStatePending)
	}
	if change.To != task.TaskStateSubmitted {
		t.Errorf("StateChange.To = %v, want %v", change.To, task.TaskStateSubmitted)
	}
	if change.Reason != "user submitted" {
		t.Errorf("StateChange.Reason = %q, want %q", change.Reason, "user submitted")
	}
	if !change.Time.Equal(now) {
		t.Errorf("StateChange.Time = %v, want %v", change.Time, now)
	}
}

// TestStateChangeZeroValue 验证 StateChange 的零值
func TestStateChangeZeroValue(t *testing.T) {
	var change task.StateChange

	// 验证零值
	if change.From != "" {
		t.Errorf("StateChange zero value From = %q, want empty string", change.From)
	}
	if change.To != "" {
		t.Errorf("StateChange zero value To = %q, want empty string", change.To)
	}
	if change.Reason != "" {
		t.Errorf("StateChange zero value Reason = %q, want empty string", change.Reason)
	}
	if !change.Time.IsZero() {
		t.Errorf("StateChange zero value Time should be zero, got %v", change.Time)
	}
}

// TestStateChangeAllStates 验证所有状态转换
func TestStateChangeAllStates(t *testing.T) {
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
			change := &task.StateChange{
				From:   from,
				To:     to,
				Reason: "test",
				Time:   time.Now(),
			}

			if change.From != from {
				t.Errorf("StateChange.From = %v, want %v", change.From, from)
			}
			if change.To != to {
				t.Errorf("StateChange.To = %v, want %v", change.To, to)
			}
		}
	}
}

