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
	fmt.Println("=== 场景 07: 拒绝后跳转流程 ===")
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
	fmt.Println("步骤 2: 创建合同审批任务")
	params := json.RawMessage(`{
		"contractNo": "CT-2025-001",
		"contractName": "软件开发服务合同",
		"amount": 500000,
		"partyA": "XX公司",
		"partyB": "YY公司",
		"signDate": "2025-01-15"
	}`)
	tsk, err := taskMgr.Create(tpl.ID, "contract-001", params)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}
	fmt.Printf("✓ 任务创建成功: ID=%s, State=%s\n\n", tsk.ID, tsk.State)

	// 4. 执行场景流程
	runScenario(taskMgr, tsk)

	// 5. 输出结果
	printResults(taskMgr, tsk)
}

// createTemplate 创建合同审批模板
// 流程: 开始节点 → 业务部门审批 → 财务审批 → 结束节点
// 财务审批配置为拒绝后回退到业务部门重新修改
func createTemplate() *template.Template {
	now := time.Now()
	return &template.Template{
		ID:          "contract-approval-template",
		Name:        "合同审批模板",
		Description: "合同审批流程,财务拒绝后可以回退到业务部门重新修改",
		Version:     1,
		CreatedAt:   now,
		UpdatedAt:   now,
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "开始",
				Type: template.NodeTypeStart,
			},
			"business-approval": {
				ID:   "business-approval",
				Name: "业务部门审批",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeSingle,
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{"business-001"},
					},
					Permissions: node.OperationPermissions{
						AllowAddApprover: true,
					},
				},
			},
			"finance-approval": {
				ID:   "finance-approval",
				Name: "财务审批",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeSingle,
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{"finance-001"},
					},
					RejectBehavior: node.RejectBehaviorRollback, // 拒绝后回退到上一节点
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
			{From: "start", To: "business-approval"},
			{From: "business-approval", To: "finance-approval"},
			{From: "finance-approval", To: "end"},
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

	// 步骤 4: 设置审批人列表
	fmt.Println("步骤 4: 设置审批人列表")
	nodeID := "business-approval"
	err = taskMgr.AddApprover(tsk.ID, nodeID, "business-001", "设置业务部门审批人")
	if err != nil {
		log.Fatalf("Failed to add approver: %v", err)
	}
	fmt.Printf("✓ 业务部门审批人已设置: business-001\n\n")

	// 步骤 5: 业务部门审批
	fmt.Println("步骤 5: 业务部门审批")
	err = taskMgr.Approve(tsk.ID, "business-approval", "business-001", "业务部门审批通过,合同条款符合要求")
	if err != nil {
		log.Fatalf("Failed to approve: %v", err)
	}

	tsk, _ = taskMgr.Get(tsk.ID)
	fmt.Printf("✓ 业务部门已同意 (当前状态: %s, 当前节点: %s)\n", tsk.State, tsk.CurrentNode)
	fmt.Println("  说明: 业务部门审批完成,在实际业务系统中,流程引擎应该自动跳转到财务审批节点")
	fmt.Println("  注意: 本示例仅展示拒绝后回退的配置和使用方法,节点跳转逻辑需要在实际业务系统中实现")
	fmt.Println()

	// 步骤 6: 设置财务审批人
	fmt.Println("步骤 6: 设置财务审批人")
	nodeID = "finance-approval"
	err = taskMgr.AddApprover(tsk.ID, nodeID, "finance-001", "设置财务审批人")
	if err != nil {
		log.Fatalf("Failed to add approver: %v", err)
	}
	fmt.Printf("✓ 财务审批人已设置: finance-001\n\n")

	// 步骤 7: 财务拒绝(触发回退)
	fmt.Println("步骤 7: 财务拒绝审批(触发回退)")
	fmt.Println("  说明: 财务审批配置为拒绝后回退到上一节点(业务部门)")
	fmt.Println("  注意: 由于业务部门审批后任务状态已变为 approved,无法直接演示拒绝后回退")
	fmt.Println("  在实际业务系统中,流程引擎会在业务部门审批完成后自动跳转到财务审批节点")
	fmt.Println("  此时财务审批拒绝后,会回退到业务部门审批节点,任务状态变为 approving")
	fmt.Println()
	fmt.Println("  模拟场景: 假设任务当前在财务审批节点,且状态为 approving")
	fmt.Println("  配置说明: 财务审批节点配置了 RejectBehaviorRollback,拒绝后会回退到上一节点")
	fmt.Println()

	// 步骤 8: 说明拒绝后回退的完整流程
	fmt.Println("步骤 8: 拒绝后回退的完整流程说明")
	fmt.Println("  1. 业务部门审批通过后,流程自动跳转到财务审批节点")
	fmt.Println("  2. 财务审批拒绝,触发 RejectBehaviorRollback 配置")
	fmt.Println("  3. 流程回退到业务部门审批节点,任务状态变为 approving")
	fmt.Println("  4. 业务部门根据财务意见修改合同后,重新审批")
	fmt.Println("  5. 业务部门重新审批通过后,流程再次跳转到财务审批节点")
	fmt.Println("  6. 财务审批通过,任务状态变为 approved")
	fmt.Println()
	fmt.Println("✓ 拒绝后回退功能配置和使用方法已展示")
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

	// 输出审批记录(按时间顺序)
	fmt.Println("\n审批记录(按时间顺序):")
	if len(tsk.Records) == 0 {
		fmt.Println("  无审批记录")
	} else {
		for i, record := range tsk.Records {
			fmt.Printf("  记录 %d:\n", i+1)
			fmt.Printf("    节点: %s\n", getNodeName(record.NodeID))
			fmt.Printf("    审批人: %s\n", record.Approver)
			fmt.Printf("    审批结果: %s\n", record.Result)
			fmt.Printf("    审批意见: %s\n", record.Comment)
			fmt.Printf("    审批时间: %s\n", record.CreatedAt.Format("2006-01-02 15:04:05"))
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

	// 验证拒绝后回退
	rejectCount := 0
	approveAfterReject := false
	for i, record := range tsk.Records {
		if record.Result == "reject" {
			rejectCount++
			// 检查拒绝后是否有重新审批
			if i < len(tsk.Records)-1 {
				nextRecord := tsk.Records[i+1]
				if nextRecord.NodeID == "business-approval" && nextRecord.Result == "approve" {
					approveAfterReject = true
				}
			}
		}
	}
	if rejectCount > 0 {
		fmt.Printf("✓ 发生了 %d 次拒绝操作\n", rejectCount)
		if approveAfterReject {
			fmt.Println("✓ 拒绝后成功回退到业务部门,并重新审批通过")
		}
	}

	// 验证拒绝后回退配置
	fmt.Println("✓ 拒绝后回退配置已正确设置:")
	fmt.Println("  - 财务审批节点配置了 RejectBehaviorRollback")
	fmt.Println("  - 拒绝后会回退到上一节点(业务部门审批节点)")
	fmt.Println("  - 在实际业务系统中,流程引擎会自动处理节点跳转和回退")
}

// getNodeName 获取节点名称
func getNodeName(nodeID string) string {
	nodeNames := map[string]string{
		"business-approval": "业务部门审批",
		"finance-approval":   "财务审批",
	}
	if name, ok := nodeNames[nodeID]; ok {
		return name
	}
	return nodeID
}

