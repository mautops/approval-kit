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
	fmt.Println("=== 场景 18: 枚举判断条件场景 ===")
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

	// 3. 演示不同枚举判断场景
	demonstrateEnumConditions(templateMgr, taskMgr, tpl)
}

// createTemplate 创建请假审批模板
// 流程: 开始节点 → 条件节点(枚举判断) → [部门审批/HR审批/总经理审批] → 结束节点
// 条件节点根据请假类型进行枚举判断
func createTemplate() *template.Template {
	now := time.Now()
	return &template.Template{
		ID:          "enum-condition-template",
		Name:        "请假审批模板",
		Description: "根据请假类型(枚举值)选择不同的审批路径",
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
									Type: "enum",
									Config: &node.EnumConditionConfig{
										Field:    "leaveType",
										Operator: "in",
										Values:   []string{"年假", "调休"},
										Source:   "task_params",
									},
								},
								{
									Type: "enum",
									Config: &node.EnumConditionConfig{
										Field:    "leaveType",
										Operator: "in",
										Values:   []string{"病假", "事假"},
										Source:   "task_params",
									},
								},
							},
						},
					},
					TrueNodeID:  "dept-approval",  // 年假/调休 或 病假/事假,走部门审批路径
					FalseNodeID: "hr-approval",    // 其他类型,走HR审批路径
				},
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
			"hr-approval": {
				ID:   "hr-approval",
				Name: "HR审批",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeSingle,
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{"hr-001"},
					},
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
			{From: "start", To: "condition-node"},
			{From: "condition-node", To: "dept-approval"},
			{From: "condition-node", To: "hr-approval"},
			{From: "dept-approval", To: "end"},
			{From: "hr-approval", To: "ceo-approval"},
			{From: "ceo-approval", To: "end"},
		},
	}
}

// demonstrateEnumConditions 演示不同枚举判断场景
func demonstrateEnumConditions(templateMgr template.TemplateManager, taskMgr task.TaskManager, tpl *template.Template) {
	// 场景1: 请假类型为"年假"
	fmt.Println("=== 场景 1: 请假类型为\"年假\" ===")
	fmt.Println()
	params1 := json.RawMessage(`{
		"requestNo": "LEAVE-2025-001",
		"leaveType": "年假",
		"applicant": "员工-001",
		"startDate": "2025-12-01",
		"endDate": "2025-12-05",
		"days": 5,
		"description": "年假申请"
	}`)
	tsk1, err := taskMgr.Create(tpl.ID, "leave-001", params1)
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
	fmt.Println("  条件判断: 请假类型 in [\"年假\", \"调休\"] OR 请假类型 in [\"病假\", \"事假\"]")
	fmt.Println("  实际结果: 请假类型=\"年假\",在第一个列表中,条件满足")
	fmt.Println("  预期路径: 部门审批路径")
	fmt.Println()

	// 场景2: 请假类型为"调休"
	fmt.Println("=== 场景 2: 请假类型为\"调休\" ===")
	fmt.Println()
	params2 := json.RawMessage(`{
		"requestNo": "LEAVE-2025-002",
		"leaveType": "调休",
		"applicant": "员工-002",
		"startDate": "2025-12-10",
		"endDate": "2025-12-10",
		"days": 1,
		"description": "调休申请"
	}`)
	tsk2, err := taskMgr.Create(tpl.ID, "leave-002", params2)
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
	fmt.Println("  条件判断: 请假类型 in [\"年假\", \"调休\"] OR 请假类型 in [\"病假\", \"事假\"]")
	fmt.Println("  实际结果: 请假类型=\"调休\",在第一个列表中,条件满足")
	fmt.Println("  预期路径: 部门审批路径")
	fmt.Println()

	// 场景3: 请假类型为"病假"
	fmt.Println("=== 场景 3: 请假类型为\"病假\" ===")
	fmt.Println()
	params3 := json.RawMessage(`{
		"requestNo": "LEAVE-2025-003",
		"leaveType": "病假",
		"applicant": "员工-003",
		"startDate": "2025-12-15",
		"endDate": "2025-12-17",
		"days": 3,
		"description": "病假申请"
	}`)
	tsk3, err := taskMgr.Create(tpl.ID, "leave-003", params3)
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
	fmt.Println("  条件判断: 请假类型 in [\"年假\", \"调休\"] OR 请假类型 in [\"病假\", \"事假\"]")
	fmt.Println("  实际结果: 请假类型=\"病假\",在第二个列表中,条件满足")
	fmt.Println("  预期路径: 部门审批路径")
	fmt.Println()

	// 场景4: 请假类型为"其他"
	fmt.Println("=== 场景 4: 请假类型为\"其他\" ===")
	fmt.Println()
	params4 := json.RawMessage(`{
		"requestNo": "LEAVE-2025-004",
		"leaveType": "其他",
		"applicant": "员工-004",
		"startDate": "2025-12-20",
		"endDate": "2025-12-25",
		"days": 6,
		"description": "其他类型请假申请"
	}`)
	tsk4, err := taskMgr.Create(tpl.ID, "leave-004", params4)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}
	fmt.Printf("✓ 任务创建成功: ID=%s\n", tsk4.ID)
	fmt.Println("  任务参数:")
	var params4Map map[string]interface{}
	json.Unmarshal(params4, &params4Map)
	for k, v := range params4Map {
		fmt.Printf("    %s: %v\n", k, v)
	}
	fmt.Println("  条件判断: 请假类型 in [\"年假\", \"调休\"] OR 请假类型 in [\"病假\", \"事假\"]")
	fmt.Println("  实际结果: 请假类型=\"其他\",不在任何列表中,条件不满足")
	fmt.Println("  预期路径: HR审批路径 → 总经理审批路径")
	fmt.Println()

	// 总结
	fmt.Println("=== 枚举判断条件说明 ===")
	fmt.Println("  枚举判断条件配置:")
	fmt.Println("    - 操作符: in (在列表中)")
	fmt.Println("    - 操作符: not_in (不在列表中)")
	fmt.Println()
	fmt.Println("  条件组合:")
	fmt.Println("    - 条件1: 请假类型 in [\"年假\", \"调休\"]")
	fmt.Println("    - 条件2: 请假类型 in [\"病假\", \"事假\"]")
	fmt.Println("    - 组合逻辑: OR (任意一个条件满足即可)")
	fmt.Println()
	fmt.Println("  审批路径:")
	fmt.Println("    - True (条件满足): 部门审批路径")
	fmt.Println("    - False (条件不满足): HR审批路径 → 总经理审批路径")
	fmt.Println()
	fmt.Println("  注意: 本示例仅展示枚举判断条件的配置方法")
	fmt.Println("  实际执行需要在实际业务系统中实现条件评估和流程跳转")
	fmt.Println()
}

