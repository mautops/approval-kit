package node_test

import (
	"encoding/json"
	"testing"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
	"github.com/mautops/approval-kit/internal/types"
)

// TestEndNodeExecutor 测试结束节点执行器
func TestEndNodeExecutor(t *testing.T) {
	executor := node.NewEndNodeExecutor()

	// 验证节点类型
	if executor.NodeType() != template.NodeTypeEnd {
		t.Errorf("EndNodeExecutor.NodeType() = %v, want %v", executor.NodeType(), template.NodeTypeEnd)
	}

	// 创建测试上下文
	tsk := &task.Task{
		ID:    "task-001",
		State: types.TaskStateApproving,
	}

	tplNode := &template.Node{
		ID:   "end",
		Name: "End Node",
		Type: template.NodeTypeEnd,
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
		t.Fatalf("EndNodeExecutor.Execute() failed: %v", err)
	}

	// 验证结果
	if result == nil {
		t.Fatal("EndNodeExecutor.Execute() should return a result")
	}

	// 结束节点的下一个节点应该为空
	if result.NextNodeID != "" {
		t.Errorf("EndNodeExecutor.NextNodeID = %q, want empty string", result.NextNodeID)
	}

	// 结束节点应该生成节点完成事件
	if len(result.Events) == 0 {
		t.Error("EndNodeExecutor should generate node_completed event")
	}

	// 验证事件类型
	foundCompleted := false
	for _, event := range result.Events {
		if event.Type == node.EventTypeNodeCompleted {
			foundCompleted = true
			break
		}
	}
	if !foundCompleted {
		t.Error("EndNodeExecutor should generate EventTypeNodeCompleted event")
	}
}

// TestEndNodeExecutorOutput 测试结束节点的输出
func TestEndNodeExecutorOutput(t *testing.T) {
	executor := node.NewEndNodeExecutor()

	ctx := &node.NodeContext{
		Task:    &task.Task{ID: "task-001"},
		Node:    &template.Node{ID: "end", Type: template.NodeTypeEnd},
		Params:  json.RawMessage(`{}`),
		Outputs: make(map[string]json.RawMessage),
		Cache:   node.NewContextCache(),
	}

	result, err := executor.Execute(ctx)
	if err != nil {
		t.Fatalf("EndNodeExecutor.Execute() failed: %v", err)
	}

	// 结束节点可以输出最终结果或空输出
	_ = result.Output
}

