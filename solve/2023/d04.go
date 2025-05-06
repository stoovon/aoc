package solve2023

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day4 struct {
}

func (d Day4) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2023, Day: 4}
}

func (d Day4) solve(input string) (int, int) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	p1 := 0
	p2 := 0
	multiplier := make([]int, len(lines))
	for i := range multiplier {
		multiplier[i] = 1
	}

	for i, line := range lines {
		parts := strings.Split(line, "|")
		winning := d.parseNumbers(strings.Split(parts[0], ":")[1])
		have := d.parseNumbers(parts[1])

		filtered := []int{}
		for idx, _ := range have {
			if _, ok := winning[idx]; ok {
				filtered = append(filtered, idx)
			}
		}

		if len(filtered) > 0 {
			p1 += 1 << (len(filtered) - 1)
		}

		mymult := multiplier[i]
		for j := i + 1; j < i+1+len(filtered) && j < len(lines); j++ {
			multiplier[j] += mymult
		}
		p2 += mymult
	}

	return p1, p2
}

func (d Day4) parseNumbers(s string) map[int]bool {
	nums := strings.Fields(s)
	result := make(map[int]bool, len(nums))
	for _, num := range nums {
		n, _ := strconv.Atoi(num)
		result[n] = true
	}
	return result
}

func (d Day4) Part1(input string) (string, error) {
	p1, _ := d.solve(input)
	return strconv.Itoa(p1), nil
}

func (d Day4) Part2(input string) (string, error) {
	_, p2 := d.solve(input)
	return strconv.Itoa(p2), nil
}

func init() {
	solve.Register(Day4{})
}
