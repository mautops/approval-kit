package template

import (
	"fmt"

	"github.com/mautops/approval-kit/internal/errors"
)

// Validate 验证模板的有效性
// 验证规则:
// 1. ID 和 Name 不能为空
// 2. 必须有且仅有一个开始节点
// 3. 所有边引用的节点必须存在
func (t *Template) Validate() error {
	// 验证 ID
	if t.ID == "" {
		return fmt.Errorf("%w: template ID is required", errors.ErrInvalidTemplate)
	}

	// 验证 Name
	if t.Name == "" {
		return fmt.Errorf("%w: template Name is required", errors.ErrInvalidTemplate)
	}

	// 验证节点不为空
	if len(t.Nodes) == 0 {
		return fmt.Errorf("%w: template must have at least one node", errors.ErrInvalidTemplate)
	}

	// 验证必须有且仅有一个开始节点
	startNodeCount := 0
	for _, node := range t.Nodes {
		if node.Type == NodeTypeStart {
			startNodeCount++
		}
	}

	if startNodeCount == 0 {
		return fmt.Errorf("%w: template must have exactly one start node, found 0", errors.ErrInvalidTemplate)
	}

	if startNodeCount > 1 {
		return fmt.Errorf("%w: template must have exactly one start node, found %d", errors.ErrInvalidTemplate, startNodeCount)
	}

	// 验证所有边引用的节点必须存在
	for i, edge := range t.Edges {
		// 验证 From 节点存在
		if _, exists := t.Nodes[edge.From]; !exists {
			return fmt.Errorf("%w: edge[%d] references non-existent node: %q", errors.ErrInvalidTemplate, i, edge.From)
		}

		// 验证 To 节点存在
		if _, exists := t.Nodes[edge.To]; !exists {
			return fmt.Errorf("%w: edge[%d] references non-existent node: %q", errors.ErrInvalidTemplate, i, edge.To)
		}
	}

	return nil
}

