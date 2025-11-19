package node_test

import (
	"encoding/json"
	"testing"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
)

// TestFixedApproverConfig 测试固定审批人配置
func TestFixedApproverConfig(t *testing.T) {
	config := &node.FixedApproverConfig{
		Approvers: []string{"user-001", "user-002", "user-003"},
	}

	// 创建测试上下文
	ctx := &node.NodeContext{
		Task:    &task.Task{ID: "task-001"},
		Node:    &template.Node{ID: "node-001"},
		Params:  json.RawMessage(`{}`),
		Outputs: make(map[string]json.RawMessage),
		Cache:   node.NewContextCache(),
	}

	// 获取审批人
	approvers, err := config.GetApprovers(ctx)
	if err != nil {
		t.Fatalf("FixedApproverConfig.GetApprovers() failed: %v", err)
	}

	// 验证审批人列表
	if len(approvers) != 3 {
		t.Errorf("GetApprovers() returned %d approvers, want 3", len(approvers))
	}
	if approvers[0] != "user-001" {
		t.Errorf("GetApprovers()[0] = %q, want %q", approvers[0], "user-001")
	}
}

// TestFixedApproverConfigIsolation 测试返回的审批人列表是隔离的
func TestFixedApproverConfigIsolation(t *testing.T) {
	config := &node.FixedApproverConfig{
		Approvers: []string{"user-001", "user-002"},
	}

	ctx := &node.NodeContext{
		Task:    &task.Task{ID: "task-001"},
		Node:    &template.Node{ID: "node-001"},
		Params:  json.RawMessage(`{}`),
		Outputs: make(map[string]json.RawMessage),
		Cache:   node.NewContextCache(),
	}

	// 获取审批人
	approvers, err := config.GetApprovers(ctx)
	if err != nil {
		t.Fatalf("FixedApproverConfig.GetApprovers() failed: %v", err)
	}

	// 修改返回的列表
	approvers[0] = "modified-user"

	// 再次获取,验证原配置未被修改
	approvers2, err := config.GetApprovers(ctx)
	if err != nil {
		t.Fatalf("FixedApproverConfig.GetApprovers() failed on second call: %v", err)
	}
	if approvers2[0] == "modified-user" {
		t.Error("GetApprovers() should return a copy, modifications should not affect config")
	}
	if approvers2[0] != "user-001" {
		t.Errorf("GetApprovers()[0] = %q, want %q", approvers2[0], "user-001")
	}
}

// TestFixedApproverConfigTiming 测试获取时机
func TestFixedApproverConfigTiming(t *testing.T) {
	config := &node.FixedApproverConfig{
		Approvers: []string{"user-001"},
	}

	// 验证默认时机
	timing := config.GetTiming()
	if timing != node.ApproverTimingOnActivate {
		t.Errorf("FixedApproverConfig.GetTiming() = %v, want %v", timing, node.ApproverTimingOnActivate)
	}
}

// TestFixedApproverConfigEmpty 测试空审批人列表
func TestFixedApproverConfigEmpty(t *testing.T) {
	config := &node.FixedApproverConfig{
		Approvers: []string{},
	}

	ctx := &node.NodeContext{
		Task:    &task.Task{ID: "task-001"},
		Node:    &template.Node{ID: "node-001"},
		Params:  json.RawMessage(`{}`),
		Outputs: make(map[string]json.RawMessage),
		Cache:   node.NewContextCache(),
	}

	approvers, err := config.GetApprovers(ctx)
	if err != nil {
		t.Fatalf("FixedApproverConfig.GetApprovers() failed: %v", err)
	}

	// 空列表应该返回空列表
	if len(approvers) != 0 {
		t.Errorf("GetApprovers() returned %d approvers for empty config, want 0", len(approvers))
	}
}

