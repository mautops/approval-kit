package template_test

import (
	"testing"

	"github.com/mautops/approval-kit/internal/template"
)

// TestEdgeStruct 验证 Edge 结构体定义
func TestEdgeStruct(t *testing.T) {
	// 验证 Edge 类型存在
	var edge *template.Edge
	if edge != nil {
		_ = edge
	}
}

// TestEdgeFields 验证 Edge 结构体的所有字段
func TestEdgeFields(t *testing.T) {
	edge := &template.Edge{
		From:      "node-001",
		To:        "node-002",
		Condition: "",
	}

	// 验证字段值
	if edge.From != "node-001" {
		t.Errorf("Edge.From = %q, want %q", edge.From, "node-001")
	}
	if edge.To != "node-002" {
		t.Errorf("Edge.To = %q, want %q", edge.To, "node-002")
	}
	if edge.Condition != "" {
		t.Errorf("Edge.Condition = %q, want empty string", edge.Condition)
	}
}

// TestEdgeZeroValue 验证 Edge 的零值
func TestEdgeZeroValue(t *testing.T) {
	var edge template.Edge

	// 验证零值
	if edge.From != "" {
		t.Errorf("Edge zero value From = %q, want empty string", edge.From)
	}
	if edge.To != "" {
		t.Errorf("Edge zero value To = %q, want empty string", edge.To)
	}
	if edge.Condition != "" {
		t.Errorf("Edge zero value Condition = %q, want empty string", edge.Condition)
	}
}

// TestEdgeWithCondition 验证 Edge 可以包含条件表达式
func TestEdgeWithCondition(t *testing.T) {
	edge := &template.Edge{
		From:      "node-001",
		To:        "node-002",
		Condition: "amount > 1000",
	}

	if edge.Condition != "amount > 1000" {
		t.Errorf("Edge.Condition = %q, want %q", edge.Condition, "amount > 1000")
	}
}

// TestEdgeWithoutCondition 验证 Edge 可以不包含条件(无条件分支)
func TestEdgeWithoutCondition(t *testing.T) {
	edge := &template.Edge{
		From:      "node-001",
		To:        "node-002",
		Condition: "",
	}

	// 无条件表达式表示默认路径
	if edge.Condition != "" {
		t.Errorf("Edge.Condition should be empty for default path, got %q", edge.Condition)
	}
}

// TestEdgeMultipleEdges 验证可以创建多条边
func TestEdgeMultipleEdges(t *testing.T) {
	edges := []*template.Edge{
		{From: "start", To: "approval-1"},
		{From: "approval-1", To: "approval-2"},
		{From: "approval-2", To: "end"},
	}

	if len(edges) != 3 {
		t.Errorf("Edges length = %d, want 3", len(edges))
	}

	// 验证边的连接关系
	if edges[0].From != "start" || edges[0].To != "approval-1" {
		t.Error("First edge connection incorrect")
	}
	if edges[1].From != "approval-1" || edges[1].To != "approval-2" {
		t.Error("Second edge connection incorrect")
	}
	if edges[2].From != "approval-2" || edges[2].To != "end" {
		t.Error("Third edge connection incorrect")
	}
}

// TestEdgeConditionalBranch 验证条件分支边
func TestEdgeConditionalBranch(t *testing.T) {
	// 条件节点可以有多个输出边,每个边有不同的条件
	edges := []*template.Edge{
		{
			From:      "condition-1",
			To:        "approval-1",
			Condition: "amount > 1000",
		},
		{
			From:      "condition-1",
			To:        "approval-2",
			Condition: "amount <= 1000",
		},
	}

	if len(edges) != 2 {
		t.Errorf("Conditional edges length = %d, want 2", len(edges))
	}

	// 验证条件表达式
	if edges[0].Condition != "amount > 1000" {
		t.Errorf("Edge[0].Condition = %q, want %q", edges[0].Condition, "amount > 1000")
	}
	if edges[1].Condition != "amount <= 1000" {
		t.Errorf("Edge[1].Condition = %q, want %q", edges[1].Condition, "amount <= 1000")
	}
}

