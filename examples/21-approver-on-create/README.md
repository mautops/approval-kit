# 场景 21: 任务创建时获取动态审批人场景

## 场景描述

这是一个任务创建时获取动态审批人场景示例,演示项目立项审批流程,在任务创建时就需要确定所有审批人,便于提前通知.

**业务背景**: 项目立项审批需要在任务创建时就确定所有审批人,便于提前通知相关人员,让审批人提前了解项目信息.系统应该支持任务创建时获取动态审批人,配置为 `ApproverTimingOnCreate` 时机的动态审批人会在任务创建时立即获取.这种模式适用于需要提前确定审批人并通知的场景.

## 流程结构

```
开始节点 → 审批节点(任务创建时获取审批人) → 结束节点
```

流程说明:
1. **开始节点**: 流程的起点,自动执行
2. **审批节点**: 项目审批节点,配置为任务创建时获取动态审批人
   - 审批模式: 多人会签 (ApprovalModeUnanimous)
   - 审批人配置: 动态审批人 (DynamicApproverConfig)
   - **获取时机**: ApproverTimingOnCreate (任务创建时获取)
   - 根据项目类型动态获取审批人
3. **结束节点**: 流程的终点

## 关键特性

- **获取时机配置**: 使用 ApproverTimingOnCreate 实现任务创建时获取审批人
- **提前准备**: 在任务提交前就确定审批人列表
- **与节点激活时获取的区别**: 创建时获取适合需要提前通知的场景
- **支持参数传递**: 从任务参数中获取审批人信息

## 代码说明

### 主要组件

1. **Mock HTTP 客户端** (`createMockHTTPClient`):
   - 创建模拟 HTTP 客户端
   - 根据项目类型返回不同的审批人列表

2. **审批人获取函数** (`approverFetcherFunc`):
   - 在任务创建时调用
   - 使用 `node.FetchApproversOnCreate` 获取配置为 on_create 时机的动态审批人

3. **模板创建** (`createTemplate`):
   - 创建包含项目审批节点的模板
   - 配置审批模式为多人会签 (ApprovalModeUnanimous)
   - 配置动态审批人,获取时机为 on_create

4. **任务创建** (`main`):
   - 创建项目立项审批任务
   - 传入任务参数(项目编号、项目名称、项目经理、项目类型、预算、描述等)
   - 验证审批人是否已在任务创建时自动获取

5. **结果输出** (`printResults`):
   - 输出任务状态、审批人列表
   - 验证审批人获取功能

### 代码结构

```go
main()
  ├── 创建 Mock HTTP 客户端
  ├── 创建管理器 (带审批人获取函数)
  ├── 创建模板 (createTemplate)
  │   └── 配置审批节点 (ApproverTimingOnCreate)
  ├── 创建任务 (Create)
  │   └── 自动调用审批人获取函数
  ├── 验证审批人是否已获取
  └── 输出结果 (printResults)
```

### 任务创建时获取审批人配置示例

```go
ApproverConfig: &node.DynamicApproverConfig{
    API: &node.HTTPAPIConfig{
        URL:    "http://example.com/api/approvers",
        Method: "POST",
        ParamMapping: &node.ParamMapping{
            Source: "task_params",
            Path:   "projectType",
            Target: "projectType",
        },
        ResponseMapping: &node.ResponseMapping{
            Path:   "approvers",
            Format: "json",
        },
    },
    Timing: node.ApproverTimingOnCreate, // 任务创建时获取
}
```

### 审批人获取函数配置示例

```go
approverFetcherFunc := func(tpl *template.Template, tsk *task.Task) error {
    return node.FetchApproversOnCreate(tpl, tsk, httpClient)
}
taskMgr := task.NewTaskManager(templateMgr, approverFetcherFunc)
```

## 执行方法

### 前置要求

- Go 1.25.4 或更高版本
- 已安装 approval-kit 依赖

### 执行步骤

1. 进入场景目录:
   ```bash
   cd examples/21-approver-on-create
   ```

2. 运行程序:
   ```bash
   go run main.go
   ```

### 预期输出

程序运行后会输出以下内容:

```
=== 场景 21: 任务创建时获取动态审批人场景 ===

步骤 1: 创建 Mock HTTP 客户端
✓ Mock HTTP 客户端已创建

步骤 2: 创建管理器
✓ 管理器已创建(带审批人获取函数)

步骤 3: 创建审批模板
✓ 模板创建成功: ID=approver-on-create-template, Name=项目立项审批模板

步骤 4: 创建项目立项审批任务
✓ 任务创建成功: ID=task-1763522796640734000-1, State=pending

步骤 5: 验证审批人是否已获取
  说明: 任务创建时,配置为 on_create 时机的动态审批人应该已自动获取

✓ 审批人已自动获取: [tech-lead-001 tech-manager-001] (共 2 人)
  说明: 审批人在任务创建时已自动获取,无需等待节点激活

=== 审批结果 ===
任务 ID: task-1763522796640734000-1
业务 ID: project-001
模板 ID: approver-on-create-template
当前状态: pending
当前节点: start
创建时间: 2025-11-19 11:26:36
更新时间: 2025-11-19 11:26:36

任务参数:
  projectNo: PRJ-2025-001
  projectName: 新产品开发项目
  projectManager: pm-001
  projectType: 研发项目
  budget: 2e+06
  description: 新产品开发项目立项审批

审批人列表:
  节点: 项目审批
    - tech-lead-001
    - tech-manager-001
  说明: 审批人在任务创建时已自动获取(ApproverTimingOnCreate)

审批记录:
  无审批记录

状态变更历史:
  无状态变更记录

=== 验证结果 ===
✓ 审批人已在任务创建时自动获取:
  - 审批人数量: 2
  - 审批人列表: [tech-lead-001 tech-manager-001]
  - 获取时机: on_create (任务创建时)

✓ 任务创建时获取动态审批人功能配置和使用方法已展示
  说明: 在实际业务系统中,审批人获取函数会在任务创建时自动调用
  配置为 on_create 时机的动态审批人会在任务创建时立即获取
```

**注意**: 
- 任务 ID 是动态生成的,每次运行都会不同
- 时间戳是实际运行时间,每次运行都会不同
- 任务参数中的字段顺序可能因 JSON 解析而有所不同
- 审批人在任务创建时已自动获取,无需等待节点激活

## 验证步骤

1. **验证模板创建**:
   - 检查模板是否包含开始节点、审批节点、结束节点
   - 检查审批节点是否配置了动态审批人
   - 检查获取时机是否正确配置(ApproverTimingOnCreate)

2. **验证任务创建**:
   - 检查任务初始状态是否为 `pending`
   - 检查任务参数是否正确设置(项目信息等)

3. **验证审批人获取**:
   - 检查审批人是否在任务创建时已自动获取
   - 检查审批人列表是否正确(根据项目类型动态获取)
   - 检查审批人数量是否正确

### 验证要点

- ✓ 审批人获取时机配置正确 (ApproverTimingOnCreate)
- ✓ 审批人在任务创建时已自动获取
- ✓ 审批人列表正确(根据项目类型动态获取)
- ✓ 审批人数量正确

## 获取时机说明

系统支持以下获取时机:

- **ApproverTimingOnCreate (任务创建时)**: 任务创建时立即获取审批人
  - 适合需要提前确定审批人并通知的场景
  - 审批人在任务创建时已确定,无需等待节点激活
- **ApproverTimingOnActivate (节点激活时)**: 节点激活时获取审批人
  - 适合审批人可能变化的场景
  - 审批人在节点激活时才获取,可以获取最新的审批人信息

### 两种时机的区别

- **on_create**: 任务创建时获取,适合需要提前通知的场景
- **on_activate**: 节点激活时获取,适合审批人可能变化的场景

## 适用场景

这个场景适用于以下实际业务场景:

- **项目立项审批**: 需要提前确定审批人并通知
- **预算申请**: 需要提前确定审批人并通知
- **其他需要提前通知的场景**: 需要提前确定审批人并通知的审批流程

## 扩展说明

如果需要扩展此场景,可以考虑:

1. **测试多个节点**: 测试多个节点都配置为 on_create 时机的情况
2. **测试获取失败**: 测试审批人获取失败时的处理
3. **测试参数传递**: 测试从任务参数中传递不同参数获取不同审批人
4. **测试固定审批人**: 测试固定审批人配置为 on_create 时机的情况
5. **测试混合时机**: 测试部分节点配置为 on_create,部分配置为 on_activate

## 相关场景

- **场景 04**: 动态审批人获取场景 - 节点激活时获取动态审批人
- **场景 08**: 高级操作场景 - 可以结合任务创建时获取审批人实现更灵活的流程

