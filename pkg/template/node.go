package template

import (
	internalTemplate "github.com/mautops/approval-kit/internal/template"
)

// Node 表示审批节点
// 节点是审批流程的基本单元,不同类型的节点有不同的执行逻辑
// 与 internal/template.Node 结构相同,但位于 pkg 目录,可以被外部导入
type Node = internalTemplate.Node

// NodeConfig 节点配置接口
// 不同类型的节点有不同的配置实现
// 与 internal/template.NodeConfig 接口定义完全一致,但位于 pkg 目录,可以被外部导入
type NodeConfig = internalTemplate.NodeConfig

// ApprovalNodeConfigAccessor 审批节点配置访问接口
// 用于在不导入 node 包的情况下访问审批节点配置的属性
// 避免循环依赖
// 与 internal/template.ApprovalNodeConfigAccessor 接口定义完全一致,但位于 pkg 目录,可以被外部导入
type ApprovalNodeConfigAccessor = internalTemplate.ApprovalNodeConfigAccessor

// OperationPermissionsAccessor 操作权限访问接口
// 用于在不导入 node 包的情况下访问操作权限配置
// 与 internal/template.OperationPermissionsAccessor 接口定义完全一致,但位于 pkg 目录,可以被外部导入
type OperationPermissionsAccessor = internalTemplate.OperationPermissionsAccessor

// NodeFromInternal 将 internal.Node 转换为 pkg.Node
func NodeFromInternal(n *internalTemplate.Node) *Node {
	return (*Node)(n)
}

// NodeToInternal 将 pkg.Node 转换为 internal.Node
func NodeToInternal(n *Node) *internalTemplate.Node {
	return (*internalTemplate.Node)(n)
}

