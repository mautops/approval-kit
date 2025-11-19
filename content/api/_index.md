---
title: "API 参考"
description: "Approval Kit 完整的 API 文档"
date: 2025-11-19
draft: false
weight: 10
---

## 模板管理 (TemplateManager)

### Create

创建审批模板.

```go
func (m *TemplateManager) Create(template *Template) error
```

### Get

获取指定版本的模板.

```go
func (m *TemplateManager) Get(id string, version int) (*Template, error)
```

### Update

更新模板,自动创建新版本.

```go
func (m *TemplateManager) Update(id string, template *Template) error
```

### Delete

删除模板.

```go
func (m *TemplateManager) Delete(id string) error
```

### ListVersions

列出模板的所有版本.

```go
func (m *TemplateManager) ListVersions(id string) ([]int, error)
```

## 任务管理 (TaskManager)

### Create

基于模板创建审批任务.

```go
func (m *TaskManager) Create(templateID string, businessID string, params json.RawMessage) (*Task, error)
```

### Get

获取任务详情.

```go
func (m *TaskManager) Get(id string) (*Task, error)
```

### Submit

提交任务进入审批流程.

```go
func (m *TaskManager) Submit(id string) error
```

### Approve

审批通过.

```go
func (m *TaskManager) Approve(id string, nodeID string, approver string, comment string) error
```

### Reject

审批拒绝.

```go
func (m *TaskManager) Reject(id string, nodeID string, approver string, comment string) error
```

### Cancel

取消任务.

```go
func (m *TaskManager) Cancel(id string, reason string) error
```

### Withdraw

撤回任务(仅限未审批状态).

```go
func (m *TaskManager) Withdraw(id string) error
```

### Query

查询任务列表.

```go
func (m *TaskManager) Query(filter *TaskFilter) ([]*Task, error)
```

支持按状态、模板、业务 ID、审批人、时间范围等条件查询.

## 状态机 (StateMachine)

### CanTransition

检查是否允许状态转换.

```go
func (f *fsm) CanTransition(from TaskState, to TaskState) bool
```

### Transition

执行状态转换.

```go
func (f *fsm) Transition(task TransitionableTask, to TaskState, reason string) (TransitionableTask, error)
```

### GetValidTransitions

获取当前状态的有效转换.

```go
func (f *fsm) GetValidTransitions(state TaskState) []TaskState
```

## 数据结构

### Template

```go
type Template struct {
    ID          string
    Name        string
    Description string
    Version     int
    CreatedAt   time.Time
    UpdatedAt   time.Time
    Nodes       map[string]*Node
    Edges       []*Edge
    Config      *TemplateConfig
}
```

### Task

```go
type Task struct {
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
```

### TaskState

```go
const (
    TaskStatePending   TaskState = "pending"
    TaskStateSubmitted TaskState = "submitted"
    TaskStateApproving TaskState = "approving"
    TaskStateApproved  TaskState = "approved"
    TaskStateRejected  TaskState = "rejected"
    TaskStateCancelled TaskState = "cancelled"
    TaskStateTimeout   TaskState = "timeout"
)
```

