package types

import (
	internalTypes "github.com/mautops/approval-kit/internal/types"
)

// TaskState 表示审批任务的状态
// 状态转换必须通过状态机进行,确保合法性
// 与 internal/types.TaskState 类型相同,但位于 pkg 目录,可以被外部导入
type TaskState = internalTypes.TaskState

// 任务状态常量
const (
	// TaskStatePending 待审批: 任务已创建,等待审批
	TaskStatePending TaskState = internalTypes.TaskStatePending

	// TaskStateSubmitted 已提交: 任务已提交进入审批流程
	TaskStateSubmitted TaskState = internalTypes.TaskStateSubmitted

	// TaskStateApproving 审批中: 任务正在审批流程中
	TaskStateApproving TaskState = internalTypes.TaskStateApproving

	// TaskStateApproved 已通过: 任务审批通过
	TaskStateApproved TaskState = internalTypes.TaskStateApproved

	// TaskStateRejected 已拒绝: 任务审批被拒绝
	TaskStateRejected TaskState = internalTypes.TaskStateRejected

	// TaskStateCancelled 已取消: 任务被取消
	TaskStateCancelled TaskState = internalTypes.TaskStateCancelled

	// TaskStateTimeout 已超时: 任务审批超时
	TaskStateTimeout TaskState = internalTypes.TaskStateTimeout

	// TaskStatePaused 已暂停: 任务被暂停,可以稍后恢复
	TaskStatePaused TaskState = internalTypes.TaskStatePaused
)

