package solve2023

import (
	"image"
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/maths"
)

func (d Day10) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2023, Day: 10}
}

var (
	UP    = image.Pt(0, -1)
	DOWN  = image.Pt(0, 1)
	LEFT  = image.Pt(-1, 0)
	RIGHT = image.Pt(1, 0)
)

type day10Grid struct {
	Width, Height int
	Cells         [][]byte
}

func (d Day10) ParseGrid(input string) *day10Grid {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	height := len(lines)
	width := len(lines[0])
	cells := make([][]byte, height)
	for y, line := range lines {
		cells[y] = []byte(line)
	}
	return &day10Grid{Width: width, Height: height, Cells: cells}
}

func (g *day10Grid) At(p image.Point) byte {
	return g.Cells[p.Y][p.X]
}

func (g *day10Grid) Find(b byte) (image.Point, bool) {
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if g.Cells[y][x] == b {
				return image.Pt(x, y), true
			}
		}
	}
	return image.Point{}, false
}

func determinant(a, b image.Point) int {
	return a.X*b.Y - a.Y*b.X
}

func (d Day10) parse(input string) (int, int) {
	grid := d.ParseGrid(input)
	corner, found := grid.Find('S')
	if !found {
		return 0, 0
	}

	var direction image.Point
	if c := grid.At(corner.Add(UP)); c == '|' || c == '7' || c == 'F' {
		direction = UP
	} else {
		direction = DOWN
	}
	position := corner.Add(direction)

	perimeter := 1
	area := 0

	for {
		for c := grid.At(position); c == '-' || c == '|'; c = grid.At(position) {
			position = position.Add(direction)
			perimeter++
		}

		c := grid.At(position)
		switch {
		case c == '7' && direction == UP:
			direction = LEFT
		case c == 'F' && direction == UP:
			direction = RIGHT
		case c == 'J' && direction == DOWN:
			direction = LEFT
		case c == 'L' && direction == DOWN:
			direction = RIGHT
		case c == 'J' || c == 'L':
			direction = UP
		case c == '7' || c == 'F':
			direction = DOWN
		default:
			area += determinant(corner, position)
			goto done
		}
		area += determinant(corner, position)
		corner = position
		position = position.Add(direction)
		perimeter++
	}
done:
	partOne := perimeter / 2
	partTwo := maths.Abs(area)/2 - perimeter/2 + 1
	return partOne, partTwo
}

type Day10 struct{}

func (d Day10) Part1(input string) (string, error) {
	partOne, _ := d.parse(input)
	return strconv.Itoa(partOne), nil
}

func (d Day10) Part2(input string) (string, error) {
	_, partTwo := d.parse(input)
	return strconv.Itoa(partTwo), nil
}

func init() {
	solve.Register(Day10{})
}
