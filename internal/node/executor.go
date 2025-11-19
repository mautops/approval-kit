package node

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
)

// NodeExecutor 节点执行器接口
// 不同类型的节点有不同的执行器实现
type NodeExecutor interface {
	// Execute 执行节点逻辑
	// ctx: 节点执行上下文,包含任务信息、节点信息、参数等
	// 返回: 节点执行结果和错误信息
	Execute(ctx *NodeContext) (*NodeResult, error)

	// NodeType 返回节点类型
	// 用于标识执行器对应的节点类型
	NodeType() template.NodeType
}

// NodeContext 节点执行上下文
// 节点执行时提供上下文信息,包含任务参数、节点输出数据等
type NodeContext struct {
	// Task 任务对象
	Task *task.Task

	// Node 当前节点
	Node *template.Node

	// Params 任务参数(JSON 格式)
	Params json.RawMessage

	// Outputs 前面节点的输出数据
	// 键为节点 ID,值为节点输出数据(JSON 格式)
	Outputs map[string]json.RawMessage

	// Cache 上下文缓存,避免重复查询
	Cache *ContextCache
}

// ContextCache 上下文缓存
// 用于缓存上下文数据,避免重复查询
type ContextCache struct {
	mu   sync.RWMutex
	data map[string]interface{}
}

// NewContextCache 创建新的上下文缓存
func NewContextCache() *ContextCache {
	return &ContextCache{
		data: make(map[string]interface{}),
	}
}

// Get 从缓存中获取数据
func (c *ContextCache) Get(key string) (interface{}, bool) {
	if c == nil {
		return nil, false
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists := c.data[key]
	return value, exists
}

// Set 设置缓存数据
func (c *ContextCache) Set(key string, value interface{}) {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	c.data[key] = value
}

// NodeResult 节点执行结果
type NodeResult struct {
	// NextNodeID 下一个节点 ID
	// 如果为空,表示流程结束
	NextNodeID string

	// Output 节点输出数据(JSON 格式)
	// 供后续节点使用
	Output json.RawMessage

	// Events 生成的事件列表
	// 节点执行过程中可能产生的事件
	Events []Event
}

// Event 事件定义
// 用于节点执行过程中产生的事件通知
type Event struct {
	// Type 事件类型
	Type EventType

	// Time 事件时间
	Time time.Time

	// Data 事件数据(JSON 格式)
	Data json.RawMessage
}

// EventType 事件类型
type EventType string

const (
	// EventTypeNodeActivated 节点激活事件
	EventTypeNodeActivated EventType = "node_activated"

	// EventTypeNodeCompleted 节点完成事件
	EventTypeNodeCompleted EventType = "node_completed"
)
