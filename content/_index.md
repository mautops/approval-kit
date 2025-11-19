---
title: "Approval Kit"
description: "审批流核心库,专注于管理审批模板和审批任务的状态流转"
date: 2025-11-19
draft: false
---

## 简介

Approval Kit 是一个业务无关的审批流核心库,专注于管理审批模板和审批任务的状态流转.本库不包含任何数据持久化逻辑,通过接口抽象与上层业务系统解耦.

## 核心特性

- **状态机管理**: 基于有限状态机(FSM)确保状态转换的合法性和一致性
- **多种审批模式**: 支持单人审批、会签、或签、比例会签、顺序审批
- **灵活的条件分支**: 支持数值、字符串、枚举等多种条件类型
- **动态审批人**: 支持通过 HTTP API 动态获取审批人
- **事件通知**: 支持 Webhook 事件通知机制
- **模板版本控制**: 支持模板版本管理,确保历史任务可追溯
- **并发安全**: 使用读写锁保证并发访问安全
- **无外部依赖**: 优先使用 Go 标准库,最小化外部依赖

## 设计原则

- **业务无关性**: 不包含任何具体业务逻辑,仅负责状态流转
- **无数据存储**: 所有数据存储由上层业务系统负责
- **测试优先**: 采用 TDD 开发,关键路径测试覆盖率 100%
- **简洁设计**: 遵循 YAGNI 原则,避免过度设计
- **API 稳定**: 遵循语义化版本控制,保证向后兼容

## 快速开始

```go
package main

import (
    "github.com/mautops/approval-kit/internal/task"
    "github.com/mautops/approval-kit/internal/template"
)

func main() {
    // 创建管理器
    templateMgr := template.NewTemplateManager()
    taskMgr := task.NewTaskManager(templateMgr, nil)
    
    // 创建模板和任务...
}
```

## 文档导航

- [快速开始](/getting-started/) - 安装和基本使用
- [核心概念](/concepts/) - 模板、任务、状态机等核心概念
- [API 参考](/api/) - 完整的 API 文档
- [使用示例](/examples/) - 25 个实际使用场景
- [架构设计](/architecture/) - 系统架构和设计理念

