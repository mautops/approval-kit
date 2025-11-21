package task

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/mautops/approval-kit/internal/types"
)

// Task 表示审批任务实例,包含状态和运行时数据
// 使用读写锁保证并发安全
type Task struct {
	mu sync.RWMutex // 读写锁,保证并发安全

	// 基本信息
	ID             string          // 任务 ID
	TemplateID     string          // 模板 ID
	TemplateVersion int            // 模板版本号
	BusinessID     string          // 关联的业务 ID
	Params         json.RawMessage // 任务参数(JSON 格式)

	// 状态信息
	State       types.TaskState // 当前状态
	CurrentNode string           // 当前节点 ID

	// 暂停相关字段
	PausedAt    *time.Time      // 暂停时间
	PausedState types.TaskState // 暂停前的状态,用于恢复时恢复到正确状态

	// 时间信息
	CreatedAt   time.Time  // 创建时间
	UpdatedAt   time.Time  // 更新时间
	SubmittedAt *time.Time // 提交时间

	// 运行时数据
	NodeOutputs map[string]json.RawMessage           // 节点 ID -> 节点输出数据
	Approvers   map[string][]string                  // 节点 ID -> 审批人列表
	Approvals   map[string]map[string]*Approval      // 节点 ID -> 审批人 -> 审批结果

	// 回退相关字段
	CompletedNodes []string // 已完成的节点 ID 列表,用于回退操作

	// 审批记录
	Records []*Record // 审批记录列表

	// 状态变更历史
	StateHistory []*StateChange // 状态变更历史
}

// Approval 审批结果
type Approval struct {
	Result    string    // 审批结果(approve/reject/transfer)
	Comment   string    // 审批意见
	CreatedAt time.Time // 审批时间
}

// Record 审批记录
// 每次审批操作时自动生成,作为状态流转的副产品
type Record struct {
	ID         string    // 记录 ID
	TaskID     string    // 任务 ID
	NodeID     string    // 节点 ID
	Approver   string    // 审批人
	Result     string    // 审批结果(approve/reject/transfer)
	Comment    string    // 审批意见
	CreatedAt  time.Time // 审批时间
	Attachments []string // 附件列表
}

// StateChange 状态变更记录
// 记录每次状态变更的详细信息,用于追溯和审计
type StateChange struct {
	From    types.TaskState // 源状态
	To      types.TaskState // 目标状态
	Reason  string          // 转换原因
	Time    time.Time       // 转换时间
}
