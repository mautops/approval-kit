package template

// NodeType 表示审批节点的类型
// 节点类型决定了节点的执行方式和行为
type NodeType string

const (
	// NodeTypeStart 开始节点: 标识审批流程的起点
	// 每个审批模板必须且只能有一个开始节点
	NodeTypeStart NodeType = "start"

	// NodeTypeApproval 审批节点: 执行审批操作的节点
	// 支持多种审批模式(单人审批、多人会签、多人或签等)
	NodeTypeApproval NodeType = "approval"

	// NodeTypeCondition 条件节点: 根据条件结果决定流程走向
	// 支持多种条件类型(数值比较、字符串匹配、枚举判断等)
	NodeTypeCondition NodeType = "condition"

	// NodeTypeEnd 结束节点: 标识审批流程的终点
	// 流程执行到此节点时,任务状态变为终态
	NodeTypeEnd NodeType = "end"
)

