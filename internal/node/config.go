package node

import (
	"github.com/mautops/approval-kit/internal/template"
)

// NodeConfig 节点配置接口
// 不同类型的节点有不同的配置实现
type NodeConfig interface {
	// NodeType 返回节点类型
	// 用于标识配置对应的节点类型
	NodeType() template.NodeType

	// Validate 验证配置的有效性
	// 返回错误信息,如果配置无效
	Validate() error
}

