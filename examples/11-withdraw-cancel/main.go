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
	fmt.Println("=== 场景 11: 任务撤回和取消场景 ===")
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

	// 3. 演示任务撤回
	demonstrateWithdraw(templateMgr, taskMgr)

	// 4. 演示任务取消
	demonstrateCancel(templateMgr, taskMgr)
}

// createTemplate 创建审批模板
func createTemplate() *template.Template {
	now := time.Now()
	return &template.Template{
		ID:          "withdraw-cancel-template",
		Name:        "撤回和取消模板",
		Description: "演示任务撤回和取消功能",
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

// demonstrateWithdraw 演示任务撤回
func demonstrateWithdraw(templateMgr template.TemplateManager, taskMgr task.TaskManager) {
	fmt.Println("=== 演示 1: 任务撤回 ===")
	fmt.Println()

	// 创建任务
	fmt.Println("步骤 1: 创建任务")
	params := json.RawMessage(`{
		"requestNo": "REQ-2025-001",
		"requestType": "撤回测试",
		"requester": "申请人-001",
		"amount": 10000,
		"description": "测试任务撤回功能"
	}`)
	tsk, err := taskMgr.Create("withdraw-cancel-template", "withdraw-001", params)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}
	fmt.Printf("✓ 任务创建成功: ID=%s, State=%s\n\n", tsk.ID, tsk.State)

	// 提交任务
	fmt.Println("步骤 2: 提交任务")
	err = taskMgr.Submit(tsk.ID)
	if err != nil {
		log.Fatalf("Failed to submit task: %v", err)
	}
	tsk, _ = taskMgr.Get(tsk.ID)
	fmt.Printf("✓ 任务已提交: State=%s\n\n", tsk.State)

	// 撤回任务
	fmt.Println("步骤 3: 撤回任务")
	fmt.Println("  说明: 在未产生审批记录前可以撤回任务")
	fmt.Println("  撤回后任务状态回退到 pending,SubmittedAt 被清空")
	fmt.Println()
	err = taskMgr.Withdraw(tsk.ID, "提交后发现错误,需要修改后重新提交")
	if err != nil {
		log.Fatalf("Failed to withdraw task: %v", err)
	}
	tsk, _ = taskMgr.Get(tsk.ID)
	fmt.Printf("✓ 任务已撤回: State=%s\n", tsk.State)
	if tsk.SubmittedAt == nil {
		fmt.Println("  SubmittedAt 已清空")
	}
	fmt.Println()

	// 输出撤回后的任务信息
	printTaskInfo(tsk, "撤回后的任务信息")
	fmt.Println()
}

// demonstrateCancel 演示任务取消
func demonstrateCancel(templateMgr template.TemplateManager, taskMgr task.TaskManager) {
	fmt.Println("=== 演示 2: 任务取消 ===")
	fmt.Println()

	// 创建任务
	fmt.Println("步骤 1: 创建任务")
	params := json.RawMessage(`{
		"requestNo": "REQ-2025-002",
		"requestType": "取消测试",
		"requester": "申请人-001",
		"amount": 20000,
		"description": "测试任务取消功能"
	}`)
	tsk, err := taskMgr.Create("withdraw-cancel-template", "cancel-001", params)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}
	fmt.Printf("✓ 任务创建成功: ID=%s, State=%s\n\n", tsk.ID, tsk.State)

	// 提交任务
	fmt.Println("步骤 2: 提交任务")
	err = taskMgr.Submit(tsk.ID)
	if err != nil {
		log.Fatalf("Failed to submit task: %v", err)
	}
	tsk, _ = taskMgr.Get(tsk.ID)
	fmt.Printf("✓ 任务已提交: State=%s\n\n", tsk.State)

	// 设置审批人
	fmt.Println("步骤 3: 设置审批人")
	err = taskMgr.AddApprover(tsk.ID, "approval", "manager-001", "设置审批人")
	if err != nil {
		log.Fatalf("Failed to add approver: %v", err)
	}
	fmt.Printf("✓ 审批人已设置: manager-001\n\n")

	// 取消任务
	fmt.Println("步骤 4: 取消任务")
	fmt.Println("  说明: 任务可以随时取消(在 pending、submitted 或 approving 状态)")
	fmt.Println("  取消后任务状态变为 cancelled,无法继续审批")
	fmt.Println()
	err = taskMgr.Cancel(tsk.ID, "申请人决定取消此申请")
	if err != nil {
		log.Fatalf("Failed to cancel task: %v", err)
	}
	tsk, _ = taskMgr.Get(tsk.ID)
	fmt.Printf("✓ 任务已取消: State=%s\n\n", tsk.State)

	// 输出取消后的任务信息
	printTaskInfo(tsk, "取消后的任务信息")
	fmt.Println()

	// 验证取消后无法继续操作
	fmt.Println("步骤 5: 验证取消后无法继续操作")
	fmt.Println("  说明: 取消后的任务无法继续审批")
	err = taskMgr.Approve(tsk.ID, "approval", "manager-001", "尝试审批已取消的任务")
	if err != nil {
		fmt.Printf("  ✓ 验证通过: 已取消的任务无法继续审批 (错误: %v)\n", err)
	} else {
		fmt.Println("  ✗ 验证失败: 已取消的任务不应该可以继续审批")
	}
	fmt.Println()
}

// printTaskInfo 输出任务信息
func printTaskInfo(tsk *task.Task, title string) {
	fmt.Printf("=== %s ===\n", title)
	fmt.Printf("任务 ID: %s\n", tsk.ID)
	fmt.Printf("业务 ID: %s\n", tsk.BusinessID)
	fmt.Printf("模板 ID: %s\n", tsk.TemplateID)
	fmt.Printf("当前状态: %s\n", tsk.State)
	fmt.Printf("当前节点: %s\n", tsk.CurrentNode)
	fmt.Printf("创建时间: %s\n", tsk.CreatedAt.Format("2006-01-02 15:04:05"))
	if tsk.SubmittedAt != nil {
		fmt.Printf("提交时间: %s\n", tsk.SubmittedAt.Format("2006-01-02 15:04:05"))
	} else {
		fmt.Printf("提交时间: 未提交\n")
	}
	fmt.Printf("更新时间: %s\n", tsk.UpdatedAt.Format("2006-01-02 15:04:05"))

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

	// 输出审批记录
	fmt.Println("\n审批记录:")
	if len(tsk.Records) == 0 {
		fmt.Println("  无审批记录")
	} else {
		for i, record := range tsk.Records {
			fmt.Printf("  记录 %d:\n", i+1)
			fmt.Printf("    节点: %s\n", record.NodeID)
			fmt.Printf("    审批人: %s\n", record.Approver)
			fmt.Printf("    操作类型: %s\n", getOperationType(record.Result))
			fmt.Printf("    操作说明: %s\n", record.Comment)
			fmt.Printf("    操作时间: %s\n", record.CreatedAt.Format("2006-01-02 15:04:05"))
		}
	}
}

// getOperationType 获取操作类型描述
func getOperationType(result string) string {
	operationTypes := map[string]string{
		"approve":         "审批通过",
		"reject":          "审批拒绝",
		"transfer":        "转交审批",
		"add_approver":    "加签",
		"remove_approver": "减签",
	}
	if name, ok := operationTypes[result]; ok {
		return name
	}
	return result
}
