package template

import "time"

// Node 表示审批节点
// 节点是审批流程的基本单元,不同类型的节点有不同的执行逻辑
type Node struct {
	ID    string   // 节点 ID
	Name  string   // 节点名称
	Type  NodeType // 节点类型
	Order int      // 节点顺序

	// 节点配置 (根据节点类型不同而不同)
	// 将在后续任务中完善
	Config NodeConfig // 节点配置接口
}

// NodeConfig 节点配置接口
// 不同类型的节点有不同的配置实现
type NodeConfig interface {
	// NodeType 返回节点类型
	NodeType() NodeType
	// Validate 验证配置的有效性
	Validate() error
}

// ApprovalNodeConfigAccessor 审批节点配置访问接口
// 用于在不导入 node 包的情况下访问审批节点配置的属性
// 避免循环依赖
type ApprovalNodeConfigAccessor interface {
	NodeConfig
	// RequireComment 返回是否必填审批意见
	RequireComment() bool
	// RequireAttachments 返回是否要求附件
	RequireAttachments() bool
	// GetTimeout 返回超时时间配置
	GetTimeout() *time.Duration
	// GetPermissions 返回操作权限配置
	GetPermissions() interface{} // 返回 OperationPermissions,但使用 interface{} 避免循环依赖
	// GetRejectBehavior 返回拒绝后行为配置
	GetRejectBehavior() string // 返回 "terminate", "rollback", "jump" 等
	// GetRejectTargetNode 返回拒绝后跳转目标节点(仅当 RejectBehavior 为 "jump" 时有效)
	GetRejectTargetNode() string
}

// OperationPermissionsAccessor 操作权限访问接口
// 用于在不导入 node 包的情况下访问操作权限配置
type OperationPermissionsAccessor interface {
	// AllowTransfer 返回是否允许转交
	AllowTransfer() bool
	// AllowAddApprover 返回是否允许加签
	AllowAddApprover() bool
	// AllowRemoveApprover 返回是否允许减签
	AllowRemoveApprover() bool
}

