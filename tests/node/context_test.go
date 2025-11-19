package node_test

import (
	"encoding/json"
	"testing"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
	"github.com/mautops/approval-kit/internal/types"
)

// TestNodeContextStruct 验证 NodeContext 结构体定义
func TestNodeContextStruct(t *testing.T) {
	// 验证 NodeContext 类型存在
	var ctx *node.NodeContext
	if ctx != nil {
		_ = ctx
	}
}

// TestNodeContextFields 验证 NodeContext 结构体的所有字段
func TestNodeContextFields(t *testing.T) {
	tsk := &task.Task{
		ID:    "task-001",
		State: types.TaskStatePending,
	}

	tplNode := &template.Node{
		ID:   "node-001",
		Name: "Test Node",
		Type: template.NodeTypeStart,
	}

	params := json.RawMessage(`{"amount": 1000}`)
	outputs := make(map[string]json.RawMessage)
	outputs["prev-node"] = json.RawMessage(`{"result": "ok"}`)

	ctx := &node.NodeContext{
		Task:    tsk,
		Node:    tplNode,
		Params:  params,
		Outputs: outputs,
		Cache:   node.NewContextCache(),
	}

	// 验证字段值
	if ctx.Task.ID != "task-001" {
		t.Errorf("NodeContext.Task.ID = %q, want %q", ctx.Task.ID, "task-001")
	}
	if ctx.Node.ID != "node-001" {
		t.Errorf("NodeContext.Node.ID = %q, want %q", ctx.Node.ID, "node-001")
	}
	if len(ctx.Params) == 0 {
		t.Error("NodeContext.Params should not be empty")
	}
	if len(ctx.Outputs) != 1 {
		t.Errorf("NodeContext.Outputs length = %d, want 1", len(ctx.Outputs))
	}
	if ctx.Cache == nil {
		t.Error("NodeContext.Cache should not be nil")
	}
}

// TestContextCache 测试上下文缓存功能
func TestContextCache(t *testing.T) {
	cache := node.NewContextCache()

	// 测试 Set 和 Get
	cache.Set("key1", "value1")
	cache.Set("key2", 123)

	value1, exists1 := cache.Get("key1")
	if !exists1 {
		t.Error("Cache.Get() should return true for existing key")
	}
	if value1 != "value1" {
		t.Errorf("Cache.Get() returned value = %v, want %v", value1, "value1")
	}

	value2, exists2 := cache.Get("key2")
	if !exists2 {
		t.Error("Cache.Get() should return true for existing key")
	}
	if value2 != 123 {
		t.Errorf("Cache.Get() returned value = %v, want %v", value2, 123)
	}

	// 测试不存在的 key
	_, exists3 := cache.Get("non-existent")
	if exists3 {
		t.Error("Cache.Get() should return false for non-existent key")
	}
}

// TestContextCacheConcurrent 在 context_cache_test.go 中定义

