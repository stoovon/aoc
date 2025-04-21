package solve2015

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day10 struct {
}

func (d Day10) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2015, Day: 10}
}

func expand(num string) string {
	var result strings.Builder
	count := 1

	for i := 1; i < len(num); i++ {
		if num[i] == num[i-1] {
			count++
		} else {
			result.WriteString(strconv.Itoa(count))
			result.WriteByte(num[i-1])
			count = 1
		}
	}

	result.WriteString(strconv.Itoa(count))
	result.WriteByte(num[len(num)-1])

	return result.String()
}

func (d Day10) lookAndSay(input string, iterations int) int {
	line := strings.TrimSpace(input)
	for i := 0; i < iterations; i++ {
		line = expand(line)
	}

	return len(line)
}

func (d Day10) Part1(input string) (string, error) {
	return strconv.Itoa(d.lookAndSay(input, 40)), nil
}

func (d Day10) Part2(input string) (string, error) {
	return strconv.Itoa(d.lookAndSay(input, 50)), nil
}

func init() {
	solve.Register(Day10{})
}
