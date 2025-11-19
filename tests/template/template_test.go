package template_test

import (
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/template"
)

// TestTemplateStruct 验证 Template 结构体定义
func TestTemplateStruct(t *testing.T) {
	// 验证 Template 类型存在
	var tmpl *template.Template
	if tmpl != nil {
		_ = tmpl
	}
}

// TestTemplateFields 验证 Template 结构体的所有字段
func TestTemplateFields(t *testing.T) {
	now := time.Now()
	tmpl := &template.Template{
		ID:          "tpl-001",
		Name:        "Test Template",
		Description: "Test Description",
		Version:     1,
		CreatedAt:   now,
		UpdatedAt:   now,
		Nodes:       make(map[string]*template.Node),
		Edges:       []*template.Edge{},
		Config:      nil,
	}

	// 验证字段值
	if tmpl.ID != "tpl-001" {
		t.Errorf("Template.ID = %q, want %q", tmpl.ID, "tpl-001")
	}
	if tmpl.Name != "Test Template" {
		t.Errorf("Template.Name = %q, want %q", tmpl.Name, "Test Template")
	}
	if tmpl.Description != "Test Description" {
		t.Errorf("Template.Description = %q, want %q", tmpl.Description, "Test Description")
	}
	if tmpl.Version != 1 {
		t.Errorf("Template.Version = %d, want %d", tmpl.Version, 1)
	}
	if tmpl.Nodes == nil {
		t.Error("Template.Nodes should not be nil")
	}
	if tmpl.Edges == nil {
		t.Error("Template.Edges should not be nil")
	}
}

// TestTemplateZeroValue 验证 Template 的零值
func TestTemplateZeroValue(t *testing.T) {
	var tmpl template.Template

	// 验证零值
	if tmpl.ID != "" {
		t.Errorf("Template zero value ID = %q, want empty string", tmpl.ID)
	}
	if tmpl.Name != "" {
		t.Errorf("Template zero value Name = %q, want empty string", tmpl.Name)
	}
	if tmpl.Version != 0 {
		t.Errorf("Template zero value Version = %d, want 0", tmpl.Version)
	}
	// 在 Go 中,map 和 slice 的零值是 nil,这是正常的
	// nil map 可以安全读取(返回零值),nil slice 可以安全读取和追加
	if tmpl.Nodes == nil {
		// nil map 是有效的零值
		_ = tmpl.Nodes
	}
	if tmpl.Edges == nil {
		// nil slice 是有效的零值
		_ = tmpl.Edges
	}
	// 验证可以安全使用零值
	if len(tmpl.Nodes) != 0 {
		t.Errorf("Template zero value Nodes length = %d, want 0", len(tmpl.Nodes))
	}
	if len(tmpl.Edges) != 0 {
		t.Errorf("Template zero value Edges length = %d, want 0", len(tmpl.Edges))
	}
}

// TestTemplateWithNodes 验证 Template 可以包含节点
func TestTemplateWithNodes(t *testing.T) {
	tmpl := &template.Template{
		ID:    "tpl-001",
		Nodes: make(map[string]*template.Node),
	}

	// 添加节点
	node := &template.Node{
		ID:   "node-001",
		Name: "Start Node",
		Type: template.NodeTypeStart,
	}
	tmpl.Nodes["node-001"] = node

	// 验证节点
	if len(tmpl.Nodes) != 1 {
		t.Errorf("Template.Nodes length = %d, want 1", len(tmpl.Nodes))
	}
	if tmpl.Nodes["node-001"] == nil {
		t.Error("Template.Nodes should contain node-001")
	}
}

// TestTemplateWithEdges 验证 Template 可以包含边
func TestTemplateWithEdges(t *testing.T) {
	tmpl := &template.Template{
		ID:    "tpl-001",
		Edges: []*template.Edge{},
	}

	// 添加边
	edge := &template.Edge{
		From: "node-001",
		To:   "node-002",
	}
	tmpl.Edges = append(tmpl.Edges, edge)

	// 验证边
	if len(tmpl.Edges) != 1 {
		t.Errorf("Template.Edges length = %d, want 1", len(tmpl.Edges))
	}
	if tmpl.Edges[0].From != "node-001" {
		t.Errorf("Template.Edges[0].From = %q, want %q", tmpl.Edges[0].From, "node-001")
	}
	if tmpl.Edges[0].To != "node-002" {
		t.Errorf("Template.Edges[0].To = %q, want %q", tmpl.Edges[0].To, "node-002")
	}
}

