package task_test

import (
	"encoding/json"
	"testing"

	pkgTask "github.com/mautops/approval-kit/pkg/task"
	internalTask "github.com/mautops/approval-kit/internal/task"
)

// TestPkgTaskManagerInterface 验证 pkg/task 包中的 TaskManager 接口定义
func TestPkgTaskManagerInterface(t *testing.T) {
	// 验证接口类型存在
	var mgr pkgTask.TaskManager
	if mgr != nil {
		_ = mgr
	}
}

// TestPkgTaskManagerCompatibility 验证 pkg/task.TaskManager 与 internal/task.TaskManager 兼容
func TestPkgTaskManagerCompatibility(t *testing.T) {
	// 验证 pkg 接口可以被 internal 实现满足
	var _ pkgTask.TaskManager = (*internalTaskManagerAdapter)(nil)
}

// internalTaskManagerAdapter 用于测试接口兼容性的适配器
type internalTaskManagerAdapter struct {
	impl internalTask.TaskManager
}

func (a *internalTaskManagerAdapter) Create(templateID string, businessID string, params json.RawMessage) (*pkgTask.Task, error) {
	task, err := a.impl.Create(templateID, businessID, params)
	if err != nil {
		return nil, err
	}
	// 类型转换将在后续任务中实现
	return pkgTask.FromInternal(task), nil
}

func (a *internalTaskManagerAdapter) Get(id string) (*pkgTask.Task, error) {
	task, err := a.impl.Get(id)
	if err != nil {
		return nil, err
	}
	return pkgTask.FromInternal(task), nil
}

func (a *internalTaskManagerAdapter) Submit(id string) error {
	return a.impl.Submit(id)
}

func (a *internalTaskManagerAdapter) Approve(id string, nodeID string, approver string, comment string) error {
	return a.impl.Approve(id, nodeID, approver, comment)
}

func (a *internalTaskManagerAdapter) ApproveWithAttachments(id string, nodeID string, approver string, comment string, attachments []string) error {
	return a.impl.ApproveWithAttachments(id, nodeID, approver, comment, attachments)
}

func (a *internalTaskManagerAdapter) Reject(id string, nodeID string, approver string, comment string) error {
	return a.impl.Reject(id, nodeID, approver, comment)
}

func (a *internalTaskManagerAdapter) RejectWithAttachments(id string, nodeID string, approver string, comment string, attachments []string) error {
	return a.impl.RejectWithAttachments(id, nodeID, approver, comment, attachments)
}

func (a *internalTaskManagerAdapter) Cancel(id string, reason string) error {
	return a.impl.Cancel(id, reason)
}

func (a *internalTaskManagerAdapter) Withdraw(id string, reason string) error {
	return a.impl.Withdraw(id, reason)
}

func (a *internalTaskManagerAdapter) Transfer(id string, nodeID string, fromApprover string, toApprover string, reason string) error {
	return a.impl.Transfer(id, nodeID, fromApprover, toApprover, reason)
}

func (a *internalTaskManagerAdapter) AddApprover(id string, nodeID string, approver string, reason string) error {
	return a.impl.AddApprover(id, nodeID, approver, reason)
}

func (a *internalTaskManagerAdapter) RemoveApprover(id string, nodeID string, approver string, reason string) error {
	return a.impl.RemoveApprover(id, nodeID, approver, reason)
}

func (a *internalTaskManagerAdapter) Query(filter *pkgTask.TaskFilter) ([]*pkgTask.Task, error) {
	internalFilter := pkgTask.TaskFilterToInternal(filter)
	tasks, err := a.impl.Query(internalFilter)
	if err != nil {
		return nil, err
	}
	result := make([]*pkgTask.Task, len(tasks))
	for i, task := range tasks {
		result[i] = pkgTask.FromInternal(task)
	}
	return result, nil
}

func (a *internalTaskManagerAdapter) HandleTimeout(id string) error {
	return a.impl.HandleTimeout(id)
}

func (a *internalTaskManagerAdapter) Pause(id string, reason string) error {
	return a.impl.Pause(id, reason)
}

func (a *internalTaskManagerAdapter) Resume(id string, reason string) error {
	return a.impl.Resume(id, reason)
}

func (a *internalTaskManagerAdapter) RollbackToNode(id string, nodeID string, reason string) error {
	return a.impl.RollbackToNode(id, nodeID, reason)
}

func (a *internalTaskManagerAdapter) ReplaceApprover(id string, nodeID string, oldApprover string, newApprover string, reason string) error {
	return a.impl.ReplaceApprover(id, nodeID, oldApprover, newApprover, reason)
}

