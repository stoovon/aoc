package solve2024

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day7 struct{}

func (d Day7) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2024, Day: 7}
}

type calibrationEquation struct {
	target  uint64
	numbers []uint64
}

func parseCalibrationNumber(text string) (uint64, error) {
	if text == "" {
		return 0, fmt.Errorf("expected a nonnegative decimal integer")
	}
	for _, digit := range text {
		if digit < '0' || digit > '9' {
			return 0, fmt.Errorf("expected a nonnegative decimal integer: %q", text)
		}
	}
	value, err := strconv.ParseUint(text, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("calibration number exceeds uint64: %q", text)
	}
	return value, nil
}

func parseCalibrationEquation(line string) (calibrationEquation, error) {
	targetText, numbersText, ok := strings.Cut(line, ":")
	if !ok {
		return calibrationEquation{}, fmt.Errorf("expected target: operands")
	}
	target, err := parseCalibrationNumber(strings.TrimSpace(targetText))
	if err != nil {
		return calibrationEquation{}, err
	}
	fields := strings.Fields(numbersText)
	if len(fields) == 0 {
		return calibrationEquation{}, fmt.Errorf("expected at least one operand")
	}
	numbers := make([]uint64, len(fields))
	for i, field := range fields {
		numbers[i], err = parseCalibrationNumber(field)
		if err != nil {
			return calibrationEquation{}, err
		}
	}
	return calibrationEquation{target: target, numbers: numbers}, nil
}

func calibrationScale(value uint64) uint64 {
	scale := uint64(10)
	for value >= 10 {
		if scale > math.MaxUint64/10 {
			return 0
		}
		scale *= 10
		value /= 10
	}
	return scale
}

func calibrationAddBound(a, b uint64) uint64 {
	if a > math.MaxUint64-b {
		return math.MaxUint64
	}
	return a + b
}

func calibrationMultiplyBound(a, b uint64) uint64 {
	if b != 0 && a > math.MaxUint64/b {
		return math.MaxUint64
	}
	return a * b
}

func calibrationMatches(equation calibrationEquation, allowConcat bool) bool {
	var scales []uint64
	if allowConcat {
		scales = make([]uint64, len(equation.numbers))
		for i, value := range equation.numbers {
			scales[i] = calibrationScale(value)
		}
	}
	maximum := make([]uint64, len(equation.numbers))
	maximum[0] = equation.numbers[0]
	for index := 1; index < len(maximum); index++ {
		previous, last := maximum[index-1], equation.numbers[index]
		maximum[index] = max(calibrationAddBound(previous, last), calibrationMultiplyBound(previous, last))
		if allowConcat {
			var concatenated uint64
			if previous == 0 {
				concatenated = last
			} else if scales[index] == 0 {
				concatenated = math.MaxUint64
			} else {
				concatenated = calibrationAddBound(calibrationMultiplyBound(previous, scales[index]), last)
			}
			maximum[index] = max(maximum[index], concatenated)
		}
	}
	var matches func(int, uint64) bool
	matches = func(index int, target uint64) bool {
		if target > maximum[index] {
			return false
		}
		if index == 0 {
			return target == equation.numbers[0]
		}
		last := equation.numbers[index]
		if last == 0 {
			if target == 0 {
				return true
			}
		} else if target%last == 0 && matches(index-1, target/last) {
			return true
		}
		if allowConcat {
			scale := scales[index]
			if scale != 0 {
				if target%scale == last && matches(index-1, target/scale) {
					return true
				}
			} else if target == last && matches(index-1, 0) {
				return true
			}
		}
		return target >= last && matches(index-1, target-last)
	}
	return matches(len(equation.numbers)-1, equation.target)
}

func solveCalibration(input string, allowConcat bool) (string, error) {
	input = strings.TrimSuffix(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	if input == "" {
		return "", fmt.Errorf("expected calibration equations")
	}
	var total uint64
	for index, line := range strings.Split(input, "\n") {
		equation, err := parseCalibrationEquation(line)
		if err != nil {
			return "", fmt.Errorf("line %d: %w", index+1, err)
		}
		if calibrationMatches(equation, allowConcat) {
			if equation.target > math.MaxUint64-total {
				return "", fmt.Errorf("line %d: calibration total exceeds uint64", index+1)
			}
			total += equation.target
		}
	}
	return strconv.FormatUint(total, 10), nil
}

func (d Day7) Part1(input string) (string, error) { return solveCalibration(input, false) }
func (d Day7) Part2(input string) (string, error) { return solveCalibration(input, true) }

func init() { solve.Register(Day7{}) }
