package solve2023

import (
	"image"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day11 struct {
}

func (d Day11) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2023, Day: 11}
}

type skySnapshot struct {
	points     []image.Point
	horizontal []int
	vertical   []int
}

func (d Day11) parse(input string) skySnapshot {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	h, w := len(lines), len(lines[0])
	points := []image.Point{}
	rows := make([]bool, h)
	cols := make([]bool, w)
	for y, line := range lines {
		for x, c := range line {
			if c == '#' {
				points = append(points, image.Pt(x, y))
				rows[y] = true
				cols[x] = true
			}
		}
	}
	// Invert: true means empty
	for i := range rows {
		rows[i] = !rows[i]
	}
	for i := range cols {
		cols[i] = !cols[i]
	}
	horizontal := make([]int, h)
	vertical := make([]int, w)
	emptyRows, emptyCols := 0, 0
	for i := 0; i < h; i++ {
		if rows[i] {
			emptyRows++
		}
		horizontal[i] = emptyRows
	}
	for i := 0; i < w; i++ {
		if cols[i] {
			emptyCols++
		}
		vertical[i] = emptyCols
	}
	return skySnapshot{points, horizontal, vertical}
}

func expand(p skySnapshot, times int) []image.Point {
	expanded := make([]image.Point, len(p.points))
	for i, pt := range p.points {
		x := pt.X + (times-1)*p.vertical[pt.X]
		y := pt.Y + (times-1)*p.horizontal[pt.Y]
		expanded[i] = image.Pt(x, y)
	}
	return expanded
}

func manhattan(a, b image.Point) int {
	dx := a.X - b.X
	if dx < 0 {
		dx = -dx
	}
	dy := a.Y - b.Y
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}

func sumDistances(points []image.Point) int64 {
	var sum int64
	for i := 1; i < len(points); i++ {
		for j := 0; j < i; j++ {
			sum += int64(manhattan(points[i], points[j]))
		}
	}
	return sum
}

func (d Day11) Part1(input string) (string, error) {
	p := d.parse(input)
	expanded := expand(p, 2)
	return strconv.FormatInt(sumDistances(expanded), 10), nil
}

func (d Day11) Part2(input string) (string, error) {
	p := d.parse(input)
	expanded := expand(p, 1000000)
	return strconv.FormatInt(sumDistances(expanded), 10), nil
}

func init() {
	solve.Register(Day11{})
}
