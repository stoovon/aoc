package solve2022

import (
	"aoc/solve"
	"sort"
	"strconv"
	"strings"
)

type Day21 struct{}

func (d Day21) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2022, Day: 21}
}

func (d Day21) Part1(input string) (string, error) {
	monkeys := make(map[string]string)
	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		parts := strings.Split(line, ": ")
		monkeys[parts[0]] = parts[1]
	}

	result := d.solve("root", monkeys)

	return strconv.Itoa(result), nil
}

func (d Day21) solve(expr string, monkeys map[string]string) int {
	if v, err := strconv.Atoi(monkeys[expr]); err == nil {
		return v
	}

	s := strings.Fields(monkeys[expr])
	return map[string]func(int, int) int{
		"+": func(a, b int) int { return a + b },
		"-": func(a, b int) int { return a - b },
		"*": func(a, b int) int { return a * b },
		"/": func(a, b int) int { return a / b },
	}[s[1]](d.solve(s[0], monkeys), d.solve(s[2], monkeys))
}

func (d Day21) Part2(input string) (string, error) {
	monkeys := make(map[string]string)

	for s := range strings.SplitSeq(strings.TrimSpace(string(input)), "\n") {
		s := strings.Split(s, ": ")
		monkeys[s[0]] = s[1]
	}

	monkeys["humn"] = "0"
	s := strings.Fields(monkeys["root"])
	if d.solve(s[0], monkeys) < d.solve(s[2], monkeys) {
		s[0], s[2] = s[2], s[0]
	}

	part2, _ := sort.Find(1e16, func(v int) int {
		monkeys["humn"] = strconv.Itoa(v)
		return d.solve(s[0], monkeys) - d.solve(s[2], monkeys)
	})

	return strconv.Itoa(part2), nil
}

func init() {
	solve.Register(Day21{})
}
