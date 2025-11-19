package task_test

import (
	"encoding/json"
	"testing"

	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
	"github.com/mautops/approval-kit/internal/types"
)

// TestQueryByState 测试按状态查询任务
func TestQueryByState(t *testing.T) {
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
	task1, err := taskMgr.Create("tpl-001", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	task2, err := taskMgr.Create("tpl-001", "biz-002", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 提交 task1
	err = taskMgr.Submit(task1.ID)
	if err != nil {
		t.Fatalf("Submit() failed: %v", err)
	}

	// 查询 pending 状态的任务
	filter := &task.TaskFilter{
		State: types.TaskStatePending,
	}
	tasks, err := taskMgr.Query(filter)
	if err != nil {
		t.Fatalf("Query() failed: %v", err)
	}

	// 应该只有 task2 是 pending 状态
	if len(tasks) != 1 {
		t.Errorf("Query() returned %d tasks, want 1", len(tasks))
	}

	if tasks[0].ID != task2.ID {
		t.Errorf("Query() returned task ID %q, want %q", tasks[0].ID, task2.ID)
	}

	// 查询 submitted 状态的任务
	filter.State = types.TaskStateSubmitted
	tasks, err = taskMgr.Query(filter)
	if err != nil {
		t.Fatalf("Query() failed: %v", err)
	}

	// 应该只有 task1 是 submitted 状态
	if len(tasks) != 1 {
		t.Errorf("Query() returned %d tasks, want 1", len(tasks))
	}

	if tasks[0].ID != task1.ID {
		t.Errorf("Query() returned task ID %q, want %q", tasks[0].ID, task1.ID)
	}
}

// TestQueryByStateEmpty 测试查询不存在的状态
func TestQueryByStateEmpty(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplate()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	// 创建任务管理器
	taskMgr := task.NewTaskManager(templateMgr, nil)

	// 查询不存在的状态
	filter := &task.TaskFilter{
		State: types.TaskStateApproved,
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

