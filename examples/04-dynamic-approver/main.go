package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
	"github.com/mautops/approval-kit/internal/types"
)

// mockHTTPClient Mock HTTP 客户端,用于模拟 API 调用
type mockHTTPClient struct {
	response *http.Response
	err      error
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.response, m.err
}

func main() {
	fmt.Println("=== 场景 04: 动态审批人获取 ===")
	fmt.Println()

	// 1. 创建管理器
	templateMgr := template.NewTemplateManager()
	taskMgr := task.NewTaskManager(templateMgr, nil)

	// 2. 创建 Mock HTTP 客户端
	fmt.Println("步骤 1: 创建 Mock HTTP 客户端")
	mockClient := createMockHTTPClient()
	fmt.Println("✓ Mock HTTP 客户端已创建")
	fmt.Println()

	// 3. 创建模板
	fmt.Println("步骤 2: 创建审批模板(配置动态审批人)")
	tpl := createTemplate(mockClient)
	err := templateMgr.Create(tpl)
	if err != nil {
		log.Fatalf("Failed to create template: %v", err)
	}
	fmt.Printf("✓ 模板创建成功: ID=%s, Name=%s\n\n", tpl.ID, tpl.Name)

	// 4. 创建任务
	fmt.Println("步骤 3: 创建项目立项任务")
	params := json.RawMessage(`{
		"projectType": "研发项目",
		"department": "技术部",
		"budget": 100000,
		"description": "新产品研发项目立项"
	}`)
	tsk, err := taskMgr.Create(tpl.ID, "project-001", params)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}
	fmt.Printf("✓ 任务创建成功: ID=%s, State=%s\n\n", tsk.ID, tsk.State)

	// 5. 执行场景流程
	runScenario(templateMgr, taskMgr, tsk, mockClient)

	// 6. 输出结果
	printResults(taskMgr, tsk)
}

// createMockHTTPClient 创建 Mock HTTP 客户端
// 模拟根据项目类型和部门返回审批人列表的 API
func createMockHTTPClient() *mockHTTPClient {
	// 模拟 API 响应: 根据项目类型和部门返回不同的审批人
	responseBody := `{
		"code": 200,
		"message": "success",
		"data": {
			"approvers": ["tech-lead-001", "tech-manager-001"]
		}
	}`

	return &mockHTTPClient{
		response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(responseBody)),
			Header:     make(http.Header),
		},
		err: nil,
	}
}

// createTemplate 创建项目立项审批模板
// 流程: 开始节点 → 审批节点(动态审批人) → 结束节点
func createTemplate(httpClient node.HTTPClient) *template.Template {
	now := time.Now()
	return &template.Template{
		ID:          "project-approval-template",
		Name:        "项目立项审批模板",
		Description: "项目立项审批,审批人根据项目类型和部门动态获取",
		Version:     1,
		CreatedAt:   now,
		UpdatedAt:   now,
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "开始",
				Type: template.NodeTypeStart,
			},
			"project-approval": {
				ID:   "project-approval",
				Name: "项目审批",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeSingle,
					ApproverConfig: &node.DynamicApproverConfig{
						API: &node.HTTPAPIConfig{
							URL:    "http://api.example.com/approvers",
							Method: "POST",
							Headers: map[string]string{
								"Content-Type": "application/json",
								"Authorization": "Bearer token-123",
							},
							// 参数映射: 从任务参数中获取项目类型和部门
							ParamMapping: &node.ParamMapping{
								Source: "task_params",
								Path:   "projectType", // 从任务参数中获取 projectType 字段
								Target: "project_type", // 映射到 API 请求参数 project_type
							},
							// 响应解析: 从 API 响应中解析审批人列表
							ResponseMapping: &node.ResponseMapping{
								Path:   "data.approvers", // 响应路径: data.approvers
								Format: "json",
							},
						},
						Timing:     node.ApproverTimingOnActivate, // 节点激活时获取
						HTTPClient: httpClient,
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
			{From: "start", To: "project-approval"},
			{From: "project-approval", To: "end"},
		},
	}
}

// runScenario 执行场景流程
func runScenario(templateMgr template.TemplateManager, taskMgr task.TaskManager, tsk *task.Task, httpClient node.HTTPClient) {
	// 步骤 4: 提交任务
	fmt.Println("步骤 4: 提交任务进入审批流程")
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

	// 步骤 5: 动态获取审批人(节点激活时)
	fmt.Println("步骤 5: 动态获取审批人(节点激活时)")
	fmt.Println("  说明: 系统会调用 HTTP API,根据项目类型和部门获取审批人列表")
	fmt.Println()

	// 获取模板和节点配置
	tpl, err := templateMgr.Get(tsk.TemplateID, 0)
	if err != nil {
		log.Fatalf("Failed to get template: %v", err)
	}

	tplNode := tpl.Nodes["project-approval"]
	config, ok := tplNode.Config.(*node.ApprovalNodeConfig)
	if !ok {
		log.Fatalf("Failed to cast node config to ApprovalNodeConfig")
	}
	dynamicConfig, ok := config.ApproverConfig.(*node.DynamicApproverConfig)
	if !ok {
		log.Fatalf("Failed to cast approver config to DynamicApproverConfig")
	}

	// 创建节点上下文
	ctx := &node.NodeContext{
		Task:    tsk,
		Node:    tplNode,
		Params:  tsk.Params,
		Outputs: make(map[string]json.RawMessage),
		Cache:   node.NewContextCache(),
	}

	// 调用 API 获取审批人
	fmt.Println("  调用 API: POST http://api.example.com/approvers")
	fmt.Println("  请求参数: {\"project_type\": \"研发项目\"}")
	approvers, err := dynamicConfig.GetApprovers(ctx)
	if err != nil {
		log.Fatalf("Failed to get approvers: %v", err)
	}

	fmt.Printf("  ✓ API 调用成功,获取到审批人: %v\n", approvers)
	fmt.Println()

	// 设置审批人列表(单人审批模式,只使用第一个审批人)
	if len(approvers) > 0 {
		approver := approvers[0]
		err = taskMgr.AddApprover(tsk.ID, "project-approval", approver, "动态获取的审批人")
		if err != nil {
			log.Fatalf("Failed to add approver: %v", err)
		}
		fmt.Printf("✓ 审批人列表已设置: %s (单人审批模式,使用第一个审批人)\n", approver)
		if len(approvers) > 1 {
			fmt.Printf("  注意: API 返回了 %d 个审批人,但当前节点配置为单人审批模式\n", len(approvers))
		}
		fmt.Println()
	}

	// 步骤 6: 执行审批操作
	fmt.Println("步骤 6: 执行审批操作")
	fmt.Println("  说明: 使用第一个审批人进行审批(单人审批模式)")
	fmt.Println()
	if len(approvers) > 0 {
		// 只使用第一个审批人(单人审批模式)
		approver := approvers[0]
		fmt.Printf("  使用审批人: %s\n", approver)
		err = taskMgr.Approve(tsk.ID, "project-approval", approver, "项目立项审批通过")
		if err != nil {
			log.Fatalf("Failed to approve: %v", err)
		}

		tsk, _ = taskMgr.Get(tsk.ID)
		fmt.Printf("  ✓ %s 已同意 (当前状态: %s)\n", approver, tsk.State)
		if len(approvers) > 1 {
			fmt.Printf("  注意: API 返回了 %d 个审批人,但当前节点配置为单人审批模式,只使用第一个审批人\n", len(approvers))
		}
	}
	fmt.Println()
	fmt.Println("✓ 审批流程完成")
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

	// 输出审批人列表
	fmt.Println("\n审批人列表:")
	approvers := tsk.Approvers["project-approval"]
	if len(approvers) > 0 {
		for i, approver := range approvers {
			fmt.Printf("  %d. %s (动态获取)\n", i+1, approver)
		}
	} else {
		fmt.Println("  无审批人")
	}

	// 输出审批记录
	fmt.Println("\n审批记录:")
	if len(tsk.Records) == 0 {
		fmt.Println("  无审批记录")
	} else {
		approveCount := 0
		for _, record := range tsk.Records {
			if record.Result == "approve" {
				approveCount++
				fmt.Printf("  记录 %d (审批):\n", approveCount)
				fmt.Printf("    节点 ID: %s\n", record.NodeID)
				fmt.Printf("    审批人: %s\n", record.Approver)
				fmt.Printf("    审批结果: %s\n", record.Result)
				fmt.Printf("    审批意见: %s\n", record.Comment)
				fmt.Printf("    审批时间: %s\n", record.CreatedAt.Format("2006-01-02 15:04:05"))
			}
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

	// 验证最终状态
	fmt.Println("\n=== 验证结果 ===")
	if tsk.State == types.TaskStateApproved {
		fmt.Println("✓ 任务已成功通过审批")
	} else {
		fmt.Printf("✗ 任务状态异常: 期望 %s, 实际 %s\n", types.TaskStateApproved, tsk.State)
	}

	if len(approvers) > 0 {
		fmt.Printf("✓ 已通过 API 动态获取 %d 个审批人\n", len(approvers))
	} else {
		fmt.Println("✗ 未获取到审批人")
	}
}

