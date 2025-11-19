package node

import (
	"fmt"
	"net/url"
	"strings"
)

// HTTPAPIConfig HTTP API 配置
// 用于动态审批人获取和条件节点评估
type HTTPAPIConfig struct {
	// URL API 地址
	URL string

	// Method 请求方法(GET/POST/PUT/DELETE 等)
	Method string

	// Headers 请求头
	Headers map[string]string

	// ParamMapping 参数映射规则
	// 定义如何从流程上下文中提取参数并映射到 API 请求参数
	ParamMapping *ParamMapping

	// ResponseMapping 响应数据解析规则
	// 定义如何从 API 响应中解析审批人列表
	ResponseMapping *ResponseMapping
}

// ParamMapping 参数映射规则
type ParamMapping struct {
	// Source 参数来源(task_params/node_outputs/context)
	Source string

	// Path 参数路径(JSONPath 或字段名)
	Path string

	// Target 目标参数名(API 请求参数名)
	Target string
}

// ResponseMapping 响应数据解析规则
type ResponseMapping struct {
	// Path 响应数据路径(JSONPath 或字段名)
	// 例如: "data.approvers" 或 "approvers"
	Path string

	// Format 响应格式(json)
	Format string
}

// Validate 验证 HTTPAPIConfig 配置
func (c *HTTPAPIConfig) Validate() error {
	// 验证 URL
	if c.URL == "" {
		return fmt.Errorf("HTTPAPIConfig.URL is required")
	}

	// 验证 URL 格式
	parsedURL, err := url.Parse(c.URL)
	if err != nil {
		return fmt.Errorf("HTTPAPIConfig.URL is invalid: %w", err)
	}

	// 验证 URL 必须包含 scheme
	if parsedURL.Scheme == "" {
		return fmt.Errorf("HTTPAPIConfig.URL must include scheme (http:// or https://)")
	}

	// 验证 scheme 必须是 http 或 https
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("HTTPAPIConfig.URL scheme must be http or https, got: %q", parsedURL.Scheme)
	}

	// 验证 Method
	if c.Method == "" {
		// 默认使用 GET
		c.Method = "GET"
	}

	// 验证 Method 值
	validMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	methodUpper := strings.ToUpper(c.Method)
	valid := false
	for _, m := range validMethods {
		if methodUpper == m {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("HTTPAPIConfig.Method %q is invalid, must be one of: %v", c.Method, validMethods)
	}

	// 标准化 Method
	c.Method = methodUpper

	// 验证 ResponseMapping
	if c.ResponseMapping != nil {
		if c.ResponseMapping.Path == "" {
			return fmt.Errorf("HTTPAPIConfig.ResponseMapping.Path is required")
		}
		if c.ResponseMapping.Format == "" {
			// 默认使用 json
			c.ResponseMapping.Format = "json"
		}
	}

	return nil
}

