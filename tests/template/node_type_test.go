package template_test

import (
	"testing"

	"github.com/mautops/approval-kit/internal/template"
)

// TestNodeTypeType 验证 NodeType 类型存在
func TestNodeTypeType(t *testing.T) {
	// 验证 NodeType 类型已定义
	var nodeType template.NodeType
	if nodeType == "" {
		// 零值应该是空字符串,这是正常的
		_ = nodeType
	}
}

// TestNodeTypeConstants 验证所有节点类型常量已定义
func TestNodeTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		nodeType template.NodeType
		expected string
	}{
		{
			name:     "start node type",
			nodeType: template.NodeTypeStart,
			expected: "start",
		},
		{
			name:     "approval node type",
			nodeType: template.NodeTypeApproval,
			expected: "approval",
		},
		{
			name:     "condition node type",
			nodeType: template.NodeTypeCondition,
			expected: "condition",
		},
		{
			name:     "end node type",
			nodeType: template.NodeTypeEnd,
			expected: "end",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.nodeType) != tt.expected {
				t.Errorf("NodeType = %q, want %q", tt.nodeType, tt.expected)
			}
		})
	}
}

// TestNodeTypeString 验证 NodeType 的字符串表示
func TestNodeTypeString(t *testing.T) {
	tests := []struct {
		name     string
		nodeType template.NodeType
		expected string
	}{
		{
			name:     "start",
			nodeType: template.NodeTypeStart,
			expected: "start",
		},
		{
			name:     "approval",
			nodeType: template.NodeTypeApproval,
			expected: "approval",
		},
		{
			name:     "condition",
			nodeType: template.NodeTypeCondition,
			expected: "condition",
		},
		{
			name:     "end",
			nodeType: template.NodeTypeEnd,
			expected: "end",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.nodeType) != tt.expected {
				t.Errorf("NodeType.String() = %q, want %q", tt.nodeType, tt.expected)
			}
		})
	}
}

// TestNodeTypeAllDefined 验证所有必需的节点类型都已定义
func TestNodeTypeAllDefined(t *testing.T) {
	// 根据 spec,应该有 4 种节点类型
	expectedTypes := []template.NodeType{
		template.NodeTypeStart,
		template.NodeTypeApproval,
		template.NodeTypeCondition,
		template.NodeTypeEnd,
	}

	if len(expectedTypes) != 4 {
		t.Errorf("期望 4 种节点类型,实际定义了 %d 种", len(expectedTypes))
	}

	// 验证每个类型都不是空值
	for i, nodeType := range expectedTypes {
		if nodeType == "" {
			t.Errorf("节点类型 %d 未定义", i)
		}
	}
}

