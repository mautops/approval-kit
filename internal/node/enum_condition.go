package node

import (
	"encoding/json"
	"fmt"
)

// EnumConditionConfig 枚举判断条件配置
type EnumConditionConfig struct {
	// Field 字段名(JSON 路径)
	Field string

	// Operator 判断操作符
	// 支持: "in"(在列表中), "not_in"(不在列表中)
	Operator string

	// Values 枚举值列表
	Values []string

	// Source 数据源
	// 支持: "task_params"(任务参数), "node_outputs"(节点输出数据)
	Source string

	// NodeID 节点 ID(当 Source 为 "node_outputs" 时必填)
	NodeID string
}

// ConditionType 返回条件类型(实现 ConditionConfig 接口)
func (c *EnumConditionConfig) ConditionType() string {
	return "enum"
}

// EnumConditionEvaluator 枚举判断条件评估器
type EnumConditionEvaluator struct{}

// NewEnumConditionEvaluator 创建新的枚举判断条件评估器
func NewEnumConditionEvaluator() ConditionEvaluator {
	return &EnumConditionEvaluator{}
}

// Supports 检查是否支持指定的条件类型(实现 ConditionEvaluator 接口)
func (e *EnumConditionEvaluator) Supports(conditionType string) bool {
	return conditionType == "enum"
}

// Evaluate 评估枚举判断条件(实现 ConditionEvaluator 接口)
func (e *EnumConditionEvaluator) Evaluate(condition *Condition, ctx *NodeContext) (bool, error) {
	config, ok := condition.Config.(*EnumConditionConfig)
	if !ok {
		return false, fmt.Errorf("invalid condition config type for enum condition")
	}

	// 获取字段值
	value, err := e.getValue(config, ctx)
	if err != nil {
		return false, err
	}

	// 执行判断
	return e.check(value, config.Operator, config.Values)
}

// getValue 从上下文中获取字段值
func (e *EnumConditionEvaluator) getValue(config *EnumConditionConfig, ctx *NodeContext) (string, error) {
	var data json.RawMessage

	switch config.Source {
	case "task_params":
		data = ctx.Task.Params
	case "node_outputs":
		if config.NodeID == "" {
			return "", fmt.Errorf("NodeID is required when Source is 'node_outputs'")
		}
		var exists bool
		data, exists = ctx.Task.NodeOutputs[config.NodeID]
		if !exists {
			return "", fmt.Errorf("node output not found: %q", config.NodeID)
		}
	default:
		return "", fmt.Errorf("unsupported source: %q", config.Source)
	}

	// 解析 JSON 获取字段值
	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	fieldValue, exists := jsonData[config.Field]
	if !exists {
		return "", fmt.Errorf("field not found: %q", config.Field)
	}

	// 转换为 string
	strValue, ok := fieldValue.(string)
	if !ok {
		return "", fmt.Errorf("field value is not string: %T", fieldValue)
	}

	return strValue, nil
}

// check 执行枚举判断
func (e *EnumConditionEvaluator) check(value string, operator string, values []string) (bool, error) {
	// 检查值是否在列表中
	inList := false
	for _, v := range values {
		if v == value {
			inList = true
			break
		}
	}

	switch operator {
	case "in":
		return inList, nil
	case "not_in":
		return !inList, nil
	default:
		return false, fmt.Errorf("unsupported operator: %q", operator)
	}
}

