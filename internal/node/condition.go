package node

import "fmt"

// Condition 条件定义
// 用于条件节点,根据条件结果决定流程走向
type Condition struct {
	// Type 条件类型
	// 支持的类型: "numeric"(数值比较), "string"(字符串匹配), "enum"(枚举判断), "custom"(自定义函数), "composite"(组合条件)
	Type string

	// Config 条件配置(根据类型不同而不同)
	Config ConditionConfig
}

// Validate 验证条件的有效性
func (c *Condition) Validate() error {
	if c.Type == "" {
		return fmt.Errorf("Condition.Type is required")
	}

	if c.Config == nil {
		return fmt.Errorf("Condition.Config is required")
	}

	// 验证配置类型与条件类型一致
	if c.Config.ConditionType() != c.Type {
		return fmt.Errorf("Condition.Config.ConditionType() = %q, want %q", c.Config.ConditionType(), c.Type)
	}

	return nil
}

// ConditionConfig 条件配置接口
// 不同类型的条件有不同的配置实现
type ConditionConfig interface {
	// ConditionType 返回条件类型
	ConditionType() string
}

