package solve2021

import (
	"fmt"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day2 struct {
}

func (d Day2) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2021, Day: 2}
}

func (d Day2) Part1(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	horizontal, depth := 0, 0
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid command: %s", line)
		}
		val, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("invalid value in command: %s", line)
		}
		switch parts[0] {
		case "forward":
			horizontal += val
		case "down":
			depth += val
		case "up":
			depth -= val
		default:
			return "", fmt.Errorf("unknown command: %s", parts[0])
		}
	}
	return strconv.Itoa(horizontal * depth), nil
}

func (d Day2) Part2(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	horizontal, depth, aim := 0, 0, 0
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid command: %s", line)
		}
		val, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("invalid value in command: %s", line)
		}
		switch parts[0] {
		case "forward":
			horizontal += val
			depth += aim * val
		case "down":
			aim += val
		case "up":
			aim -= val
		default:
			return "", fmt.Errorf("unknown command: %s", parts[0])
		}
	}
	return strconv.Itoa(horizontal * depth), nil
}

func init() {
	solve.Register(Day2{})
}
