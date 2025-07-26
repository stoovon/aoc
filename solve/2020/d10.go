package solve2020

import (
	"aoc/solve"
	"sort"
	"strconv"
	"strings"
)

type Day10 struct{}

func (d Day10) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2020, Day: 10}
}

func parseAdapters(input string) ([]int, error) {
	lines := strings.Fields(strings.ReplaceAll(input, "\r\n", "\n"))
	nums := make([]int, len(lines))
	for i, line := range lines {
		n, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}
		nums[i] = n
	}
	sort.Ints(nums)
	return nums, nil
}

func (d Day10) Part1(input string) (string, error) {
	nums, err := parseAdapters(input)
	if err != nil {
		return "", err
	}
	ones, threes := 0, 0
	prev := 0
	for _, n := range nums {
		diff := n - prev
		if diff == 1 {
			ones++
		} else if diff == 3 {
			threes++
		}
		prev = n
	}
	threes++ // device's built-in adapter
	return strconv.Itoa(ones * threes), nil
}

func (d Day10) Part2(input string) (string, error) {
	nums, err := parseAdapters(input)
	if err != nil {
		return "", err
	}
	device := nums[len(nums)-1] + 3
	nums = append([]int{0}, nums...)
	nums = append(nums, device)
	ways := make(map[int]int)
	ways[0] = 1
	for _, n := range nums[1:] {
		ways[n] = ways[n-1] + ways[n-2] + ways[n-3]
	}
	return strconv.Itoa(ways[device]), nil
}

func init() {
	solve.Register(Day10{})
}
