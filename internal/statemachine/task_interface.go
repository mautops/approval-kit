package statemachine

import (
	"time"

	"github.com/mautops/approval-kit/internal/types"
)

// TransitionableTask 定义状态转换所需的任务接口
// 这个接口抽象了状态转换操作所需的最小方法集
type TransitionableTask interface {
	// GetState 获取当前状态
	GetState() types.TaskState

	// SetState 设置新状态(用于不可变实现,应该返回新对象)
	SetState(state types.TaskState)

	// GetUpdatedAt 获取更新时间
	GetUpdatedAt() time.Time

	// SetUpdatedAt 设置更新时间
	SetUpdatedAt(t time.Time)

	// GetStateHistory 获取状态变更历史
	GetStateHistory() []*StateChange

	// AddStateChange 添加状态变更记录
	AddStateChange(change *StateChange)

	// Clone 创建任务副本(用于不可变实现)
	Clone() TransitionableTask
}

// StateChange 状态变更记录
type StateChange struct {
	From    types.TaskState // 源状态
	To      types.TaskState // 目标状态
	Reason  string          // 转换原因
	Time    time.Time       // 转换时间
}

