package orchestrator

import (
	"fmt"
	"strconv"
	"strings"
)

type Token struct {
	Type  string
	Value string
	Num   float64
}

func parseExpression(expr string) ([]*Task, string, error) {
	tokens, err := tokenize(expr)
	if err != nil {
		return nil, "", fmt.Errorf("tokenize error: %v", err)
	}

	postfix, err := shuntingYard(tokens)
	if err != nil {
		return nil, "", fmt.Errorf("shuntingYard error: %v", err)
	}

	tasks, finalID, err := buildTasks(postfix)
	if err != nil {
		return nil, "", fmt.Errorf("buildTasks error: %v", err)
	}

	return tasks, finalID, nil
}

func tokenize(expr string) ([]Token, error) {
	var tokens []Token
	expr = strings.ReplaceAll(expr, " ", "")

	for i := 0; i < len(expr); {
		c := expr[i]
		switch {
		case c >= '0' && c <= '9' || c == '.':
			j := i
			for ; j < len(expr); j++ {
				if !(expr[j] >= '0' && expr[j] <= '9' || expr[j] == '.') {
					break
				}
			}
			num, err := strconv.ParseFloat(expr[i:j], 64)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, Token{Type: "num", Value: expr[i:j], Num: num})
			i = j
		case c == '+' || c == '-' || c == '*' || c == '/':
			tokens = append(tokens, Token{Type: "op", Value: string(c)})
			i++
		case c == '(' || c == ')':
			tokens = append(tokens, Token{Type: "paren", Value: string(c)})
			i++
		default:
			return nil, fmt.Errorf("invalid character: %c", c)
		}
	}
	return tokens, nil
}

func shuntingYard(tokens []Token) ([]Token, error) {
	var output []Token
	var stack []Token

	for _, token := range tokens {
		switch token.Type {
		case "num":
			output = append(output, token)
		case "op":
			for len(stack) > 0 && stack[len(stack)-1].Type == "op" && precedence(token.Value) <= precedence(stack[len(stack)-1].Value) {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		case "paren":
			if token.Value == "(" {
				stack = append(stack, token)
			} else {
				for len(stack) > 0 && stack[len(stack)-1].Value != "(" {
					output = append(output, stack[len(stack)-1])
					stack = stack[:len(stack)-1]
				}
				if len(stack) == 0 {
					return nil, fmt.Errorf("mismatched parentheses")
				}
				stack = stack[:len(stack)-1]
			}
		}
	}

	for len(stack) > 0 {
		if stack[len(stack)-1].Value == "(" {
			return nil, fmt.Errorf("mismatched parentheses")
		}
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return output, nil
}

// internal/orchestrator/parser.go
func buildTasks(postfix []Token) ([]*Task, string, error) {
	var stack []string
	tasks := make([]*Task, 0)
	values := make(map[string]float64)

	for _, token := range postfix {
		if token.Type == "num" {
			// Создаем завершенную задачу для числа
			id := fmt.Sprintf("num_%d", len(tasks)+1)
			task := &Task{
				ID:     id,
				Result: token.Num,
				Status: "completed",
			}
			stack = append(stack, id)
			values[id] = token.Num
			tasks = append(tasks, task)
			continue
		}

		if len(stack) < 2 {
			return nil, "", fmt.Errorf("invalid expression")
		}

		b := stack[len(stack)-1]
		a := stack[len(stack)-2]
		stack = stack[:len(stack)-2]

		taskID := fmt.Sprintf("op_%d", len(tasks)+1)
		task := &Task{
			ID:        taskID,
			Arg1ID:    a,
			Arg2ID:    b,
			Operation: token.Value,
			Status:    "pending",
			DependsOn: []string{a, b},
		}

		stack = append(stack, taskID)
		tasks = append(tasks, task)
	}

	if len(stack) != 1 {
		return nil, "", fmt.Errorf("invalid expression structure")
	}

	return tasks, stack[0], nil
}

func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	}
	return 0
}
