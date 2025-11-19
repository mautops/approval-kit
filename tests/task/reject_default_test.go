package task_test

import (
	"testing"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
	"github.com/mautops/approval-kit/internal/types"
)

// TestRejectDefaultBehavior 测试拒绝后默认行为(RejectBehavior 为空字符串)
func TestRejectDefaultBehavior(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	
	// 创建 RejectBehavior 为空字符串的模板(应该使用默认行为: terminate)
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
					Mode: node.ApprovalModeSingle,
					// RejectBehavior 为空字符串,应该使用默认行为
					RejectBehavior: node.RejectBehavior(""),
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

	// 拒绝审批(应该终止流程,使用默认行为)
	err = taskMgr.Reject(tsk.ID, "approval-001", "user-001", "rejected")
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

// TestRejectConfigNotAccessor 测试节点配置不能转换为 ApprovalNodeConfigAccessor 的情况
func TestRejectConfigNotAccessor(t *testing.T) {
	templateMgr := template.NewTemplateManager()
	
	// 创建一个 Config 为 nil 的审批节点
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

	// 拒绝审批(应该终止流程,因为无法获取配置)
	err = taskMgr.Reject(tsk.ID, "approval-001", "user-001", "rejected")
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

