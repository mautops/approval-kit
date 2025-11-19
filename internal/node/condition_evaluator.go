package node

// ConditionEvaluator 条件评估器接口
// 用于评估条件节点的条件表达式,决定流程走向
type ConditionEvaluator interface {
	// Evaluate 评估条件
	// condition: 条件定义
	// ctx: 节点执行上下文,包含任务参数、节点输出数据等
	// 返回: 条件评估结果(true/false)和错误信息
	Evaluate(condition *Condition, ctx *NodeContext) (bool, error)

	// Supports 检查是否支持指定的条件类型
	// conditionType: 条件类型(如 "numeric", "string", "enum", "custom", "composite")
	// 返回: 是否支持该条件类型
	Supports(conditionType string) bool
}

