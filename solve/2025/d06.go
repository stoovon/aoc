package solve2025

import (
	"errors"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day6 struct {
}

func (d Day6) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2025, Day: 6}
}

func (d Day6) Part1(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")

	if len(lines) == 0 {
		return "", errors.New("empty input")
	}

	columns := parseColumns(lines)

	if len(columns) == 0 {
		return "", errors.New("no columns found")
	}

	grandTotal := int64(0)
	for _, col := range columns {
		result, err := evaluateColumn(col)
		if err != nil {
			return "", err
		}
		grandTotal += result
	}

	return strconv.FormatInt(grandTotal, 10), nil
}

func parseColumns(lines []string) [][]string {
	if len(lines) == 0 {
		return nil
	}

	var columns [][]string

	for _, line := range lines {
		fields := strings.Fields(line)

		if len(columns) == 0 {
			columns = make([][]string, len(fields))
		}

		if len(fields) != len(columns) {
		}

		for i, field := range fields {
			if i < len(columns) {
				columns[i] = append(columns[i], field)
			}
		}
	}

	return columns
}

// The column contains numbers and ends with an operator
func evaluateColumn(col []string) (int64, error) {
	if len(col) == 0 {
		return 0, errors.New("empty column")
	}

	operator := strings.TrimSpace(col[len(col)-1])
	if operator != "+" && operator != "*" {
		return 0, errors.New("invalid operator: " + operator)
	}

	var numbers []int64
	for i := 0; i < len(col)-1; i++ {
		val, err := strconv.ParseInt(strings.TrimSpace(col[i]), 10, 64)
		if err != nil {
			return 0, errors.New("failed to parse number: " + col[i])
		}
		numbers = append(numbers, val)
	}

	if len(numbers) == 0 {
		return 0, errors.New("no numbers in column")
	}

	if operator == "+" {
		result := int64(0)
		for _, num := range numbers {
			result += num
		}
		return result, nil
	} else { // operator == "*"
		result := int64(1)
		for _, num := range numbers {
			result *= num
		}
		return result, nil
	}
}

func (d Day6) Part2(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")

	if len(lines) == 0 {
		return "", errors.New("empty input")
	}

	// Columns are separated by spaces
	columns := parseFixedWidthColumns(lines)

	if len(columns) == 0 {
		return "", errors.New("no columns found")
	}

	var grandTotal int64

	for _, col := range columns {
		result, err := evaluateColumnAsProblems(col)
		if err != nil {
			return "", err
		}
		grandTotal += result
	}

	return strconv.FormatInt(grandTotal, 10), nil
}

// Columns are identified by whitespace separation across ALL lines
func parseFixedWidthColumns(lines []string) [][]string {
	if len(lines) == 0 {
		return nil
	}

	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	isContentAtPos := make([]bool, maxLen)
	for _, line := range lines {
		for i := 0; i < len(line); i++ {
			if line[i] != ' ' {
				isContentAtPos[i] = true
			}
		}
	}

	var columnStarts []int
	for i := 0; i < maxLen; i++ {
		if isContentAtPos[i] && (i == 0 || !isContentAtPos[i-1]) {
			columnStarts = append(columnStarts, i)
		}
	}

	columnStarts = append(columnStarts, maxLen)

	var columns [][]string
	for i := 0; i < len(columnStarts)-1; i++ {
		start := columnStarts[i]
		end := columnStarts[i+1]

		var col []string
		for _, line := range lines {
			var colVal string
			if start < len(line) {
				if end <= len(line) {
					colVal = line[start:end]
				} else {
					colVal = line[start:] + strings.Repeat(" ", end-len(line))
				}
			} else {
				colVal = strings.Repeat(" ", end-start)
			}
			col = append(col, colVal)
		}
		columns = append(columns, col)
	}

	return columns
}

// evaluateColumnAsProblems reads a column right-to-left, top-to-bottom
// to extract numbers and an operator, then evaluates the result
func evaluateColumnAsProblems(col []string) (int64, error) {
	if len(col) < 2 {
		return 0, errors.New("column too short")
	}

	operatorStr := strings.TrimSpace(col[len(col)-1])

	if operatorStr == "" {
		return 0, nil
	}

	if !isOperator(operatorStr) {
		return 0, errors.New("expected operator in last row, got: " + operatorStr)
	}

	width := len(col[0])
	if width == 0 {
		return 0, errors.New("empty column")
	}

	// Each position (right to left) represents a digit position in a number
	var numbers []int64
	for charPos := width - 1; charPos >= 0; charPos-- {
		var digitStr string
		for rowIdx := 0; rowIdx < len(col)-1; rowIdx++ {
			ch := byte(' ')
			if charPos < len(col[rowIdx]) {
				ch = col[rowIdx][charPos]
			}
			if ch != ' ' {
				digitStr += string(ch)
			}
		}

		if digitStr != "" {
			num, err := strconv.ParseInt(digitStr, 10, 64)
			if err != nil {
				return 0, errors.New("invalid number: " + digitStr)
			}
			numbers = append(numbers, num)
		}
	}

	if len(numbers) == 0 {
		return 0, errors.New("no numbers found in column")
	}

	result := evaluateNumbers(numbers, operatorStr)
	return result, nil
}

func isOperator(s string) bool {
	return s == "+" || s == "*"
}

func evaluateNumbers(numbers []int64, operator string) int64 {
	if len(numbers) == 0 {
		return 0
	}

	switch operator {
	case "+":
		result := int64(0)
		for _, num := range numbers {
			result += num
		}
		return result
	case "*":
		result := int64(1)
		for _, num := range numbers {
			result *= num
		}
		return result
	}

	return 0
}

func init() {
	solve.Register(Day6{})
}
