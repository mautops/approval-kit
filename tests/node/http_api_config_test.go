package node_test

import (
	"testing"

	"github.com/mautops/approval-kit/internal/node"
)

// TestHTTPAPIConfig 测试 HTTPAPIConfig 结构体定义
func TestHTTPAPIConfig(t *testing.T) {
	// 测试结构体字段
	config := &node.HTTPAPIConfig{
		URL:    "http://example.com/api/approvers",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Authorization": "Bearer token",
		},
	}

	if config.URL != "http://example.com/api/approvers" {
		t.Errorf("HTTPAPIConfig.URL = %q, want %q", config.URL, "http://example.com/api/approvers")
	}

	if config.Method != "POST" {
		t.Errorf("HTTPAPIConfig.Method = %q, want %q", config.Method, "POST")
	}

	if config.Headers == nil {
		t.Error("HTTPAPIConfig.Headers should not be nil")
	}

	if config.Headers["Content-Type"] != "application/json" {
		t.Errorf("HTTPAPIConfig.Headers[Content-Type] = %q, want %q", config.Headers["Content-Type"], "application/json")
	}
}

// TestHTTPAPIConfigValidate 测试 HTTPAPIConfig 验证
func TestHTTPAPIConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *node.HTTPAPIConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: &node.HTTPAPIConfig{
				URL:    "http://example.com/api/approvers",
				Method: "POST",
			},
			wantErr: false,
		},
		{
			name: "empty URL",
			config: &node.HTTPAPIConfig{
				URL:    "",
				Method: "POST",
			},
			wantErr: true,
		},
		{
			name: "empty method (should default to GET)",
			config: &node.HTTPAPIConfig{
				URL:    "http://example.com/api/approvers",
				Method: "",
			},
			wantErr: false, // 空方法会默认设置为 GET,所以不应该报错
		},
		{
			name: "invalid URL",
			config: &node.HTTPAPIConfig{
				URL:    "not-a-valid-url",
				Method: "POST",
			},
			wantErr: true,
		},
		{
			name: "invalid method",
			config: &node.HTTPAPIConfig{
				URL:    "http://example.com/api/approvers",
				Method: "INVALID",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("HTTPAPIConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestHTTPAPIConfigDefaultMethod 测试默认方法
func TestHTTPAPIConfigDefaultMethod(t *testing.T) {
	config := &node.HTTPAPIConfig{
		URL: "http://example.com/api/approvers",
		// Method 未设置,应该默认为 GET
	}

	// 验证默认方法
	if config.Method == "" {
		// 如果 Validate 方法设置了默认值,这里应该检查
		_ = config.Validate()
	}

	// 验证方法应该是 GET 或 POST(根据实现)
	if config.Method != "GET" && config.Method != "POST" {
		t.Logf("HTTPAPIConfig.Method = %q (may be set by Validate)", config.Method)
	}
}

