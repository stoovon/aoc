package solve2020

import (
	"aoc/solve"
	"strconv"
	"strings"
)

type Day18 struct{}

func (d Day18) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2020, Day: 18}
}

// Evaluate an expression left-to-right, with parentheses
func evalHomework(input string, evalFunc func([]string, int) (int64, int, error)) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var sum int64
	for _, line := range lines {
		tokens := tokenize(line)
		val, _, err := evalFunc(tokens, 0)
		if err != nil {
			return "", err
		}
		sum += val
	}
	return strconv.FormatInt(sum, 10), nil
}

func tokenize(expr string) []string {
	var tokens []string
	i := 0
	for i < len(expr) {
		switch expr[i] {
		case ' ':
			i++
		case '(', ')', '+', '*':
			tokens = append(tokens, string(expr[i]))
			i++
		default:
			j := i
			for j < len(expr) && expr[j] >= '0' && expr[j] <= '9' {
				j++
			}
			tokens = append(tokens, expr[i:j])
			i = j
		}
	}
	return tokens
}

// Evaluate tokens left-to-right, handling parentheses
// Evaluate tokens left-to-right, handling parentheses (Part 1)
func evalTokensLeftToRight(tokens []string, pos int) (int64, int, error) {
	var acc int64
	var op string
	for pos < len(tokens) {
		tok := tokens[pos]
		switch tok {
		case "+", "*":
			op = tok
			pos++
		case "(":
			val, newPos, err := evalTokensLeftToRight(tokens, pos+1)
			if err != nil {
				return 0, 0, err
			}
			if op == "" {
				acc = val
			} else if op == "+" {
				acc += val
			} else if op == "*" {
				acc *= val
			}
			pos = newPos
		case ")":
			return acc, pos + 1, nil
		default: // number
			val, err := strconv.ParseInt(tok, 10, 64)
			if err != nil {
				return 0, 0, err
			}
			if op == "" {
				acc = val
			} else if op == "+" {
				acc += val
			} else if op == "*" {
				acc *= val
			}
			pos++
		}
	}
	return acc, pos, nil
}

// Evaluate tokens with addition precedence over multiplication (Part 2)
func evalTokensAddPrecedence(tokens []string, pos int) (int64, int, error) {
	var values []int64
	var ops []string
	var applyAdd = func() {
		for i := 0; i < len(ops); {
			if ops[i] == "+" {
				values[i] += values[i+1]
				values = append(values[:i+1], values[i+2:]...)
				ops = append(ops[:i], ops[i+1:]...)
			} else {
				i++
			}
		}
	}
	for pos < len(tokens) {
		tok := tokens[pos]
		switch tok {
		case "+", "*":
			ops = append(ops, tok)
			pos++
		case "(":
			val, newPos, err := evalTokensAddPrecedence(tokens, pos+1)
			if err != nil {
				return 0, 0, err
			}
			values = append(values, val)
			pos = newPos
		case ")":
			applyAdd()
			res := values[0]
			for i, op := range ops {
				if op == "*" {
					res *= values[i+1]
				}
			}
			return res, pos + 1, nil
		default: // number
			val, err := strconv.ParseInt(tok, 10, 64)
			if err != nil {
				return 0, 0, err
			}
			values = append(values, val)
			pos++
		}
	}
	applyAdd()
	res := values[0]
	for i, op := range ops {
		if op == "*" {
			res *= values[i+1]
		}
	}
	return res, pos, nil
}

func (d Day18) Part1(input string) (string, error) {
	return evalHomework(input, evalTokensLeftToRight)
}

func (d Day18) Part2(input string) (string, error) {
	return evalHomework(input, evalTokensAddPrecedence)
}

func init() {
	solve.Register(Day18{})
}
