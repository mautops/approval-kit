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
	fmt.Println("=== 场景 17: 字符串匹配条件场景 ===")
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

	// 3. 演示不同字符串匹配场景
	demonstrateStringConditions(templateMgr, taskMgr, tpl)
}

// createTemplate 创建合同审批模板
// 流程: 开始节点 → 条件节点(字符串匹配) → [采购审批/法务审批/标准审批] → 结束节点
// 条件节点根据合同类型进行字符串匹配
func createTemplate() *template.Template {
	now := time.Now()
	return &template.Template{
		ID:          "string-condition-template",
		Name:        "合同审批模板",
		Description: "根据合同类型(字符串匹配)选择不同的审批路径",
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
				Name: "条件判断",
				Type: template.NodeTypeCondition,
				Config: &node.ConditionNodeConfig{
					Condition: &node.Condition{
						Type: "composite",
						Config: &node.CompositeConditionConfig{
							Operator: "or", // OR 逻辑: 任意一个条件满足即可
							Conditions: []*node.Condition{
								{
									Type: "string",
									Config: &node.StringConditionConfig{
										Field:    "contractType",
										Operator: "contains",
										Value:    "采购",
										Source:   "task_params",
									},
								},
								{
									Type: "string",
									Config: &node.StringConditionConfig{
										Field:    "contractType",
										Operator: "starts_with",
										Value:    "服务",
										Source:   "task_params",
									},
								},
							},
						},
					},
					TrueNodeID:  "purchase-approval",  // 包含"采购"或以"服务"开头,走采购审批路径
					FalseNodeID: "standard-approval", // 其他情况,走标准审批路径
				},
			},
			"purchase-approval": {
				ID:   "purchase-approval",
				Name: "采购审批",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeSingle,
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{"purchase-manager-001"},
					},
					Permissions: node.OperationPermissions{
						AllowAddApprover: true,
					},
				},
			},
			"legal-approval": {
				ID:   "legal-approval",
				Name: "法务审批",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeSingle,
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{"legal-manager-001"},
					},
					Permissions: node.OperationPermissions{
						AllowAddApprover: true,
					},
				},
			},
			"standard-approval": {
				ID:   "standard-approval",
				Name: "标准审批",
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
			{From: "start", To: "condition-node"},
			{From: "condition-node", To: "purchase-approval"},
			{From: "condition-node", To: "standard-approval"},
			{From: "purchase-approval", To: "end"},
			{From: "legal-approval", To: "end"},
			{From: "standard-approval", To: "end"},
		},
	}
}

// demonstrateStringConditions 演示不同字符串匹配场景
func demonstrateStringConditions(templateMgr template.TemplateManager, taskMgr task.TaskManager, tpl *template.Template) {
	// 场景1: 合同类型包含"采购"
	fmt.Println("=== 场景 1: 合同类型包含\"采购\" ===")
	fmt.Println()
	params1 := json.RawMessage(`{
		"contractNo": "CT-2025-001",
		"contractType": "设备采购合同",
		"contractor": "供应商-001",
		"amount": 100000,
		"description": "设备采购合同审批"
	}`)
	tsk1, err := taskMgr.Create(tpl.ID, "contract-001", params1)
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
	fmt.Println("  条件判断: 合同类型包含\"采购\" OR 合同类型以\"服务\"开头")
	fmt.Println("  实际结果: 合同类型=\"设备采购合同\",包含\"采购\",条件满足")
	fmt.Println("  预期路径: 采购审批路径")
	fmt.Println()

	// 场景2: 合同类型以"服务"开头
	fmt.Println("=== 场景 2: 合同类型以\"服务\"开头 ===")
	fmt.Println()
	params2 := json.RawMessage(`{
		"contractNo": "CT-2025-002",
		"contractType": "服务外包合同",
		"contractor": "服务商-001",
		"amount": 50000,
		"description": "服务外包合同审批"
	}`)
	tsk2, err := taskMgr.Create(tpl.ID, "contract-002", params2)
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
	fmt.Println("  条件判断: 合同类型包含\"采购\" OR 合同类型以\"服务\"开头")
	fmt.Println("  实际结果: 合同类型=\"服务外包合同\",以\"服务\"开头,条件满足")
	fmt.Println("  预期路径: 采购审批路径")
	fmt.Println()

	// 场景3: 其他合同类型
	fmt.Println("=== 场景 3: 其他合同类型 ===")
	fmt.Println()
	params3 := json.RawMessage(`{
		"contractNo": "CT-2025-003",
		"contractType": "租赁合同",
		"contractor": "出租方-001",
		"amount": 30000,
		"description": "租赁合同审批"
	}`)
	tsk3, err := taskMgr.Create(tpl.ID, "contract-003", params3)
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
	fmt.Println("  条件判断: 合同类型包含\"采购\" OR 合同类型以\"服务\"开头")
	fmt.Println("  实际结果: 合同类型=\"租赁合同\",不包含\"采购\"且不以\"服务\"开头,条件不满足")
	fmt.Println("  预期路径: 标准审批路径")
	fmt.Println()

	// 总结
	fmt.Println("=== 字符串匹配条件说明 ===")
	fmt.Println("  字符串匹配条件配置:")
	fmt.Println("    - 操作符: contains (包含)")
	fmt.Println("    - 操作符: starts_with (以...开始)")
	fmt.Println("    - 操作符: ends_with (以...结束)")
	fmt.Println("    - 操作符: eq (等于)")
	fmt.Println()
	fmt.Println("  条件组合:")
	fmt.Println("    - 条件1: 合同类型包含\"采购\"")
	fmt.Println("    - 条件2: 合同类型以\"服务\"开头")
	fmt.Println("    - 组合逻辑: OR (任意一个条件满足即可)")
	fmt.Println()
	fmt.Println("  审批路径:")
	fmt.Println("    - True (条件满足): 采购审批路径")
	fmt.Println("    - False (条件不满足): 标准审批路径")
	fmt.Println()
	fmt.Println("  注意: 本示例仅展示字符串匹配条件的配置方法")
	fmt.Println("  实际执行需要在实际业务系统中实现条件评估和流程跳转")
	fmt.Println()
}

