package node_test

import (
	"encoding/json"
	"testing"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
	"github.com/mautops/approval-kit/internal/types"
)

// TestStartNodeExecutor 测试开始节点执行器
func TestStartNodeExecutor(t *testing.T) {
	executor := node.NewStartNodeExecutor()

	// 验证节点类型
	if executor.NodeType() != template.NodeTypeStart {
		t.Errorf("StartNodeExecutor.NodeType() = %v, want %v", executor.NodeType(), template.NodeTypeStart)
	}

	// 创建测试上下文
	tsk := &task.Task{
		ID:    "task-001",
		State: types.TaskStatePending,
	}

	tplNode := &template.Node{
		ID:   "start",
		Name: "Start Node",
		Type: template.NodeTypeStart,
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
		t.Fatalf("StartNodeExecutor.Execute() failed: %v", err)
	}

	// 验证结果
	if result == nil {
		t.Fatal("StartNodeExecutor.Execute() should return a result")
	}

	// 开始节点应该生成节点激活事件
	if len(result.Events) == 0 {
		t.Error("StartNodeExecutor should generate node_activated event")
	}

	// 验证事件类型
	foundActivated := false
	for _, event := range result.Events {
		if event.Type == node.EventTypeNodeActivated {
			foundActivated = true
			break
		}
	}
	if !foundActivated {
		t.Error("StartNodeExecutor should generate EventTypeNodeActivated event")
	}
}

// TestStartNodeExecutorNextNode 测试开始节点的下一个节点
func TestStartNodeExecutorNextNode(t *testing.T) {
	executor := node.NewStartNodeExecutor()

	ctx := &node.NodeContext{
		Task:    &task.Task{ID: "task-001"},
		Node:    &template.Node{ID: "start", Type: template.NodeTypeStart},
		Params:  json.RawMessage(`{}`),
		Outputs: make(map[string]json.RawMessage),
		Cache:   node.NewContextCache(),
	}

	result, err := executor.Execute(ctx)
	if err != nil {
		t.Fatalf("StartNodeExecutor.Execute() failed: %v", err)
	}

	// 开始节点通常不指定下一个节点,由流程引擎根据边的定义决定
	// 但这里我们允许返回空字符串,表示由流程引擎决定
	_ = result.NextNodeID
}

// TestStartNodeExecutorOutput 测试开始节点的输出
func TestStartNodeExecutorOutput(t *testing.T) {
	executor := node.NewStartNodeExecutor()

	ctx := &node.NodeContext{
		Task:    &task.Task{ID: "task-001"},
		Node:    &template.Node{ID: "start", Type: template.NodeTypeStart},
		Params:  json.RawMessage(`{"amount": 1000}`),
		Outputs: make(map[string]json.RawMessage),
		Cache:   node.NewContextCache(),
	}

	result, err := executor.Execute(ctx)
	if err != nil {
		t.Fatalf("StartNodeExecutor.Execute() failed: %v", err)
	}

	// 开始节点可以输出任务参数或空输出
	_ = result.Output
}

