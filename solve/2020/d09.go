package solve2020

import (
	"aoc/solve"
	"errors"
	"strconv"
	"strings"
)

type Day9 struct{}

func (d Day9) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2020, Day: 9}
}

func parseNums(input string) ([]int, error) {
	lines := strings.Fields(strings.ReplaceAll(input, "\r\n", "\n"))
	nums := make([]int, len(lines))
	for i, line := range lines {
		n, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}
		nums[i] = n
	}
	return nums, nil
}

func findInvalid(nums []int, preamble int) (int, error) {
	for i := preamble; i < len(nums); i++ {
		valid := false
		for j := i - preamble; j < i-1; j++ {
			for k := j + 1; k < i; k++ {
				if nums[j] != nums[k] && nums[j]+nums[k] == nums[i] {
					valid = true
					break
				}
			}
			if valid {
				break
			}
		}
		if !valid {
			return nums[i], nil
		}
	}
	return -1, errors.New("all numbers valid")
}

func (d Day9) Part1(input string) (string, error) {
	nums, err := parseNums(input)
	if err != nil {
		return "", err
	}
	invalid, err := findInvalid(nums, 25)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(invalid), nil
}

func (d Day9) Part2(input string) (string, error) {
	nums, err := parseNums(input)
	if err != nil {
		return "", err
	}
	invalid, err := findInvalid(nums, 25)
	if err != nil {
		return "", err
	}
	for start := 0; start < len(nums); start++ {
		sum := nums[start]
		min, max := nums[start], nums[start]
		for end := start + 1; end < len(nums); end++ {
			sum += nums[end]
			if nums[end] < min {
				min = nums[end]
			}
			if nums[end] > max {
				max = nums[end]
			}
			if sum == invalid {
				return strconv.Itoa(min + max), nil
			}
			if sum > invalid {
				break
			}
		}
	}
	return "", errors.New("no contiguous range found")
}

func init() {
	solve.Register(Day9{})
}
