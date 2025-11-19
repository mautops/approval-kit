package task_test

import (
	"encoding/json"
	"testing"

	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
)

// TestTaskManagerGet 测试任务查询功能
func TestTaskManagerGet(t *testing.T) {
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
	createdTask, err := taskMgr.Create("tpl-001", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 查询任务
	got, err := taskMgr.Get(createdTask.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	// 验证任务信息
	if got.ID != createdTask.ID {
		t.Errorf("Get() returned task ID = %q, want %q", got.ID, createdTask.ID)
	}
	if got.TemplateID != "tpl-001" {
		t.Errorf("Get() returned task TemplateID = %q, want %q", got.TemplateID, "tpl-001")
	}
	if got.BusinessID != "biz-001" {
		t.Errorf("Get() returned task BusinessID = %q, want %q", got.BusinessID, "biz-001")
	}
}

// TestTaskManagerGetNotFound 测试查询不存在的任务
func TestTaskManagerGetNotFound(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	taskMgr := task.NewTaskManager(templateMgr, nil)

	// 查询不存在的任务
	_, err := taskMgr.Get("non-existent")
	if err == nil {
		t.Error("Get() should fail when task does not exist")
	}
}

// TestTaskManagerGetIsolation 测试查询返回的任务是隔离的
func TestTaskManagerGetIsolation(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplate()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	taskMgr := task.NewTaskManager(templateMgr, nil)

	// 创建任务
	createdTask, err := taskMgr.Create("tpl-001", "biz-001", nil)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 获取任务
	got, err := taskMgr.Get(createdTask.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	// 修改返回的任务
	got.BusinessID = "modified-biz-id"

	// 再次获取,验证原任务未被修改
	got2, err := taskMgr.Get(createdTask.ID)
	if err != nil {
		t.Fatalf("Get() failed on second call: %v", err)
	}
	if got2.BusinessID == "modified-biz-id" {
		t.Error("Get() should return a copy, modifications should not affect stored task")
	}
	if got2.BusinessID != "biz-001" {
		t.Errorf("Get() returned task BusinessID = %q, want %q", got2.BusinessID, "biz-001")
	}
}

// TestTaskManagerGetConcurrent 测试并发查询
func TestTaskManagerGetConcurrent(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplate()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	taskMgr := task.NewTaskManager(templateMgr, nil)

	// 创建任务
	createdTask, err := taskMgr.Create("tpl-001", "biz-001", nil)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 并发查询
	concurrency := 100
	results := make(chan *task.Task, concurrency)
	errors := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			tsk, err := taskMgr.Get(createdTask.ID)
			if err != nil {
				errors <- err
				return
			}
			results <- tsk
		}()
	}

	// 收集结果
	successCount := 0
	errorCount := 0
	for i := 0; i < concurrency; i++ {
		select {
		case <-results:
			successCount++
		case <-errors:
			errorCount++
		}
	}

	if errorCount > 0 {
		t.Errorf("Get() failed %d times in concurrent access", errorCount)
	}
	if successCount != concurrency {
		t.Errorf("Get() succeeded %d times, want %d", successCount, concurrency)
	}
}

