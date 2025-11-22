package task

import (
	internalTask "github.com/mautops/approval-kit/internal/task"
)

// Approval 审批结果
// 与 internal/task.Approval 结构相同,但位于 pkg 目录,可以被外部导入
type Approval = internalTask.Approval

// Record 审批记录
// 每次审批操作时自动生成,作为状态流转的副产品
// 与 internal/task.Record 结构相同,但位于 pkg 目录,可以被外部导入
type Record = internalTask.Record

// StateChange 状态变更记录
// 记录每次状态变更的详细信息,用于追溯和审计
// 与 internal/task.StateChange 结构相同,但位于 pkg 目录,可以被外部导入
type StateChange = internalTask.StateChange

