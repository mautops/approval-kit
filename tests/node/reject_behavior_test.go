package node_test

import (
	"testing"

	"github.com/mautops/approval-kit/internal/node"
)

// TestRejectBehavior 测试 RejectBehavior 类型
func TestRejectBehavior(t *testing.T) {
	behaviors := []node.RejectBehavior{
		node.RejectBehaviorTerminate,
		node.RejectBehaviorRollback,
		node.RejectBehaviorJump,
	}

	for _, behavior := range behaviors {
		if behavior == "" {
			t.Errorf("RejectBehavior should not be empty")
		}
	}
}

// TestRejectBehaviorTerminate 测试终止行为
func TestRejectBehaviorTerminate(t *testing.T) {
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

// TestRejectBehaviorRollback 测试回滚行为
func TestRejectBehaviorRollback(t *testing.T) {
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

// TestRejectBehaviorJump 测试跳转行为
func TestRejectBehaviorJump(t *testing.T) {
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
