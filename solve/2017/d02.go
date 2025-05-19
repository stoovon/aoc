package solve2017

import (
	"fmt"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day2 struct {
}

func (d Day2) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2017, Day: 2}
}

func parseRows(input string) [][]int {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	rows := make([][]int, len(lines))
	for i, line := range lines {
		fields := strings.Fields(line)
		nums := make([]int, len(fields))
		for j, f := range fields {
			nums[j], _ = strconv.Atoi(f)
		}
		rows[i] = nums
	}
	return rows
}

func (d Day2) Part1(input string) (string, error) {
	rows := parseRows(input)
	sum := 0
	for _, row := range rows {
		minVal, maxVal := row[0], row[0]
		for _, n := range row[1:] {
			if n < minVal {
				minVal = n
			}
			if n > maxVal {
				maxVal = n
			}
		}
		sum += maxVal - minVal
	}
	return fmt.Sprintf("%d", sum), nil
}

func (d Day2) Part2(input string) (string, error) {
	rows := parseRows(input)
	sum := 0
	for _, row := range rows {
		found := false
		for i := 0; i < len(row); i++ {
			for j := 0; j < len(row); j++ {
				if i != j && row[i]%row[j] == 0 {
					sum += row[i] / row[j]
					found = true
					break
				}
			}
			if found {
				break
			}
		}
	}
	return fmt.Sprintf("%d", sum), nil
}

func init() {
	solve.Register(Day2{})
}
