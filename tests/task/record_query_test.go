package task_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
)

// TestRecordQuery 测试审批记录查询
func TestRecordQuery(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	// 使用多人会签模式,允许多个审批人
	tpl := createTestTemplateWithMultipleApprovers()
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

	// 提交任务
	err = taskMgr.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Submit() failed: %v", err)
	}

	// 手动设置审批人列表(因为固定审批人需要在节点激活时设置,这里为了测试手动设置)
	// 通过 AddApprover 添加审批人(需要模板允许加签)
	err = taskMgr.AddApprover(tsk.ID, "approval-001", "user-001", "setup approver")
	if err != nil {
		t.Fatalf("AddApprover() failed: %v", err)
	}

	err = taskMgr.AddApprover(tsk.ID, "approval-001", "user-002", "setup approver")
	if err != nil {
		t.Fatalf("AddApprover() failed: %v", err)
	}

	// 执行多次审批操作(多人会签模式,需要所有审批人都同意)
	err = taskMgr.Approve(tsk.ID, "approval-001", "user-001", "first approval")
	if err != nil {
		t.Fatalf("Approve() failed: %v", err)
	}

	err = taskMgr.Approve(tsk.ID, "approval-001", "user-002", "second approval")
	if err != nil {
		t.Fatalf("Approve() failed: %v", err)
	}

	// 通过任务对象查询所有记录
	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	records := tsk.GetRecords()
	// 应该有 4 条记录: 2 条 add_approver 记录 + 2 条 approve 记录
	if len(records) < 2 {
		t.Errorf("GetRecords() returned %d records, want at least 2", len(records))
	}

	// 验证有审批记录(approve 类型的记录)
	approveRecords := 0
	for _, record := range records {
		if record.Result == "approve" {
			approveRecords++
		}
	}
	if approveRecords != 2 {
		t.Errorf("Expected 2 approve records, got %d", approveRecords)
	}

	// 验证记录内容
	foundUser001 := false
	foundUser002 := false
	for _, record := range records {
		if record.Result == "approve" {
			if record.Approver == "user-001" {
				foundUser001 = true
			}
			if record.Approver == "user-002" {
				foundUser002 = true
			}
		}
	}
	if !foundUser001 {
		t.Error("Expected to find approve record for user-001")
	}
	if !foundUser002 {
		t.Error("Expected to find approve record for user-002")
	}
}

// TestRecordQueryByNodeID 测试按节点 ID 查询审批记录
func TestRecordQueryByNodeID(t *testing.T) {
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

	// 提交任务
	err = taskMgr.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Submit() failed: %v", err)
	}

	// 在节点执行审批操作
	err = taskMgr.Approve(tsk.ID, "approval-001", "user-001", "approval on node 1")
	if err != nil {
		t.Fatalf("Approve() failed: %v", err)
	}

	// 通过任务对象查询记录
	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	// 查询特定节点的记录
	records := tsk.GetRecordsByNodeID("approval-001")
	if len(records) != 1 {
		t.Errorf("GetRecordsByNodeID() returned %d records, want 1", len(records))
	}

	if records[0].NodeID != "approval-001" {
		t.Errorf("GetRecordsByNodeID()[0].NodeID = %q, want %q", records[0].NodeID, "approval-001")
	}
}

// TestRecordQueryEmpty 测试空记录查询
func TestRecordQueryEmpty(t *testing.T) {
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

	// 通过任务对象查询记录(还没有审批操作)
	records := tsk.GetRecords()
	if len(records) != 0 {
		t.Errorf("GetRecords() returned %d records, want 0", len(records))
	}

	// 查询特定节点的记录
	records = tsk.GetRecordsByNodeID("approval-001")
	if len(records) != 0 {
		t.Errorf("GetRecordsByNodeID() returned %d records, want 0", len(records))
	}
}

// createTestTemplateWithMultipleApprovers 创建支持多个审批人的测试模板(多人会签模式)
func createTestTemplateWithMultipleApprovers() *template.Template {
	return &template.Template{
		ID:          "tpl-001",
		Name:        "Test Template With Multiple Approvers",
		Description: "Test template with multiple approvers",
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
					Mode: node.ApprovalModeUnanimous, // 多人会签模式
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{"user-001", "user-002"},
					},
					Permissions: node.OperationPermissions{
						AllowTransfer:    true,
						AllowAddApprover: true,
						AllowRemoveApprover: true,
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

