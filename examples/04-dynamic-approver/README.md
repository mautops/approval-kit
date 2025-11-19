# 场景 04: 动态审批人获取

## 场景描述

这是一个动态审批人获取示例,模拟项目立项审批场景.审批人需要根据项目类型和部门动态获取,而不是在模板中固定配置.

**业务背景**: 员工提交项目立项申请,系统需要根据项目类型(研发项目、市场项目等)和部门(技术部、市场部等)动态确定审批人.不同类型的项目需要不同层级的审批人,不同部门也有不同的审批流程.

## 流程结构

```
开始节点 → 审批节点(动态审批人) → 结束节点
```

流程说明:
1. **开始节点**: 流程的起点,自动执行
2. **审批节点**: 使用动态审批人配置
   - 审批模式: 单人审批 (ApprovalModeSingle)
   - 审批人配置: 动态审批人 (DynamicApproverConfig)
   - 获取时机: 节点激活时 (ApproverTimingOnActivate)
   - API 调用: 通过 HTTP API 获取审批人列表
3. **结束节点**: 流程的终点

## 关键特性

- **动态审批人配置**: 使用 DynamicApproverConfig 通过 HTTP API 获取审批人
- **HTTP API 配置**: 配置 API 地址、请求方法、请求头等
- **参数映射**: 从任务参数中提取数据并映射到 API 请求参数
- **响应解析**: 从 API 响应中解析审批人列表
- **获取时机**: 支持节点激活时获取 (ApproverTimingOnActivate)
- **重试机制**: API 调用失败时自动重试(指数退避)
- **Mock HTTP 客户端**: 使用 Mock 客户端模拟 API 调用

## 代码说明

### 主要组件

1. **Mock HTTP 客户端** (`createMockHTTPClient`):
   - 实现 `node.HTTPClient` 接口
   - 模拟 API 响应,返回审批人列表
   - 用于演示,不进行真实的网络请求

2. **模板创建** (`createTemplate`):
   - 创建包含动态审批人配置的模板
   - 配置 HTTP API 地址、请求方法、请求头
   - 配置参数映射规则(从任务参数中获取项目类型)
   - 配置响应解析规则(从响应中解析审批人列表)
   - 设置获取时机为节点激活时

3. **任务创建** (`main`):
   - 创建项目立项任务
   - 传入任务参数(项目类型、部门、预算、描述)

4. **流程执行** (`runScenario`):
   - 提交任务进入审批流程
   - 在节点激活时调用 API 获取审批人
   - 设置审批人列表
   - 执行审批操作

5. **结果输出** (`printResults`):
   - 输出任务状态、审批人列表、审批记录
   - 验证动态审批人获取是否成功

### 代码结构

```go
main()
  ├── 创建管理器 (TemplateManager, TaskManager)
  ├── 创建 Mock HTTP 客户端 (createMockHTTPClient)
  ├── 创建模板 (createTemplate)
  │   └── 配置动态审批人 (DynamicApproverConfig)
  │       ├── HTTP API 配置
  │       ├── 参数映射规则
  │       ├── 响应解析规则
  │       └── 获取时机配置
  ├── 创建任务 (Create)
  ├── 执行流程 (runScenario)
  │   ├── 提交任务 (Submit)
  │   ├── 动态获取审批人 (GetApprovers)
  │   ├── 设置审批人列表 (AddApprover)
  │   └── 执行审批 (Approve)
  └── 输出结果 (printResults)
```

### 动态审批人配置示例

```go
ApproverConfig: &node.DynamicApproverConfig{
    API: &node.HTTPAPIConfig{
        URL:    "http://api.example.com/approvers",
        Method: "POST",
        Headers: map[string]string{
            "Content-Type": "application/json",
            "Authorization": "Bearer token-123",
        },
        // 参数映射: 从任务参数中获取项目类型
        ParamMapping: &node.ParamMapping{
            Source: "task_params",
            Path:   "projectType",
            Target: "project_type",
        },
        // 响应解析: 从 API 响应中解析审批人列表
        ResponseMapping: &node.ResponseMapping{
            Path:   "data.approvers",
            Format: "json",
        },
    },
    Timing:     node.ApproverTimingOnActivate,
    HTTPClient: httpClient,
}
```

## 执行方法

### 前置要求

- Go 1.25.4 或更高版本
- 已安装 approval-kit 依赖

### 执行步骤

1. 进入场景目录:
   ```bash
   cd examples/04-dynamic-approver
   ```

2. 运行程序:
   ```bash
   go run main.go
   ```

### 预期输出

程序运行后会输出以下内容:

```
=== 场景 04: 动态审批人获取 ===

步骤 1: 创建 Mock HTTP 客户端
✓ Mock HTTP 客户端已创建

步骤 2: 创建审批模板(配置动态审批人)
✓ 模板创建成功: ID=project-approval-template, Name=项目立项审批模板

步骤 3: 创建项目立项任务
✓ 任务创建成功: ID=task-1763519480908463000-1, State=pending

步骤 4: 提交任务进入审批流程
✓ 任务已提交: State=submitted

步骤 5: 动态获取审批人(节点激活时)
  说明: 系统会调用 HTTP API,根据项目类型和部门获取审批人列表

  调用 API: POST http://api.example.com/approvers
  请求参数: {"project_type": "研发项目"}
  ✓ API 调用成功,获取到审批人: [tech-lead-001 tech-manager-001]

✓ 审批人列表已设置: tech-lead-001 (单人审批模式,使用第一个审批人)
  注意: API 返回了 2 个审批人,但当前节点配置为单人审批模式

步骤 6: 执行审批操作
  说明: 使用第一个审批人进行审批(单人审批模式)

  使用审批人: tech-lead-001
  ✓ tech-lead-001 已同意 (当前状态: approved)
  注意: API 返回了 2 个审批人,但当前节点配置为单人审批模式,只使用第一个审批人

✓ 审批流程完成

=== 审批结果 ===
任务 ID: task-1763519480908463000-1
业务 ID: project-001
模板 ID: project-approval-template
当前状态: approved
当前节点: start
创建时间: 2025-11-19 10:31:20
提交时间: 2025-11-19 10:31:20
更新时间: 2025-11-19 10:31:20

任务参数:
  projectType: 研发项目
  department: 技术部
  budget: 100000
  description: 新产品研发项目立项

审批人列表:
  1. tech-lead-001 (动态获取)

审批记录:
  记录 1 (审批):
    节点 ID: project-approval
    审批人: tech-lead-001
    审批结果: approve
    审批意见: 项目立项审批通过
    审批时间: 2025-11-19 10:31:20

状态变更历史:
  变更 1: pending -> submitted (原因: task submitted, 时间: 2025-11-19 10:31:20)
  变更 2: approving -> approved (原因: all approvers approved, 时间: 2025-11-19 10:31:20)

=== 验证结果 ===
✓ 任务已成功通过审批
✓ 已通过 API 动态获取 1 个审批人
```

**注意**: 
- 任务 ID 是动态生成的,每次运行都会不同
- 时间戳是实际运行时间,每次运行都会不同
- 任务参数中的字段顺序可能因 JSON 解析而有所不同
- 本示例使用 Mock HTTP 客户端模拟 API 调用,实际使用时需要提供真实的 HTTP 客户端实现
- API 返回了 2 个审批人,但当前节点配置为单人审批模式,只使用第一个审批人

## 验证步骤

1. **验证模板创建**:
   - 检查模板是否包含开始节点、审批节点、结束节点
   - 检查审批节点配置是否正确(动态审批人配置)
   - 检查 HTTP API 配置是否正确(URL、方法、请求头)
   - 检查参数映射和响应解析配置是否正确

2. **验证任务创建**:
   - 检查任务初始状态是否为 `pending`
   - 检查任务参数是否正确设置(项目类型、部门等)

3. **验证任务提交**:
   - 检查任务状态是否从 `pending` 变为 `submitted`
   - 检查提交时间是否已设置

4. **验证动态审批人获取**:
   - 检查 API 调用是否成功
   - 检查是否从 API 响应中正确解析出审批人列表
   - 检查参数映射是否正确(请求参数是否包含项目类型)

5. **验证审批人设置**:
   - 检查审批人列表是否正确设置
   - 检查审批人是否来自 API 响应

6. **验证审批操作**:
   - 检查审批操作是否成功执行
   - 检查任务状态是否正确更新

7. **验证审批记录**:
   - 检查是否生成了审批记录
   - 检查审批记录内容是否正确

### 验证要点

- ✓ 任务最终状态为 `approved`
- ✓ 成功通过 API 获取到审批人列表
- ✓ 参数映射正确(请求参数包含项目类型)
- ✓ 响应解析正确(从响应中解析出审批人)
- ✓ 审批人列表已正确设置
- ✓ 审批记录已生成
- ✓ 状态变更历史包含状态转换记录

## 适用场景

这个场景适用于以下实际业务场景:

- **项目立项审批**: 根据项目类型和部门动态确定审批人
- **预算申请**: 根据申请金额和部门动态确定审批人
- **合同审批**: 根据合同类型和金额动态确定审批人
- **采购审批**: 根据采购类型和部门动态确定审批人

## 扩展说明

如果需要扩展此场景,可以考虑:

1. **测试不同项目类型**: 修改任务参数,测试不同项目类型返回不同的审批人
2. **任务创建时获取**: 配置获取时机为 `ApproverTimingOnCreate`,在任务创建时获取审批人
3. **多个参数映射**: 配置多个参数映射规则,从任务参数中提取多个字段
4. **复杂响应解析**: 配置更复杂的响应路径,解析嵌套的 JSON 结构
5. **API 认证**: 配置更复杂的请求头,支持 Token 认证等
6. **重试机制演示**: 模拟 API 调用失败,展示重试机制

## 相关场景

- **场景 01**: 最简单的单人审批流程 - 固定审批人场景
- **场景 21**: 任务创建时获取动态审批人场景 - 展示不同的获取时机
- **场景 02**: 多人会签审批流程 - 可以结合动态审批人实现多人会签

