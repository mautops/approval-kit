package task_test

import (
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/task"
)

// TestRecordValidation 测试审批记录数据结构验证
func TestRecordValidation(t *testing.T) {
	tests := []struct {
		name      string
		record    *task.Record
		wantErr   bool
		wantErrMsg string
	}{
		{
			name: "valid record",
			record: &task.Record{
				ID:          "record-001",
				TaskID:      "task-001",
				NodeID:      "node-001",
				Approver:    "user-001",
				Result:      "approve",
				Comment:     "approved",
				CreatedAt:   time.Now(),
				Attachments: []string{"file1.pdf"},
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			record: &task.Record{
				TaskID:    "task-001",
				NodeID:    "node-001",
				Approver:  "user-001",
				Result:    "approve",
				CreatedAt: time.Now(),
			},
			wantErr:   true,
			wantErrMsg: "record ID is required",
		},
		{
			name: "missing TaskID",
			record: &task.Record{
				ID:        "record-001",
				NodeID:    "node-001",
				Approver:  "user-001",
				Result:    "approve",
				CreatedAt: time.Now(),
			},
			wantErr:   true,
			wantErrMsg: "record TaskID is required",
		},
		{
			name: "missing NodeID",
			record: &task.Record{
				ID:        "record-001",
				TaskID:    "task-001",
				Approver:  "user-001",
				Result:    "approve",
				CreatedAt: time.Now(),
			},
			wantErr:   true,
			wantErrMsg: "record NodeID is required",
		},
		{
			name: "missing Approver",
			record: &task.Record{
				ID:        "record-001",
				TaskID:    "task-001",
				NodeID:    "node-001",
				Result:    "approve",
				CreatedAt: time.Now(),
			},
			wantErr:   true,
			wantErrMsg: "record Approver is required",
		},
		{
			name: "missing Result",
			record: &task.Record{
				ID:        "record-001",
				TaskID:    "task-001",
				NodeID:    "node-001",
				Approver:  "user-001",
				CreatedAt: time.Now(),
			},
			wantErr:   true,
			wantErrMsg: "record Result is required",
		},
		{
			name: "invalid Result",
			record: &task.Record{
				ID:        "record-001",
				TaskID:    "task-001",
				NodeID:    "node-001",
				Approver:  "user-001",
				Result:    "invalid",
				CreatedAt: time.Now(),
			},
			wantErr:   true,
			wantErrMsg: "invalid record Result",
		},
		{
			name: "zero CreatedAt",
			record: &task.Record{
				ID:        "record-001",
				TaskID:    "task-001",
				NodeID:    "node-001",
				Approver:  "user-001",
				Result:    "approve",
				CreatedAt: time.Time{},
			},
			wantErr:   true,
			wantErrMsg: "record CreatedAt is required",
		},
		{
			name: "valid result: approve",
			record: &task.Record{
				ID:        "record-001",
				TaskID:    "task-001",
				NodeID:    "node-001",
				Approver:  "user-001",
				Result:    "approve",
				CreatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "valid result: reject",
			record: &task.Record{
				ID:        "record-001",
				TaskID:    "task-001",
				NodeID:    "node-001",
				Approver:  "user-001",
				Result:    "reject",
				CreatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "valid result: transfer",
			record: &task.Record{
				ID:        "record-001",
				TaskID:    "task-001",
				NodeID:    "node-001",
				Approver:  "user-001",
				Result:    "transfer",
				CreatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "record with attachments",
			record: &task.Record{
				ID:          "record-001",
				TaskID:      "task-001",
				NodeID:      "node-001",
				Approver:    "user-001",
				Result:      "approve",
				Comment:     "approved with attachments",
				CreatedAt:   time.Now(),
				Attachments: []string{"file1.pdf", "file2.pdf"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.record.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Record.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				if tt.wantErrMsg != "" && err.Error() != tt.wantErrMsg {
					// 检查错误消息是否包含期望的关键词
					if !contains(err.Error(), tt.wantErrMsg) {
						t.Errorf("Record.Validate() error message = %q, want contains %q", err.Error(), tt.wantErrMsg)
					}
				}
			}
		})
	}
}

// contains 检查字符串是否包含子字符串
func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

