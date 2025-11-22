package task

import (
	"encoding/json"
)

// TaskManager 任务管理接口
// 负责审批任务的创建、查询、提交、审批等操作
// 与 internal/task.TaskManager 接口定义完全一致,但位于 pkg 目录,可以被外部导入
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

	// Pause 暂停任务
	// id: 任务 ID
	// reason: 暂停原因
	// 返回: 错误信息
	// 注意: 只有 pending、submitted、approving 状态可以暂停
	// 暂停时会记录暂停前的状态,用于恢复时恢复到正确状态
	Pause(id string, reason string) error

	// Resume 恢复任务
	// id: 任务 ID
	// reason: 恢复原因
	// 返回: 错误信息
	// 注意: 只有 paused 状态可以恢复
	// 恢复时会恢复到暂停前的状态(pending、submitted 或 approving)
	Resume(id string, reason string) error

	// RollbackToNode 回退到指定节点
	// id: 任务 ID
	// nodeID: 目标节点 ID
	// reason: 回退原因
	// 返回: 错误信息
	// 注意: 只能回退到已完成的节点
	// 回退时会清理回退节点之后的审批记录和状态
	RollbackToNode(id string, nodeID string, reason string) error

	// ReplaceApprover 替换审批人
	// id: 任务 ID
	// nodeID: 节点 ID
	// oldApprover: 原审批人 ID
	// newApprover: 新审批人 ID
	// reason: 替换原因
	// 返回: 错误信息
	// 注意: 只能替换尚未审批的审批人
	// 替换后会保留原审批人的审批记录(如果有),新审批人可以继续审批
	ReplaceApprover(id string, nodeID string, oldApprover string, newApprover string, reason string) error
}

