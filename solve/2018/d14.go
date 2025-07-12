package solve2018

import (
	"aoc/solve"
	"strconv"
	"strings"
)

type Day14 struct{}

func (d Day14) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2018, Day: 14}
}

// shared recipe simulation
func makeRecipes(target int, pattern []int, findPattern bool) (string, int) {
	scores := []int{3, 7}
	elf1, elf2 := 0, 1
	for {
		sum := scores[elf1] + scores[elf2]
		if sum >= 10 {
			scores = append(scores, sum/10)
		}
		scores = append(scores, sum%10)
		elf1 = (elf1 + 1 + scores[elf1]) % len(scores)
		elf2 = (elf2 + 1 + scores[elf2]) % len(scores)
		if !findPattern && len(scores) >= target+10 {
			res := ""
			for i := target; i < target+10; i++ {
				res += strconv.Itoa(scores[i])
			}
			return res, -1
		}
		if findPattern {
			for k := 1; k <= 2; k++ {
				if len(scores) >= len(pattern)+k-1 {
					start := len(scores) - len(pattern) - (2 - k)
					if start >= 0 {
						match := true
						for i := 0; i < len(pattern); i++ {
							if scores[start+i] != pattern[i] {
								match = false
								break
							}
						}
						if match {
							return "", start
						}
					}
				}
			}
		}
	}
}

func (d Day14) Part1(input string) (string, error) {
	input = strings.TrimSpace(input)
	n, err := strconv.Atoi(input)
	if err != nil {
		return "", err
	}
	res, _ := makeRecipes(n, nil, false)
	return res, nil
}

func (d Day14) Part2(input string) (string, error) {
	input = strings.TrimSpace(input)
	pattern := make([]int, len(input))
	for i := range input {
		pattern[i] = int(input[i] - '0')
	}
	_, idx := makeRecipes(0, pattern, true)
	return strconv.Itoa(idx), nil
}

func init() {
	solve.Register(Day14{})
}
