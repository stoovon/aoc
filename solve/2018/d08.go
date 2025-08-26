package solve2018

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day8 struct {
}

func (d Day8) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2018, Day: 8}
}

func (d Day8) Part1(input string) (string, error) {
	nums := d.parseInput(input)
	sum, _ := sumMetadata(nums, 0)
	return strconv.Itoa(sum), nil
}

func (d Day8) parseInput(input string) []int {
	fields := strings.Fields(input)
	nums := make([]int, len(fields))
	for i, f := range fields {
		nums[i], _ = strconv.Atoi(f)
	}
	return nums
}

func sumMetadata(nums []int, pos int) (sum int, next int) {
	children := nums[pos]
	metadata := nums[pos+1]
	pos += 2
	for i := 0; i < children; i++ {
		var childSum int
		childSum, pos = sumMetadata(nums, pos)
		sum += childSum
	}
	for i := 0; i < metadata; i++ {
		sum += nums[pos]
		pos++
	}
	return sum, pos
}

func (d Day8) Part2(input string) (string, error) {
	nums := d.parseInput(input)
	value, _ := nodeValue(nums, 0)
	return strconv.Itoa(value), nil
}

func nodeValue(nums []int, pos int) (value int, next int) {
	children := nums[pos]
	metadata := nums[pos+1]
	pos += 2
	childValues := make([]int, children)
	for i := 0; i < children; i++ {
		childValues[i], pos = nodeValue(nums, pos)
	}
	if children == 0 {
		for i := 0; i < metadata; i++ {
			value += nums[pos]
			pos++
		}
	} else {
		for i := 0; i < metadata; i++ {
			idx := nums[pos] - 1
			if idx >= 0 && idx < children {
				value += childValues[idx]
			}
			pos++
		}
	}
	return value, pos
}

func init() {
	solve.Register(Day8{})
}
