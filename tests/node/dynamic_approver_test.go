package node_test

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
	"github.com/mautops/approval-kit/internal/types"
)

// TestDynamicApproverConfig 测试 DynamicApproverConfig 结构体
func TestDynamicApproverConfig(t *testing.T) {
	mockClient := &mockHTTPClientForDynamicApprover{
		response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`{"approvers": ["user-001", "user-002"]}`)),
		},
		err: nil,
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

	if config.API == nil {
		t.Error("DynamicApproverConfig.API should not be nil")
	}

	if config.GetTiming() != node.ApproverTimingOnActivate {
		t.Errorf("DynamicApproverConfig.GetTiming() = %v, want %v", config.GetTiming(), node.ApproverTimingOnActivate)
	}

	if config.HTTPClient == nil {
		t.Error("DynamicApproverConfig.HTTPClient should not be nil")
	}
}

// TestDynamicApproverConfigGetApprovers 测试动态审批人获取
func TestDynamicApproverConfigGetApprovers(t *testing.T) {
	mockClient := &mockHTTPClientForDynamicApprover{
		response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`{"approvers": ["user-001", "user-002", "user-003"]}`)),
		},
		err: nil,
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
			Params: json.RawMessage(`{"department": "engineering"}`),
		},
		Node: &template.Node{
			ID:   "approval-001",
			Name: "Approval Node",
			Type: template.NodeTypeApproval,
		},
		Params:  json.RawMessage(`{"department": "engineering"}`),
		Outputs: make(map[string]json.RawMessage),
		Cache:   node.NewContextCache(),
	}

	approvers, err := config.GetApprovers(ctx)
	if err != nil {
		t.Fatalf("DynamicApproverConfig.GetApprovers() failed: %v", err)
	}

	if len(approvers) != 3 {
		t.Errorf("DynamicApproverConfig.GetApprovers() returned %d approvers, want 3", len(approvers))
	}

	expectedApprovers := []string{"user-001", "user-002", "user-003"}
	for i, approver := range approvers {
		if approver != expectedApprovers[i] {
			t.Errorf("DynamicApproverConfig.GetApprovers() approvers[%d] = %q, want %q", i, approver, expectedApprovers[i])
		}
	}
}

// TestDynamicApproverConfigGetApproversError 测试 API 调用失败
func TestDynamicApproverConfigGetApproversError(t *testing.T) {
	mockClient := &mockHTTPClientForDynamicApprover{
		response: nil,
		err:      http.ErrHandlerTimeout,
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

	_, err := config.GetApprovers(ctx)
	if err == nil {
		t.Error("DynamicApproverConfig.GetApprovers() should return error when API call fails")
	}
}

// TestDynamicApproverConfigGetApproversInvalidResponse 测试无效响应
func TestDynamicApproverConfigGetApproversInvalidResponse(t *testing.T) {
	mockClient := &mockHTTPClientForDynamicApprover{
		response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`{"invalid": "response"}`)),
		},
		err: nil,
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

	_, err := config.GetApprovers(ctx)
	if err == nil {
		t.Error("DynamicApproverConfig.GetApprovers() should return error when response is invalid")
	}
}

// TestDynamicApproverConfigTiming 测试获取时机
func TestDynamicApproverConfigTiming(t *testing.T) {
	config := &node.DynamicApproverConfig{
		API: &node.HTTPAPIConfig{
			URL:    "http://example.com/api/approvers",
			Method: "POST",
		},
		Timing: node.ApproverTimingOnCreate,
	}

	if config.GetTiming() != node.ApproverTimingOnCreate {
		t.Errorf("DynamicApproverConfig.GetTiming() = %v, want %v", config.GetTiming(), node.ApproverTimingOnCreate)
	}
}

// mockHTTPClientForDynamicApprover 用于测试的 HTTPClient 实现
type mockHTTPClientForDynamicApprover struct {
	response *http.Response
	err      error
}

func (m *mockHTTPClientForDynamicApprover) Do(req *http.Request) (*http.Response, error) {
	return m.response, m.err
}

