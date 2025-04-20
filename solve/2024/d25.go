package solve2024

import (
	"image"
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/grids"
)

type Day25 struct {
}

func (d Day25) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2024, Day: 25}
}

var (
	DOWN   = image.Point{Y: 1}
	UP     = image.Point{Y: -1}
	ORIGIN = image.Point{}
)

// calculateHeights calculates the heights for locks or keys based on the direction.
func calculateHeights(grid [][]rune, start image.Point, direction image.Point) int {
	heights := 0
	for x := 0; x < 5; x++ {
		position := image.Point{X: x, Y: start.Y}
		for position.Y >= 0 && position.Y < len(grid) && grid[position.Y][position.X] == '#' {
			position = position.Add(direction)
		}
		if direction == DOWN {
			heights = (heights << 4) + (position.Y - 1)
		} else {
			heights = (heights << 4) + (start.Y - position.Y)
		}
	}
	return heights
}

func (d Day25) Part1(input string) (string, error) {
	data := strings.TrimSpace(input)
	locks := make([]int, 0)
	keys := make([]int, 0)
	result := 0

	blocks := strings.Split(data, "\n\n")
	for _, block := range blocks {
		grid := grids.NewGridOptions().Parse(block).ColumnsByRows()

		if grid[ORIGIN.Y][ORIGIN.X] == '#' {
			locks = append(locks, calculateHeights(grid, image.Point{Y: 1}, DOWN))
		} else {
			keys = append(keys, calculateHeights(grid, image.Point{Y: 5}, UP))
		}
	}

	for _, lock := range locks {
		for _, key := range keys {
			if (lock+key+0x22222)&0x88888 == 0 {
				result++
			}
		}
	}

	return strconv.Itoa(result), nil
}

func (d Day25) Part2(_ string) (string, error) {
	return "", nil
}

func init() {
	solve.Register(Day25{})
}
