package event

import (
	"errors"
	"log"
	"sync"
	"time"
)

// EventNotifier 事件通知器
// 使用 channel 和 goroutine 实现异步事件推送,不阻塞主流程
type EventNotifier struct {
	handlers []EventHandler
	queue    chan *Event
	wg       sync.WaitGroup
	stop     chan struct{}
	once     sync.Once
}

// NewEventNotifier 创建新的事件通知器
// handlers: 事件处理器列表
// queueSize: 事件队列大小
func NewEventNotifier(handlers []EventHandler, queueSize int) *EventNotifier {
	if queueSize <= 0 {
		queueSize = 100 // 默认队列大小
	}

	notifier := &EventNotifier{
		handlers: handlers,
		queue:    make(chan *Event, queueSize),
		stop:     make(chan struct{}),
	}

	// 启动 worker goroutine
	notifier.wg.Add(1)
	go notifier.worker()

	return notifier
}

// Notify 异步推送事件
// 如果队列满,记录日志但不阻塞
func (n *EventNotifier) Notify(evt *Event) {
	select {
	case n.queue <- evt:
		// 事件成功入队
	default:
		// 队列满时记录日志,不阻塞
		log.Printf("event queue full, dropping event: type=%q, task=%q", evt.Type, evt.Task.ID)
	}
}

// worker 事件处理 worker
func (n *EventNotifier) worker() {
	defer n.wg.Done()

	for {
		select {
		case evt := <-n.queue:
			n.pushEvent(evt)
		case <-n.stop:
			return
		}
	}
}

// pushEvent 推送事件到所有处理器
func (n *EventNotifier) pushEvent(evt *Event) {
	for _, handler := range n.handlers {
		go func(h EventHandler) {
			if err := n.pushWithRetry(h, evt); err != nil {
				log.Printf("failed to push event: type=%q, task=%q, error=%v", evt.Type, evt.Task.ID, err)
			}
		}(handler)
	}
}

// pushWithRetry 带重试的推送
// 使用指数退避策略进行重试
func (n *EventNotifier) pushWithRetry(handler EventHandler, evt *Event) error {
	maxRetries := 3
	backoff := time.Second // 初始退避时间

	for i := 0; i < maxRetries; i++ {
		if err := handler.Handle(evt); err == nil {
			return nil
		}
		if i < maxRetries-1 {
			// 指数退避
			time.Sleep(backoff)
			backoff *= 2
		}
	}
	return ErrEventPushFailed
}

// Stop 停止事件通知器
func (n *EventNotifier) Stop() {
	n.once.Do(func() {
		close(n.stop)
		n.wg.Wait()
	})
}

// ErrEventPushFailed 事件推送失败错误
var ErrEventPushFailed = errors.New("event push failed after retries")
