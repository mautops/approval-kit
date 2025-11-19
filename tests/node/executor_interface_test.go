package node_test

import (
	"testing"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/template"
)

// TestNodeExecutorInterface 验证 NodeExecutor 接口定义
func TestNodeExecutorInterface(t *testing.T) {
	// 验证接口类型存在
	var executor node.NodeExecutor
	if executor != nil {
		_ = executor
	}
}

// TestNodeExecutorMethods 验证 NodeExecutor 接口方法签名
func TestNodeExecutorMethods(t *testing.T) {
	// 验证接口包含所有必需的方法
	// 通过编译时检查,如果方法不存在会编译失败
	var _ node.NodeExecutor = (*nodeExecutorImpl)(nil)
}

// nodeExecutorImpl 用于测试接口方法签名的实现
type nodeExecutorImpl struct{}

func (n *nodeExecutorImpl) Execute(ctx *node.NodeContext) (*node.NodeResult, error) {
	return nil, nil
}

func (n *nodeExecutorImpl) NodeType() template.NodeType {
	return template.NodeTypeStart
}
