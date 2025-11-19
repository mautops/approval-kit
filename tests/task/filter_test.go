package task_test

import (
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/types"
)

// TestTaskFilter 测试 TaskFilter 结构体
func TestTaskFilter(t *testing.T) {
	filter := &task.TaskFilter{
		State:      types.TaskStateApproving,
		TemplateID: "tpl-001",
		BusinessID: "biz-001",
		Approver:   "user-001",
		StartTime:  time.Now().Add(-24 * time.Hour),
		EndTime:    time.Now(),
	}

	if filter.State != types.TaskStateApproving {
		t.Errorf("TaskFilter.State = %q, want %q", filter.State, types.TaskStateApproving)
	}

	if filter.TemplateID != "tpl-001" {
		t.Errorf("TaskFilter.TemplateID = %q, want %q", filter.TemplateID, "tpl-001")
	}

	if filter.BusinessID != "biz-001" {
		t.Errorf("TaskFilter.BusinessID = %q, want %q", filter.BusinessID, "biz-001")
	}

	if filter.Approver != "user-001" {
		t.Errorf("TaskFilter.Approver = %q, want %q", filter.Approver, "user-001")
	}
}

// TestTaskFilterEmpty 测试空过滤器
func TestTaskFilterEmpty(t *testing.T) {
	filter := &task.TaskFilter{}

	if filter.State != "" {
		t.Errorf("TaskFilter.State should be empty, got %q", filter.State)
	}

	if filter.TemplateID != "" {
		t.Errorf("TaskFilter.TemplateID should be empty, got %q", filter.TemplateID)
	}

	if filter.BusinessID != "" {
		t.Errorf("TaskFilter.BusinessID should be empty, got %q", filter.BusinessID)
	}

	if filter.Approver != "" {
		t.Errorf("TaskFilter.Approver should be empty, got %q", filter.Approver)
	}
}

