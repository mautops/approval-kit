package event_test

import (
	"sync"

	"github.com/mautops/approval-kit/internal/event"
)

// mockEventHandler 用于测试的模拟事件处理器
type mockEventHandler struct {
	events []*event.Event
	mu     sync.Mutex
}

func (m *mockEventHandler) Handle(evt *event.Event) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = append(m.events, evt)
	return nil
}

func (m *mockEventHandler) GetEvents() []*event.Event {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.events
}

func (m *mockEventHandler) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = make([]*event.Event, 0)
}
