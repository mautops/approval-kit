package template_test

import (
	"testing"

	"github.com/mautops/approval-kit/internal/template"
)

// TestNodeStruct 验证 Node 结构体定义
func TestNodeStruct(t *testing.T) {
	// 验证 Node 类型存在
	var node *template.Node
	if node != nil {
		_ = node
	}
}

// TestNodeFields 验证 Node 结构体的所有字段
func TestNodeFields(t *testing.T) {
	node := &template.Node{
		ID:     "node-001",
		Name:   "Start Node",
		Type:   template.NodeTypeStart,
		Order:  1,
		Config: nil, // 配置将在后续任务中实现
	}

	// 验证字段值
	if node.ID != "node-001" {
		t.Errorf("Node.ID = %q, want %q", node.ID, "node-001")
	}
	if node.Name != "Start Node" {
		t.Errorf("Node.Name = %q, want %q", node.Name, "Start Node")
	}
	if node.Type != template.NodeTypeStart {
		t.Errorf("Node.Type = %v, want %v", node.Type, template.NodeTypeStart)
	}
	if node.Order != 1 {
		t.Errorf("Node.Order = %d, want %d", node.Order, 1)
	}
}

// TestNodeZeroValue 验证 Node 的零值
func TestNodeZeroValue(t *testing.T) {
	var node template.Node

	// 验证零值
	if node.ID != "" {
		t.Errorf("Node zero value ID = %q, want empty string", node.ID)
	}
	if node.Name != "" {
		t.Errorf("Node zero value Name = %q, want empty string", node.Name)
	}
	if node.Type != "" {
		t.Errorf("Node zero value Type = %q, want empty string", node.Type)
	}
	if node.Order != 0 {
		t.Errorf("Node zero value Order = %d, want 0", node.Order)
	}
}

// TestNodeTypes 验证 Node 可以设置为不同的节点类型
func TestNodeTypes(t *testing.T) {
	tests := []struct {
		name     string
		nodeType template.NodeType
	}{
		{
			name:     "start node",
			nodeType: template.NodeTypeStart,
		},
		{
			name:     "approval node",
			nodeType: template.NodeTypeApproval,
		},
		{
			name:     "condition node",
			nodeType: template.NodeTypeCondition,
		},
		{
			name:     "end node",
			nodeType: template.NodeTypeEnd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &template.Node{
				ID:    "node-001",
				Name:  "Test Node",
				Type:  tt.nodeType,
				Order: 1,
			}

			if node.Type != tt.nodeType {
				t.Errorf("Node.Type = %v, want %v", node.Type, tt.nodeType)
			}
		})
	}
}

// TestNodeConfigInterface 验证 NodeConfig 接口存在
func TestNodeConfigInterface(t *testing.T) {
	// 验证接口类型存在
	var config template.NodeConfig
	if config != nil {
		_ = config
	}
}

// TestNodeOrder 验证节点顺序字段
func TestNodeOrder(t *testing.T) {
	nodes := []*template.Node{
		{ID: "node-1", Order: 1},
		{ID: "node-2", Order: 2},
		{ID: "node-3", Order: 3},
	}

	for i, node := range nodes {
		if node.Order != i+1 {
			t.Errorf("Node[%d].Order = %d, want %d", i, node.Order, i+1)
		}
	}
}

