package task

import (
	internalTask "github.com/mautops/approval-kit/internal/task"
)

// TaskFilter 任务查询过滤器
// 与 internal/task.TaskFilter 结构相同,但位于 pkg 目录,可以被外部导入
type TaskFilter = internalTask.TaskFilter

// TaskFilterToInternal 将 pkg.TaskFilter 转换为 internal.TaskFilter
func TaskFilterToInternal(f *TaskFilter) *internalTask.TaskFilter {
	return (*internalTask.TaskFilter)(f)
}

// TaskFilterFromInternal 将 internal.TaskFilter 转换为 pkg.TaskFilter
func TaskFilterFromInternal(f *internalTask.TaskFilter) *TaskFilter {
	return (*TaskFilter)(f)
}

