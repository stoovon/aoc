package solve2016

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day18 struct {
}

func (d Day18) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2016, Day: 18}
}

const (
	safe = '.'
	trap = '^'
)

// rows generates the first n rows of tiles given the initial row.
func (d Day18) rows(n int, initialRow string) []string {
	result := []string{initialRow}
	for i := 1; i < n; i++ {
		previous := string(safe) + result[len(result)-1] + string(safe)
		var nextRow strings.Builder
		for j := 1; j < len(previous)-1; j++ {
			if previous[j-1] != previous[j+1] {
				nextRow.WriteByte(trap)
			} else {
				nextRow.WriteByte(safe)
			}
		}
		result = append(result, nextRow.String())
	}
	return result
}

// countSafeTiles counts the number of safe tiles in all rows.
func (d Day18) countSafeTiles(rows []string) int {
	count := 0
	for _, row := range rows {
		count += strings.Count(row, string(safe))
	}
	return count
}

func (d Day18) Part1(input string) (string, error) {
	initialRow := strings.TrimSpace(input)
	safeCount := d.countSafeTiles(d.rows(40, initialRow))
	return strconv.Itoa(safeCount), nil
}

func (d Day18) Part2(input string) (string, error) {
	initialRow := strings.TrimSpace(input)
	safeCount := d.countSafeTiles(d.rows(400000, initialRow))
	return strconv.Itoa(safeCount), nil
}

func init() {
	solve.Register(Day18{})
}
