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
)

func main() {
	fmt.Println("=== 场景 21: 任务创建时获取动态审批人场景 ===")
	fmt.Println()

	// 1. 创建 Mock HTTP 客户端
	fmt.Println("步骤 1: 创建 Mock HTTP 客户端")
	mockHTTPClient := createMockHTTPClient()
	fmt.Printf("✓ Mock HTTP 客户端已创建\n\n")

	// 2. 创建管理器
	fmt.Println("步骤 2: 创建管理器")
	templateMgr := template.NewTemplateManager()
	approverFetcherFunc := func(tpl *template.Template, tsk *task.Task) error {
		return node.FetchApproversOnCreate(tpl, tsk, mockHTTPClient)
	}
	taskMgr := task.NewTaskManager(templateMgr, approverFetcherFunc)
	fmt.Printf("✓ 管理器已创建(带审批人获取函数)\n\n")

	// 3. 创建模板
	fmt.Println("步骤 3: 创建审批模板")
	tpl := createTemplate()
	err := templateMgr.Create(tpl)
	if err != nil {
		log.Fatalf("Failed to create template: %v", err)
	}
	fmt.Printf("✓ 模板创建成功: ID=%s, Name=%s\n\n", tpl.ID, tpl.Name)

	// 4. 创建任务
	fmt.Println("步骤 4: 创建项目立项审批任务")
	params := json.RawMessage(`{
		"projectNo": "PRJ-2025-001",
		"projectName": "新产品开发项目",
		"projectManager": "pm-001",
		"projectType": "研发项目",
		"budget": 2000000,
		"description": "新产品开发项目立项审批"
	}`)
	tsk, err := taskMgr.Create(tpl.ID, "project-001", params)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}
	fmt.Printf("✓ 任务创建成功: ID=%s, State=%s\n\n", tsk.ID, tsk.State)

	// 5. 验证审批人是否已获取
	fmt.Println("步骤 5: 验证审批人是否已获取")
	fmt.Println("  说明: 任务创建时,配置为 on_create 时机的动态审批人应该已自动获取")
	fmt.Println()
	nodeID := "project-approval"
	approvers := tsk.Approvers[nodeID]
	if len(approvers) > 0 {
		fmt.Printf("✓ 审批人已自动获取: %v (共 %d 人)\n", approvers, len(approvers))
		fmt.Println("  说明: 审批人在任务创建时已自动获取,无需等待节点激活")
	} else {
		fmt.Println("✗ 审批人未获取")
	}
	fmt.Println()

	// 6. 输出结果
	printResults(taskMgr, tsk)
}

// createMockHTTPClient 创建 Mock HTTP 客户端
func createMockHTTPClient() node.HTTPClient {
	return &mockHTTPClient{
		response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`{"approvers": ["tech-lead-001", "tech-manager-001"]}`)),
		},
		err: nil,
	}
}

// mockHTTPClient Mock HTTP 客户端实现
type mockHTTPClient struct {
	response *http.Response
	err      error
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	// 模拟 API 调用
	// 根据项目类型返回不同的审批人
	body, _ := io.ReadAll(req.Body)
	bodyStr := string(body)

	if strings.Contains(bodyStr, "研发项目") {
		// 研发项目返回技术负责人和技术经理
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`{"approvers": ["tech-lead-001", "tech-manager-001"]}`)),
		}, nil
	}

	// 默认返回
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{"approvers": ["manager-001"]}`)),
	}, nil
}

// createTemplate 创建项目立项审批模板
// 流程: 开始节点 → 审批节点(任务创建时获取审批人) → 结束节点
// 审批节点配置为任务创建时获取动态审批人
func createTemplate() *template.Template {
	now := time.Now()
	return &template.Template{
		ID:          "approver-on-create-template",
		Name:        "项目立项审批模板",
		Description: "项目立项审批流程,任务创建时获取动态审批人",
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
					Mode: node.ApprovalModeUnanimous, // 多人会签模式
					ApproverConfig: &node.DynamicApproverConfig{
						API: &node.HTTPAPIConfig{
							URL:    "http://example.com/api/approvers",
							Method: "POST",
							ParamMapping: &node.ParamMapping{
								Source: "task_params",
								Path:   "projectType",
								Target: "projectType",
							},
							ResponseMapping: &node.ResponseMapping{
								Path:   "approvers",
								Format: "json",
							},
						},
						Timing: node.ApproverTimingOnCreate, // 任务创建时获取
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
	nodeID := "project-approval"
	approvers := tsk.Approvers[nodeID]
	if len(approvers) > 0 {
		fmt.Printf("  节点: %s\n", getNodeName(nodeID))
		for _, approver := range approvers {
			fmt.Printf("    - %s\n", approver)
		}
		fmt.Printf("  说明: 审批人在任务创建时已自动获取(ApproverTimingOnCreate)\n")
	} else {
		fmt.Println("  无审批人")
	}

	// 输出审批记录
	fmt.Println("\n审批记录:")
	if len(tsk.Records) == 0 {
		fmt.Println("  无审批记录")
	} else {
		for i, record := range tsk.Records {
			fmt.Printf("  记录 %d:\n", i+1)
			fmt.Printf("    节点: %s\n", getNodeName(record.NodeID))
			fmt.Printf("    审批人: %s\n", record.Approver)
			fmt.Printf("    操作类型: %s\n", getOperationType(record.Result))
			fmt.Printf("    操作说明: %s\n", record.Comment)
			fmt.Printf("    操作时间: %s\n", record.CreatedAt.Format("2006-01-02 15:04:05"))
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
	if len(approvers) > 0 {
		fmt.Println("✓ 审批人已在任务创建时自动获取:")
		fmt.Printf("  - 审批人数量: %d\n", len(approvers))
		fmt.Printf("  - 审批人列表: %v\n", approvers)
		fmt.Println("  - 获取时机: on_create (任务创建时)")
	} else {
		fmt.Println("✗ 审批人未在任务创建时获取")
	}

	fmt.Println()
	fmt.Println("✓ 任务创建时获取动态审批人功能配置和使用方法已展示")
	fmt.Println("  说明: 在实际业务系统中,审批人获取函数会在任务创建时自动调用")
	fmt.Println("  配置为 on_create 时机的动态审批人会在任务创建时立即获取")
}

// getNodeName 获取节点名称
func getNodeName(nodeID string) string {
	nodeNames := map[string]string{
		"project-approval": "项目审批",
	}
	if name, ok := nodeNames[nodeID]; ok {
		return name
	}
	return nodeID
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

