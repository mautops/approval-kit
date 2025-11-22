package event

import (
	internalEvent "github.com/mautops/approval-kit/internal/event"
)

// EventType 事件类型
// 与 internal/event.EventType 类型相同,但位于 pkg 目录,可以被外部导入
type EventType = internalEvent.EventType

// 事件类型常量
const (
	// EventTypeTaskCreated 任务创建事件
	EventTypeTaskCreated EventType = internalEvent.EventTypeTaskCreated

	// EventTypeTaskSubmitted 任务提交事件
	EventTypeTaskSubmitted EventType = internalEvent.EventTypeTaskSubmitted

	// EventTypeNodeActivated 节点激活事件
	EventTypeNodeActivated EventType = internalEvent.EventTypeNodeActivated

	// EventTypeApprovalOp 审批操作事件
	EventTypeApprovalOp EventType = internalEvent.EventTypeApprovalOp

	// EventTypeTaskApproved 任务通过事件
	EventTypeTaskApproved EventType = internalEvent.EventTypeTaskApproved

	// EventTypeTaskRejected 任务拒绝事件
	EventTypeTaskRejected EventType = internalEvent.EventTypeTaskRejected

	// EventTypeTaskTimeout 任务超时事件
	EventTypeTaskTimeout EventType = internalEvent.EventTypeTaskTimeout

	// EventTypeTaskCancelled 任务取消事件
	EventTypeTaskCancelled EventType = internalEvent.EventTypeTaskCancelled

	// EventTypeTaskWithdrawn 任务撤回事件
	EventTypeTaskWithdrawn EventType = internalEvent.EventTypeTaskWithdrawn

	// EventTypeTaskPaused 任务暂停事件
	EventTypeTaskPaused EventType = internalEvent.EventTypeTaskPaused

	// EventTypeTaskResumed 任务恢复事件
	EventTypeTaskResumed EventType = internalEvent.EventTypeTaskResumed

	// EventTypeTaskRollback 任务回退事件
	EventTypeTaskRollback EventType = internalEvent.EventTypeTaskRollback

	// EventTypeApproverReplaced 审批人替换事件
	EventTypeApproverReplaced EventType = internalEvent.EventTypeApproverReplaced

	// EventTypeNodeCompleted 节点完成事件
	EventTypeNodeCompleted EventType = internalEvent.EventTypeNodeCompleted
)

// Event 事件定义
// 用于通知上层业务系统审批流程中的关键事件
// 与 internal/event.Event 结构相同,但位于 pkg 目录,可以被外部导入
type Event = internalEvent.Event

// TaskInfo 任务信息
// 与 internal/event.TaskInfo 结构相同,但位于 pkg 目录,可以被外部导入
type TaskInfo = internalEvent.TaskInfo

// NodeInfo 节点信息
// 与 internal/event.NodeInfo 结构相同,但位于 pkg 目录,可以被外部导入
type NodeInfo = internalEvent.NodeInfo

// ApprovalInfo 审批信息
// 与 internal/event.ApprovalInfo 结构相同,但位于 pkg 目录,可以被外部导入
type ApprovalInfo = internalEvent.ApprovalInfo

// BusinessInfo 业务信息
// 与 internal/event.BusinessInfo 结构相同,但位于 pkg 目录,可以被外部导入
type BusinessInfo = internalEvent.BusinessInfo

// EventFromInternal 将 internal.Event 转换为 pkg.Event
func EventFromInternal(e *internalEvent.Event) *Event {
	return (*Event)(e)
}

// EventToInternal 将 pkg.Event 转换为 internal.Event
func EventToInternal(e *Event) *internalEvent.Event {
	return (*internalEvent.Event)(e)
}

