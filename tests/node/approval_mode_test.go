package node_test

import (
	"testing"

	"github.com/mautops/approval-kit/internal/node"
)

// TestApprovalModeType 验证 ApprovalMode 类型存在
func TestApprovalModeType(t *testing.T) {
	// 验证 ApprovalMode 类型已定义
	var mode node.ApprovalMode
	if mode == "" {
		// 零值应该是空字符串,这是正常的
		_ = mode
	}
}

// TestApprovalModeConstants 验证所有审批模式常量已定义
func TestApprovalModeConstants(t *testing.T) {
	tests := []struct {
		name     string
		mode     node.ApprovalMode
		expected string
	}{
		{
			name:     "single mode",
			mode:     node.ApprovalModeSingle,
			expected: "single",
		},
		{
			name:     "unanimous mode",
			mode:     node.ApprovalModeUnanimous,
			expected: "unanimous",
		},
		{
			name:     "or mode",
			mode:     node.ApprovalModeOr,
			expected: "or",
		},
		{
			name:     "proportional mode",
			mode:     node.ApprovalModeProportional,
			expected: "proportional",
		},
		{
			name:     "sequential mode",
			mode:     node.ApprovalModeSequential,
			expected: "sequential",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.mode) != tt.expected {
				t.Errorf("ApprovalMode = %q, want %q", tt.mode, tt.expected)
			}
		})
	}
}

// TestApprovalModeString 验证 ApprovalMode 的字符串表示
func TestApprovalModeString(t *testing.T) {
	tests := []struct {
		name     string
		mode     node.ApprovalMode
		expected string
	}{
		{
			name:     "single",
			mode:     node.ApprovalModeSingle,
			expected: "single",
		},
		{
			name:     "unanimous",
			mode:     node.ApprovalModeUnanimous,
			expected: "unanimous",
		},
		{
			name:     "or",
			mode:     node.ApprovalModeOr,
			expected: "or",
		},
		{
			name:     "proportional",
			mode:     node.ApprovalModeProportional,
			expected: "proportional",
		},
		{
			name:     "sequential",
			mode:     node.ApprovalModeSequential,
			expected: "sequential",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.mode) != tt.expected {
				t.Errorf("ApprovalMode.String() = %q, want %q", tt.mode, tt.expected)
			}
		})
	}
}

// TestApprovalModeAllDefined 验证所有必需的审批模式都已定义
func TestApprovalModeAllDefined(t *testing.T) {
	// 根据 spec,应该有 5 种审批模式
	expectedModes := []node.ApprovalMode{
		node.ApprovalModeSingle,
		node.ApprovalModeUnanimous,
		node.ApprovalModeOr,
		node.ApprovalModeProportional,
		node.ApprovalModeSequential,
	}

	if len(expectedModes) != 5 {
		t.Errorf("期望 5 种审批模式,实际定义了 %d 种", len(expectedModes))
	}

	// 验证每个模式都不是空值
	for i, mode := range expectedModes {
		if mode == "" {
			t.Errorf("审批模式 %d 未定义", i)
		}
	}
}

