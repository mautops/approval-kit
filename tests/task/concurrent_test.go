package task_test

import (
	"sync"
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/task"
)

// TestTaskConcurrentRead 测试并发读取的安全性
func TestTaskConcurrentRead(t *testing.T) {
	tsk := &task.Task{
		ID:    "task-001",
		State: task.TaskStatePending,
	}

	var wg sync.WaitGroup
	concurrency := 100
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			// 并发读取
			_ = tsk.GetState()
			_ = tsk.GetCurrentNode()
			_ = tsk.GetUpdatedAt()
			_ = tsk.GetStateHistory()
			_ = tsk.GetRecords()
			_ = tsk.Snapshot()
		}()
	}

	wg.Wait()
	// 如果没有 race condition,测试应该通过
}

// TestTaskConcurrentUpdate 测试并发更新的安全性
func TestTaskConcurrentUpdate(t *testing.T) {
	tsk := &task.Task{
		ID:    "task-001",
		State: task.TaskStatePending,
	}

	var wg sync.WaitGroup
	concurrency := 50
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func(id int) {
			defer wg.Done()
			// 并发更新
			err := tsk.Update(func(t *task.Task) error {
				t.UpdatedAt = time.Now()
				return nil
			})
			if err != nil {
				t.Errorf("Update() failed in goroutine %d: %v", id, err)
			}
		}(i)
	}

	wg.Wait()
	// 如果没有 race condition,测试应该通过
}

// TestTaskSnapshotIsolation 测试快照的隔离性
func TestTaskSnapshotIsolation(t *testing.T) {
	tsk := &task.Task{
		ID:    "task-001",
		State: task.TaskStatePending,
	}

	// 创建快照
	snapshot := tsk.Snapshot()

	// 修改原任务
	err := tsk.Update(func(t *task.Task) error {
		t.State = task.TaskStateSubmitted
		return nil
	})
	if err != nil {
		t.Fatalf("Update() failed: %v", err)
	}

	// 验证快照未被修改
	if snapshot.State != task.TaskStatePending {
		t.Error("Snapshot should not be affected by task updates")
	}

	// 验证原任务已修改
	if tsk.GetState() != task.TaskStateSubmitted {
		t.Error("Original task should be updated")
	}
}

