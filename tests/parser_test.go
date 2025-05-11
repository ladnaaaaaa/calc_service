package tests

import (
	"testing"

	"github.com/ladnaaaaaa/calc_service/internal/models"
	"github.com/ladnaaaaaa/calc_service/internal/orchestrator"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		expectError bool
		expected    []*models.Task
	}{
		{
			name:        "Simple addition",
			expression:  "2+3",
			expectError: false,
			expected: []*models.Task{
				{
					Arg1:      2,
					Arg2:      3,
					Operation: models.OperationAdd,
					Status:    models.StatusPending,
					OrderNum:  1,
				},
			},
		},
		{
			name:        "Complex expression",
			expression:  "2+3*4",
			expectError: false,
			expected: []*models.Task{
				{
					Arg1:      3,
					Arg2:      4,
					Operation: models.OperationMultiply,
					Status:    models.StatusPending,
					OrderNum:  1,
				},
				{
					Arg1:      2,
					Arg2:      12,
					Operation: models.OperationAdd,
					Status:    models.StatusPending,
					OrderNum:  2,
				},
			},
		},
		{
			name:        "Expression with parentheses",
			expression:  "(2+3)*4",
			expectError: false,
			expected: []*models.Task{
				{
					Arg1:      2,
					Arg2:      3,
					Operation: models.OperationAdd,
					Status:    models.StatusPending,
					OrderNum:  1,
				},
				{
					Arg1:      5,
					Arg2:      4,
					Operation: models.OperationMultiply,
					Status:    models.StatusPending,
					OrderNum:  2,
				},
			},
		},
		{
			name:        "Invalid expression",
			expression:  "2++3",
			expectError: true,
		},
		{
			name:        "Invalid characters",
			expression:  "2+a",
			expectError: true,
		},
		{
			name:        "Unbalanced parentheses",
			expression:  "(2+3",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, err := orchestrator.ParseExpression(tt.expression)
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, len(tt.expected), len(tasks))

			for i, task := range tasks {
				expected := tt.expected[i]
				assert.Equal(t, expected.Arg1, task.Arg1)
				assert.Equal(t, expected.Arg2, task.Arg2)
				assert.Equal(t, expected.Operation, task.Operation)
				assert.Equal(t, expected.Status, task.Status)
				assert.Equal(t, expected.OrderNum, task.OrderNum)
			}
		})
	}
}

func TestTaskOrder(t *testing.T) {
	expression := "2+3*4-5"
	tasks, err := orchestrator.ParseExpression(expression)
	assert.NoError(t, err)

	// Verify task order
	assert.Equal(t, 3, len(tasks))
	assert.Equal(t, models.OperationMultiply, tasks[0].Operation) // 3*4
	assert.Equal(t, models.OperationAdd, tasks[1].Operation)      // 2+12
	assert.Equal(t, models.OperationSubtract, tasks[2].Operation) // 14-5
}

func TestTaskDependencies(t *testing.T) {
	expression := "(2+3)*(4-5)"
	tasks, err := orchestrator.ParseExpression(expression)
	assert.NoError(t, err)

	// Verify task dependencies
	assert.Equal(t, 3, len(tasks))
	assert.Equal(t, models.OperationAdd, tasks[0].Operation)      // 2+3
	assert.Equal(t, models.OperationSubtract, tasks[1].Operation) // 4-5
	assert.Equal(t, models.OperationMultiply, tasks[2].Operation) // 5*(-1)
}
