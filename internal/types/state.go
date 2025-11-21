package types

// TaskState 表示审批任务的状态
// 状态转换必须通过状态机进行,确保合法性
type TaskState string

const (
	// TaskStatePending 待审批: 任务已创建,等待审批
	TaskStatePending TaskState = "pending"

	// TaskStateSubmitted 已提交: 任务已提交进入审批流程
	TaskStateSubmitted TaskState = "submitted"

	// TaskStateApproving 审批中: 任务正在审批流程中
	TaskStateApproving TaskState = "approving"

	// TaskStateApproved 已通过: 任务审批通过
	TaskStateApproved TaskState = "approved"

	// TaskStateRejected 已拒绝: 任务审批被拒绝
	TaskStateRejected TaskState = "rejected"

	// TaskStateCancelled 已取消: 任务被取消
	TaskStateCancelled TaskState = "cancelled"

	// TaskStateTimeout 已超时: 任务审批超时
	TaskStateTimeout TaskState = "timeout"

	// TaskStatePaused 已暂停: 任务被暂停,可以稍后恢复
	TaskStatePaused TaskState = "paused"
)

