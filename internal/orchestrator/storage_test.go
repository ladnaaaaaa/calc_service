package orchestrator

import (
	"sync"
	"testing"
)

func TestStoreConcurrency(t *testing.T) {
	store := NewStore()
	wg := sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			expr := &Expression{
				ID:     string(rune(n)),
				Status: "pending",
			}
			store.AddExpression(expr)
		}(i)
	}
	wg.Wait()

	if len(store.GetAllExpressions()) != 100 {
		t.Errorf("Expected 100 expressions, got %d", len(store.GetAllExpressions()))
	}
}

func TestTaskDependencies(t *testing.T) {
	store := NewStore()

	task1 := &Task{ID: "1", Status: "pending"}
	task2 := &Task{ID: "2", Status: "pending", DependsOn: []string{"1"}}

	store.AddTask(task1)
	store.AddTask(task2)

	if store.IsTaskReady(task2) {
		t.Error("Task2 should not be ready")
	}

	task1.Status = "completed"
	store.UpdateTask(task1)

	if !store.IsTaskReady(task2) {
		t.Error("Task2 should be ready")
	}
}

func (s *Store) isTaskReady(task *Task) bool {
	for _, depID := range task.DependsOn {
		depTask, exists := s.GetTask(depID)
		if !exists || depTask.Status != "completed" {
			return false
		}
	}
	return true
}
