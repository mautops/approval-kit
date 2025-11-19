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
	fmt.Println("=== 场景 22: 状态变更历史查询场景 ===")
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

	// 3. 创建任务并执行完整流程
	fmt.Println("步骤 2: 创建任务并执行完整流程")
	tsk := createAndProcessTask(templateMgr, taskMgr, tpl)
	fmt.Println()

	// 4. 查询和分析状态变更历史
	analyzeStateHistory(taskMgr, tsk)
}

// createTemplate 创建审批模板
func createTemplate() *template.Template {
	now := time.Now()
	return &template.Template{
		ID:          "state-history-template",
		Name:        "状态变更历史模板",
		Description: "演示状态变更历史查询功能",
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

// createAndProcessTask 创建任务并执行完整流程
func createAndProcessTask(templateMgr template.TemplateManager, taskMgr task.TaskManager, tpl *template.Template) *task.Task {
	// 创建任务
	fmt.Println("  1. 创建任务")
	params := json.RawMessage(`{
		"requestNo": "REQ-2025-001",
		"requestType": "状态变更历史测试",
		"requester": "申请人-001",
		"amount": 50000,
		"description": "测试状态变更历史功能"
	}`)
	tsk, err := taskMgr.Create(tpl.ID, "state-history-001", params)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}
	fmt.Printf("    ✓ 任务创建成功: ID=%s, State=%s\n", tsk.ID, tsk.State)

	// 提交任务
	fmt.Println("  2. 提交任务")
	err = taskMgr.Submit(tsk.ID)
	if err != nil {
		log.Fatalf("Failed to submit task: %v", err)
	}
	tsk, _ = taskMgr.Get(tsk.ID)
	fmt.Printf("    ✓ 任务已提交: State=%s\n", tsk.State)

	// 设置审批人
	fmt.Println("  3. 设置审批人")
	err = taskMgr.AddApprover(tsk.ID, "approval", "manager-001", "设置审批人")
	if err != nil {
		log.Fatalf("Failed to add approver: %v", err)
	}
	fmt.Printf("    ✓ 审批人已设置: manager-001\n")

	// 审批操作
	fmt.Println("  4. 审批操作")
	err = taskMgr.Approve(tsk.ID, "approval", "manager-001", "审批通过")
	if err != nil {
		log.Fatalf("Failed to approve: %v", err)
	}
	tsk, _ = taskMgr.Get(tsk.ID)
	fmt.Printf("    ✓ 审批已通过: State=%s\n", tsk.State)

	return tsk
}

// analyzeStateHistory 查询和分析状态变更历史
func analyzeStateHistory(taskMgr task.TaskManager, tsk *task.Task) {
	fmt.Println("=== 状态变更历史查询和分析 ===")
	fmt.Println()

	// 获取最新任务状态
	tsk, err := taskMgr.Get(tsk.ID)
	if err != nil {
		log.Fatalf("Failed to get task: %v", err)
	}

	// 输出任务基本信息
	fmt.Println("任务基本信息:")
	fmt.Printf("  任务 ID: %s\n", tsk.ID)
	fmt.Printf("  业务 ID: %s\n", tsk.BusinessID)
	fmt.Printf("  模板 ID: %s\n", tsk.TemplateID)
	fmt.Printf("  当前状态: %s\n", tsk.State)
	fmt.Printf("  创建时间: %s\n", tsk.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println()

	// 输出状态变更历史
	fmt.Println("状态变更历史:")
	history := tsk.GetStateHistory()
	if len(history) == 0 {
		fmt.Println("  无状态变更记录")
	} else {
		fmt.Printf("  共 %d 条状态变更记录:\n\n", len(history))
		for i, change := range history {
			fmt.Printf("  变更 %d:\n", i+1)
			fmt.Printf("    源状态: %s\n", change.From)
			fmt.Printf("    目标状态: %s\n", change.To)
			fmt.Printf("    变更原因: %s\n", change.Reason)
			fmt.Printf("    变更时间: %s\n", change.Time.Format("2006-01-02 15:04:05"))
			fmt.Printf("    状态流转: %s → %s\n", change.From, change.To)
			fmt.Println()
		}
	}

	// 状态变更分析
	fmt.Println("=== 状态变更分析 ===")
	fmt.Println()

	if len(history) > 0 {
		// 状态流转路径
		fmt.Println("状态流转路径:")
		path := fmt.Sprintf("%s", history[0].From)
		for _, change := range history {
			path += fmt.Sprintf(" → %s", change.To)
		}
		fmt.Printf("  %s\n", path)
		fmt.Println()

		// 状态变更统计
		fmt.Println("状态变更统计:")
		stateCounts := make(map[types.TaskState]int)
		for _, change := range history {
			stateCounts[change.To]++
		}
		for state, count := range stateCounts {
			fmt.Printf("  - %s: %d 次\n", state, count)
		}
		fmt.Println()

		// 状态变更时长分析
		fmt.Println("状态变更时长分析:")
		if len(history) > 0 {
			firstChange := history[0]
			lastChange := history[len(history)-1]
			totalDuration := lastChange.Time.Sub(firstChange.Time)
			fmt.Printf("  首次状态变更: %s\n", firstChange.Time.Format("2006-01-02 15:04:05"))
			fmt.Printf("  最后状态变更: %s\n", lastChange.Time.Format("2006-01-02 15:04:05"))
			fmt.Printf("  总时长: %s\n", formatDuration(totalDuration))
			fmt.Println()

			// 各状态停留时长
			if len(history) > 1 {
				fmt.Println("  各状态停留时长:")
				for i := 0; i < len(history)-1; i++ {
					currentChange := history[i]
					nextChange := history[i+1]
					duration := nextChange.Time.Sub(currentChange.Time)
					fmt.Printf("    %s: %s\n", currentChange.To, formatDuration(duration))
				}
			}
		}
		fmt.Println()

		// 状态变更原因分析
		fmt.Println("状态变更原因分析:")
		reasonCounts := make(map[string]int)
		for _, change := range history {
			reasonCounts[change.Reason]++
		}
		for reason, count := range reasonCounts {
			fmt.Printf("  - %s: %d 次\n", reason, count)
		}
		fmt.Println()
	}

	// 审计信息
	fmt.Println("=== 审计信息 ===")
	fmt.Println()
	fmt.Println("✓ 状态变更历史完整记录:")
	fmt.Println("  - 每次状态变更都记录了源状态、目标状态、变更原因、变更时间")
	fmt.Println("  - 可以完整追溯任务的状态流转过程")
	fmt.Println("  - 支持审计和合规要求")
	fmt.Println()
	fmt.Println("✓ 状态变更历史查询功能已展示")
	fmt.Println("  说明: 在实际业务系统中,可以通过 StateHistory 字段查询完整的状态变更历史")
}

// formatDuration 格式化时间间隔
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0f秒", d.Seconds())
	} else if d < time.Hour {
		return fmt.Sprintf("%.0f分钟", d.Minutes())
	} else if d < 24*time.Hour {
		return fmt.Sprintf("%.1f小时", d.Hours())
	} else {
		days := d.Hours() / 24
		hours := d.Hours() - days*24
		return fmt.Sprintf("%.0f天%.0f小时", days, hours)
	}
}

