package node_test

import (
	"testing"

	"github.com/mautops/approval-kit/internal/node"
)

// TestRejectJumpConfig 测试拒绝后跳转配置
func TestRejectJumpConfig(t *testing.T) {
	config := &node.ApprovalNodeConfig{
		Mode:           node.ApprovalModeSingle,
		RejectBehavior: node.RejectBehaviorJump,
		ApproverConfig: &node.FixedApproverConfig{
			Approvers: []string{"user-001"},
		},
	}

	if config.RejectBehavior != node.RejectBehaviorJump {
		t.Errorf("RejectBehavior = %q, want %q", config.RejectBehavior, node.RejectBehaviorJump)
	}
}

// TestRejectJumpConfigWithTargetNode 测试拒绝后跳转到指定节点
func TestRejectJumpConfigWithTargetNode(t *testing.T) {
	config := &node.ApprovalNodeConfig{
		Mode:            node.ApprovalModeSingle,
		RejectBehavior:  node.RejectBehaviorJump,
		RejectTargetNode: "node-002",
		ApproverConfig: &node.FixedApproverConfig{
			Approvers: []string{"user-001"},
		},
	}

	if config.RejectBehavior != node.RejectBehaviorJump {
		t.Errorf("RejectBehavior = %q, want %q", config.RejectBehavior, node.RejectBehaviorJump)
	}

	if config.RejectTargetNode != "node-002" {
		t.Errorf("RejectTargetNode = %q, want %q", config.RejectTargetNode, "node-002")
	}

	// 验证配置
	err := config.Validate()
	if err != nil {
		t.Errorf("Validate() should pass for valid jump config, got error: %v", err)
	}
}

// TestRejectJumpConfigWithoutTargetNode 测试拒绝后跳转但未指定目标节点(应该验证失败)
func TestRejectJumpConfigWithoutTargetNode(t *testing.T) {
	config := &node.ApprovalNodeConfig{
		Mode:            node.ApprovalModeSingle,
		RejectBehavior:  node.RejectBehaviorJump,
		RejectTargetNode: "", // 未指定目标节点
		ApproverConfig: &node.FixedApproverConfig{
			Approvers: []string{"user-001"},
		},
	}

	// 验证应该失败
	err := config.Validate()
	if err == nil {
		t.Error("Validate() should fail when RejectBehavior is Jump but RejectTargetNode is empty")
	}
}

// TestRejectRollbackConfig 测试拒绝后回滚配置
func TestRejectRollbackConfig(t *testing.T) {
	config := &node.ApprovalNodeConfig{
		Mode:           node.ApprovalModeSingle,
		RejectBehavior: node.RejectBehaviorRollback,
		ApproverConfig: &node.FixedApproverConfig{
			Approvers: []string{"user-001"},
		},
	}

	if config.RejectBehavior != node.RejectBehaviorRollback {
		t.Errorf("RejectBehavior = %q, want %q", config.RejectBehavior, node.RejectBehaviorRollback)
	}
}

// TestRejectTerminateConfig 测试拒绝后终止配置
func TestRejectTerminateConfig(t *testing.T) {
	config := &node.ApprovalNodeConfig{
		Mode:           node.ApprovalModeSingle,
		RejectBehavior: node.RejectBehaviorTerminate,
		ApproverConfig: &node.FixedApproverConfig{
			Approvers: []string{"user-001"},
		},
	}

	if config.RejectBehavior != node.RejectBehaviorTerminate {
		t.Errorf("RejectBehavior = %q, want %q", config.RejectBehavior, node.RejectBehaviorTerminate)
	}
}

