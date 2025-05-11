package orchestrator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ladnaaaaaa/calc_service/internal/models"
)

type Token struct {
	Type  string
	Value string
	Num   float64
}

func ParseExpression(expr string) ([]*models.Task, error) {
	tokens, err := tokenize(expr)
	if err != nil {
		return nil, fmt.Errorf("tokenize error: %v", err)
	}

	postfix, err := shuntingYard(tokens)
	if err != nil {
		return nil, fmt.Errorf("shuntingYard error: %v", err)
	}

	tasks, err := buildTasks(postfix)
	if err != nil {
		return nil, fmt.Errorf("buildTasks error: %v", err)
	}

	return tasks, nil
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

func buildTasks(postfix []Token) ([]*models.Task, error) {
	var stack []float64
	tasks := make([]*models.Task, 0)
	orderNum := 0

	for _, token := range postfix {
		if token.Type == "num" {
			stack = append(stack, token.Num)
			continue
		}

		if len(stack) < 2 {
			return nil, fmt.Errorf("invalid expression")
		}

		arg2 := stack[len(stack)-1]
		arg1 := stack[len(stack)-2]
		stack = stack[:len(stack)-2]

		task := &models.Task{
			Arg1:      arg1,
			Arg2:      arg2,
			Operation: models.Operation(token.Value),
			Status:    models.StatusPending,
			OrderNum:  orderNum,
		}
		orderNum++

		tasks = append(tasks, task)

		var res float64
		switch token.Value {
		case "+":
			res = arg1 + arg2
		case "-":
			res = arg1 - arg2
		case "*":
			res = arg1 * arg2
		case "/":
			res = arg1 / arg2
		}
		stack = append(stack, res)
	}

	if len(stack) != 1 {
		return nil, fmt.Errorf("invalid expression structure")
	}

	return tasks, nil
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
