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
	fmt.Println("=== 场景 03: 条件分支审批流程 ===")
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

	// 3. 创建任务 - 测试中额报销场景
	amount := 3000.0
	params := json.RawMessage(fmt.Sprintf(`{
		"amount": %.0f,
		"category": "差旅费",
		"description": "出差费用报销"
	}`, amount))

	tsk, err := taskMgr.Create(tpl.ID, "expense-001", params)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}
	fmt.Printf("✓ 任务创建成功: ID=%s, State=%s\n\n", tsk.ID, tsk.State)

	// 4. 执行场景流程
	runScenario(taskMgr, tsk, amount)

	// 5. 输出结果
	printResults(taskMgr, tsk, amount)
}

// createTemplate 创建费用报销审批模板
// 流程: 开始节点 → 条件节点1(金额<1000?) → [部门经理审批/条件节点2] → [结束/条件节点2] → ...
// 根据金额大小选择不同的审批路径:
// - 金额 < 1000: 部门经理审批 → 结束
// - 金额 >= 1000 且 < 5000: 部门经理审批 → 财务审批 → 结束
// - 金额 >= 5000: 部门经理审批 → 财务审批 → 总经理审批 → 结束
func createTemplate() *template.Template {
	now := time.Now()
	return &template.Template{
		ID:          "expense-approval-template",
		Name:        "费用报销审批模板",
		Description: "根据金额大小选择不同的审批路径",
		Version:     1,
		CreatedAt:   now,
		UpdatedAt:   now,
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "开始",
				Type: template.NodeTypeStart,
			},
			// 条件节点1: 判断金额是否 < 1000
			"condition-amount-1000": {
				ID:   "condition-amount-1000",
				Name: "金额判断(< 1000)",
				Type: template.NodeTypeCondition,
				Config: &node.ConditionNodeConfig{
					Condition: &node.Condition{
						Type: "numeric",
						Config: &node.NumericConditionConfig{
							Field:    "amount",
							Operator: "lt", // 小于
							Value:    1000.0,
							Source:   "task_params",
						},
					},
					TrueNodeID:  "manager-approval-low",  // 金额 < 1000: 部门经理审批(低额)
					FalseNodeID: "condition-amount-5000", // 金额 >= 1000: 继续判断
				},
			},
			// 条件节点2: 判断金额是否 < 5000
			"condition-amount-5000": {
				ID:   "condition-amount-5000",
				Name: "金额判断(< 5000)",
				Type: template.NodeTypeCondition,
				Config: &node.ConditionNodeConfig{
					Condition: &node.Condition{
						Type: "numeric",
						Config: &node.NumericConditionConfig{
							Field:    "amount",
							Operator: "lt", // 小于
							Value:    5000.0,
							Source:   "task_params",
						},
					},
					TrueNodeID:  "manager-approval-medium", // 金额 < 5000: 部门经理审批(中额)
					FalseNodeID: "manager-approval-high",   // 金额 >= 5000: 部门经理审批(高额)
				},
			},
			// 部门经理审批(低额, < 1000)
			"manager-approval-low": {
				ID:   "manager-approval-low",
				Name: "部门经理审批(低额)",
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
			// 部门经理审批(中额, 1000-5000)
			"manager-approval-medium": {
				ID:   "manager-approval-medium",
				Name: "部门经理审批(中额)",
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
			// 部门经理审批(高额, >= 5000)
			"manager-approval-high": {
				ID:   "manager-approval-high",
				Name: "部门经理审批(高额)",
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
			// 财务审批(中额和高额都需要)
			"finance-approval": {
				ID:   "finance-approval",
				Name: "财务审批",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeSingle,
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{"finance-001"},
					},
					Permissions: node.OperationPermissions{
						AllowAddApprover: true,
					},
				},
			},
			// 总经理审批(仅高额需要)
			"ceo-approval": {
				ID:   "ceo-approval",
				Name: "总经理审批",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeSingle,
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{"ceo-001"},
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
			{From: "start", To: "condition-amount-1000"},
			// 条件节点1的分支
			{From: "condition-amount-1000", To: "manager-approval-low"},      // True: 低额路径
			{From: "condition-amount-1000", To: "condition-amount-5000"},    // False: 继续判断
			// 条件节点2的分支
			{From: "condition-amount-5000", To: "manager-approval-medium"},   // True: 中额路径
			{From: "condition-amount-5000", To: "manager-approval-high"},     // False: 高额路径
			// 低额路径: 部门经理审批 → 结束
			{From: "manager-approval-low", To: "end"},
			// 中额路径: 部门经理审批 → 财务审批 → 结束
			{From: "manager-approval-medium", To: "finance-approval"},
			{From: "finance-approval", To: "end"},
			// 高额路径: 部门经理审批 → 财务审批 → 总经理审批 → 结束
			{From: "manager-approval-high", To: "finance-approval"},
			{From: "finance-approval", To: "ceo-approval"},
			{From: "ceo-approval", To: "end"},
		},
	}
}

// runScenario 执行场景流程
func runScenario(taskMgr task.TaskManager, tsk *task.Task, amount float64) {
	// 步骤 2: 提交任务
	fmt.Println("步骤 2: 提交任务进入审批流程")
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

	// 步骤 3: 根据金额判断审批路径
	fmt.Println("步骤 3: 根据金额判断审批路径")
	fmt.Printf("  金额: %.0f\n", amount)
	
	var approvalPath []struct {
		nodeID   string
		approver string
		comment  string
	}

	if amount < 1000 {
		// 低额路径: 部门经理审批
		fmt.Printf("  金额 %.0f < 1000, 走低额路径: 部门经理审批 → 结束\n", amount)
		approvalPath = []struct {
			nodeID   string
			approver string
			comment  string
		}{
			{"manager-approval-low", "manager-001", "部门经理审批通过"},
		}
	} else if amount < 5000 {
		// 中额路径: 部门经理 → 财务审批
		fmt.Printf("  金额 %.0f >= 1000 且 < 5000, 走中额路径: 部门经理审批 → 财务审批 → 结束\n", amount)
		approvalPath = []struct {
			nodeID   string
			approver string
			comment  string
		}{
			{"manager-approval-medium", "manager-001", "部门经理审批通过"},
			{"finance-approval", "finance-001", "财务审批通过"},
		}
	} else {
		// 高额路径: 部门经理 → 财务 → 总经理审批
		fmt.Printf("  金额 %.0f >= 5000, 走高额路径: 部门经理审批 → 财务审批 → 总经理审批 → 结束\n", amount)
		approvalPath = []struct {
			nodeID   string
			approver string
			comment  string
		}{
			{"manager-approval-high", "manager-001", "部门经理审批通过"},
			{"finance-approval", "finance-001", "财务审批通过"},
			{"ceo-approval", "ceo-001", "总经理审批通过"},
		}
	}
	fmt.Println()

	// 步骤 4: 执行审批操作
	fmt.Println("步骤 4: 执行审批操作")
	fmt.Println("  说明: 条件节点会根据金额自动判断路径,这里手动执行对应路径的审批")
	fmt.Println()

	for i, step := range approvalPath {
		// 设置审批人(如果需要)
		tsk, _ = taskMgr.Get(tsk.ID)
		if tsk.Approvers[step.nodeID] == nil || len(tsk.Approvers[step.nodeID]) == 0 {
			err = taskMgr.AddApprover(tsk.ID, step.nodeID, step.approver, "设置审批人")
			if err != nil {
				log.Fatalf("Failed to add approver: %v", err)
			}
		}

		// 执行审批
		err = taskMgr.Approve(tsk.ID, step.nodeID, step.approver, step.comment)
		if err != nil {
			log.Fatalf("Failed to approve: %v", err)
		}

		tsk, _ = taskMgr.Get(tsk.ID)
		fmt.Printf("  %d. %s 已同意 (节点: %s, 状态: %s)\n", i+1, step.approver, step.nodeID, tsk.State)

		// 如果任务已完成,退出循环
		if tsk.State == types.TaskStateApproved {
			break
		}
	}
	fmt.Println()
	fmt.Println("✓ 审批流程完成")
	fmt.Println()
}

// printResults 输出结果
func printResults(taskMgr task.TaskManager, tsk *task.Task, amount float64) {
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

	// 输出审批记录
	fmt.Println("\n审批记录:")
	if len(tsk.Records) == 0 {
		fmt.Println("  无审批记录")
	} else {
		approveCount := 0
		for _, record := range tsk.Records {
			if record.Result == "approve" {
				approveCount++
				fmt.Printf("  记录 %d (审批):\n", approveCount)
				fmt.Printf("    节点 ID: %s\n", record.NodeID)
				fmt.Printf("    审批人: %s\n", record.Approver)
				fmt.Printf("    审批结果: %s\n", record.Result)
				fmt.Printf("    审批意见: %s\n", record.Comment)
				fmt.Printf("    审批时间: %s\n", record.CreatedAt.Format("2006-01-02 15:04:05"))
			}
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

	// 验证审批路径
	expectedApprovals := 0
	if amount < 1000 {
		expectedApprovals = 1
		fmt.Printf("✓ 低额路径: 1 个审批节点(部门经理)\n")
	} else if amount < 5000 {
		expectedApprovals = 2
		fmt.Printf("✓ 中额路径: 2 个审批节点(部门经理、财务)\n")
	} else {
		expectedApprovals = 3
		fmt.Printf("✓ 高额路径: 3 个审批节点(部门经理、财务、总经理)\n")
	}

	approveCount := 0
	for _, record := range tsk.Records {
		if record.Result == "approve" {
			approveCount++
		}
	}
	if approveCount == expectedApprovals {
		fmt.Printf("✓ 审批记录数量正确: %d 条\n", approveCount)
	} else {
		fmt.Printf("✗ 审批记录数量异常: 期望 %d 条, 实际 %d 条\n", expectedApprovals, approveCount)
	}
}

