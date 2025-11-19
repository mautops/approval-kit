package statemachine_test

import (
	"sync"
	"testing"

	"github.com/mautops/approval-kit/internal/statemachine"
	"github.com/mautops/approval-kit/internal/task"
)

// TestStateMachineConcurrentRead 测试并发读取的安全性
func TestStateMachineConcurrentRead(t *testing.T) {
	sm := statemachine.NewStateMachine()

	// 并发读取 CanTransition
	var wg sync.WaitGroup
	concurrency := 100
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			// 并发调用 CanTransition
			_ = sm.CanTransition(task.TaskStatePending, task.TaskStateSubmitted)
			_ = sm.CanTransition(task.TaskStateSubmitted, task.TaskStateApproving)
			_ = sm.CanTransition(task.TaskStateApproving, task.TaskStateApproved)
		}()
	}

	wg.Wait()
	// 如果没有 race condition,测试应该通过
}

// TestStateMachineConcurrentGetValidTransitions 测试并发获取有效转换的安全性
func TestStateMachineConcurrentGetValidTransitions(t *testing.T) {
	sm := statemachine.NewStateMachine()

	// 并发读取 GetValidTransitions
	var wg sync.WaitGroup
	concurrency := 100
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			// 并发调用 GetValidTransitions
			_ = sm.GetValidTransitions(task.TaskStatePending)
			_ = sm.GetValidTransitions(task.TaskStateSubmitted)
			_ = sm.GetValidTransitions(task.TaskStateApproving)
		}()
	}

	wg.Wait()
	// 如果没有 race condition,测试应该通过
}

// TestStateMachineConcurrentTransition 测试并发状态转换的安全性
func TestStateMachineConcurrentTransition(t *testing.T) {
	sm := statemachine.NewStateMachine()

	// 创建多个任务,并发执行状态转换
	var wg sync.WaitGroup
	concurrency := 50
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func(id int) {
			defer wg.Done()
			// 每个 goroutine 创建自己的任务并执行转换
			testTask := createTestTask(task.TaskStatePending)
			newTask, err := sm.Transition(testTask, task.TaskStateSubmitted, "concurrent test")
			if err != nil {
				t.Errorf("Transition() failed in goroutine %d: %v", id, err)
				return
			}
			if newTask == nil {
				t.Errorf("Transition() returned nil in goroutine %d", id)
				return
			}
			if newTask.GetState() != task.TaskStateSubmitted {
				t.Errorf("Transition() wrong state in goroutine %d: got %v, want %v",
					id, newTask.GetState(), task.TaskStateSubmitted)
			}
		}(i)
	}

	wg.Wait()
	// 如果没有 race condition,测试应该通过
}

// TestStateMachineConcurrentMixedOperations 测试混合并发操作的安全性
func TestStateMachineConcurrentMixedOperations(t *testing.T) {
	sm := statemachine.NewStateMachine()

	var wg sync.WaitGroup
	concurrency := 100
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func(id int) {
			defer wg.Done()

			// 混合调用不同的方法
			switch id % 3 {
			case 0:
				// CanTransition
				_ = sm.CanTransition(task.TaskStatePending, task.TaskStateSubmitted)
			case 1:
				// GetValidTransitions
				_ = sm.GetValidTransitions(task.TaskStatePending)
			case 2:
				// Transition
				testTask := createTestTask(task.TaskStatePending)
				_, _ = sm.Transition(testTask, task.TaskStateSubmitted, "mixed test")
			}
		}(i)
	}

	wg.Wait()
	// 如果没有 race condition,测试应该通过
}

// TestStateMachineMultipleInstances 测试多个状态机实例的并发安全性
func TestStateMachineMultipleInstances(t *testing.T) {
	// 创建多个状态机实例
	instances := make([]statemachine.StateMachine, 10)
	for i := range instances {
		instances[i] = statemachine.NewStateMachine()
	}

	var wg sync.WaitGroup
	concurrency := 100
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func(id int) {
			defer wg.Done()
			// 使用不同的状态机实例
			sm := instances[id%len(instances)]
			_ = sm.CanTransition(task.TaskStatePending, task.TaskStateSubmitted)
			_ = sm.GetValidTransitions(task.TaskStatePending)
		}(i)
	}

	wg.Wait()
	// 如果没有 race condition,测试应该通过
}

// TestStateMachineTransitionIsolation 测试状态转换的隔离性
// 验证不同任务的状态转换不会相互影响
func TestStateMachineTransitionIsolation(t *testing.T) {
	sm := statemachine.NewStateMachine()

	var wg sync.WaitGroup
	concurrency := 50
	wg.Add(concurrency)

	// 每个 goroutine 创建独立的任务并执行转换
	for i := 0; i < concurrency; i++ {
		go func(id int) {
			defer wg.Done()

			// 创建独立的任务
			task1 := createTestTask(task.TaskStatePending)
			task2 := createTestTask(task.TaskStatePending)

			// 对 task1 执行转换
			newTask1, err1 := sm.Transition(task1, task.TaskStateSubmitted, "task1")
			if err1 != nil {
				t.Errorf("Task1 transition failed in goroutine %d: %v", id, err1)
				return
			}

			// 对 task2 执行转换
			newTask2, err2 := sm.Transition(task2, task.TaskStateSubmitted, "task2")
			if err2 != nil {
				t.Errorf("Task2 transition failed in goroutine %d: %v", id, err2)
				return
			}

			// 验证原任务未被修改
			if task1.GetState() != task.TaskStatePending {
				t.Errorf("Task1 original state modified in goroutine %d", id)
			}
			if task2.GetState() != task.TaskStatePending {
				t.Errorf("Task2 original state modified in goroutine %d", id)
			}

			// 验证新任务状态正确
			if newTask1.GetState() != task.TaskStateSubmitted {
				t.Errorf("Task1 new state incorrect in goroutine %d", id)
			}
			if newTask2.GetState() != task.TaskStateSubmitted {
				t.Errorf("Task2 new state incorrect in goroutine %d", id)
			}
		}(i)
	}

	wg.Wait()
}

