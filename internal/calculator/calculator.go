package calculator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Expression struct {
	value    string
	position int
}

func Calc(expression string) (float64, error) {
	return calculate(expression)
}

func calculate(exprString string) (float64, error) {
	expression := &Expression{
		value: strings.ReplaceAll(exprString, " ", ""),
	}
	result, err := expression.parse()
	if err != nil {
		return 0, err
	}
	if expression.position < len(expression.value) {
		return 0, fmt.Errorf("неожиданный символ в позиции %d", expression.position)
	}
	return result, nil
}

func (expression *Expression) parse() (float64, error) {
	result, err := expression.parseOperation()
	if err != nil {
		return 0, err
	}
	for expression.position < len(expression.value) {
		char := expression.value[expression.position]
		if char == '+' {
			expression.position++
			value, err := expression.parseOperation()
			if err != nil {
				return 0, err
			}
			result += value
		} else if char == '-' {
			expression.position++
			value, err := expression.parseOperation()
			if err != nil {
				return 0, err
			}
			result -= value
		} else {
			break
		}
	}
	return result, nil
}

func (expression *Expression) parseOperation() (float64, error) {
	result, err := expression.parseBrackets()
	if err != nil {
		return 0, err
	}
	for expression.position < len(expression.value) {
		ch := expression.value[expression.position]
		if ch == '*' {
			expression.position++
			value, err := expression.parseBrackets()
			if err != nil {
				return 0, err
			}
			result *= value
		} else if ch == '/' {
			expression.position++
			value, err := expression.parseBrackets()
			if err != nil {
				return 0, err
			}
			if value == 0 {
				return 0, errors.New("деление на ноль")
			}
			result /= value
		} else {
			break
		}
	}
	return result, nil
}

func (expression *Expression) parseBrackets() (float64, error) {
	if expression.position >= len(expression.value) {
		return 0, fmt.Errorf("неожиданный конец выражения")
	}

	if expression.value[expression.position] == '+' {
		expression.position++
		return expression.parseBrackets()
	} else if expression.value[expression.position] == '-' {
		expression.position++
		value, err := expression.parseBrackets()
		if err != nil {
			return 0, err
		}
		return -value, nil
	}

	char := expression.value[expression.position]
	if char == '(' {
		expression.position++
		result, err := expression.parse()
		if err != nil {
			return 0, err
		}
		if expression.position >= len(expression.value) || expression.value[expression.position] != ')' {
			return 0, fmt.Errorf("отсутствует закрывающая скобка на позиции %d", expression.position)
		}
		expression.position++
		return result, nil
	} else if unicode.IsDigit(rune(char)) || char == '.' {
		start := expression.position
		for expression.position < len(expression.value) &&
			(unicode.IsDigit(rune(expression.value[expression.position])) || expression.value[expression.position] == '.') {
			expression.position++
		}
		numberStr := expression.value[start:expression.position]
		number, err := strconv.ParseFloat(numberStr, 64)
		if err != nil {
			return 0, fmt.Errorf("некорректное число '%s'", numberStr)
		}
		return number, nil
	} else {
		return 0, fmt.Errorf("неожиданный символ '%c' в позиции %d", char, expression.position)
	}
}
