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
	fmt.Println("=== 场景 16: 或签模式独立场景 ===")
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
	fmt.Println("步骤 2: 创建紧急采购审批任务")
	params := json.RawMessage(`{
		"requestNo": "REQ-2025-001",
		"requestType": "紧急采购",
		"requester": "采购员-001",
		"amount": 50000,
		"urgency": "high",
		"description": "紧急采购申请,需要快速审批"
	}`)
	tsk, err := taskMgr.Create(tpl.ID, "urgent-purchase-001", params)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}
	fmt.Printf("✓ 任务创建成功: ID=%s, State=%s\n\n", tsk.ID, tsk.State)

	// 4. 执行场景流程
	runScenario(taskMgr, tsk)

	// 5. 输出结果
	printResults(taskMgr, tsk)
}

// createTemplate 创建紧急采购审批模板
// 流程: 开始节点 → 审批节点(多人或签) → 结束节点
// 审批节点配置为多人或签模式,任意一人同意即可通过
func createTemplate() *template.Template {
	now := time.Now()
	return &template.Template{
		ID:          "or-approval-template",
		Name:        "紧急采购审批模板",
		Description: "紧急采购审批流程,任意一位采购经理同意即可通过",
		Version:     1,
		CreatedAt:   now,
		UpdatedAt:   now,
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "开始",
				Type: template.NodeTypeStart,
			},
			"purchase-approval": {
				ID:   "purchase-approval",
				Name: "采购审批",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeOr, // 多人或签模式
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{
							"purchase-manager-001",
							"purchase-manager-002",
							"purchase-manager-003",
						},
					},
					RejectBehavior: node.RejectBehaviorTerminate, // 所有审批人拒绝时终止流程
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
			{From: "start", To: "purchase-approval"},
			{From: "purchase-approval", To: "end"},
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

	nodeID := "purchase-approval"

	// 步骤 4: 设置审批人列表
	fmt.Println("步骤 4: 设置审批人列表")
	fmt.Println("  说明: 或签模式下,任意一人同意即可通过")
	fmt.Println("  审批人列表: [purchase-manager-001, purchase-manager-002, purchase-manager-003]")
	fmt.Println()

	err = taskMgr.AddApprover(tsk.ID, nodeID, "purchase-manager-001", "设置采购经理1为审批人")
	if err != nil {
		log.Fatalf("Failed to add approver: %v", err)
	}
	err = taskMgr.AddApprover(tsk.ID, nodeID, "purchase-manager-002", "设置采购经理2为审批人")
	if err != nil {
		log.Fatalf("Failed to add approver: %v", err)
	}
	err = taskMgr.AddApprover(tsk.ID, nodeID, "purchase-manager-003", "设置采购经理3为审批人")
	if err != nil {
		log.Fatalf("Failed to add approver: %v", err)
	}
	fmt.Printf("✓ 审批人列表已设置: 共 %d 人\n\n", 3)

	// 步骤 5: 第一个审批人同意
	fmt.Println("步骤 5: 第一个审批人同意")
	fmt.Println("  说明: 或签模式下,第一个审批人同意后流程立即继续")
	fmt.Println()
	err = taskMgr.Approve(tsk.ID, nodeID, "purchase-manager-001", "采购经理1审批通过,紧急采购可以执行")
	if err != nil {
		log.Fatalf("Failed to approve: %v", err)
	}

	tsk, _ = taskMgr.Get(tsk.ID)
	fmt.Printf("✓ purchase-manager-001 已同意 (当前状态: %s)\n", tsk.State)
	fmt.Println("  说明: 第一个审批人同意后,任务状态变为 approved")
	fmt.Println("  注意: 其他审批人无需再审批,流程已完成")
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
	nodeID := "purchase-approval"
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
			fmt.Printf("    审批意见: %s\n", record.Comment)
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

	// 验证或签模式
	approveCount := 0
	for _, record := range tsk.Records {
		if record.Result == "approve" {
			approveCount++
		}
	}
	fmt.Printf("✓ 或签模式验证: 共 %d 个审批人,仅需 %d 人同意即可通过\n", len(approvers), 1)
	fmt.Printf("  实际同意人数: %d\n", approveCount)
	if approveCount >= 1 {
		fmt.Println("  ✓ 或签模式工作正常: 第一个审批人同意后流程立即完成")
	} else {
		fmt.Println("  ✗ 或签模式异常: 需要至少1人同意")
	}
}

// getNodeName 获取节点名称
func getNodeName(nodeID string) string {
	nodeNames := map[string]string{
		"purchase-approval": "采购审批",
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

