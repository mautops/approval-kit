package event_test

import (
	"sync"
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/event"
)

// TestEventRetry 测试事件重试机制
func TestEventRetry(t *testing.T) {
	// 创建会失败的事件处理器(前两次失败,第三次成功)
	handler := &failingEventHandler{
		failCount:   2,
		currentFail: 0,
		mu:          sync.Mutex{},
		events:      make([]*event.Event, 0),
	}

	notifier := event.NewEventNotifier([]event.EventHandler{handler}, 10)

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

	// 推送事件
	notifier.Notify(evt)

	// 等待重试完成(最多重试 3 次,每次间隔 1s, 2s, 所以最多等待 4s)
	time.Sleep(5 * time.Second)

	// 停止通知器
	notifier.Stop()

	// 验证事件最终被处理成功
	handler.mu.Lock()
	defer handler.mu.Unlock()

	if len(handler.events) != 1 {
		t.Errorf("Event should be handled after retries, got %d events", len(handler.events))
	}

	// 验证重试次数
	if handler.currentFail < handler.failCount {
		t.Errorf("Handler should fail %d times, but only failed %d times", handler.failCount, handler.currentFail)
	}
}

// TestEventRetryFailure 测试重试后仍然失败的情况
func TestEventRetryFailure(t *testing.T) {
	// 创建总是失败的事件处理器
	handler := &alwaysFailingEventHandler{
		mu:     sync.Mutex{},
		events: make([]*event.Event, 0),
	}

	notifier := event.NewEventNotifier([]event.EventHandler{handler}, 10)

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

	// 推送事件
	notifier.Notify(evt)

	// 等待重试完成
	time.Sleep(5 * time.Second)

	// 停止通知器
	notifier.Stop()

	// 验证事件没有被处理(因为一直失败)
	handler.mu.Lock()
	defer handler.mu.Unlock()

	if len(handler.events) != 0 {
		t.Errorf("Event should not be handled after all retries failed, got %d events", len(handler.events))
	}

	// 验证重试了 3 次
	if handler.attemptCount != 3 {
		t.Errorf("Handler should be called 3 times (max retries), but was called %d times", handler.attemptCount)
	}
}

// failingEventHandler 会失败指定次数的事件处理器
type failingEventHandler struct {
	failCount   int
	currentFail int
	mu          sync.Mutex
	events      []*event.Event
}

func (f *failingEventHandler) Handle(evt *event.Event) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.currentFail++
	if f.currentFail <= f.failCount {
		return event.ErrEventPushFailed
	}

	f.events = append(f.events, evt)
	return nil
}

// alwaysFailingEventHandler 总是失败的事件处理器
type alwaysFailingEventHandler struct {
	mu           sync.Mutex
	events       []*event.Event
	attemptCount int
}

func (a *alwaysFailingEventHandler) Handle(evt *event.Event) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.attemptCount++
	return event.ErrEventPushFailed
}
