# 开发指南

## 当前状态

✅ **Phase 1.1: 项目初始化** - 已完成

项目基础结构已创建:
- 目录结构已建立
- 基础包文件已创建
- 项目结构验证测试已通过
- 项目可以正常编译

## TDD 开发流程

每个功能开发必须遵循以下流程:

### 1. Red 阶段 (编写失败的测试)

```bash
# 在 tests/ 目录下创建测试文件
# 编写测试用例,运行测试确认失败
go test ./tests/... -v
```

### 2. Green 阶段 (编写最小实现)

```bash
# 在 internal/ 目录下实现功能
# 编写最少的代码使测试通过
go test ./tests/... -v
```

### 3. Refactor 阶段 (重构优化)

```bash
# 优化代码结构,确保测试仍然通过
go test ./tests/... -v
go vet ./...
```

## 下一步开发任务

根据 `tasks.mdc` 中的任务清单,下一步应该:

### Phase 1.2: Go 模块配置验证

1. **Red**: 编写模块依赖验证测试
2. **Green**: 验证 go.mod 配置正确
3. **Refactor**: 优化模块配置

### Phase 2: 核心数据结构定义

按照以下顺序实现:

1. **任务状态定义** (TaskState)
2. **节点类型定义** (NodeType)
3. **审批模式定义** (ApprovalMode)
4. **错误定义** (Errors)

## 开发命令

```bash
# 运行所有测试
go test ./tests/... -v

# 运行特定测试
go test ./tests/... -run TestName -v

# 查看测试覆盖率
go test ./tests/... -cover

# 生成覆盖率报告
go test ./tests/... -coverprofile=coverage.out
go tool cover -html=coverage.out

# 代码检查
go vet ./...

# 编译检查
go build ./...
```

## 项目结构说明

```
approval-kit/
├── internal/          # 内部实现代码
│   ├── task/         # 任务管理模块
│   ├── template/     # 模板管理模块
│   ├── statemachine/ # 状态机引擎
│   ├── node/         # 节点执行引擎
│   └── event/        # 事件通知模块
├── tests/            # 测试代码(独立目录)
│   ├── task/         # 任务管理测试
│   ├── template/     # 模板管理测试
│   ├── statemachine/ # 状态机测试
│   ├── node/         # 节点执行测试
│   └── event/        # 事件通知测试
├── examples/         # 示例代码
├── go.mod           # Go 模块定义
├── README.md        # 项目说明
└── DEVELOPMENT.md   # 开发指南(本文件)
```

## 开发规范

1. **严格遵循 TDD**: 先写测试,再写实现
2. **小步快跑**: 每个任务足够小,可以在一个 TDD 循环中完成
3. **测试覆盖率**: 关键路径 100%,整体 80%+
4. **代码规范**: 通过 go vet 和 golangci-lint 检查
5. **文档完整**: 所有公共 API 必须有文档注释

## 参考文档

- `tasks.mdc`: 详细的任务清单
- `spec.mdc`: 技术实现方案
- `prd.mdc`: 功能需求文档
- `constitution.mdc`: 项目开发原则

