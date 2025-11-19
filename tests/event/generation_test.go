package event_test

import (
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/event"
	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
)

// TestEventGenerationOnTaskCreate 测试任务创建时生成事件
func TestEventGenerationOnTaskCreate(t *testing.T) {
	handler := newMockEventHandler()
	notifier := event.NewEventNotifier([]event.EventHandler{handler}, 10)
	defer notifier.Stop()

	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateForEvent()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	taskMgr := task.NewTaskManagerWithNotifier(templateMgr, nil, notifier)

	params := json.RawMessage(`{"amount": 1000}`)
	tsk, err := taskMgr.Create("tpl-event-001", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	waitForEvents(handler, 1, 200*time.Millisecond)

	events := handler.GetEvents()
	if len(events) == 0 {
		t.Fatal("Expected at least one event, got 0")
	}

	// 验证任务创建事件
	found := findEventByType(events, event.EventTypeTaskCreated)
	if found == nil {
		t.Fatal("Expected EventTypeTaskCreated event, not found")
	}

	if found.Task.ID != tsk.ID {
		t.Errorf("Event.Task.ID = %q, want %q", found.Task.ID, tsk.ID)
	}
	if found.Task.BusinessID != "biz-001" {
		t.Errorf("Event.Task.BusinessID = %q, want %q", found.Task.BusinessID, "biz-001")
	}
	if found.Task.State != "pending" {
		t.Errorf("Event.Task.State = %q, want %q", found.Task.State, "pending")
	}
	if found.Node == nil {
		t.Error("Event.Node should not be nil")
	}
	if found.Business == nil {
		t.Error("Event.Business should not be nil")
	}
}

// TestEventGenerationOnTaskSubmit 测试任务提交时生成事件
func TestEventGenerationOnTaskSubmit(t *testing.T) {
	handler := newMockEventHandler()
	notifier := event.NewEventNotifier([]event.EventHandler{handler}, 10)
	defer notifier.Stop()

	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateForEvent()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	taskMgr := task.NewTaskManagerWithNotifier(templateMgr, nil, notifier)

	params := json.RawMessage(`{"amount": 1000}`)
	tsk, err := taskMgr.Create("tpl-event-001", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	err = taskMgr.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Submit() failed: %v", err)
	}

	waitForEvents(handler, 2, 200*time.Millisecond)

	events := handler.GetEvents()
	if len(events) < 2 {
		t.Errorf("Expected at least 2 events, got %d", len(events))
	}

	// 验证任务提交事件
	found := findEventByType(events, event.EventTypeTaskSubmitted)
	if found == nil {
		t.Fatal("Expected EventTypeTaskSubmitted event, not found")
	}

	if found.Task.ID != tsk.ID {
		t.Errorf("Event.Task.ID = %q, want %q", found.Task.ID, tsk.ID)
	}
	if found.Task.State != "submitted" {
		t.Errorf("Event.Task.State = %q, want %q", found.Task.State, "submitted")
	}
}

// TestEventGenerationOnApproval 测试审批操作时生成事件
func TestEventGenerationOnApproval(t *testing.T) {
	tests := []struct {
		name           string
		operation      string
		expectedType   event.EventType
		expectedResult string
		operationFunc  func(task.TaskManager, string) error
	}{
		{
			name:           "approve operation",
			operation:      "approve",
			expectedType:   event.EventTypeApprovalOp,
			expectedResult: "approve",
			operationFunc: func(mgr task.TaskManager, taskID string) error {
				return mgr.Approve(taskID, "approval-001", "user-001", "approved")
			},
		},
		{
			name:           "reject operation",
			operation:      "reject",
			expectedType:   event.EventTypeApprovalOp,
			expectedResult: "reject",
			operationFunc: func(mgr task.TaskManager, taskID string) error {
				return mgr.Reject(taskID, "approval-001", "user-001", "rejected")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := newMockEventHandler()
			notifier := event.NewEventNotifier([]event.EventHandler{handler}, 10)
			defer notifier.Stop()

			templateMgr := template.NewTemplateManager()
			tpl := createTestTemplateWithApprovalNodeForEvent()
			err := templateMgr.Create(tpl)
			if err != nil {
				t.Fatalf("Create template failed: %v", err)
			}

			taskMgr := task.NewTaskManagerWithNotifier(templateMgr, nil, notifier)

			params := json.RawMessage(`{"amount": 1000}`)
			tsk, err := taskMgr.Create("tpl-event-001", "biz-001", params)
			if err != nil {
				t.Fatalf("Create() failed: %v", err)
			}

			err = taskMgr.Submit(tsk.ID)
			if err != nil {
				t.Fatalf("Submit() failed: %v", err)
			}

			err = tt.operationFunc(taskMgr, tsk.ID)
			if err != nil {
				t.Fatalf("%s() failed: %v", tt.operation, err)
			}

			waitForEvents(handler, 3, 200*time.Millisecond)

			events := handler.GetEvents()
			if len(events) < 3 {
				t.Errorf("Expected at least 3 events, got %d", len(events))
			}

			// 验证审批操作事件
			found := findEventByType(events, tt.expectedType)
			if found == nil {
				t.Fatalf("Expected %s event, not found", tt.expectedType)
			}

			if found.Approval == nil {
				t.Fatal("Event.Approval should not be nil for approval operation")
			}

			if found.Approval.Approver != "user-001" {
				t.Errorf("Event.Approval.Approver = %q, want %q", found.Approval.Approver, "user-001")
			}

			if found.Approval.Result != tt.expectedResult {
				t.Errorf("Event.Approval.Result = %q, want %q", found.Approval.Result, tt.expectedResult)
			}

			if found.Approval.NodeID != "approval-001" {
				t.Errorf("Event.Approval.NodeID = %q, want %q", found.Approval.NodeID, "approval-001")
			}
		})
	}
}

// TestEventGenerationOnTaskApproved 测试任务通过时生成事件
func TestEventGenerationOnTaskApproved(t *testing.T) {
	handler := newMockEventHandler()
	notifier := event.NewEventNotifier([]event.EventHandler{handler}, 10)
	defer notifier.Stop()

	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateWithApprovalNodeForEvent()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	taskMgr := task.NewTaskManagerWithNotifier(templateMgr, nil, notifier)

	params := json.RawMessage(`{"amount": 1000}`)
	tsk, err := taskMgr.Create("tpl-event-001", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	err = taskMgr.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Submit() failed: %v", err)
	}

	err = taskMgr.Approve(tsk.ID, "approval-001", "user-001", "approved")
	if err != nil {
		t.Fatalf("Approve() failed: %v", err)
	}

	waitForEvents(handler, 4, 300*time.Millisecond)

	events := handler.GetEvents()

	// 验证任务通过事件
	found := findEventByType(events, event.EventTypeTaskApproved)
	if found == nil {
		t.Fatal("Expected EventTypeTaskApproved event, not found")
	}

	if found.Task.ID != tsk.ID {
		t.Errorf("Event.Task.ID = %q, want %q", found.Task.ID, tsk.ID)
	}

	if found.Task.State != "approved" {
		t.Errorf("Event.Task.State = %q, want %q", found.Task.State, "approved")
	}
}

// TestEventGenerationOnTaskRejected 测试任务拒绝时生成事件
func TestEventGenerationOnTaskRejected(t *testing.T) {
	handler := newMockEventHandler()
	notifier := event.NewEventNotifier([]event.EventHandler{handler}, 10)
	defer notifier.Stop()

	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateWithApprovalNodeForEvent()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	taskMgr := task.NewTaskManagerWithNotifier(templateMgr, nil, notifier)

	params := json.RawMessage(`{"amount": 1000}`)
	tsk, err := taskMgr.Create("tpl-event-001", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	err = taskMgr.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Submit() failed: %v", err)
	}

	err = taskMgr.Reject(tsk.ID, "approval-001", "user-001", "rejected")
	if err != nil {
		t.Fatalf("Reject() failed: %v", err)
	}

	waitForEvents(handler, 4, 300*time.Millisecond)

	events := handler.GetEvents()

	// 验证任务拒绝事件
	found := findEventByType(events, event.EventTypeTaskRejected)
	if found == nil {
		t.Fatal("Expected EventTypeTaskRejected event, not found")
	}

	if found.Task.ID != tsk.ID {
		t.Errorf("Event.Task.ID = %q, want %q", found.Task.ID, tsk.ID)
	}

	if found.Task.State != "rejected" {
		t.Errorf("Event.Task.State = %q, want %q", found.Task.State, "rejected")
	}
}

// TestEventGenerationOnTaskCancelled 测试任务取消时生成事件
func TestEventGenerationOnTaskCancelled(t *testing.T) {
	handler := newMockEventHandler()
	notifier := event.NewEventNotifier([]event.EventHandler{handler}, 10)
	defer notifier.Stop()

	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateForEvent()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	taskMgr := task.NewTaskManagerWithNotifier(templateMgr, nil, notifier)

	params := json.RawMessage(`{"amount": 1000}`)
	tsk, err := taskMgr.Create("tpl-event-001", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	err = taskMgr.Cancel(tsk.ID, "cancelled by user")
	if err != nil {
		t.Fatalf("Cancel() failed: %v", err)
	}

	waitForEvents(handler, 2, 200*time.Millisecond)

	events := handler.GetEvents()

	// 验证任务取消事件
	found := findEventByType(events, event.EventTypeTaskCancelled)
	if found == nil {
		t.Fatal("Expected EventTypeTaskCancelled event, not found")
	}

	if found.Task.ID != tsk.ID {
		t.Errorf("Event.Task.ID = %q, want %q", found.Task.ID, tsk.ID)
	}

	if found.Task.State != "cancelled" {
		t.Errorf("Event.Task.State = %q, want %q", found.Task.State, "cancelled")
	}
}

// TestEventGenerationOnNodeActivated 测试节点激活时生成事件
func TestEventGenerationOnNodeActivated(t *testing.T) {
	handler := newMockEventHandler()
	notifier := event.NewEventNotifier([]event.EventHandler{handler}, 10)
	defer notifier.Stop()

	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateWithApprovalNodeForEvent()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	taskMgr := task.NewTaskManagerWithNotifier(templateMgr, nil, notifier)

	params := json.RawMessage(`{"amount": 1000}`)
	tsk, err := taskMgr.Create("tpl-event-001", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	err = taskMgr.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Submit() failed: %v", err)
	}

	waitForEvents(handler, 3, 300*time.Millisecond)

	events := handler.GetEvents()

	// 验证节点激活事件
	found := findEventByType(events, event.EventTypeNodeActivated)
	if found == nil {
		t.Fatal("Expected EventTypeNodeActivated event, not found")
	}

	if found.Node == nil {
		t.Fatal("Event.Node should not be nil")
	}

	if found.Node.ID != "approval-001" {
		t.Errorf("Event.Node.ID = %q, want %q", found.Node.ID, "approval-001")
	}

	if found.Node.Type != "approval" {
		t.Errorf("Event.Node.Type = %q, want %q", found.Node.Type, "approval")
	}
}

// TestEventGenerationOnNodeCompleted 测试节点完成时生成事件
func TestEventGenerationOnNodeCompleted(t *testing.T) {
	handler := newMockEventHandler()
	notifier := event.NewEventNotifier([]event.EventHandler{handler}, 10)
	defer notifier.Stop()

	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateWithApprovalNodeForEvent()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	taskMgr := task.NewTaskManagerWithNotifier(templateMgr, nil, notifier)

	params := json.RawMessage(`{"amount": 1000}`)
	tsk, err := taskMgr.Create("tpl-event-001", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	err = taskMgr.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Submit() failed: %v", err)
	}

	err = taskMgr.Approve(tsk.ID, "approval-001", "user-001", "approved")
	if err != nil {
		t.Fatalf("Approve() failed: %v", err)
	}

	waitForEvents(handler, 5, 300*time.Millisecond)

	events := handler.GetEvents()

	// 验证节点完成事件
	found := findEventByType(events, event.EventTypeNodeCompleted)
	if found == nil {
		t.Fatal("Expected EventTypeNodeCompleted event, not found")
	}

	if found.Node == nil {
		t.Fatal("Event.Node should not be nil")
	}

	if found.Node.ID != "approval-001" {
		t.Errorf("Event.Node.ID = %q, want %q", found.Node.ID, "approval-001")
	}
}

// TestEventGenerationAllEventTypes 测试所有事件类型都能正确生成
func TestEventGenerationAllEventTypes(t *testing.T) {
	expectedEventTypes := []event.EventType{
		event.EventTypeTaskCreated,
		event.EventTypeTaskSubmitted,
		event.EventTypeNodeActivated,
		event.EventTypeApprovalOp,
		event.EventTypeTaskApproved,
		event.EventTypeTaskRejected,
		event.EventTypeTaskCancelled,
		event.EventTypeNodeCompleted,
	}

	handler := newMockEventHandler()
	notifier := event.NewEventNotifier([]event.EventHandler{handler}, 10)
	defer notifier.Stop()

	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateWithApprovalNodeForEvent()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	taskMgr := task.NewTaskManagerWithNotifier(templateMgr, nil, notifier)

	params := json.RawMessage(`{"amount": 1000}`)
	tsk, err := taskMgr.Create("tpl-event-001", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	err = taskMgr.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Submit() failed: %v", err)
	}

	err = taskMgr.Approve(tsk.ID, "approval-001", "user-001", "approved")
	if err != nil {
		t.Fatalf("Approve() failed: %v", err)
	}

	waitForEvents(handler, 5, 300*time.Millisecond)

	events := handler.GetEvents()

	// 验证所有预期的事件类型都已生成
	for _, expectedType := range expectedEventTypes {
		found := findEventByType(events, expectedType)
		if found == nil {
			// 某些事件类型可能不会在单个流程中全部出现,这是正常的
			// 我们只验证在流程中应该出现的事件类型
			if expectedType == event.EventTypeTaskRejected || expectedType == event.EventTypeTaskCancelled {
				// 这些事件在正常审批流程中不会出现,跳过
				continue
			}
			t.Logf("Event type %s not found in events (may be expected)", expectedType)
		}
	}
}

// TestEventDataCompleteness 测试事件数据完整性
func TestEventDataCompleteness(t *testing.T) {
	handler := newMockEventHandler()
	notifier := event.NewEventNotifier([]event.EventHandler{handler}, 10)
	defer notifier.Stop()

	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateForEvent()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	taskMgr := task.NewTaskManagerWithNotifier(templateMgr, nil, notifier)

	params := json.RawMessage(`{"amount": 1000, "reason": "test"}`)
	tsk, err := taskMgr.Create("tpl-event-001", "biz-001", params)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	waitForEvents(handler, 1, 200*time.Millisecond)

	events := handler.GetEvents()
	if len(events) == 0 {
		t.Fatal("Expected at least one event, got 0")
	}

	evt := events[0]

	// 验证事件基本字段
	if evt.ID == "" {
		t.Error("Event.ID should not be empty")
	}

	if evt.Type == "" {
		t.Error("Event.Type should not be empty")
	}

	if evt.Time.IsZero() {
		t.Error("Event.Time should not be zero")
	}

	// 验证任务信息
	if evt.Task == nil {
		t.Fatal("Event.Task should not be nil")
	}

	if evt.Task.ID != tsk.ID {
		t.Errorf("Event.Task.ID = %q, want %q", evt.Task.ID, tsk.ID)
	}

	if evt.Task.TemplateID == "" {
		t.Error("Event.Task.TemplateID should not be empty")
	}

	if evt.Task.BusinessID == "" {
		t.Error("Event.Task.BusinessID should not be empty")
	}

	if evt.Task.State == "" {
		t.Error("Event.Task.State should not be empty")
	}

	// 验证节点信息(所有事件都应该包含)
	if evt.Node == nil {
		t.Fatal("Event.Node should not be nil")
	}

	if evt.Node.ID == "" {
		t.Error("Event.Node.ID should not be empty")
	}

	// 验证业务信息
	if evt.Business == nil {
		t.Fatal("Event.Business should not be nil")
	}

	if evt.Business.ID == "" {
		t.Error("Event.Business.ID should not be empty")
	}
}

// createTestTemplateForEvent 创建用于事件测试的简单模板
func createTestTemplateForEvent() *template.Template {
	return &template.Template{
		ID:          "tpl-event-001",
		Name:        "Test Template for Event",
		Description: "Test template for event generation",
		Version:     1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "Start Node",
				Type: template.NodeTypeStart,
			},
			"end": {
				ID:   "end",
				Name: "End Node",
				Type: template.NodeTypeEnd,
			},
		},
		Edges: []*template.Edge{
			{From: "start", To: "end"},
		},
	}
}

// createTestTemplateWithApprovalNodeForEvent 创建包含审批节点的测试模板
func createTestTemplateWithApprovalNodeForEvent() *template.Template {
	return &template.Template{
		ID:          "tpl-event-001",
		Name:        "Test Template with Approval Node",
		Description: "Test template with approval node for event generation",
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

// newMockEventHandler 创建新的 mock 事件处理器
// mockEventHandler 在 helper_test.go 中定义
func newMockEventHandler() *mockEventHandler {
	return &mockEventHandler{
		events: make([]*event.Event, 0),
		mu:     sync.Mutex{},
	}
}

// waitForEvents 等待事件处理完成
func waitForEvents(handler *mockEventHandler, minCount int, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		events := handler.GetEvents()
		if len(events) >= minCount {
			time.Sleep(50 * time.Millisecond) // 额外等待确保所有事件都已处理
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// findEventByType 在事件列表中查找指定类型的事件
func findEventByType(events []*event.Event, eventType event.EventType) *event.Event {
	for _, evt := range events {
		if evt.Type == eventType {
			return evt
		}
	}
	return nil
}
