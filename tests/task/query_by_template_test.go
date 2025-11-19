package task_test

import (
	"encoding/json"
	"testing"

	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
)

// TestQueryByTemplate 测试按模板查询任务
func TestQueryByTemplate(t *testing.T) {
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

	// 创建多个任务(使用不同模板)
	params := json.RawMessage(`{"amount": 1000}`)
	task1, err := taskMgr.Create("tpl-001", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	task2, err := taskMgr.Create("tpl-001", "biz-002", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	task3, err := taskMgr.Create("tpl-002", "biz-003", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 查询 tpl-001 模板的任务
	filter := &task.TaskFilter{
		TemplateID: "tpl-001",
	}
	tasks, err := taskMgr.Query(filter)
	if err != nil {
		t.Fatalf("Query() failed: %v", err)
	}

	// 应该返回 task1 和 task2
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

	if !taskIDs[task2.ID] {
		t.Errorf("Query() should include task %q", task2.ID)
	}

	if taskIDs[task3.ID] {
		t.Errorf("Query() should not include task %q", task3.ID)
	}
}

// TestQueryByTemplateEmpty 测试查询不存在的模板
func TestQueryByTemplateEmpty(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplate()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	// 创建任务管理器
	taskMgr := task.NewTaskManager(templateMgr, nil)

	// 查询不存在的模板
	filter := &task.TaskFilter{
		TemplateID: "non-existent",
	}
	tasks, err := taskMgr.Query(filter)
	if err != nil {
		t.Fatalf("Query() failed: %v", err)
	}

	// 应该返回空列表
	if len(tasks) != 0 {
		t.Errorf("Query() returned %d tasks, want 0", len(tasks))
	}
}

