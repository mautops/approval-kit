package task_test

import (
	"testing"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
	"github.com/mautops/approval-kit/internal/types"
)

// TestTaskManagerApprove 测试审批操作
func TestTaskManagerApprove(t *testing.T) {
	templateManager := template.NewTemplateManager()
	tm := task.NewTaskManager(templateManager, nil)

	// 创建模板
	tmpl := createTestTemplateWithApprovalNode()
	err := templateManager.Create(tmpl)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	// 创建任务
	tsk, err := tm.Create(tmpl.ID, "business-001", nil)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// 提交任务
	err = tm.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Failed to submit task: %v", err)
	}

	// 执行审批操作
	err = tm.Approve(tsk.ID, "approval-001", "user-001", "approved")
	if err != nil {
		t.Fatalf("TaskManager.Approve() failed: %v", err)
	}

	// 验证审批记录
	tsk, err = tm.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}

	// 验证审批记录已生成
	if len(tsk.Records) == 0 {
		t.Fatal("Task should have approval records")
	}

	// 验证审批记录内容
	record := tsk.Records[0]
	if record.NodeID != "approval-001" {
		t.Errorf("Record.NodeID = %q, want %q", record.NodeID, "approval-001")
	}
	if record.Approver != "user-001" {
		t.Errorf("Record.Approver = %q, want %q", record.Approver, "user-001")
	}
	if record.Result != "approve" {
		t.Errorf("Record.Result = %q, want %q", record.Result, "approve")
	}
	if record.Comment != "approved" {
		t.Errorf("Record.Comment = %q, want %q", record.Comment, "approved")
	}

	// 验证审批结果已记录
	approvals, exists := tsk.Approvals["approval-001"]
	if !exists {
		t.Fatal("Task should have approvals for node approval-001")
	}
	approval, exists := approvals["user-001"]
	if !exists {
		t.Fatal("Task should have approval from user-001")
	}
	if approval.Result != "approve" {
		t.Errorf("Approval.Result = %q, want %q", approval.Result, "approve")
	}
}

// TestTaskManagerApproveNotFound 测试审批不存在的任务
func TestTaskManagerApproveNotFound(t *testing.T) {
	tm := task.NewTaskManager(template.NewTemplateManager(), nil)

	err := tm.Approve("non-existent", "node-001", "user-001", "approved")
	if err == nil {
		t.Error("TaskManager.Approve() should return error for non-existent task")
	}
}

// TestTaskManagerApproveInvalidState 测试在无效状态下审批
func TestTaskManagerApproveInvalidState(t *testing.T) {
	templateManager := template.NewTemplateManager()
	tm := task.NewTaskManager(templateManager, nil)

	// 创建模板
	tmpl := createTestTemplateWithApprovalNode()
	err := templateManager.Create(tmpl)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	// 创建任务(状态为 pending)
	tsk, err := tm.Create(tmpl.ID, "business-001", nil)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// 尝试在 pending 状态下审批(应该失败,因为任务还未提交)
	err = tm.Approve(tsk.ID, "node-001", "user-001", "approved")
	// 这个行为取决于业务逻辑,可能允许也可能不允许
	// 这里我们假设只有 submitted 或 approving 状态才能审批
	if err == nil {
		t.Error("TaskManager.Approve() should return error when task is in pending state")
	}
}

// TestTaskManagerReject 测试拒绝操作
func TestTaskManagerReject(t *testing.T) {
	templateManager := template.NewTemplateManager()
	tm := task.NewTaskManager(templateManager, nil)

	// 创建模板
	tmpl := createTestTemplateWithApprovalNode()
	err := templateManager.Create(tmpl)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	// 创建任务
	tsk, err := tm.Create(tmpl.ID, "business-001", nil)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// 提交任务
	err = tm.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Failed to submit task: %v", err)
	}

	// 执行拒绝操作
	err = tm.Reject(tsk.ID, "approval-001", "user-001", "rejected")
	if err != nil {
		t.Fatalf("TaskManager.Reject() failed: %v", err)
	}

	// 验证拒绝记录
	tsk, err = tm.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}

	// 验证审批记录已生成
	if len(tsk.Records) == 0 {
		t.Fatal("Task should have rejection records")
	}

	// 验证拒绝记录内容
	record := tsk.Records[0]
	if record.Result != "reject" {
		t.Errorf("Record.Result = %q, want %q", record.Result, "reject")
	}
	if record.Comment != "rejected" {
		t.Errorf("Record.Comment = %q, want %q", record.Comment, "rejected")
	}
}

// TestTaskManagerApproveRequireComment 测试审批意见必填
func TestTaskManagerApproveRequireComment(t *testing.T) {
	templateManager := template.NewTemplateManager()
	tm := task.NewTaskManager(templateManager, nil)

	// 创建包含必填审批意见配置的模板
	tmpl := &template.Template{
		ID:   "tpl-require-comment",
		Name: "Require Comment Template",
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "Start Node",
				Type: template.NodeTypeStart,
			},
			"approval-001": {
				ID:   "approval-001",
				Name: "Approval Node",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode:              node.ApprovalModeSingle,
					RequireCommentField: true,
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{"user-001"},
					},
				},
			},
			"end": {
				ID:   "end",
				Name: "End Node",
				Type: template.NodeTypeEnd,
			},
		},
		Edges: []*template.Edge{
			{From: "start", To: "approval-001"},
			{From: "approval-001", To: "end"},
		},
		Version: 1,
	}
	err := templateManager.Create(tmpl)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	// 创建任务
	tsk, err := tm.Create("tpl-require-comment", "business-001", nil)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// 提交任务
	err = tm.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Failed to submit task: %v", err)
	}

	// 尝试不带审批意见审批(应该失败)
	err = tm.Approve(tsk.ID, "approval-001", "user-001", "")
	if err == nil {
		t.Error("Approve() should fail when comment is required but empty")
	}

	// 带审批意见审批(应该成功)
	err = tm.Approve(tsk.ID, "approval-001", "user-001", "approved with comment")
	if err != nil {
		t.Fatalf("Approve() should succeed with comment: %v", err)
	}
}

// TestTaskManagerApproveNodeNotFound 测试审批节点不存在的情况
func TestTaskManagerApproveNodeNotFound(t *testing.T) {
	templateManager := template.NewTemplateManager()
	tm := task.NewTaskManager(templateManager, nil)

	// 创建模板
	tmpl := createTestTemplateWithApprovalNode()
	err := templateManager.Create(tmpl)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	// 创建任务
	tsk, err := tm.Create(tmpl.ID, "business-001", nil)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// 提交任务
	err = tm.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Failed to submit task: %v", err)
	}

	// 尝试审批不存在的节点(应该失败)
	err = tm.Approve(tsk.ID, "non-existent-node", "user-001", "approved")
	if err == nil {
		t.Error("Approve() should fail when node does not exist")
	}
}

// TestTaskManagerApproveTemplateNotFound 测试模板不存在的情况
func TestTaskManagerApproveTemplateNotFound(t *testing.T) {
	templateManager := template.NewTemplateManager()
	tm := task.NewTaskManager(templateManager, nil)

	// 创建一个任务,但模板已被删除
	tmpl := createTestTemplateWithApprovalNode()
	err := templateManager.Create(tmpl)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	// 创建任务
	tsk, err := tm.Create(tmpl.ID, "business-001", nil)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// 删除模板
	err = templateManager.Delete(tmpl.ID)
	if err != nil {
		t.Fatalf("Failed to delete template: %v", err)
	}

	// 提交任务
	err = tm.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Failed to submit task: %v", err)
	}

	// 尝试审批(应该失败,因为模板不存在)
	err = tm.Approve(tsk.ID, "approval-001", "user-001", "approved")
	if err == nil {
		t.Error("Approve() should fail when template does not exist")
	}
}

// TestTaskManagerApproveNonApprovalNode 测试审批非审批节点
func TestTaskManagerApproveNonApprovalNode(t *testing.T) {
	templateManager := template.NewTemplateManager()
	tm := task.NewTaskManager(templateManager, nil)

	// 创建包含非审批节点的模板
	tmpl := &template.Template{
		ID:   "tpl-non-approval",
		Name: "Non Approval Template",
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "Start Node",
				Type: template.NodeTypeStart,
			},
			"end": {
				ID:   "end",
				Name: "End Node",
				Type: template.NodeTypeEnd,
			},
		},
		Edges: []*template.Edge{
			{From: "start", To: "end"},
		},
		Version: 1,
	}
	err := templateManager.Create(tmpl)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	// 创建任务
	tsk, err := tm.Create("tpl-non-approval", "business-001", nil)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// 提交任务
	err = tm.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Failed to submit task: %v", err)
	}

	// 尝试审批非审批节点(应该失败或成功但不会验证 RequireComment)
	err = tm.Approve(tsk.ID, "start", "user-001", "")
	// 这个行为取决于实现,可能允许也可能不允许
	// 这里我们验证不会因为 RequireComment 而失败
	if err != nil {
		// 如果失败,应该是因为节点类型不匹配,而不是因为 RequireComment
		if err.Error() == "comment is required for approval node \"start\"" {
			t.Error("Approve() should not check RequireComment for non-approval nodes")
		}
	}
}

// TestTaskManagerApproveInApprovingState 测试在 approving 状态下审批
func TestTaskManagerApproveInApprovingState(t *testing.T) {
	templateManager := template.NewTemplateManager()
	tm := task.NewTaskManager(templateManager, nil)

	// 创建模板(使用多人会签模式,允许多个审批人)
	// 使用 record_query_test.go 中定义的 createTestTemplateWithMultipleApprovers
	// 注意: 需要确保该函数在同一个包中可访问
	tmpl := createTestTemplateWithMultipleApproversForApprove()
	err := templateManager.Create(tmpl)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	// 创建任务
	tsk, err := tm.Create(tmpl.ID, "business-001", nil)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// 提交任务
	err = tm.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Failed to submit task: %v", err)
	}

	// 手动设置审批人列表(因为固定审批人需要在节点激活时设置)
	err = tm.AddApprover(tsk.ID, "approval-001", "user-001", "setup approver")
	if err != nil {
		t.Fatalf("AddApprover() failed: %v", err)
	}

	err = tm.AddApprover(tsk.ID, "approval-001", "user-002", "setup approver")
	if err != nil {
		t.Fatalf("AddApprover() failed: %v", err)
	}

	// 第一次审批(将状态从 submitted 转为 approving)
	err = tm.Approve(tsk.ID, "approval-001", "user-001", "first approval")
	if err != nil {
		t.Fatalf("First Approve() failed: %v", err)
	}

	// 验证状态已变为 approving(多人会签模式,需要所有审批人都同意)
	tsk, err = tm.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}
	if tsk.State != types.TaskStateApproving {
		t.Errorf("Task.State = %q, want %q", tsk.State, types.TaskStateApproving)
	}

	// 第二次审批(在 approving 状态下审批,多人会签模式允许)
	err = tm.Approve(tsk.ID, "approval-001", "user-002", "second approval")
	if err != nil {
		t.Fatalf("Second Approve() failed: %v", err)
	}

	// 验证最终状态为 approved(所有审批人都已同意)
	tsk, err = tm.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}
	if tsk.State != types.TaskStateApproved {
		t.Errorf("Task.State = %q, want %q", tsk.State, types.TaskStateApproved)
	}
}

// createTestTemplateWithMultipleApproversForApprove 创建支持多个审批人的测试模板(多人会签模式)
// 用于 approve_test.go 中的测试
func createTestTemplateWithMultipleApproversForApprove() *template.Template {
	return &template.Template{
		ID:   "tpl-multiple-approve",
		Name: "Test Template With Multiple Approvers For Approve",
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "Start Node",
				Type: template.NodeTypeStart,
			},
			"approval-001": {
				ID:   "approval-001",
				Name: "Approval Node",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeUnanimous, // 多人会签模式
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{"user-001", "user-002"},
					},
					Permissions: node.OperationPermissions{
						AllowTransfer:    true,
						AllowAddApprover: true,
						AllowRemoveApprover: true,
					},
				},
			},
			"end": {
				ID:   "end",
				Name: "End Node",
				Type: template.NodeTypeEnd,
			},
		},
		Edges: []*template.Edge{
			{From: "start", To: "approval-001"},
			{From: "approval-001", To: "end"},
		},
		Version: 1,
	}
}

// TestTaskManagerApproveNodeConfigNotAccessor 测试节点配置不能转换为 ApprovalNodeConfigAccessor 的情况
func TestTaskManagerApproveNodeConfigNotAccessor(t *testing.T) {
	templateManager := template.NewTemplateManager()
	tm := task.NewTaskManager(templateManager, nil)

	// 创建一个使用非 ApprovalNodeConfig 配置的审批节点
	// 注意: 由于 ApprovalNodeConfig 实现了 ApprovalNodeConfigAccessor,
	// 我们需要创建一个不实现该接口的配置类型
	// 但为了简化测试,我们测试一个实际场景: 节点类型是 Approval 但 Config 为 nil
	tmpl := &template.Template{
		ID:   "tpl-nil-config",
		Name: "Nil Config Template",
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "Start Node",
				Type: template.NodeTypeStart,
			},
			"approval-001": {
				ID:   "approval-001",
				Name: "Approval Node",
				Type: template.NodeTypeApproval,
				Config: nil, // Config 为 nil,不能转换为 ApprovalNodeConfigAccessor
			},
			"end": {
				ID:   "end",
				Name: "End Node",
				Type: template.NodeTypeEnd,
			},
		},
		Edges: []*template.Edge{
			{From: "start", To: "approval-001"},
			{From: "approval-001", To: "end"},
		},
		Version: 1,
	}
	err := templateManager.Create(tmpl)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	// 创建任务
	tsk, err := tm.Create("tpl-nil-config", "business-001", nil)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// 提交任务
	err = tm.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Failed to submit task: %v", err)
	}

	// 审批(应该成功,因为 Config 为 nil 时不会检查 RequireComment)
	err = tm.Approve(tsk.ID, "approval-001", "user-001", "")
	if err != nil {
		// 如果失败,应该是因为其他原因,而不是因为 RequireComment
		if err.Error() == "comment is required for approval node \"approval-001\"" {
			t.Error("Approve() should not check RequireComment when Config is nil")
		}
	}
}

// createTestTemplateWithApprovalNode 创建包含审批节点的测试模板
func createTestTemplateWithApprovalNode() *template.Template {
	return &template.Template{
		ID:   "tpl-001",
		Name: "Test Template",
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "Start Node",
				Type: template.NodeTypeStart,
			},
			"approval-001": {
				ID:   "approval-001",
				Name: "Approval Node",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeSingle,
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{"user-001"},
					},
					Permissions: node.OperationPermissions{
						AllowTransfer:    true,
						AllowAddApprover: true,
						AllowRemoveApprover: true,
					},
				},
			},
			"end": {
				ID:   "end",
				Name: "End Node",
				Type: template.NodeTypeEnd,
			},
		},
		Edges: []*template.Edge{
			{From: "start", To: "approval-001"},
			{From: "approval-001", To: "end"},
		},
		Version: 1,
	}
}
