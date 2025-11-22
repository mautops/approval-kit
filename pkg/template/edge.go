package template

import (
	internalTemplate "github.com/mautops/approval-kit/internal/template"
)

// Edge 表示节点间的连接
// 定义审批流程中节点之间的流转关系,支持条件分支
// 与 internal/template.Edge 结构相同,但位于 pkg 目录,可以被外部导入
type Edge = internalTemplate.Edge

// EdgeFromInternal 将 internal.Edge 转换为 pkg.Edge
func EdgeFromInternal(e *internalTemplate.Edge) *Edge {
	return (*Edge)(e)
}

// EdgeToInternal 将 pkg.Edge 转换为 internal.Edge
func EdgeToInternal(e *Edge) *internalTemplate.Edge {
	return (*internalTemplate.Edge)(e)
}

