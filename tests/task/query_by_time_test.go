package task_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
)

// TestQueryByTimeRange 测试按时间范围查询任务
func TestQueryByTimeRange(t *testing.T) {
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

	// 创建第一个任务
	params := json.RawMessage(`{"amount": 1000}`)
	task1, err := taskMgr.Create("tpl-001", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	time.Sleep(10 * time.Millisecond)

	// 创建第二个任务
	task2, err := taskMgr.Create("tpl-001", "biz-002", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	time.Sleep(10 * time.Millisecond)
	endTime := time.Now()

	// 查询 startTime 到 endTime 之间的任务
	filter := &task.TaskFilter{
		StartTime: startTime,
		EndTime:   endTime,
	}
	tasks, err := taskMgr.Query(filter)
	if err != nil {
		t.Fatalf("Query() failed: %v", err)
	}

	// 应该返回 task1 和 task2
	if len(tasks) < 2 {
		t.Errorf("Query() returned %d tasks, want at least 2", len(tasks))
	}

	// 验证任务 ID
	taskIDs := make(map[string]bool)
	for _, tsk := range tasks {
		taskIDs[tsk.ID] = true
	}

	if !taskIDs[task1.ID] {
		t.Errorf("Query() should include task %q", task1.ID)
	}

	if !taskIDs[task2.ID] {
		t.Errorf("Query() should include task %q", task2.ID)
	}
}

// TestQueryByStartTime 测试只设置开始时间
func TestQueryByStartTime(t *testing.T) {
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

	// 创建任务
	params := json.RawMessage(`{"amount": 1000}`)
	task1, err := taskMgr.Create("tpl-001", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 查询 startTime 之后的任务
	filter := &task.TaskFilter{
		StartTime: startTime,
	}
	tasks, err := taskMgr.Query(filter)
	if err != nil {
		t.Fatalf("Query() failed: %v", err)
	}

	// 应该返回 task1
	if len(tasks) < 1 {
		t.Errorf("Query() returned %d tasks, want at least 1", len(tasks))
	}

	// 验证任务 ID
	found := false
	for _, tsk := range tasks {
		if tsk.ID == task1.ID {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Query() should include task %q", task1.ID)
	}
}

// TestQueryByEndTime 测试只设置结束时间
func TestQueryByEndTime(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplate()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	// 创建任务管理器
	taskMgr := task.NewTaskManager(templateMgr, nil)

	// 创建任务
	params := json.RawMessage(`{"amount": 1000}`)
	task1, err := taskMgr.Create("tpl-001", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	time.Sleep(10 * time.Millisecond)
	endTime := time.Now()

	// 查询 endTime 之前的任务
	filter := &task.TaskFilter{
		EndTime: endTime,
	}
	tasks, err := taskMgr.Query(filter)
	if err != nil {
		t.Fatalf("Query() failed: %v", err)
	}

	// 应该返回 task1
	if len(tasks) < 1 {
		t.Errorf("Query() returned %d tasks, want at least 1", len(tasks))
	}

	// 验证任务 ID
	found := false
	for _, tsk := range tasks {
		if tsk.ID == task1.ID {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Query() should include task %q", task1.ID)
	}
}

