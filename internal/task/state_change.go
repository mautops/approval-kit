package task

import (
	"fmt"
	"time"
)

// Validate 验证状态变更记录的有效性
func (sc *StateChange) Validate() error {
	if sc.From == "" {
		return fmt.Errorf("state change From is required")
	}
	if sc.To == "" {
		return fmt.Errorf("state change To is required")
	}
	// 验证状态不能相同
	if sc.From == sc.To {
		return fmt.Errorf("state change From and To cannot be the same: %q", sc.From)
	}
	// 验证时间
	if sc.Time.IsZero() {
		return fmt.Errorf("state change Time is required")
	}
	return nil
}

// NewStateChange 创建新的状态变更记录
func NewStateChange(from, to TaskState, reason string) *StateChange {
	return &StateChange{
		From:   from,
		To:     to,
		Reason: reason,
		Time:   time.Now(),
	}
}

