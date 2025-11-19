package task_test

import (
	"encoding/json"
	"testing"

	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
)

// TestQueryByApprover 测试按审批人查询任务
func TestQueryByApprover(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateWithApprovalNode()
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

	// 提交任务
	err = taskMgr.Submit(task1.ID)
	if err != nil {
		t.Fatalf("Submit() failed: %v", err)
	}

	err = taskMgr.Submit(task2.ID)
	if err != nil {
		t.Fatalf("Submit() failed: %v", err)
	}

	// 通过 AddApprover 添加审批人(模板需要允许加签)
	// 注意: createTestTemplateWithApprovalNode 需要配置 AllowAddApprover 权限
	err = taskMgr.AddApprover(task1.ID, "approval-001", "user-001", "test add approver")
	if err != nil {
		t.Fatalf("AddApprover() failed: %v", err)
	}

	err = taskMgr.AddApprover(task2.ID, "approval-001", "user-001", "test add approver")
	if err != nil {
		t.Fatalf("AddApprover() failed: %v", err)
	}

	// 查询 user-001 的待审批任务
	filter := &task.TaskFilter{
		Approver: "user-001",
	}
	tasks, err := taskMgr.Query(filter)
	if err != nil {
		t.Fatalf("Query() failed: %v", err)
	}

	// 应该返回 task1 和 task2(因为模板中配置了 user-001 作为审批人)
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

// TestQueryByApproverEmpty 测试查询不存在的审批人
func TestQueryByApproverEmpty(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplate()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	// 创建任务管理器
	taskMgr := task.NewTaskManager(templateMgr, nil)

	// 查询不存在的审批人
	filter := &task.TaskFilter{
		Approver: "non-existent-user",
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


