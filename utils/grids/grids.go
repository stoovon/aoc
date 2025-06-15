package grids

import (
	"image"
	"strings"
)

// WidePoint works just like image.Point, but uses int64 for coordinates in order
// to avoid performance issues with repeated conversions between int and int64, and
// also to handle larger grids without overflow.
type WidePoint struct {
	X int64
	Y int64
}

type GridOptions struct {
	Separator string
	Cutset    string
}

func NewGridOptions() *GridOptions {
	return &GridOptions{
		Separator: "\n",
		Cutset:    "#",
	}
}

func (o *GridOptions) WithSeparator(separator string) *GridOptions {
	o.Separator = separator
	return o
}

type Grid struct {
	rawData [][]rune
	cutset  string
}

func (o *GridOptions) Parse(block string) Grid {
	lines := strings.Split(block, o.Separator)
	grid := make([][]rune, 0)

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		grid = append(grid, []rune(line))
	}
	return Grid{rawData: grid, cutset: o.Cutset}
}

func (g Grid) ColumnsByRows() [][]rune {
	return g.rawData
}

func (g Grid) PointsFromTopLeft() map[rune]image.Point {
	result := make(map[rune]image.Point)
	for y, row := range g.rawData {
		for x, char := range row {
			if strings.Contains(g.cutset, string(char)) {
				continue
			}
			result[char] = image.Point{X: y, Y: x}
		}
	}
	return result
}

func (g Grid) CharsFromTopLeft() map[image.Point]rune {
	grid := map[image.Point]rune{}
	for y, row := range g.rawData {
		for x, char := range row {
			grid[image.Point{X: x, Y: y}] = char
		}
	}

	return grid
}

func ParseCharsFromTopLeft(block string) map[image.Point]rune {
	return NewGridOptions().Parse(block).CharsFromTopLeft()
}
