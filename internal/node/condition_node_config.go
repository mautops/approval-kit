package node

import (
	"fmt"

	"github.com/mautops/approval-kit/internal/template"
)

// ConditionNodeConfig 条件节点配置
// 实现 NodeConfig 接口
type ConditionNodeConfig struct {
	// Condition 条件定义
	Condition *Condition

	// TrueNodeID 条件为 true 时跳转的节点 ID
	TrueNodeID string

	// FalseNodeID 条件为 false 时跳转的节点 ID
	FalseNodeID string
}

// NodeType 返回节点类型(实现 NodeConfig 接口)
func (c *ConditionNodeConfig) NodeType() template.NodeType {
	return template.NodeTypeCondition
}

// Validate 验证配置的有效性(实现 NodeConfig 接口)
func (c *ConditionNodeConfig) Validate() error {
	if c.Condition == nil {
		return fmt.Errorf("ConditionNodeConfig.Condition is required")
	}

	if err := c.Condition.Validate(); err != nil {
		return fmt.Errorf("ConditionNodeConfig.Condition validation failed: %w", err)
	}

	if c.TrueNodeID == "" {
		return fmt.Errorf("ConditionNodeConfig.TrueNodeID is required")
	}

	if c.FalseNodeID == "" {
		return fmt.Errorf("ConditionNodeConfig.FalseNodeID is required")
	}

	return nil
}

