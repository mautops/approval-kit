package task

import (
	"encoding/json"
	"time"

	"github.com/mautops/approval-kit/internal/types"
)

// GetState 并发安全地获取任务状态
func (t *Task) GetState() types.TaskState {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.State
}

// GetCurrentNode 并发安全地获取当前节点 ID
func (t *Task) GetCurrentNode() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.CurrentNode
}

// GetUpdatedAt 并发安全地获取更新时间
func (t *Task) GetUpdatedAt() time.Time {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.UpdatedAt
}

// GetStateHistory 并发安全地获取状态变更历史
func (t *Task) GetStateHistory() []*StateChange {
	t.mu.RLock()
	defer t.mu.RUnlock()
	// 返回副本,避免外部修改
	result := make([]*StateChange, len(t.StateHistory))
	copy(result, t.StateHistory)
	return result
}

// GetRecords 并发安全地获取审批记录
// 返回记录的副本,确保隔离性
func (t *Task) GetRecords() []*Record {
	t.mu.RLock()
	defer t.mu.RUnlock()
	// 返回深拷贝,避免外部修改
	result := make([]*Record, len(t.Records))
	for i, record := range t.Records {
		result[i] = record.cloneRecord()
	}
	return result
}

// GetRecordsByNodeID 按节点 ID 获取审批记录
// nodeID: 节点 ID
// 返回: 该节点的所有审批记录(副本)
func (t *Task) GetRecordsByNodeID(nodeID string) []*Record {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var records []*Record
	for _, record := range t.Records {
		if record.NodeID == nodeID {
			records = append(records, record.cloneRecord())
		}
	}
	return records
}

// Update 并发安全地更新任务
// 使用函数式更新模式,确保原子性
func (t *Task) Update(fn func(*Task) error) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return fn(t)
}

// Snapshot 创建任务快照(并发安全)
// 返回任务的不可变快照,用于读取操作
func (t *Task) Snapshot() *TaskSnapshot {
	t.mu.RLock()
	defer t.mu.RUnlock()

	// 复制所有字段
	snapshot := &TaskSnapshot{
		ID:             t.ID,
		TemplateID:     t.TemplateID,
		TemplateVersion: t.TemplateVersion,
		BusinessID:     t.BusinessID,
		State:          t.State,
		CurrentNode:    t.CurrentNode,
		CreatedAt:      t.CreatedAt,
		UpdatedAt:      t.UpdatedAt,
		SubmittedAt:    t.SubmittedAt,
	}

	// 复制 Params
	if t.Params != nil {
		snapshot.Params = make(json.RawMessage, len(t.Params))
		copy(snapshot.Params, t.Params)
	}

	// 复制 NodeOutputs
	snapshot.NodeOutputs = make(map[string]json.RawMessage)
	for k, v := range t.NodeOutputs {
		output := make(json.RawMessage, len(v))
		copy(output, v)
		snapshot.NodeOutputs[k] = output
	}

	// 复制 Approvers
	snapshot.Approvers = make(map[string][]string)
	for k, v := range t.Approvers {
		approvers := make([]string, len(v))
		copy(approvers, v)
		snapshot.Approvers[k] = approvers
	}

	// 复制 Approvals
	snapshot.Approvals = make(map[string]map[string]*Approval)
	for k, v := range t.Approvals {
		approvals := make(map[string]*Approval)
		for k2, v2 := range v {
			approvals[k2] = &Approval{
				Result:    v2.Result,
				Comment:   v2.Comment,
				CreatedAt: v2.CreatedAt,
			}
		}
		snapshot.Approvals[k] = approvals
	}

	// 复制 Records
	snapshot.Records = make([]*Record, len(t.Records))
	for i, r := range t.Records {
		snapshot.Records[i] = &Record{
			ID:         r.ID,
			TaskID:     r.TaskID,
			NodeID:     r.NodeID,
			Approver:   r.Approver,
			Result:     r.Result,
			Comment:    r.Comment,
			CreatedAt:  r.CreatedAt,
			Attachments: make([]string, len(r.Attachments)),
		}
		copy(snapshot.Records[i].Attachments, r.Attachments)
	}

	// 复制 StateHistory
	snapshot.StateHistory = make([]*StateChange, len(t.StateHistory))
	for i, sc := range t.StateHistory {
		snapshot.StateHistory[i] = &StateChange{
			From:   sc.From,
			To:     sc.To,
			Reason: sc.Reason,
			Time:   sc.Time,
		}
	}

	return snapshot
}

// TaskSnapshot 任务快照
// 不可变的任务数据副本,用于安全读取
type TaskSnapshot struct {
	ID             string
	TemplateID     string
	TemplateVersion int
	BusinessID     string
	Params         json.RawMessage
	State          TaskState
	CurrentNode    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	SubmittedAt    *time.Time
	NodeOutputs    map[string]json.RawMessage
	Approvers      map[string][]string
	Approvals      map[string]map[string]*Approval
	Records        []*Record
	StateHistory   []*StateChange
}

