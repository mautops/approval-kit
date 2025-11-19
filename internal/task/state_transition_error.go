package task

import (
	"fmt"
)

// StateTransitionError 表示状态转换错误
// 包含状态转换的上下文信息,支持错误链追踪
type StateTransitionError struct {
	From TaskState // 源状态
	To   TaskState // 目标状态
	Err  error     // 底层错误
}

// Error 实现 error 接口
func (e *StateTransitionError) Error() string {
	return fmt.Sprintf("state transition failed: %s -> %s: %v", e.From, e.To, e.Err)
}

// Unwrap 实现错误展开接口,支持错误链追踪
func (e *StateTransitionError) Unwrap() error {
	return e.Err
}

