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
	fmt.Println("=== 场景 14: 审批记录查询和分析场景 ===")
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

	// 3. 创建多个任务并执行审批
	fmt.Println("步骤 2: 创建多个任务并执行审批")
	tasks := createAndProcessTasks(templateMgr, taskMgr, tpl)
	fmt.Println()

	// 4. 演示多维度查询
	demonstrateQuery(taskMgr, tasks)
}

// createTemplate 创建审批模板
func createTemplate() *template.Template {
	now := time.Now()
	return &template.Template{
		ID:          "record-query-template",
		Name:        "审批记录查询模板",
		Description: "演示审批记录查询和分析功能",
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
					Mode: node.ApprovalModeUnanimous, // 多人会签模式
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

// createAndProcessTasks 创建多个任务并执行审批
func createAndProcessTasks(templateMgr template.TemplateManager, taskMgr task.TaskManager, tpl *template.Template) []*task.Task {
	var tasks []*task.Task

	// 任务1: 已通过
	fmt.Println("  创建任务1: 已通过")
	params1 := json.RawMessage(`{"requestNo": "REQ-001", "amount": 10000}`)
	tsk1, err := taskMgr.Create(tpl.ID, "biz-001", params1)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}
	taskMgr.Submit(tsk1.ID)
	taskMgr.AddApprover(tsk1.ID, "approval", "manager-001", "设置审批人")
	taskMgr.AddApprover(tsk1.ID, "approval", "finance-001", "设置审批人")
	taskMgr.Approve(tsk1.ID, "approval", "manager-001", "经理审批通过")
	taskMgr.Approve(tsk1.ID, "approval", "finance-001", "财务审批通过")
	tsk1, _ = taskMgr.Get(tsk1.ID)
	tasks = append(tasks, tsk1)
	fmt.Printf("    ✓ 任务1已创建并审批通过: ID=%s, State=%s\n", tsk1.ID, tsk1.State)

	// 任务2: 审批中
	fmt.Println("  创建任务2: 审批中")
	params2 := json.RawMessage(`{"requestNo": "REQ-002", "amount": 20000}`)
	tsk2, err := taskMgr.Create(tpl.ID, "biz-002", params2)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}
	taskMgr.Submit(tsk2.ID)
	taskMgr.AddApprover(tsk2.ID, "approval", "manager-001", "设置审批人")
	taskMgr.AddApprover(tsk2.ID, "approval", "finance-001", "设置审批人")
	taskMgr.Approve(tsk2.ID, "approval", "manager-001", "经理审批通过")
	tsk2, _ = taskMgr.Get(tsk2.ID)
	tasks = append(tasks, tsk2)
	fmt.Printf("    ✓ 任务2已创建并部分审批: ID=%s, State=%s\n", tsk2.ID, tsk2.State)

	// 任务3: 已拒绝
	fmt.Println("  创建任务3: 已拒绝")
	params3 := json.RawMessage(`{"requestNo": "REQ-003", "amount": 30000}`)
	tsk3, err := taskMgr.Create(tpl.ID, "biz-003", params3)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}
	taskMgr.Submit(tsk3.ID)
	taskMgr.AddApprover(tsk3.ID, "approval", "manager-001", "设置审批人")
	taskMgr.Reject(tsk3.ID, "approval", "manager-001", "经理拒绝")
	tsk3, _ = taskMgr.Get(tsk3.ID)
	tasks = append(tasks, tsk3)
	fmt.Printf("    ✓ 任务3已创建并拒绝: ID=%s, State=%s\n", tsk3.ID, tsk3.State)

	// 任务4: 待审批
	fmt.Println("  创建任务4: 待审批")
	params4 := json.RawMessage(`{"requestNo": "REQ-004", "amount": 40000}`)
	tsk4, err := taskMgr.Create(tpl.ID, "biz-004", params4)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}
	taskMgr.Submit(tsk4.ID)
	tsk4, _ = taskMgr.Get(tsk4.ID)
	tasks = append(tasks, tsk4)
	fmt.Printf("    ✓ 任务4已创建: ID=%s, State=%s\n", tsk4.ID, tsk4.State)

	return tasks
}

// demonstrateQuery 演示多维度查询
func demonstrateQuery(taskMgr task.TaskManager, tasks []*task.Task) {
	fmt.Println("=== 多维度查询演示 ===")
	fmt.Println()

	// 查询1: 按状态查询
	fmt.Println("查询1: 按状态查询(已通过)")
	filter1 := &task.TaskFilter{
		State: types.TaskStateApproved,
	}
	results1, err := taskMgr.Query(filter1)
	if err != nil {
		log.Fatalf("Failed to query: %v", err)
	}
	fmt.Printf("  结果: 找到 %d 个已通过的任务\n", len(results1))
	for _, tsk := range results1 {
		fmt.Printf("    - 任务 ID: %s, 业务 ID: %s\n", tsk.ID, tsk.BusinessID)
	}
	fmt.Println()

	// 查询2: 按模板查询
	fmt.Println("查询2: 按模板查询")
	filter2 := &task.TaskFilter{
		TemplateID: "record-query-template",
	}
	results2, err := taskMgr.Query(filter2)
	if err != nil {
		log.Fatalf("Failed to query: %v", err)
	}
	fmt.Printf("  结果: 找到 %d 个任务\n", len(results2))
	fmt.Println()

	// 查询3: 按业务ID查询
	fmt.Println("查询3: 按业务ID查询(biz-001)")
	filter3 := &task.TaskFilter{
		BusinessID: "biz-001",
	}
	results3, err := taskMgr.Query(filter3)
	if err != nil {
		log.Fatalf("Failed to query: %v", err)
	}
	fmt.Printf("  结果: 找到 %d 个任务\n", len(results3))
	if len(results3) > 0 {
		tsk := results3[0]
		fmt.Printf("    - 任务 ID: %s, 状态: %s\n", tsk.ID, tsk.State)
		fmt.Printf("    - 审批记录数: %d\n", len(tsk.Records))
	}
	fmt.Println()

	// 查询4: 按审批人查询
	fmt.Println("查询4: 按审批人查询(manager-001)")
	filter4 := &task.TaskFilter{
		Approver: "manager-001",
	}
	results4, err := taskMgr.Query(filter4)
	if err != nil {
		log.Fatalf("Failed to query: %v", err)
	}
	fmt.Printf("  结果: 找到 %d 个任务\n", len(results4))
	for _, tsk := range results4 {
		fmt.Printf("    - 任务 ID: %s, 状态: %s\n", tsk.ID, tsk.State)
	}
	fmt.Println()

	// 查询5: 按时间范围查询
	fmt.Println("查询5: 按时间范围查询(最近1小时)")
	now := time.Now()
	oneHourAgo := now.Add(-1 * time.Hour)
	filter5 := &task.TaskFilter{
		StartTime: oneHourAgo,
		EndTime:   now,
	}
	results5, err := taskMgr.Query(filter5)
	if err != nil {
		log.Fatalf("Failed to query: %v", err)
	}
	fmt.Printf("  结果: 找到 %d 个任务\n", len(results5))
	fmt.Println()

	// 查询6: 综合查询
	fmt.Println("查询6: 综合查询(模板+状态)")
	filter6 := &task.TaskFilter{
		TemplateID: "record-query-template",
		State:      types.TaskStateApproving,
	}
	results6, err := taskMgr.Query(filter6)
	if err != nil {
		log.Fatalf("Failed to query: %v", err)
	}
	fmt.Printf("  结果: 找到 %d 个审批中的任务\n", len(results6))
	fmt.Println()

	// 审批记录分析
	fmt.Println("=== 审批记录分析 ===")
	fmt.Println()
	for i, tsk := range tasks {
		fmt.Printf("任务 %d: ID=%s, 状态=%s\n", i+1, tsk.ID, tsk.State)
		fmt.Printf("  审批记录数: %d\n", len(tsk.Records))
		if len(tsk.Records) > 0 {
			fmt.Println("  审批记录详情:")
			for j, record := range tsk.Records {
				fmt.Printf("    记录 %d:\n", j+1)
				fmt.Printf("      节点: %s\n", record.NodeID)
				fmt.Printf("      审批人: %s\n", record.Approver)
				fmt.Printf("      操作类型: %s\n", getOperationType(record.Result))
				fmt.Printf("      审批意见: %s\n", record.Comment)
				fmt.Printf("      操作时间: %s\n", record.CreatedAt.Format("2006-01-02 15:04:05"))
			}
		}
		fmt.Println()
	}

	// 统计信息
	fmt.Println("=== 统计信息 ===")
	fmt.Println()
	totalTasks := len(tasks)
	approvedTasks := 0
	approvingTasks := 0
	rejectedTasks := 0
	pendingTasks := 0
	totalRecords := 0

	for _, tsk := range tasks {
		switch tsk.State {
		case types.TaskStateApproved:
			approvedTasks++
		case types.TaskStateApproving:
			approvingTasks++
		case types.TaskStateRejected:
			rejectedTasks++
		case types.TaskStatePending, types.TaskStateSubmitted:
			pendingTasks++
		}
		totalRecords += len(tsk.Records)
	}

	fmt.Printf("  总任务数: %d\n", totalTasks)
	fmt.Printf("  已通过: %d\n", approvedTasks)
	fmt.Printf("  审批中: %d\n", approvingTasks)
	fmt.Printf("  已拒绝: %d\n", rejectedTasks)
	fmt.Printf("  待审批: %d\n", pendingTasks)
	fmt.Printf("  总审批记录数: %d\n", totalRecords)
	if totalTasks > 0 {
		fmt.Printf("  平均审批记录数: %.2f\n", float64(totalRecords)/float64(totalTasks))
	}
}

// getOperationType 获取操作类型描述
func getOperationType(result string) string {
	operationTypes := map[string]string{
		"approve":        "审批通过",
		"reject":         "审批拒绝",
		"transfer":       "转交审批",
		"add_approver":   "加签",
		"remove_approver": "减签",
	}
	if name, ok := operationTypes[result]; ok {
		return name
	}
	return result
}

