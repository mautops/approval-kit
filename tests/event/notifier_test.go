package event_test

import (
	"sync"
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/event"
)

// mockEventHandler 在 helper.go 中定义

// TestEventNotifier 测试事件通知器
func TestEventNotifier(t *testing.T) {
	handler := &mockEventHandler{
		events: make([]*event.Event, 0),
		mu:     sync.Mutex{},
	}

	notifier := event.NewEventNotifier([]event.EventHandler{handler}, 10)

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

	// 异步推送事件
	notifier.Notify(evt)

	// 等待事件处理完成
	time.Sleep(100 * time.Millisecond)

	// 停止通知器
	notifier.Stop()

	// 验证事件已被处理
	handler.mu.Lock()
	defer handler.mu.Unlock()

	if len(handler.events) != 1 {
		t.Errorf("EventNotifier should handle 1 event, got %d", len(handler.events))
	}
}

// TestEventNotifierMultipleHandlers 测试多个事件处理器
func TestEventNotifierMultipleHandlers(t *testing.T) {
	handler1 := &mockEventHandler{
		events: make([]*event.Event, 0),
		mu:     sync.Mutex{},
	}
	handler2 := &mockEventHandler{
		events: make([]*event.Event, 0),
		mu:     sync.Mutex{},
	}

	notifier := event.NewEventNotifier([]event.EventHandler{handler1, handler2}, 10)

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

	// 异步推送事件
	notifier.Notify(evt)

	// 等待事件处理完成
	time.Sleep(100 * time.Millisecond)

	// 停止通知器
	notifier.Stop()

	// 验证两个处理器都收到了事件
	handler1.mu.Lock()
	handler1Count := len(handler1.events)
	handler1.mu.Unlock()

	handler2.mu.Lock()
	handler2Count := len(handler2.events)
	handler2.mu.Unlock()

	if handler1Count != 1 {
		t.Errorf("Handler1 should receive 1 event, got %d", handler1Count)
	}

	if handler2Count != 1 {
		t.Errorf("Handler2 should receive 1 event, got %d", handler2Count)
	}
}

// TestEventNotifierNonBlocking 测试非阻塞推送
func TestEventNotifierNonBlocking(t *testing.T) {
	handler := &slowEventHandler{
		events: make([]*event.Event, 0),
		mu:     sync.Mutex{},
		delay:  50 * time.Millisecond,
	}

	notifier := event.NewEventNotifier([]event.EventHandler{handler}, 2)

	// 推送多个事件(队列大小为 2)
	for i := 0; i < 5; i++ {
		evt := &event.Event{
			Type: event.EventTypeTaskCreated,
			Time: time.Now(),
			Task: &event.TaskInfo{
				ID: "task-001",
			},
			Node: &event.NodeInfo{
				ID: "node-001",
			},
			Business: &event.BusinessInfo{
				ID: "biz-001",
			},
		}
		notifier.Notify(evt)
	}

	// 等待事件处理完成
	time.Sleep(500 * time.Millisecond)

	// 停止通知器
	notifier.Stop()

	// 验证事件已被处理(可能部分事件被丢弃,因为队列满)
	handler.mu.Lock()
	defer handler.mu.Unlock()

	if len(handler.events) == 0 {
		t.Error("EventNotifier should handle at least some events")
	}
}

// slowEventHandler 慢速事件处理器(用于测试非阻塞)
type slowEventHandler struct {
	events []*event.Event
	mu     sync.Mutex
	delay  time.Duration
}

func (s *slowEventHandler) Handle(evt *event.Event) error {
	time.Sleep(s.delay)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, evt)
	return nil
}
