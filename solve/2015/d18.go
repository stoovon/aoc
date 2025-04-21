package solve2015

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day18 struct {
}

func (d Day18) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2015, Day: 18}
}

func round(grid []string, brokenCorners bool) []string {
	next := make([]string, len(grid))
	for i := range grid {
		var temp strings.Builder
		for j := range grid[i] {
			count := 0
			up := i > 0
			down := i < len(grid)-1
			left := j > 0
			right := j < len(grid[i])-1

			if up {
				if left && grid[i-1][j-1] == '#' {
					count++
				}
				if grid[i-1][j] == '#' {
					count++
				}
				if right && grid[i-1][j+1] == '#' {
					count++
				}
			}
			if right && grid[i][j+1] == '#' {
				count++
			}
			if left && grid[i][j-1] == '#' {
				count++
			}
			if down {
				if left && grid[i+1][j-1] == '#' {
					count++
				}
				if grid[i+1][j] == '#' {
					count++
				}
				if right && grid[i+1][j+1] == '#' {
					count++
				}
			}

			if grid[i][j] == '.' && count == 3 {
				temp.WriteByte('#')
			} else if grid[i][j] == '#' && (count == 2 || count == 3) {
				temp.WriteByte('#')
			} else {
				temp.WriteByte('.')
			}
		}
		next[i] = temp.String()
	}

	if brokenCorners {
		next[0] = "#" + next[0][1:len(next[0])-1] + "#"
		next[len(next)-1] = "#" + next[len(next)-1][1:len(next[len(next)-1])-1] + "#"
	}

	return next
}

func countNums(rounds int, grid []string, brokenCorners bool) int {
	if brokenCorners {
		grid[0] = "#" + grid[0][1:len(grid[0])-1] + "#"
		grid[len(grid)-1] = "#" + grid[len(grid)-1][1:len(grid[len(grid)-1])-1] + "#"
	}

	for i := 0; i < rounds; i++ {
		grid = round(grid, brokenCorners)
	}

	count := 0
	for _, row := range grid {
		for _, cell := range row {
			if cell == '#' {
				count++
			}
		}
	}

	return count
}

func (d Day18) Part1(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	return strconv.Itoa(countNums(100, lines, false)), nil
}

func (d Day18) Part2(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	return strconv.Itoa(countNums(100, lines, true)), nil
}

func init() {
	solve.Register(Day18{})
}
