package task_test

import (
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/task"
)

// TestStateChangeValidate 测试状态变更验证逻辑
func TestStateChangeValidate(t *testing.T) {
	tests := []struct {
		name    string
		change  *task.StateChange
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid state change",
			change: &task.StateChange{
				From:   task.TaskStatePending,
				To:     task.TaskStateSubmitted,
				Reason: "user submitted",
				Time:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "missing From",
			change: &task.StateChange{
				To:     task.TaskStateSubmitted,
				Reason: "test",
				Time:   time.Now(),
			},
			wantErr: true,
			errMsg:  "state change From is required",
		},
		{
			name: "missing To",
			change: &task.StateChange{
				From:   task.TaskStatePending,
				Reason: "test",
				Time:   time.Now(),
			},
			wantErr: true,
			errMsg:  "state change To is required",
		},
		{
			name: "same From and To",
			change: &task.StateChange{
				From:   task.TaskStatePending,
				To:     task.TaskStatePending,
				Reason: "test",
				Time:   time.Now(),
			},
			wantErr: true,
			errMsg:  "state change From and To cannot be the same",
		},
		{
			name: "zero Time",
			change: &task.StateChange{
				From:   task.TaskStatePending,
				To:     task.TaskStateSubmitted,
				Reason: "test",
				// Time is zero value
			},
			wantErr: true,
			errMsg:  "state change Time is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.change.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("StateChange.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				if tt.errMsg != "" && err.Error()[:len(tt.errMsg)] != tt.errMsg {
					t.Errorf("StateChange.Validate() error message = %q, want prefix %q", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

// TestNewStateChange 测试创建新状态变更记录
func TestNewStateChange(t *testing.T) {
	change := task.NewStateChange(
		task.TaskStatePending,
		task.TaskStateSubmitted,
		"user submitted",
	)

	if change.From != task.TaskStatePending {
		t.Errorf("NewStateChange().From = %v, want %v", change.From, task.TaskStatePending)
	}
	if change.To != task.TaskStateSubmitted {
		t.Errorf("NewStateChange().To = %v, want %v", change.To, task.TaskStateSubmitted)
	}
	if change.Reason != "user submitted" {
		t.Errorf("NewStateChange().Reason = %q, want %q", change.Reason, "user submitted")
	}
	if change.Time.IsZero() {
		t.Error("NewStateChange().Time should be set")
	}

	// 验证状态变更记录是有效的
	if err := change.Validate(); err != nil {
		t.Errorf("NewStateChange() created invalid state change: %v", err)
	}
}

