package task_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/task"
)

// TestTaskStruct 验证 Task 结构体定义
func TestTaskStruct(t *testing.T) {
	// 验证 Task 类型存在
	var tsk *task.Task
	if tsk != nil {
		_ = tsk
	}
}

// TestTaskFields 验证 Task 结构体的所有字段
func TestTaskFields(t *testing.T) {
	now := time.Now()
	params := json.RawMessage(`{"amount": 1000}`)

	tsk := &task.Task{
		ID:             "task-001",
		TemplateID:     "tpl-001",
		TemplateVersion: 1,
		BusinessID:     "biz-001",
		Params:         params,
		State:          task.TaskStatePending,
		CurrentNode:    "start",
		CreatedAt:      now,
		UpdatedAt:      now,
		SubmittedAt:    nil,
		NodeOutputs:    make(map[string]json.RawMessage),
		Approvers:      make(map[string][]string),
		Approvals:      make(map[string]map[string]*task.Approval),
		Records:        []*task.Record{},
		StateHistory:   []*task.StateChange{},
	}

	// 验证字段值
	if tsk.ID != "task-001" {
		t.Errorf("Task.ID = %q, want %q", tsk.ID, "task-001")
	}
	if tsk.TemplateID != "tpl-001" {
		t.Errorf("Task.TemplateID = %q, want %q", tsk.TemplateID, "tpl-001")
	}
	if tsk.TemplateVersion != 1 {
		t.Errorf("Task.TemplateVersion = %d, want %d", tsk.TemplateVersion, 1)
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
}

// TestTaskZeroValue 验证 Task 的零值
func TestTaskZeroValue(t *testing.T) {
	var tsk task.Task

	// 验证零值
	if tsk.ID != "" {
		t.Errorf("Task zero value ID = %q, want empty string", tsk.ID)
	}
	if tsk.State != "" {
		t.Errorf("Task zero value State = %q, want empty string", tsk.State)
	}
	// map 和 slice 的零值是 nil,这是正常的
	if tsk.NodeOutputs == nil {
		_ = tsk.NodeOutputs
	}
	if tsk.Approvers == nil {
		_ = tsk.Approvers
	}
	if tsk.Records == nil {
		_ = tsk.Records
	}
}

// TestTaskParams 验证任务参数字段
func TestTaskParams(t *testing.T) {
	params := json.RawMessage(`{"amount": 1000, "department": "IT"}`)
	tsk := &task.Task{
		ID:     "task-001",
		Params: params,
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

// TestTaskStateHistory 验证状态变更历史字段
func TestTaskStateHistory(t *testing.T) {
	tsk := &task.Task{
		ID:           "task-001",
		StateHistory: []*task.StateChange{},
	}

	// 添加状态变更记录
	change := &task.StateChange{
		From:    task.TaskStatePending,
		To:      task.TaskStateSubmitted,
		Reason:  "user submitted",
		Time:    time.Now(),
	}
	tsk.StateHistory = append(tsk.StateHistory, change)

	if len(tsk.StateHistory) != 1 {
		t.Errorf("Task.StateHistory length = %d, want 1", len(tsk.StateHistory))
	}
}

// TestTaskRecords 验证审批记录字段
func TestTaskRecords(t *testing.T) {
	tsk := &task.Task{
		ID:      "task-001",
		Records: []*task.Record{},
	}

	// 添加审批记录
	record := &task.Record{
		ID:        "record-001",
		TaskID:    "task-001",
		NodeID:    "node-001",
		Approver:  "user-001",
		Result:    "approve",
		Comment:   "approved",
		CreatedAt: time.Now(),
	}
	tsk.Records = append(tsk.Records, record)

	if len(tsk.Records) != 1 {
		t.Errorf("Task.Records length = %d, want 1", len(tsk.Records))
	}
}

// TestTaskRuntimeData 验证运行时数据字段
func TestTaskRuntimeData(t *testing.T) {
	tsk := &task.Task{
		ID:          "task-001",
		NodeOutputs: make(map[string]json.RawMessage),
		Approvers:   make(map[string][]string),
		Approvals:   make(map[string]map[string]*task.Approval),
	}

	// 添加节点输出
	output := json.RawMessage(`{"result": "success"}`)
	tsk.NodeOutputs["node-001"] = output

	// 添加审批人
	tsk.Approvers["node-001"] = []string{"user-001", "user-002"}

	// 添加审批结果
	tsk.Approvals["node-001"] = make(map[string]*task.Approval)
	tsk.Approvals["node-001"]["user-001"] = &task.Approval{
		Result:    "approve",
		Comment:   "approved",
		CreatedAt: time.Now(),
	}

	// 验证数据
	if len(tsk.NodeOutputs) != 1 {
		t.Errorf("Task.NodeOutputs length = %d, want 1", len(tsk.NodeOutputs))
	}
	if len(tsk.Approvers) != 1 {
		t.Errorf("Task.Approvers length = %d, want 1", len(tsk.Approvers))
	}
	if len(tsk.Approvals) != 1 {
		t.Errorf("Task.Approvals length = %d, want 1", len(tsk.Approvals))
	}
}

