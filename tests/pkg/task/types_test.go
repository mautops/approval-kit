package task_test

import (
	"testing"
	"time"

	pkgTask "github.com/mautops/approval-kit/pkg/task"
	internalTask "github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/pkg/types"
)

// TestPkgApprovalType 验证 pkg/task.Approval 类型暴露
func TestPkgApprovalType(t *testing.T) {
	var approval pkgTask.Approval
	approval.Result = "approve"
	approval.Comment = "test"
	approval.CreatedAt = time.Now()

	if approval.Result != "approve" {
		t.Errorf("Approval.Result = %v, want approve", approval.Result)
	}
}

// TestPkgRecordType 验证 pkg/task.Record 类型暴露
func TestPkgRecordType(t *testing.T) {
	var record pkgTask.Record
	record.ID = "record-1"
	record.TaskID = "task-1"
	record.NodeID = "node-1"
	record.Approver = "user-1"
	record.Result = "approve"
	record.Comment = "test"
	record.CreatedAt = time.Now()
	record.Attachments = []string{"file1.pdf"}

	if record.ID != "record-1" {
		t.Errorf("Record.ID = %v, want record-1", record.ID)
	}
}

// TestPkgStateChangeType 验证 pkg/task.StateChange 类型暴露
func TestPkgStateChangeType(t *testing.T) {
	var stateChange pkgTask.StateChange
	stateChange.From = types.TaskStatePending
	stateChange.To = types.TaskStateSubmitted
	stateChange.Reason = "test"
	stateChange.Time = time.Now()

	if stateChange.From != types.TaskStatePending {
		t.Errorf("StateChange.From = %v, want pending", stateChange.From)
	}
}

// TestPkgTypesCompatibility 验证 pkg 类型与 internal 类型兼容
func TestPkgTypesCompatibility(t *testing.T) {
	// 验证类型别名可以互相转换
	internalApproval := &internalTask.Approval{
		Result:    "approve",
		Comment:   "test",
		CreatedAt: time.Now(),
	}

	var pkgApproval pkgTask.Approval = *internalApproval
	if pkgApproval.Result != "approve" {
		t.Errorf("Type alias conversion failed")
	}

	internalRecord := &internalTask.Record{
		ID:        "record-1",
		TaskID:    "task-1",
		NodeID:    "node-1",
		Approver:  "user-1",
		Result:    "approve",
		Comment:   "test",
		CreatedAt: time.Now(),
	}

	var pkgRecord pkgTask.Record = *internalRecord
	if pkgRecord.ID != "record-1" {
		t.Errorf("Type alias conversion failed")
	}

	internalStateChange := &internalTask.StateChange{
		From:   types.TaskStatePending,
		To:     types.TaskStateSubmitted,
		Reason: "test",
		Time:   time.Now(),
	}

	var pkgStateChange pkgTask.StateChange = *internalStateChange
	if pkgStateChange.From != types.TaskStatePending {
		t.Errorf("Type alias conversion failed")
	}
}

