package task

import (
	"encoding/json"
)

// Clone 创建任务的深拷贝
// 用于确保任务对象的隔离性
func (t *Task) Clone() *Task {
	if t == nil {
		return nil
	}

	// 创建新任务对象
	clone := &Task{
		ID:             t.ID,
		TemplateID:     t.TemplateID,
		TemplateVersion: t.TemplateVersion,
		BusinessID:     t.BusinessID,
		State:          t.State,
		CurrentNode:    t.CurrentNode,
		PausedState:    t.PausedState,
		CreatedAt:      t.CreatedAt,
		UpdatedAt:      t.UpdatedAt,
	}

	// 复制 Params
	if t.Params != nil {
		clone.Params = make(json.RawMessage, len(t.Params))
		copy(clone.Params, t.Params)
	}

	// 复制 SubmittedAt
	if t.SubmittedAt != nil {
		submittedAt := *t.SubmittedAt
		clone.SubmittedAt = &submittedAt
	}

	// 复制 PausedAt
	if t.PausedAt != nil {
		pausedAt := *t.PausedAt
		clone.PausedAt = &pausedAt
	}

	// 复制 NodeOutputs
	clone.NodeOutputs = make(map[string]json.RawMessage)
	for k, v := range t.NodeOutputs {
		output := make(json.RawMessage, len(v))
		copy(output, v)
		clone.NodeOutputs[k] = output
	}

	// 复制 Approvers
	clone.Approvers = make(map[string][]string)
	for k, v := range t.Approvers {
		approvers := make([]string, len(v))
		copy(approvers, v)
		clone.Approvers[k] = approvers
	}

	// 复制 Approvals
	clone.Approvals = make(map[string]map[string]*Approval)
	for k, v := range t.Approvals {
		approvals := make(map[string]*Approval)
		for k2, v2 := range v {
			approvals[k2] = &Approval{
				Result:    v2.Result,
				Comment:   v2.Comment,
				CreatedAt: v2.CreatedAt,
			}
		}
		clone.Approvals[k] = approvals
	}

	// 复制 CompletedNodes
	clone.CompletedNodes = make([]string, len(t.CompletedNodes))
	copy(clone.CompletedNodes, t.CompletedNodes)

	// 复制 Records
	clone.Records = make([]*Record, len(t.Records))
	for i, r := range t.Records {
		clone.Records[i] = &Record{
			ID:         r.ID,
			TaskID:     r.TaskID,
			NodeID:     r.NodeID,
			Approver:   r.Approver,
			Result:     r.Result,
			Comment:    r.Comment,
			CreatedAt:  r.CreatedAt,
			Attachments: make([]string, len(r.Attachments)),
		}
		copy(clone.Records[i].Attachments, r.Attachments)
	}

	// 复制 StateHistory
	clone.StateHistory = make([]*StateChange, len(t.StateHistory))
	for i, sc := range t.StateHistory {
		clone.StateHistory[i] = &StateChange{
			From:   sc.From,
			To:     sc.To,
			Reason: sc.Reason,
			Time:   sc.Time,
		}
	}

	return clone
}

// cloneRecord 创建审批记录的副本
func (r *Record) cloneRecord() *Record {
	// 复制附件列表
	attachments := make([]string, len(r.Attachments))
	copy(attachments, r.Attachments)

	return &Record{
		ID:          r.ID,
		TaskID:      r.TaskID,
		NodeID:      r.NodeID,
		Approver:    r.Approver,
		Result:      r.Result,
		Comment:     r.Comment,
		CreatedAt:   r.CreatedAt,
		Attachments: attachments,
	}
}

