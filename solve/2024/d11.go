package solve2024

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day11 struct {
}

func (d Day11) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2024, Day: 11}
}

func (d Day11) run(stones map[int]int, blinks int) (r int) {
	for range blinks {
		next := map[int]int{}
		for k, v := range stones {
			if k == 0 {
				next[1] += v
			} else if s := strconv.Itoa(k); len(s)%2 == 0 {
				n1, _ := strconv.Atoi(s[:len(s)/2])
				n2, _ := strconv.Atoi(s[len(s)/2:])
				next[n1] += v
				next[n2] += v
			} else {
				next[k*2024] += v
			}
		}
		stones = next
	}
	for _, v := range stones {
		r += v
	}
	return r
}

func (d Day11) Part1(input string) (string, error) {
	stones := map[int]int{}
	for _, s := range strings.Split(strings.TrimSpace(string(input)), " ") {
		n, _ := strconv.Atoi(s)
		stones[n]++
	}

	return strconv.Itoa(d.run(stones, 25)), nil
}

func (d Day11) Part2(input string) (string, error) {
	stones := map[int]int{}
	for _, s := range strings.Split(strings.TrimSpace(string(input)), " ") {
		n, _ := strconv.Atoi(s)
		stones[n]++
	}

	return strconv.Itoa(d.run(stones, 75)), nil
}

func init() {
	solve.Register(Day11{})
}
