package node_test

import (
	"testing"

	"github.com/mautops/approval-kit/internal/node"
)

// TestConditionEvaluatorInterface 测试 ConditionEvaluator 接口定义
func TestConditionEvaluatorInterface(t *testing.T) {
	// 这个测试验证 ConditionEvaluator 接口是否存在以及方法签名是否正确
	// 接口应该包含 Evaluate 方法用于评估条件

	// 验证接口存在
	var _ node.ConditionEvaluator = (*mockConditionEvaluator)(nil)
}

// mockConditionEvaluator 用于测试的 ConditionEvaluator 实现
type mockConditionEvaluator struct{}

func (m *mockConditionEvaluator) Evaluate(condition *node.Condition, ctx *node.NodeContext) (bool, error) {
	return false, nil
}

func (m *mockConditionEvaluator) Supports(conditionType string) bool {
	return false
}

// TestConditionEvaluatorEvaluate 测试 Evaluate 方法
func TestConditionEvaluatorEvaluate(t *testing.T) {
	evaluator := &mockConditionEvaluator{}

	// 创建测试条件
	condition := &node.Condition{
		Type: "test",
	}

	// 创建测试上下文
	ctx := &node.NodeContext{}

	// 调用 Evaluate 方法
	result, err := evaluator.Evaluate(condition, ctx)
	if err != nil {
		t.Fatalf("Evaluate() should not return error: %v", err)
	}

	// 验证返回值类型
	_ = result // 验证返回 bool 类型
}

// TestConditionEvaluatorSupports 测试 Supports 方法
func TestConditionEvaluatorSupports(t *testing.T) {
	evaluator := &mockConditionEvaluator{}

	// 调用 Supports 方法
	supports := evaluator.Supports("test")
	_ = supports // 验证返回 bool 类型
}

