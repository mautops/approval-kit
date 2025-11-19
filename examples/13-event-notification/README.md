# 场景 13: 事件通知集成场景

## 场景描述

这是一个事件通知集成场景示例,演示审批流程中的关键事件通过 Webhook 推送到业务系统.

**业务背景**: 在实际业务中,审批流程中的关键事件需要通知到外部系统,如通知系统、日志系统、数据分析系统等.系统应该支持异步事件推送,不阻塞主流程,并提供重试机制和幂等性保证.

## 流程结构

```
开始节点 → 审批节点 → 结束节点
```

流程说明:
1. **开始节点**: 流程的起点,自动执行
2. **审批节点**: 审批节点
   - 审批模式: 单人审批 (ApprovalModeSingle)
   - 审批人: 审批人
3. **结束节点**: 流程的终点

## 关键特性

- **多种事件类型**: 支持任务创建、提交、节点激活、审批操作、任务完成等事件
- **Webhook 配置**: 支持配置 Webhook 地址、HTTP 方法、请求头、超时时间等
- **异步推送**: 使用 channel 和 goroutine 实现异步事件推送,不阻塞主流程
- **重试机制**: 推送失败自动重试,使用指数退避策略
- **幂等性保证**: 确保事件不重复处理

## 代码说明

### 主要组件

1. **Mock Webhook 服务器**:
   - 使用 httptest.NewServer 创建测试服务器
   - 接收并记录所有推送的事件

2. **Webhook 处理器** (`webhookHandler`):
   - 配置 Webhook 地址、HTTP 方法、请求头、超时时间
   - 将事件序列化为 JSON 并推送到 Webhook URL

3. **事件通知器** (`notifier`):
   - 使用 channel 和 goroutine 实现异步事件推送
   - 支持多个事件处理器
   - 提供重试机制和幂等性保证

4. **任务管理器** (`taskMgr`):
   - 使用 NewTaskManagerWithNotifier 创建带事件通知器的任务管理器
   - 在关键操作点自动生成并推送事件

5. **流程执行** (`runScenario`):
   - 创建任务 → 触发 task_created 事件
   - 提交任务 → 触发 task_submitted 和 node_activated 事件
   - 设置审批人 → 触发 approval_operation 事件
   - 审批操作 → 触发 approval_operation、node_completed、task_approved 事件

### 代码结构

```go
main()
  ├── 创建 Mock Webhook 服务器
  ├── 创建 Webhook 处理器
  ├── 创建事件通知器
  ├── 创建管理器 (带事件通知器)
  ├── 创建模板
  ├── 执行流程 (runScenario)
  │   ├── 创建任务 (触发 task_created 事件)
  │   ├── 提交任务 (触发 task_submitted、node_activated 事件)
  │   ├── 设置审批人 (触发 approval_operation 事件)
  │   └── 审批操作 (触发 approval_operation、node_completed、task_approved 事件)
  ├── 等待事件推送完成
  ├── 输出接收到的事件
  └── 停止事件通知器
```

### Webhook 配置示例

```go
webhookConfig := &event.WebhookConfig{
    URL:    "https://example.com/webhook",
    Method: "POST",
    Headers: map[string]string{
        "Content-Type": "application/json",
        "X-Auth-Token":  "your-token",
    },
    Timeout: 30,
}
webhookHandler := event.NewWebhookHandler(webhookConfig)
```

### 事件通知器配置示例

```go
notifier := event.NewEventNotifier([]event.EventHandler{webhookHandler}, 100)
taskMgr := task.NewTaskManagerWithNotifier(templateMgr, nil, notifier)
```

## 执行方法

### 前置要求

- Go 1.25.4 或更高版本
- 已安装 approval-kit 依赖

### 执行步骤

1. 进入场景目录:
   ```bash
   cd examples/13-event-notification
   ```

2. 运行程序:
   ```bash
   go run main.go
   ```

### 预期输出

程序运行后会输出以下内容:

```
=== 场景 13: 事件通知集成场景 ===

步骤 1: 创建 Mock Webhook 服务器
✓ Webhook 服务器已启动: http://127.0.0.1:61224

步骤 2: 创建 Webhook 处理器
✓ Webhook 处理器已创建

步骤 3: 创建事件通知器
✓ 事件通知器已创建

步骤 4: 创建管理器
✓ 管理器已创建

步骤 5: 创建审批模板
✓ 模板创建成功: ID=event-notification-template, Name=事件通知模板

=== 执行审批流程 ===

步骤 1: 创建任务
✓ 任务创建成功: ID=task-1763521652523572000-1
  说明: 任务创建时会触发 task_created 事件

  [Webhook] 收到事件: Type=task_created, TaskID=task-1763521652523572000-1, Time=2025-11-19 11:07:32
步骤 2: 提交任务
✓ 任务已提交: ID=task-1763521652523572000-1
  说明: 任务提交时会触发 task_submitted 事件

  [Webhook] 收到事件: Type=task_submitted, TaskID=task-1763521652523572000-1, Time=2025-11-19 11:07:32
  [Webhook] 收到事件: Type=node_activated, TaskID=task-1763521652523572000-1, Time=2025-11-19 11:07:32
步骤 3: 设置审批人
✓ 审批人已设置: manager-001

  [Webhook] 收到事件: Type=approval_operation, TaskID=task-1763521652523572000-1, Time=2025-11-19 11:07:32
步骤 4: 审批操作
✓ 审批已通过: ID=task-1763521652523572000-1, State=approved
  说明: 审批操作时会触发以下事件:
    - approval_operation: 审批操作事件
    - node_completed: 节点完成事件
    - task_approved: 任务通过事件

步骤 6: 等待事件推送完成
  [Webhook] 收到事件: Type=approval_operation, TaskID=task-1763521652523572000-1, Time=2025-11-19 11:07:32
  [Webhook] 收到事件: Type=task_approved, TaskID=task-1763521652523572000-1, Time=2025-11-19 11:07:32
  [Webhook] 收到事件: Type=node_completed, TaskID=task-1763521652523572000-1, Time=2025-11-19 11:07:32

=== 接收到的 Webhook 事件 ===
  共接收到 7 个事件:

  事件 1:
    事件类型: task_created
    事件时间: 2025-11-19 11:07:32
    任务 ID: task-1763521652523572000-1
    业务 ID: event-001
    任务状态: pending
    节点信息: 开始 (start)

  事件 2:
    事件类型: task_submitted
    事件时间: 2025-11-19 11:07:32
    任务 ID: task-1763521652523572000-1
    业务 ID: event-001
    任务状态: submitted
    节点信息: 开始 (start)

  事件 3:
    事件类型: node_activated
    事件时间: 2025-11-19 11:07:32
    任务 ID: task-1763521652523572000-1
    业务 ID: event-001
    任务状态: submitted
    节点信息: 审批 (approval)

  事件 4:
    事件类型: approval_operation
    事件时间: 2025-11-19 11:07:32
    任务 ID: task-1763521652523572000-1
    业务 ID: event-001
    任务状态: submitted
    节点信息: 审批 (approval)
    审批信息: manager-001 - add_approver

  事件 5:
    事件类型: approval_operation
    事件时间: 2025-11-19 11:07:32
    任务 ID: task-1763521652523572000-1
    业务 ID: event-001
    任务状态: approved
    节点信息: 审批 (approval)
    审批信息: manager-001 - approve

  事件 6:
    事件类型: task_approved
    事件时间: 2025-11-19 11:07:32
    任务 ID: task-1763521652523572000-1
    业务 ID: event-001
    任务状态: approved
    节点信息: 审批 (approval)

  事件 7:
    事件类型: node_completed
    事件时间: 2025-11-19 11:07:32
    任务 ID: task-1763521652523572000-1
    业务 ID: event-001
    任务状态: approved
    节点信息: 审批 (approval)

✓ 事件通知器已停止
```

**注意**: 
- Webhook 服务器地址是动态生成的,每次运行都会不同
- 时间戳是实际运行时间,每次运行都会不同
- 事件推送是异步的,需要等待一段时间才能看到所有事件

## 验证步骤

1. **验证 Webhook 服务器**:
   - 检查 Webhook 服务器是否成功启动
   - 检查 Webhook 服务器是否能接收事件

2. **验证事件推送**:
   - 检查任务创建时是否触发 task_created 事件
   - 检查任务提交时是否触发 task_submitted 和 node_activated 事件
   - 检查设置审批人时是否触发 approval_operation 事件
   - 检查审批操作时是否触发 approval_operation、node_completed、task_approved 事件

3. **验证事件内容**:
   - 检查事件是否包含完整的信息(事件类型、时间、任务信息、节点信息等)
   - 检查审批操作事件是否包含审批信息

### 验证要点

- ✓ Webhook 服务器成功接收所有事件
- ✓ 事件类型正确(7 种事件类型)
- ✓ 事件内容完整(任务信息、节点信息、审批信息等)
- ✓ 事件推送是异步的,不阻塞主流程
- ✓ 事件通知器正确停止

## 支持的事件类型

系统支持以下事件类型:

- **task_created**: 任务创建事件
- **task_submitted**: 任务提交事件
- **node_activated**: 节点激活事件
- **approval_operation**: 审批操作事件(包括审批、拒绝、转交、加签、减签等)
- **task_approved**: 任务通过事件
- **task_rejected**: 任务拒绝事件
- **task_timeout**: 任务超时事件
- **task_cancelled**: 任务取消事件
- **task_withdrawn**: 任务撤回事件
- **node_completed**: 节点完成事件

## 事件数据结构

每个事件包含以下信息:

- **ID**: 事件 ID(用于幂等性保证)
- **Type**: 事件类型
- **Time**: 事件时间
- **Task**: 任务信息(任务 ID、模板 ID、业务 ID、任务状态)
- **Node**: 节点信息(节点 ID、节点名称、节点类型)
- **Approval**: 审批信息(如适用,包括节点 ID、审批人、审批结果、审批意见)
- **Business**: 业务信息(业务 ID)

## 适用场景

这个场景适用于以下实际业务场景:

- **通知系统集成**: 审批事件推送到通知系统,及时通知相关人员
- **日志系统集成**: 审批事件推送到日志系统,记录审批历史
- **数据分析系统集成**: 审批事件推送到数据分析系统,进行审批效率分析
- **其他外部系统集成**: 需要与外部系统集成的场景

## 扩展说明

如果需要扩展此场景,可以考虑:

1. **测试重试机制**: 模拟 Webhook 服务器失败,验证重试机制
2. **测试幂等性**: 验证相同事件不会重复处理
3. **测试多个 Webhook**: 配置多个 Webhook 地址,验证事件推送到所有地址
4. **测试认证**: 配置 Webhook 认证信息,验证认证是否正确
5. **测试超时**: 配置超时时间,验证超时处理是否正确

## 相关场景

- **场景 01**: 最简单的单人审批流程 - 可以结合事件通知实现完整的审批流程
- **场景 08**: 高级操作场景 - 高级操作也会触发相应的事件

