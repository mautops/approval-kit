package statemachine

import (
	"github.com/mautops/approval-kit/pkg/types"
)

// StateMachine 状态机接口
// 管理审批任务的状态流转,确保状态转换的合法性和一致性
// 与 internal/statemachine.StateMachine 接口定义完全一致,但位于 pkg 目录,可以被外部导入
type StateMachine interface {
	// CanTransition 检查是否允许状态转换
	// from: 源状态
	// to: 目标状态
	// 返回: true 表示允许转换,false 表示不允许
	CanTransition(from types.TaskState, to types.TaskState) bool

	// Transition 执行状态转换
	// task: 任务对象(必须实现 TransitionableTask 接口)
	// to: 目标状态
	// reason: 转换原因
	// 返回: 新的任务对象和错误信息
	// 注意: 返回新任务对象,原任务对象不会被修改(不可变实现)
	Transition(task TransitionableTask, to types.TaskState, reason string) (TransitionableTask, error)

	// GetValidTransitions 获取当前状态的有效转换
	// state: 当前状态
	// 返回: 允许转换到的状态列表
	GetValidTransitions(state types.TaskState) []types.TaskState
}

