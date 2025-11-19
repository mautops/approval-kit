package task

import "github.com/mautops/approval-kit/internal/types"

// TaskState 是 types.TaskState 的类型别名,保持向后兼容
type TaskState = types.TaskState

// 导出常量,保持向后兼容
const (
	TaskStatePending   = types.TaskStatePending
	TaskStateSubmitted = types.TaskStateSubmitted
	TaskStateApproving = types.TaskStateApproving
	TaskStateApproved  = types.TaskStateApproved
	TaskStateRejected  = types.TaskStateRejected
	TaskStateCancelled = types.TaskStateCancelled
	TaskStateTimeout   = types.TaskStateTimeout
)

