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
	fmt.Println("=== 场景 12: 多条件组合判断场景 ===")
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

	// 3. 演示不同条件组合
	demonstrateCompositeConditions(templateMgr, taskMgr)
}

// createTemplate 创建审批模板
// 流程: 开始节点 → 条件节点(组合条件) → [审批路径 A/审批路径 B] → 结束节点
// 条件节点使用组合条件,根据金额和部门组合判断
func createTemplate() *template.Template {
	now := time.Now()
	return &template.Template{
		ID:          "composite-condition-template",
		Name:        "组合条件审批模板",
		Description: "根据金额和部门组合判断审批路径",
		Version:     1,
		CreatedAt:   now,
		UpdatedAt:   now,
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "开始",
				Type: template.NodeTypeStart,
			},
			"condition-node": {
				ID:   "condition-node",
				Name: "组合条件判断",
				Type: template.NodeTypeCondition,
				Config: &node.ConditionNodeConfig{
					Condition: &node.Condition{
						Type: "composite",
						Config: &node.CompositeConditionConfig{
							Operator: "and", // AND 逻辑: 所有条件都必须满足
							Conditions: []*node.Condition{
								{
									Type: "numeric",
									Config: &node.NumericConditionConfig{
										Field:    "amount",
										Operator: ">=",
										Value:    50000,
										Source:   "task_params",
									},
								},
								{
									Type: "string",
									Config: &node.StringConditionConfig{
										Field:    "department",
										Operator: "eq",
										Value:    "技术部",
										Source:   "task_params",
									},
								},
							},
						},
					},
					TrueNodeID:  "high-approval",  // 金额>=50000 且 部门=技术部,走高级审批路径
					FalseNodeID: "normal-approval", // 其他情况,走普通审批路径
				},
			},
			"high-approval": {
				ID:   "high-approval",
				Name: "高级审批",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeUnanimous, // 多人会签模式
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{
							"cto-001",
							"finance-001",
							"ceo-001",
						},
					},
					Permissions: node.OperationPermissions{
						AllowAddApprover: true,
					},
				},
			},
			"normal-approval": {
				ID:   "normal-approval",
				Name: "普通审批",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeSingle, // 单人审批模式
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
			{From: "start", To: "condition-node"},
			{From: "condition-node", To: "high-approval"},
			{From: "condition-node", To: "normal-approval"},
			{From: "high-approval", To: "end"},
			{From: "normal-approval", To: "end"},
		},
	}
}

// demonstrateCompositeConditions 演示不同条件组合
func demonstrateCompositeConditions(templateMgr template.TemplateManager, taskMgr task.TaskManager) {
	// 场景 1: 金额>=50000 且 部门=技术部 (走高级审批路径)
	fmt.Println("=== 场景 1: 金额>=50000 且 部门=技术部 ===")
	fmt.Println()
	params1 := json.RawMessage(`{
		"requestNo": "REQ-2025-001",
		"requestType": "组合条件测试",
		"requester": "申请人-001",
		"amount": 100000,
		"department": "技术部",
		"description": "金额>=50000 且 部门=技术部,应该走高级审批路径"
	}`)
	tsk1, err := taskMgr.Create("composite-condition-template", "composite-001", params1)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}
	fmt.Printf("✓ 任务创建成功: ID=%s\n", tsk1.ID)
	fmt.Println("  任务参数:")
	var params1Map map[string]interface{}
	json.Unmarshal(params1, &params1Map)
	for k, v := range params1Map {
		fmt.Printf("    %s: %v\n", k, v)
	}
	fmt.Println("  条件判断: 金额>=50000 AND 部门=技术部")
	fmt.Println("  预期路径: 高级审批路径(多人会签)")
	fmt.Println()

	// 场景 2: 金额>=50000 但 部门!=技术部 (走普通审批路径)
	fmt.Println("=== 场景 2: 金额>=50000 但 部门!=技术部 ===")
	fmt.Println()
	params2 := json.RawMessage(`{
		"requestNo": "REQ-2025-002",
		"requestType": "组合条件测试",
		"requester": "申请人-001",
		"amount": 100000,
		"department": "市场部",
		"description": "金额>=50000 但 部门!=技术部,应该走普通审批路径"
	}`)
	tsk2, err := taskMgr.Create("composite-condition-template", "composite-002", params2)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}
	fmt.Printf("✓ 任务创建成功: ID=%s\n", tsk2.ID)
	fmt.Println("  任务参数:")
	var params2Map map[string]interface{}
	json.Unmarshal(params2, &params2Map)
	for k, v := range params2Map {
		fmt.Printf("    %s: %v\n", k, v)
	}
	fmt.Println("  条件判断: 金额>=50000 AND 部门=技术部")
	fmt.Println("  实际结果: 金额>=50000 但 部门!=技术部,条件不满足")
	fmt.Println("  预期路径: 普通审批路径(单人审批)")
	fmt.Println()

	// 场景 3: 金额<50000 且 部门=技术部 (走普通审批路径)
	fmt.Println("=== 场景 3: 金额<50000 且 部门=技术部 ===")
	fmt.Println()
	params3 := json.RawMessage(`{
		"requestNo": "REQ-2025-003",
		"requestType": "组合条件测试",
		"requester": "申请人-001",
		"amount": 30000,
		"department": "技术部",
		"description": "金额<50000 且 部门=技术部,应该走普通审批路径"
	}`)
	tsk3, err := taskMgr.Create("composite-condition-template", "composite-003", params3)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}
	fmt.Printf("✓ 任务创建成功: ID=%s\n", tsk3.ID)
	fmt.Println("  任务参数:")
	var params3Map map[string]interface{}
	json.Unmarshal(params3, &params3Map)
	for k, v := range params3Map {
		fmt.Printf("    %s: %v\n", k, v)
	}
	fmt.Println("  条件判断: 金额>=50000 AND 部门=技术部")
	fmt.Println("  实际结果: 金额<50000 但 部门=技术部,条件不满足")
	fmt.Println("  预期路径: 普通审批路径(单人审批)")
	fmt.Println()

	// 总结
	fmt.Println("=== 组合条件说明 ===")
	fmt.Println("  组合条件配置:")
	fmt.Println("    - 操作符: AND (所有条件都必须满足)")
	fmt.Println("    - 条件1: 金额 >= 50000")
	fmt.Println("    - 条件2: 部门 = 技术部")
	fmt.Println()
	fmt.Println("  审批路径:")
	fmt.Println("    - True (所有条件满足): 高级审批路径(多人会签)")
	fmt.Println("    - False (任一条件不满足): 普通审批路径(单人审批)")
	fmt.Println()
	fmt.Println("  注意: 本示例仅展示组合条件的配置方法")
	fmt.Println("  实际执行需要在实际业务系统中实现条件评估和流程跳转")
	fmt.Println()
}

