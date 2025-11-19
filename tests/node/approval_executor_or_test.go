package node_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
	"github.com/mautops/approval-kit/internal/types"
)

// TestApprovalNodeExecutorOrMode 测试多人或签模式
func TestApprovalNodeExecutorOrMode(t *testing.T) {
	config := &node.ApprovalNodeConfig{
		Mode: node.ApprovalModeOr,
		ApproverConfig: &node.FixedApproverConfig{
			Approvers: []string{"user-001", "user-002", "user-003"},
		},
	}

	executor := node.NewApprovalNodeExecutor(config)

	// 验证节点类型
	if executor.NodeType() != template.NodeTypeApproval {
		t.Errorf("ApprovalNodeExecutor.NodeType() = %v, want %v", executor.NodeType(), template.NodeTypeApproval)
	}

	// 测试: 所有审批人都未审批(未完成)
	tsk := &task.Task{
		ID:    "task-001",
		State: types.TaskStateApproving,
		Approvals: map[string]map[string]*task.Approval{
			"approval-001": {},
		},
	}

	tplNode := &template.Node{
		ID:   "approval-001",
		Name: "Approval Node",
		Type: template.NodeTypeApproval,
	}

	ctx := &node.NodeContext{
		Task:    tsk,
		Node:    tplNode,
		Params:  json.RawMessage(`{}`),
		Outputs: make(map[string]json.RawMessage),
		Cache:   node.NewContextCache(),
	}

	// 执行节点,应该返回错误(审批未完成)
	result, err := executor.Execute(ctx)
	if err == nil {
		t.Error("ApprovalNodeExecutor.Execute() should return error when approval is pending")
	}
	if result != nil {
		t.Error("ApprovalNodeExecutor.Execute() should return nil result when approval is pending")
	}
}

// TestApprovalNodeExecutorOrModeFirstApproved 测试多人或签模式第一个审批人同意
func TestApprovalNodeExecutorOrModeFirstApproved(t *testing.T) {
	config := &node.ApprovalNodeConfig{
		Mode: node.ApprovalModeOr,
		ApproverConfig: &node.FixedApproverConfig{
			Approvers: []string{"user-001", "user-002"},
		},
	}

	executor := node.NewApprovalNodeExecutor(config)

	// 创建任务,第一个审批人同意
	tsk := &task.Task{
		ID:    "task-001",
		State: types.TaskStateApproving,
		Approvals: map[string]map[string]*task.Approval{
			"approval-001": {
				"user-001": {
					Result:    "approve",
					Comment:   "approved",
					CreatedAt: time.Now(),
				},
				// user-002 还未审批
			},
		},
	}

	tplNode := &template.Node{
		ID:   "approval-001",
		Name: "Approval Node",
		Type: template.NodeTypeApproval,
	}

	ctx := &node.NodeContext{
		Task:    tsk,
		Node:    tplNode,
		Params:  json.RawMessage(`{}`),
		Outputs: make(map[string]json.RawMessage),
		Cache:   node.NewContextCache(),
	}

	// 执行节点,应该成功(第一个审批人同意即可)
	result, err := executor.Execute(ctx)
	if err != nil {
		t.Fatalf("ApprovalNodeExecutor.Execute() failed: %v", err)
	}

	// 验证结果
	if result == nil {
		t.Fatal("ApprovalNodeExecutor.Execute() should return a result when first approved")
	}

	// 验证生成了事件
	if len(result.Events) == 0 {
		t.Error("ApprovalNodeExecutor should generate events")
	}
}

// TestApprovalNodeExecutorOrModeAllRejected 测试多人或签模式全部拒绝
func TestApprovalNodeExecutorOrModeAllRejected(t *testing.T) {
	config := &node.ApprovalNodeConfig{
		Mode: node.ApprovalModeOr,
		ApproverConfig: &node.FixedApproverConfig{
			Approvers: []string{"user-001", "user-002"},
		},
		RejectBehavior: node.RejectBehaviorTerminate,
	}

	executor := node.NewApprovalNodeExecutor(config)

	// 创建任务,所有审批人都拒绝
	tsk := &task.Task{
		ID:    "task-001",
		State: types.TaskStateApproving,
		Approvals: map[string]map[string]*task.Approval{
			"approval-001": {
				"user-001": {
					Result:    "reject",
					Comment:   "rejected",
					CreatedAt: time.Now(),
				},
				"user-002": {
					Result:    "reject",
					Comment:   "rejected",
					CreatedAt: time.Now(),
				},
			},
		},
	}

	tplNode := &template.Node{
		ID:   "approval-001",
		Name: "Approval Node",
		Type: template.NodeTypeApproval,
	}

	ctx := &node.NodeContext{
		Task:    tsk,
		Node:    tplNode,
		Params:  json.RawMessage(`{}`),
		Outputs: make(map[string]json.RawMessage),
		Cache:   node.NewContextCache(),
	}

	// 执行节点
	result, err := executor.Execute(ctx)
	if err != nil {
		t.Fatalf("ApprovalNodeExecutor.Execute() failed: %v", err)
	}

	// 验证结果
	if result == nil {
		t.Fatal("ApprovalNodeExecutor.Execute() should return a result when all rejected")
	}

	// 验证结果是拒绝
	if result.Output == nil {
		t.Error("ApprovalNodeExecutor should return output with result")
	}
}

// TestApprovalNodeExecutorOrModeMixed 测试多人或签模式混合情况(一人同意一人拒绝)
func TestApprovalNodeExecutorOrModeMixed(t *testing.T) {
	config := &node.ApprovalNodeConfig{
		Mode: node.ApprovalModeOr,
		ApproverConfig: &node.FixedApproverConfig{
			Approvers: []string{"user-001", "user-002"},
		},
	}

	executor := node.NewApprovalNodeExecutor(config)

	// 创建任务,一人同意一人拒绝
	tsk := &task.Task{
		ID:    "task-001",
		State: types.TaskStateApproving,
		Approvals: map[string]map[string]*task.Approval{
			"approval-001": {
				"user-001": {
					Result:    "approve",
					Comment:   "approved",
					CreatedAt: time.Now(),
				},
				"user-002": {
					Result:    "reject",
					Comment:   "rejected",
					CreatedAt: time.Now(),
				},
			},
		},
	}

	tplNode := &template.Node{
		ID:   "approval-001",
		Name: "Approval Node",
		Type: template.NodeTypeApproval,
	}

	ctx := &node.NodeContext{
		Task:    tsk,
		Node:    tplNode,
		Params:  json.RawMessage(`{}`),
		Outputs: make(map[string]json.RawMessage),
		Cache:   node.NewContextCache(),
	}

	// 执行节点,应该成功(第一个审批人同意即可)
	result, err := executor.Execute(ctx)
	if err != nil {
		t.Fatalf("ApprovalNodeExecutor.Execute() failed: %v", err)
	}

	// 验证结果
	if result == nil {
		t.Fatal("ApprovalNodeExecutor.Execute() should return a result when first approved")
	}
}

