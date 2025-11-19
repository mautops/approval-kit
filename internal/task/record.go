package task

import (
	"fmt"
	"time"
)

// Validate 验证审批记录的有效性
func (r *Record) Validate() error {
	if r.ID == "" {
		return fmt.Errorf("record ID is required")
	}
	if r.TaskID == "" {
		return fmt.Errorf("record TaskID is required")
	}
	if r.NodeID == "" {
		return fmt.Errorf("record NodeID is required")
	}
	if r.Approver == "" {
		return fmt.Errorf("record Approver is required")
	}
	if r.Result == "" {
		return fmt.Errorf("record Result is required")
	}
	// 验证审批结果类型
	validResults := []string{"approve", "reject", "transfer", "add_approver", "remove_approver"}
	valid := false
	for _, validResult := range validResults {
		if r.Result == validResult {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid record Result: %q, must be one of: approve, reject, transfer, add_approver, remove_approver", r.Result)
	}
	// 验证时间
	if r.CreatedAt.IsZero() {
		return fmt.Errorf("record CreatedAt is required")
	}
	return nil
}

// NewRecord 创建新的审批记录
func NewRecord(id, taskID, nodeID, approver, result, comment string, attachments []string) *Record {
	return &Record{
		ID:          id,
		TaskID:     taskID,
		NodeID:     nodeID,
		Approver:   approver,
		Result:     result,
		Comment:    comment,
		CreatedAt:  time.Now(),
		Attachments: attachments,
	}
}

