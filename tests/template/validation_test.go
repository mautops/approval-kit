package template_test

import (
	stderrors "errors"
	"testing"

	"github.com/mautops/approval-kit/internal/errors"
	"github.com/mautops/approval-kit/internal/template"
)

// TestTemplateValidate 测试模板验证
// 使用表驱动测试,覆盖各种验证场景
func TestTemplateValidate(t *testing.T) {
	tests := []struct {
		name    string
		tmpl    *template.Template
		wantErr bool
		errType error
	}{
		{
			name: "valid template with one start node",
			tmpl: &template.Template{
				ID:   "tpl-001",
				Name: "Valid Template",
				Nodes: map[string]*template.Node{
					"start": {
						ID:   "start",
						Type: template.NodeTypeStart,
					},
				},
				Edges: []*template.Edge{},
			},
			wantErr: false,
			errType: nil,
		},
		{
			name: "missing start node",
			tmpl: &template.Template{
				ID:   "tpl-002",
				Name: "No Start Node",
				Nodes: map[string]*template.Node{
					"approval": {
						ID:   "approval",
						Type: template.NodeTypeApproval,
					},
				},
				Edges: []*template.Edge{},
			},
			wantErr: true,
			errType: errors.ErrInvalidTemplate,
		},
		{
			name: "no nodes",
			tmpl: &template.Template{
				ID:    "tpl-003",
				Name:  "Empty Template",
				Nodes: map[string]*template.Node{},
				Edges: []*template.Edge{},
			},
			wantErr: true,
			errType: errors.ErrInvalidTemplate,
		},
		{
			name: "multiple start nodes",
			tmpl: &template.Template{
				ID:   "tpl-004",
				Name: "Multiple Start Nodes",
				Nodes: map[string]*template.Node{
					"start-1": {
						ID:   "start-1",
						Type: template.NodeTypeStart,
					},
					"start-2": {
						ID:   "start-2",
						Type: template.NodeTypeStart,
					},
				},
				Edges: []*template.Edge{},
			},
			wantErr: true,
			errType: errors.ErrInvalidTemplate,
		},
		{
			name: "valid template with edges",
			tmpl: &template.Template{
				ID:   "tpl-005",
				Name: "Template With Edges",
				Nodes: map[string]*template.Node{
					"start": {
						ID:   "start",
						Type: template.NodeTypeStart,
					},
					"approval": {
						ID:   "approval",
						Type: template.NodeTypeApproval,
					},
					"end": {
						ID:   "end",
						Type: template.NodeTypeEnd,
					},
				},
				Edges: []*template.Edge{
					{From: "start", To: "approval"},
					{From: "approval", To: "end"},
				},
			},
			wantErr: false,
			errType: nil,
		},
		{
			name: "edge references non-existent node (from)",
			tmpl: &template.Template{
				ID:   "tpl-006",
				Name: "Invalid Edge From",
				Nodes: map[string]*template.Node{
					"start": {
						ID:   "start",
						Type: template.NodeTypeStart,
					},
				},
				Edges: []*template.Edge{
					{From: "non-existent", To: "start"},
				},
			},
			wantErr: true,
			errType: errors.ErrInvalidTemplate,
		},
		{
			name: "edge references non-existent node (to)",
			tmpl: &template.Template{
				ID:   "tpl-007",
				Name: "Invalid Edge To",
				Nodes: map[string]*template.Node{
					"start": {
						ID:   "start",
						Type: template.NodeTypeStart,
					},
				},
				Edges: []*template.Edge{
					{From: "start", To: "non-existent"},
				},
			},
			wantErr: true,
			errType: errors.ErrInvalidTemplate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.tmpl.Validate()

			// 验证错误
			if (err != nil) != tt.wantErr {
				t.Errorf("Template.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				// 验证错误类型
				if tt.errType != nil {
					if !stderrors.Is(err, tt.errType) {
						t.Errorf("Template.Validate() error type = %v, want %v", err, tt.errType)
					}
				}
			}
		})
	}
}

// TestTemplateValidateStartNodeRequirement 专门测试开始节点的要求
func TestTemplateValidateStartNodeRequirement(t *testing.T) {
	t.Run("exactly one start node required", func(t *testing.T) {
		// 零个开始节点
		tmpl1 := &template.Template{
			Nodes: map[string]*template.Node{
				"approval": {ID: "approval", Type: template.NodeTypeApproval},
			},
		}
		if err := tmpl1.Validate(); err == nil {
			t.Error("Template with no start node should fail validation")
		}

		// 一个开始节点(正确)
		tmpl2 := &template.Template{
			ID:    "tpl-002",
			Name:  "Template 2",
			Nodes: map[string]*template.Node{
				"start": {ID: "start", Type: template.NodeTypeStart},
			},
		}
		if err := tmpl2.Validate(); err != nil {
			t.Errorf("Template with one start node should pass validation, got error: %v", err)
		}

		// 多个开始节点
		tmpl3 := &template.Template{
			Nodes: map[string]*template.Node{
				"start-1": {ID: "start-1", Type: template.NodeTypeStart},
				"start-2": {ID: "start-2", Type: template.NodeTypeStart},
			},
		}
		if err := tmpl3.Validate(); err == nil {
			t.Error("Template with multiple start nodes should fail validation")
		}
	})
}

// TestTemplateValidateEdgeReferences 测试边的节点引用验证
func TestTemplateValidateEdgeReferences(t *testing.T) {
	t.Run("all edge references must exist", func(t *testing.T) {
		// 有效的边引用
		tmpl1 := &template.Template{
			ID:    "tpl-001",
			Name:  "Template 1",
			Nodes: map[string]*template.Node{
				"start":    {ID: "start", Type: template.NodeTypeStart},
				"approval": {ID: "approval", Type: template.NodeTypeApproval},
			},
			Edges: []*template.Edge{
				{From: "start", To: "approval"},
			},
		}
		if err := tmpl1.Validate(); err != nil {
			t.Errorf("Template with valid edge references should pass validation, got error: %v", err)
		}

		// 无效的 From 引用
		tmpl2 := &template.Template{
			Nodes: map[string]*template.Node{
				"start": {ID: "start", Type: template.NodeTypeStart},
			},
			Edges: []*template.Edge{
				{From: "non-existent", To: "start"},
			},
		}
		if err := tmpl2.Validate(); err == nil {
			t.Error("Template with invalid edge From reference should fail validation")
		}

		// 无效的 To 引用
		tmpl3 := &template.Template{
			Nodes: map[string]*template.Node{
				"start": {ID: "start", Type: template.NodeTypeStart},
			},
			Edges: []*template.Edge{
				{From: "start", To: "non-existent"},
			},
		}
		if err := tmpl3.Validate(); err == nil {
			t.Error("Template with invalid edge To reference should fail validation")
		}
	})
}

