package event

import (
	"encoding/json"
	"time"
)

// EventType 事件类型
type EventType string

const (
	// EventTypeTaskCreated 任务创建事件
	EventTypeTaskCreated EventType = "task_created"

	// EventTypeTaskSubmitted 任务提交事件
	EventTypeTaskSubmitted EventType = "task_submitted"

	// EventTypeNodeActivated 节点激活事件
	EventTypeNodeActivated EventType = "node_activated"

	// EventTypeApprovalOp 审批操作事件
	EventTypeApprovalOp EventType = "approval_operation"

	// EventTypeTaskApproved 任务通过事件
	EventTypeTaskApproved EventType = "task_approved"

	// EventTypeTaskRejected 任务拒绝事件
	EventTypeTaskRejected EventType = "task_rejected"

	// EventTypeTaskTimeout 任务超时事件
	EventTypeTaskTimeout EventType = "task_timeout"

	// EventTypeTaskCancelled 任务取消事件
	EventTypeTaskCancelled EventType = "task_cancelled"

	// EventTypeTaskWithdrawn 任务撤回事件
	EventTypeTaskWithdrawn EventType = "task_withdrawn"

	// EventTypeNodeCompleted 节点完成事件
	EventTypeNodeCompleted EventType = "node_completed"
)

// Event 事件定义
// 用于通知上层业务系统审批流程中的关键事件
type Event struct {
	// ID 事件 ID(用于幂等性保证)
	ID string

	// Type 事件类型
	Type EventType

	// Time 事件时间
	Time time.Time

	// Task 任务信息
	Task *TaskInfo

	// Node 节点信息(所有事件都包含)
	Node *NodeInfo

	// Approval 审批信息(如适用)
	Approval *ApprovalInfo

	// Business 业务信息
	Business *BusinessInfo
}

// TaskInfo 任务信息
type TaskInfo struct {
	// ID 任务 ID
	ID string

	// TemplateID 模板 ID
	TemplateID string

	// BusinessID 业务 ID
	BusinessID string

	// State 任务状态
	State string
}

// NodeInfo 节点信息
type NodeInfo struct {
	// ID 节点 ID
	ID string

	// Name 节点名称
	Name string

	// Type 节点类型
	Type string
}

// ApprovalInfo 审批信息
type ApprovalInfo struct {
	// NodeID 节点 ID
	NodeID string

	// Approver 审批人
	Approver string

	// Result 审批结果(approve/reject/transfer)
	Result string

	// Comment 审批意见
	Comment string
}

// BusinessInfo 业务信息
type BusinessInfo struct {
	// ID 业务 ID
	ID string

	// Data 业务数据(JSON 格式)
	Data json.RawMessage
}

