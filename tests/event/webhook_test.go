package event_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/event"
)

// TestWebhookHandler 测试 Webhook 事件处理器
func TestWebhookHandler(t *testing.T) {
	// 创建测试服务器
	var receivedEvents []*event.Event
	var mu sync.Mutex

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var evt event.Event
		if err := json.NewDecoder(r.Body).Decode(&evt); err != nil {
			t.Errorf("Failed to decode event: %v", err)
			return
		}

		mu.Lock()
		receivedEvents = append(receivedEvents, &evt)
		mu.Unlock()

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// 创建 Webhook 配置
	config := &event.WebhookConfig{
		URL:    server.URL,
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	// 创建 Webhook 处理器
	handler := event.NewWebhookHandler(config)

	// 创建事件
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

	// 处理事件
	err := handler.Handle(evt)
	if err != nil {
		t.Fatalf("Handle() failed: %v", err)
	}

	// 等待请求完成
	time.Sleep(100 * time.Millisecond)

	// 验证事件已发送
	mu.Lock()
	defer mu.Unlock()

	if len(receivedEvents) != 1 {
		t.Errorf("Expected 1 event, got %d", len(receivedEvents))
	}

	if receivedEvents[0].Type != event.EventTypeTaskCreated {
		t.Errorf("Event type = %q, want %q", receivedEvents[0].Type, event.EventTypeTaskCreated)
	}
}

// TestWebhookHandlerWithAuth 测试带认证的 Webhook
func TestWebhookHandlerWithAuth(t *testing.T) {
	// 创建测试服务器
	var receivedToken string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedToken = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// 创建带认证的 Webhook 配置
	config := &event.WebhookConfig{
		URL:    server.URL,
		Method: "POST",
		Headers: map[string]string{
			"Authorization": "Bearer test-token",
			"Content-Type":  "application/json",
		},
	}

	// 创建 Webhook 处理器
	handler := event.NewWebhookHandler(config)

	// 创建事件
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

	// 处理事件
	err := handler.Handle(evt)
	if err != nil {
		t.Fatalf("Handle() failed: %v", err)
	}

	// 等待请求完成
	time.Sleep(100 * time.Millisecond)

	// 验证认证头已发送
	if receivedToken != "Bearer test-token" {
		t.Errorf("Authorization header = %q, want %q", receivedToken, "Bearer test-token")
	}
}

// TestWebhookHandlerError 测试 Webhook 错误处理
func TestWebhookHandlerError(t *testing.T) {
	// 创建返回错误的测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// 创建 Webhook 配置
	config := &event.WebhookConfig{
		URL:    server.URL,
		Method: "POST",
	}

	// 创建 Webhook 处理器
	handler := event.NewWebhookHandler(config)

	// 创建事件
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

	// 处理事件(应该返回错误)
	err := handler.Handle(evt)
	if err == nil {
		t.Error("Handle() should return error for 500 status")
	}
}
