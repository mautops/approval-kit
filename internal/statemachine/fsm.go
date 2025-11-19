package statemachine

import (
	"time"

	"github.com/mautops/approval-kit/internal/errors"
	"github.com/mautops/approval-kit/internal/types"
)

// StateMachine 状态机接口
// 管理审批任务的状态流转,确保状态转换的合法性和一致性
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

// fsm 状态机实现
type fsm struct {
	transitions map[types.TaskState][]types.TaskState
}

// NewStateMachine 创建新的状态机实例
func NewStateMachine() StateMachine {
	return &fsm{
		transitions: stateTransitions,
	}
}

// CanTransition 检查是否允许状态转换
func (f *fsm) CanTransition(from types.TaskState, to types.TaskState) bool {
	// 相同状态不允许转换
	if from == to {
		return false
	}

	// 查找源状态的有效转换列表
	validTos, exists := f.transitions[from]
	if !exists {
		// 如果源状态没有定义转换规则,不允许转换
		return false
	}

	// 检查目标状态是否在有效转换列表中
	for _, validTo := range validTos {
		if validTo == to {
			return true
		}
	}

	return false
}

// GetValidTransitions 获取当前状态的有效转换
func (f *fsm) GetValidTransitions(state types.TaskState) []types.TaskState {
	validTos, exists := f.transitions[state]
	if !exists {
		return []types.TaskState{}
	}

	// 返回副本,避免外部修改
	result := make([]types.TaskState, len(validTos))
	copy(result, validTos)
	return result
}

// Transition 执行状态转换(纯函数实现)
// 使用不可变模式,返回新的任务对象,不修改原任务
func (f *fsm) Transition(t TransitionableTask, to types.TaskState, reason string) (TransitionableTask, error) {
	// 1. 验证状态转换合法性
	from := t.GetState()
	if !f.CanTransition(from, to) {
		return nil, errors.ErrInvalidStateTransition
	}

	// 2. 创建新任务对象(不可变)
	newTask := t.Clone()

	// 3. 更新状态
	newTask.SetState(to)
	newTask.SetUpdatedAt(time.Now())

	// 4. 记录状态变更历史
	change := &StateChange{
		From:   from,
		To:     to,
		Reason: reason,
		Time:   time.Now(),
	}
	newTask.AddStateChange(change)

	return newTask, nil
}
