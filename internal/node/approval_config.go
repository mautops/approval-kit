package node

import (
	"fmt"
	"time"

	"github.com/mautops/approval-kit/internal/errors"
	"github.com/mautops/approval-kit/internal/template"
)

// ApprovalNodeConfig 审批节点配置
// 实现 NodeConfig 接口
type ApprovalNodeConfig struct {
	// 审批模式
	Mode ApprovalMode

	// 审批人配置
	ApproverConfig ApproverConfig

	// 其他配置
	Timeout         *time.Duration      // 超时时间
	RejectBehavior  RejectBehavior      // 拒绝后行为
	RejectTargetNode string             // 拒绝后跳转目标节点(仅当 RejectBehavior 为 RejectBehaviorJump 时有效)
	Permissions     OperationPermissions // 操作权限
	RequireCommentField  bool                // 是否必填审批意见
	RequireAttachmentsField bool             // 是否要求附件

	// 比例会签配置(仅用于 ApprovalModeProportional)
	ProportionalThreshold *ProportionalThreshold // 比例阈值
}

// ProportionalThreshold 比例会签阈值配置
type ProportionalThreshold struct {
	// Required 需要同意的审批人数量
	Required int

	// Total 总审批人数量
	Total int
}

// NodeType 返回节点类型(实现 NodeConfig 接口)
func (c *ApprovalNodeConfig) NodeType() template.NodeType {
	return template.NodeTypeApproval
}

// Validate 验证配置的有效性(实现 NodeConfig 接口)
func (c *ApprovalNodeConfig) Validate() error {
	// 验证审批模式
	if c.Mode == "" {
		return fmt.Errorf("%w: approval mode is required", errors.ErrInvalidTemplate)
	}

	// 验证审批模式值
	validModes := map[ApprovalMode]bool{
		ApprovalModeSingle:       true,
		ApprovalModeUnanimous:    true,
		ApprovalModeOr:           true,
		ApprovalModeProportional: true,
		ApprovalModeSequential:   true,
	}
	if !validModes[c.Mode] {
		return fmt.Errorf("%w: invalid approval mode: %q", errors.ErrInvalidTemplate, c.Mode)
	}

	// 验证审批人配置
	if c.ApproverConfig == nil {
		return fmt.Errorf("%w: approver config is required", errors.ErrInvalidTemplate)
	}

	// 验证比例会签配置(如果使用比例会签模式)
	if c.Mode == ApprovalModeProportional {
		if c.ProportionalThreshold == nil {
			return fmt.Errorf("%w: proportional threshold is required for proportional mode", errors.ErrInvalidTemplate)
		}
		if c.ProportionalThreshold.Required <= 0 {
			return fmt.Errorf("%w: proportional threshold required must be greater than 0", errors.ErrInvalidTemplate)
		}
		if c.ProportionalThreshold.Total <= 0 {
			return fmt.Errorf("%w: proportional threshold total must be greater than 0", errors.ErrInvalidTemplate)
		}
		if c.ProportionalThreshold.Required > c.ProportionalThreshold.Total {
			return fmt.Errorf("%w: proportional threshold required cannot be greater than total", errors.ErrInvalidTemplate)
		}
	}

	// 验证拒绝后跳转配置(如果使用跳转行为)
	if c.RejectBehavior == RejectBehaviorJump {
		if c.RejectTargetNode == "" {
			return fmt.Errorf("%w: reject target node is required when reject behavior is jump", errors.ErrInvalidTemplate)
		}
	}

	// 验证超时配置(如果设置了超时,必须大于 0)
	if c.Timeout != nil && *c.Timeout <= 0 {
		return fmt.Errorf("%w: timeout must be greater than 0", errors.ErrInvalidTemplate)
	}

	return nil
}

// RequireComment 返回是否必填审批意见(实现 ApprovalNodeConfigAccessor 接口)
func (c *ApprovalNodeConfig) RequireComment() bool {
	return c.RequireCommentField
}

// RequireAttachments 返回是否要求附件(实现 ApprovalNodeConfigAccessor 接口)
func (c *ApprovalNodeConfig) RequireAttachments() bool {
	return c.RequireAttachmentsField
}

// GetTimeout 返回超时时间配置(实现 ApprovalNodeConfigAccessor 接口)
func (c *ApprovalNodeConfig) GetTimeout() *time.Duration {
	return c.Timeout
}

// GetPermissions 返回操作权限配置(实现 ApprovalNodeConfigAccessor 接口)
func (c *ApprovalNodeConfig) GetPermissions() interface{} {
	return &permissionsAccessor{perms: c.Permissions}
}

// GetRejectBehavior 返回拒绝后行为配置(实现 ApprovalNodeConfigAccessor 接口)
func (c *ApprovalNodeConfig) GetRejectBehavior() string {
	return string(c.RejectBehavior)
}

// GetRejectTargetNode 返回拒绝后跳转目标节点(实现 ApprovalNodeConfigAccessor 接口)
func (c *ApprovalNodeConfig) GetRejectTargetNode() string {
	return c.RejectTargetNode
}

// permissionsAccessor 权限访问器,实现 OperationPermissionsAccessor 接口
type permissionsAccessor struct {
	perms OperationPermissions
}

// AllowTransfer 返回是否允许转交
func (p *permissionsAccessor) AllowTransfer() bool {
	return p.perms.AllowTransfer
}

// AllowAddApprover 返回是否允许加签
func (p *permissionsAccessor) AllowAddApprover() bool {
	return p.perms.AllowAddApprover
}

// AllowRemoveApprover 返回是否允许减签
func (p *permissionsAccessor) AllowRemoveApprover() bool {
	return p.perms.AllowRemoveApprover
}

// RejectBehavior 拒绝后行为
type RejectBehavior string

const (
	// RejectBehaviorTerminate 拒绝后终止流程
	RejectBehaviorTerminate RejectBehavior = "terminate"

	// RejectBehaviorRollback 拒绝后回退到上一节点
	RejectBehaviorRollback RejectBehavior = "rollback"

	// RejectBehaviorJump 拒绝后跳转到指定节点
	RejectBehaviorJump RejectBehavior = "jump"
)

// OperationPermissions 操作权限配置
type OperationPermissions struct {
	// AllowTransfer 允许转交审批
	AllowTransfer bool

	// AllowAddApprover 允许加签
	AllowAddApprover bool

	// AllowRemoveApprover 允许减签
	AllowRemoveApprover bool
}

// ApproverConfig 审批人配置接口
type ApproverConfig interface {
	// GetApprovers 获取审批人列表
	// ctx: 节点执行上下文
	// 返回: 审批人 ID 列表和错误信息
	GetApprovers(ctx *NodeContext) ([]string, error)

	// GetTiming 返回获取时机
	// 返回: 获取时机(任务创建时或节点激活时)
	GetTiming() ApproverTiming
}

// ApproverTiming 审批人获取时机
type ApproverTiming string

const (
	// ApproverTimingOnCreate 任务创建时获取
	ApproverTimingOnCreate ApproverTiming = "on_create"

	// ApproverTimingOnActivate 节点激活时获取
	ApproverTimingOnActivate ApproverTiming = "on_activate"
)

// FixedApproverConfig 固定审批人配置
// 在模板中预设审批人列表
type FixedApproverConfig struct {
	Approvers []string // 审批人列表(用户 ID)
}

// GetApprovers 获取审批人列表(实现 ApproverConfig 接口)
func (c *FixedApproverConfig) GetApprovers(ctx *NodeContext) ([]string, error) {
	// 返回固定审批人列表的副本
	result := make([]string, len(c.Approvers))
	copy(result, c.Approvers)
	return result, nil
}

// GetTiming 返回获取时机(实现 ApproverConfig 接口)
func (c *FixedApproverConfig) GetTiming() ApproverTiming {
	// 固定审批人在任务创建时即可确定,但为了灵活性,默认在节点激活时获取
	// 如果需要任务创建时获取,可以在配置中指定
	return ApproverTimingOnActivate
}

