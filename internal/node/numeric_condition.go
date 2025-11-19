package node

import (
	"encoding/json"
	"fmt"
)

// NumericConditionConfig 数值比较条件配置
type NumericConditionConfig struct {
	// Field 字段名(JSON 路径)
	Field string

	// Operator 比较操作符
	// 支持: "gt"(大于), "lt"(小于), "eq"(等于), "gte"(大于等于), "lte"(小于等于)
	Operator string

	// Value 比较值
	Value float64

	// Source 数据源
	// 支持: "task_params"(任务参数), "node_outputs"(节点输出数据)
	Source string

	// NodeID 节点 ID(当 Source 为 "node_outputs" 时必填)
	NodeID string
}

// ConditionType 返回条件类型(实现 ConditionConfig 接口)
func (c *NumericConditionConfig) ConditionType() string {
	return "numeric"
}

// NumericConditionEvaluator 数值比较条件评估器
type NumericConditionEvaluator struct{}

// NewNumericConditionEvaluator 创建新的数值比较条件评估器
func NewNumericConditionEvaluator() ConditionEvaluator {
	return &NumericConditionEvaluator{}
}

// Supports 检查是否支持指定的条件类型(实现 ConditionEvaluator 接口)
func (e *NumericConditionEvaluator) Supports(conditionType string) bool {
	return conditionType == "numeric"
}

// Evaluate 评估数值比较条件(实现 ConditionEvaluator 接口)
func (e *NumericConditionEvaluator) Evaluate(condition *Condition, ctx *NodeContext) (bool, error) {
	config, ok := condition.Config.(*NumericConditionConfig)
	if !ok {
		return false, fmt.Errorf("invalid condition config type for numeric condition")
	}

	// 获取字段值
	value, err := e.getValue(config, ctx)
	if err != nil {
		return false, err
	}

	// 执行比较
	return e.compare(value, config.Operator, config.Value)
}

// getValue 从上下文中获取字段值
func (e *NumericConditionEvaluator) getValue(config *NumericConditionConfig, ctx *NodeContext) (float64, error) {
	var data json.RawMessage

	switch config.Source {
	case "task_params":
		data = ctx.Task.Params
	case "node_outputs":
		if config.NodeID == "" {
			return 0, fmt.Errorf("NodeID is required when Source is 'node_outputs'")
		}
		var exists bool
		data, exists = ctx.Task.NodeOutputs[config.NodeID]
		if !exists {
			return 0, fmt.Errorf("node output not found: %q", config.NodeID)
		}
	default:
		return 0, fmt.Errorf("unsupported source: %q", config.Source)
	}

	// 解析 JSON 获取字段值
	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return 0, fmt.Errorf("failed to parse JSON: %w", err)
	}

	fieldValue, exists := jsonData[config.Field]
	if !exists {
		return 0, fmt.Errorf("field not found: %q", config.Field)
	}

	// 转换为 float64
	switch v := fieldValue.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("field value is not numeric: %T", fieldValue)
	}
}

// compare 执行数值比较
func (e *NumericConditionEvaluator) compare(value float64, operator string, target float64) (bool, error) {
	switch operator {
	case "gt":
		return value > target, nil
	case "lt":
		return value < target, nil
	case "eq":
		return value == target, nil
	case "gte":
		return value >= target, nil
	case "lte":
		return value <= target, nil
	default:
		return false, fmt.Errorf("unsupported operator: %q", operator)
	}
}

