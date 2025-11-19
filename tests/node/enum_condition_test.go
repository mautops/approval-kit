package node_test

import (
	"encoding/json"
	"testing"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
)

// TestEnumConditionConfig 测试枚举判断条件配置
func TestEnumConditionConfig(t *testing.T) {
	config := &node.EnumConditionConfig{
		Field:    "status",
		Operator: "in",
		Values:   []string{"pending", "approved"},
		Source:   "task_params",
	}

	if config.ConditionType() != "enum" {
		t.Errorf("EnumConditionConfig.ConditionType() = %q, want %q", config.ConditionType(), "enum")
	}

	if config.Field != "status" {
		t.Errorf("EnumConditionConfig.Field = %q, want %q", config.Field, "status")
	}

	if len(config.Values) != 2 {
		t.Errorf("EnumConditionConfig.Values length = %d, want %d", len(config.Values), 2)
	}
}

// TestEnumConditionEvaluator 测试枚举判断条件评估器
func TestEnumConditionEvaluator(t *testing.T) {
	evaluator := node.NewEnumConditionEvaluator()

	// 测试 Supports 方法
	if !evaluator.Supports("enum") {
		t.Error("EnumConditionEvaluator should support 'enum' type")
	}

	if evaluator.Supports("numeric") {
		t.Error("EnumConditionEvaluator should not support 'numeric' type")
	}
}

// TestEnumConditionEvaluate 测试枚举判断条件评估
func TestEnumConditionEvaluate(t *testing.T) {
	evaluator := node.NewEnumConditionEvaluator()

	tests := []struct {
		name      string
		condition *node.Condition
		ctx       *node.NodeContext
		want      bool
		wantErr   bool
	}{
		{
			name: "in operator: value in list",
			condition: &node.Condition{
				Type: "enum",
				Config: &node.EnumConditionConfig{
					Field:    "status",
					Operator: "in",
					Values:   []string{"pending", "approved", "rejected"},
					Source:   "task_params",
				},
			},
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"status": "approved"}`),
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "in operator: value not in list",
			condition: &node.Condition{
				Type: "enum",
				Config: &node.EnumConditionConfig{
					Field:    "status",
					Operator: "in",
					Values:   []string{"pending", "approved"},
					Source:   "task_params",
				},
			},
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"status": "rejected"}`),
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "not_in operator: value not in list",
			condition: &node.Condition{
				Type: "enum",
				Config: &node.EnumConditionConfig{
					Field:    "status",
					Operator: "not_in",
					Values:   []string{"pending", "approved"},
					Source:   "task_params",
				},
			},
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"status": "rejected"}`),
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "not_in operator: value in list",
			condition: &node.Condition{
				Type: "enum",
				Config: &node.EnumConditionConfig{
					Field:    "status",
					Operator: "not_in",
					Values:   []string{"pending", "approved"},
					Source:   "task_params",
				},
			},
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"status": "approved"}`),
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "field not found",
			condition: &node.Condition{
				Type: "enum",
				Config: &node.EnumConditionConfig{
					Field:    "missing",
					Operator: "in",
					Values:   []string{"value1", "value2"},
					Source:   "task_params",
				},
			},
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"status": "approved"}`),
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid operator",
			condition: &node.Condition{
				Type: "enum",
				Config: &node.EnumConditionConfig{
					Field:    "status",
					Operator: "invalid",
					Values:   []string{"pending", "approved"},
					Source:   "task_params",
				},
			},
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"status": "approved"}`),
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "non-string value",
			condition: &node.Condition{
				Type: "enum",
				Config: &node.EnumConditionConfig{
					Field:    "amount",
					Operator: "in",
					Values:   []string{"value1"},
					Source:   "task_params",
				},
			},
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"amount": 1000}`),
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "from node_outputs",
			condition: &node.Condition{
				Type: "enum",
				Config: &node.EnumConditionConfig{
					Field:    "result",
					Operator: "in",
					Values:   []string{"success", "failed"},
					Source:   "node_outputs",
					NodeID:   "node-001",
				},
			},
			ctx: &node.NodeContext{
				Task: &task.Task{
					NodeOutputs: map[string]json.RawMessage{
						"node-001": json.RawMessage(`{"result": "success"}`),
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
				t.Errorf("EnumConditionEvaluator.Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.want {
				t.Errorf("EnumConditionEvaluator.Evaluate() = %v, want %v", result, tt.want)
			}
		})
	}
}

