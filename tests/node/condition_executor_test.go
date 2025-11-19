package node_test

import (
	"encoding/json"
	"testing"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
	"github.com/mautops/approval-kit/internal/template"
)

// TestConditionNodeExecutor 测试条件节点执行器
func TestConditionNodeExecutor(t *testing.T) {
	executor := node.NewConditionNodeExecutor()

	// 测试 NodeType 方法
	if executor.NodeType() != template.NodeTypeCondition {
		t.Errorf("ConditionNodeExecutor.NodeType() = %q, want %q", executor.NodeType(), template.NodeTypeCondition)
	}
}

// TestConditionNodeExecutorExecute 测试条件节点执行
func TestConditionNodeExecutorExecute(t *testing.T) {
	executor := node.NewConditionNodeExecutor()

	tests := []struct {
		name      string
		ctx       *node.NodeContext
		wantNode  string
		wantErr   bool
	}{
		{
			name: "condition true: jump to true node",
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"amount": 2000}`),
				},
				Node: &template.Node{
					ID:   "condition-001",
					Type: template.NodeTypeCondition,
					Config: &node.ConditionNodeConfig{
						Condition: &node.Condition{
							Type: "numeric",
							Config: &node.NumericConditionConfig{
								Field:    "amount",
								Operator: "gt",
								Value:    1000.0,
								Source:   "task_params",
							},
						},
						TrueNodeID:  "approval-001",
						FalseNodeID: "reject-001",
					},
				},
				Params:  json.RawMessage(`{"amount": 2000}`),
				Outputs: make(map[string]json.RawMessage),
				Cache:   node.NewContextCache(),
			},
			wantNode: "approval-001",
			wantErr:  false,
		},
		{
			name: "condition false: jump to false node",
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"amount": 500}`),
				},
				Node: &template.Node{
					ID:   "condition-001",
					Type: template.NodeTypeCondition,
					Config: &node.ConditionNodeConfig{
						Condition: &node.Condition{
							Type: "numeric",
							Config: &node.NumericConditionConfig{
								Field:    "amount",
								Operator: "gt",
								Value:    1000.0,
								Source:   "task_params",
							},
						},
						TrueNodeID:  "approval-001",
						FalseNodeID: "reject-001",
					},
				},
				Params:  json.RawMessage(`{"amount": 500}`),
				Outputs: make(map[string]json.RawMessage),
				Cache:   node.NewContextCache(),
			},
			wantNode: "reject-001",
			wantErr:  false,
		},
		{
			name: "missing condition config",
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"amount": 2000}`),
				},
				Node: &template.Node{
					ID:   "condition-001",
					Type: template.NodeTypeCondition,
					Config: nil,
				},
				Params:  json.RawMessage(`{"amount": 2000}`),
				Outputs: make(map[string]json.RawMessage),
				Cache:   node.NewContextCache(),
			},
			wantNode: "",
			wantErr:  true,
		},
		{
			name: "condition evaluation error",
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{}`),
				},
				Node: &template.Node{
					ID:   "condition-001",
					Type: template.NodeTypeCondition,
					Config: &node.ConditionNodeConfig{
						Condition: &node.Condition{
							Type: "numeric",
							Config: &node.NumericConditionConfig{
								Field:    "missing",
								Operator: "gt",
								Value:    1000.0,
								Source:   "task_params",
							},
						},
						TrueNodeID:  "approval-001",
						FalseNodeID: "reject-001",
					},
				},
				Params:  json.RawMessage(`{}`),
				Outputs: make(map[string]json.RawMessage),
				Cache:   node.NewContextCache(),
			},
			wantNode: "",
			wantErr:  true,
		},
		{
			name: "composite condition",
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"amount": 2000, "department": "engineering"}`),
				},
				Node: &template.Node{
					ID:   "condition-001",
					Type: template.NodeTypeCondition,
					Config: &node.ConditionNodeConfig{
						Condition: &node.Condition{
							Type: "composite",
							Config: &node.CompositeConditionConfig{
								Operator: "and",
								Conditions: []*node.Condition{
									{
										Type: "numeric",
										Config: &node.NumericConditionConfig{
											Field:    "amount",
											Operator: "gt",
											Value:    1000.0,
											Source:   "task_params",
										},
									},
									{
										Type: "string",
										Config: &node.StringConditionConfig{
											Field:    "department",
											Operator: "eq",
											Value:    "engineering",
											Source:   "task_params",
										},
									},
								},
							},
						},
						TrueNodeID:  "approval-001",
						FalseNodeID: "reject-001",
					},
				},
				Params:  json.RawMessage(`{"amount": 2000, "department": "engineering"}`),
				Outputs: make(map[string]json.RawMessage),
				Cache:   node.NewContextCache(),
			},
			wantNode: "approval-001",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := executor.Execute(tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConditionNodeExecutor.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if result == nil {
					t.Error("ConditionNodeExecutor.Execute() returned nil result")
					return
				}
				if result.NextNodeID != tt.wantNode {
					t.Errorf("ConditionNodeExecutor.Execute() NextNodeID = %q, want %q", result.NextNodeID, tt.wantNode)
				}
			}
		})
	}
}

