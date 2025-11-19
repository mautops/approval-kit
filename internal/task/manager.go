package task

import (
	"encoding/json"
)

// TaskManager 任务管理接口
// 负责审批任务的创建、查询、提交、审批等操作
type TaskManager interface {
	// Create 基于模板创建审批任务实例
	// templateID: 模板 ID
	// businessID: 关联的业务 ID
	// params: 任务参数(JSON 格式),用于条件判断和动态审批人获取
	// 返回: 任务对象和错误信息
	Create(templateID string, businessID string, params json.RawMessage) (*Task, error)

	// Get 获取审批任务详情
	// id: 任务 ID
	// 返回: 任务对象和错误信息
	Get(id string) (*Task, error)

	// Submit 提交任务进入审批流程
	// id: 任务 ID
	// 返回: 错误信息
	// 注意: 提交会触发状态转换,从 pending 转换为 submitted
	Submit(id string) error

	// Approve 审批人进行同意操作
	// id: 任务 ID
	// nodeID: 节点 ID
	// approver: 审批人 ID
	// comment: 审批意见
	// 返回: 错误信息
	Approve(id string, nodeID string, approver string, comment string) error

	// ApproveWithAttachments 审批人进行同意操作(带附件)
	// id: 任务 ID
	// nodeID: 节点 ID
	// approver: 审批人 ID
	// comment: 审批意见
	// attachments: 附件列表
	// 返回: 错误信息
	ApproveWithAttachments(id string, nodeID string, approver string, comment string, attachments []string) error

	// Reject 审批人进行拒绝操作
	// id: 任务 ID
	// nodeID: 节点 ID
	// approver: 审批人 ID
	// comment: 审批意见
	// 返回: 错误信息
	Reject(id string, nodeID string, approver string, comment string) error

	// RejectWithAttachments 审批人进行拒绝操作(带附件)
	// id: 任务 ID
	// nodeID: 节点 ID
	// approver: 审批人 ID
	// comment: 审批意见
	// attachments: 附件列表
	// 返回: 错误信息
	RejectWithAttachments(id string, nodeID string, approver string, comment string, attachments []string) error

	// Cancel 取消任务
	// id: 任务 ID
	// reason: 取消原因
	// 返回: 错误信息
	Cancel(id string, reason string) error

	// Withdraw 撤回任务
	// id: 任务 ID
	// reason: 撤回原因
	// 返回: 错误信息
	// 注意: 撤回会将任务从 submitted 或 approving 状态撤回回 pending 状态
	// 如果任务已有审批记录,不允许撤回
	Withdraw(id string, reason string) error

	// Transfer 转交审批
	// id: 任务 ID
	// nodeID: 节点 ID
	// fromApprover: 原审批人 ID
	// toApprover: 新审批人 ID
	// reason: 转交原因
	// 返回: 错误信息
	// 注意: 转交需要节点配置允许转交,且原审批人必须是当前审批人
	Transfer(id string, nodeID string, fromApprover string, toApprover string, reason string) error

	// AddApprover 加签
	// id: 任务 ID
	// nodeID: 节点 ID
	// approver: 新审批人 ID
	// reason: 加签原因
	// 返回: 错误信息
	// 注意: 加签需要节点配置允许加签
	AddApprover(id string, nodeID string, approver string, reason string) error

	// RemoveApprover 减签
	// id: 任务 ID
	// nodeID: 节点 ID
	// approver: 要移除的审批人 ID
	// reason: 减签原因
	// 返回: 错误信息
	// 注意: 减签需要节点配置允许减签,且审批人必须在审批人列表中
	RemoveApprover(id string, nodeID string, approver string, reason string) error

	// Query 查询任务列表
	// filter: 查询过滤器
	// 返回: 任务列表和错误信息
	Query(filter *TaskFilter) ([]*Task, error)

	// HandleTimeout 处理任务超时
	// id: 任务 ID
	// 返回: 错误信息
	// 注意: 如果任务已超时,将任务状态转换为 timeout
	HandleTimeout(id string) error
}

