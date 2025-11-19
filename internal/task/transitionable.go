package task

import (
	"time"

	"github.com/mautops/approval-kit/internal/types"
)

// SetState 设置新状态
func (t *Task) SetState(state types.TaskState) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.State = state
}

// SetUpdatedAt 设置更新时间
func (t *Task) SetUpdatedAt(tm time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.UpdatedAt = tm
}

// AddStateChangeRecord 添加状态变更记录
func (t *Task) AddStateChangeRecord(from types.TaskState, to types.TaskState, reason string, tm time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.StateHistory = append(t.StateHistory, &StateChange{
		From:   from,
		To:     to,
		Reason: reason,
		Time:   tm,
	})
}

