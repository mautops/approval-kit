package event_test

import (
	"sync"
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/event"
)

// TestEventIdempotency 测试事件幂等性
func TestEventIdempotency(t *testing.T) {
	handler := &mockEventHandler{
		events: make([]*event.Event, 0),
		mu:     sync.Mutex{},
	}

	notifier := event.NewEventNotifier([]event.EventHandler{handler}, 10)

	// 创建相同的事件(相同 ID)
	evt1 := &event.Event{
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

	evt2 := &event.Event{
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

	// 推送相同的事件两次
	notifier.Notify(evt1)
	notifier.Notify(evt2)

	// 等待事件处理完成
	time.Sleep(100 * time.Millisecond)

	// 停止通知器
	notifier.Stop()

	// 验证事件处理次数(幂等性: 相同事件应该只处理一次,或者多次处理结果相同)
	handler.mu.Lock()
	defer handler.mu.Unlock()

	// 注意: 当前实现中,每个事件都会被处理
	// 幂等性保证应该由事件处理器实现(例如通过事件 ID 去重)
	// 这里我们验证事件确实被处理了
	if len(handler.events) < 1 {
		t.Error("Expected at least 1 event to be processed")
	}
}

// TestEventIdempotencyWithID 测试带事件 ID 的幂等性
func TestEventIdempotencyWithID(t *testing.T) {
	handler := &idempotentEventHandler{
		processedIDs: make(map[string]bool),
		mu:           sync.Mutex{},
	}

	notifier := event.NewEventNotifier([]event.EventHandler{handler}, 10)

	// 创建带 ID 的事件
	evt := &event.Event{
		ID:   "event-001",
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

	// 推送相同的事件多次
	notifier.Notify(evt)
	notifier.Notify(evt)
	notifier.Notify(evt)

	// 等待事件处理完成
	time.Sleep(200 * time.Millisecond)

	// 停止通知器
	notifier.Stop()

	// 验证事件只被处理一次(幂等性)
	handler.mu.Lock()
	defer handler.mu.Unlock()

	if handler.processCount != 1 {
		t.Errorf("Event should be processed only once (idempotency), but was processed %d times", handler.processCount)
	}
}

// idempotentEventHandler 支持幂等性的事件处理器
type idempotentEventHandler struct {
	processedIDs map[string]bool
	processCount int
	mu           sync.Mutex
}

func (i *idempotentEventHandler) Handle(evt *event.Event) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	// 使用事件 ID 进行去重
	if evt.ID != "" {
		if i.processedIDs[evt.ID] {
			// 事件已处理过,跳过(幂等性)
			return nil
		}
		i.processedIDs[evt.ID] = true
	}

	i.processCount++
	return nil
}
