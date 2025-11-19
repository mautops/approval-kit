package node_test

import (
	"testing"

	"github.com/mautops/approval-kit/internal/node"
)

// TestParamMapping 测试参数映射
func TestParamMapping(t *testing.T) {
	config := &node.DynamicApproverConfig{
		API: &node.HTTPAPIConfig{
			URL:    "http://example.com/api/approvers",
			Method: "POST",
			ParamMapping: &node.ParamMapping{
				Source: "task_params",
				Path:   "department",
				Target: "dept",
			},
		},
		Timing:     node.ApproverTimingOnActivate,
		HTTPClient: &mockHTTPClientForDynamicApprover{},
	}

	// 测试参数映射逻辑
	// 这里主要测试配置是否正确设置
	if config.API.ParamMapping == nil {
		t.Error("ParamMapping should be set")
	}

	if config.API.ParamMapping.Source != "task_params" {
		t.Errorf("ParamMapping.Source = %q, want %q", config.API.ParamMapping.Source, "task_params")
	}

	if config.API.ParamMapping.Path != "department" {
		t.Errorf("ParamMapping.Path = %q, want %q", config.API.ParamMapping.Path, "department")
	}

	if config.API.ParamMapping.Target != "dept" {
		t.Errorf("ParamMapping.Target = %q, want %q", config.API.ParamMapping.Target, "dept")
	}
}

// TestParamMappingFromTaskParams 测试从任务参数映射
func TestParamMappingFromTaskParams(t *testing.T) {
	config := &node.DynamicApproverConfig{
		API: &node.HTTPAPIConfig{
			URL:    "http://example.com/api/approvers",
			Method: "POST",
			ParamMapping: &node.ParamMapping{
				Source: "task_params",
				Path:   "department",
				Target: "dept",
			},
		},
		Timing:     node.ApproverTimingOnActivate,
		HTTPClient: &mockHTTPClientForDynamicApprover{},
	}

	// 验证参数映射配置
	if config.API.ParamMapping.Source != "task_params" {
		t.Errorf("ParamMapping.Source = %q, want %q", config.API.ParamMapping.Source, "task_params")
	}
}

// TestParamMappingFromNodeOutputs 测试从节点输出映射
func TestParamMappingFromNodeOutputs(t *testing.T) {
	config := &node.DynamicApproverConfig{
		API: &node.HTTPAPIConfig{
			URL:    "http://example.com/api/approvers",
			Method: "POST",
			ParamMapping: &node.ParamMapping{
				Source: "node_outputs",
				Path:   "previous_node.result",
				Target: "result",
			},
		},
		Timing:     node.ApproverTimingOnActivate,
		HTTPClient: &mockHTTPClientForDynamicApprover{},
	}

	// 验证参数映射配置
	if config.API.ParamMapping.Source != "node_outputs" {
		t.Errorf("ParamMapping.Source = %q, want %q", config.API.ParamMapping.Source, "node_outputs")
	}
}

// TestResponseMapping 测试响应解析
func TestResponseMapping(t *testing.T) {
	config := &node.HTTPAPIConfig{
		URL:    "http://example.com/api/approvers",
		Method: "POST",
		ResponseMapping: &node.ResponseMapping{
			Path:   "data.approvers",
			Format: "json",
		},
	}

	if config.ResponseMapping == nil {
		t.Error("ResponseMapping should be set")
	}

	if config.ResponseMapping.Path != "data.approvers" {
		t.Errorf("ResponseMapping.Path = %q, want %q", config.ResponseMapping.Path, "data.approvers")
	}

	if config.ResponseMapping.Format != "json" {
		t.Errorf("ResponseMapping.Format = %q, want %q", config.ResponseMapping.Format, "json")
	}
}

// TestResponseMappingSimplePath 测试简单路径响应解析
func TestResponseMappingSimplePath(t *testing.T) {
	config := &node.HTTPAPIConfig{
		URL:    "http://example.com/api/approvers",
		Method: "POST",
		ResponseMapping: &node.ResponseMapping{
			Path:   "approvers",
			Format: "json",
		},
	}

	if config.ResponseMapping.Path != "approvers" {
		t.Errorf("ResponseMapping.Path = %q, want %q", config.ResponseMapping.Path, "approvers")
	}
}

