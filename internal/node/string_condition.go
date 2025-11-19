package node

import (
	"encoding/json"
	"fmt"
	"strings"
)

// StringConditionConfig 字符串匹配条件配置
type StringConditionConfig struct {
	// Field 字段名(JSON 路径)
	Field string

	// Operator 匹配操作符
	// 支持: "eq"(等于), "contains"(包含), "starts_with"(以...开始), "ends_with"(以...结束)
	Operator string

	// Value 匹配值
	Value string

	// Source 数据源
	// 支持: "task_params"(任务参数), "node_outputs"(节点输出数据)
	Source string

	// NodeID 节点 ID(当 Source 为 "node_outputs" 时必填)
	NodeID string
}

// ConditionType 返回条件类型(实现 ConditionConfig 接口)
func (c *StringConditionConfig) ConditionType() string {
	return "string"
}

// StringConditionEvaluator 字符串匹配条件评估器
type StringConditionEvaluator struct{}

// NewStringConditionEvaluator 创建新的字符串匹配条件评估器
func NewStringConditionEvaluator() ConditionEvaluator {
	return &StringConditionEvaluator{}
}

// Supports 检查是否支持指定的条件类型(实现 ConditionEvaluator 接口)
func (e *StringConditionEvaluator) Supports(conditionType string) bool {
	return conditionType == "string"
}

// Evaluate 评估字符串匹配条件(实现 ConditionEvaluator 接口)
func (e *StringConditionEvaluator) Evaluate(condition *Condition, ctx *NodeContext) (bool, error) {
	config, ok := condition.Config.(*StringConditionConfig)
	if !ok {
		return false, fmt.Errorf("invalid condition config type for string condition")
	}

	// 获取字段值
	value, err := e.getValue(config, ctx)
	if err != nil {
		return false, err
	}

	// 执行匹配
	return e.match(value, config.Operator, config.Value)
}

// getValue 从上下文中获取字段值
func (e *StringConditionEvaluator) getValue(config *StringConditionConfig, ctx *NodeContext) (string, error) {
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

// match 执行字符串匹配
func (e *StringConditionEvaluator) match(value string, operator string, target string) (bool, error) {
	switch operator {
	case "eq":
		return value == target, nil
	case "contains":
		return strings.Contains(value, target), nil
	case "starts_with":
		return strings.HasPrefix(value, target), nil
	case "ends_with":
		return strings.HasSuffix(value, target), nil
	default:
		return false, fmt.Errorf("unsupported operator: %q", operator)
	}
}

