package tests

import (
	stderrors "errors"
	"fmt"
	"testing"

	"github.com/mautops/approval-kit/internal/errors"
)

// TestErrorConstants 验证所有错误常量已定义
func TestErrorConstants(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "ErrInvalidTemplate",
			err:  errors.ErrInvalidTemplate,
			want: "invalid template",
		},
		{
			name: "ErrInvalidStateTransition",
			err:  errors.ErrInvalidStateTransition,
			want: "invalid state transition",
		},
		{
			name: "ErrNodeNotFound",
			err:  errors.ErrNodeNotFound,
			want: "node not found",
		},
		{
			name: "ErrApproverNotFound",
			err:  errors.ErrApproverNotFound,
			want: "approver not found",
		},
		{
			name: "ErrApprovalPending",
			err:  errors.ErrApprovalPending,
			want: "approval pending",
		},
		{
			name: "ErrConcurrentModification",
			err:  errors.ErrConcurrentModification,
			want: "concurrent modification",
		},
		{
			name: "ErrEventPushFailed",
			err:  errors.ErrEventPushFailed,
			want: "event push failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Errorf("%s 未定义", tt.name)
				return
			}
			if tt.err.Error() != tt.want {
				t.Errorf("%s.Error() = %q, want %q", tt.name, tt.err.Error(), tt.want)
			}
		})
	}
}

// TestStateTransitionError 验证 StateTransitionError 错误类型
// 注意: StateTransitionError 的详细测试在 tests/task 包中
func TestStateTransitionError(t *testing.T) {
	// 这个测试已经移到 tests/task 包中,避免循环依赖
	t.Skip("StateTransitionError 测试已移到 tests/task 包中")
}

// TestErrorWrapping 验证错误包装机制
func TestErrorWrapping(t *testing.T) {
	baseErr := errors.ErrInvalidTemplate
	wrappedErr := fmt.Errorf("template validation failed: %w", baseErr)

	// 验证错误可以被展开
	if !stderrors.Is(wrappedErr, baseErr) {
		t.Errorf("wrapped error should be unwrappable to base error")
	}

	// 验证错误消息包含原始错误
	if !contains(wrappedErr.Error(), baseErr.Error()) {
		t.Errorf("wrapped error message should contain base error message")
	}
}

// TestErrorComparison 验证错误比较
func TestErrorComparison(t *testing.T) {
	err1 := errors.ErrInvalidTemplate
	err2 := errors.ErrInvalidTemplate
	err3 := errors.ErrInvalidStateTransition

	// 相同错误应该相等
	if err1 != err2 {
		t.Errorf("相同错误常量应该相等")
	}

	// 不同错误应该不相等
	if err1 == err3 {
		t.Errorf("不同错误常量应该不相等")
	}
}

// 辅助函数
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

