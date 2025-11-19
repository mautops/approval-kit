package node

import (
	"fmt"
)

// CompositeConditionConfig 组合条件配置
type CompositeConditionConfig struct {
	// Operator 组合操作符
	// 支持: "and"(与), "or"(或)
	Operator string

	// Conditions 子条件列表
	Conditions []*Condition
}

// ConditionType 返回条件类型(实现 ConditionConfig 接口)
func (c *CompositeConditionConfig) ConditionType() string {
	return "composite"
}

// CompositeConditionEvaluator 组合条件评估器
type CompositeConditionEvaluator struct {
	registry *ConditionEvaluatorRegistry
}

// NewCompositeConditionEvaluator 创建新的组合条件评估器
// registry: 条件评估器注册表(必须已包含所有基础评估器)
func NewCompositeConditionEvaluator(registry *ConditionEvaluatorRegistry) ConditionEvaluator {
	// 创建组合条件评估器实例
	compositeEvaluator := &CompositeConditionEvaluator{
		registry: registry,
	}
	// 注册组合条件评估器自身(支持嵌套)
	registry.Register(compositeEvaluator)
	return compositeEvaluator
}

// Supports 检查是否支持指定的条件类型(实现 ConditionEvaluator 接口)
func (e *CompositeConditionEvaluator) Supports(conditionType string) bool {
	return conditionType == "composite"
}

// Evaluate 评估组合条件(实现 ConditionEvaluator 接口)
func (e *CompositeConditionEvaluator) Evaluate(condition *Condition, ctx *NodeContext) (bool, error) {
	config, ok := condition.Config.(*CompositeConditionConfig)
	if !ok {
		return false, fmt.Errorf("invalid condition config type for composite condition")
	}

	// 验证条件列表不为空
	if len(config.Conditions) == 0 {
		return false, fmt.Errorf("composite condition must have at least one sub-condition")
	}

	// 评估所有子条件
	results := make([]bool, len(config.Conditions))
	for i, subCondition := range config.Conditions {
		// 获取对应的评估器
		evaluator := e.registry.GetEvaluator(subCondition.Type)
		if evaluator == nil {
			return false, fmt.Errorf("no evaluator found for condition type: %q", subCondition.Type)
		}

		// 评估子条件
		result, err := evaluator.Evaluate(subCondition, ctx)
		if err != nil {
			return false, fmt.Errorf("failed to evaluate sub-condition %d: %w", i, err)
		}
		results[i] = result
	}

	// 根据操作符组合结果
	return e.combine(results, config.Operator)
}

// combine 组合多个条件结果
func (e *CompositeConditionEvaluator) combine(results []bool, operator string) (bool, error) {
	if len(results) == 0 {
		return false, fmt.Errorf("no results to combine")
	}

	switch operator {
	case "and":
		// AND: 所有结果都为 true 时返回 true
		for _, result := range results {
			if !result {
				return false, nil
			}
		}
		return true, nil
	case "or":
		// OR: 任意一个结果为 true 时返回 true
		for _, result := range results {
			if result {
				return true, nil
			}
		}
		return false, nil
	default:
		return false, fmt.Errorf("unsupported operator: %q", operator)
	}
}

// ConditionEvaluatorRegistry 条件评估器注册表
type ConditionEvaluatorRegistry struct {
	evaluators map[string]ConditionEvaluator
}

// NewConditionEvaluatorRegistry 创建新的条件评估器注册表
func NewConditionEvaluatorRegistry() *ConditionEvaluatorRegistry {
	registry := &ConditionEvaluatorRegistry{
		evaluators: make(map[string]ConditionEvaluator),
	}

	// 注册默认评估器
	registry.Register(NewNumericConditionEvaluator())
	registry.Register(NewStringConditionEvaluator())
	registry.Register(NewEnumConditionEvaluator())
	// 注册组合条件评估器(支持嵌套,传入已注册基础评估器的 registry)
	compositeEvaluator := NewCompositeConditionEvaluator(registry)
	registry.Register(compositeEvaluator)

	return registry
}

// Register 注册条件评估器
func (r *ConditionEvaluatorRegistry) Register(evaluator ConditionEvaluator) {
	// 注册所有支持的条件类型
	for _, conditionType := range []string{"numeric", "string", "enum", "composite"} {
		if evaluator.Supports(conditionType) {
			r.evaluators[conditionType] = evaluator
		}
	}
}

// GetEvaluator 获取条件评估器
func (r *ConditionEvaluatorRegistry) GetEvaluator(conditionType string) ConditionEvaluator {
	return r.evaluators[conditionType]
}

