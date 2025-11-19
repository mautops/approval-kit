package node_test

import (
	"testing"

	"github.com/mautops/approval-kit/internal/node"
)

// TestApproverTimingOnCreate 测试任务创建时获取审批人
func TestApproverTimingOnCreate(t *testing.T) {
	config := &node.DynamicApproverConfig{
		API: &node.HTTPAPIConfig{
			URL:    "http://example.com/api/approvers",
			Method: "POST",
			ResponseMapping: &node.ResponseMapping{
				Path:   "approvers",
				Format: "json",
			},
		},
		Timing: node.ApproverTimingOnCreate,
	}

	// 验证获取时机
	timing := config.GetTiming()
	if timing != node.ApproverTimingOnCreate {
		t.Errorf("DynamicApproverConfig.GetTiming() = %v, want %v", timing, node.ApproverTimingOnCreate)
	}
}

// TestApproverTimingOnActivate 测试节点激活时获取审批人
func TestApproverTimingOnActivate(t *testing.T) {
	config := &node.DynamicApproverConfig{
		API: &node.HTTPAPIConfig{
			URL:    "http://example.com/api/approvers",
			Method: "POST",
			ResponseMapping: &node.ResponseMapping{
				Path:   "approvers",
				Format: "json",
			},
		},
		Timing: node.ApproverTimingOnActivate,
	}

	// 验证获取时机
	timing := config.GetTiming()
	if timing != node.ApproverTimingOnActivate {
		t.Errorf("DynamicApproverConfig.GetTiming() = %v, want %v", timing, node.ApproverTimingOnActivate)
	}
}

// TestApproverTimingDefault 测试默认获取时机
func TestApproverTimingDefault(t *testing.T) {
	config := &node.DynamicApproverConfig{
		API: &node.HTTPAPIConfig{
			URL:    "http://example.com/api/approvers",
			Method: "POST",
			ResponseMapping: &node.ResponseMapping{
				Path:   "approvers",
				Format: "json",
			},
		},
		// Timing 未设置,应该默认为节点激活时
	}

	// 验证默认获取时机
	timing := config.GetTiming()
	if timing != node.ApproverTimingOnActivate {
		t.Errorf("DynamicApproverConfig.GetTiming() = %v, want %v (default)", timing, node.ApproverTimingOnActivate)
	}
}


// TestApproverTimingIntegration 测试获取时机在任务创建时的集成场景
func TestApproverTimingIntegration(t *testing.T) {
	// 这个测试验证在任务创建时,如果配置了 on_create 时机,应该能够获取审批人
	// 注意: 实际的集成测试需要在 TaskManager.Create 中实现
	config := &node.DynamicApproverConfig{
		API: &node.HTTPAPIConfig{
			URL:    "http://example.com/api/approvers",
			Method: "POST",
			ResponseMapping: &node.ResponseMapping{
				Path:   "approvers",
				Format: "json",
			},
		},
		Timing: node.ApproverTimingOnCreate,
	}

	// 验证配置正确
	if config.GetTiming() != node.ApproverTimingOnCreate {
		t.Errorf("Expected timing to be on_create, got %v", config.GetTiming())
	}

	// 注意: 实际的获取逻辑需要在 TaskManager.Create 中实现
	// 这里只验证配置是否正确
}


