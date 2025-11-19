package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
	"github.com/mautops/approval-kit/internal/types"
)

func main() {
	fmt.Println("=== 场景 08: 高级操作场景(转交、加签、减签) ===")
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
	fmt.Println("步骤 2: 创建项目审批任务")
	params := json.RawMessage(`{
		"projectNo": "PRJ-2025-001",
		"projectName": "新产品开发项目",
		"projectManager": "pm-001",
		"budget": 2000000,
		"description": "新产品开发项目审批"
	}`)
	tsk, err := taskMgr.Create(tpl.ID, "project-001", params)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}
	fmt.Printf("✓ 任务创建成功: ID=%s, State=%s\n\n", tsk.ID, tsk.State)

	// 4. 执行场景流程
	runScenario(taskMgr, tsk)

	// 5. 输出结果
	printResults(taskMgr, tsk)
}

// createTemplate 创建项目审批模板
// 流程: 开始节点 → 审批节点(支持高级操作) → 结束节点
// 审批节点配置了转交、加签、减签权限
func createTemplate() *template.Template {
	now := time.Now()
	return &template.Template{
		ID:          "project-approval-template",
		Name:        "项目审批模板",
		Description: "项目审批流程,支持转交、加签、减签等高级操作",
		Version:     1,
		CreatedAt:   now,
		UpdatedAt:   now,
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "开始",
				Type: template.NodeTypeStart,
			},
			"project-approval": {
				ID:   "project-approval",
				Name: "项目审批",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeUnanimous, // 多人会签模式
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{
							"manager-001", // 项目经理
							"finance-001", // 财务审批人
						},
					},
					Permissions: node.OperationPermissions{
						AllowTransfer:      true, // 允许转交
						AllowAddApprover:   true, // 允许加签
						AllowRemoveApprover: true, // 允许减签
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
			{From: "start", To: "project-approval"},
			{From: "project-approval", To: "end"},
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
	fmt.Printf("✓ 任务已提交: State=%s\n\n", tsk.State)

	nodeID := "project-approval"

	// 步骤 4: 设置初始审批人列表
	fmt.Println("步骤 4: 设置初始审批人列表")
	err = taskMgr.AddApprover(tsk.ID, nodeID, "manager-001", "设置项目经理为审批人")
	if err != nil {
		log.Fatalf("Failed to add approver: %v", err)
	}
	err = taskMgr.AddApprover(tsk.ID, nodeID, "finance-001", "设置财务审批人")
	if err != nil {
		log.Fatalf("Failed to add approver: %v", err)
	}
	fmt.Printf("✓ 初始审批人列表已设置: [manager-001, finance-001]\n\n")

	// 步骤 5: 加签操作
	fmt.Println("步骤 5: 加签操作")
	fmt.Println("  说明: 在审批过程中添加额外的审批人")
	err = taskMgr.AddApprover(tsk.ID, nodeID, "tech-lead-001", "加签: 添加技术负责人参与审批")
	if err != nil {
		log.Fatalf("Failed to add approver: %v", err)
	}
	tsk, _ = taskMgr.Get(tsk.ID)
	approvers := tsk.Approvers[nodeID]
	fmt.Printf("✓ 加签成功: 已添加 tech-lead-001\n")
	fmt.Printf("  当前审批人列表: %v (共 %d 人)\n\n", approvers, len(approvers))

	// 步骤 6: 转交操作
	fmt.Println("步骤 6: 转交操作")
	fmt.Println("  说明: 审批人可以将任务转交给其他审批人")
	err = taskMgr.Transfer(tsk.ID, nodeID, "manager-001", "manager-002", "项目经理出差,转交给副经理审批")
	if err != nil {
		log.Fatalf("Failed to transfer: %v", err)
	}
	tsk, _ = taskMgr.Get(tsk.ID)
	approvers = tsk.Approvers[nodeID]
	fmt.Printf("✓ 转交成功: manager-001 已转交给 manager-002\n")
	fmt.Printf("  当前审批人列表: %v (共 %d 人)\n\n", approvers, len(approvers))

	// 步骤 7: 减签操作
	fmt.Println("步骤 7: 减签操作")
	fmt.Println("  说明: 移除部分审批人(需权限控制)")
	err = taskMgr.RemoveApprover(tsk.ID, nodeID, "tech-lead-001", "减签: 技术负责人无需参与审批")
	if err != nil {
		log.Fatalf("Failed to remove approver: %v", err)
	}
	tsk, _ = taskMgr.Get(tsk.ID)
	approvers = tsk.Approvers[nodeID]
	fmt.Printf("✓ 减签成功: 已移除 tech-lead-001\n")
	fmt.Printf("  当前审批人列表: %v (共 %d 人)\n\n", approvers, len(approvers))

	// 步骤 8: 执行审批操作
	fmt.Println("步骤 8: 执行审批操作")
	fmt.Println("  说明: 当前审批人列表为 [manager-002, finance-001],需要全部同意")
	fmt.Println()

	// manager-002 审批
	fmt.Println("  1. manager-002 审批")
	err = taskMgr.Approve(tsk.ID, nodeID, "manager-002", "副经理审批通过,项目方案可行")
	if err != nil {
		log.Fatalf("Failed to approve: %v", err)
	}
	tsk, _ = taskMgr.Get(tsk.ID)
	fmt.Printf("    ✓ manager-002 已同意 (当前状态: %s)\n\n", tsk.State)

	// finance-001 审批
	fmt.Println("  2. finance-001 审批")
	err = taskMgr.Approve(tsk.ID, nodeID, "finance-001", "财务审批通过,预算合理")
	if err != nil {
		log.Fatalf("Failed to approve: %v", err)
	}
	tsk, _ = taskMgr.Get(tsk.ID)
	fmt.Printf("    ✓ finance-001 已同意 (当前状态: %s)\n\n", tsk.State)

	fmt.Println("✓ 所有审批人已同意,审批流程完成")
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
	nodeID := "project-approval"
	approvers := tsk.Approvers[nodeID]
	if len(approvers) > 0 {
		for _, approver := range approvers {
			fmt.Printf("  - %s\n", approver)
		}
	} else {
		fmt.Println("  无审批人")
	}

	// 输出审批记录(按时间顺序)
	fmt.Println("\n审批记录(按时间顺序):")
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
	if tsk.State == types.TaskStateApproved {
		fmt.Println("✓ 任务已成功通过审批")
	} else {
		fmt.Printf("✗ 任务状态异常: 期望 %s, 实际 %s\n", types.TaskStateApproved, tsk.State)
	}

	// 验证高级操作
	transferCount := 0
	addApproverCount := 0
	removeApproverCount := 0
	approveCount := 0

	for _, record := range tsk.Records {
		switch record.Result {
		case "transfer":
			transferCount++
		case "add_approver":
			addApproverCount++
		case "remove_approver":
			removeApproverCount++
		case "approve":
			approveCount++
		}
	}

	fmt.Println("✓ 高级操作统计:")
	fmt.Printf("  - 转交操作: %d 次\n", transferCount)
	fmt.Printf("  - 加签操作: %d 次\n", addApproverCount)
	fmt.Printf("  - 减签操作: %d 次\n", removeApproverCount)
	fmt.Printf("  - 审批操作: %d 次\n", approveCount)

	// 验证最终审批人列表
	expectedApprovers := []string{"manager-002", "finance-001"}
	if len(approvers) == len(expectedApprovers) {
		fmt.Printf("✓ 最终审批人列表正确: %v\n", approvers)
	} else {
		fmt.Printf("✗ 最终审批人列表异常: 期望 %v, 实际 %v\n", expectedApprovers, approvers)
	}
}

// getNodeName 获取节点名称
func getNodeName(nodeID string) string {
	nodeNames := map[string]string{
		"project-approval": "项目审批",
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

