package solve2022

import (
	"aoc/solve"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Day11 struct{}

func (d Day11) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2022, Day: 11}
}

func (d Day11) Part1(input string) (string, error) {
	// Parse the input to extract monkey data
	monkeys, err := parseMonkeys(input)
	if err != nil {
		return "", err
	}

	// Simulate 20 rounds with division by 3
	inspected := inspectMonkeys(monkeys, 20, true)

	// Calculate the level of monkey business
	level := monkeyBusiness(inspected)
	return strconv.Itoa(level), nil
}

func (d Day11) Part2(input string) (string, error) {
	// Parse the input to extract monkey data
	monkeys, err := parseMonkeys(input)
	if err != nil {
		return "", err
	}

	// Simulate 10000 rounds without division by 3
	inspected := inspectMonkeys(monkeys, 10000, false)

	// Calculate the level of monkey business
	level := monkeyBusiness(inspected)
	return strconv.Itoa(level), nil
}

type Monkey struct {
	N       int
	Items   []int
	Op      string
	Arg     int
	Test    int
	IfTrue  int
	IfFalse int
}

func parseMonkeys(input string) ([]Monkey, error) {
	lines := strings.Split(input, "\n\n")
	monkeys := []Monkey{}

	for _, block := range lines {
		var n, test, t, f int
		var itemsStr, op, argStr string
		// Extract the 'Starting items' line manually
		lines := strings.Split(block, "\n")
		if len(lines) < 6 {
			return nil, fmt.Errorf("monkey block has insufficient lines: %v", block)
		}
		itemsStr = strings.TrimPrefix(lines[1], "  Starting items: ")
		fmt.Sscanf(lines[0], "Monkey %d:", &n)
		fmt.Sscanf(lines[2], "  Operation: new = old %s %s", &op, &argStr)
		fmt.Sscanf(lines[3], "  Test: divisible by %d", &test)
		fmt.Sscanf(lines[4], "    If true: throw to monkey %d", &t)
		fmt.Sscanf(lines[5], "    If false: throw to monkey %d", &f)

		if test == 0 {
			return nil, fmt.Errorf("monkey %d has a test divisor of 0, which is invalid", n)
		}

		items := []int{}
		for _, item := range strings.Split(itemsStr, ", ") {
			var val int
			fmt.Sscanf(item, "%d", &val)
			items = append(items, val)
		}

		arg := 0
		if argStr != "old" {
			fmt.Sscanf(argStr, "%d", &arg)
		}

		monkeys = append(monkeys, Monkey{N: n, Items: items, Op: op, Arg: arg, Test: test, IfTrue: t, IfFalse: f})
	}

	return monkeys, nil
}

func inspectMonkeys(monkeys []Monkey, rounds int, divideByThree bool) map[int]int {
	inspected := map[int]int{}
	items := map[int][]int{}
	for _, monkey := range monkeys {
		items[monkey.N] = append([]int{}, monkey.Items...)
	}

	// Calculate the least common multiple (LCM) of all Test values
	lcm := 1
	for _, monkey := range monkeys {
		lcm = lcmValue(lcm, monkey.Test)
	}

	for round := 0; round < rounds; round++ {
		for _, monkey := range monkeys {
			inspected[monkey.N] += len(items[monkey.N])
			for _, old := range items[monkey.N] {
				arg := old
				if monkey.Arg != 0 {
					arg = monkey.Arg
				}
				newVal := applyOperation(monkey.Op, old, arg)
				if divideByThree {
					newVal /= 3
				}
				newVal %= lcm // Keep worry levels manageable
				throwTo := monkey.IfFalse
				if newVal%monkey.Test == 0 {
					throwTo = monkey.IfTrue
				}
				items[throwTo] = append(items[throwTo], newVal)
			}
			items[monkey.N] = []int{}
		}
	}

	return inspected
}

func applyOperation(op string, old, arg int) int {
	switch op {
	case "+":
		return old + arg
	case "*":
		return old * arg
	default:
		return old
	}
}

func monkeyBusiness(inspected map[int]int) int {
	counts := []int{}
	for _, count := range inspected {
		counts = append(counts, count)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(counts)))
	return counts[0] * counts[1]
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcmValue(a, b int) int {
	return a / gcd(a, b) * b
}

func init() {
	solve.Register(Day11{})
}
