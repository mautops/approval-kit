package task_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
)

// createTestTemplateWithReplacePermission 创建允许替换审批人的测试模板
func createTestTemplateWithReplacePermission() *template.Template {
	return &template.Template{
		ID:          "tpl-001",
		Name:        "Test Template With Replace Permission",
		Description: "Test template with replace permission",
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
					Permissions: node.OperationPermissions{
						AllowAddApprover: true, // 允许加签,用于设置审批人列表
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

// TestReplaceApprover 测试替换审批人功能
func TestReplaceApprover(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateWithReplacePermission()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	taskMgr := task.NewTaskManager(templateMgr, nil)

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

	// 验证提交后当前节点已更新为审批节点
	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}
	if tsk.CurrentNode != "approval-001" {
		t.Fatalf("Task.CurrentNode = %q, want %q", tsk.CurrentNode, "approval-001")
	}

	// 手动设置审批人列表(因为固定审批人需要在节点激活时设置)
	err = taskMgr.AddApprover(tsk.ID, "approval-001", "user-001", "setup approver")
	if err != nil {
		t.Fatalf("AddApprover() failed: %v", err)
	}

	// 替换审批人
	err = taskMgr.ReplaceApprover(tsk.ID, "approval-001", "user-001", "user-002", "user replaced")
	if err != nil {
		t.Fatalf("ReplaceApprover() failed: %v", err)
	}

	// 验证审批人已替换
	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	approvers := tsk.Approvers["approval-001"]
	if len(approvers) != 1 {
		t.Fatalf("Task.Approvers length = %d, want 1", len(approvers))
	}
	if approvers[0] != "user-002" {
		t.Errorf("Task.Approvers[0] = %q, want %q", approvers[0], "user-002")
	}

	// 验证原审批人不在审批人列表中
	for _, approver := range approvers {
		if approver == "user-001" {
			t.Error("Old approver should not be in approvers list")
		}
	}
}

// TestReplaceApproverNotFound 测试替换不存在的审批人
func TestReplaceApproverNotFound(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateWithReplacePermission()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	taskMgr := task.NewTaskManager(templateMgr, nil)

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

	// 手动设置审批人列表
	err = taskMgr.AddApprover(tsk.ID, "approval-001", "user-001", "setup approver")
	if err != nil {
		t.Fatalf("AddApprover() failed: %v", err)
	}

	// 尝试替换不存在的审批人(应该失败)
	err = taskMgr.ReplaceApprover(tsk.ID, "approval-001", "non-existent", "user-002", "user replaced")
	if err == nil {
		t.Error("ReplaceApprover() should fail for non-existent approver")
	}
}

// TestReplaceApproverAlreadyApproved 测试替换已审批的审批人
func TestReplaceApproverAlreadyApproved(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateWithReplacePermission()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	taskMgr := task.NewTaskManager(templateMgr, nil)

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

	// 手动设置审批人列表
	err = taskMgr.AddApprover(tsk.ID, "approval-001", "user-001", "setup approver")
	if err != nil {
		t.Fatalf("AddApprover() failed: %v", err)
	}

	// 执行审批操作
	err = taskMgr.Approve(tsk.ID, "approval-001", "user-001", "approved")
	if err != nil {
		t.Fatalf("Approve() failed: %v", err)
	}

	// 尝试替换已审批的审批人(应该失败)
	err = taskMgr.ReplaceApprover(tsk.ID, "approval-001", "user-001", "user-002", "user replaced")
	if err == nil {
		t.Error("ReplaceApprover() should fail for already approved approver")
	}
}

// TestReplaceApproverTaskNotFound 测试替换不存在的任务的审批人
func TestReplaceApproverTaskNotFound(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	taskMgr := task.NewTaskManager(templateMgr, nil)

	err := taskMgr.ReplaceApprover("non-existent", "approval-001", "user-001", "user-002", "user replaced")
	if err == nil {
		t.Error("ReplaceApprover() should fail for non-existent task")
	}
}

// TestReplaceApproverRecord 测试替换审批人时生成记录
func TestReplaceApproverRecord(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateWithReplacePermission()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	taskMgr := task.NewTaskManager(templateMgr, nil)

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

	// 验证提交后当前节点已更新为审批节点
	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}
	if tsk.CurrentNode != "approval-001" {
		t.Fatalf("Task.CurrentNode = %q, want %q", tsk.CurrentNode, "approval-001")
	}

	// 手动设置审批人列表
	err = taskMgr.AddApprover(tsk.ID, "approval-001", "user-001", "setup approver")
	if err != nil {
		t.Fatalf("AddApprover() failed: %v", err)
	}

	// 记录初始记录数
	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}
	initialRecordCount := len(tsk.Records)

	// 替换审批人
	err = taskMgr.ReplaceApprover(tsk.ID, "approval-001", "user-001", "user-002", "user replaced")
	if err != nil {
		t.Fatalf("ReplaceApprover() failed: %v", err)
	}

	// 验证生成了替换记录
	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	if len(tsk.Records) != initialRecordCount+1 {
		t.Errorf("Task.Records length = %d, want %d", len(tsk.Records), initialRecordCount+1)
	}

	// 检查最后一条记录是否是替换记录
	lastRecord := tsk.Records[len(tsk.Records)-1]
	if lastRecord.Result != "replace" {
		t.Errorf("Last record.Result = %q, want %q", lastRecord.Result, "replace")
	}
	if lastRecord.Approver != "user-001" {
		t.Errorf("Last record.Approver = %q, want %q", lastRecord.Approver, "user-001")
	}
}

