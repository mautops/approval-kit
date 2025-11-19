package node_test

import (
	"encoding/json"
	"testing"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
)

// TestCompositeConditionConfig 测试组合条件配置
func TestCompositeConditionConfig(t *testing.T) {
	config := &node.CompositeConditionConfig{
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
	}

	if config.ConditionType() != "composite" {
		t.Errorf("CompositeConditionConfig.ConditionType() = %q, want %q", config.ConditionType(), "composite")
	}

	if config.Operator != "and" {
		t.Errorf("CompositeConditionConfig.Operator = %q, want %q", config.Operator, "and")
	}

	if len(config.Conditions) != 2 {
		t.Errorf("CompositeConditionConfig.Conditions length = %d, want %d", len(config.Conditions), 2)
	}
}

// TestCompositeConditionEvaluator 测试组合条件评估器
func TestCompositeConditionEvaluator(t *testing.T) {
	registry := node.NewConditionEvaluatorRegistry()
	evaluator := node.NewCompositeConditionEvaluator(registry)

	// 测试 Supports 方法
	if !evaluator.Supports("composite") {
		t.Error("CompositeConditionEvaluator should support 'composite' type")
	}

	if evaluator.Supports("numeric") {
		t.Error("CompositeConditionEvaluator should not support 'numeric' type")
	}
}

// TestCompositeConditionEvaluate 测试组合条件评估
func TestCompositeConditionEvaluate(t *testing.T) {
	registry := node.NewConditionEvaluatorRegistry()
	evaluator := node.NewCompositeConditionEvaluator(registry)

	tests := []struct {
		name      string
		condition *node.Condition
		ctx       *node.NodeContext
		want      bool
		wantErr   bool
	}{
		{
			name: "and operator: all true",
			condition: &node.Condition{
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
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"amount": 2000, "department": "engineering"}`),
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "and operator: one false",
			condition: &node.Condition{
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
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"amount": 500, "department": "engineering"}`),
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "or operator: one true",
			condition: &node.Condition{
				Type: "composite",
				Config: &node.CompositeConditionConfig{
					Operator: "or",
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
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"amount": 500, "department": "engineering"}`),
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "or operator: all false",
			condition: &node.Condition{
				Type: "composite",
				Config: &node.CompositeConditionConfig{
					Operator: "or",
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
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"amount": 500, "department": "sales"}`),
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "nested composite conditions",
			condition: &node.Condition{
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
							Type: "composite",
							Config: &node.CompositeConditionConfig{
								Operator: "or",
								Conditions: []*node.Condition{
									{
										Type: "string",
										Config: &node.StringConditionConfig{
											Field:    "department",
											Operator: "eq",
											Value:    "engineering",
											Source:   "task_params",
										},
									},
									{
										Type: "string",
										Config: &node.StringConditionConfig{
											Field:    "department",
											Operator: "eq",
											Value:    "product",
											Source:   "task_params",
										},
									},
								},
							},
						},
					},
				},
			},
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"amount": 2000, "department": "engineering"}`),
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "empty conditions",
			condition: &node.Condition{
				Type: "composite",
				Config: &node.CompositeConditionConfig{
					Operator:   "and",
					Conditions: []*node.Condition{},
				},
			},
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{}`),
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid operator",
			condition: &node.Condition{
				Type: "composite",
				Config: &node.CompositeConditionConfig{
					Operator: "invalid",
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
					},
				},
			},
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"amount": 2000}`),
				},
			},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.Evaluate(tt.condition, tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompositeConditionEvaluator.Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.want {
				t.Errorf("CompositeConditionEvaluator.Evaluate() = %v, want %v", result, tt.want)
			}
		})
	}
}

