package node

import (
	"encoding/json"
	"time"

	"github.com/mautops/approval-kit/internal/template"
)

// EndNodeExecutor 结束节点执行器
// 结束节点用于标识审批流程的终点
type EndNodeExecutor struct{}

// NewEndNodeExecutor 创建新的结束节点执行器
func NewEndNodeExecutor() NodeExecutor {
	return &EndNodeExecutor{}
}

// NodeType 返回节点类型
func (e *EndNodeExecutor) NodeType() template.NodeType {
	return template.NodeTypeEnd
}

// Execute 执行结束节点逻辑
// 结束节点表示流程完成,不指定下一个节点
func (e *EndNodeExecutor) Execute(ctx *NodeContext) (*NodeResult, error) {
	// 生成节点完成事件
	event := Event{
		Type: EventTypeNodeCompleted,
		Time: time.Now(),
		Data: json.RawMessage(`{"node_id": "` + ctx.Node.ID + `"}`),
	}

	// 结束节点可以输出最终结果(或空对象)
	output := json.RawMessage("{}")

	return &NodeResult{
		NextNodeID: "", // 结束节点没有下一个节点
		Output:     output,
		Events:     []Event{event},
	}, nil
}

