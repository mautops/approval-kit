package event_test

import (
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/event"
)

// TestEventHandlerInterface 测试 EventHandler 接口
func TestEventHandlerInterface(t *testing.T) {
	// 验证接口类型存在(通过编译时检查,如果接口不存在会编译失败)
	var handler event.EventHandler
	_ = handler // 使用变量避免未使用变量警告
}

// TestEventHandlerMethods 测试 EventHandler 接口方法签名
func TestEventHandlerMethods(t *testing.T) {
	// 验证接口包含所有必需的方法
	// 通过编译时检查,如果方法不存在会编译失败
	var _ event.EventHandler = (*mockEventHandler)(nil)
}

// mockEventHandler 在 helper.go 中定义

// TestEventHandlerHandle 测试事件处理
func TestEventHandlerHandle(t *testing.T) {
	handler := &mockEventHandler{
		events: make([]*event.Event, 0),
	}

	evt := &event.Event{
		Type: event.EventTypeTaskCreated,
		Time: time.Now(),
		Task: &event.TaskInfo{
			ID:         "task-001",
			TemplateID: "tpl-001",
			BusinessID: "biz-001",
			State:      "pending",
		},
		Node: &event.NodeInfo{
			ID:   "node-001",
			Name: "Start Node",
			Type: "start",
		},
		Business: &event.BusinessInfo{
			ID: "biz-001",
		},
	}

	err := handler.Handle(evt)
	if err != nil {
		t.Fatalf("Handle() failed: %v", err)
	}

	if len(handler.events) != 1 {
		t.Errorf("Handle() should add 1 event, got %d", len(handler.events))
	}

	if handler.events[0].Type != event.EventTypeTaskCreated {
		t.Errorf("Handle() event type = %q, want %q", handler.events[0].Type, event.EventTypeTaskCreated)
	}
}
