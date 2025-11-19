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

// TestApprovalNodeExecutorSequentialMode 测试顺序审批模式
func TestApprovalNodeExecutorSequentialMode(t *testing.T) {
	config := &node.ApprovalNodeConfig{
		Mode: node.ApprovalModeSequential,
		ApproverConfig: &node.FixedApproverConfig{
			Approvers: []string{"user-001", "user-002", "user-003"},
		},
	}

	executor := node.NewApprovalNodeExecutor(config)

	// 验证节点类型
	if executor.NodeType() != template.NodeTypeApproval {
		t.Errorf("ApprovalNodeExecutor.NodeType() = %v, want %v", executor.NodeType(), template.NodeTypeApproval)
	}

	// 测试: 第一个审批人还未审批(未完成)
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

// TestApprovalNodeExecutorSequentialModeFirstApproved 测试顺序审批模式第一个审批人同意
func TestApprovalNodeExecutorSequentialModeFirstApproved(t *testing.T) {
	config := &node.ApprovalNodeConfig{
		Mode: node.ApprovalModeSequential,
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
				// user-002 还未审批,但顺序审批模式下需要等待下一个审批人
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

	// 执行节点,应该返回错误(还有下一个审批人需要审批)
	result, err := executor.Execute(ctx)
	if err == nil {
		t.Error("ApprovalNodeExecutor.Execute() should return error when next approver needs to approve")
	}
	if result != nil {
		t.Error("ApprovalNodeExecutor.Execute() should return nil result when next approver needs to approve")
	}
}

// TestApprovalNodeExecutorSequentialModeAllApproved 测试顺序审批模式全部同意
func TestApprovalNodeExecutorSequentialModeAllApproved(t *testing.T) {
	config := &node.ApprovalNodeConfig{
		Mode: node.ApprovalModeSequential,
		ApproverConfig: &node.FixedApproverConfig{
			Approvers: []string{"user-001", "user-002"},
		},
	}

	executor := node.NewApprovalNodeExecutor(config)

	// 创建任务,所有审批人都按顺序同意
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

	// 执行节点,应该成功(所有审批人都已按顺序同意)
	result, err := executor.Execute(ctx)
	if err != nil {
		t.Fatalf("ApprovalNodeExecutor.Execute() failed: %v", err)
	}

	// 验证结果
	if result == nil {
		t.Fatal("ApprovalNodeExecutor.Execute() should return a result when all approved")
	}

	// 验证生成了事件
	if len(result.Events) == 0 {
		t.Error("ApprovalNodeExecutor should generate events")
	}
}

// TestApprovalNodeExecutorSequentialModeRejected 测试顺序审批模式被拒绝
func TestApprovalNodeExecutorSequentialModeRejected(t *testing.T) {
	config := &node.ApprovalNodeConfig{
		Mode: node.ApprovalModeSequential,
		ApproverConfig: &node.FixedApproverConfig{
			Approvers: []string{"user-001", "user-002"},
		},
		RejectBehavior: node.RejectBehaviorTerminate,
	}

	executor := node.NewApprovalNodeExecutor(config)

	// 创建任务,第一个审批人拒绝
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
				// user-002 还未审批,但第一个审批人拒绝后流程应该终止
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

	// 验证结果是拒绝
	if result.Output == nil {
		t.Error("ApprovalNodeExecutor should return output with result")
	}
}

// TestApprovalNodeExecutorSequentialModeOutOfOrder 测试顺序审批模式乱序审批
func TestApprovalNodeExecutorSequentialModeOutOfOrder(t *testing.T) {
	config := &node.ApprovalNodeConfig{
		Mode: node.ApprovalModeSequential,
		ApproverConfig: &node.FixedApproverConfig{
			Approvers: []string{"user-001", "user-002", "user-003"},
		},
	}

	executor := node.NewApprovalNodeExecutor(config)

	// 创建任务,第二个审批人先审批(乱序)
	tsk := &task.Task{
		ID:    "task-001",
		State: types.TaskStateApproving,
		Approvals: map[string]map[string]*task.Approval{
			"approval-001": {
				"user-002": {
					Result:    "approve",
					Comment:   "approved",
					CreatedAt: time.Now(),
				},
				// user-001 还未审批,但 user-002 已经审批了(乱序)
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

	// 执行节点,应该返回错误(第一个审批人还未审批)
	result, err := executor.Execute(ctx)
	if err == nil {
		t.Error("ApprovalNodeExecutor.Execute() should return error when previous approver has not approved")
	}
	if result != nil {
		t.Error("ApprovalNodeExecutor.Execute() should return nil result when previous approver has not approved")
	}
}

