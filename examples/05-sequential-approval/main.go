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
	fmt.Println("=== 场景 05: 顺序审批流程 ===")
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
	fmt.Println("步骤 2: 创建员工转正任务")
	params := json.RawMessage(`{
		"employeeId": "emp-001",
		"employeeName": "张三",
		"department": "技术部",
		"position": "高级工程师",
		"probationPeriod": 6,
		"performance": "优秀"
	}`)
	tsk, err := taskMgr.Create(tpl.ID, "promotion-001", params)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}
	fmt.Printf("✓ 任务创建成功: ID=%s, State=%s\n\n", tsk.ID, tsk.State)

	// 4. 执行场景流程
	runScenario(taskMgr, tsk)

	// 5. 输出结果
	printResults(taskMgr, tsk)
}

// createTemplate 创建员工转正审批模板
// 流程: 开始节点 → 审批节点(顺序审批模式) → 结束节点
// 使用顺序审批模式(ApprovalModeSequential),多个审批人按顺序依次审批
func createTemplate() *template.Template {
	now := time.Now()
	return &template.Template{
		ID:          "promotion-approval-template",
		Name:        "员工转正审批模板",
		Description: "员工转正需要按顺序经过直属主管、部门经理、HR、总经理审批",
		Version:     1,
		CreatedAt:   now,
		UpdatedAt:   now,
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "开始",
				Type: template.NodeTypeStart,
			},
			"sequential-approval": {
				ID:   "sequential-approval",
				Name: "顺序审批",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeSequential, // 顺序审批模式
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{
							"direct-manager-001", // 直属主管
							"dept-manager-001",   // 部门经理
							"hr-001",             // HR
							"ceo-001",            // 总经理
						},
					},
					RejectBehavior: node.RejectBehaviorTerminate, // 任一审批人拒绝,流程终止
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
			{From: "start", To: "sequential-approval"},
			{From: "sequential-approval", To: "end"},
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
	fmt.Println("  说明: 顺序审批模式需要设置所有审批人,按顺序依次审批")
	fmt.Println()

	nodeID := "sequential-approval"
	approvers := []string{
		"direct-manager-001", // 直属主管
		"dept-manager-001",   // 部门经理
		"hr-001",             // HR
		"ceo-001",            // 总经理
	}

	for _, approver := range approvers {
		err = taskMgr.AddApprover(tsk.ID, nodeID, approver, "设置审批人")
		if err != nil {
			log.Fatalf("Failed to add approver: %v", err)
		}
	}
	fmt.Printf("✓ 审批人列表已设置: %v\n\n", approvers)

	// 步骤 5: 按顺序执行审批操作
	fmt.Println("步骤 5: 按顺序执行审批操作")
	fmt.Println("  说明: 顺序审批模式下,前一个审批人同意后,下一个审批人才能审批")
	fmt.Println()

	// 定义审批顺序和意见
	approvalSteps := []struct {
		approver string
		comment  string
		name     string
	}{
		{"direct-manager-001", "直属主管审批通过,员工表现优秀", "直属主管"},
		{"dept-manager-001", "部门经理审批通过,同意转正", "部门经理"},
		{"hr-001", "HR审批通过,符合转正条件", "HR"},
		{"ceo-001", "总经理审批通过,同意转正", "总经理"},
	}

	for i, step := range approvalSteps {
		fmt.Printf("  %d. %s审批\n", i+1, step.name)

		// 执行审批
		err = taskMgr.Approve(tsk.ID, nodeID, step.approver, step.comment)
		if err != nil {
			log.Fatalf("Failed to approve: %v", err)
		}

		tsk, _ = taskMgr.Get(tsk.ID)
		fmt.Printf("    ✓ %s已同意 (当前状态: %s, 已审批: %d/%d)\n", step.name, tsk.State, i+1, len(approvalSteps))

		// 如果任务已完成,退出循环
		if tsk.State == types.TaskStateApproved {
			break
		}
		fmt.Println()
	}

	fmt.Println()
	fmt.Println("✓ 所有审批人已按顺序完成审批,审批流程完成")
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

	// 输出审批记录(按审批顺序)
	fmt.Println("\n审批记录(按审批顺序):")
	if len(tsk.Records) == 0 {
		fmt.Println("  无审批记录")
	} else {
		// 定义审批人顺序
		approverOrder := []string{
			"direct-manager-001",
			"dept-manager-001",
			"hr-001",
			"ceo-001",
		}

		recordIndex := 1
		for _, approverID := range approverOrder {
			for _, record := range tsk.Records {
				if record.NodeID == "sequential-approval" && record.Approver == approverID && record.Result == "approve" {
					fmt.Printf("  记录 %d:\n", recordIndex)
					fmt.Printf("    审批人: %s (%s)\n", record.Approver, getApproverName(record.Approver))
					fmt.Printf("    审批结果: %s\n", record.Result)
					fmt.Printf("    审批意见: %s\n", record.Comment)
					fmt.Printf("    审批时间: %s\n", record.CreatedAt.Format("2006-01-02 15:04:05"))
					recordIndex++
					break
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

	// 验证审批顺序
	expectedApprovers := []string{
		"direct-manager-001",
		"dept-manager-001",
		"hr-001",
		"ceo-001",
	}
	approveCount := 0
	for _, approverID := range expectedApprovers {
		for _, record := range tsk.Records {
			if record.NodeID == "sequential-approval" && record.Approver == approverID && record.Result == "approve" {
				approveCount++
				break
			}
		}
	}
	if approveCount == len(expectedApprovers) {
		fmt.Printf("✓ 已按顺序完成 %d 个审批人的审批\n", approveCount)
	} else {
		fmt.Printf("✗ 审批人数量异常: 期望 %d 个, 实际 %d 个\n", len(expectedApprovers), approveCount)
	}
}

// getApproverName 获取审批人名称
func getApproverName(approverID string) string {
	approverNames := map[string]string{
		"direct-manager-001": "直属主管",
		"dept-manager-001":    "部门经理",
		"hr-001":              "HR",
		"ceo-001":             "总经理",
	}
	if name, ok := approverNames[approverID]; ok {
		return name
	}
	return approverID
}

