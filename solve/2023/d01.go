package solve2023

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day1 struct {
}

func (d Day1) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2023, Day: 1}
}

func (d Day1) applyReplacer(r *strings.Replacer, input string) (result int) {
	for _, s := range strings.Fields(string(input)) {
		s = r.Replace(r.Replace(s))
		result += 10 * int(s[strings.IndexAny(s, "123456789")]-'0')
		result += int(s[strings.LastIndexAny(s, "123456789")] - '0')
	}
	return result
}

func (d Day1) Part1(input string) (string, error) {
	sum := d.applyReplacer(strings.NewReplacer(), input)

	return strconv.Itoa(sum), nil
}

func (d Day1) Part2(input string) (string, error) {
	sum := d.applyReplacer(strings.NewReplacer("one", "o1e", "two", "t2o", "three", "t3e", "four",
		"f4r", "five", "f5e", "six", "s6x", "seven", "s7n", "eight", "e8t", "nine", "n9e"), input)

	return strconv.Itoa(sum), nil
}

func init() {
	solve.Register(Day1{})
}
