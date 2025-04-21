package solve2022

import (
	"sort"
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/maths"
)

type Day1 struct {
}

func (d Day1) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2022, Day: 1}
}

func calculateSums(input string) ([]int, error) {
	invs := strings.Split(strings.TrimSpace(input), "\n\n")
	invSum := make([]int, len(invs))

	for i, inv := range invs {
		lines := strings.Split(strings.TrimSpace(inv), "\n")
		sum := 0
		for _, line := range lines {
			calorie, err := strconv.Atoi(line)
			if err != nil {
				return nil, err
			}
			sum += calorie
		}
		invSum[i] = sum
	}
	return invSum, nil
}

func (d Day1) Part1(input string) (string, error) {
	invSum, err := calculateSums(input)
	if err != nil {
		return "", err
	}
	maxSum := maths.Max(invSum...)
	return strconv.Itoa(maxSum), nil
}

func (d Day1) Part2(input string) (string, error) {
	invSum, err := calculateSums(input)
	if err != nil {
		return "", err
	}
	sort.Sort(sort.Reverse(sort.IntSlice(invSum)))
	topThreeSum := maths.Sum(invSum[:3]...)
	return strconv.Itoa(topThreeSum), nil
}

func init() {
	solve.Register(Day1{})
}
