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

// TestNodeOutputStorage 测试节点输出数据存储
func TestNodeOutputStorage(t *testing.T) {
	// 创建任务
	tsk := &task.Task{
		ID:          "task-001",
		TemplateID:  "tpl-001",
		State:       types.TaskStateApproving,
		CurrentNode: "node-1",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		NodeOutputs: make(map[string]json.RawMessage),
	}

	// 设置节点输出数据
	output1 := json.RawMessage(`{"result": "approved", "amount": 1000}`)
	tsk.NodeOutputs["node-1"] = output1

	// 验证输出数据已存储
	stored, exists := tsk.NodeOutputs["node-1"]
	if !exists {
		t.Error("Node output not found")
	}

	if string(stored) != string(output1) {
		t.Errorf("Node output = %q, want %q", string(stored), string(output1))
	}
}

// TestNodeOutputInContext 测试节点输出数据在上下文中的使用
func TestNodeOutputInContext(t *testing.T) {
	// 创建任务
	tsk := &task.Task{
		ID:          "task-001",
		TemplateID:  "tpl-001",
		State:       types.TaskStateApproving,
		CurrentNode: "node-2",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		NodeOutputs: make(map[string]json.RawMessage),
	}

	// 设置前面节点的输出数据
	output1 := json.RawMessage(`{"result": "approved", "amount": 1000}`)
	tsk.NodeOutputs["node-1"] = output1

	// 创建节点上下文
	ctx := &node.NodeContext{
		Task:    tsk,
		Node:    &template.Node{ID: "node-2", Type: template.NodeTypeApproval},
		Params:  json.RawMessage(`{}`),
		Outputs: make(map[string]json.RawMessage),
		Cache:   node.NewContextCache(),
	}

	// 从任务中获取前面节点的输出数据
	ctx.Outputs["node-1"] = tsk.NodeOutputs["node-1"]

	// 验证输出数据在上下文中可用
	output, exists := ctx.Outputs["node-1"]
	if !exists {
		t.Error("Node output not found in context")
	}

	if string(output) != string(output1) {
		t.Errorf("Context output = %q, want %q", string(output), string(output1))
	}
}

// TestNodeOutputDataPassing 测试节点输出数据传递
func TestNodeOutputDataPassing(t *testing.T) {
	// 创建任务
	tsk := &task.Task{
		ID:          "task-001",
		TemplateID:  "tpl-001",
		State:       types.TaskStateApproving,
		CurrentNode: "node-2",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		NodeOutputs: make(map[string]json.RawMessage),
	}

	// 模拟节点1的输出
	output1 := json.RawMessage(`{"result": "approved", "amount": 1000}`)
	tsk.NodeOutputs["node-1"] = output1

	// 模拟节点2使用节点1的输出数据
	ctx := &node.NodeContext{
		Task:    tsk,
		Node:    &template.Node{ID: "node-2", Type: template.NodeTypeApproval},
		Params:  json.RawMessage(`{}`),
		Outputs: make(map[string]json.RawMessage),
		Cache:   node.NewContextCache(),
	}

	// 从任务中获取前面节点的输出
	ctx.Outputs["node-1"] = tsk.NodeOutputs["node-1"]

	// 解析输出数据
	var data map[string]interface{}
	if err := json.Unmarshal(ctx.Outputs["node-1"], &data); err != nil {
		t.Fatalf("Failed to unmarshal output: %v", err)
	}

	// 验证数据内容
	if data["result"] != "approved" {
		t.Errorf("Output result = %v, want %v", data["result"], "approved")
	}
	if data["amount"] != float64(1000) {
		t.Errorf("Output amount = %v, want %v", data["amount"], float64(1000))
	}
}

// TestNodeOutputMultipleNodes 测试多个节点的输出数据
func TestNodeOutputMultipleNodes(t *testing.T) {
	// 创建任务
	tsk := &task.Task{
		ID:          "task-001",
		TemplateID:  "tpl-001",
		State:       types.TaskStateApproving,
		CurrentNode: "node-3",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		NodeOutputs: make(map[string]json.RawMessage),
	}

	// 设置多个节点的输出数据
	output1 := json.RawMessage(`{"result": "approved"}`)
	output2 := json.RawMessage(`{"amount": 1000}`)
	tsk.NodeOutputs["node-1"] = output1
	tsk.NodeOutputs["node-2"] = output2

	// 验证所有输出数据都已存储
	if len(tsk.NodeOutputs) != 2 {
		t.Errorf("NodeOutputs count = %d, want %d", len(tsk.NodeOutputs), 2)
	}

	// 验证每个节点的输出数据
	if string(tsk.NodeOutputs["node-1"]) != string(output1) {
		t.Errorf("Node-1 output = %q, want %q", string(tsk.NodeOutputs["node-1"]), string(output1))
	}
	if string(tsk.NodeOutputs["node-2"]) != string(output2) {
		t.Errorf("Node-2 output = %q, want %q", string(tsk.NodeOutputs["node-2"]), string(output2))
	}
}

