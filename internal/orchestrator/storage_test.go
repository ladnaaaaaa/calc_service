package orchestrator

import (
	"testing"

	"github.com/ladnaaaaaa/calc_service/internal/database"
	"github.com/ladnaaaaaa/calc_service/internal/models"
)

func TestStoreConcurrency(t *testing.T) {
	database.InitTestDB(t)
	database.ClearDB()
	store := NewStore()
	userID := uint(1)

	for i := 0; i < 100; i++ {
		expr := &models.Expression{
			Expression: "expr",
			Status:     models.StatusPending,
			UserID:     userID,
		}
		err := store.AddExpression(expr)
		if err != nil {
			t.Errorf("failed to add expression: %v", err)
		}
	}

	exprs, err := store.GetAllExpressions(userID)
	if err != nil {
		t.Fatalf("failed to get expressions: %v", err)
	}
	if len(exprs) != 100 {
		t.Errorf("Expected 100 expressions, got %d", len(exprs))
	}
}

func TestTaskDependencies(t *testing.T) {
	database.InitTestDB(t)
	database.ClearDB()
	store := NewStore()
	userID := uint(1)
	expr := &models.Expression{
		Expression: "expr",
		Status:     models.StatusPending,
		UserID:     userID,
	}
	err := store.AddExpression(expr)
	if err != nil {
		t.Fatalf("failed to add expression: %v", err)
	}
	task1 := &models.Task{
		ExpressionID: expr.ID,
		Arg1:         1,
		Arg2:         2,
		Operation:    models.OperationAdd,
		Status:       models.StatusPending,
		OrderNum:     0,
	}
	task2 := &models.Task{
		ExpressionID: expr.ID,
		Arg1:         3,
		Arg2:         4,
		Operation:    models.OperationMultiply,
		Status:       models.StatusPending,
		OrderNum:     1,
	}
	err = store.AddTask(task1)
	if err != nil {
		t.Fatalf("failed to add task1: %v", err)
	}
	err = store.AddTask(task2)
	if err != nil {
		t.Fatalf("failed to add task2: %v", err)
	}
	tasks, err := store.GetTasksByExpressionID(expr.ID)
	if err != nil {
		t.Fatalf("failed to get tasks: %v", err)
	}
	if !store.IsTaskReady(&tasks[0]) {
		t.Errorf("task1 should be ready")
	}
	if store.IsTaskReady(&tasks[1]) {
		t.Errorf("task2 should not be ready before task1 is completed")
	}
	tasks[0].Status = models.StatusCompleted
	err = store.UpdateTask(&tasks[0])
	if err != nil {
		t.Fatalf("failed to update task1: %v", err)
	}
	if !store.IsTaskReady(&tasks[1]) {
		t.Errorf("task2 should be ready after task1 is completed")
	}
}
