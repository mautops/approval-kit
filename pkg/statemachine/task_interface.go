package statemachine

import (
	"time"

	"github.com/mautops/approval-kit/pkg/types"
	internalSM "github.com/mautops/approval-kit/internal/statemachine"
)

// TransitionableTask 定义状态转换所需的任务接口
// 这个接口抽象了状态转换操作所需的最小方法集
// 与 internal/statemachine.TransitionableTask 接口定义完全一致,但位于 pkg 目录,可以被外部导入
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
// 与 internal/statemachine.StateChange 结构相同,但位于 pkg 目录,可以被外部导入
type StateChange = internalSM.StateChange

// StateChangeFromInternal 将 internal.StateChange 转换为 pkg.StateChange
func StateChangeFromInternal(sc *internalSM.StateChange) *StateChange {
	return (*StateChange)(sc)
}

// StateChangeToInternal 将 pkg.StateChange 转换为 internal.StateChange
func StateChangeToInternal(sc *StateChange) *internalSM.StateChange {
	return (*internalSM.StateChange)(sc)
}

