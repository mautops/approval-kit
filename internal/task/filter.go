package task

import (
	"time"

	"github.com/mautops/approval-kit/internal/types"
)

// TaskFilter 任务查询过滤器
// 用于按不同条件查询任务
type TaskFilter struct {
	// State 任务状态(可选)
	State types.TaskState

	// TemplateID 模板 ID(可选)
	TemplateID string

	// BusinessID 业务 ID(可选)
	BusinessID string

	// Approver 审批人(可选,用于查询待审批任务)
	Approver string

	// StartTime 开始时间(可选,用于时间范围查询)
	StartTime time.Time

	// EndTime 结束时间(可选,用于时间范围查询)
	EndTime time.Time
}

