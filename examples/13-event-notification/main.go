package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"

	"github.com/mautops/approval-kit/internal/event"
	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
)

func main() {
	fmt.Println("=== 场景 13: 事件通知集成场景 ===")
	fmt.Println()

	// 1. 创建 Mock Webhook 服务器
	fmt.Println("步骤 1: 创建 Mock Webhook 服务器")
	receivedEvents := make([]*event.Event, 0)
	var mu sync.Mutex

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var evt event.Event
		if err := json.NewDecoder(r.Body).Decode(&evt); err != nil {
			log.Printf("Failed to decode event: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		mu.Lock()
		receivedEvents = append(receivedEvents, &evt)
		mu.Unlock()

		fmt.Printf("  [Webhook] 收到事件: Type=%s, TaskID=%s, Time=%s\n",
			evt.Type, evt.Task.ID, evt.Time.Format("2006-01-02 15:04:05"))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	fmt.Printf("✓ Webhook 服务器已启动: %s\n\n", server.URL)

	// 2. 创建 Webhook 处理器
	fmt.Println("步骤 2: 创建 Webhook 处理器")
	webhookConfig := &event.WebhookConfig{
		URL:    server.URL,
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/json",
			"X-Auth-Token":  "test-token-123",
		},
		Timeout: 30,
	}
	webhookHandler := event.NewWebhookHandler(webhookConfig)
	fmt.Printf("✓ Webhook 处理器已创建\n\n")

	// 3. 创建事件通知器
	fmt.Println("步骤 3: 创建事件通知器")
	notifier := event.NewEventNotifier([]event.EventHandler{webhookHandler}, 100)
	fmt.Printf("✓ 事件通知器已创建\n\n")

	// 4. 创建管理器
	fmt.Println("步骤 4: 创建管理器")
	templateMgr := template.NewTemplateManager()
	taskMgr := task.NewTaskManagerWithNotifier(templateMgr, nil, notifier)
	fmt.Printf("✓ 管理器已创建\n\n")

	// 5. 创建模板
	fmt.Println("步骤 5: 创建审批模板")
	tpl := createTemplate()
	err := templateMgr.Create(tpl)
	if err != nil {
		log.Fatalf("Failed to create template: %v", err)
	}
	fmt.Printf("✓ 模板创建成功: ID=%s, Name=%s\n\n", tpl.ID, tpl.Name)

	// 6. 执行场景流程
	runScenario(taskMgr, tpl)

	// 7. 等待事件推送完成
	fmt.Println("步骤 6: 等待事件推送完成")
	time.Sleep(500 * time.Millisecond) // 等待异步事件推送完成
	fmt.Println()

	// 8. 输出接收到的事件
	fmt.Println("=== 接收到的 Webhook 事件 ===")
	mu.Lock()
	if len(receivedEvents) == 0 {
		fmt.Println("  无事件接收")
	} else {
		fmt.Printf("  共接收到 %d 个事件:\n\n", len(receivedEvents))
		for i, evt := range receivedEvents {
			fmt.Printf("  事件 %d:\n", i+1)
			fmt.Printf("    事件类型: %s\n", evt.Type)
			fmt.Printf("    事件时间: %s\n", evt.Time.Format("2006-01-02 15:04:05"))
			fmt.Printf("    任务 ID: %s\n", evt.Task.ID)
			fmt.Printf("    业务 ID: %s\n", evt.Task.BusinessID)
			fmt.Printf("    任务状态: %s\n", evt.Task.State)
			if evt.Node != nil {
				fmt.Printf("    节点信息: %s (%s)\n", evt.Node.Name, evt.Node.ID)
			}
			if evt.Approval != nil {
				fmt.Printf("    审批信息: %s - %s\n", evt.Approval.Approver, evt.Approval.Result)
			}
			fmt.Println()
		}
	}
	mu.Unlock()

	// 9. 停止事件通知器
	notifier.Stop()
	fmt.Println("✓ 事件通知器已停止")
}

// createTemplate 创建审批模板
func createTemplate() *template.Template {
	now := time.Now()
	return &template.Template{
		ID:          "event-notification-template",
		Name:        "事件通知模板",
		Description: "演示事件通知功能",
		Version:     1,
		CreatedAt:   now,
		UpdatedAt:   now,
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "开始",
				Type: template.NodeTypeStart,
			},
			"approval": {
				ID:   "approval",
				Name: "审批",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeSingle,
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{"manager-001"},
					},
					Permissions: node.OperationPermissions{
						AllowAddApprover: true,
					},
				},
			},
			"end": {
				ID:   "end",
				Name: "结束",
				Type: template.NodeTypeEnd,
			},
		},
		Edges: []*template.Edge{
			{From: "start", To: "approval"},
			{From: "approval", To: "end"},
		},
	}
}

// runScenario 执行场景流程
func runScenario(taskMgr task.TaskManager, tpl *template.Template) {
	fmt.Println("=== 执行审批流程 ===")
	fmt.Println()

	// 创建任务
	fmt.Println("步骤 1: 创建任务")
	params := json.RawMessage(`{
		"requestNo": "REQ-2025-001",
		"requestType": "事件通知测试",
		"requester": "申请人-001",
		"amount": 50000,
		"description": "测试事件通知功能"
	}`)
	tsk, err := taskMgr.Create(tpl.ID, "event-001", params)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}
	fmt.Printf("✓ 任务创建成功: ID=%s\n", tsk.ID)
	fmt.Println("  说明: 任务创建时会触发 task_created 事件")
	fmt.Println()

	// 等待事件推送
	time.Sleep(100 * time.Millisecond)

	// 提交任务
	fmt.Println("步骤 2: 提交任务")
	err = taskMgr.Submit(tsk.ID)
	if err != nil {
		log.Fatalf("Failed to submit task: %v", err)
	}
	fmt.Printf("✓ 任务已提交: ID=%s\n", tsk.ID)
	fmt.Println("  说明: 任务提交时会触发 task_submitted 事件")
	fmt.Println()

	// 等待事件推送
	time.Sleep(100 * time.Millisecond)

	// 设置审批人
	fmt.Println("步骤 3: 设置审批人")
	err = taskMgr.AddApprover(tsk.ID, "approval", "manager-001", "设置审批人")
	if err != nil {
		log.Fatalf("Failed to add approver: %v", err)
	}
	fmt.Printf("✓ 审批人已设置: manager-001\n\n")

	// 等待事件推送
	time.Sleep(100 * time.Millisecond)

	// 审批操作
	fmt.Println("步骤 4: 审批操作")
	err = taskMgr.Approve(tsk.ID, "approval", "manager-001", "审批通过")
	if err != nil {
		log.Fatalf("Failed to approve: %v", err)
	}
	tsk, _ = taskMgr.Get(tsk.ID)
	fmt.Printf("✓ 审批已通过: ID=%s, State=%s\n", tsk.ID, tsk.State)
	fmt.Println("  说明: 审批操作时会触发以下事件:")
	fmt.Println("    - approval_operation: 审批操作事件")
	fmt.Println("    - node_completed: 节点完成事件")
	fmt.Println("    - task_approved: 任务通过事件")
	fmt.Println()
}

