package solve2021

import (
	"aoc/solve"
	"errors"
	"sort"
	"strconv"
	"strings"
)

type Day10 struct{}

func parseChunks(line string, openToClose map[rune]rune) (isCorrupted bool, firstIllegalChar rune, stack []rune) {
	var workingStack []rune
	for _, c := range line {
		switch c {
		case '(', '[', '{', '<':
			workingStack = append(workingStack, c)
		case ')', ']', '}', '>':
			if len(workingStack) == 0 || openToClose[workingStack[len(workingStack)-1]] != c {
				return true, c, nil
			}
			workingStack = workingStack[:len(workingStack)-1]
		}
	}
	return false, 0, workingStack
}

func (d Day10) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2021, Day: 10}
}

func (d Day10) Part1(input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", errors.New("empty input")
	}
	lines := strings.Split(input, "\n")
	scoreTable := map[rune]int{
		')': 3,
		']': 57,
		'}': 1197,
		'>': 25137,
	}
	openToClose := map[rune]rune{
		'(': ')',
		'[': ']',
		'{': '}',
		'<': '>',
	}
	totalScore := 0
	for _, line := range lines {
		corrupted, illegal, _ := parseChunks(line, openToClose)
		if corrupted {
			totalScore += scoreTable[illegal]
		}
	}
	return strconv.Itoa(totalScore), nil
}

func (d Day10) Part2(input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", errors.New("empty input")
	}
	lines := strings.Split(input, "\n")
	openToClose := map[rune]rune{
		'(': ')',
		'[': ']',
		'{': '}',
		'<': '>',
	}
	scoreTable := map[rune]int{
		')': 1,
		']': 2,
		'}': 3,
		'>': 4,
	}
	var scores []int
	for _, line := range lines {
		corrupted, _, stack := parseChunks(line, openToClose)
		if corrupted || len(stack) == 0 {
			continue
		}

		score := 0
		for i := len(stack) - 1; i >= 0; i-- {
			closeChar := openToClose[stack[i]]
			score = score*5 + scoreTable[closeChar]
		}
		scores = append(scores, score)
	}
	if len(scores) == 0 {
		return "", errors.New("no incomplete lines")
	}
	sort.Ints(scores)
	mid := scores[len(scores)/2]
	return strconv.Itoa(mid), nil
}

func init() {
	solve.Register(Day10{})
}
