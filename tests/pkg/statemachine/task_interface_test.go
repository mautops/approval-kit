package statemachine_test

import (
	"testing"
	"time"

	pkgSM "github.com/mautops/approval-kit/pkg/statemachine"
	"github.com/mautops/approval-kit/pkg/types"
)

// TestPkgTransitionableTaskInterface 验证 pkg/statemachine 包中的 TransitionableTask 接口定义
func TestPkgTransitionableTaskInterface(t *testing.T) {
	// 验证接口类型存在
	var task pkgSM.TransitionableTask
	if task != nil {
		_ = task
	}
}

// TestPkgStateChangeType 验证 pkg/statemachine.StateChange 类型暴露
func TestPkgStateChangeType(t *testing.T) {
	var stateChange pkgSM.StateChange
	stateChange.From = types.TaskStatePending
	stateChange.To = types.TaskStateSubmitted
	stateChange.Reason = "test"
	stateChange.Time = time.Now()

	if stateChange.From != types.TaskStatePending {
		t.Errorf("StateChange.From = %v, want pending", stateChange.From)
	}
}

// TestPkgTransitionableTaskCompatibility 验证 pkg/statemachine.TransitionableTask 接口定义
func TestPkgTransitionableTaskCompatibility(t *testing.T) {
	// 验证接口类型存在
	var task pkgSM.TransitionableTask
	if task != nil {
		_ = task
	}
}

