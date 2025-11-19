package template_test

import (
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/errors"
	"github.com/mautops/approval-kit/internal/template"
)

// createTestTemplate 创建测试用的模板
func createTestTemplate() *template.Template {
	return &template.Template{
		ID:          "tpl-001",
		Name:        "Test Template",
		Description: "Test template description",
		Version:     1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Nodes: map[string]*template.Node{
			"start": {
				ID:   "start",
				Name: "Start Node",
				Type: template.NodeTypeStart,
			},
		},
		Edges: []*template.Edge{},
	}
}

// TestTemplateManagerCreate 测试模板创建功能
func TestTemplateManagerCreate(t *testing.T) {
	mgr := template.NewTemplateManager()

	tpl := createTestTemplate()

	// 测试创建成功
	err := mgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 验证模板已创建
	got, err := mgr.Get(tpl.ID, tpl.Version)
	if err != nil {
		t.Fatalf("Get() failed after Create: %v", err)
	}
	if got.ID != tpl.ID {
		t.Errorf("Get() returned template ID = %q, want %q", got.ID, tpl.ID)
	}
}

// TestTemplateManagerCreateInvalidTemplate 测试创建无效模板
func TestTemplateManagerCreateInvalidTemplate(t *testing.T) {
	mgr := template.NewTemplateManager()

	tests := []struct {
		name    string
		tpl     *template.Template
		wantErr bool
		errType error
	}{
		{
			name: "missing start node",
			tpl: &template.Template{
				ID:    "tpl-invalid",
				Name:  "Invalid Template",
				Nodes: map[string]*template.Node{},
			},
			wantErr: true,
			errType: errors.ErrInvalidTemplate,
		},
		{
			name: "empty ID",
			tpl: &template.Template{
				ID:    "",
				Name:  "Template",
				Nodes: map[string]*template.Node{
					"start": {
						ID:   "start",
						Type: template.NodeTypeStart,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "empty name",
			tpl: &template.Template{
				ID:    "tpl-001",
				Name:  "",
				Nodes: map[string]*template.Node{
					"start": {
						ID:   "start",
						Type: template.NodeTypeStart,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mgr.Create(tt.tpl)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errType != nil {
				if err != tt.errType {
					t.Errorf("Create() error = %v, want %v", err, tt.errType)
				}
			}
		})
	}
}

// TestTemplateManagerCreateDuplicate 测试创建重复模板
func TestTemplateManagerCreateDuplicate(t *testing.T) {
	mgr := template.NewTemplateManager()

	tpl := createTestTemplate()

	// 第一次创建应该成功
	err := mgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create() failed on first create: %v", err)
	}

	// 第二次创建相同 ID 和版本应该失败
	err = mgr.Create(tpl)
	if err == nil {
		t.Error("Create() should fail when template with same ID and version already exists")
	}
}

// TestTemplateManagerCreateDifferentVersions 测试创建不同版本的模板
func TestTemplateManagerCreateDifferentVersions(t *testing.T) {
	mgr := template.NewTemplateManager()

	tpl1 := createTestTemplate()
	tpl1.Version = 1

	tpl2 := createTestTemplate()
	tpl2.Version = 2

	// 创建版本 1
	err := mgr.Create(tpl1)
	if err != nil {
		t.Fatalf("Create() failed for version 1: %v", err)
	}

	// 创建版本 2 应该成功
	err = mgr.Create(tpl2)
	if err != nil {
		t.Fatalf("Create() failed for version 2: %v", err)
	}

	// 验证两个版本都存在
	got1, err := mgr.Get(tpl1.ID, 1)
	if err != nil {
		t.Fatalf("Get() failed for version 1: %v", err)
	}
	if got1.Version != 1 {
		t.Errorf("Get() returned version = %d, want 1", got1.Version)
	}

	got2, err := mgr.Get(tpl2.ID, 2)
	if err != nil {
		t.Fatalf("Get() failed for version 2: %v", err)
	}
	if got2.Version != 2 {
		t.Errorf("Get() returned version = %d, want 2", got2.Version)
	}
}

