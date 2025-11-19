package node_test

import (
	"encoding/json"
	"testing"

	"github.com/mautops/approval-kit/internal/node"
	"github.com/mautops/approval-kit/internal/task"
)

// TestStringConditionConfig 测试字符串匹配条件配置
func TestStringConditionConfig(t *testing.T) {
	config := &node.StringConditionConfig{
		Field:    "department",
		Operator: "eq",
		Value:    "engineering",
		Source:   "task_params",
	}

	if config.ConditionType() != "string" {
		t.Errorf("StringConditionConfig.ConditionType() = %q, want %q", config.ConditionType(), "string")
	}

	if config.Field != "department" {
		t.Errorf("StringConditionConfig.Field = %q, want %q", config.Field, "department")
	}

	if config.Operator != "eq" {
		t.Errorf("StringConditionConfig.Operator = %q, want %q", config.Operator, "eq")
	}

	if config.Value != "engineering" {
		t.Errorf("StringConditionConfig.Value = %q, want %q", config.Value, "engineering")
	}
}

// TestStringConditionEvaluator 测试字符串匹配条件评估器
func TestStringConditionEvaluator(t *testing.T) {
	evaluator := node.NewStringConditionEvaluator()

	// 测试 Supports 方法
	if !evaluator.Supports("string") {
		t.Error("StringConditionEvaluator should support 'string' type")
	}

	if evaluator.Supports("numeric") {
		t.Error("StringConditionEvaluator should not support 'numeric' type")
	}
}

// TestStringConditionEvaluate 测试字符串匹配条件评估
func TestStringConditionEvaluate(t *testing.T) {
	evaluator := node.NewStringConditionEvaluator()

	tests := []struct {
		name      string
		condition *node.Condition
		ctx       *node.NodeContext
		want      bool
		wantErr   bool
	}{
		{
			name: "eq operator: exact match",
			condition: &node.Condition{
				Type: "string",
				Config: &node.StringConditionConfig{
					Field:    "department",
					Operator: "eq",
					Value:    "engineering",
					Source:   "task_params",
				},
			},
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"department": "engineering"}`),
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "eq operator: no match",
			condition: &node.Condition{
				Type: "string",
				Config: &node.StringConditionConfig{
					Field:    "department",
					Operator: "eq",
					Value:    "engineering",
					Source:   "task_params",
				},
			},
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"department": "sales"}`),
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "contains operator: contains substring",
			condition: &node.Condition{
				Type: "string",
				Config: &node.StringConditionConfig{
					Field:    "description",
					Operator: "contains",
					Value:    "urgent",
					Source:   "task_params",
				},
			},
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"description": "This is urgent task"}`),
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "contains operator: no substring",
			condition: &node.Condition{
				Type: "string",
				Config: &node.StringConditionConfig{
					Field:    "description",
					Operator: "contains",
					Value:    "urgent",
					Source:   "task_params",
				},
			},
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"description": "This is normal task"}`),
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "starts_with operator",
			condition: &node.Condition{
				Type: "string",
				Config: &node.StringConditionConfig{
					Field:    "code",
					Operator: "starts_with",
					Value:    "APP",
					Source:   "task_params",
				},
			},
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"code": "APP-001"}`),
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "ends_with operator",
			condition: &node.Condition{
				Type: "string",
				Config: &node.StringConditionConfig{
					Field:    "code",
					Operator: "ends_with",
					Value:    "001",
					Source:   "task_params",
				},
			},
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"code": "APP-001"}`),
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "field not found",
			condition: &node.Condition{
				Type: "string",
				Config: &node.StringConditionConfig{
					Field:    "missing",
					Operator: "eq",
					Value:    "value",
					Source:   "task_params",
				},
			},
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"department": "engineering"}`),
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid operator",
			condition: &node.Condition{
				Type: "string",
				Config: &node.StringConditionConfig{
					Field:    "department",
					Operator: "invalid",
					Value:    "engineering",
					Source:   "task_params",
				},
			},
			ctx: &node.NodeContext{
				Task: &task.Task{
					Params: json.RawMessage(`{"department": "engineering"}`),
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "non-string value",
			condition: &node.Condition{
				Type: "string",
				Config: &node.StringConditionConfig{
					Field:    "amount",
					Operator: "eq",
					Value:    "engineering",
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
				Type: "string",
				Config: &node.StringConditionConfig{
					Field:    "status",
					Operator: "eq",
					Value:    "approved",
					Source:   "node_outputs",
					NodeID:   "node-001",
				},
			},
			ctx: &node.NodeContext{
				Task: &task.Task{
					NodeOutputs: map[string]json.RawMessage{
						"node-001": json.RawMessage(`{"status": "approved"}`),
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
				t.Errorf("StringConditionEvaluator.Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.want {
				t.Errorf("StringConditionEvaluator.Evaluate() = %v, want %v", result, tt.want)
			}
		})
	}
}

