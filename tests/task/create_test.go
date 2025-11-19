package task_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
)

// createTestTemplate 创建测试用的模板
func createTestTemplate() *template.Template {
	return &template.Template{
		ID:          "tpl-001",
		Name:        "Test Template",
		Description: "Test template description",
		Version:     1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "Start Node",
				Type: template.NodeTypeStart,
			},
			"approval-001": {
				ID:   "approval-001",
				Name: "Approval Node",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeSingle,
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{"user-001"},
					},
				},
			},
			"end": {
				ID:   "end",
				Name: "End Node",
				Type: template.NodeTypeEnd,
			},
		},
		Edges: []*template.Edge{
			{From: "start", To: "approval-001"},
			{From: "approval-001", To: "end"},
		},
	}
}

// TestTaskManagerCreate 测试任务创建功能
func TestTaskManagerCreate(t *testing.T) {
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
	tsk, err := taskMgr.Create("tpl-001", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 验证任务基本信息
	if tsk.ID == "" {
		t.Error("Task.ID should be set")
	}
	if tsk.TemplateID != "tpl-001" {
		t.Errorf("Task.TemplateID = %q, want %q", tsk.TemplateID, "tpl-001")
	}
	if tsk.BusinessID != "biz-001" {
		t.Errorf("Task.BusinessID = %q, want %q", tsk.BusinessID, "biz-001")
	}
	if tsk.State != task.TaskStatePending {
		t.Errorf("Task.State = %v, want %v", tsk.State, task.TaskStatePending)
	}
	if tsk.CurrentNode != "start" {
		t.Errorf("Task.CurrentNode = %q, want %q", tsk.CurrentNode, "start")
	}
	if tsk.CreatedAt.IsZero() {
		t.Error("Task.CreatedAt should be set")
	}
}

// TestTaskManagerCreateWithParams 测试带参数的任务创建
func TestTaskManagerCreateWithParams(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplate()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	taskMgr := task.NewTaskManager(templateMgr, nil)

	params := json.RawMessage(`{"amount": 1000, "department": "IT"}`)
	tsk, err := taskMgr.Create("tpl-001", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 验证参数
	if len(tsk.Params) == 0 {
		t.Error("Task.Params should not be empty")
	}

	// 验证可以解析 JSON
	var data map[string]interface{}
	if err := json.Unmarshal(tsk.Params, &data); err != nil {
		t.Errorf("Task.Params should be valid JSON, got error: %v", err)
	}
}

// TestTaskManagerCreateTemplateNotFound 测试模板不存在的情况
func TestTaskManagerCreateTemplateNotFound(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	taskMgr := task.NewTaskManager(templateMgr, nil)

	params := json.RawMessage(`{}`)
	_, err := taskMgr.Create("non-existent", "biz-001", params)
	if err == nil {
		t.Error("Create() should fail when template does not exist")
	}
}

// TestTaskManagerCreateWithNilParams 测试空参数的任务创建
func TestTaskManagerCreateWithNilParams(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplate()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	taskMgr := task.NewTaskManager(templateMgr, nil)

	tsk, err := taskMgr.Create("tpl-001", "biz-001", nil)
	if err != nil {
		t.Fatalf("Create() failed with nil params: %v", err)
	}

	if tsk == nil {
		t.Fatal("Create() should return a task even with nil params")
	}
}

// TestTaskManagerCreateInitializesRuntimeData 测试运行时数据初始化
func TestTaskManagerCreateInitializesRuntimeData(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplate()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	taskMgr := task.NewTaskManager(templateMgr, nil)

	tsk, err := taskMgr.Create("tpl-001", "biz-001", nil)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 验证运行时数据已初始化
	if tsk.NodeOutputs == nil {
		t.Error("Task.NodeOutputs should be initialized")
	}
	if tsk.Approvers == nil {
		t.Error("Task.Approvers should be initialized")
	}
	if tsk.Approvals == nil {
		t.Error("Task.Approvals should be initialized")
	}
	if tsk.Records == nil {
		t.Error("Task.Records should be initialized")
	}
	if tsk.StateHistory == nil {
		t.Error("Task.StateHistory should be initialized")
	}
}

