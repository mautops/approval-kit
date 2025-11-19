package template

// Edge 表示节点间的连接
// 定义审批流程中节点之间的流转关系,支持条件分支
type Edge struct {
	From      string // 源节点 ID
	To        string // 目标节点 ID
	Condition string // 条件表达式(可选,用于条件节点)
}

