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
	fmt.Println("=== 场景 23: 基于节点输出的条件判断场景 ===")

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

	// 运行场景: 测试不同评分的情况
	scenarios := []struct {
		name   string
		score  int
		expect string
	}{
		{"高分场景(>=80)", 85, "直接进入最终审批"},
		{"中分场景(60-80)", 70, "需要修改后重新评审"},
		{"低分场景(<60)", 45, "直接拒绝"},
	}

	for _, scenario := range scenarios {
		fmt.Printf("--- %s ---\n", scenario.name)
		runScenario(taskMgr, tpl.ID, scenario.score, scenario.expect)
		fmt.Println()
	}
}

// createTemplate 创建技术方案评审模板
func createTemplate() *template.Template {
	now := time.Now()

	// 技术评审节点(输出评分)
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

	// 条件节点: 根据评分判断路径
	// 使用组合条件实现三个分支:
	// 1. 评分 >= 80: 直接进入最终审批
	// 2. 评分 < 80 且 >= 60: 需要修改后重新评审
	// 3. 评分 < 60: 直接拒绝
	//
	// 注意: 由于条件节点只支持 true/false 两个分支,我们需要使用两个条件节点来实现三个分支
	// 第一个条件节点: 评分 >= 80 -> 最终审批 / 否则进入第二个条件节点
	// 第二个条件节点: 评分 >= 60 -> 重新评审 / 否则拒绝

	// 条件节点1: 评分 >= 80
	conditionNode1 := &template.Node{
		ID:   "score-condition-1",
		Name: "评分条件判断1(>=80)",
		Type: template.NodeTypeCondition,
		Config: &node.ConditionNodeConfig{
			Condition: &node.Condition{
				Type: "numeric",
				Config: &node.NumericConditionConfig{
					Field:    "score",
					Operator: "gte",
					Value:    80,
					Source:   "node_outputs",
					NodeID:   "tech-review", // 从技术评审节点的输出中读取
				},
			},
			TrueNodeID:  "final-approval",      // 评分 >= 80: 最终审批
			FalseNodeID: "score-condition-2",  // 评分 < 80: 进入第二个条件节点
		},
		Order: 2,
	}

	// 条件节点2: 评分 >= 60
	conditionNode2 := &template.Node{
		ID:   "score-condition-2",
		Name: "评分条件判断2(>=60)",
		Type: template.NodeTypeCondition,
		Config: &node.ConditionNodeConfig{
			Condition: &node.Condition{
				Type: "numeric",
				Config: &node.NumericConditionConfig{
					Field:    "score",
					Operator: "gte",
					Value:    60,
					Source:   "node_outputs",
					NodeID:   "tech-review", // 从技术评审节点的输出中读取
				},
			},
			TrueNodeID:  "re-review",    // 评分 >= 60: 重新评审
			FalseNodeID: "reject-end",   // 评分 < 60: 拒绝
		},
		Order: 3,
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
		Order: 4,
	}

	// 重新评审节点
	reReviewNode := &template.Node{
		ID:   "re-review",
		Name: "重新评审",
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
		Order: 5,
	}

	// 拒绝结束节点
	rejectEndNode := &template.Node{
		ID:   "reject-end",
		Name: "拒绝结束",
		Type: template.NodeTypeEnd,
		Order: 6,
	}

	// 正常结束节点
	normalEndNode := &template.Node{
		ID:   "normal-end",
		Name: "正常结束",
		Type: template.NodeTypeEnd,
		Order: 7,
	}

	// 开始节点
	startNode := &template.Node{
		ID:   "start",
		Name: "开始",
		Type: template.NodeTypeStart,
		Order: 0,
	}

	return &template.Template{
		ID:          "tech-review-template",
		Name:        "技术方案评审模板",
		Description: "根据技术评审的评分决定后续审批路径",
		Version:     1,
		CreatedAt:   now,
		UpdatedAt:   now,
		Nodes: map[string]*template.Node{
			"start":              startNode,
			"tech-review":        techReviewNode,
			"score-condition-1":  conditionNode1,
			"score-condition-2":  conditionNode2,
			"final-approval":     finalApprovalNode,
			"re-review":          reReviewNode,
			"reject-end":         rejectEndNode,
			"normal-end":         normalEndNode,
		},
		Edges: []*template.Edge{
			{From: "start", To: "tech-review"},
			{From: "tech-review", To: "score-condition-1"},
			{From: "score-condition-1", To: "final-approval", Condition: "score >= 80"},
			{From: "score-condition-1", To: "score-condition-2", Condition: "score < 80"},
			{From: "score-condition-2", To: "re-review", Condition: "score >= 60"},
			{From: "score-condition-2", To: "reject-end", Condition: "score < 60"},
			{From: "final-approval", To: "normal-end"},
			{From: "re-review", To: "normal-end"},
		},
	}
}

// runScenario 运行场景
func runScenario(taskMgr task.TaskManager, templateID string, score int, expectPath string) {
	// 1. 创建任务
	taskParams := map[string]interface{}{
		"projectName": "新功能开发项目",
		"description":  "开发新的用户管理功能",
	}
	paramsJSON, _ := json.Marshal(taskParams)

	taskID := fmt.Sprintf("task-%d-%d", time.Now().UnixNano(), score)
	tsk, err := taskMgr.Create(templateID, fmt.Sprintf("project-%d", score), paramsJSON)
	if err != nil {
		fmt.Printf("❌ 任务创建失败: %v\n", err)
		return
	}
	taskID = tsk.ID
	fmt.Printf("✓ 任务创建成功: ID=%s, 评分=%d\n", taskID, score)

	// 2. 提交任务
	if err := taskMgr.Submit(taskID); err != nil {
		fmt.Printf("❌ 任务提交失败: %v\n", err)
		return
	}
	fmt.Printf("✓ 任务已提交\n")

	// 3. 模拟技术评审节点输出评分
	// 注意: 在实际业务系统中,节点执行时会自动生成输出数据
	// 这里我们手动设置节点输出数据来模拟这个过程
	tsk, _ = taskMgr.Get(taskID)
	tsk.Update(func(t *task.Task) error {
		// 设置技术评审节点的输出数据
		outputData := map[string]interface{}{
			"score":      score,
			"comment":    fmt.Sprintf("技术评审评分: %d", score),
			"reviewer":   "tech-lead-001",
			"reviewTime": time.Now().Format("2006-01-02 15:04:05"),
		}
		outputJSON, _ := json.Marshal(outputData)
		if t.NodeOutputs == nil {
			t.NodeOutputs = make(map[string]json.RawMessage)
		}
		t.NodeOutputs["tech-review"] = outputJSON
		return nil
	})
	fmt.Printf("✓ 技术评审节点输出已设置: score=%d\n", score)

	// 4. 模拟条件节点判断
	// 注意: 在实际业务系统中,条件节点会自动执行并跳转到下一个节点
	// 这里我们手动判断并说明预期的路径
	tsk, _ = taskMgr.Get(taskID)
	var pathDescription string

	if score >= 80 {
		pathDescription = "评分 >= 80, 进入最终审批"
	} else if score >= 60 {
		pathDescription = "评分 < 80 且 >= 60, 需要修改后重新评审"
	} else {
		pathDescription = "评分 < 60, 直接拒绝"
	}

	fmt.Printf("✓ 条件节点判断结果: %s\n", pathDescription)
	fmt.Printf("  预期路径: %s\n", expectPath)
	fmt.Printf("  实际路径: %s\n", pathDescription)

	// 5. 验证节点输出数据
	tsk, _ = taskMgr.Get(taskID)
	if output, exists := tsk.NodeOutputs["tech-review"]; exists {
		var outputData map[string]interface{}
		if err := json.Unmarshal(output, &outputData); err == nil {
			if outputScore, ok := outputData["score"].(float64); ok {
				fmt.Printf("✓ 节点输出数据验证: tech-review.score = %.0f\n", outputScore)
			}
		}
	} else {
		fmt.Printf("⚠ 节点输出数据未找到: tech-review\n")
	}

	// 6. 说明完整流程
	fmt.Printf("\n完整流程说明:\n")
	fmt.Printf("  1. 技术评审节点(tech-review)执行,输出评分: %d\n", score)
	fmt.Printf("  2. 条件节点1(score-condition-1)读取 tech-review 的输出,判断评分是否 >= 80\n")
	if score >= 80 {
		fmt.Printf("  3. 条件满足,跳转到最终审批节点(final-approval)\n")
		fmt.Printf("  4. 最终审批通过后,流程结束\n")
	} else {
		fmt.Printf("  3. 条件不满足,进入条件节点2(score-condition-2),判断评分是否 >= 60\n")
		if score >= 60 {
			fmt.Printf("  4. 条件满足,跳转到重新评审节点(re-review)\n")
			fmt.Printf("  5. 重新评审通过后,流程结束\n")
		} else {
			fmt.Printf("  4. 条件不满足,跳转到拒绝结束节点(reject-end)\n")
			fmt.Printf("  5. 流程结束,任务被拒绝\n")
		}
	}

	fmt.Printf("\n注意: 此示例主要展示节点输出数据作为条件判断的配置方式.\n")
	fmt.Printf("在实际业务系统中,条件节点会自动执行并跳转到下一个节点.\n")
}

