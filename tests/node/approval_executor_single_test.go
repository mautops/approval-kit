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

// TestApprovalNodeExecutorSingleMode 测试单人审批模式执行器
func TestApprovalNodeExecutorSingleMode(t *testing.T) {
	config := &node.ApprovalNodeConfig{
		Mode: node.ApprovalModeSingle,
		ApproverConfig: &node.FixedApproverConfig{
			Approvers: []string{"user-001"},
		},
	}

	executor := node.NewApprovalNodeExecutor(config)

	// 验证节点类型
	if executor.NodeType() != template.NodeTypeApproval {
		t.Errorf("ApprovalNodeExecutor.NodeType() = %v, want %v", executor.NodeType(), template.NodeTypeApproval)
	}

	// 创建测试上下文
	tsk := &task.Task{
		ID:    "task-001",
		State: types.TaskStateApproving,
		Approvals: make(map[string]map[string]*task.Approval),
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

	// 测试: 审批未完成(没有审批记录)
	result, err := executor.Execute(ctx)
	if err == nil {
		t.Error("ApprovalNodeExecutor.Execute() should return error when approval is pending")
	}
	if result != nil {
		t.Error("ApprovalNodeExecutor.Execute() should return nil result when approval is pending")
	}
}

// TestApprovalNodeExecutorSingleModeApproved 测试单人审批模式已同意
func TestApprovalNodeExecutorSingleModeApproved(t *testing.T) {
	config := &node.ApprovalNodeConfig{
		Mode: node.ApprovalModeSingle,
		ApproverConfig: &node.FixedApproverConfig{
			Approvers: []string{"user-001"},
		},
	}

	executor := node.NewApprovalNodeExecutor(config)

	// 创建任务,包含审批记录
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
		t.Fatal("ApprovalNodeExecutor.Execute() should return a result when approved")
	}

	// 验证生成了事件
	if len(result.Events) == 0 {
		t.Error("ApprovalNodeExecutor should generate events")
	}
}

// TestApprovalNodeExecutorSingleModeRejected 测试单人审批模式已拒绝
func TestApprovalNodeExecutorSingleModeRejected(t *testing.T) {
	config := &node.ApprovalNodeConfig{
		Mode: node.ApprovalModeSingle,
		ApproverConfig: &node.FixedApproverConfig{
			Approvers: []string{"user-001"},
		},
		RejectBehavior: node.RejectBehaviorTerminate,
	}

	executor := node.NewApprovalNodeExecutor(config)

	// 创建任务,包含拒绝记录
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
		t.Fatal("ApprovalNodeExecutor.Execute() should return a result when rejected")
	}
}

