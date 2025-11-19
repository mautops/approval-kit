---
title: "状态机"
description: "状态机的设计和实现"
date: 2025-11-19
draft: false
weight: 40
---

## 设计理念

Approval Kit 使用有限状态机(FSM)管理任务状态流转,确保:

- 状态转换的合法性
- 状态转换的一致性
- 状态转换的可追溯性

## 状态定义

任务状态包括:

- `pending`: 待审批
- `submitted`: 已提交
- `approving`: 审批中
- `approved`: 已通过
- `rejected`: 已拒绝
- `cancelled`: 已取消
- `timeout`: 已超时

## 转换规则

状态转换必须遵循预定义的规则:

```
pending → submitted → approving → approved/rejected
         ↓
      cancelled
```

## 转换执行

状态转换通过 `StateMachine.Transition` 方法执行:

1. 验证转换合法性
2. 创建新任务对象(不可变实现)
3. 更新状态和时间戳
4. 记录状态变更历史

## 并发安全

状态转换使用版本号机制确保原子性,防止并发修改冲突.

