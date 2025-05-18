package solve2023

import (
	"image"
	"math"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day21 struct {
}

func (d Day21) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2023, Day: 21}
}

type Map struct {
	Width, Height int
	Tiles         [][]rune
}

func NewMap(input string) *Map {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	tiles := make([][]rune, len(lines))
	for i, line := range lines {
		tiles[i] = []rune(line)
	}
	return &Map{
		Width:  len(tiles[0]),
		Height: len(tiles),
		Tiles:  tiles,
	}
}

func (m *Map) FindChar(ch rune) image.Point {
	for y, row := range m.Tiles {
		for x, c := range row {
			if c == ch {
				return image.Point{X: x, Y: y}
			}
		}
	}
	panic("start not found")
}

func (m *Map) Translate(p image.Point) image.Point {
	x := ((p.X % m.Width) + m.Width) % m.Width
	y := ((p.Y % m.Height) + m.Height) % m.Height
	return image.Point{X: x, Y: y}
}

func (m *Map) Get(p image.Point) rune {
	pt := m.Translate(p)
	return m.Tiles[pt.Y][pt.X]
}

func (m *Map) Moves(p image.Point) []image.Point {
	dirs := []image.Point{
		{X: -1, Y: 0}, {X: 1, Y: 0},
		{X: 0, Y: -1}, {X: 0, Y: 1},
	}
	var res []image.Point
	for _, d := range dirs {
		np := image.Point{X: p.X + d.X, Y: p.Y + d.Y}
		if m.Get(np) != '#' {
			res = append(res, np)
		}
	}
	return res
}

func (d Day21) Part1(input string) (string, error) {
	m := NewMap(input)
	cur := map[image.Point]struct{}{m.FindChar('S'): {}}
	for step := 0; step < 64; step++ {
		next := make(map[image.Point]struct{})
		for p := range cur {
			for _, np := range m.Moves(p) {
				next[np] = struct{}{}
			}
		}
		cur = next
	}
	return strconv.Itoa(len(cur)), nil
}

// Helper: solve quadratic a*x^2 + b*x + c = y for three points
func fitQuadratic(x, y [3]float64) (a, b, c float64) {
	det := func(a, b, c, d, e, f, g, h, i float64) float64 {
		return a*e*i + b*f*g + c*d*h - c*e*g - b*d*i - a*f*h
	}
	A := [3][3]float64{
		{x[0] * x[0], x[0], 1},
		{x[1] * x[1], x[1], 1},
		{x[2] * x[2], x[2], 1},
	}
	Y := y
	D := det(
		A[0][0], A[0][1], A[0][2],
		A[1][0], A[1][1], A[1][2],
		A[2][0], A[2][1], A[2][2],
	)
	Da := det(
		Y[0], A[0][1], A[0][2],
		Y[1], A[1][1], A[1][2],
		Y[2], A[2][1], A[2][2],
	)
	Db := det(
		A[0][0], Y[0], A[0][2],
		A[1][0], Y[1], A[1][2],
		A[2][0], Y[2], A[2][2],
	)
	Dc := det(
		A[0][0], A[0][1], Y[0],
		A[1][0], A[1][1], Y[1],
		A[2][0], A[2][1], Y[2],
	)
	return Da / D, Db / D, Dc / D
}

func (d Day21) Part2(input string) (string, error) {
	m := NewMap(input)
	maxSteps := 26501365
	start := m.FindChar('S')
	cur := map[image.Point]struct{}{start: {}}
	width := m.Width

	var res []float64
	var i int
	for {
		next := make(map[image.Point]struct{})
		for p := range cur {
			for _, np := range m.Moves(p) {
				next[np] = struct{}{}
			}
		}
		cur = next
		i++
		if i%width == maxSteps%width {
			res = append(res, float64(len(cur)))
			if len(res) == 3 {
				break
			}
		}
	}
	x := [3]float64{0, 1, 2}
	y := [3]float64{res[0], res[1], res[2]}
	a, b, c := fitQuadratic(x, y)
	X := float64(maxSteps / width)
	ans := a*X*X + b*X + c
	return strconv.FormatInt(int64(math.Round(ans)), 10), nil
}

func init() {
	solve.Register(Day21{})
}
