package node_test

import (
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/template"
)

// TestApprovalNodeConfigStruct 验证 ApprovalNodeConfig 结构体定义
func TestApprovalNodeConfigStruct(t *testing.T) {
	// 验证结构体类型存在
	var config *node.ApprovalNodeConfig
	if config != nil {
		_ = config
	}
}

// TestApprovalNodeConfigFields 验证 ApprovalNodeConfig 结构体的所有字段
func TestApprovalNodeConfigFields(t *testing.T) {
	timeout := 24 * time.Hour
	config := &node.ApprovalNodeConfig{
		Mode:            node.ApprovalModeSingle,
		ApproverConfig: nil, // 将在后续任务中实现
		Timeout:         &timeout,
		RejectBehavior:  node.RejectBehaviorTerminate,
		Permissions:     node.OperationPermissions{},
		RequireCommentField:  true,
		RequireAttachmentsField: false,
	}

	// 验证字段值
	if config.Mode != node.ApprovalModeSingle {
		t.Errorf("ApprovalNodeConfig.Mode = %v, want %v", config.Mode, node.ApprovalModeSingle)
	}
	if config.Timeout == nil || *config.Timeout != timeout {
		t.Error("ApprovalNodeConfig.Timeout should be set correctly")
	}
	if config.RejectBehavior != node.RejectBehaviorTerminate {
		t.Errorf("ApprovalNodeConfig.RejectBehavior = %v, want %v", config.RejectBehavior, node.RejectBehaviorTerminate)
	}
	if !config.RequireComment() {
		t.Error("ApprovalNodeConfig.RequireComment() should be true")
	}
	if config.RequireAttachments() {
		t.Error("ApprovalNodeConfig.RequireAttachments() should be false")
	}
}

// TestApprovalNodeConfigNodeType 验证节点类型
func TestApprovalNodeConfigNodeType(t *testing.T) {
	config := &node.ApprovalNodeConfig{
		Mode: node.ApprovalModeSingle,
	}

	if config.NodeType() != template.NodeTypeApproval {
		t.Errorf("ApprovalNodeConfig.NodeType() = %v, want %v", config.NodeType(), template.NodeTypeApproval)
	}
}

// TestApprovalNodeConfigValidate 测试配置验证
func TestApprovalNodeConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *node.ApprovalNodeConfig
		wantErr bool
	}{
		{
			name: "valid config with fixed approvers",
			config: &node.ApprovalNodeConfig{
				Mode: node.ApprovalModeSingle,
				ApproverConfig: &node.FixedApproverConfig{
					Approvers: []string{"user-001"},
				},
			},
			wantErr: false,
		},
		{
			name: "missing approver config",
			config: &node.ApprovalNodeConfig{
				Mode:          node.ApprovalModeSingle,
				ApproverConfig: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ApprovalNodeConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

