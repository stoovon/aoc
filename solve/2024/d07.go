package solve2024

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day7 struct {
}

func (d Day7) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2024, Day: 7}
}

func (d Day7) value(target int, numbers []int, allowConcat bool) int {
	if len(numbers) == 1 {
		if numbers[0] == target {
			return target
		}
		return 0
	}

	// Try addition
	if result := d.value(target, append([]int{numbers[0] + numbers[1]}, numbers[2:]...), allowConcat); result != 0 {
		return result
	}

	// Try multiplication
	if result := d.value(target, append([]int{numbers[0] * numbers[1]}, numbers[2:]...), allowConcat); result != 0 {
		return result
	}

	// Try concatenation if allowed
	if allowConcat {
		concatValue, _ := strconv.Atoi(fmt.Sprintf("%d%d", numbers[0], numbers[1]))
		if result := d.value(target, append([]int{concatValue}, numbers[2:]...), allowConcat); result != 0 {
			return result
		}
	}

	return 0
}

func (d Day7) processInput(input string, processFunc func(test int, numbers []int) int) string {
	total := 0

	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		parts := strings.Split(line, ": ")
		target, _ := strconv.Atoi(parts[0])

		var numbers []int
		if err := json.Unmarshal([]byte("["+strings.ReplaceAll(parts[1], " ", ",")+"]"), &numbers); err != nil {
			return fmt.Sprintf("%d: %s", target, err.Error())
		}

		total += processFunc(target, numbers)
	}

	return strconv.Itoa(total)
}

func (d Day7) Part1(input string) (string, error) {
	result := d.processInput(input, func(test int, numbers []int) int {
		return d.value(test, numbers, false)
	})

	return result, nil
}

func (d Day7) Part2(input string) (string, error) {
	result := d.processInput(input, func(test int, numbers []int) int {
		return d.value(test, numbers, true)
	})

	return result, nil
}

func init() {
	solve.Register(Day7{})
}
