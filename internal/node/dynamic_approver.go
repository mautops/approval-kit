package node

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// DynamicApproverConfig 动态审批人配置
// 通过 HTTP API 动态获取审批人列表
type DynamicApproverConfig struct {
	// API HTTP API 配置
	API *HTTPAPIConfig

	// Timing 获取时机
	Timing ApproverTiming

	// HTTPClient HTTP 客户端(依赖注入)
	HTTPClient HTTPClient
}

// GetApprovers 获取审批人列表(实现 ApproverConfig 接口)
func (c *DynamicApproverConfig) GetApprovers(ctx *NodeContext) ([]string, error) {
	// 1. 验证配置
	if c.API == nil {
		return nil, fmt.Errorf("DynamicApproverConfig.API is required")
	}

	if err := c.API.Validate(); err != nil {
		return nil, fmt.Errorf("DynamicApproverConfig.API validation failed: %w", err)
	}

	if c.HTTPClient == nil {
		return nil, fmt.Errorf("DynamicApproverConfig.HTTPClient is required")
	}

	// 2. 构建 HTTP 请求
	req, err := c.buildRequest(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	// 3. 执行 HTTP 请求(带重试机制)
	resp, err := c.doWithRetry(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// 4. 检查响应状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	// 5. 解析响应数据
	approvers, err := c.parseResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return approvers, nil
}

// GetTiming 返回获取时机(实现 ApproverConfig 接口)
func (c *DynamicApproverConfig) GetTiming() ApproverTiming {
	if c.Timing == "" {
		// 默认值为节点激活时
		return ApproverTimingOnActivate
	}
	return c.Timing
}

// buildRequest 构建 HTTP 请求
func (c *DynamicApproverConfig) buildRequest(ctx *NodeContext) (*http.Request, error) {
	var req *http.Request
	var err error

	// 根据请求方法构建请求
	if c.API.Method == "GET" {
		// GET 请求: 参数放在 URL 中
		url := c.API.URL
		if c.API.ParamMapping != nil {
			// 添加查询参数
			params := c.mapParams(ctx)
			if len(params) > 0 {
				urlWithParams, err := c.addQueryParams(url, params)
				if err == nil {
					url = urlWithParams
				}
			}
		}
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
	} else {
		// POST/PUT/DELETE 等请求: 参数放在 Body 中
		body := c.buildRequestBody(ctx)
		req, err = http.NewRequest(c.API.Method, c.API.URL, body)
		if err != nil {
			return nil, err
		}
		// 设置 Content-Type
		if c.API.Headers == nil || c.API.Headers["Content-Type"] == "" {
			req.Header.Set("Content-Type", "application/json")
		}
	}

	// 设置请求头
	if c.API.Headers != nil {
		for key, value := range c.API.Headers {
			req.Header.Set(key, value)
		}
	}

	return req, nil
}

// buildRequestBody 构建请求体
func (c *DynamicApproverConfig) buildRequestBody(ctx *NodeContext) io.Reader {
	// 如果有参数映射配置,使用映射逻辑
	if c.API.ParamMapping != nil {
		params := c.mapParams(ctx)
		body, err := json.Marshal(params)
		if err != nil {
			// 如果序列化失败,返回空对象
			return strings.NewReader("{}")
		}
		return strings.NewReader(string(body))
	}

	// 如果没有参数映射配置,直接使用任务参数
	if ctx.Params != nil && len(ctx.Params) > 0 {
		return strings.NewReader(string(ctx.Params))
	}
	return strings.NewReader("{}")
}

// mapParams 根据参数映射规则映射参数
func (c *DynamicApproverConfig) mapParams(ctx *NodeContext) map[string]interface{} {
	result := make(map[string]interface{})

	// 如果只有一个参数映射,直接映射
	if c.API.ParamMapping != nil {
		value := c.getValueBySource(ctx, c.API.ParamMapping.Source, c.API.ParamMapping.Path)
		if value != nil {
			result[c.API.ParamMapping.Target] = value
		}
	}

	return result
}

// getValueBySource 根据数据源获取值
func (c *DynamicApproverConfig) getValueBySource(ctx *NodeContext, source, path string) interface{} {
	switch source {
	case "task_params":
		return c.getValueFromJSON(ctx.Params, path)
	case "node_outputs":
		// 从节点输出中获取值
		// path 格式: "node_id.field" 或 "node_id"
		parts := splitPath(path)
		if len(parts) == 0 {
			return nil
		}
		nodeID := parts[0]
		output, exists := ctx.Outputs[nodeID]
		if !exists {
			return nil
		}
		// 如果有子路径,继续解析
		if len(parts) > 1 {
			subPath := strings.Join(parts[1:], ".")
			return c.getValueFromJSON(output, subPath)
		}
		// 如果没有子路径,返回整个输出
		var data interface{}
		if err := json.Unmarshal(output, &data); err != nil {
			return nil
		}
		return data
	case "context":
		// 从上下文中获取值(使用缓存)
		value, _ := ctx.Cache.Get(path)
		return value
	default:
		return nil
	}
}

// getValueFromJSON 从 JSON 数据中根据路径获取值
func (c *DynamicApproverConfig) getValueFromJSON(data json.RawMessage, path string) interface{} {
	if len(data) == 0 {
		return nil
	}

	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil
	}

	parts := splitPath(path)
	current := interface{}(obj)

	for _, part := range parts {
		m, ok := current.(map[string]interface{})
		if !ok {
			return nil
		}

		value, exists := m[part]
		if !exists {
			return nil
		}

		current = value
	}

	return current
}

// addQueryParams 添加查询参数到 URL
func (c *DynamicApproverConfig) addQueryParams(baseURL string, params map[string]interface{}) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return baseURL, err
	}

	q := u.Query()
	for key, value := range params {
		// 将值转换为字符串
		var strValue string
		switch v := value.(type) {
		case string:
			strValue = v
		case int, int64, float64:
			strValue = fmt.Sprintf("%v", v)
		case bool:
			if v {
				strValue = "true"
			} else {
				strValue = "false"
			}
		default:
			// 对于复杂类型,序列化为 JSON
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				continue
			}
			strValue = string(jsonBytes)
		}
		q.Set(key, strValue)
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}

// doWithRetry 带重试机制的 HTTP 请求执行
// 使用指数退避策略
func (c *DynamicApproverConfig) doWithRetry(req *http.Request) (*http.Response, error) {
	maxRetries := 3
	baseDelay := 100 * time.Millisecond

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		// 执行请求
		resp, err := c.HTTPClient.Do(req)
		if err == nil {
			// 请求成功,检查状态码
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				return resp, nil
			}
			// 状态码错误,关闭响应体
			resp.Body.Close()
			lastErr = fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
		} else {
			lastErr = err
		}

		// 如果不是最后一次尝试,等待后重试
		if attempt < maxRetries-1 {
			// 指数退避: delay = baseDelay * 2^attempt
			delay := baseDelay * time.Duration(1<<uint(attempt))
			time.Sleep(delay)
		}
	}

	return nil, lastErr
}

// parseResponse 解析响应数据
func (c *DynamicApproverConfig) parseResponse(resp *http.Response) ([]string, error) {
	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// 解析 JSON
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// 根据 ResponseMapping 解析审批人列表
	if c.API.ResponseMapping == nil {
		return nil, fmt.Errorf("HTTPAPIConfig.ResponseMapping is required")
	}

	// 简单的路径解析(支持简单的字段访问,如 "approvers" 或 "data.approvers")
	path := c.API.ResponseMapping.Path
	value, err := getValueByPath(data, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get value by path %q: %w", path, err)
	}

	// 转换为字符串数组
	approvers, err := convertToStringSlice(value)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to string slice: %w", err)
	}

	return approvers, nil
}

// getValueByPath 根据路径获取值
// 支持简单的字段访问,如 "approvers" 或 "data.approvers"
func getValueByPath(data map[string]interface{}, path string) (interface{}, error) {
	// 简单的路径解析(不支持复杂的 JSONPath)
	parts := splitPath(path)
	current := interface{}(data)

	for _, part := range parts {
		m, ok := current.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("path %q is invalid: expected object at %q", path, part)
		}

		value, exists := m[part]
		if !exists {
			return nil, fmt.Errorf("path %q is invalid: field %q not found", path, part)
		}

		current = value
	}

	return current, nil
}

// splitPath 分割路径
func splitPath(path string) []string {
	// 简单的分割,支持 "approvers" 或 "data.approvers"
	parts := []string{}
	current := ""
	for _, char := range path {
		if char == '.' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

// convertToStringSlice 转换为字符串数组
func convertToStringSlice(value interface{}) ([]string, error) {
	// 检查是否是数组
	arr, ok := value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("value is not an array: %T", value)
	}

	// 转换为字符串数组
	result := make([]string, 0, len(arr))
	for i, item := range arr {
		str, ok := item.(string)
		if !ok {
			return nil, fmt.Errorf("array item at index %d is not a string: %T", i, item)
		}
		result = append(result, str)
	}

	return result, nil
}

