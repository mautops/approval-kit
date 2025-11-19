package node

import (
	"encoding/json"

	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
)

// FetchApproversOnCreate 在任务创建时获取审批人
// 遍历模板中的所有节点,查找配置为 on_create 时机的动态审批人节点并获取审批人
// 这个函数应该在 TaskManager.Create 中调用,但由于循环依赖,需要通过依赖注入的方式调用
func FetchApproversOnCreate(tpl *template.Template, tsk *task.Task, httpClient HTTPClient) error {
	// 遍历模板中的所有节点
	for _, tplNode := range tpl.Nodes {
		if tplNode.Type != template.NodeTypeApproval {
			continue
		}

		// 检查节点配置
		config, ok := tplNode.Config.(*ApprovalNodeConfig)
		if !ok {
			continue
		}

		// 检查审批人配置
		if config.ApproverConfig == nil {
			continue
		}

		// 检查获取时机
		if config.ApproverConfig.GetTiming() != ApproverTimingOnCreate {
			continue
		}

		// 如果是动态审批人配置,需要设置 HTTPClient
		if dynamicConfig, ok := config.ApproverConfig.(*DynamicApproverConfig); ok {
			if dynamicConfig.HTTPClient == nil {
				dynamicConfig.HTTPClient = httpClient
			}
		}

		// 创建节点上下文
		ctx := &NodeContext{
			Task:    tsk,
			Node:    tplNode,
			Params:  tsk.Params,
			Outputs: make(map[string]json.RawMessage),
			Cache:   NewContextCache(),
		}

		// 获取审批人列表
		approvers, err := config.ApproverConfig.GetApprovers(ctx)
		if err != nil {
			// 如果获取失败,记录错误但不阻止任务创建
			// 审批人可以在节点激活时重新获取
			continue
		}

		// 保存审批人列表
		if tsk.Approvers == nil {
			tsk.Approvers = make(map[string][]string)
		}
		tsk.Approvers[tplNode.ID] = approvers
	}

	return nil
}

