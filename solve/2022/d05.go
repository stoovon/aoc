package solve2022

import (
	"aoc/solve"
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type Day5 struct{}

func (d Day5) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2022, Day: 5}
}

func parseInput(input string) ([][]rune, []string) {
	parts := strings.Split(input, "\n\n")
	stacksStr, procedures := parts[0], parts[1]

	lines := strings.Split(stacksStr, "\n")
	numStacks := len(strings.Fields(lines[len(lines)-1]))
	stacks := make([][]rune, numStacks)

	for i := range stacks {
		stacks[i] = []rune{}
	}

	for i := len(lines) - 2; i >= 0; i-- {
		line := lines[i]
		for stackIdx, j := 0, 1; j < len(line); j, stackIdx = j+4, stackIdx+1 {
			if line[j] != ' ' {
				stacks[stackIdx] = append(stacks[stackIdx], rune(line[j]))
			}
		}
	}

	return stacks, strings.Split(procedures, "\n")
}

func parseInstruction(instruction string) (int, int, int, error) {
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(instruction, -1)
	if len(matches) != 3 {
		return 0, 0, 0, errors.New("invalid instruction format")
	}

	cratesToMove, _ := strconv.Atoi(matches[0])
	fromStack, _ := strconv.Atoi(matches[1])
	toStack, _ := strconv.Atoi(matches[2])

	return cratesToMove, fromStack - 1, toStack - 1, nil
}

func executeProcedures(stacks [][]rune, procedures []string, moveCrates func([][]rune, int, int, int) error) (string, error) {
	for _, procedure := range procedures {
		if procedure == "" {
			continue
		}
		cratesToMove, fromStack, toStack, err := parseInstruction(procedure)
		if err != nil {
			return "", err
		}
		if err := moveCrates(stacks, cratesToMove, fromStack, toStack); err != nil {
			return "", err
		}
	}

	result := ""
	for _, stack := range stacks {
		if len(stack) > 0 {
			result += string(stack[len(stack)-1])
		}
	}

	return result, nil
}

func moveCratesPart1(stacks [][]rune, cratesToMove, fromStack, toStack int) error {
	for i := 0; i < cratesToMove; i++ {
		if len(stacks[fromStack]) == 0 {
			return errors.New("stack underflow")
		}
		stacks[toStack] = append(stacks[toStack], stacks[fromStack][len(stacks[fromStack])-1])
		stacks[fromStack] = stacks[fromStack][:len(stacks[fromStack])-1]
	}
	return nil
}

func moveCratesPart2(stacks [][]rune, cratesToMove, fromStack, toStack int) error {
	if len(stacks[fromStack]) < cratesToMove {
		return errors.New("stack underflow")
	}
	stacks[toStack] = append(stacks[toStack], stacks[fromStack][len(stacks[fromStack])-cratesToMove:]...)
	stacks[fromStack] = stacks[fromStack][:len(stacks[fromStack])-cratesToMove]
	return nil
}

func (d Day5) Part1(input string) (string, error) {
	stacks, procedures := parseInput(input)
	return executeProcedures(stacks, procedures, moveCratesPart1)
}

func (d Day5) Part2(input string) (string, error) {
	stacks, procedures := parseInput(input)
	return executeProcedures(stacks, procedures, moveCratesPart2)
}

func init() {
	solve.Register(Day5{})
}
