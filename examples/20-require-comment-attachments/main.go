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
	fmt.Println("=== 场景 20: 审批意见和附件必填场景 ===")
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

	// 3. 创建任务
	fmt.Println("步骤 2: 创建重要合同审批任务")
	params := json.RawMessage(`{
		"contractNo": "CT-2025-001",
		"contractType": "重要合同",
		"contractor": "供应商-001",
		"amount": 5000000,
		"description": "重要合同审批,要求审批意见和附件必填"
	}`)
	tsk, err := taskMgr.Create(tpl.ID, "contract-001", params)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}
	fmt.Printf("✓ 任务创建成功: ID=%s, State=%s\n\n", tsk.ID, tsk.State)

	// 4. 执行场景流程
	runScenario(taskMgr, tsk)

	// 5. 输出结果
	printResults(taskMgr, tsk)
}

// createTemplate 创建重要合同审批模板
// 流程: 开始节点 → 审批节点(配置必填) → 结束节点
// 审批节点配置了审批意见和附件必填
func createTemplate() *template.Template {
	now := time.Now()
	return &template.Template{
		ID:          "require-comment-attachments-template",
		Name:        "重要合同审批模板",
		Description: "重要合同审批流程,要求审批意见和附件必填",
		Version:     1,
		CreatedAt:   now,
		UpdatedAt:   now,
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "开始",
				Type: template.NodeTypeStart,
			},
			"contract-approval": {
				ID:   "contract-approval",
				Name: "合同审批",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeSingle,
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{"legal-manager-001"},
					},
					RequireCommentField:    true, // 审批意见必填
					RequireAttachmentsField: true, // 附件必填
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
			{From: "start", To: "contract-approval"},
			{From: "contract-approval", To: "end"},
		},
	}
}

// runScenario 执行场景流程
func runScenario(taskMgr task.TaskManager, tsk *task.Task) {
	// 步骤 3: 提交任务
	fmt.Println("步骤 3: 提交任务进入审批流程")
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

	nodeID := "contract-approval"

	// 步骤 4: 设置审批人
	fmt.Println("步骤 4: 设置审批人")
	err = taskMgr.AddApprover(tsk.ID, nodeID, "legal-manager-001", "设置法务经理为审批人")
	if err != nil {
		log.Fatalf("Failed to add approver: %v", err)
	}
	fmt.Printf("✓ 审批人已设置: legal-manager-001\n\n")

	// 步骤 5: 测试审批意见必填验证
	fmt.Println("步骤 5: 测试审批意见必填验证")
	fmt.Println("  说明: 节点配置了 RequireComment=true,审批时必须填写审批意见")
	fmt.Println("  测试: 尝试不填写审批意见进行审批")
	fmt.Println()

	err = taskMgr.Approve(tsk.ID, nodeID, "legal-manager-001", "") // 空审批意见
	if err != nil {
		fmt.Printf("  ✓ 验证通过: 审批意见为空时被拒绝 (错误: %v)\n", err)
	} else {
		fmt.Println("  ✗ 验证失败: 审批意见为空时不应该允许审批")
	}
	fmt.Println()

	// 步骤 6: 测试附件必填验证
	fmt.Println("步骤 6: 测试附件必填验证")
	fmt.Println("  说明: 节点配置了 RequireAttachments=true,审批时必须上传附件")
	fmt.Println("  测试: 尝试不上传附件进行审批")
	fmt.Println()

	err = taskMgr.ApproveWithAttachments(tsk.ID, nodeID, "legal-manager-001", "审批通过", []string{}) // 空附件列表
	if err != nil {
		fmt.Printf("  ✓ 验证通过: 附件为空时被拒绝 (错误: %v)\n", err)
	} else {
		fmt.Println("  ✗ 验证失败: 附件为空时不应该允许审批")
	}
	fmt.Println()

	// 步骤 7: 正确审批(带审批意见和附件)
	fmt.Println("步骤 7: 正确审批(带审批意见和附件)")
	fmt.Println("  说明: 填写审批意见并上传附件后,审批操作成功")
	fmt.Println()

	attachments := []string{
		"contract-review-report.pdf",
		"legal-opinion.pdf",
	}
	err = taskMgr.ApproveWithAttachments(tsk.ID, nodeID, "legal-manager-001", "合同条款审查通过,已上传审查报告和法律意见书", attachments)
	if err != nil {
		log.Fatalf("Failed to approve: %v", err)
	}

	tsk, _ = taskMgr.Get(tsk.ID)
	fmt.Printf("✓ 审批已通过: ID=%s, State=%s\n", tsk.ID, tsk.State)
	fmt.Printf("  审批意见: 合同条款审查通过,已上传审查报告和法律意见书\n")
	fmt.Printf("  附件数量: %d\n", len(attachments))
	for i, attachment := range attachments {
		fmt.Printf("    附件 %d: %s\n", i+1, attachment)
	}
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

	// 输出审批记录(按时间顺序)
	fmt.Println("\n审批记录(按时间顺序):")
	if len(tsk.Records) == 0 {
		fmt.Println("  无审批记录")
	} else {
		for i, record := range tsk.Records {
			fmt.Printf("  记录 %d:\n", i+1)
			fmt.Printf("    节点: %s\n", getNodeName(record.NodeID))
			fmt.Printf("    审批人: %s\n", record.Approver)
			fmt.Printf("    操作类型: %s\n", getOperationType(record.Result))
			fmt.Printf("    审批意见: %s\n", record.Comment)
			if len(record.Attachments) > 0 {
				fmt.Printf("    附件数量: %d\n", len(record.Attachments))
				for j, attachment := range record.Attachments {
					fmt.Printf("      附件 %d: %s\n", j+1, attachment)
				}
			}
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
	if tsk.State == types.TaskStateApproved {
		fmt.Println("✓ 任务已成功通过审批")
	} else {
		fmt.Printf("✗ 任务状态异常: 期望 %s, 实际 %s\n", types.TaskStateApproved, tsk.State)
	}

	// 验证审批意见和附件必填
	hasComment := false
	hasAttachments := false
	for _, record := range tsk.Records {
		if record.Result == "approve" {
			if record.Comment != "" {
				hasComment = true
			}
			if len(record.Attachments) > 0 {
				hasAttachments = true
			}
		}
	}

	fmt.Println("✓ 审批意见和附件必填验证:")
	if hasComment {
		fmt.Println("  - 审批意见已填写")
	} else {
		fmt.Println("  - ✗ 审批意见未填写")
	}
	if hasAttachments {
		fmt.Println("  - 附件已上传")
	} else {
		fmt.Println("  - ✗ 附件未上传")
	}
}

// getNodeName 获取节点名称
func getNodeName(nodeID string) string {
	nodeNames := map[string]string{
		"contract-approval": "合同审批",
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

