package task_test

import (
	"encoding/json"
	"testing"

	"github.com/mautops/approval-kit/internal/task"
)

// TestTaskManagerInterface 验证 TaskManager 接口定义
func TestTaskManagerInterface(t *testing.T) {
	// 验证接口类型存在
	var mgr task.TaskManager
	if mgr != nil {
		_ = mgr
	}
}

// TestTaskManagerMethods 验证 TaskManager 接口方法签名
func TestTaskManagerMethods(t *testing.T) {
	// 验证接口包含所有必需的方法
	// 通过编译时检查,如果方法不存在会编译失败
	var _ task.TaskManager = (*taskManagerImpl)(nil)
}

// taskManagerImpl 用于测试接口方法签名的实现
type taskManagerImpl struct{}

func (m *taskManagerImpl) Create(templateID string, businessID string, params json.RawMessage) (*task.Task, error) {
	return nil, nil
}

func (m *taskManagerImpl) Get(id string) (*task.Task, error) {
	return nil, nil
}

func (m *taskManagerImpl) Submit(id string) error {
	return nil
}

func (m *taskManagerImpl) Approve(id string, nodeID string, approver string, comment string) error {
	return nil
}

func (m *taskManagerImpl) Reject(id string, nodeID string, approver string, comment string) error {
	return nil
}

func (m *taskManagerImpl) Cancel(id string, reason string) error {
	return nil
}

func (m *taskManagerImpl) Query(filter *task.TaskFilter) ([]*task.Task, error) {
	return nil, nil
}

func (m *taskManagerImpl) ApproveWithAttachments(id string, nodeID string, approver string, comment string, attachments []string) error {
	return nil
}

func (m *taskManagerImpl) RejectWithAttachments(id string, nodeID string, approver string, comment string, attachments []string) error {
	return nil
}

func (m *taskManagerImpl) Withdraw(id string, reason string) error {
	return nil
}

func (m *taskManagerImpl) Transfer(id string, nodeID string, fromApprover string, toApprover string, reason string) error {
	return nil
}

func (m *taskManagerImpl) AddApprover(id string, nodeID string, approver string, reason string) error {
	return nil
}

func (m *taskManagerImpl) RemoveApprover(id string, nodeID string, approver string, reason string) error {
	return nil
}

func (m *taskManagerImpl) HandleTimeout(id string) error {
	return nil
}

