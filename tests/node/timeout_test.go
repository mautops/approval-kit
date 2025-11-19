package node_test

import (
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/node"
)

// TestTimeoutConfig 测试超时配置
func TestTimeoutConfig(t *testing.T) {
	timeout := 30 * time.Minute
	config := &node.ApprovalNodeConfig{
		Mode:    node.ApprovalModeSingle,
		Timeout: &timeout,
	}

	if config.Timeout == nil {
		t.Error("Timeout should be set")
	}

	if *config.Timeout != timeout {
		t.Errorf("Timeout = %v, want %v", *config.Timeout, timeout)
	}
}

// TestTimeoutConfigNil 测试空超时配置
func TestTimeoutConfigNil(t *testing.T) {
	config := &node.ApprovalNodeConfig{
		Mode:    node.ApprovalModeSingle,
		Timeout: nil,
	}

	if config.Timeout != nil {
		t.Error("Timeout should be nil when not configured")
	}
}

// TestTimeoutConfigValidation 测试超时配置验证
func TestTimeoutConfigValidation(t *testing.T) {
	// 测试有效的超时配置
	timeout := 30 * time.Minute
	config := &node.ApprovalNodeConfig{
		Mode:    node.ApprovalModeSingle,
		Timeout: &timeout,
		ApproverConfig: &node.FixedApproverConfig{
			Approvers: []string{"user-001"},
		},
	}

	err := config.Validate()
	if err != nil {
		t.Errorf("Validate() should pass for valid timeout config, got error: %v", err)
	}

	// 测试零值超时(应该无效)
	zeroTimeout := time.Duration(0)
	config.Timeout = &zeroTimeout
	err = config.Validate()
	if err == nil {
		t.Error("Validate() should fail for zero timeout")
	}
}
