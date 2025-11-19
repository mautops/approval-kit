package node

import (
	"encoding/json"
	"fmt"

	"github.com/mautops/approval-kit/internal/template"
)

// ConditionNodeExecutor 条件节点执行器
type ConditionNodeExecutor struct {
	registry *ConditionEvaluatorRegistry
}

// NewConditionNodeExecutor 创建新的条件节点执行器
func NewConditionNodeExecutor() NodeExecutor {
	return &ConditionNodeExecutor{
		registry: NewConditionEvaluatorRegistry(),
	}
}

// NodeType 返回节点类型(实现 NodeExecutor 接口)
func (e *ConditionNodeExecutor) NodeType() template.NodeType {
	return template.NodeTypeCondition
}

// Execute 执行条件节点逻辑(实现 NodeExecutor 接口)
func (e *ConditionNodeExecutor) Execute(ctx *NodeContext) (*NodeResult, error) {
	// 1. 获取节点配置
	config, ok := ctx.Node.Config.(*ConditionNodeConfig)
	if !ok {
		return nil, fmt.Errorf("invalid node config type for condition node")
	}

	// 2. 验证配置
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("condition node config validation failed: %w", err)
	}

	// 3. 获取条件评估器
	evaluator := e.registry.GetEvaluator(config.Condition.Type)
	if evaluator == nil {
		return nil, fmt.Errorf("no evaluator found for condition type: %q", config.Condition.Type)
	}

	// 4. 评估条件
	result, err := evaluator.Evaluate(config.Condition, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate condition: %w", err)
	}

	// 5. 根据条件结果决定下一个节点
	var nextNodeID string
	if result {
		nextNodeID = config.TrueNodeID
	} else {
		nextNodeID = config.FalseNodeID
	}

	// 6. 生成输出数据
	output := json.RawMessage(`{"condition_result": ` + boolToJSON(result) + `, "next_node_id": "` + nextNodeID + `"}`)

	return &NodeResult{
		NextNodeID: nextNodeID,
		Output:     output,
		Events:     []Event{},
	}, nil
}

// boolToJSON 将 bool 转换为 JSON 字符串
func boolToJSON(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

