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
	fmt.Println("=== 场景 15: 模板版本管理场景 ===")
	fmt.Println()

	// 1. 创建管理器
	templateMgr := template.NewTemplateManager()
	taskMgr := task.NewTaskManager(templateMgr, nil)

	// 2. 创建模板版本1
	fmt.Println("步骤 1: 创建模板版本1")
	tpl1 := createTemplateV1()
	err := templateMgr.Create(tpl1)
	if err != nil {
		log.Fatalf("Failed to create template: %v", err)
	}
	fmt.Printf("✓ 模板版本1创建成功: ID=%s, Version=%d, Name=%s\n\n", tpl1.ID, tpl1.Version, tpl1.Name)

	// 3. 使用版本1创建任务
	fmt.Println("步骤 2: 使用版本1创建任务")
	params1 := json.RawMessage(`{"requestNo": "REQ-001", "amount": 10000}`)
	tsk1, err := taskMgr.Create(tpl1.ID, "biz-001", params1)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}
	fmt.Printf("✓ 任务1创建成功: ID=%s, TemplateVersion=%d\n\n", tsk1.ID, tsk1.TemplateVersion)

	// 4. 更新模板到版本2
	fmt.Println("步骤 3: 更新模板到版本2")
	tpl2 := createTemplateV2()
	err = templateMgr.Update(tpl1.ID, tpl2)
	if err != nil {
		log.Fatalf("Failed to update template: %v", err)
	}
	fmt.Printf("✓ 模板已更新到版本2: ID=%s, Version=%d, Name=%s\n\n", tpl2.ID, tpl2.Version, tpl2.Name)

	// 5. 使用版本2创建任务
	fmt.Println("步骤 4: 使用版本2创建任务")
	params2 := json.RawMessage(`{"requestNo": "REQ-002", "amount": 20000}`)
	tsk2, err := taskMgr.Create(tpl2.ID, "biz-002", params2)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}
	fmt.Printf("✓ 任务2创建成功: ID=%s, TemplateVersion=%d\n\n", tsk2.ID, tsk2.TemplateVersion)

	// 6. 更新模板到版本3
	fmt.Println("步骤 5: 更新模板到版本3")
	tpl3 := createTemplateV3()
	err = templateMgr.Update(tpl1.ID, tpl3)
	if err != nil {
		log.Fatalf("Failed to update template: %v", err)
	}
	fmt.Printf("✓ 模板已更新到版本3: ID=%s, Version=%d, Name=%s\n\n", tpl3.ID, tpl3.Version, tpl3.Name)

	// 7. 使用版本3创建任务
	fmt.Println("步骤 6: 使用版本3创建任务")
	params3 := json.RawMessage(`{"requestNo": "REQ-003", "amount": 30000}`)
	tsk3, err := taskMgr.Create(tpl3.ID, "biz-003", params3)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}
	fmt.Printf("✓ 任务3创建成功: ID=%s, TemplateVersion=%d\n\n", tsk3.ID, tsk3.TemplateVersion)

	// 8. 查询所有版本
	fmt.Println("步骤 7: 查询所有版本")
	versions, err := templateMgr.ListVersions(tpl1.ID)
	if err != nil {
		log.Fatalf("Failed to list versions: %v", err)
	}
	fmt.Printf("✓ 模板共有 %d 个版本: %v\n\n", len(versions), versions)

	// 9. 获取不同版本的模板
	fmt.Println("步骤 8: 获取不同版本的模板")
	for _, version := range versions {
		tpl, err := templateMgr.Get(tpl1.ID, version)
		if err != nil {
			log.Fatalf("Failed to get template version %d: %v", version, err)
		}
		fmt.Printf("  版本 %d: Name=%s, CreatedAt=%s, UpdatedAt=%s\n",
			tpl.Version,
			tpl.Name,
			tpl.CreatedAt.Format("2006-01-02 15:04:05"),
			tpl.UpdatedAt.Format("2006-01-02 15:04:05"),
		)
	}
	fmt.Println()

	// 10. 验证历史任务使用历史版本
	fmt.Println("步骤 9: 验证历史任务使用历史版本")
	fmt.Println("  任务1使用模板版本1:")
	verifyTaskTemplate(taskMgr, templateMgr, tsk1.ID, 1)
	fmt.Println()
	fmt.Println("  任务2使用模板版本2:")
	verifyTaskTemplate(taskMgr, templateMgr, tsk2.ID, 2)
	fmt.Println()
	fmt.Println("  任务3使用模板版本3:")
	verifyTaskTemplate(taskMgr, templateMgr, tsk3.ID, 3)
	fmt.Println()

	// 11. 版本对比
	fmt.Println("步骤 10: 版本对比")
	compareVersions(templateMgr, tpl1.ID, versions)
}

// createTemplateV1 创建模板版本1
func createTemplateV1() *template.Template {
	now := time.Now()
	return &template.Template{
		ID:          "versioning-template",
		Name:        "版本管理模板 V1",
		Description: "版本1: 单人审批",
		Version:     1,
		CreatedAt:   now,
		UpdatedAt:   now,
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "开始",
				Type: template.NodeTypeStart,
			},
			"approval": {
				ID:   "approval",
				Name: "审批",
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
			{From: "start", To: "approval"},
			{From: "approval", To: "end"},
		},
	}
}

// createTemplateV2 创建模板版本2
func createTemplateV2() *template.Template {
	now := time.Now()
	return &template.Template{
		ID:          "versioning-template",
		Name:        "版本管理模板 V2",
		Description: "版本2: 多人会签",
		Version:     0, // Update 方法会自动递增版本号
		CreatedAt:   now,
		UpdatedAt:   now,
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "开始",
				Type: template.NodeTypeStart,
			},
			"approval": {
				ID:   "approval",
				Name: "审批",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeUnanimous, // 改为多人会签
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{
							"manager-001",
							"finance-001",
						},
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
			{From: "start", To: "approval"},
			{From: "approval", To: "end"},
		},
	}
}

// createTemplateV3 创建模板版本3
func createTemplateV3() *template.Template {
	now := time.Now()
	return &template.Template{
		ID:          "versioning-template",
		Name:        "版本管理模板 V3",
		Description: "版本3: 添加条件节点",
		Version:     0, // Update 方法会自动递增版本号
		CreatedAt:   now,
		UpdatedAt:   now,
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "开始",
				Type: template.NodeTypeStart,
			},
			"condition": {
				ID:   "condition",
				Name: "条件判断",
				Type: template.NodeTypeCondition,
				Config: &node.ConditionNodeConfig{
					Condition: &node.Condition{
						Type: "numeric",
						Config: &node.NumericConditionConfig{
							Field:    "amount",
							Operator: ">=",
							Value:    20000,
							Source:   "task_params",
						},
					},
					TrueNodeID:  "high-approval",
					FalseNodeID: "normal-approval",
				},
			},
			"high-approval": {
				ID:   "high-approval",
				Name: "高级审批",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeUnanimous,
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{
							"manager-001",
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
			{From: "start", To: "condition"},
			{From: "condition", To: "high-approval"},
			{From: "condition", To: "normal-approval"},
			{From: "high-approval", To: "end"},
			{From: "normal-approval", To: "end"},
		},
	}
}

// verifyTaskTemplate 验证任务使用的模板版本
func verifyTaskTemplate(taskMgr task.TaskManager, templateMgr template.TemplateManager, taskID string, expectedVersion int) {
	tsk, err := taskMgr.Get(taskID)
	if err != nil {
		log.Fatalf("Failed to get task: %v", err)
	}

	tpl, err := templateMgr.Get(tsk.TemplateID, tsk.TemplateVersion)
	if err != nil {
		log.Fatalf("Failed to get template: %v", err)
	}

	fmt.Printf("    任务 ID: %s\n", tsk.ID)
	fmt.Printf("    模板 ID: %s\n", tsk.TemplateID)
	fmt.Printf("    模板版本: %d (期望: %d)\n", tsk.TemplateVersion, expectedVersion)
	fmt.Printf("    模板名称: %s\n", tpl.Name)
	if tsk.TemplateVersion == expectedVersion {
		fmt.Printf("    ✓ 版本匹配正确\n")
	} else {
		fmt.Printf("    ✗ 版本不匹配\n")
	}
}

// compareVersions 对比不同版本的模板
func compareVersions(templateMgr template.TemplateManager, templateID string, versions []int) {
	fmt.Println("  版本对比:")
	fmt.Println()

	for i, version := range versions {
		tpl, err := templateMgr.Get(templateID, version)
		if err != nil {
			log.Fatalf("Failed to get template version %d: %v", version, err)
		}

		fmt.Printf("  版本 %d: %s\n", tpl.Version, tpl.Name)
		fmt.Printf("    描述: %s\n", tpl.Description)
		fmt.Printf("    节点数: %d\n", len(tpl.Nodes))
		fmt.Printf("    边数: %d\n", len(tpl.Edges))

		// 列出节点类型
		nodeTypes := make(map[string]int)
		for _, node := range tpl.Nodes {
			nodeTypes[string(node.Type)]++
		}
		fmt.Printf("    节点类型分布:")
		for nodeType, count := range nodeTypes {
			fmt.Printf(" %s:%d", nodeType, count)
		}
		fmt.Println()

		if i < len(versions)-1 {
			fmt.Println()
		}
	}
}

