package node_test

import (
	"encoding/json"
	"testing"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
)

// TestNumericConditionConfig 测试数值比较条件配置
func TestNumericConditionConfig(t *testing.T) {
	config := &node.NumericConditionConfig{
		Field:    "amount",
		Operator: "gt",
		Value:    1000.0,
		Source:   "task_params",
	}

	if config.ConditionType() != "numeric" {
		t.Errorf("NumericConditionConfig.ConditionType() = %q, want %q", config.ConditionType(), "numeric")
	}

	if config.Field != "amount" {
		t.Errorf("NumericConditionConfig.Field = %q, want %q", config.Field, "amount")
	}

	if config.Operator != "gt" {
		t.Errorf("NumericConditionConfig.Operator = %q, want %q", config.Operator, "gt")
	}

	if config.Value != 1000.0 {
		t.Errorf("NumericConditionConfig.Value = %v, want %v", config.Value, 1000.0)
	}
}

// TestNumericConditionEvaluator 测试数值比较条件评估器
func TestNumericConditionEvaluator(t *testing.T) {
	evaluator := node.NewNumericConditionEvaluator()

	// 测试 Supports 方法
	if !evaluator.Supports("numeric") {
		t.Error("NumericConditionEvaluator should support 'numeric' type")
	}

	if evaluator.Supports("string") {
		t.Error("NumericConditionEvaluator should not support 'string' type")
	}
}

// TestNumericConditionEvaluate 测试数值比较条件评估
func TestNumericConditionEvaluate(t *testing.T) {
	evaluator := node.NewNumericConditionEvaluator()

	tests := []struct {
		name      string
		condition *node.Condition
		ctx       *node.NodeContext
		want      bool
		wantErr   bool
	}{
		{
			name: "gt operator: 2000 > 1000",
			condition: &node.Condition{
				Type: "numeric",
				Config: &node.NumericConditionConfig{
					Field:    "amount",
					Operator: "gt",
					Value:    1000.0,
					Source:   "task_params",
				},
			},
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"amount": 2000}`),
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "gt operator: 500 < 1000",
			condition: &node.Condition{
				Type: "numeric",
				Config: &node.NumericConditionConfig{
					Field:    "amount",
					Operator: "gt",
					Value:    1000.0,
					Source:   "task_params",
				},
			},
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"amount": 500}`),
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "lt operator: 500 < 1000",
			condition: &node.Condition{
				Type: "numeric",
				Config: &node.NumericConditionConfig{
					Field:    "amount",
					Operator: "lt",
					Value:    1000.0,
					Source:   "task_params",
				},
			},
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"amount": 500}`),
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "eq operator: 1000 == 1000",
			condition: &node.Condition{
				Type: "numeric",
				Config: &node.NumericConditionConfig{
					Field:    "amount",
					Operator: "eq",
					Value:    1000.0,
					Source:   "task_params",
				},
			},
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"amount": 1000}`),
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "gte operator: 1000 >= 1000",
			condition: &node.Condition{
				Type: "numeric",
				Config: &node.NumericConditionConfig{
					Field:    "amount",
					Operator: "gte",
					Value:    1000.0,
					Source:   "task_params",
				},
			},
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"amount": 1000}`),
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "lte operator: 1000 <= 1000",
			condition: &node.Condition{
				Type: "numeric",
				Config: &node.NumericConditionConfig{
					Field:    "amount",
					Operator: "lte",
					Value:    1000.0,
					Source:   "task_params",
				},
			},
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"amount": 1000}`),
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "field not found",
			condition: &node.Condition{
				Type: "numeric",
				Config: &node.NumericConditionConfig{
					Field:    "missing",
					Operator: "gt",
					Value:    1000.0,
					Source:   "task_params",
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
		{
			name: "invalid operator",
			condition: &node.Condition{
				Type: "numeric",
				Config: &node.NumericConditionConfig{
					Field:    "amount",
					Operator: "invalid",
					Value:    1000.0,
					Source:   "task_params",
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
		{
			name: "non-numeric value",
			condition: &node.Condition{
				Type: "numeric",
				Config: &node.NumericConditionConfig{
					Field:    "amount",
					Operator: "gt",
					Value:    1000.0,
					Source:   "task_params",
				},
			},
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"amount": "not-a-number"}`),
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "from node_outputs",
			condition: &node.Condition{
				Type: "numeric",
				Config: &node.NumericConditionConfig{
					Field:    "result",
					Operator: "gt",
					Value:    50.0,
					Source:   "node_outputs",
					NodeID:   "node-001",
				},
			},
			ctx: &node.NodeContext{
				Task: &task.Task{
					NodeOutputs: map[string]json.RawMessage{
						"node-001": json.RawMessage(`{"result": 75}`),
					},
				},
			},
			want:    true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.Evaluate(tt.condition, tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NumericConditionEvaluator.Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.want {
				t.Errorf("NumericConditionEvaluator.Evaluate() = %v, want %v", result, tt.want)
			}
		})
	}
}

