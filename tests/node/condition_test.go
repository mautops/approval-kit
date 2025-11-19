package node_test

import (
	"testing"

	"github.com/mautops/approval-kit/internal/node"
)

// TestConditionStructure 测试 Condition 结构体定义
func TestConditionStructure(t *testing.T) {
	// 测试创建 Condition 结构体
	condition := &node.Condition{
		Type:   "numeric",
		Config: nil, // 配置将在后续任务中完善
	}

	if condition.Type != "numeric" {
		t.Errorf("Condition.Type = %q, want %q", condition.Type, "numeric")
	}
}

// TestConditionConfigInterface 测试 ConditionConfig 接口
func TestConditionConfigInterface(t *testing.T) {
	// 验证 ConditionConfig 接口存在
	var _ node.ConditionConfig = (*mockConditionConfig)(nil)
}

// mockConditionConfig 用于测试的 ConditionConfig 实现
type mockConditionConfig struct {
	conditionType string
}

func (m *mockConditionConfig) ConditionType() string {
	return m.conditionType
}

// TestConditionConfigConditionType 测试 ConditionType 方法
func TestConditionConfigConditionType(t *testing.T) {
	config := &mockConditionConfig{
		conditionType: "numeric",
	}

	conditionType := config.ConditionType()
	if conditionType != "numeric" {
		t.Errorf("ConditionConfig.ConditionType() = %q, want %q", conditionType, "numeric")
	}
}

// TestConditionValidation 测试条件验证
func TestConditionValidation(t *testing.T) {
	tests := []struct {
		name      string
		condition *node.Condition
		wantErr   bool
	}{
		{
			name: "valid condition with type and config",
			condition: &node.Condition{
				Type:   "numeric",
				Config: &mockConditionConfig{conditionType: "numeric"},
			},
			wantErr: false,
		},
		{
			name: "empty condition type",
			condition: &node.Condition{
				Type:   "",
				Config: &mockConditionConfig{conditionType: ""},
			},
			wantErr: true,
		},
		{
			name: "missing config",
			condition: &node.Condition{
				Type:   "numeric",
				Config: nil,
			},
			wantErr: true,
		},
		{
			name: "type mismatch",
			condition: &node.Condition{
				Type:   "numeric",
				Config: &mockConditionConfig{conditionType: "string"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.condition.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Condition.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

