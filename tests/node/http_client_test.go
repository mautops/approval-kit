package node_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/mautops/approval-kit/internal/node"
)

// TestHTTPClientInterface 验证 HTTPClient 接口定义
func TestHTTPClientInterface(t *testing.T) {
	// 验证接口类型存在
	var client node.HTTPClient
	if client != nil {
		_ = client
	}
}

// TestHTTPClientMethods 验证 HTTPClient 接口方法签名
func TestHTTPClientMethods(t *testing.T) {
	// 验证接口包含所有必需的方法
	// 通过编译时检查,如果方法不存在会编译失败
	var _ node.HTTPClient = (*httpClientImpl)(nil)
}

// httpClientImpl 用于测试接口方法签名的实现
type httpClientImpl struct{}

func (c *httpClientImpl) Do(req *http.Request) (*http.Response, error) {
	return nil, nil
}

// TestHTTPClientDo 测试 HTTPClient.Do 方法
func TestHTTPClientDo(t *testing.T) {
	// 创建一个 mock HTTPClient
	mockClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(io.Reader(nil)),
		},
		err: nil,
	}

	// 创建一个测试请求
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// 调用 Do 方法
	resp, err := mockClient.Do(req)
	if err != nil {
		t.Fatalf("HTTPClient.Do() failed: %v", err)
	}

	if resp == nil {
		t.Error("HTTPClient.Do() should return a response")
	}

	if resp.StatusCode != 200 {
		t.Errorf("HTTPClient.Do() returned status code %d, want 200", resp.StatusCode)
	}
}

// mockHTTPClient 用于测试的 HTTPClient 实现
type mockHTTPClient struct {
	response *http.Response
	err      error
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.response, m.err
}

