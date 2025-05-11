package orchestrator

import (
	"testing"

	"github.com/ladnaaaaaa/calc_service/internal/models"
)

func TestParseExpression(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected float64
	}{
		{"simple_addition", "2+3", 5},
		{"multiplication_priority", "2+3*4", 14},
		{"with_parentheses", "(2+3)*4", 20},
		{"division", "10/2", 5},
		{"complex_expression", "(10-4)*3/(2+1)", 6},
		{"invalid_characters", "2+a", 0},
		{"empty_expression", "", 0},
		{"invalid_parentheses", "(2+3", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, err := ParseExpression(tt.expr)
			if tt.name == "invalid_characters" || tt.name == "empty_expression" || tt.name == "invalid_parentheses" {
				if err == nil {
					t.Errorf("expected error for invalid input, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Эмулируем выполнение задач по принципу постфиксного вычисления
			stack := make([]float64, 0)
			for _, task := range tasks {
				stack = append(stack, task.Arg1)
				stack = append(stack, task.Arg2)
				arg2 := stack[len(stack)-1]
				arg1 := stack[len(stack)-2]
				stack = stack[:len(stack)-2]
				var res float64
				switch task.Operation {
				case models.OperationAdd:
					res = arg1 + arg2
				case models.OperationSubtract:
					res = arg1 - arg2
				case models.OperationMultiply:
					res = arg1 * arg2
				case models.OperationDivide:
					res = arg1 / arg2
				}
				stack = append(stack, res)
			}
			if len(stack) > 0 && stack[len(stack)-1] != tt.expected {
				t.Errorf("expected result %v, got %v", tt.expected, stack[len(stack)-1])
			}
		})
	}
}
