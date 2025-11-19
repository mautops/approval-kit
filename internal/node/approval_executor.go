package node

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mautops/approval-kit/internal/errors"
	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
)

// ApprovalNodeExecutor 审批节点执行器
// 根据审批模式执行不同的审批逻辑
type ApprovalNodeExecutor struct {
	config   *ApprovalNodeConfig
	registry ApprovalModeHandlerRegistry
}

// NewApprovalNodeExecutor 创建新的审批节点执行器
func NewApprovalNodeExecutor(config *ApprovalNodeConfig) NodeExecutor {
	return &ApprovalNodeExecutor{
		config:   config,
		registry: NewApprovalModeHandlerRegistry(),
	}
}

// NodeType 返回节点类型(实现 NodeExecutor 接口)
func (e *ApprovalNodeExecutor) NodeType() template.NodeType {
	return template.NodeTypeApproval
}

// Execute 执行审批节点逻辑(实现 NodeExecutor 接口)
func (e *ApprovalNodeExecutor) Execute(ctx *NodeContext) (*NodeResult, error) {
	// 1. 获取审批人列表
	approvers, err := e.config.ApproverConfig.GetApprovers(ctx)
	if err != nil {
		return nil, err
	}

	// 2. 获取当前节点的审批记录
	nodeID := ctx.Node.ID
	approvals, exists := ctx.Task.Approvals[nodeID]
	if !exists {
		approvals = make(map[string]*task.Approval)
	}

	// 3. 根据审批模式检查审批状态
	handler := e.registry.GetHandler(e.config.Mode)
	if handler == nil {
		return nil, fmt.Errorf("no handler found for approval mode: %q", e.config.Mode)
	}

	completed, result := handler.CheckCompletion(approvers, approvals, e.config)

	// 4. 如果审批未完成,返回错误
	if !completed {
		return nil, errors.ErrApprovalPending
	}

	// 5. 审批完成,生成节点完成事件
	event := Event{
		Type: EventTypeNodeCompleted,
		Time: time.Now(),
		Data: json.RawMessage(`{"node_id": "` + nodeID + `", "result": "` + result.Result + `"}`),
	}

	// 6. 生成输出数据
	output := json.RawMessage(`{"result": "` + result.Result + `"}`)

	return &NodeResult{
		NextNodeID: result.NextNodeID,
		Output:     output,
		Events:     []Event{event},
	}, nil
}

// ApprovalResult 审批结果
type ApprovalResult struct {
	// Completed 是否完成
	Completed bool

	// Result 审批结果(approve/reject)
	Result string

	// NextNodeID 下一个节点 ID
	NextNodeID string
}

