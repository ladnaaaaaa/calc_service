package orchestrator

import (
	"testing"

	"github.com/ladnaaaaaa/calc_service/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestParseExpression(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		want    float64
		wantErr bool
	}{
		{
			name:    "simple addition",
			expr:    "2+3",
			want:    5,
			wantErr: false,
		},
		{
			name:    "multiplication priority",
			expr:    "2+3*4",
			want:    14,
			wantErr: false,
		},
		{
			name:    "with parentheses",
			expr:    "(2+3)*4",
			want:    20,
			wantErr: false,
		},
		{
			name:    "division",
			expr:    "10/2",
			want:    5,
			wantErr: false,
		},
		{
			name:    "complex expression",
			expr:    "(10-4)*3/(2+1)",
			want:    6,
			wantErr: false,
		},
		{
			name:    "invalid characters",
			expr:    "2+a",
			wantErr: true,
		},
		{
			name:    "empty expression",
			expr:    "",
			wantErr: true,
		},
		{
			name:    "invalid parentheses",
			expr:    "2+(3*4",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, err := parseExpression(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseExpression() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			store := NewStore()
			expr := &models.Expression{
				Expression: tt.expr,
				Status:     models.StatusPending,
				UserID:     1,
			}

			if err := store.AddExpression(expr); err != nil {
				t.Fatalf("failed to add expression: %v", err)
			}

			for i, task := range tasks {
				task.ExpressionID = expr.ID
				task.OrderNum = i
				if err := store.AddTask(task); err != nil {
					t.Fatalf("failed to add task: %v", err)
				}
			}

			// Process tasks in order
			for i := 0; i < len(tasks); i++ {
				tasks, err := store.GetTasksByExpressionID(expr.ID)
				if err != nil {
					t.Fatalf("failed to get tasks: %v", err)
				}

				var currentTask *models.Task
				for j := range tasks {
					if tasks[j].OrderNum == i {
						currentTask = &tasks[j]
						break
					}
				}

				if currentTask == nil {
					t.Fatalf("task not found")
				}

				switch currentTask.Operation {
				case models.OperationAdd:
					currentTask.Result = currentTask.Arg1 + currentTask.Arg2
				case models.OperationSubtract:
					currentTask.Result = currentTask.Arg1 - currentTask.Arg2
				case models.OperationMultiply:
					currentTask.Result = currentTask.Arg1 * currentTask.Arg2
				case models.OperationDivide:
					currentTask.Result = currentTask.Arg1 / currentTask.Arg2
				}
				currentTask.Status = models.StatusCompleted

				if err := store.UpdateTask(currentTask); err != nil {
					t.Fatalf("failed to update task: %v", err)
				}
			}

			// Get the final result
			dbTasks, err := store.GetTasksByExpressionID(expr.ID)
			if err != nil {
				t.Fatalf("failed to get tasks: %v", err)
			}

			// Find the last task
			var lastTask models.Task
			for _, task := range dbTasks {
				if task.OrderNum == len(dbTasks)-1 {
					lastTask = task
					break
				}
			}

			assert.Equal(t, tt.want, lastTask.Result)
			assert.Equal(t, models.StatusCompleted, lastTask.Status)
		})
	}
}
