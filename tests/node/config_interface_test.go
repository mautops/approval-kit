package node_test

import (
	"testing"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/template"
)

// TestNodeConfigInterface 验证 NodeConfig 接口定义
func TestNodeConfigInterface(t *testing.T) {
	// 验证接口类型存在
	var config node.NodeConfig
	if config != nil {
		_ = config
	}
}

// TestNodeConfigMethods 验证 NodeConfig 接口方法签名
func TestNodeConfigMethods(t *testing.T) {
	// 验证接口包含所有必需的方法
	// 通过编译时检查,如果方法不存在会编译失败
	var _ node.NodeConfig = (*nodeConfigImpl)(nil)
}

// nodeConfigImpl 用于测试接口方法签名的实现
type nodeConfigImpl struct{}

func (n *nodeConfigImpl) NodeType() template.NodeType {
	return template.NodeTypeStart
}

func (n *nodeConfigImpl) Validate() error {
	return nil
}

