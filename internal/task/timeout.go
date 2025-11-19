package task

import (
	"time"

	"github.com/mautops/approval-kit/internal/template"
	"github.com/mautops/approval-kit/internal/types"
)

// CheckTimeout 检查任务是否超时
// 返回: 是否超时,超时的节点 ID(如果超时)
func (m *memoryTaskManager) CheckTimeout(tsk *Task) (bool, string) {
	// 获取当前节点
	tsk.mu.RLock()
	currentNodeID := tsk.CurrentNode
	state := tsk.State
	tsk.mu.RUnlock()

	// 只有 submitted 或 approving 状态的任务才需要检查超时
	if state != types.TaskStateSubmitted && state != types.TaskStateApproving {
		return false, ""
	}

	// 获取模板
	tpl, err := m.templateMgr.Get(tsk.TemplateID, 0)
	if err != nil {
		return false, ""
	}

	// 获取当前节点
	node, exists := tpl.Nodes[currentNodeID]
	if !exists {
		return false, ""
	}

	// 检查节点配置
	if node.Type != template.NodeTypeApproval {
		return false, ""
	}

	// 获取审批节点配置
	config, ok := node.Config.(template.ApprovalNodeConfigAccessor)
	if !ok {
		return false, ""
	}

	// 检查是否配置了超时
	timeout := config.GetTimeout()
	if timeout == nil {
		return false, ""
	}

	// 检查节点是否已超时
	// 需要找到节点激活时间或任务提交时间
	tsk.mu.RLock()
	submittedAt := tsk.SubmittedAt
	tsk.mu.RUnlock()

	// 如果没有提交时间,使用创建时间
	startTime := tsk.CreatedAt
	if submittedAt != nil {
		startTime = *submittedAt
	}

	// 检查是否超时
	if time.Since(startTime) > *timeout {
		return true, currentNodeID
	}

	return false, ""
}
