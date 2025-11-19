package task_test

import (
	"testing"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
	"github.com/mautops/approval-kit/internal/types"
)

// TestRejectJumpTargetNodeNotFound 测试拒绝后跳转但目标节点不存在
func TestRejectJumpTargetNodeNotFound(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	
	// 创建包含跳转配置但目标节点不存在的模板
	tpl := &template.Template{
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
					Mode:            node.ApprovalModeSingle,
					RejectBehavior:  node.RejectBehaviorJump,
					RejectTargetNode: "non-existent-node",
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
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	taskMgr := task.NewTaskManager(templateMgr, nil)
	tsk, err := taskMgr.Create("tpl-001", "biz-001", nil)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	err = taskMgr.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Submit() failed: %v", err)
	}

	// 拒绝审批(应该失败,因为目标节点不存在)
	err = taskMgr.Reject(tsk.ID, "approval-001", "user-001", "rejected")
	if err == nil {
		t.Error("Reject() should fail when target node does not exist")
	}
}

// TestRejectJumpNoTargetNode 测试拒绝后跳转但未指定目标节点
func TestRejectJumpNoTargetNode(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	
	// 创建包含跳转配置但未指定目标节点的模板
	tpl := &template.Template{
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
					Mode:            node.ApprovalModeSingle,
					RejectBehavior:  node.RejectBehaviorJump,
					RejectTargetNode: "", // 未指定目标节点
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
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	taskMgr := task.NewTaskManager(templateMgr, nil)
	tsk, err := taskMgr.Create("tpl-001", "biz-001", nil)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	err = taskMgr.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Submit() failed: %v", err)
	}

	// 拒绝审批(应该终止流程,因为未指定目标节点)
	err = taskMgr.Reject(tsk.ID, "approval-001", "user-001", "rejected")
	if err != nil {
		t.Fatalf("Reject() should succeed and terminate: %v", err)
	}

	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	// 验证任务状态已变为已拒绝
	if tsk.State != types.TaskStateRejected {
		t.Errorf("Task.State = %q, want %q", tsk.State, types.TaskStateRejected)
	}
}

// TestRejectRollbackToStartNode 测试拒绝后回滚到 start 节点
// 这是一个实际的场景: 当审批节点回滚时,会回滚到 start 节点
func TestRejectRollbackToStartNode(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	
	tpl := &template.Template{
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
					Mode:           node.ApprovalModeSingle,
					RejectBehavior: node.RejectBehaviorRollback,
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
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	taskMgr := task.NewTaskManager(templateMgr, nil)
	tsk, err := taskMgr.Create("tpl-001", "biz-001", nil)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	err = taskMgr.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Submit() failed: %v", err)
	}

	// 拒绝审批(应该回滚到 start 节点)
	err = taskMgr.Reject(tsk.ID, "approval-001", "user-001", "rejected")
	if err != nil {
		t.Fatalf("Reject() should succeed: %v", err)
	}

	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	// 验证任务已回滚到 start 节点
	if tsk.CurrentNode != "start" {
		t.Errorf("Task.CurrentNode = %q, want %q", tsk.CurrentNode, "start")
	}
	if tsk.State != types.TaskStateApproving {
		t.Errorf("Task.State = %q, want %q", tsk.State, types.TaskStateApproving)
	}
}

// TestRejectNonApprovalNode 测试拒绝非审批节点
func TestRejectNonApprovalNode(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	
	// 创建包含非审批节点的模板
	tpl := &template.Template{
		ID:   "tpl-001",
		Name: "Test Template",
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
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	taskMgr := task.NewTaskManager(templateMgr, nil)
	tsk, err := taskMgr.Create("tpl-001", "biz-001", nil)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	err = taskMgr.Submit(tsk.ID)
	if err != nil {
		t.Fatalf("Submit() failed: %v", err)
	}

	// 拒绝非审批节点(应该终止流程)
	err = taskMgr.Reject(tsk.ID, "start", "user-001", "rejected")
	if err != nil {
		t.Fatalf("Reject() should succeed: %v", err)
	}

	tsk, err = taskMgr.Get(tsk.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	// 验证任务状态已变为已拒绝
	if tsk.State != types.TaskStateRejected {
		t.Errorf("Task.State = %q, want %q", tsk.State, types.TaskStateRejected)
	}
}

