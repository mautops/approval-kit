package task_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
)

// TestRecordGeneration 测试审批记录自动生成
func TestRecordGeneration(t *testing.T) {
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

	// 执行审批操作
	err = taskMgr.Approve(tsk.ID, "approval-001", "user-001", "approved")
	if err != nil {
		t.Fatalf("Approve() failed: %v", err)
	}

	// 获取任务,验证审批记录已生成
	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	// 验证审批记录数量
	if len(tsk.Records) == 0 {
		t.Error("Expected at least one approval record, got 0")
	}

	// 验证审批记录内容
	record := tsk.Records[0]
	if record.TaskID != tsk.ID {
		t.Errorf("Record.TaskID = %q, want %q", record.TaskID, tsk.ID)
	}

	if record.NodeID != "approval-001" {
		t.Errorf("Record.NodeID = %q, want %q", record.NodeID, "approval-001")
	}

	if record.Approver != "user-001" {
		t.Errorf("Record.Approver = %q, want %q", record.Approver, "user-001")
	}

	if record.Result != "approve" {
		t.Errorf("Record.Result = %q, want %q", record.Result, "approve")
	}

	if record.Comment != "approved" {
		t.Errorf("Record.Comment = %q, want %q", record.Comment, "approved")
	}

	if record.CreatedAt.IsZero() {
		t.Error("Record.CreatedAt should not be zero")
	}
}

// TestRecordGenerationOnReject 测试拒绝操作时生成审批记录
func TestRecordGenerationOnReject(t *testing.T) {
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

	// 执行拒绝操作
	err = taskMgr.Reject(tsk.ID, "approval-001", "user-001", "rejected")
	if err != nil {
		t.Fatalf("Reject() failed: %v", err)
	}

	// 获取任务,验证审批记录已生成
	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	// 验证审批记录数量
	if len(tsk.Records) == 0 {
		t.Error("Expected at least one approval record, got 0")
	}

	// 验证审批记录内容
	record := tsk.Records[0]
	if record.Result != "reject" {
		t.Errorf("Record.Result = %q, want %q", record.Result, "reject")
	}

	if record.Comment != "rejected" {
		t.Errorf("Record.Comment = %q, want %q", record.Comment, "rejected")
	}
}

// TestRecordMultipleApprovals 测试多次审批操作生成多条记录
func TestRecordMultipleApprovals(t *testing.T) {
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

	// 获取任务,验证审批记录数量
	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	// 验证审批记录数量
	if len(tsk.Records) < 2 {
		t.Errorf("Expected at least 2 approval records, got %d", len(tsk.Records))
	}

	// 验证每条记录都有唯一 ID
	recordIDs := make(map[string]bool)
	for _, record := range tsk.Records {
		if recordIDs[record.ID] {
			t.Errorf("Duplicate record ID: %q", record.ID)
		}
		recordIDs[record.ID] = true
	}
}

// TestRecordDataStructure 测试审批记录数据结构
func TestRecordDataStructure(t *testing.T) {
	record := &task.Record{
		ID:         "record-001",
		TaskID:     "task-001",
		NodeID:     "node-001",
		Approver:   "user-001",
		Result:     "approve",
		Comment:    "approved",
		CreatedAt:  time.Now(),
		Attachments: []string{"file1.pdf", "file2.pdf"},
	}

	// 验证所有字段都已设置
	if record.ID == "" {
		t.Error("Record.ID should not be empty")
	}

	if record.TaskID == "" {
		t.Error("Record.TaskID should not be empty")
	}

	if record.NodeID == "" {
		t.Error("Record.NodeID should not be empty")
	}

	if record.Approver == "" {
		t.Error("Record.Approver should not be empty")
	}

	if record.Result == "" {
		t.Error("Record.Result should not be empty")
	}

	if record.CreatedAt.IsZero() {
		t.Error("Record.CreatedAt should not be zero")
	}

	if len(record.Attachments) != 2 {
		t.Errorf("Record.Attachments length = %d, want %d", len(record.Attachments), 2)
	}
}
