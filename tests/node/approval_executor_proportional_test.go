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

// TestApprovalNodeExecutorProportionalMode 测试比例会签模式
func TestApprovalNodeExecutorProportionalMode(t *testing.T) {
	config := &node.ApprovalNodeConfig{
		Mode: node.ApprovalModeProportional,
		ApproverConfig: &node.FixedApproverConfig{
			Approvers: []string{"user-001", "user-002", "user-003", "user-004", "user-005"},
		},
		ProportionalThreshold: &node.ProportionalThreshold{
			Required: 3, // 5 人中需要 3 人同意
			Total:    5,
		},
	}

	executor := node.NewApprovalNodeExecutor(config)

	// 验证节点类型
	if executor.NodeType() != template.NodeTypeApproval {
		t.Errorf("ApprovalNodeExecutor.NodeType() = %v, want %v", executor.NodeType(), template.NodeTypeApproval)
	}

	// 测试: 同意人数未达到阈值(未完成)
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
				// 只有 2 人同意,未达到 3 人阈值
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

	// 执行节点,应该返回错误(审批未完成)
	result, err := executor.Execute(ctx)
	if err == nil {
		t.Error("ApprovalNodeExecutor.Execute() should return error when approval threshold not met")
	}
	if result != nil {
		t.Error("ApprovalNodeExecutor.Execute() should return nil result when approval threshold not met")
	}
}

// TestApprovalNodeExecutorProportionalModeThresholdMet 测试比例会签模式达到阈值
func TestApprovalNodeExecutorProportionalModeThresholdMet(t *testing.T) {
	config := &node.ApprovalNodeConfig{
		Mode: node.ApprovalModeProportional,
		ApproverConfig: &node.FixedApproverConfig{
			Approvers: []string{"user-001", "user-002", "user-003", "user-004", "user-005"},
		},
		ProportionalThreshold: &node.ProportionalThreshold{
			Required: 3, // 5 人中需要 3 人同意
			Total:    5,
		},
	}

	executor := node.NewApprovalNodeExecutor(config)

	// 创建任务,3 人同意(达到阈值)
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
				"user-003": {
					Result:    "approve",
					Comment:   "approved",
					CreatedAt: time.Now(),
				},
				// user-004 和 user-005 还未审批,但已达到阈值
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

	// 执行节点,应该成功(达到阈值)
	result, err := executor.Execute(ctx)
	if err != nil {
		t.Fatalf("ApprovalNodeExecutor.Execute() failed: %v", err)
	}

	// 验证结果
	if result == nil {
		t.Fatal("ApprovalNodeExecutor.Execute() should return a result when threshold met")
	}

	// 验证生成了事件
	if len(result.Events) == 0 {
		t.Error("ApprovalNodeExecutor should generate events")
	}
}

// TestApprovalNodeExecutorProportionalModeWithRejections 测试比例会签模式包含拒绝
func TestApprovalNodeExecutorProportionalModeWithRejections(t *testing.T) {
	config := &node.ApprovalNodeConfig{
		Mode: node.ApprovalModeProportional,
		ApproverConfig: &node.FixedApproverConfig{
			Approvers: []string{"user-001", "user-002", "user-003", "user-004", "user-005"},
		},
		ProportionalThreshold: &node.ProportionalThreshold{
			Required: 3, // 5 人中需要 3 人同意
			Total:    5,
		},
		RejectBehavior: node.RejectBehaviorTerminate,
	}

	executor := node.NewApprovalNodeExecutor(config)

	// 创建任务,3 人同意 2 人拒绝(达到阈值)
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
				"user-003": {
					Result:    "approve",
					Comment:   "approved",
					CreatedAt: time.Now(),
				},
				"user-004": {
					Result:    "reject",
					Comment:   "rejected",
					CreatedAt: time.Now(),
				},
				"user-005": {
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

	// 执行节点,应该成功(达到阈值,即使有拒绝)
	result, err := executor.Execute(ctx)
	if err != nil {
		t.Fatalf("ApprovalNodeExecutor.Execute() failed: %v", err)
	}

	// 验证结果
	if result == nil {
		t.Fatal("ApprovalNodeExecutor.Execute() should return a result when threshold met")
	}
}

