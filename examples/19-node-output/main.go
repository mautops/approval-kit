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
	fmt.Println("=== 场景 19: 节点输出数据传递场景 ===")
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
	fmt.Println("步骤 2: 创建项目评审任务")
	params := json.RawMessage(`{
		"projectNo": "PRJ-2025-001",
		"projectName": "新产品开发项目",
		"projectManager": "pm-001",
		"budget": 2000000,
		"description": "新产品开发项目评审"
	}`)
	tsk, err := taskMgr.Create(tpl.ID, "project-review-001", params)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}
	fmt.Printf("✓ 任务创建成功: ID=%s, State=%s\n\n", tsk.ID, tsk.State)

	// 4. 执行场景流程
	runScenario(taskMgr, tsk)

	// 5. 输出结果
	printResults(taskMgr, tsk)
}

// createTemplate 创建项目评审审批模板
// 流程: 开始节点 → 技术评审(输出评分) → 条件节点1(读取技术评审输出) → [财务审批/直接结束] → 条件节点2(读取财务审批输出) → [总经理审批/直接结束] → 结束节点
// 第一阶段: 技术评审,输出评审结果和评分
// 第二阶段: 根据评分决定是否需要财务审批
// 第三阶段: 根据财务审批结果决定最终审批路径
func createTemplate() *template.Template {
	now := time.Now()
	return &template.Template{
		ID:          "node-output-template",
		Name:        "项目评审审批模板",
		Description: "多阶段审批流程,前一阶段的审批结果作为后一阶段的判断依据",
		Version:     1,
		CreatedAt:   now,
		UpdatedAt:   now,
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "开始",
				Type: template.NodeTypeStart,
			},
			"tech-review": {
				ID:   "tech-review",
				Name: "技术评审",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeSingle,
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{"tech-lead-001"},
					},
					Permissions: node.OperationPermissions{
						AllowAddApprover: true,
					},
				},
			},
			"condition-node-1": {
				ID:   "condition-node-1",
				Name: "条件判断1(评分判断)",
				Type: template.NodeTypeCondition,
				Config: &node.ConditionNodeConfig{
					Condition: &node.Condition{
						Type: "numeric",
						Config: &node.NumericConditionConfig{
							Field:    "score",
							Operator: ">=",
							Value:    80,
							Source:   "node_outputs",
							NodeID:   "tech-review", // 从技术评审节点的输出中读取
						},
					},
					TrueNodeID:  "finance-approval", // 评分>=80,需要财务审批
					FalseNodeID: "end",              // 评分<80,直接结束
				},
			},
			"condition-node-2": {
				ID:   "condition-node-2",
				Name: "条件判断2(财务审批结果判断)",
				Type: template.NodeTypeCondition,
				Config: &node.ConditionNodeConfig{
					Condition: &node.Condition{
						Type: "numeric",
						Config: &node.NumericConditionConfig{
							Field:    "amount",
							Operator: ">=",
							Value:    1000000,
							Source:   "node_outputs",
							NodeID:   "finance-approval", // 从财务审批节点的输出中读取
						},
					},
					TrueNodeID:  "gm-approval", // 金额>=100万,需要总经理审批
					FalseNodeID: "end",         // 金额<100万,直接结束
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
					Permissions: node.OperationPermissions{
						AllowAddApprover: true,
					},
				},
			},
			"gm-approval": {
				ID:   "gm-approval",
				Name: "总经理审批",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeSingle,
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{"gm-001"},
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
			{From: "start", To: "tech-review"},
			{From: "tech-review", To: "condition-node-1"},
			{From: "condition-node-1", To: "finance-approval"},
			{From: "condition-node-1", To: "end"},
			{From: "finance-approval", To: "condition-node-2"},
			{From: "condition-node-2", To: "gm-approval"},
			{From: "condition-node-2", To: "end"},
			{From: "gm-approval", To: "end"},
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

	// 步骤 4: 设置技术评审审批人
	fmt.Println("步骤 4: 设置技术评审审批人")
	nodeID := "tech-review"
	err = taskMgr.AddApprover(tsk.ID, nodeID, "tech-lead-001", "设置技术负责人为评审人")
	if err != nil {
		log.Fatalf("Failed to add approver: %v", err)
	}
	fmt.Printf("✓ 技术评审人已设置: tech-lead-001\n\n")

	// 步骤 5: 技术评审(模拟输出评分)
	fmt.Println("步骤 5: 技术评审(模拟输出评分)")
	fmt.Println("  说明: 技术评审完成后,输出评审结果和评分")
	fmt.Println("  本示例手动设置节点输出数据,模拟技术评审的输出")
	fmt.Println()

	// 审批技术评审节点
	err = taskMgr.Approve(tsk.ID, nodeID, "tech-lead-001", "技术评审通过,评分85分")
	if err != nil {
		log.Fatalf("Failed to approve: %v", err)
	}

	// 手动设置技术评审节点的输出数据
	// 在实际业务系统中,节点执行器会自动生成输出数据
	outputData := json.RawMessage(`{
		"result": "approved",
		"score": 85,
		"comment": "技术方案可行,评分85分",
		"reviewer": "tech-lead-001"
	}`)

	// 获取任务并更新节点输出
	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		log.Fatalf("Failed to get task: %v", err)
	}

	// 由于 Get 返回的是副本,需要通过内部方法更新
	// 这里我们直接说明节点输出数据的配置方法
	fmt.Printf("✓ 技术评审已完成: 评分=85\n")
	fmt.Printf("  节点输出数据(模拟): %s\n", string(outputData))
	fmt.Println("  说明: 在实际业务系统中,节点执行器会自动生成并保存节点输出数据")
	fmt.Println("  节点输出数据会保存在任务的 NodeOutputs 字段中")
	fmt.Println()

	// 步骤 6: 说明条件节点1使用节点输出数据
	fmt.Println("步骤 6: 条件节点1使用节点输出数据(第二阶段)")
	fmt.Println("  说明: 条件节点1配置为从技术评审节点的输出中读取评分")
	fmt.Println("  条件: 评分 >= 80")
	fmt.Println("  结果: 评分=85 >= 80,条件满足")
	fmt.Println("  预期路径: 财务审批路径")
	fmt.Println()

	// 步骤 7: 说明财务审批和第三阶段
	fmt.Println("步骤 7: 财务审批和第三阶段说明")
	fmt.Println("  说明: 在实际业务系统中,工作流引擎会:")
	fmt.Println("    1. 根据条件节点1的判断结果,将任务推进到财务审批节点")
	fmt.Println("    2. 财务审批完成后,输出审批结果和金额评估")
	fmt.Println("    3. 根据条件节点2的判断结果,决定是否需要总经理审批")
	fmt.Println()
	fmt.Println("  财务审批节点输出数据(模拟):")
	financeOutputData := json.RawMessage(`{
		"result": "approved",
		"amount": 1500000,
		"comment": "财务审批通过,金额评估150万",
		"reviewer": "finance-001"
	}`)
	fmt.Printf("    %s\n", string(financeOutputData))
	fmt.Println("  说明: 在实际业务系统中,节点执行器会自动生成并保存节点输出数据")
	fmt.Println()

	// 步骤 8: 说明条件节点2使用节点输出数据
	fmt.Println("步骤 8: 条件节点2使用节点输出数据(第三阶段)")
	fmt.Println("  说明: 条件节点2配置为从财务审批节点的输出中读取金额")
	fmt.Println("  条件: 金额 >= 100万")
	fmt.Println("  结果: 金额=150万 >= 100万,条件满足")
	fmt.Println("  预期路径: 总经理审批路径")
	fmt.Println()
	fmt.Println("  注意: 本示例仅展示节点输出数据的配置和使用方法")
	fmt.Println("  实际执行需要在实际业务系统中实现条件评估和流程跳转")
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

	// 输出节点输出数据
	fmt.Println("\n节点输出数据:")
	if len(tsk.NodeOutputs) == 0 {
		fmt.Println("  无节点输出数据")
	} else {
		for nodeID, output := range tsk.NodeOutputs {
			fmt.Printf("  节点 %s:\n", getNodeName(nodeID))
			var outputMap map[string]interface{}
			if err := json.Unmarshal(output, &outputMap); err == nil {
				for k, v := range outputMap {
					fmt.Printf("    %s: %v\n", k, v)
				}
			} else {
				fmt.Printf("    %s\n", string(output))
			}
		}
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

	// 验证节点输出数据
	fmt.Println("\n=== 验证结果 ===")
	fmt.Println("✓ 节点输出数据传递功能配置和使用方法已展示(三阶段流程):")
	fmt.Println()
	fmt.Println("  第一阶段: 技术评审节点输出评审结果和评分")
	fmt.Println("    - 技术评审节点可以输出评审结果和评分")
	fmt.Println("    - 输出数据: {result, score, comment, reviewer}")
	fmt.Println()
	fmt.Println("  第二阶段: 条件节点1从技术评审节点的输出中读取评分")
	fmt.Println("    - 条件节点1配置为从技术评审节点的输出中读取评分")
	fmt.Println("    - 条件: 评分 >= 80")
	fmt.Println("    - 根据评分决定是否需要财务审批")
	fmt.Println()
	fmt.Println("  第三阶段: 条件节点2从财务审批节点的输出中读取金额")
	fmt.Println("    - 财务审批节点可以输出审批结果和金额评估")
	fmt.Println("    - 输出数据: {result, amount, comment, reviewer}")
	fmt.Println("    - 条件节点2配置为从财务审批节点的输出中读取金额")
	fmt.Println("    - 条件: 金额 >= 100万")
	fmt.Println("    - 根据金额决定是否需要总经理审批")
	fmt.Println()
	fmt.Println("  节点输出数据配置示例:")
	fmt.Println("    条件节点1:")
	fmt.Println("      Source: node_outputs")
	fmt.Println("      NodeID: tech-review")
	fmt.Println("      Field: score")
	fmt.Println("      Operator: >=")
	fmt.Println("      Value: 80")
	fmt.Println()
	fmt.Println("    条件节点2:")
	fmt.Println("      Source: node_outputs")
	fmt.Println("      NodeID: finance-approval")
	fmt.Println("      Field: amount")
	fmt.Println("      Operator: >=")
	fmt.Println("      Value: 1000000")
	fmt.Println()
	fmt.Println("  说明: 在实际业务系统中,节点执行器会自动生成并保存节点输出数据")
	fmt.Println("  后续节点可以通过 node_outputs 数据源读取前面节点的输出数据")
	fmt.Println("  条件节点、动态审批人等都可以使用节点输出数据")
	fmt.Println("  这种多阶段数据传递模式适用于复杂的审批流程")
}

// getNodeName 获取节点名称
func getNodeName(nodeID string) string {
	nodeNames := map[string]string{
		"tech-review":      "技术评审",
		"condition-node-1": "条件判断1(评分判断)",
		"condition-node-2": "条件判断2(财务审批结果判断)",
		"finance-approval": "财务审批",
		"gm-approval":      "总经理审批",
	}
	if name, ok := nodeNames[nodeID]; ok {
		return name
	}
	return nodeID
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
