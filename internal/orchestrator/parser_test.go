package orchestrator

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
		{
			name:    "invalid characters",
			expr:    "2+a",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, finalID, err := parseExpression(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseExpression() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			store := NewStore()
			expr := &Expression{
				ID:          "test",
				Tasks:       tasks,
				FinalTaskID: finalID,
				Status:      "processing",
			}

			for _, task := range tasks {
				store.AddTask(task)
			}

			for {
				var readyTask *Task
				for _, task := range store.GetPendingTasks() {
					if store.IsTaskReady(task) {
						readyTask = task
						break
					}
				}
				if readyTask == nil {
					break
				}

				arg1, _ := store.GetTask(readyTask.Arg1ID)
				arg2, _ := store.GetTask(readyTask.Arg2ID)

				switch readyTask.Operation {
				case "+":
					readyTask.Result = arg1.Result + arg2.Result
				case "-":
					readyTask.Result = arg1.Result - arg2.Result
				case "*":
					readyTask.Result = arg1.Result * arg2.Result
				case "/":
					readyTask.Result = arg1.Result / arg2.Result
				}
				readyTask.Status = "completed"
				store.UpdateTask(readyTask)
			}

			finalTask, _ := store.GetTask(finalID)
			expr.Result = finalTask.Result
			expr.Status = "completed"
			store.AddExpression(expr)

			updatedExpr, _ := store.GetExpression("test")
			assert.Equal(t, tt.want, updatedExpr.Result)
			assert.Equal(t, "completed", updatedExpr.Status)
		})
	}
}
