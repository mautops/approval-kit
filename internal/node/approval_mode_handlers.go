package node

import (
	"github.com/mautops/approval-kit/internal/task"
)

// singleModeHandler 单人审批模式处理器
type singleModeHandler struct{}

func (h *singleModeHandler) Mode() ApprovalMode {
	return ApprovalModeSingle
}

func (h *singleModeHandler) CheckCompletion(approvers []string, approvals map[string]*task.Approval, config *ApprovalNodeConfig) (bool, *ApprovalResult) {
	// 单人审批模式必须只有一个审批人
	if len(approvers) != 1 {
		return false, nil
	}

	approver := approvers[0]
	approval, exists := approvals[approver]

	// 如果没有审批记录,审批未完成
	if !exists || approval == nil {
		return false, nil
	}

	// 根据审批结果决定下一步
	result := &ApprovalResult{
		Completed: true,
		Result:    approval.Result,
		NextNodeID: "", // 下一个节点由流程引擎根据边的定义决定
	}

	// 如果被拒绝,根据配置决定行为
	if approval.Result == "reject" {
		result = createRejectResult(config)
	}

	return true, result
}

// unanimousModeHandler 多人会签模式处理器
type unanimousModeHandler struct{}

func (h *unanimousModeHandler) Mode() ApprovalMode {
	return ApprovalModeUnanimous
}

func (h *unanimousModeHandler) CheckCompletion(approvers []string, approvals map[string]*task.Approval, config *ApprovalNodeConfig) (bool, *ApprovalResult) {
	// 多人会签模式必须至少有两个审批人
	if len(approvers) < 2 {
		return false, nil
	}

	// 统计审批结果
	approvedCount := 0
	hasRejection := false

	// 检查每个审批人的审批结果
	for _, approver := range approvers {
		approval, exists := approvals[approver]
		if !exists || approval == nil {
			// 有审批人还未审批,审批未完成
			return false, nil
		}

		if approval.Result == "approve" {
			approvedCount++
		} else if approval.Result == "reject" {
			hasRejection = true
		}
	}

	// 如果所有审批人都同意,审批完成
	if approvedCount == len(approvers) {
		return true, &ApprovalResult{
			Completed: true,
			Result:    "approve",
			NextNodeID: "", // 下一个节点由流程引擎根据边的定义决定
		}
	}

	// 如果有审批人拒绝,根据配置决定行为
	if hasRejection {
		return true, createRejectResult(config)
	}

	// 其他情况(不应该发生),审批未完成
	return false, nil
}

// orModeHandler 多人或签模式处理器
type orModeHandler struct{}

func (h *orModeHandler) Mode() ApprovalMode {
	return ApprovalModeOr
}

func (h *orModeHandler) CheckCompletion(approvers []string, approvals map[string]*task.Approval, config *ApprovalNodeConfig) (bool, *ApprovalResult) {
	// 多人或签模式必须至少有两个审批人
	if len(approvers) < 2 {
		return false, nil
	}

	// 统计审批结果
	rejectedCount := 0
	hasApproval := false

	// 检查每个审批人的审批结果
	for _, approver := range approvers {
		approval, exists := approvals[approver]
		if !exists || approval == nil {
			// 有审批人还未审批,继续等待
			continue
		}

		if approval.Result == "approve" {
			hasApproval = true
		} else if approval.Result == "reject" {
			rejectedCount++
		}
	}

	// 如果任意一人同意,审批完成
	if hasApproval {
		return true, &ApprovalResult{
			Completed: true,
			Result:    "approve",
			NextNodeID: "", // 下一个节点由流程引擎根据边的定义决定
		}
	}

	// 如果所有审批人都已审批且全部拒绝,根据配置决定行为
	if rejectedCount == len(approvers) {
		return true, createRejectResult(config)
	}

	// 其他情况(有审批人还未审批),审批未完成
	return false, nil
}

// proportionalModeHandler 比例会签模式处理器
type proportionalModeHandler struct{}

func (h *proportionalModeHandler) Mode() ApprovalMode {
	return ApprovalModeProportional
}

func (h *proportionalModeHandler) CheckCompletion(approvers []string, approvals map[string]*task.Approval, config *ApprovalNodeConfig) (bool, *ApprovalResult) {
	// 比例会签模式必须至少有两个审批人
	if len(approvers) < 2 {
		return false, nil
	}

	// 验证比例阈值配置
	if config.ProportionalThreshold == nil {
		return false, nil
	}

	threshold := config.ProportionalThreshold
	if threshold.Required <= 0 || threshold.Total <= 0 {
		return false, nil
	}

	// 统计审批结果
	approvedCount := 0

	// 检查每个审批人的审批结果
	for _, approver := range approvers {
		approval, exists := approvals[approver]
		if !exists || approval == nil {
			// 审批人还未审批,继续等待
			continue
		}

		if approval.Result == "approve" {
			approvedCount++
		}
		// 拒绝的数量不需要统计,因为只要达到阈值即可通过
	}

	// 如果同意的数量达到阈值,审批完成
	if approvedCount >= threshold.Required {
		return true, &ApprovalResult{
			Completed: true,
			Result:    "approve",
			NextNodeID: "", // 下一个节点由流程引擎根据边的定义决定
		}
	}

	// 其他情况(有审批人还未审批或未达到阈值),审批未完成
	// 比例会签模式下,只要达到阈值即可通过,未达到阈值时继续等待
	return false, nil
}

// sequentialModeHandler 顺序审批模式处理器
type sequentialModeHandler struct{}

func (h *sequentialModeHandler) Mode() ApprovalMode {
	return ApprovalModeSequential
}

func (h *sequentialModeHandler) CheckCompletion(approvers []string, approvals map[string]*task.Approval, config *ApprovalNodeConfig) (bool, *ApprovalResult) {
	// 顺序审批模式必须至少有两个审批人
	if len(approvers) < 2 {
		return false, nil
	}

	// 按顺序检查每个审批人的审批结果
	for i, approver := range approvers {
		approval, exists := approvals[approver]

		// 如果当前审批人还未审批
		if !exists || approval == nil {
			// 如果是第一个审批人,审批未完成
			if i == 0 {
				return false, nil
			}

			// 检查前一个审批人是否已同意
			prevApprover := approvers[i-1]
			prevApproval, prevExists := approvals[prevApprover]

			// 如果前一个审批人还未审批,审批未完成(必须按顺序)
			if !prevExists || prevApproval == nil {
				return false, nil
			}

			// 如果前一个审批人已同意,当前审批人还未审批,审批未完成(等待当前审批人)
			if prevApproval.Result == "approve" {
				return false, nil
			}

			// 如果前一个审批人已拒绝,根据配置决定行为
			if prevApproval.Result == "reject" {
				return true, createRejectResult(config)
			}
		} else {
			// 当前审批人已审批
			if approval.Result == "reject" {
				// 如果当前审批人拒绝,根据配置决定行为
				return true, createRejectResult(config)
			}

			// 如果当前审批人同意
			if approval.Result == "approve" {
				// 如果是最后一个审批人,所有审批人都已按顺序同意
				if i == len(approvers)-1 {
					return true, &ApprovalResult{
						Completed: true,
						Result:    "approve",
						NextNodeID: "", // 下一个节点由流程引擎根据边的定义决定
					}
				}
				// 如果不是最后一个审批人,继续检查下一个审批人
			}
		}
	}

	// 其他情况,审批未完成
	return false, nil
}

// createRejectResult 创建拒绝结果
// 根据拒绝后行为配置决定下一步
func createRejectResult(config *ApprovalNodeConfig) *ApprovalResult {
	result := &ApprovalResult{
		Completed: true,
		Result:    "reject",
		NextNodeID: "", // 下一个节点由流程引擎根据边的定义决定
	}

	// 根据拒绝后行为配置决定下一步
	switch config.RejectBehavior {
	case RejectBehaviorTerminate:
		// 终止流程,不指定下一个节点
		result.NextNodeID = ""
	case RejectBehaviorRollback:
		// 回退到上一节点(由流程引擎处理)
		result.NextNodeID = ""
	case RejectBehaviorJump:
		// 跳转到指定节点(由流程引擎处理)
		result.NextNodeID = ""
	}

	return result
}

