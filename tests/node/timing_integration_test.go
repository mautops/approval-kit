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

// TestApproverTimingOnCreateIntegration 测试任务创建时获取审批人的集成场景
func TestApproverTimingOnCreateIntegration(t *testing.T) {
	// 创建 mock HTTPClient
	mockClient := &mockHTTPClientForDynamicApprover{
		response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`{"approvers": ["user-001", "user-002"]}`)),
		},
		err: nil,
	}

	// 创建模板,包含配置为 on_create 时机的动态审批人节点
	tpl := &template.Template{
		ID:          "tpl-001",
		Name:        "Test Template",
		Description: "Test template with dynamic approver on create",
		Version:     1,
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "Start Node",
				Type: template.NodeTypeStart,
			},
			"approval": {
				ID:   "approval-001",
				Name: "Approval Node",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeSingle,
					ApproverConfig: &node.DynamicApproverConfig{
						API: &node.HTTPAPIConfig{
							URL:    "http://example.com/api/approvers",
							Method: "POST",
							ResponseMapping: &node.ResponseMapping{
								Path:   "approvers",
								Format: "json",
							},
						},
						Timing:     node.ApproverTimingOnCreate,
						HTTPClient: mockClient,
					},
				},
			},
		},
		Edges: []*template.Edge{},
	}

	// 创建任务
	tsk := &task.Task{
		ID:    "task-001",
		State: types.TaskStatePending,
		Params: json.RawMessage(`{"department": "engineering"}`),
		Approvers: make(map[string][]string),
	}

	// 调用 FetchApproversOnCreate
	err := node.FetchApproversOnCreate(tpl, tsk, mockClient)
	if err != nil {
		t.Fatalf("FetchApproversOnCreate() failed: %v", err)
	}

	// 验证审批人已获取
	approvers, exists := tsk.Approvers["approval-001"]
	if !exists {
		t.Error("Approvers should be fetched for approval-001 node")
	}

	if len(approvers) != 2 {
		t.Errorf("Expected 2 approvers, got %d", len(approvers))
	}

	expectedApprovers := []string{"user-001", "user-002"}
	for i, approver := range approvers {
		if approver != expectedApprovers[i] {
			t.Errorf("Approvers[%d] = %q, want %q", i, approver, expectedApprovers[i])
		}
	}
}

// TestApproverTimingOnActivateIntegration 测试节点激活时获取审批人的集成场景
func TestApproverTimingOnActivateIntegration(t *testing.T) {
	// 创建 mock HTTPClient
	mockClient := &mockHTTPClientForDynamicApprover{
		response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`{"approvers": ["user-003"]}`)),
		},
		err: nil,
	}

	// 创建配置为 on_activate 时机的动态审批人配置
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

	// 调用 GetApprovers (节点激活时)
	approvers, err := config.GetApprovers(ctx)
	if err != nil {
		t.Fatalf("GetApprovers() failed: %v", err)
	}

	if len(approvers) != 1 {
		t.Errorf("Expected 1 approver, got %d", len(approvers))
	}

	if approvers[0] != "user-003" {
		t.Errorf("Approvers[0] = %q, want %q", approvers[0], "user-003")
	}
}

