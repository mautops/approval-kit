package task

import (
	"fmt"
	"time"

	"github.com/mautops/approval-kit/internal/event"
	"github.com/mautops/approval-kit/internal/template"
)

// generateEventID 生成事件 ID(用于幂等性保证)
func generateEventID(taskID string, eventType event.EventType, t time.Time) string {
	return fmt.Sprintf("%s-%s-%d", taskID, eventType, t.UnixNano())
}

// generateEvent 生成事件
func (m *memoryTaskManager) generateEvent(eventType event.EventType, tsk *Task, node *template.Node, approval *event.ApprovalInfo) {
	if m.eventNotifier == nil {
		return
	}

	// 获取节点信息
	var nodeInfo *event.NodeInfo
	if node != nil {
		nodeInfo = &event.NodeInfo{
			ID:   node.ID,
			Name: node.Name,
			Type: string(node.Type),
		}
	} else {
		// 如果没有节点信息,使用当前节点
		tsk.mu.RLock()
		currentNodeID := tsk.CurrentNode
		tsk.mu.RUnlock()

		// 从模板获取节点信息
		if tpl, err := m.templateMgr.Get(tsk.TemplateID, 0); err == nil {
			if n, exists := tpl.Nodes[currentNodeID]; exists {
				nodeInfo = &event.NodeInfo{
					ID:   n.ID,
					Name: n.Name,
					Type: string(n.Type),
				}
			}
		}

		// 如果仍然没有节点信息,使用默认值
		if nodeInfo == nil {
			nodeInfo = &event.NodeInfo{
				ID:   currentNodeID,
				Name: currentNodeID,
				Type: "unknown",
			}
		}
	}

	// 构建任务信息
	tsk.mu.RLock()
	taskInfo := &event.TaskInfo{
		ID:         tsk.ID,
		TemplateID: tsk.TemplateID,
		BusinessID: tsk.BusinessID,
		State:      string(tsk.State),
	}
	tsk.mu.RUnlock()

	// 构建业务信息
	businessInfo := &event.BusinessInfo{
		ID:   tsk.BusinessID,
		Data: tsk.Params,
	}

	// 生成事件 ID(用于幂等性保证)
	eventID := generateEventID(tsk.ID, eventType, time.Now())

	// 创建事件
	evt := &event.Event{
		ID:        eventID,
		Type:      eventType,
		Time:      time.Now(),
		Task:      taskInfo,
		Node:      nodeInfo,
		Approval:  approval,
		Business:  businessInfo,
	}

	// 异步推送事件
	m.eventNotifier.Notify(evt)
}

