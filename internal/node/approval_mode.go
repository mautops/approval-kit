package node

// ApprovalMode 表示审批节点的审批模式
// 不同的审批模式决定了审批完成的判断条件
type ApprovalMode string

const (
	// ApprovalModeSingle 单人审批: 单个审批人进行审批
	// 审批人同意/拒绝后流程继续
	ApprovalModeSingle ApprovalMode = "single"

	// ApprovalModeUnanimous 多人会签: 多个审批人需全部同意
	// 所有审批人同意后流程继续;任一审批人拒绝时,根据配置决定流程走向
	ApprovalModeUnanimous ApprovalMode = "unanimous"

	// ApprovalModeOr 多人或签: 多个审批人中任意一人同意即可
	// 第一个审批人同意后流程继续;所有审批人拒绝后,根据配置决定流程走向
	ApprovalModeOr ApprovalMode = "or"

	// ApprovalModeProportional 比例会签: 多个审批人中达到指定比例同意即可
	// 例如 5 人中 3 人同意即可通过
	ApprovalModeProportional ApprovalMode = "proportional"

	// ApprovalModeSequential 顺序审批: 多个审批人按顺序依次审批
	// 前一个审批人同意后下一个审批人才能审批
	ApprovalModeSequential ApprovalMode = "sequential"
)

