package node

import (
	"encoding/json"
	"time"

	"github.com/mautops/approval-kit/internal/template"
)

// StartNodeExecutor 开始节点执行器
// 开始节点用于标识审批流程的起点,通常自动执行
type StartNodeExecutor struct{}

// NewStartNodeExecutor 创建新的开始节点执行器
func NewStartNodeExecutor() NodeExecutor {
	return &StartNodeExecutor{}
}

// NodeType 返回节点类型
func (e *StartNodeExecutor) NodeType() template.NodeType {
	return template.NodeTypeStart
}

// Execute 执行开始节点逻辑
// 开始节点通常自动执行,输出任务参数供后续节点使用
func (e *StartNodeExecutor) Execute(ctx *NodeContext) (*NodeResult, error) {
	// 生成节点激活事件
	event := Event{
		Type: EventTypeNodeActivated,
		Time: time.Now(),
		Data: json.RawMessage(`{"node_id": "` + ctx.Node.ID + `"}`),
	}

	// 开始节点输出任务参数(或空对象)
	output := ctx.Params
	if output == nil || len(output) == 0 {
		output = json.RawMessage("{}")
	}

	return &NodeResult{
		NextNodeID: "", // 下一个节点由流程引擎根据边的定义决定
		Output:     output,
		Events:     []Event{event},
	}, nil
}

