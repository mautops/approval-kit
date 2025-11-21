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
	fmt.Println("=== 场景 02: 多人会签审批流程 ===")
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
	fmt.Println("步骤 2: 创建采购申请任务")
	params := json.RawMessage(`{
		"item": "办公设备",
		"amount": 50000,
		"reason": "部门办公设备采购",
		"supplier": "XX供应商"
	}`)
	tsk, err := taskMgr.Create(tpl.ID, "purchase-001", params)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}
	fmt.Printf("✓ 任务创建成功: ID=%s, State=%s\n\n", tsk.ID, tsk.State)

	// 4. 执行场景流程
	runScenario(taskMgr, tsk)

	// 5. 输出结果
	printResults(taskMgr, tsk)
}

// createTemplate 创建采购审批模板
// 流程: 开始节点 → 审批节点(多人会签) → 结束节点
func createTemplate() *template.Template {
	now := time.Now()
	return &template.Template{
		ID:          "purchase-approval-template",
		Name:        "采购审批模板",
		Description: "采购申请需要财务、采购、总经理全部同意才能通过",
		Version:     1,
		CreatedAt:   now,
		UpdatedAt:   now,
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "开始",
				Type: template.NodeTypeStart,
			},
			"unanimous-approval": {
				ID:   "unanimous-approval",
				Name: "会签审批",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeUnanimous, // 多人会签模式
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{
							"finance-001",    // 财务审批人
							"purchase-001",   // 采购审批人
							"manager-001",    // 总经理审批人
						},
					},
					RejectBehavior: node.RejectBehaviorTerminate, // 任一审批人拒绝,流程终止
					Permissions: node.OperationPermissions{
						AllowAddApprover: true, // 允许加签,用于设置初始审批人列表
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
			{From: "start", To: "unanimous-approval"},
			{From: "unanimous-approval", To: "end"},
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
	fmt.Printf("✓ 任务已提交: State=%s, CurrentNode=%s\n\n", tsk.State, tsk.CurrentNode)

	// 步骤 3.5: 设置审批人列表(固定审批人需要在节点激活前设置)
	fmt.Println("步骤 3.5: 设置审批人列表")
	err = taskMgr.AddApprover(tsk.ID, "unanimous-approval", "finance-001", "设置财务审批人")
	if err != nil {
		log.Fatalf("Failed to add finance approver: %v", err)
	}
	err = taskMgr.AddApprover(tsk.ID, "unanimous-approval", "purchase-001", "设置采购审批人")
	if err != nil {
		log.Fatalf("Failed to add purchase approver: %v", err)
	}
	err = taskMgr.AddApprover(tsk.ID, "unanimous-approval", "manager-001", "设置总经理审批人")
	if err != nil {
		log.Fatalf("Failed to add manager approver: %v", err)
	}
	fmt.Println("✓ 审批人列表已设置: finance-001, purchase-001, manager-001")

	// 步骤 4: 执行多人会签审批
	fmt.Println("步骤 4: 执行多人会签审批")
	fmt.Println("  说明: 需要财务、采购、总经理全部同意才能通过")
	fmt.Println()

	// 4.1 财务审批
	fmt.Println("  4.1 财务审批")
	err = taskMgr.Approve(tsk.ID, "unanimous-approval", "finance-001", "财务审核通过,预算充足")
	if err != nil {
		log.Fatalf("Failed to approve by finance: %v", err)
	}
	tsk, _ = taskMgr.Get(tsk.ID)
	fmt.Printf("  ✓ 财务已同意 (当前状态: %s, 已审批: 1/3)\n", tsk.State)

	// 4.2 采购审批
	fmt.Println("  4.2 采购审批")
	err = taskMgr.Approve(tsk.ID, "unanimous-approval", "purchase-001", "采购审核通过,供应商资质合格")
	if err != nil {
		log.Fatalf("Failed to approve by purchase: %v", err)
	}
	tsk, _ = taskMgr.Get(tsk.ID)
	fmt.Printf("  ✓ 采购已同意 (当前状态: %s, 已审批: 2/3)\n", tsk.State)

	// 4.3 总经理审批
	fmt.Println("  4.3 总经理审批")
	err = taskMgr.Approve(tsk.ID, "unanimous-approval", "manager-001", "总经理审批通过")
	if err != nil {
		log.Fatalf("Failed to approve by manager: %v", err)
	}
	tsk, _ = taskMgr.Get(tsk.ID)
	fmt.Printf("  ✓ 总经理已同意 (当前状态: %s, 已审批: 3/3)\n", tsk.State)
	fmt.Println()
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

	// 输出审批记录
	fmt.Println("\n审批记录:")
	if len(tsk.Records) == 0 {
		fmt.Println("  无审批记录")
	} else {
		for i, record := range tsk.Records {
			fmt.Printf("  记录 %d:\n", i+1)
			fmt.Printf("    节点 ID: %s\n", record.NodeID)
			fmt.Printf("    审批人: %s\n", record.Approver)
			fmt.Printf("    审批结果: %s\n", record.Result)
			fmt.Printf("    审批意见: %s\n", record.Comment)
			fmt.Printf("    审批时间: %s\n", record.CreatedAt.Format("2006-01-02 15:04:05"))
		}
	}

	// 输出审批人列表和审批状态
	fmt.Println("\n审批人列表和审批状态:")
	approvers := tsk.Approvers["unanimous-approval"]
	if len(approvers) > 0 {
		for _, approver := range approvers {
			approvals := tsk.Approvals["unanimous-approval"]
			if approvals != nil {
				approval := approvals[approver]
				if approval != nil {
					fmt.Printf("  %s: %s (%s)\n", approver, approval.Result, approval.Comment)
				} else {
					fmt.Printf("  %s: 未审批\n", approver)
				}
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

	if len(tsk.Records) >= 3 {
		fmt.Printf("✓ 已生成 %d 条审批记录(财务、采购、总经理)\n", len(tsk.Records))
	} else {
		fmt.Printf("✗ 审批记录数量不足: 期望至少 3 条, 实际 %d 条\n", len(tsk.Records))
	}

	// 验证所有审批人都已同意
	allApproved := true
	approvals := tsk.Approvals["unanimous-approval"]
	if approvals != nil {
		for _, approver := range approvers {
			approval := approvals[approver]
			if approval == nil || approval.Result != "approve" {
				allApproved = false
				break
			}
		}
	}
	if allApproved {
		fmt.Println("✓ 所有审批人都已同意(会签完成)")
	} else {
		fmt.Println("✗ 部分审批人未同意或已拒绝")
	}
}

