package node

import (
	"github.com/mautops/approval-kit/internal/task"
)

// ApprovalModeHandler 审批模式处理器接口
// 使用策略模式,将不同审批模式的逻辑封装到独立的处理器中
type ApprovalModeHandler interface {
	// CheckCompletion 检查审批是否完成
	// approvers: 审批人列表
	// approvals: 审批结果映射(节点 ID -> 审批人 -> 审批结果)
	// config: 审批节点配置
	// 返回: 是否完成和审批结果
	CheckCompletion(approvers []string, approvals map[string]*task.Approval, config *ApprovalNodeConfig) (bool, *ApprovalResult)

	// Mode 返回处理器对应的审批模式
	Mode() ApprovalMode
}

// ApprovalModeHandlerRegistry 审批模式处理器注册表
// 用于管理和查找不同审批模式的处理器
type ApprovalModeHandlerRegistry interface {
	// GetHandler 根据审批模式获取对应的处理器
	// mode: 审批模式
	// 返回: 处理器实例,如果模式不存在则返回 nil
	GetHandler(mode ApprovalMode) ApprovalModeHandler

	// RegisterHandler 注册审批模式处理器
	// mode: 审批模式
	// handler: 处理器实例
	RegisterHandler(mode ApprovalMode, handler ApprovalModeHandler)
}

// approvalModeHandlerRegistry 审批模式处理器注册表实现
type approvalModeHandlerRegistry struct {
	handlers map[ApprovalMode]ApprovalModeHandler
}

// NewApprovalModeHandlerRegistry 创建新的审批模式处理器注册表
func NewApprovalModeHandlerRegistry() ApprovalModeHandlerRegistry {
	registry := &approvalModeHandlerRegistry{
		handlers: make(map[ApprovalMode]ApprovalModeHandler),
	}

	// 注册默认处理器
	registry.RegisterHandler(ApprovalModeSingle, &singleModeHandler{})
	registry.RegisterHandler(ApprovalModeUnanimous, &unanimousModeHandler{})
	registry.RegisterHandler(ApprovalModeOr, &orModeHandler{})
	registry.RegisterHandler(ApprovalModeProportional, &proportionalModeHandler{})
	registry.RegisterHandler(ApprovalModeSequential, &sequentialModeHandler{})

	return registry
}

// GetHandler 根据审批模式获取对应的处理器
func (r *approvalModeHandlerRegistry) GetHandler(mode ApprovalMode) ApprovalModeHandler {
	return r.handlers[mode]
}

// RegisterHandler 注册审批模式处理器
func (r *approvalModeHandlerRegistry) RegisterHandler(mode ApprovalMode, handler ApprovalModeHandler) {
	r.handlers[mode] = handler
}

