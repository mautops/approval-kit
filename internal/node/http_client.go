package node

import (
	"net/http"
)

// HTTPClient HTTP 客户端接口
// 用于动态审批人获取和条件节点评估
// 抽象标准库的 http.Client,便于测试和扩展
type HTTPClient interface {
	// Do 执行 HTTP 请求
	// req: HTTP 请求对象
	// 返回: HTTP 响应对象和错误
	Do(req *http.Request) (*http.Response, error)
}

// DefaultHTTPClient 默认 HTTP 客户端实现
// 使用标准库的 http.Client
type DefaultHTTPClient struct {
	client *http.Client
}

// NewDefaultHTTPClient 创建默认 HTTP 客户端
func NewDefaultHTTPClient() HTTPClient {
	return &DefaultHTTPClient{
		client: &http.Client{},
	}
}

// Do 执行 HTTP 请求
func (c *DefaultHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

