package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
)

func main() {
	fmt.Println("=== 场景 24: 并发安全场景 ===\n")

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

	// 运行并发安全测试场景
	runConcurrentScenarios(taskMgr, tpl.ID)
}

// createTemplate 创建多人会签审批模板
func createTemplate() *template.Template {
	now := time.Now()

	// 审批节点(多人会签)
	approvalNode := &template.Node{
		ID:   "approval",
		Name: "审批节点",
		Type: template.NodeTypeApproval,
		Config: &node.ApprovalNodeConfig{
			Mode: node.ApprovalModeUnanimous, // 多人会签模式
			ApproverConfig: &node.FixedApproverConfig{
				Approvers: []string{"approver-001", "approver-002", "approver-003"},
			},
			Permissions: node.OperationPermissions{
				AllowAddApprover: true,
			},
		},
		Order: 1,
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
		Order: 2,
	}

	return &template.Template{
		ID:          "concurrent-safety-template",
		Name:        "并发安全测试模板",
		Description: "测试多个审批人同时审批同一个任务的并发安全性",
		Version:     1,
		CreatedAt:   now,
		UpdatedAt:   now,
		Nodes: map[string]*template.Node{
			"start":    startNode,
			"approval": approvalNode,
			"end":      endNode,
		},
		Edges: []*template.Edge{
			{From: "start", To: "approval"},
			{From: "approval", To: "end"},
		},
	}
}

// runConcurrentScenarios 运行并发安全测试场景
func runConcurrentScenarios(taskMgr task.TaskManager, templateID string) {
	// 场景1: 并发读取测试
	fmt.Println("--- 场景1: 并发读取测试 ---")
	testConcurrentRead(taskMgr, templateID)
	fmt.Println()

	// 场景2: 并发更新测试
	fmt.Println("--- 场景2: 并发更新测试 ---")
	testConcurrentUpdate(taskMgr, templateID)
	fmt.Println()

	// 场景3: 并发审批测试
	fmt.Println("--- 场景3: 并发审批测试 ---")
	testConcurrentApproval(taskMgr, templateID)
	fmt.Println()

	// 场景4: 并发读取和更新混合测试
	fmt.Println("--- 场景4: 并发读取和更新混合测试 ---")
	testConcurrentReadAndUpdate(taskMgr, templateID)
	fmt.Println()
}

// testConcurrentRead 测试并发读取的安全性
func testConcurrentRead(taskMgr task.TaskManager, templateID string) {
	// 创建任务
	paramsJSON, _ := json.Marshal(map[string]interface{}{
		"test": "concurrent-read",
	})
	tsk, err := taskMgr.Create(templateID, "concurrent-read-001", paramsJSON)
	if err != nil {
		fmt.Printf("❌ 任务创建失败: %v\n", err)
		return
	}
	taskID := tsk.ID
	fmt.Printf("✓ 任务创建成功: ID=%s\n", taskID)

	// 提交任务
	if err := taskMgr.Submit(taskID); err != nil {
		fmt.Printf("❌ 任务提交失败: %v\n", err)
		return
	}
	fmt.Printf("✓ 任务已提交\n")

	// 并发读取
	concurrency := 50
	var wg sync.WaitGroup
	wg.Add(concurrency)

	startTime := time.Now()
	for i := 0; i < concurrency; i++ {
		go func(id int) {
			defer wg.Done()
			// 并发读取任务状态
			tsk, err := taskMgr.Get(taskID)
			if err != nil {
				fmt.Printf("❌ Goroutine %d: Get() 失败: %v\n", id, err)
				return
			}
			// 读取各种字段
			_ = tsk.GetState()
			_ = tsk.GetCurrentNode()
			_ = tsk.GetUpdatedAt()
			_ = tsk.GetStateHistory()
			_ = tsk.GetRecords()
			_ = tsk.Snapshot()
		}(i)
	}
	wg.Wait()
	duration := time.Since(startTime)

	fmt.Printf("✓ 并发读取测试完成: %d 个 goroutine, 耗时 %v\n", concurrency, duration)
	fmt.Printf("  说明: 使用读写锁保证并发读取的安全性,多个 goroutine 可以同时读取任务数据\n")
}

// testConcurrentUpdate 测试并发更新的安全性
func testConcurrentUpdate(taskMgr task.TaskManager, templateID string) {
	// 创建任务
	paramsJSON, _ := json.Marshal(map[string]interface{}{
		"test": "concurrent-update",
	})
	tsk, err := taskMgr.Create(templateID, "concurrent-update-001", paramsJSON)
	if err != nil {
		fmt.Printf("❌ 任务创建失败: %v\n", err)
		return
	}
	taskID := tsk.ID
	fmt.Printf("✓ 任务创建成功: ID=%s\n", taskID)

	// 提交任务
	if err := taskMgr.Submit(taskID); err != nil {
		fmt.Printf("❌ 任务提交失败: %v\n", err)
		return
	}
	fmt.Printf("✓ 任务已提交\n")

	// 并发更新
	concurrency := 30
	var wg sync.WaitGroup
	wg.Add(concurrency)

	startTime := time.Now()
	for i := 0; i < concurrency; i++ {
		go func(id int) {
			defer wg.Done()
			// 并发更新任务
			tsk, err := taskMgr.Get(taskID)
			if err != nil {
				fmt.Printf("❌ Goroutine %d: Get() 失败: %v\n", id, err)
				return
			}
			err = tsk.Update(func(t *task.Task) error {
				t.UpdatedAt = time.Now()
				return nil
			})
			if err != nil {
				fmt.Printf("❌ Goroutine %d: Update() 失败: %v\n", id, err)
			}
		}(i)
	}
	wg.Wait()
	duration := time.Since(startTime)

	// 验证任务状态一致性
	tsk, _ = taskMgr.Get(taskID)
	fmt.Printf("✓ 并发更新测试完成: %d 个 goroutine, 耗时 %v\n", concurrency, duration)
	fmt.Printf("  最终状态: %s, 更新时间: %s\n", tsk.GetState(), tsk.GetUpdatedAt().Format("2006-01-02 15:04:05"))
	fmt.Printf("  说明: 使用写锁保证并发更新的原子性,更新操作是串行化的\n")
}

// testConcurrentApproval 测试并发审批的安全性
func testConcurrentApproval(taskMgr task.TaskManager, templateID string) {
	// 创建任务
	paramsJSON, _ := json.Marshal(map[string]interface{}{
		"test": "concurrent-approval",
	})
	tsk, err := taskMgr.Create(templateID, "concurrent-approval-001", paramsJSON)
	if err != nil {
		fmt.Printf("❌ 任务创建失败: %v\n", err)
		return
	}
	taskID := tsk.ID
	fmt.Printf("✓ 任务创建成功: ID=%s\n", taskID)

	// 提交任务
	if err := taskMgr.Submit(taskID); err != nil {
		fmt.Printf("❌ 任务提交失败: %v\n", err)
		return
	}
	fmt.Printf("✓ 任务已提交\n")

	// 设置审批人
	approvers := []string{"approver-001", "approver-002", "approver-003"}
	for _, approver := range approvers {
		if err := taskMgr.AddApprover(taskID, "approval", approver, "并发测试添加审批人"); err != nil {
			fmt.Printf("❌ 设置审批人失败: %v\n", err)
			return
		}
	}
	fmt.Printf("✓ 审批人已设置: %v\n", approvers)

	// 并发审批
	var wg sync.WaitGroup
	wg.Add(len(approvers))

	startTime := time.Now()
	for i, approver := range approvers {
		go func(id int, apr string) {
			defer wg.Done()
			// 并发审批
			err := taskMgr.Approve(taskID, "approval", apr, fmt.Sprintf("审批意见 from %s", apr))
			if err != nil {
				fmt.Printf("⚠ Goroutine %d: Approve() 失败: %v (可能是状态已变更)\n", id, err)
			} else {
				fmt.Printf("✓ Goroutine %d: %s 审批成功\n", id, apr)
			}
		}(i, approver)
	}
	wg.Wait()
	duration := time.Since(startTime)

	// 验证任务状态
	tsk, _ = taskMgr.Get(taskID)
	records := tsk.GetRecords()
	approveCount := 0
	for _, record := range records {
		if record.Result == "approve" {
			approveCount++
		}
	}

	fmt.Printf("✓ 并发审批测试完成: %d 个审批人, 耗时 %v\n", len(approvers), duration)
	fmt.Printf("  最终状态: %s\n", tsk.GetState())
	fmt.Printf("  审批记录数: %d (期望: %d)\n", approveCount, len(approvers))
	fmt.Printf("  说明: 使用写锁保证审批操作的原子性,状态转换是串行化的\n")
}

// testConcurrentReadAndUpdate 测试并发读取和更新混合场景
func testConcurrentReadAndUpdate(taskMgr task.TaskManager, templateID string) {
	// 创建任务
	paramsJSON, _ := json.Marshal(map[string]interface{}{
		"test": "concurrent-read-update",
	})
	tsk, err := taskMgr.Create(templateID, "concurrent-read-update-001", paramsJSON)
	if err != nil {
		fmt.Printf("❌ 任务创建失败: %v\n", err)
		return
	}
	taskID := tsk.ID
	fmt.Printf("✓ 任务创建成功: ID=%s\n", taskID)

	// 提交任务
	if err := taskMgr.Submit(taskID); err != nil {
		fmt.Printf("❌ 任务提交失败: %v\n", err)
		return
	}
	fmt.Printf("✓ 任务已提交\n")

	// 并发读取和更新混合
	readConcurrency := 30
	updateConcurrency := 10
	var wg sync.WaitGroup
	wg.Add(readConcurrency + updateConcurrency)

	startTime := time.Now()

	// 启动读取 goroutine
	for i := 0; i < readConcurrency; i++ {
		go func(id int) {
			defer wg.Done()
			tsk, err := taskMgr.Get(taskID)
			if err != nil {
				return
			}
			_ = tsk.GetState()
			_ = tsk.GetCurrentNode()
			_ = tsk.Snapshot()
		}(i)
	}

	// 启动更新 goroutine
	for i := 0; i < updateConcurrency; i++ {
		go func(id int) {
			defer wg.Done()
			tsk, err := taskMgr.Get(taskID)
			if err != nil {
				return
			}
			_ = tsk.Update(func(t *task.Task) error {
				t.UpdatedAt = time.Now()
				return nil
			})
		}(i)
	}

	wg.Wait()
	duration := time.Since(startTime)

	// 验证任务状态一致性
	tsk, _ = taskMgr.Get(taskID)
	fmt.Printf("✓ 并发读取和更新混合测试完成: %d 个读取 + %d 个更新, 耗时 %v\n",
		readConcurrency, updateConcurrency, duration)
	fmt.Printf("  最终状态: %s\n", tsk.GetState())
	fmt.Printf("  说明: 读写锁支持并发读取,写操作会阻塞所有读写操作,保证数据一致性\n")
}

