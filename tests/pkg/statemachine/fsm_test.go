package statemachine_test

import (
	"testing"
	"time"

	pkgSM "github.com/mautops/approval-kit/pkg/statemachine"
	internalSM "github.com/mautops/approval-kit/internal/statemachine"
	internalTypes "github.com/mautops/approval-kit/internal/types"
	"github.com/mautops/approval-kit/pkg/types"
)

// TestPkgStateMachineInterface 验证 pkg/statemachine 包中的 StateMachine 接口定义
func TestPkgStateMachineInterface(t *testing.T) {
	// 验证接口类型存在
	var sm pkgSM.StateMachine
	if sm != nil {
		_ = sm
	}
}

// TestPkgStateMachineCompatibility 验证 pkg/statemachine.StateMachine 与 internal/statemachine.StateMachine 兼容
func TestPkgStateMachineCompatibility(t *testing.T) {
	// 验证 pkg 接口可以被 internal 实现通过适配器满足
	var _ pkgSM.StateMachine = (*internalStateMachineAdapter)(nil)
}

// TestPkgStateMachineMethods 验证 StateMachine 接口方法
func TestPkgStateMachineMethods(t *testing.T) {
	internalSMImpl := internalSM.NewStateMachine()
	sm := &internalStateMachineAdapter{impl: internalSMImpl}

	// 测试 CanTransition
	canTransition := sm.CanTransition(types.TaskStatePending, types.TaskStateSubmitted)
	if !canTransition {
		t.Errorf("CanTransition(pending, submitted) = false, want true")
	}

	// 测试 GetValidTransitions
	validTransitions := sm.GetValidTransitions(types.TaskStatePending)
	if len(validTransitions) == 0 {
		t.Errorf("GetValidTransitions(pending) returned empty list")
	}
}

// internalStateMachineAdapter 用于测试接口兼容性的适配器
type internalStateMachineAdapter struct {
	impl internalSM.StateMachine
}

func (a *internalStateMachineAdapter) CanTransition(from types.TaskState, to types.TaskState) bool {
	// 类型别名可以直接转换
	return a.impl.CanTransition(from, to)
}

func (a *internalStateMachineAdapter) Transition(task pkgSM.TransitionableTask, to types.TaskState, reason string) (pkgSM.TransitionableTask, error) {
	// 需要将 pkg.TransitionableTask 转换为 internal.TransitionableTask
	// 这里简化处理,实际使用时需要适配器
	internalTask := &internalTransitionableTaskAdapter{impl: task}
	_, err := a.impl.Transition(internalTask, to, reason)
	if err != nil {
		return nil, err
	}
	// 将 internal.TransitionableTask 转换回 pkg.TransitionableTask
	// 这里简化处理,返回原任务
	return task, nil
}

func (a *internalStateMachineAdapter) GetValidTransitions(state types.TaskState) []types.TaskState {
	return a.impl.GetValidTransitions(state)
}

// internalTransitionableTaskAdapter 用于适配 pkg.TransitionableTask 到 internal.TransitionableTask
type internalTransitionableTaskAdapter struct {
	impl pkgSM.TransitionableTask
}

func (a *internalTransitionableTaskAdapter) GetState() internalTypes.TaskState {
	return a.impl.GetState()
}

func (a *internalTransitionableTaskAdapter) SetState(state internalTypes.TaskState) {
	a.impl.SetState(state)
}

func (a *internalTransitionableTaskAdapter) GetUpdatedAt() time.Time {
	return a.impl.GetUpdatedAt()
}

func (a *internalTransitionableTaskAdapter) SetUpdatedAt(t time.Time) {
	a.impl.SetUpdatedAt(t)
}

func (a *internalTransitionableTaskAdapter) GetStateHistory() []*internalSM.StateChange {
	pkgHistory := a.impl.GetStateHistory()
	result := make([]*internalSM.StateChange, len(pkgHistory))
	for i, change := range pkgHistory {
		result[i] = pkgSM.StateChangeToInternal(change)
	}
	return result
}

func (a *internalTransitionableTaskAdapter) AddStateChange(change *internalSM.StateChange) {
	pkgChange := pkgSM.StateChangeFromInternal(change)
	a.impl.AddStateChange(pkgChange)
}

func (a *internalTransitionableTaskAdapter) Clone() internalSM.TransitionableTask {
	cloned := a.impl.Clone()
	// 将 pkg.TransitionableTask 包装为 internal.TransitionableTask
	return &internalTransitionableTaskAdapter{impl: cloned}
}

