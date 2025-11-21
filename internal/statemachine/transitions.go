package statemachine

import (
	"github.com/mautops/approval-kit/internal/types"
)

// stateTransitions 定义状态转换规则
// 键为源状态,值为允许转换到的目标状态列表
	// 空列表表示终态,不允许转换到其他状态
	var stateTransitions = map[types.TaskState][]types.TaskState{
	// 待审批状态: 可以提交、取消或暂停
	types.TaskStatePending: {
		types.TaskStateSubmitted,
		types.TaskStateCancelled,
		types.TaskStatePaused,
	},

	// 已提交状态: 可以进入审批中、取消、撤回或暂停
	types.TaskStateSubmitted: {
		types.TaskStateApproving,
		types.TaskStateCancelled,
		types.TaskStatePending, // 撤回
		types.TaskStatePaused,
	},

	// 审批中状态: 可以批准、拒绝、取消、超时、撤回或暂停
	types.TaskStateApproving: {
		types.TaskStateApproved,
		types.TaskStateRejected,
		types.TaskStateCancelled,
		types.TaskStateTimeout,
		types.TaskStatePending, // 撤回
		types.TaskStatePaused,
	},

		// 已通过状态: 终态,不允许转换
		types.TaskStateApproved: {},

		// 已拒绝状态: 终态,不允许转换
		types.TaskStateRejected: {},

		// 已取消状态: 终态,不允许转换
		types.TaskStateCancelled: {},

	// 已超时状态: 终态,不允许转换
	types.TaskStateTimeout: {},

	// 已暂停状态: 可以恢复到暂停前的状态(待审批、已提交或审批中)
	types.TaskStatePaused: {
		types.TaskStatePending,
		types.TaskStateSubmitted,
		types.TaskStateApproving,
	},
}

// GetStateTransitions 返回状态转换规则映射表
// 用于测试和验证状态转换规则
func GetStateTransitions() map[types.TaskState][]types.TaskState {
	// 返回副本,避免外部修改
	result := make(map[types.TaskState][]types.TaskState, len(stateTransitions))
	for k, v := range stateTransitions {
		// 复制切片
		sliceCopy := make([]types.TaskState, len(v))
		copy(sliceCopy, v)
		result[k] = sliceCopy
	}
	return result
}

