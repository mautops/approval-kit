package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
)

func main() {
	fmt.Println("=== 场景 25: 多路径并行审批后汇聚场景 ===")

	// 创建管理器
	templateMgr := template.NewTemplateManager()
	taskMgr := task.NewTaskManager(templateMgr, nil)

	// 创建模板
	tpl := createTemplate()
	if err := templateMgr.Create(tpl); err != nil {
		fmt.Printf("❌ 模板创建失败: %v\n", err)
		return
	}
	fmt.Printf("✓ 模板创建成功: ID=%s, Name=%s\n\n", tpl.ID, tpl.Name)

	// 运行场景
	runScenario(taskMgr, tpl.ID)
}

// createTemplate 创建多路径并行审批模板
// 流程: 开始节点 → [技术评审节点, 财务评审节点] (并行) → 最终审批节点 → 结束节点
func createTemplate() *template.Template {
	now := time.Now()

	// 技术评审节点
	techReviewNode := &template.Node{
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
		Order: 1,
	}

	// 财务评审节点
	financeReviewNode := &template.Node{
		ID:   "finance-review",
		Name: "财务评审",
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
		Order: 2,
	}

	// 最终审批节点
	finalApprovalNode := &template.Node{
		ID:   "final-approval",
		Name: "最终审批",
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
		Order: 3,
	}

	// 开始节点
	startNode := &template.Node{
		ID:   "start",
		Name: "开始",
		Type: template.NodeTypeStart,
		Order: 0,
	}

	// 结束节点
	endNode := &template.Node{
		ID:   "end",
		Name: "结束",
		Type: template.NodeTypeEnd,
		Order: 4,
	}

	return &template.Template{
		ID:          "parallel-merge-template",
		Name:        "多路径并行审批模板",
		Description: "技术评审和财务评审并行进行,都完成后进入最终审批",
		Version:     1,
		CreatedAt:   now,
		UpdatedAt:   now,
		Nodes: map[string]*template.Node{
			"start":         startNode,
			"tech-review":   techReviewNode,
			"finance-review": financeReviewNode,
			"final-approval": finalApprovalNode,
			"end":           endNode,
		},
		Edges: []*template.Edge{
			// 开始节点分支到两个并行审批节点
			{From: "start", To: "tech-review"},
			{From: "start", To: "finance-review"},
			// 两个并行审批节点汇聚到最终审批节点
			{From: "tech-review", To: "final-approval"},
			{From: "finance-review", To: "final-approval"},
			// 最终审批节点到结束节点
			{From: "final-approval", To: "end"},
		},
	}
}

// runScenario 运行场景
func runScenario(taskMgr task.TaskManager, templateID string) {
	// 1. 创建任务
	taskParams := map[string]interface{}{
		"projectName": "新产品开发项目",
		"budget":      5000000,
		"description": "开发新产品,需要技术评审和财务评审",
	}
	paramsJSON, _ := json.Marshal(taskParams)

	tsk, err := taskMgr.Create(templateID, "project-001", paramsJSON)
	if err != nil {
		fmt.Printf("❌ 任务创建失败: %v\n", err)
		return
	}
	taskID := tsk.ID
	fmt.Printf("✓ 任务创建成功: ID=%s\n", taskID)

	// 2. 提交任务
	if err := taskMgr.Submit(taskID); err != nil {
		fmt.Printf("❌ 任务提交失败: %v\n", err)
		return
	}
	fmt.Printf("✓ 任务已提交\n\n")

	// 3. 说明并行审批流程
	fmt.Println("=== 并行审批流程说明 ===")
	fmt.Println("流程结构:")
	fmt.Println("  开始节点")
	fmt.Println("    ↓")
	fmt.Println("  ├─→ 技术评审节点 (tech-review)")
	fmt.Println("  └─→ 财务评审节点 (finance-review)")
	fmt.Println("    ↓ (汇聚)")
	fmt.Println("  最终审批节点 (final-approval)")
	fmt.Println("    ↓")
	fmt.Println("  结束节点")
	fmt.Println()

	fmt.Println("并行审批说明:")
	fmt.Println("  1. 任务提交后,技术评审和财务评审可以并行进行")
	fmt.Println("  2. 两个评审节点都完成后,才能进入最终审批")
	fmt.Println("  3. 最终审批通过后,流程结束")
	fmt.Println()

	// 4. 模拟并行审批
	fmt.Println("=== 模拟并行审批 ===")
	fmt.Println("注意: 当前示例主要展示多路径并行审批的配置方式.")
	fmt.Println("由于当前库不包含工作流引擎,无法自动处理并行分支和汇聚逻辑.")
	fmt.Println("在实际业务系统中,需要工作流引擎来:")
	fmt.Println("  1. 监控并行节点的完成状态")
	fmt.Println("  2. 等待所有并行节点完成后触发汇聚")
	fmt.Println("  3. 将任务推进到下一个节点")
	fmt.Println()

	// 设置审批人
	fmt.Println("步骤 1: 设置审批人")
	if err := taskMgr.AddApprover(taskID, "tech-review", "tech-lead-001", "设置技术评审人"); err != nil {
		fmt.Printf("❌ 设置技术评审人失败: %v\n", err)
		return
	}
	fmt.Printf("✓ 技术评审人已设置: tech-lead-001\n")

	if err := taskMgr.AddApprover(taskID, "finance-review", "finance-001", "设置财务评审人"); err != nil {
		fmt.Printf("❌ 设置财务评审人失败: %v\n", err)
		return
	}
	fmt.Printf("✓ 财务评审人已设置: finance-001\n")
	fmt.Println()

	// 说明并行审批的限制
	fmt.Println("步骤 2: 并行审批说明")
	fmt.Println("  在实际业务系统中:")
	fmt.Println("    - 技术评审和财务评审可以同时进行")
	fmt.Println("    - 两个评审节点都完成后,才能进入最终审批")
	fmt.Println("    - 工作流引擎会监控节点完成状态并触发汇聚")
	fmt.Println()
	fmt.Println("  在当前示例中:")
	fmt.Println("    - 由于缺少工作流引擎,无法实现真正的并行分支和汇聚")
	fmt.Println("    - 第一个节点审批完成后,任务状态会立即变为 approved")
	fmt.Println("    - 这导致第二个节点无法继续审批")
	fmt.Println("    - 此示例主要展示模板配置和节点结构")
	fmt.Println()

	// 5. 说明汇聚逻辑
	fmt.Println("=== 汇聚逻辑说明 ===")
	fmt.Println("在实际业务系统中,工作流引擎会:")
	fmt.Println("  1. 监控两个并行审批节点的完成状态")
	fmt.Println("  2. 当两个节点都完成时,自动触发汇聚逻辑")
	fmt.Println("  3. 将任务推进到最终审批节点")
	fmt.Println("  4. 最终审批节点激活,等待审批")
	fmt.Println()
	fmt.Println("汇聚逻辑的关键点:")
	fmt.Println("  - 需要等待所有并行路径都完成")
	fmt.Println("  - 只有当所有前置节点都完成时,才能进入下一个节点")
	fmt.Println("  - 这需要工作流引擎的状态管理和节点依赖跟踪")
	fmt.Println()

	// 7. 输出结果
	printResults(taskMgr, taskID)
}

// printResults 输出结果
func printResults(taskMgr task.TaskManager, taskID string) {
	tsk, err := taskMgr.Get(taskID)
	if err != nil {
		fmt.Printf("❌ 获取任务失败: %v\n", err)
		return
	}

	fmt.Println("=== 审批结果 ===")
	fmt.Printf("任务 ID: %s\n", tsk.ID)
	fmt.Printf("业务 ID: %s\n", tsk.BusinessID)
	fmt.Printf("模板 ID: %s\n", tsk.TemplateID)
	fmt.Printf("当前状态: %s\n", tsk.GetState())
	fmt.Printf("当前节点: %s\n", tsk.GetCurrentNode())
	fmt.Println()

	// 输出审批记录
	records := tsk.GetRecords()
	fmt.Printf("审批记录 (共 %d 条):\n", len(records))
	for i, record := range records {
		fmt.Printf("  %d. 节点: %s, 审批人: %s, 结果: %s, 意见: %s\n",
			i+1,
			record.NodeID,
			record.Approver,
			record.Result,
			record.Comment,
		)
	}
	fmt.Println()

	// 验证并行审批
	fmt.Println("=== 验证并行审批 ===")
	techReviewCount := 0
	financeReviewCount := 0
	finalApprovalCount := 0

	for _, record := range records {
		if record.NodeID == "tech-review" && record.Result == "approve" {
			techReviewCount++
		}
		if record.NodeID == "finance-review" && record.Result == "approve" {
			financeReviewCount++
		}
		if record.NodeID == "final-approval" && record.Result == "approve" {
			finalApprovalCount++
		}
	}

	fmt.Printf("技术评审通过数: %d (期望: 1)\n", techReviewCount)
	fmt.Printf("财务评审通过数: %d (期望: 1)\n", financeReviewCount)
	fmt.Printf("最终审批通过数: %d (期望: 1)\n", finalApprovalCount)

	if techReviewCount == 1 && financeReviewCount == 1 && finalApprovalCount == 1 {
		fmt.Println("✓ 并行审批流程验证通过")
	} else {
		fmt.Println("⚠ 并行审批流程验证未完全通过")
	}
	fmt.Println()

	fmt.Println("注意: 此示例主要展示多路径并行审批的配置方式.")
	fmt.Println("在实际业务系统中,并行分支和汇聚逻辑需要工作流引擎的支持.")
}

