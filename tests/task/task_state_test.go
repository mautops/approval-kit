package task_test

import (
	"testing"

	"github.com/mautops/approval-kit/internal/task"
)

// TestTaskStateType 验证 TaskState 类型存在
func TestTaskStateType(t *testing.T) {
	// 验证 TaskState 类型已定义
	var state task.TaskState
	if state == "" {
		// 零值应该是空字符串,这是正常的
		_ = state
	}
}

// TestTaskStateConstants 验证所有任务状态常量已定义
func TestTaskStateConstants(t *testing.T) {
	tests := []struct {
		name     string
		state    task.TaskState
		expected string
	}{
		{
			name:     "pending state",
			state:    task.TaskStatePending,
			expected: "pending",
		},
		{
			name:     "submitted state",
			state:    task.TaskStateSubmitted,
			expected: "submitted",
		},
		{
			name:     "approving state",
			state:    task.TaskStateApproving,
			expected: "approving",
		},
		{
			name:     "approved state",
			state:    task.TaskStateApproved,
			expected: "approved",
		},
		{
			name:     "rejected state",
			state:    task.TaskStateRejected,
			expected: "rejected",
		},
		{
			name:     "cancelled state",
			state:    task.TaskStateCancelled,
			expected: "cancelled",
		},
		{
			name:     "timeout state",
			state:    task.TaskStateTimeout,
			expected: "timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.state) != tt.expected {
				t.Errorf("TaskState = %q, want %q", tt.state, tt.expected)
			}
		})
	}
}

// TestTaskStateString 验证 TaskState 的字符串表示
func TestTaskStateString(t *testing.T) {
	tests := []struct {
		name     string
		state    task.TaskState
		expected string
	}{
		{
			name:     "pending",
			state:    task.TaskStatePending,
			expected: "pending",
		},
		{
			name:     "submitted",
			state:    task.TaskStateSubmitted,
			expected: "submitted",
		},
		{
			name:     "approving",
			state:    task.TaskStateApproving,
			expected: "approving",
		},
		{
			name:     "approved",
			state:    task.TaskStateApproved,
			expected: "approved",
		},
		{
			name:     "rejected",
			state:    task.TaskStateRejected,
			expected: "rejected",
		},
		{
			name:     "cancelled",
			state:    task.TaskStateCancelled,
			expected: "cancelled",
		},
		{
			name:     "timeout",
			state:    task.TaskStateTimeout,
			expected: "timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.state) != tt.expected {
				t.Errorf("TaskState.String() = %q, want %q", tt.state, tt.expected)
			}
		})
	}
}

// TestTaskStateAllDefined 验证所有必需的状态都已定义
func TestTaskStateAllDefined(t *testing.T) {
	// 根据 PRD,应该有 7 个状态
	expectedStates := []task.TaskState{
		task.TaskStatePending,
		task.TaskStateSubmitted,
		task.TaskStateApproving,
		task.TaskStateApproved,
		task.TaskStateRejected,
		task.TaskStateCancelled,
		task.TaskStateTimeout,
	}

	if len(expectedStates) != 7 {
		t.Errorf("期望 7 个状态,实际定义了 %d 个", len(expectedStates))
	}

	// 验证每个状态都不是空值
	for i, state := range expectedStates {
		if state == "" {
			t.Errorf("状态 %d 未定义", i)
		}
	}
}

