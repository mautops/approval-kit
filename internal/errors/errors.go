package errors

import (
	"fmt"
)

// 错误定义
// 所有错误常量遵循 Go 语言错误处理最佳实践
var (
	// ErrInvalidTemplate 表示模板无效
	ErrInvalidTemplate = fmt.Errorf("invalid template")

	// ErrInvalidStateTransition 表示状态转换无效
	ErrInvalidStateTransition = fmt.Errorf("invalid state transition")

	// ErrNodeNotFound 表示节点未找到
	ErrNodeNotFound = fmt.Errorf("node not found")

	// ErrApproverNotFound 表示审批人未找到
	ErrApproverNotFound = fmt.Errorf("approver not found")

	// ErrApprovalPending 表示审批待处理
	ErrApprovalPending = fmt.Errorf("approval pending")

	// ErrConcurrentModification 表示并发修改冲突
	ErrConcurrentModification = fmt.Errorf("concurrent modification")

	// ErrEventPushFailed 表示事件推送失败
	ErrEventPushFailed = fmt.Errorf("event push failed")
)

