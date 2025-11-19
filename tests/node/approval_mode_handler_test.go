package node_test

import (
	"testing"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
)

// TestApprovalModeHandlerInterface 验证 ApprovalModeHandler 接口定义
func TestApprovalModeHandlerInterface(t *testing.T) {
	// 验证接口类型存在
	var handler node.ApprovalModeHandler
	if handler != nil {
		_ = handler
	}
}

// TestApprovalModeHandlerMethods 验证 ApprovalModeHandler 接口方法签名
func TestApprovalModeHandlerMethods(t *testing.T) {
	// 验证接口包含所有必需的方法
	// 通过编译时检查,如果方法不存在会编译失败
	var _ node.ApprovalModeHandler = (*approvalModeHandlerImpl)(nil)
}

// approvalModeHandlerImpl 用于测试接口方法签名的实现
type approvalModeHandlerImpl struct{}

func (h *approvalModeHandlerImpl) CheckCompletion(approvers []string, approvals map[string]*task.Approval, config *node.ApprovalNodeConfig) (bool, *node.ApprovalResult) {
	return false, nil
}

func (h *approvalModeHandlerImpl) Mode() node.ApprovalMode {
	return node.ApprovalModeSingle
}

// TestApprovalModeHandlerRegistry 测试审批模式处理器注册表
func TestApprovalModeHandlerRegistry(t *testing.T) {
	// 验证可以通过模式获取处理器
	registry := node.NewApprovalModeHandlerRegistry()

	// 测试获取各种模式的处理器
	handler := registry.GetHandler(node.ApprovalModeSingle)
	if handler == nil {
		t.Error("Registry should return handler for ApprovalModeSingle")
	}

	handler = registry.GetHandler(node.ApprovalModeUnanimous)
	if handler == nil {
		t.Error("Registry should return handler for ApprovalModeUnanimous")
	}

	handler = registry.GetHandler(node.ApprovalModeOr)
	if handler == nil {
		t.Error("Registry should return handler for ApprovalModeOr")
	}

	handler = registry.GetHandler(node.ApprovalModeProportional)
	if handler == nil {
		t.Error("Registry should return handler for ApprovalModeProportional")
	}

	handler = registry.GetHandler(node.ApprovalModeSequential)
	if handler == nil {
		t.Error("Registry should return handler for ApprovalModeSequential")
	}
}

