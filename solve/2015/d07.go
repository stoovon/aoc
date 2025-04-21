package solve2015

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/acstrings"
)

type Day7 struct {
}

var (
	lowerCase = regexp.MustCompile(`^[a-z]+$`)
)

func (d Day7) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2015, Day: 7}
}

const maxInt = 0xFFFF

func createValueFunc(rules map[string]string) func(string) int {
	cache := make(map[string]int)

	var getValue func(string) int
	getValue = func(key string) int {
		if val, ok := cache[key]; ok {
			return val
		}

		if num, err := strconv.Atoi(key); err == nil {
			return num
		}

		value := rules[key]
		if num, err := strconv.Atoi(value); err == nil {
			cache[key] = num
			return num
		}

		switch {
		case lowerCase.MatchString(value):
			cache[key] = getValue(value)
		case strings.HasPrefix(value, "NOT "):
			arg := getValue(value[4:])
			cache[key] = maxInt - arg
		case strings.Contains(value, " AND "):
			parts := strings.Split(value, " AND ")
			arg1, arg2 := getValue(parts[0]), getValue(parts[1])
			cache[key] = arg1 & arg2
		case strings.Contains(value, " OR "):
			parts := strings.Split(value, " OR ")
			arg1, arg2 := getValue(parts[0]), getValue(parts[1])
			cache[key] = arg1 | arg2
		case strings.Contains(value, " LSHIFT "):
			parts := strings.Split(value, " LSHIFT ")
			arg1, shift := getValue(parts[0]), acstrings.MustInt(parts[1])
			cache[key] = (arg1 << shift) & maxInt
		case strings.Contains(value, " RSHIFT "):
			parts := strings.Split(value, " RSHIFT ")
			arg1, shift := getValue(parts[0]), acstrings.MustInt(parts[1])
			cache[key] = arg1 >> shift
		default:
			panic(fmt.Sprintf("unsupported operation: %s", value))
		}

		return cache[key]
	}

	return getValue
}

func (d Day7) solve(input string) (part1, part2 int) {
	lines := strings.Split(input, "\n")
	rules := make(map[string]string, len(lines))

	for _, line := range lines {
		parts := strings.Split(line, " -> ")

		if len(parts) != 2 {
			continue
		}

		rules[parts[1]] = parts[0]
	}

	getValuePart1 := createValueFunc(rules)
	part1 = getValuePart1("a")

	rules["b"] = strconv.Itoa(part1)
	getValuePart2 := createValueFunc(rules)
	part2 = getValuePart2("a")

	return part1, part2
}

func (d Day7) Part1(input string) (string, error) {
	part1, _ := d.solve(input)
	return strconv.Itoa(part1), nil
}

func (d Day7) Part2(input string) (string, error) {
	_, part2 := d.solve(input)
	return strconv.Itoa(part2), nil
}

func init() {
	solve.Register(Day7{})
}
