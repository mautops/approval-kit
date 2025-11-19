package task_test

import (
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
)

// TestOperationPermissions 测试操作权限控制
func TestOperationPermissions(t *testing.T) {
	tests := []struct {
		name           string
		permissions    node.OperationPermissions
		allowTransfer  bool
		allowAddApprover bool
		allowRemoveApprover bool
	}{
		{
			name: "all permissions allowed",
			permissions: node.OperationPermissions{
				AllowTransfer:      true,
				AllowAddApprover:   true,
				AllowRemoveApprover: true,
			},
			allowTransfer:      true,
			allowAddApprover:   true,
			allowRemoveApprover: true,
		},
		{
			name: "no permissions",
			permissions: node.OperationPermissions{
				AllowTransfer:      false,
				AllowAddApprover:   false,
				AllowRemoveApprover: false,
			},
			allowTransfer:      false,
			allowAddApprover:   false,
			allowRemoveApprover: false,
		},
		{
			name: "only transfer allowed",
			permissions: node.OperationPermissions{
				AllowTransfer:      true,
				AllowAddApprover:   false,
				AllowRemoveApprover: false,
			},
			allowTransfer:      true,
			allowAddApprover:   false,
			allowRemoveApprover: false,
		},
		{
			name: "only add approver allowed",
			permissions: node.OperationPermissions{
				AllowTransfer:      false,
				AllowAddApprover:   true,
				AllowRemoveApprover: false,
			},
			allowTransfer:      false,
			allowAddApprover:   true,
			allowRemoveApprover: false,
		},
		{
			name: "only remove approver allowed",
			permissions: node.OperationPermissions{
				AllowTransfer:      false,
				AllowAddApprover:   false,
				AllowRemoveApprover: true,
			},
			allowTransfer:      false,
			allowAddApprover:   false,
			allowRemoveApprover: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 验证权限配置
			if tt.permissions.AllowTransfer != tt.allowTransfer {
				t.Errorf("AllowTransfer = %v, want %v", tt.permissions.AllowTransfer, tt.allowTransfer)
			}
			if tt.permissions.AllowAddApprover != tt.allowAddApprover {
				t.Errorf("AllowAddApprover = %v, want %v", tt.permissions.AllowAddApprover, tt.allowAddApprover)
			}
			if tt.permissions.AllowRemoveApprover != tt.allowRemoveApprover {
				t.Errorf("AllowRemoveApprover = %v, want %v", tt.permissions.AllowRemoveApprover, tt.allowRemoveApprover)
			}
		})
	}
}

// TestPermissionsInApprovalNodeConfig 测试审批节点配置中的权限控制
func TestPermissionsInApprovalNodeConfig(t *testing.T) {
	// 创建带权限配置的模板
	templateMgr := template.NewTemplateManager()
	tpl := createTestTemplateWithPermissions()
	err := templateMgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create template failed: %v", err)
	}

	// 获取模板
	createdTpl, err := templateMgr.Get("tpl-001", 0)
	if err != nil {
		t.Fatalf("Get template failed: %v", err)
	}

	// 验证节点配置中的权限
	node := createdTpl.Nodes["approval-node"]
	if node == nil {
		t.Fatal("Approval node not found")
	}

	config, ok := node.Config.(template.ApprovalNodeConfigAccessor)
	if !ok {
		t.Fatal("Node config is not ApprovalNodeConfig")
	}

	// 验证权限配置
	perms, ok := config.GetPermissions().(template.OperationPermissionsAccessor)
	if !ok {
		t.Fatal("Permissions is not OperationPermissionsAccessor")
	}
	if !perms.AllowTransfer() {
		t.Error("AllowTransfer should be true")
	}
	if !perms.AllowAddApprover() {
		t.Error("AllowAddApprover should be true")
	}
	if !perms.AllowRemoveApprover() {
		t.Error("AllowRemoveApprover should be true")
	}
}

// TestPermissionsValidation 测试权限验证逻辑
func TestPermissionsValidation(t *testing.T) {
	// 创建模板管理器
	templateMgr := template.NewTemplateManager()
	taskMgr := task.NewTaskManager(templateMgr, nil)

	// 测试转交权限验证
	t.Run("transfer permission check", func(t *testing.T) {
		tpl := createTestTemplateWithoutTransferPermission()
		err := templateMgr.Create(tpl)
		if err != nil {
			t.Fatalf("Create template failed: %v", err)
		}

		tsk, err := taskMgr.Create("tpl-001", "biz-001", nil)
		if err != nil {
			t.Fatalf("Create() failed: %v", err)
		}

		err = taskMgr.Submit(tsk.ID)
		if err != nil {
			t.Fatalf("Submit() failed: %v", err)
		}

		// 尝试转交(应该失败,因为没有权限)
		err = taskMgr.Transfer(tsk.ID, "approval-node", "user-1", "user-2", "reason")
		if err == nil {
			t.Error("Transfer() should fail when permission is not allowed")
		}
	})

	// 测试加签权限验证
	t.Run("add approver permission check", func(t *testing.T) {
		tpl := createTestTemplateWithoutAddApproverPermission()
		err := templateMgr.Create(tpl)
		if err != nil {
			t.Fatalf("Create template failed: %v", err)
		}

		tsk, err := taskMgr.Create("tpl-002", "biz-002", nil)
		if err != nil {
			t.Fatalf("Create() failed: %v", err)
		}

		err = taskMgr.Submit(tsk.ID)
		if err != nil {
			t.Fatalf("Submit() failed: %v", err)
		}

		// 尝试加签(应该失败,因为没有权限)
		err = taskMgr.AddApprover(tsk.ID, "approval-node", "user-2", "reason")
		if err == nil {
			t.Error("AddApprover() should fail when permission is not allowed")
		}
	})

	// 测试减签权限验证
	t.Run("remove approver permission check", func(t *testing.T) {
		tpl := createTestTemplateWithoutRemoveApproverPermission()
		err := templateMgr.Create(tpl)
		if err != nil {
			t.Fatalf("Create template failed: %v", err)
		}

		tsk, err := taskMgr.Create("tpl-003", "biz-003", nil)
		if err != nil {
			t.Fatalf("Create() failed: %v", err)
		}

		err = taskMgr.Submit(tsk.ID)
		if err != nil {
			t.Fatalf("Submit() failed: %v", err)
		}

		// 尝试减签(应该失败,因为没有权限)
		err = taskMgr.RemoveApprover(tsk.ID, "approval-node", "user-1", "reason")
		if err == nil {
			t.Error("RemoveApprover() should fail when permission is not allowed")
		}
	})
}

// createTestTemplateWithPermissions 创建带所有权限的测试模板
func createTestTemplateWithPermissions() *template.Template {
	return &template.Template{
		ID:          "tpl-001",
		Name:        "Test Template With Permissions",
		Description: "Test template with all permissions",
		Version:     1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "Start Node",
				Type: template.NodeTypeStart,
			},
			"approval-node": {
				ID:   "approval-node",
				Name: "Approval Node",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeSingle,
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{"user-1"},
					},
					Permissions: node.OperationPermissions{
						AllowTransfer:      true,
						AllowAddApprover:   true,
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
			{From: "start", To: "approval-node"},
			{From: "approval-node", To: "end"},
		},
	}
}

// createTestTemplateWithoutTransferPermission 创建不允许转交的测试模板
func createTestTemplateWithoutTransferPermission() *template.Template {
	return &template.Template{
		ID:          "tpl-001",
		Name:        "Test Template Without Transfer",
		Description: "Test template without transfer permission",
		Version:     1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "Start Node",
				Type: template.NodeTypeStart,
			},
			"approval-node": {
				ID:   "approval-node",
				Name: "Approval Node",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeSingle,
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{"user-1"},
					},
					Permissions: node.OperationPermissions{
						AllowTransfer: false,
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
			{From: "start", To: "approval-node"},
			{From: "approval-node", To: "end"},
		},
	}
}

// createTestTemplateWithoutAddApproverPermission 创建不允许加签的测试模板
func createTestTemplateWithoutAddApproverPermission() *template.Template {
	return &template.Template{
		ID:          "tpl-002",
		Name:        "Test Template Without Add Approver",
		Description: "Test template without add approver permission",
		Version:     1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "Start Node",
				Type: template.NodeTypeStart,
			},
			"approval-node": {
				ID:   "approval-node",
				Name: "Approval Node",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeSingle,
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{"user-1"},
					},
					Permissions: node.OperationPermissions{
						AllowAddApprover: false,
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
			{From: "start", To: "approval-node"},
			{From: "approval-node", To: "end"},
		},
	}
}

// createTestTemplateWithoutRemoveApproverPermission 创建不允许减签的测试模板
func createTestTemplateWithoutRemoveApproverPermission() *template.Template {
	return &template.Template{
		ID:          "tpl-003",
		Name:        "Test Template Without Remove Approver",
		Description: "Test template without remove approver permission",
		Version:     1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "Start Node",
				Type: template.NodeTypeStart,
			},
			"approval-node": {
				ID:   "approval-node",
				Name: "Approval Node",
				Type: template.NodeTypeApproval,
				Config: &node.ApprovalNodeConfig{
					Mode: node.ApprovalModeSingle,
					ApproverConfig: &node.FixedApproverConfig{
						Approvers: []string{"user-1"},
					},
					Permissions: node.OperationPermissions{
						AllowRemoveApprover: false,
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
			{From: "start", To: "approval-node"},
			{From: "approval-node", To: "end"},
		},
	}
}

