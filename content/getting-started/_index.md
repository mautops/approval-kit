---
title: "快速开始"
description: "安装和使用 Approval Kit 的快速指南"
date: 2025-11-19
draft: false
weight: 10
---

## 安装

### 使用 Go Modules

```bash
go get github.com/mautops/approval-kit
```

### 版本要求

- Go 1.25.4 或更高版本

## 基本使用

### 1. 创建管理器

```go
import (
    "github.com/mautops/approval-kit/internal/task"
    "github.com/mautops/approval-kit/internal/template"
)

// 创建模板管理器
templateMgr := template.NewTemplateManager()

// 创建任务管理器
taskMgr := task.NewTaskManager(templateMgr, nil)
```

### 2. 创建审批模板

```go
import (
    "time"
    "github.com/mautops/approval-kit/internal/node"
    "github.com/mautops/approval-kit/internal/template"
)

now := time.Now()
tpl := &template.Template{
    ID:          "simple-approval",
    Name:        "简单审批模板",
    Description: "单人审批流程",
    Version:     1,
    CreatedAt:   now,
    UpdatedAt:   now,
    Nodes: map[string]*template.Node{
        "start": {
            ID:   "start",
            Name: "开始",
            Type: template.NodeTypeStart,
        },
        "approval": {
            ID:   "approval",
            Name: "审批节点",
            Type: template.NodeTypeApproval,
            Config: &node.ApprovalNodeConfig{
                Mode: node.ApprovalModeSingle,
                ApproverConfig: &node.FixedApproverConfig{
                    Approvers: []string{"approver-001"},
                },
            },
        },
        "end": {
            ID:   "end",
            Name: "结束",
            Type: template.NodeTypeEnd,
        },
    },
    Edges: []*template.Edge{
        {From: "start", To: "approval"},
        {From: "approval", To: "end"},
    },
}

err := templateMgr.Create(tpl)
```

### 3. 创建审批任务

```go
import (
    "encoding/json"
    "github.com/mautops/approval-kit/internal/task"
)

params := json.RawMessage(`{"amount": 1000}`)
tsk, err := taskMgr.Create("simple-approval", "business-001", params)
```

### 4. 提交和审批

```go
// 提交任务
err = taskMgr.Submit(tsk.ID)

// 审批通过
err = taskMgr.Approve(tsk.ID, "approval", "approver-001", "同意")
```

## 下一步

- 了解[核心概念](/concepts/) - 深入理解模板、任务和状态机
- 查看[使用示例](/examples/) - 25 个实际使用场景
- 阅读[API 参考](/api/) - 完整的 API 文档

