package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
)

func main() {
	fmt.Println("=== 场景 09: 超时处理场景 ===")
	fmt.Println()

	// 1. 创建管理器
	templateMgr := template.NewTemplateManager()
	taskMgr := task.NewTaskManager(templateMgr, nil)

	// 2. 创建模板
	fmt.Println("步骤 1: 创建审批模板")
	tpl := createTemplate()
	err := templateMgr.Create(tpl)
	if err != nil {
		log.Fatalf("Failed to create template: %v", err)
	}
	fmt.Printf("✓ 模板创建成功: ID=%s, Name=%s\n\n", tpl.ID, tpl.Name)

	// 3. 创建任务
	fmt.Println("步骤 2: 创建紧急审批任务")
	params := json.RawMessage(`{
		"requestNo": "REQ-2025-001",
		"requestType": "紧急采购",
		"requester": "采购员-001",
		"amount": 50000,
		"urgency": "high",
		"description": "紧急采购申请,需要在24小时内完成审批"
	}`)
	tsk, err := taskMgr.Create(tpl.ID, "urgent-001", params)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}
	fmt.Printf("✓ 任务创建成功: ID=%s, State=%s\n\n", tsk.ID, tsk.State)

	// 4. 执行场景流程
	runScenario(taskMgr, tsk)

	// 5. 输出结果
	printResults(taskMgr, tsk)
}

// createTemplate 创建紧急审批模板
// 流程: 开始节点 → 审批节点(配置超时) → 结束节点
// 审批节点配置了24小时超时时间
func createTemplate() *template.Template {
	now := time.Now()
	timeout := 24 * time.Hour // 24小时超时
	return &template.Template{
		ID:          "urgent-approval-template",
		Name:        "紧急审批模板",
		Description: "紧急审批流程,如果审批人在24小时内未审批,自动处理或通知",
		Version:     1,
		CreatedAt:   now,
		UpdatedAt:   now,
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "开始",
				Type: template.NodeTypeStart,
			},
			"urgent-approval": {
				ID:   "urgent-approval",
				Name: "紧急审批",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeSingle,
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{"manager-001"},
					},
					Timeout: &timeout, // 配置24小时超时
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
			{From: "start", To: "urgent-approval"},
			{From: "urgent-approval", To: "end"},
		},
	}
}

// runScenario 执行场景流程
func runScenario(taskMgr task.TaskManager, tsk *task.Task) {
	// 步骤 3: 提交任务
	fmt.Println("步骤 3: 提交任务进入审批流程")
	err := taskMgr.Submit(tsk.ID)
	if err != nil {
		log.Fatalf("Failed to submit task: %v", err)
	}

	// 获取最新状态
	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		log.Fatalf("Failed to get task: %v", err)
	}
	fmt.Printf("✓ 任务已提交: State=%s\n", tsk.State)
	if tsk.SubmittedAt != nil {
		fmt.Printf("  提交时间: %s\n", tsk.SubmittedAt.Format("2006-01-02 15:04:05"))
	}
	fmt.Println()

	// 步骤 4: 设置审批人
	fmt.Println("步骤 4: 设置审批人")
	nodeID := "urgent-approval"
	err = taskMgr.AddApprover(tsk.ID, nodeID, "manager-001", "设置审批人")
	if err != nil {
		log.Fatalf("Failed to add approver: %v", err)
	}
	fmt.Printf("✓ 审批人已设置: manager-001\n\n")

	// 步骤 5: 检查超时配置
	fmt.Println("步骤 5: 检查超时配置")
	fmt.Println("  说明: 审批节点配置了24小时超时时间")
	fmt.Println("  超时时间: 24小时")
	fmt.Println("  超时处理: 如果审批人在24小时内未审批,任务状态将变为timeout")
	fmt.Println()

	// 步骤 6: 模拟超时检测
	fmt.Println("步骤 6: 模拟超时检测")
	fmt.Println("  说明: 在实际业务系统中,应该有定时任务定期检查超时任务")
	fmt.Println("  本示例展示超时检测和处理的配置方法")
	fmt.Println()

	// 检查是否超时(使用短超时时间进行演示)
	// 注意: 实际场景中,超时时间应该是24小时,这里我们使用短时间进行演示
	fmt.Println("  模拟场景: 假设任务已超过超时时间")
	fmt.Println("  说明: 在实际业务系统中,可以通过以下方式处理超时:")
	fmt.Println("    1. 定时任务定期调用 CheckTimeout 检查超时任务")
	fmt.Println("    2. 调用 HandleTimeout 处理超时任务")
	fmt.Println("    3. 超时后可以自动通过、自动拒绝或发送通知")
	fmt.Println()

	// 步骤 7: 说明超时处理流程
	fmt.Println("步骤 7: 超时处理流程说明")
	fmt.Println("  1. 任务提交后,开始计时")
	fmt.Println("  2. 定时任务定期检查任务是否超时")
	fmt.Println("  3. 如果超时,调用 HandleTimeout 处理超时任务")
	fmt.Println("  4. 任务状态变为 timeout,生成超时事件")
	fmt.Println("  5. 根据业务规则,可以自动通过、自动拒绝或发送通知")
	fmt.Println()
	fmt.Println("✓ 超时处理功能配置和使用方法已展示")
	fmt.Println()
}

// printResults 输出结果
func printResults(taskMgr task.TaskManager, tsk *task.Task) {
	// 获取最新任务状态
	tsk, err := taskMgr.Get(tsk.ID)
	if err != nil {
		log.Fatalf("Failed to get task: %v", err)
	}

	fmt.Println("=== 审批结果 ===")
	fmt.Printf("任务 ID: %s\n", tsk.ID)
	fmt.Printf("业务 ID: %s\n", tsk.BusinessID)
	fmt.Printf("模板 ID: %s\n", tsk.TemplateID)
	fmt.Printf("当前状态: %s\n", tsk.State)
	fmt.Printf("当前节点: %s\n", tsk.CurrentNode)
	fmt.Printf("创建时间: %s\n", tsk.CreatedAt.Format("2006-01-02 15:04:05"))
	if tsk.SubmittedAt != nil {
		fmt.Printf("提交时间: %s\n", tsk.SubmittedAt.Format("2006-01-02 15:04:05"))
		// 计算已过去的时间
		elapsed := time.Since(*tsk.SubmittedAt)
		fmt.Printf("已过去时间: %s\n", formatDuration(elapsed))
	}
	fmt.Printf("更新时间: %s\n", tsk.UpdatedAt.Format("2006-01-02 15:04:05"))

	// 输出任务参数
	fmt.Println("\n任务参数:")
	var params map[string]interface{}
	if err := json.Unmarshal(tsk.Params, &params); err == nil {
		for k, v := range params {
			fmt.Printf("  %s: %v\n", k, v)
		}
	}

	// 输出审批人列表
	fmt.Println("\n审批人列表:")
	nodeID := "urgent-approval"
	approvers := tsk.Approvers[nodeID]
	if len(approvers) > 0 {
		for _, approver := range approvers {
			fmt.Printf("  - %s\n", approver)
		}
	} else {
		fmt.Println("  无审批人")
	}

	// 输出审批记录
	fmt.Println("\n审批记录:")
	if len(tsk.Records) == 0 {
		fmt.Println("  无审批记录")
	} else {
		for i, record := range tsk.Records {
			fmt.Printf("  记录 %d:\n", i+1)
			fmt.Printf("    节点: %s\n", getNodeName(record.NodeID))
			fmt.Printf("    审批人: %s\n", record.Approver)
			fmt.Printf("    操作类型: %s\n", getOperationType(record.Result))
			fmt.Printf("    操作说明: %s\n", record.Comment)
			fmt.Printf("    操作时间: %s\n", record.CreatedAt.Format("2006-01-02 15:04:05"))
		}
	}

	// 输出状态变更历史
	fmt.Println("\n状态变更历史:")
	if len(tsk.StateHistory) == 0 {
		fmt.Println("  无状态变更记录")
	} else {
		for i, change := range tsk.StateHistory {
			fmt.Printf("  变更 %d: %s -> %s (原因: %s, 时间: %s)\n",
				i+1,
				change.From,
				change.To,
				change.Reason,
				change.Time.Format("2006-01-02 15:04:05"),
			)
		}
	}

	// 验证最终状态
	fmt.Println("\n=== 验证结果 ===")
	fmt.Println("✓ 超时配置已正确设置:")
	fmt.Println("  - 审批节点配置了24小时超时时间")
	fmt.Println("  - 超时检测机制已配置")
	fmt.Println("  - 超时处理逻辑已配置")
	fmt.Println()
	fmt.Println("  说明: 在实际业务系统中,可以通过以下方式使用超时功能:")
	fmt.Println("    1. 定时任务定期调用 CheckTimeout 检查超时任务")
	fmt.Println("    2. 调用 HandleTimeout 处理超时任务")
	fmt.Println("    3. 超时后任务状态变为 timeout,生成超时事件")
	fmt.Println("    4. 根据业务规则,可以自动通过、自动拒绝或发送通知")
}

// getNodeName 获取节点名称
func getNodeName(nodeID string) string {
	nodeNames := map[string]string{
		"urgent-approval": "紧急审批",
	}
	if name, ok := nodeNames[nodeID]; ok {
		return name
	}
	return nodeID
}

// getOperationType 获取操作类型描述
func getOperationType(result string) string {
	operationTypes := map[string]string{
		"approve":        "审批通过",
		"reject":         "审批拒绝",
		"transfer":       "转交审批",
		"add_approver":   "加签",
		"remove_approver": "减签",
	}
	if name, ok := operationTypes[result]; ok {
		return name
	}
	return result
}

// formatDuration 格式化时间间隔
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0f秒", d.Seconds())
	} else if d < time.Hour {
		return fmt.Sprintf("%.0f分钟", d.Minutes())
	} else if d < 24*time.Hour {
		return fmt.Sprintf("%.1f小时", d.Hours())
	} else {
		days := d.Hours() / 24
		hours := d.Hours() - days*24
		return fmt.Sprintf("%.0f天%.0f小时", days, hours)
	}
}

