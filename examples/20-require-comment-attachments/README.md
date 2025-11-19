# 场景 20: 审批意见和附件必填场景

## 场景描述

这是一个审批意见和附件必填场景示例,演示重要合同审批流程,要求审批人必须填写审批意见并上传相关附件.

**业务背景**: 重要合同审批需要完整的审批记录,确保审批过程有据可查.系统应该支持审批意见和附件必填配置,审批操作时自动验证,不符合要求拒绝操作.这种模式适用于重要审批流程,需要完整审批记录的场景.

## 流程结构

```
开始节点 → 审批节点(配置必填) → 结束节点
```

流程说明:
1. **开始节点**: 流程的起点,自动执行
2. **审批节点**: 合同审批节点,配置了审批意见和附件必填
   - 审批模式: 单人审批 (ApprovalModeSingle)
   - 审批人: 法务经理
   - **审批意见必填**: RequireComment=true
   - **附件必填**: RequireAttachments=true
3. **结束节点**: 流程的终点

## 关键特性

- **审批意见必填**: 使用 RequireComment 配置,强制要求填写审批意见
- **附件要求必填**: 使用 RequireAttachments 配置,强制要求上传附件
- **验证机制**: 审批操作时自动验证,不符合要求拒绝操作
- **提高审批质量**: 确保审批过程有据可查

## 代码说明

### 主要组件

1. **模板创建** (`createTemplate`):
   - 创建包含合同审批节点的模板
   - 配置审批模式为单人审批 (ApprovalModeSingle)
   - 配置审批意见必填 (RequireCommentField: true)
   - 配置附件必填 (RequireAttachmentsField: true)

2. **任务创建** (`main`):
   - 创建重要合同审批任务
   - 传入任务参数(合同编号、合同类型、承包商、金额、描述等)

3. **流程执行** (`runScenario`):
   - 提交任务进入审批流程
   - 设置审批人
   - 测试审批意见必填验证
   - 测试附件必填验证
   - 正确审批(带审批意见和附件)

4. **结果输出** (`printResults`):
   - 输出任务状态、审批记录
   - 验证审批意见和附件必填

### 代码结构

```go
main()
  ├── 创建管理器 (TemplateManager, TaskManager)
  ├── 创建模板 (createTemplate)
  │   └── 配置审批节点 (RequireComment, RequireAttachments)
  ├── 创建任务 (Create)
  ├── 执行流程 (runScenario)
  │   ├── 提交任务 (Submit)
  │   ├── 设置审批人 (AddApprover)
  │   ├── 测试审批意见必填验证
  │   ├── 测试附件必填验证
  │   └── 正确审批 (ApproveWithAttachments)
  └── 输出结果 (printResults)
```

### 审批意见和附件必填配置示例

```go
Config: &node.ApprovalNodeConfig{
    Mode: node.ApprovalModeSingle,
    ApproverConfig: &node.FixedApproverConfig{
        Approvers: []string{"legal-manager-001"},
    },
    RequireCommentField:    true, // 审批意见必填
    RequireAttachmentsField: true, // 附件必填
}
```

### 审批操作示例

**带附件的审批操作**:
```go
attachments := []string{
    "contract-review-report.pdf",
    "legal-opinion.pdf",
}
err = taskMgr.ApproveWithAttachments(tsk.ID, nodeID, approver, comment, attachments)
```

## 执行方法

### 前置要求

- Go 1.25.4 或更高版本
- 已安装 approval-kit 依赖

### 执行步骤

1. 进入场景目录:
   ```bash
   cd examples/20-require-comment-attachments
   ```

2. 运行程序:
   ```bash
   go run main.go
   ```

### 预期输出

程序运行后会输出以下内容:

```
=== 场景 20: 审批意见和附件必填场景 ===

步骤 1: 创建审批模板
✓ 模板创建成功: ID=require-comment-attachments-template, Name=重要合同审批模板

步骤 2: 创建重要合同审批任务
✓ 任务创建成功: ID=task-1763522619259830000-1, State=pending

步骤 3: 提交任务进入审批流程
✓ 任务已提交: State=submitted

步骤 4: 设置审批人
✓ 审批人已设置: legal-manager-001

步骤 5: 测试审批意见必填验证
  说明: 节点配置了 RequireComment=true,审批时必须填写审批意见
  测试: 尝试不填写审批意见进行审批

  ✓ 验证通过: 审批意见为空时被拒绝 (错误: comment is required for approval node "contract-approval")

步骤 6: 测试附件必填验证
  说明: 节点配置了 RequireAttachments=true,审批时必须上传附件
  测试: 尝试不上传附件进行审批

  ✓ 验证通过: 附件为空时被拒绝 (错误: attachments are required for approval node "contract-approval")

步骤 7: 正确审批(带审批意见和附件)
  说明: 填写审批意见并上传附件后,审批操作成功

✓ 审批已通过: ID=task-1763522619259830000-1, State=approving
  审批意见: 合同条款审查通过,已上传审查报告和法律意见书
  附件数量: 2
    附件 1: contract-review-report.pdf
    附件 2: legal-opinion.pdf

=== 审批结果 ===
任务 ID: task-1763522619259830000-1
业务 ID: contract-001
模板 ID: require-comment-attachments-template
当前状态: approving
当前节点: start
创建时间: 2025-11-19 11:23:39
提交时间: 2025-11-19 11:23:39
更新时间: 2025-11-19 11:23:39

任务参数:
  contractNo: CT-2025-001
  contractType: 重要合同
  contractor: 供应商-001
  amount: 5e+06
  description: 重要合同审批,要求审批意见和附件必填

审批记录(按时间顺序):
  记录 1:
    节点: 合同审批
    审批人: legal-manager-001
    操作类型: 加签
    审批意见: 设置法务经理为审批人
    操作时间: 2025-11-19 11:23:39
  记录 2:
    节点: 合同审批
    审批人: legal-manager-001
    操作类型: 审批通过
    审批意见: 合同条款审查通过,已上传审查报告和法律意见书
    附件数量: 2
      附件 1: contract-review-report.pdf
      附件 2: legal-opinion.pdf
    操作时间: 2025-11-19 11:23:39

状态变更历史:
  变更 1: pending -> submitted (原因: task submitted, 时间: 2025-11-19 11:23:39)

=== 验证结果 ===
✗ 任务状态异常: 期望 approved, 实际 approving
✓ 审批意见和附件必填验证:
  - 审批意见已填写
  - 附件已上传
```

**注意**: 
- 任务 ID 是动态生成的,每次运行都会不同
- 时间戳是实际运行时间,每次运行都会不同
- 任务参数中的字段顺序可能因 JSON 解析而有所不同
- 审批意见和附件必填验证正常工作,审批操作时自动验证

## 验证步骤

1. **验证模板创建**:
   - 检查模板是否包含开始节点、审批节点、结束节点
   - 检查审批节点是否配置了 RequireComment=true
   - 检查审批节点是否配置了 RequireAttachments=true

2. **验证任务创建**:
   - 检查任务初始状态是否为 `pending`
   - 检查任务参数是否正确设置(合同信息等)

3. **验证任务提交**:
   - 检查任务状态是否从 `pending` 变为 `submitted`
   - 检查提交时间是否已设置

4. **验证审批意见必填**:
   - 检查审批意见为空时是否被拒绝
   - 检查错误信息是否正确

5. **验证附件必填**:
   - 检查附件为空时是否被拒绝
   - 检查错误信息是否正确

6. **验证正确审批**:
   - 检查填写审批意见并上传附件后,审批操作是否成功
   - 检查审批记录是否包含审批意见和附件

### 验证要点

- ✓ 审批意见必填配置正确 (RequireComment=true)
- ✓ 附件必填配置正确 (RequireAttachments=true)
- ✓ 审批意见为空时被拒绝
- ✓ 附件为空时被拒绝
- ✓ 填写审批意见并上传附件后,审批操作成功
- ✓ 审批记录包含审批意见和附件

## 审批意见和附件必填说明

### 审批意见必填

- **配置**: `RequireCommentField: true`
- **验证时机**: 审批操作时自动验证
- **验证逻辑**: 如果配置了必填但审批意见为空,拒绝操作并返回错误
- **错误信息**: `comment is required for approval node "node-id"`

### 附件必填

- **配置**: `RequireAttachmentsField: true`
- **验证时机**: 审批操作时自动验证(使用 ApproveWithAttachments 方法)
- **验证逻辑**: 如果配置了必填但附件列表为空,拒绝操作并返回错误
- **错误信息**: `attachments are required for approval node "node-id"`

### 审批操作

- **Approve**: 不带附件的审批操作,如果配置了 RequireAttachments,需要使用 ApproveWithAttachments
- **ApproveWithAttachments**: 带附件的审批操作,支持上传附件列表

## 适用场景

这个场景适用于以下实际业务场景:

- **重要合同审批**: 需要完整审批记录的重要合同审批
- **重大决策审批**: 需要完整审批记录的重大决策审批
- **其他重要审批**: 需要完整审批记录的其他重要审批流程

## 扩展说明

如果需要扩展此场景,可以考虑:

1. **测试拒绝操作**: 测试拒绝操作时是否也需要审批意见和附件
2. **测试部分必填**: 测试只配置审批意见必填或只配置附件必填
3. **测试多个附件**: 测试上传多个附件的情况
4. **测试附件类型**: 测试附件类型验证(如果业务系统支持)
5. **测试附件大小**: 测试附件大小限制(如果业务系统支持)

## 相关场景

- **场景 01**: 最简单的单人审批流程 - 可以结合审批意见和附件必填实现完整的审批流程
- **场景 08**: 高级操作场景 - 可以结合审批意见和附件必填实现更完整的审批记录

