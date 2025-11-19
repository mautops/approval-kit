package task_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
	"github.com/mautops/approval-kit/internal/types"
)

// TestQueryCombined 测试综合查询功能(多条件组合)
func TestQueryCombined(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	tpl1 := createTestTemplate()
	tpl1.ID = "tpl-001"
	err := templateMgr.Create(tpl1)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	tpl2 := createTestTemplate()
	tpl2.ID = "tpl-002"
	err = templateMgr.Create(tpl2)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	// 创建任务管理器
	taskMgr := task.NewTaskManager(templateMgr, nil)

	// 记录开始时间
	startTime := time.Now()
	time.Sleep(10 * time.Millisecond)

	// 创建多个任务
	params := json.RawMessage(`{"amount": 1000}`)
	task1, err := taskMgr.Create("tpl-001", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	_, err = taskMgr.Create("tpl-001", "biz-002", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	_, err = taskMgr.Create("tpl-002", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 提交 task1
	err = taskMgr.Submit(task1.ID)
	if err != nil {
		t.Fatalf("Submit() failed: %v", err)
	}

	time.Sleep(10 * time.Millisecond)
	endTime := time.Now()

	// 综合查询: tpl-001 + submitted 状态 + 时间范围
	filter := &task.TaskFilter{
		TemplateID: "tpl-001",
		State:      types.TaskStateSubmitted,
		StartTime:  startTime,
		EndTime:    endTime,
	}
	tasks, err := taskMgr.Query(filter)
	if err != nil {
		t.Fatalf("Query() failed: %v", err)
	}

	// 应该只返回 task1
	if len(tasks) != 1 {
		t.Errorf("Query() returned %d tasks, want 1", len(tasks))
	}

	if tasks[0].ID != task1.ID {
		t.Errorf("Query() returned task ID %q, want %q", tasks[0].ID, task1.ID)
	}
}

// TestQueryCombinedMultipleConditions 测试多个条件组合查询
func TestQueryCombinedMultipleConditions(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplate()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	// 创建任务管理器
	taskMgr := task.NewTaskManager(templateMgr, nil)

	// 记录开始时间
	startTime := time.Now()
	time.Sleep(10 * time.Millisecond)

	// 创建多个任务
	params := json.RawMessage(`{"amount": 1000}`)
	task1, err := taskMgr.Create("tpl-001", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	task2, err := taskMgr.Create("tpl-001", "biz-002", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	task3, err := taskMgr.Create("tpl-001", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	time.Sleep(10 * time.Millisecond)
	endTime := time.Now()

	// 综合查询: tpl-001 + biz-001 + 时间范围
	filter := &task.TaskFilter{
		TemplateID: "tpl-001",
		BusinessID: "biz-001",
		StartTime:  startTime,
		EndTime:    endTime,
	}
	tasks, err := taskMgr.Query(filter)
	if err != nil {
		t.Fatalf("Query() failed: %v", err)
	}

	// 应该返回 task1 和 task3
	if len(tasks) != 2 {
		t.Errorf("Query() returned %d tasks, want 2", len(tasks))
	}

	// 验证任务 ID
	taskIDs := make(map[string]bool)
	for _, tsk := range tasks {
		taskIDs[tsk.ID] = true
	}

	if !taskIDs[task1.ID] {
		t.Errorf("Query() should include task %q", task1.ID)
	}

	if !taskIDs[task3.ID] {
		t.Errorf("Query() should include task %q", task3.ID)
	}

	if taskIDs[task2.ID] {
		t.Errorf("Query() should not include task %q", task2.ID)
	}
}

// TestQueryNoFilter 测试无过滤条件查询
func TestQueryNoFilter(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplate()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	// 创建任务管理器
	taskMgr := task.NewTaskManager(templateMgr, nil)

	// 创建多个任务
	params := json.RawMessage(`{"amount": 1000}`)
	_, err = taskMgr.Create("tpl-001", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	_, err = taskMgr.Create("tpl-001", "biz-002", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 无过滤条件查询
	filter := &task.TaskFilter{}
	tasks, err := taskMgr.Query(filter)
	if err != nil {
		t.Fatalf("Query() failed: %v", err)
	}

	// 应该返回所有任务
	if len(tasks) < 2 {
		t.Errorf("Query() returned %d tasks, want at least 2", len(tasks))
	}
}

