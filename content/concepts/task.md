---
title: "审批任务"
description: "审批任务的结构和生命周期"
date: 2025-11-19
draft: false
weight: 30
---

## 任务结构

审批任务是基于模板创建的实例,包含:

- **基本信息**: ID、模板 ID、业务 ID、参数
- **状态信息**: 当前状态、当前节点
- **运行时数据**: 节点输出、审批人列表、审批结果
- **审批记录**: 所有审批操作的记录
- **状态历史**: 状态变更的完整历史

## 任务状态

- **pending**: 待审批
- **submitted**: 已提交
- **approving**: 审批中
- **approved**: 已通过
- **rejected**: 已拒绝
- **cancelled**: 已取消
- **timeout**: 已超时

## 状态流转

任务状态通过状态机管理,确保状态转换的合法性和一致性.

状态转换规则:
- `pending` → `submitted` → `approving` → `approved`/`rejected`
- 任何状态都可以转换为 `cancelled`
- 终态(`approved`、`rejected`、`cancelled`、`timeout`)不能转换到其他状态

## 审批记录

每次审批操作都会生成审批记录,包含:

- 节点 ID
- 审批人
- 审批结果(approve/reject/transfer)
- 审批意见
- 审批时间
- 附件信息

