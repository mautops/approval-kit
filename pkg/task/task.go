package task

import (
	internalTask "github.com/mautops/approval-kit/internal/task"
)

// Task 表示审批任务实例
// 与 internal/task.Task 结构相同,但位于 pkg 目录,可以被外部导入
type Task = internalTask.Task

// FromInternal 将 internal.Task 转换为 pkg.Task
func FromInternal(t *internalTask.Task) *Task {
	return (*Task)(t)
}

// ToInternal 将 pkg.Task 转换为 internal.Task
func ToInternal(t *Task) *internalTask.Task {
	return (*internalTask.Task)(t)
}

