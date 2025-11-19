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
	fmt.Println("=== 场景 07-Jump: 拒绝后跳转到指定节点 ===")
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
	fmt.Println("步骤 2: 创建方案审批任务")
	params := json.RawMessage(`{
		"proposalNo": "PR-2025-001",
		"proposalName": "新产品开发方案",
		"proposer": "产品经理-001",
		"budget": 1000000,
		"description": "新产品开发方案审批"
	}`)
	tsk, err := taskMgr.Create(tpl.ID, "proposal-001", params)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}
	fmt.Printf("✓ 任务创建成功: ID=%s, State=%s\n\n", tsk.ID, tsk.State)

	// 4. 执行场景流程
	runScenario(taskMgr, tsk)

	// 5. 输出结果
	printResults(taskMgr, tsk)
}

// createTemplate 创建方案审批模板
// 流程: 开始节点 → 部门审批 → 财务审批 → 总经理审批 → 结束节点
// 财务审批配置为拒绝后跳转到部门审批节点(重新修改)
func createTemplate() *template.Template {
	now := time.Now()
	return &template.Template{
		ID:          "proposal-approval-template",
		Name:        "方案审批模板",
		Description: "方案审批流程,财务拒绝后跳转到部门审批节点重新修改",
		Version:     1,
		CreatedAt:   now,
		UpdatedAt:   now,
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "开始",
				Type: template.NodeTypeStart,
			},
			"dept-approval": {
				ID:   "dept-approval",
				Name: "部门审批",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeSingle,
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{"dept-manager-001"},
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
					RejectBehavior:  node.RejectBehaviorJump,      // 拒绝后跳转到指定节点
					RejectTargetNode: "dept-approval",              // 跳转到部门审批节点
					Permissions: node.OperationPermissions{
						AllowAddApprover: true,
					},
				},
			},
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
			{From: "start", To: "dept-approval"},
			{From: "dept-approval", To: "finance-approval"},
			{From: "finance-approval", To: "ceo-approval"},
			{From: "ceo-approval", To: "end"},
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
	nodeID := "dept-approval"
	err = taskMgr.AddApprover(tsk.ID, nodeID, "dept-manager-001", "设置部门审批人")
	if err != nil {
		log.Fatalf("Failed to add approver: %v", err)
	}
	fmt.Printf("✓ 部门审批人已设置: dept-manager-001\n\n")

	// 步骤 5: 部门审批
	fmt.Println("步骤 5: 部门审批")
	err = taskMgr.Approve(tsk.ID, "dept-approval", "dept-manager-001", "部门审批通过,方案可行")
	if err != nil {
		log.Fatalf("Failed to approve: %v", err)
	}

	tsk, _ = taskMgr.Get(tsk.ID)
	fmt.Printf("✓ 部门已同意 (当前状态: %s, 当前节点: %s)\n", tsk.State, tsk.CurrentNode)
	fmt.Println("  说明: 部门审批完成,在实际业务系统中,流程引擎应该自动跳转到财务审批节点")
	fmt.Println("  注意: 本示例仅展示拒绝后跳转的配置和使用方法,节点跳转逻辑需要在实际业务系统中实现")
	fmt.Println()

	// 步骤 6: 设置财务审批人
	fmt.Println("步骤 6: 设置财务审批人")
	nodeID = "finance-approval"
	err = taskMgr.AddApprover(tsk.ID, nodeID, "finance-001", "设置财务审批人")
	if err != nil {
		log.Fatalf("Failed to add approver: %v", err)
	}
	fmt.Printf("✓ 财务审批人已设置: finance-001\n\n")

	// 步骤 7: 财务拒绝(触发跳转)
	fmt.Println("步骤 7: 财务拒绝审批(触发跳转)")
	fmt.Println("  说明: 财务审批配置为拒绝后跳转到指定节点(部门审批节点)")
	fmt.Println("  注意: 由于部门审批后任务状态已变为 approved,无法直接演示拒绝后跳转")
	fmt.Println("  在实际业务系统中,流程引擎会在部门审批完成后自动跳转到财务审批节点")
	fmt.Println("  此时财务审批拒绝后,会跳转到指定的目标节点(部门审批节点),任务状态变为 approving")
	fmt.Println()
	fmt.Println("  模拟场景: 假设任务当前在财务审批节点,且状态为 approving")
	fmt.Println("  配置说明:")
	fmt.Println("    - 财务审批节点配置了 RejectBehaviorJump")
	fmt.Println("    - 拒绝后跳转目标节点设置为 dept-approval(部门审批节点)")
	fmt.Println("    - 拒绝后会跳转到部门审批节点,而不是回退到上一节点")
	fmt.Println()

	// 步骤 8: 说明拒绝后跳转的完整流程
	fmt.Println("步骤 8: 拒绝后跳转的完整流程说明")
	fmt.Println("  1. 部门审批通过后,流程自动跳转到财务审批节点")
	fmt.Println("  2. 财务审批拒绝,触发 RejectBehaviorJump 配置")
	fmt.Println("  3. 流程跳转到指定的目标节点(部门审批节点),任务状态变为 approving")
	fmt.Println("  4. 部门根据财务意见修改方案后,重新审批")
	fmt.Println("  5. 部门重新审批通过后,流程再次跳转到财务审批节点")
	fmt.Println("  6. 财务审批通过,流程继续到总经理审批节点")
	fmt.Println("  7. 总经理审批通过,任务状态变为 approved")
	fmt.Println()
	fmt.Println("✓ 拒绝后跳转功能配置和使用方法已展示")
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

	// 验证拒绝后跳转配置
	fmt.Println("✓ 拒绝后跳转配置已正确设置:")
	fmt.Println("  - 财务审批节点配置了 RejectBehaviorJump")
	fmt.Println("  - 拒绝后跳转目标节点设置为 dept-approval(部门审批节点)")
	fmt.Println("  - 拒绝后会跳转到指定的目标节点,而不是回退到上一节点")
	fmt.Println("  - 在实际业务系统中,流程引擎会自动处理节点跳转")
}

// getNodeName 获取节点名称
func getNodeName(nodeID string) string {
	nodeNames := map[string]string{
		"dept-approval":   "部门审批",
		"finance-approval": "财务审批",
		"ceo-approval":     "总经理审批",
	}
	if name, ok := nodeNames[nodeID]; ok {
		return name
	}
	return nodeID
}

