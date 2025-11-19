package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// WebhookConfig Webhook 配置
type WebhookConfig struct {
	// URL Webhook 地址
	URL string

	// Method HTTP 方法(默认 POST)
	Method string

	// Headers HTTP 请求头(用于认证等)
	Headers map[string]string

	// Timeout 超时时间(秒,默认 30)
	Timeout int
}

// WebhookHandler Webhook 事件处理器
// 实现 EventHandler 接口,将事件推送到 Webhook URL
type WebhookHandler struct {
	config     *WebhookConfig
	httpClient *http.Client
}

// NewWebhookHandler 创建新的 Webhook 处理器
func NewWebhookHandler(config *WebhookConfig) EventHandler {
	if config.Method == "" {
		config.Method = "POST"
	}
	if config.Timeout <= 0 {
		config.Timeout = 30
	}

	return &WebhookHandler{
		config: config,
		httpClient: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
	}
}

// Handle 处理事件(实现 EventHandler 接口)
func (h *WebhookHandler) Handle(evt *Event) error {
	// 序列化事件为 JSON
	data, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest(h.config.Method, h.config.URL, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	for k, v := range h.config.Headers {
		req.Header.Set(k, v)
	}

	// 发送请求
	resp, err := h.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("webhook returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

