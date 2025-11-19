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
	fmt.Println("=== 场景 10: 复杂综合审批流程 ===")
	fmt.Println()

	// 1. 创建 Mock HTTP 客户端
	mockHTTPClient := createMockHTTPClient()

	// 2. 创建管理器
	templateMgr := template.NewTemplateManager()
	taskMgr := task.NewTaskManager(templateMgr, nil)

	// 3. 创建模板
	fmt.Println("步骤 1: 创建审批模板")
	tpl := createTemplate(mockHTTPClient)
	err := templateMgr.Create(tpl)
	if err != nil {
		log.Fatalf("Failed to create template: %v", err)
	}
	fmt.Printf("✓ 模板创建成功: ID=%s, Name=%s\n\n", tpl.ID, tpl.Name)

	// 4. 运行多个场景
	scenarios := []struct {
		name        string
		projectType string
		expectPath  string
	}{
		{"研发项目", "研发项目", "技术负责人(动态获取) → CTO(顺序审批)"},
		{"市场项目", "市场项目", "市场总监(单人审批) → 财务(会签)"},
		{"其他项目", "其他项目", "部门经理(或签) → 总经理(比例会签)"},
	}

	for _, scenario := range scenarios {
		fmt.Printf("--- %s场景 ---\n", scenario.name)
		runScenario(taskMgr, tpl.ID, scenario.projectType, scenario.expectPath)
		fmt.Println()
	}
}

// createMockHTTPClient 创建 Mock HTTP 客户端
func createMockHTTPClient() node.HTTPClient {
	return &MockHTTPClient{}
}

// MockHTTPClient Mock HTTP 客户端实现
type MockHTTPClient struct{}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	// 模拟 API 响应: 根据项目类型返回不同的技术负责人
	var approvers []string
	if req.Method == "POST" {
		// 读取请求体
		body, _ := io.ReadAll(req.Body)
		req.Body.Close()

		var requestData map[string]interface{}
		if err := json.Unmarshal(body, &requestData); err == nil {
			if projectType, ok := requestData["projectType"].(string); ok {
				if projectType == "研发项目" {
					// 研发项目返回技术负责人
					approvers = []string{"tech-lead-001", "tech-manager-001"}
				}
			}
		}
	}

	// 构建响应
	responseData := map[string]interface{}{
		"approvers": approvers,
	}
	responseJSON, _ := json.Marshal(responseData)
	responseBody := strings.NewReader(string(responseJSON))

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(responseBody),
		Header:     make(http.Header),
	}, nil
}

// createTemplate 创建复杂综合审批模板
// 流程: 开始节点 → 条件节点1(项目类型判断) → [研发项目路径/条件节点2] → 条件节点2(市场项目判断) → [市场项目路径/其他项目路径] → 结束节点
// 研发项目路径: 技术负责人(动态获取) → CTO(顺序审批)
// 市场项目路径: 市场总监(单人审批) → 财务(会签)
// 其他项目路径: 部门经理(或签) → 总经理(比例会签)
func createTemplate(httpClient node.HTTPClient) *template.Template {
	now := time.Now()

	// 条件节点1: 判断是否为研发项目
	conditionNode1 := &template.Node{
		ID:   "condition-node-1",
		Name: "条件判断1(研发项目)",
		Type: template.NodeTypeCondition,
		Config: &node.ConditionNodeConfig{
			Condition: &node.Condition{
				Type: "enum",
				Config: &node.EnumConditionConfig{
					Field:    "projectType",
					Operator: "in",
					Values:   []string{"研发项目"},
					Source:   "task_params",
				},
			},
			TrueNodeID:  "tech-lead-approval", // 研发项目: 技术负责人审批
			FalseNodeID: "condition-node-2",   // 非研发项目: 进入条件节点2
		},
		Order: 1,
	}

	// 条件节点2: 判断是否为市场项目
	conditionNode2 := &template.Node{
		ID:   "condition-node-2",
		Name: "条件判断2(市场项目)",
		Type: template.NodeTypeCondition,
		Config: &node.ConditionNodeConfig{
			Condition: &node.Condition{
				Type: "enum",
				Config: &node.EnumConditionConfig{
					Field:    "projectType",
					Operator: "in",
					Values:   []string{"市场项目"},
					Source:   "task_params",
				},
			},
			TrueNodeID:  "market-director-approval", // 市场项目: 市场总监审批
			FalseNodeID: "dept-manager-approval",    // 其他项目: 部门经理审批
		},
		Order: 2,
	}

	// 超时时间配置
	timeout24h := 24 * time.Hour
	timeout48h := 48 * time.Hour

	// 研发项目路径: 技术负责人审批(动态获取)
	techLeadApproval := &template.Node{
		ID:   "tech-lead-approval",
		Name: "技术负责人审批",
		Type: template.NodeTypeApproval,
		Config: &node.ApprovalNodeConfig{
			Mode: node.ApprovalModeSingle,
			ApproverConfig: &node.DynamicApproverConfig{
				API: &node.HTTPAPIConfig{
					URL:    "http://api.example.com/approvers",
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
				Timing:     node.ApproverTimingOnActivate,
				HTTPClient: httpClient,
			},
			Timeout:        &timeout24h, // 超时配置: 24小时
			RejectBehavior: node.RejectBehaviorRollback, // 拒绝后回退
			Permissions: node.OperationPermissions{
				AllowAddApprover: true,
			},
		},
		Order: 3,
	}

	// 研发项目路径: CTO审批(顺序审批)
	ctoApproval := &template.Node{
		ID:   "cto-approval",
		Name: "CTO审批",
		Type: template.NodeTypeApproval,
		Config: &node.ApprovalNodeConfig{
			Mode: node.ApprovalModeSequential,
			ApproverConfig: &node.FixedApproverConfig{
				Approvers: []string{"cto-001"},
			},
			Timeout: &timeout48h, // 超时配置: 48小时
			Permissions: node.OperationPermissions{
				AllowAddApprover: true,
			},
		},
		Order: 4,
	}

	// 市场项目路径: 市场总监审批(单人审批)
	marketDirectorApproval := &template.Node{
		ID:   "market-director-approval",
		Name: "市场总监审批",
		Type: template.NodeTypeApproval,
		Config: &node.ApprovalNodeConfig{
			Mode: node.ApprovalModeSingle,
			ApproverConfig: &node.FixedApproverConfig{
				Approvers: []string{"market-director-001"},
			},
			Timeout: &timeout24h,
			Permissions: node.OperationPermissions{
				AllowAddApprover: true,
			},
		},
		Order: 5,
	}

	// 市场项目路径: 财务审批(会签)
	financeApproval := &template.Node{
		ID:   "finance-approval",
		Name: "财务审批",
		Type: template.NodeTypeApproval,
		Config: &node.ApprovalNodeConfig{
			Mode: node.ApprovalModeUnanimous,
			ApproverConfig: &node.FixedApproverConfig{
				Approvers: []string{"finance-001", "finance-002", "finance-003"},
			},
			Timeout: &timeout24h,
			Permissions: node.OperationPermissions{
				AllowAddApprover: true,
			},
		},
		Order: 6,
	}

	// 其他项目路径: 部门经理审批(或签)
	deptManagerApproval := &template.Node{
		ID:   "dept-manager-approval",
		Name: "部门经理审批",
		Type: template.NodeTypeApproval,
		Config: &node.ApprovalNodeConfig{
			Mode: node.ApprovalModeOr,
			ApproverConfig: &node.FixedApproverConfig{
				Approvers: []string{"dept-manager-001", "dept-manager-002"},
			},
			Timeout: &timeout24h,
			Permissions: node.OperationPermissions{
				AllowAddApprover: true,
			},
		},
		Order: 7,
	}

	// 其他项目路径: 总经理审批(比例会签)
	gmApproval := &template.Node{
		ID:   "gm-approval",
		Name: "总经理审批",
		Type: template.NodeTypeApproval,
		Config: &node.ApprovalNodeConfig{
			Mode: node.ApprovalModeProportional,
			ApproverConfig: &node.FixedApproverConfig{
				Approvers: []string{"gm-001", "gm-002", "gm-003", "gm-004", "gm-005"},
			},
			ProportionalThreshold: &node.ProportionalThreshold{
				Required: 3, // 5 人中需要 3 人同意
				Total:    5,
			},
			Timeout: &timeout48h,
			Permissions: node.OperationPermissions{
				AllowAddApprover: true,
			},
		},
		Order: 8,
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
		Order: 9,
	}

	return &template.Template{
		ID:          "complex-approval-template",
		Name:        "复杂综合审批模板",
		Description: "大型项目立项审批,包含多种审批模式、条件分支、动态审批人、超时处理等",
		Version:     1,
		CreatedAt:   now,
		UpdatedAt:   now,
		Nodes: map[string]*template.Node{
			"start":                    startNode,
			"condition-node-1":         conditionNode1,
			"condition-node-2":         conditionNode2,
			"tech-lead-approval":       techLeadApproval,
			"cto-approval":             ctoApproval,
			"market-director-approval":  marketDirectorApproval,
			"finance-approval":         financeApproval,
			"dept-manager-approval":    deptManagerApproval,
			"gm-approval":              gmApproval,
			"end":                      endNode,
		},
		Edges: []*template.Edge{
			{From: "start", To: "condition-node-1"},
			{From: "condition-node-1", To: "tech-lead-approval", Condition: "projectType == 研发项目"},
			{From: "condition-node-1", To: "condition-node-2", Condition: "projectType != 研发项目"},
			{From: "condition-node-2", To: "market-director-approval", Condition: "projectType == 市场项目"},
			{From: "condition-node-2", To: "dept-manager-approval", Condition: "projectType != 市场项目"},
			{From: "tech-lead-approval", To: "cto-approval"},
			{From: "cto-approval", To: "end"},
			{From: "market-director-approval", To: "finance-approval"},
			{From: "finance-approval", To: "end"},
			{From: "dept-manager-approval", To: "gm-approval"},
			{From: "gm-approval", To: "end"},
		},
		Config: &template.TemplateConfig{
			Webhooks: []*template.WebhookConfig{
				{
					URL:     "http://webhook.example.com/events",
					Method:  "POST",
					Headers: map[string]string{"Authorization": "Bearer token-123"},
				},
			},
		},
	}
}

// runScenario 执行场景流程
func runScenario(taskMgr task.TaskManager, templateID string, projectType string, expectPath string) {
	// 1. 创建任务
	taskParams := map[string]interface{}{
		"projectNo":   fmt.Sprintf("PRJ-2025-%s", projectType),
		"projectName": fmt.Sprintf("%s项目", projectType),
		"projectType": projectType,
		"requester":   "申请人-001",
		"amount":      2000000,
		"description": fmt.Sprintf("%s项目立项审批", projectType),
	}
	paramsJSON, _ := json.Marshal(taskParams)

	tsk, err := taskMgr.Create(templateID, fmt.Sprintf("project-%s-001", projectType), paramsJSON)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}
	taskID := tsk.ID
	fmt.Printf("✓ 任务创建成功: ID=%s, 项目类型=%s\n", taskID, projectType)

	// 2. 提交任务
	if err := taskMgr.Submit(taskID); err != nil {
		log.Fatalf("Failed to submit task: %v", err)
	}
	fmt.Printf("✓ 任务已提交\n")

	// 3. 说明审批流程
	fmt.Println("\n=== 审批流程说明 ===")
	fmt.Printf("项目类型: %s\n", projectType)
	fmt.Printf("预期路径: %s\n", expectPath)
	fmt.Println()

	// 根据项目类型说明不同的审批路径
	switch projectType {
	case "研发项目":
		fmt.Println("研发项目审批路径:")
		fmt.Println("  1. 条件节点1判断: 项目类型 == 研发项目 → True")
		fmt.Println("  2. 技术负责人审批(动态获取): 通过 HTTP API 获取技术负责人")
		fmt.Println("  3. CTO审批(顺序审批): 技术负责人同意后,CTO按顺序审批")
		fmt.Println("  4. 流程结束")
		fmt.Println()
		fmt.Println("关键特性:")
		fmt.Println("  - 动态审批人: 技术负责人根据项目领域动态获取")
		fmt.Println("  - 顺序审批: CTO审批必须等技术负责人同意后")
		fmt.Println("  - 超时处理: 技术负责人审批24小时,CTO审批48小时")
		fmt.Println("  - 拒绝后回退: 技术负责人拒绝后回退到上一节点")

	case "市场项目":
		fmt.Println("市场项目审批路径:")
		fmt.Println("  1. 条件节点1判断: 项目类型 == 研发项目 → False")
		fmt.Println("  2. 条件节点2判断: 项目类型 == 市场项目 → True")
		fmt.Println("  3. 市场总监审批(单人审批): 市场总监审批")
		fmt.Println("  4. 财务审批(会签): 需要所有财务审批人全部同意")
		fmt.Println("  5. 流程结束")
		fmt.Println()
		fmt.Println("关键特性:")
		fmt.Println("  - 单人审批: 市场总监单人审批")
		fmt.Println("  - 会签模式: 财务审批需要多人全部同意")
		fmt.Println("  - 超时处理: 市场总监审批24小时,财务审批24小时")

	case "其他项目":
		fmt.Println("其他项目审批路径:")
		fmt.Println("  1. 条件节点1判断: 项目类型 == 研发项目 → False")
		fmt.Println("  2. 条件节点2判断: 项目类型 == 市场项目 → False")
		fmt.Println("  3. 部门经理审批(或签): 任意一个部门经理同意即可")
		fmt.Println("  4. 总经理审批(比例会签): 5人中需要3人同意")
		fmt.Println("  5. 流程结束")
		fmt.Println()
		fmt.Println("关键特性:")
		fmt.Println("  - 或签模式: 部门经理审批任意一人同意即可")
		fmt.Println("  - 比例会签: 总经理审批达到比例(3/5)即可")
		fmt.Println("  - 超时处理: 部门经理审批24小时,总经理审批48小时")
	}

	fmt.Println("\n注意: 本示例主要展示复杂综合审批流程的配置方式.")
	fmt.Println("在实际业务系统中,工作流引擎会根据条件自动选择审批路径.")
	fmt.Println("并根据审批模式自动处理审批流程.")

	// 输出结果
	printResults(taskMgr, taskID)
}

// printResults 输出结果
func printResults(taskMgr task.TaskManager, taskID string) {
	tsk, err := taskMgr.Get(taskID)
	if err != nil {
		log.Fatalf("Failed to get task: %v", err)
	}

	fmt.Println("\n=== 审批结果 ===")
	fmt.Printf("任务 ID: %s\n", tsk.ID)
	fmt.Printf("业务 ID: %s\n", tsk.BusinessID)
	fmt.Printf("模板 ID: %s\n", tsk.TemplateID)
	fmt.Printf("当前状态: %s\n", tsk.GetState())
	fmt.Printf("当前节点: %s\n", tsk.GetCurrentNode())

	// 输出任务参数
	fmt.Println("\n任务参数:")
	var params map[string]interface{}
	if err := json.Unmarshal(tsk.Params, &params); err == nil {
		for k, v := range params {
			fmt.Printf("  %s: %v\n", k, v)
		}
	}

	// 输出审批记录
	fmt.Println("\n审批记录:")
	records := tsk.GetRecords()
	if len(records) == 0 {
		fmt.Println("  无审批记录")
	} else {
		for i, record := range records {
			fmt.Printf("  记录 %d: 节点=%s, 审批人=%s, 操作=%s\n",
				i+1, record.NodeID, record.Approver, record.Result)
		}
	}

	// 验证配置
	fmt.Println("\n=== 验证结果 ===")
	fmt.Println("✓ 复杂综合审批流程配置已正确设置:")
	fmt.Println("  - 条件节点1(研发项目判断)已配置")
	fmt.Println("  - 条件节点2(市场项目判断)已配置")
	fmt.Println("  - 技术负责人审批(动态获取)已配置")
	fmt.Println("  - CTO审批(顺序审批)已配置")
	fmt.Println("  - 市场总监审批(单人审批)已配置")
	fmt.Println("  - 财务审批(会签)已配置")
	fmt.Println("  - 部门经理审批(或签)已配置")
	fmt.Println("  - 总经理审批(比例会签)已配置")
	fmt.Println("  - 超时处理配置已设置")
	fmt.Println("  - 拒绝后跳转配置已设置")
	fmt.Println("  - Webhook事件通知配置已设置")
	fmt.Println()
	fmt.Println("  说明: 本场景综合展示了多种审批模式和功能")
	fmt.Println("  在实际业务系统中,流程引擎会根据条件自动选择审批路径")
	fmt.Println("  并根据审批模式自动处理审批流程")
}
