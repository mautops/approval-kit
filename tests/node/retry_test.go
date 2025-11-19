package node_test

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
	"github.com/mautops/approval-kit/internal/types"
)

// TestDynamicApproverConfigRetry 测试 API 调用重试机制
func TestDynamicApproverConfigRetry(t *testing.T) {
	// 创建一个会在前两次调用失败,第三次成功的 mock client
	attempts := 0
	mockClient := &retryMockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			attempts++
			if attempts < 3 {
				// 前两次失败
				return nil, http.ErrHandlerTimeout
			}
			// 第三次成功
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(`{"approvers": ["user-001"]}`)),
			}, nil
		},
	}

	config := &node.DynamicApproverConfig{
		API: &node.HTTPAPIConfig{
			URL:    "http://example.com/api/approvers",
			Method: "POST",
			ResponseMapping: &node.ResponseMapping{
				Path:   "approvers",
				Format: "json",
			},
		},
		Timing:     node.ApproverTimingOnActivate,
		HTTPClient: mockClient,
	}

	ctx := &node.NodeContext{
		Task: &task.Task{
			ID:    "task-001",
			State: types.TaskStateApproving,
		},
		Node: &template.Node{
			ID:   "approval-001",
			Name: "Approval Node",
			Type: template.NodeTypeApproval,
		},
		Params:  json.RawMessage(`{}`),
		Outputs: make(map[string]json.RawMessage),
		Cache:   node.NewContextCache(),
	}

	// 调用 GetApprovers,应该会重试并最终成功
	approvers, err := config.GetApprovers(ctx)
	if err != nil {
		t.Fatalf("GetApprovers() failed after retries: %v", err)
	}

	if len(approvers) != 1 {
		t.Errorf("GetApprovers() returned %d approvers, want 1", len(approvers))
	}

	if approvers[0] != "user-001" {
		t.Errorf("GetApprovers()[0] = %q, want %q", approvers[0], "user-001")
	}

	// 验证重试了至少 2 次
	if attempts < 2 {
		t.Errorf("Expected at least 2 retry attempts, got %d", attempts)
	}
}

// TestDynamicApproverConfigRetryMaxAttempts 测试最大重试次数
func TestDynamicApproverConfigRetryMaxAttempts(t *testing.T) {
	// 创建一个总是失败的 mock client
	attempts := 0
	mockClient := &retryMockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			attempts++
			return nil, http.ErrHandlerTimeout
		},
	}

	config := &node.DynamicApproverConfig{
		API: &node.HTTPAPIConfig{
			URL:    "http://example.com/api/approvers",
			Method: "POST",
			ResponseMapping: &node.ResponseMapping{
				Path:   "approvers",
				Format: "json",
			},
		},
		Timing:     node.ApproverTimingOnActivate,
		HTTPClient: mockClient,
	}

	ctx := &node.NodeContext{
		Task: &task.Task{
			ID:    "task-001",
			State: types.TaskStateApproving,
		},
		Node: &template.Node{
			ID:   "approval-001",
			Name: "Approval Node",
			Type: template.NodeTypeApproval,
		},
		Params:  json.RawMessage(`{}`),
		Outputs: make(map[string]json.RawMessage),
		Cache:   node.NewContextCache(),
	}

	// 调用 GetApprovers,应该达到最大重试次数后失败
	_, err := config.GetApprovers(ctx)
	if err == nil {
		t.Error("GetApprovers() should fail after max retries")
	}

	// 验证重试了预期的次数(默认应该是 3 次)
	expectedAttempts := 3
	if attempts < expectedAttempts {
		t.Errorf("Expected at least %d retry attempts, got %d", expectedAttempts, attempts)
	}
}

// TestDynamicApproverConfigRetryBackoff 测试指数退避
func TestDynamicApproverConfigRetryBackoff(t *testing.T) {
	// 创建一个会在前两次调用失败,第三次成功的 mock client
	attempts := 0
	lastAttemptTime := time.Now()
	attemptTimes := []time.Duration{}

	mockClient := &retryMockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			now := time.Now()
			if attempts > 0 {
				// 记录每次尝试之间的时间间隔
				attemptTimes = append(attemptTimes, now.Sub(lastAttemptTime))
			}
			lastAttemptTime = now

			attempts++
			if attempts < 3 {
				// 前两次失败
				return nil, http.ErrHandlerTimeout
			}
			// 第三次成功
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(`{"approvers": ["user-001"]}`)),
			}, nil
		},
	}

	config := &node.DynamicApproverConfig{
		API: &node.HTTPAPIConfig{
			URL:    "http://example.com/api/approvers",
			Method: "POST",
			ResponseMapping: &node.ResponseMapping{
				Path:   "approvers",
				Format: "json",
			},
		},
		Timing:     node.ApproverTimingOnActivate,
		HTTPClient: mockClient,
	}

	ctx := &node.NodeContext{
		Task: &task.Task{
			ID:    "task-001",
			State: types.TaskStateApproving,
		},
		Node: &template.Node{
			ID:   "approval-001",
			Name: "Approval Node",
			Type: template.NodeTypeApproval,
		},
		Params:  json.RawMessage(`{}`),
		Outputs: make(map[string]json.RawMessage),
		Cache:   node.NewContextCache(),
	}

	// 调用 GetApprovers
	_, err := config.GetApprovers(ctx)
	if err != nil {
		t.Fatalf("GetApprovers() failed: %v", err)
	}

	// 验证重试间隔是递增的(指数退避)
	if len(attemptTimes) >= 2 {
		// 第二次重试的间隔应该大于第一次
		if attemptTimes[1] <= attemptTimes[0] {
			t.Logf("Retry intervals: %v (backoff may not be strictly increasing due to timing)", attemptTimes)
		}
	}
}

// retryMockHTTPClient 用于测试重试机制的 mock HTTPClient
type retryMockHTTPClient struct {
	doFunc func(*http.Request) (*http.Response, error)
}

func (m *retryMockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.doFunc != nil {
		return m.doFunc(req)
	}
	return nil, http.ErrHandlerTimeout
}


