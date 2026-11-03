package solve2025

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day9 struct{}

type theaterPoint struct{ x, y int64 }

func (d Day9) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2025, Day: 9}
}

func parseTheaterPoints(input string) ([]theaterPoint, error) {
	var points []theaterPoint
	for index, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		xText, yText, ok := strings.Cut(line, ",")
		if !ok {
			return nil, fmt.Errorf("line %d: expected x,y coordinates", index+1)
		}
		x, err := strconv.ParseInt(strings.TrimSpace(xText), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: x coordinate: %w", index+1, err)
		}
		y, err := strconv.ParseInt(strings.TrimSpace(yText), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: y coordinate: %w", index+1, err)
		}
		points = append(points, theaterPoint{x, y})
	}
	if len(points) < 2 {
		return nil, errors.New("at least two red tiles are required")
	}
	return points, nil
}

func theaterArea(a, b theaterPoint) (int64, error) {
	width := uint64(max(a.x, b.x)) - uint64(min(a.x, b.x))
	height := uint64(max(a.y, b.y)) - uint64(min(a.y, b.y))
	if width >= math.MaxInt64 || height >= math.MaxInt64 {
		return 0, errors.New("rectangle dimensions exceed int64")
	}
	w, h := int64(width)+1, int64(height)+1
	if w > math.MaxInt64/h {
		return 0, errors.New("rectangle area exceeds int64")
	}
	return w * h, nil
}

func theaterEdgesIntersect(a, b, c, d theaterPoint) bool {
	if a.x == b.x && c.x == d.x {
		return a.x == c.x && max(min(a.y, b.y), min(c.y, d.y)) <= min(max(a.y, b.y), max(c.y, d.y))
	}
	if a.y == b.y && c.y == d.y {
		return a.y == c.y && max(min(a.x, b.x), min(c.x, d.x)) <= min(max(a.x, b.x), max(c.x, d.x))
	}
	if a.y == b.y {
		a, b, c, d = c, d, a, b
	}
	return a.x >= min(c.x, d.x) && a.x <= max(c.x, d.x) && c.y >= min(a.y, b.y) && c.y <= max(a.y, b.y)
}

func validateTheaterPolygon(points []theaterPoint) error {
	if len(points) < 4 {
		return errors.New("polygon requires at least four vertices")
	}
	seen := make(map[theaterPoint]bool, len(points))
	for i, a := range points {
		b := points[(i+1)%len(points)]
		if seen[a] {
			return errors.New("polygon contains a repeated vertex")
		}
		seen[a] = true
		if a == b || (a.x != b.x && a.y != b.y) {
			return errors.New("polygon edges must be nonzero horizontal or vertical segments")
		}
	}
	for i, a := range points {
		b := points[(i+1)%len(points)]
		for j := i + 1; j < len(points); j++ {
			c, d := points[j], points[(j+1)%len(points)]
			if !theaterEdgesIntersect(a, b, c, d) {
				continue
			}
			adjacent := j == i+1 || (i == 0 && j == len(points)-1)
			if !adjacent {
				return errors.New("polygon edges intersect or touch away from adjacent vertices")
			}
			if a.x == b.x && c.x == d.x && max(min(a.y, b.y), min(c.y, d.y)) < min(max(a.y, b.y), max(c.y, d.y)) {
				return errors.New("adjacent polygon edges overlap")
			}
			if a.y == b.y && c.y == d.y && max(min(a.x, b.x), min(c.x, d.x)) < min(max(a.x, b.x), max(c.x, d.x)) {
				return errors.New("adjacent polygon edges overlap")
			}
		}
	}
	return nil
}

func theaterAxis(values []int64) ([]int64, map[int64]int) {
	slices.Sort(values)
	values = slices.Compact(values)
	axis := make([]int64, 0, len(values))
	indices := make(map[int64]int, len(values))
	for i, value := range values {
		if i > 0 && values[i-1] < value-1 {
			axis = append(axis, values[i-1]+1)
		}
		indices[value] = len(axis)
		axis = append(axis, value)
	}
	return axis, indices
}

type theaterMask struct {
	prefix []int
	stride int
	xs, ys map[int64]int
}

func theaterMaskSize(columns, rows int) (int, error) {
	const maxEntries = 4 * 1024 * 1024
	if columns < 1 || rows < 1 || columns >= maxEntries || rows >= maxEntries || columns+1 > maxEntries/(rows+1) {
		return 0, errors.New("compressed polygon exceeds the prefix table budget")
	}
	return (columns + 1) * (rows + 1), nil
}

func buildTheaterMask(points []theaterPoint) (*theaterMask, error) {
	if err := validateTheaterPolygon(points); err != nil {
		return nil, err
	}
	xValues, yValues := make([]int64, len(points)), make([]int64, len(points))
	for i, p := range points {
		xValues[i], yValues[i] = p.x, p.y
	}
	xs, xIndex := theaterAxis(xValues)
	ys, yIndex := theaterAxis(yValues)
	entries, err := theaterMaskSize(len(xs), len(ys))
	if err != nil {
		return nil, err
	}
	mask := &theaterMask{prefix: make([]int, entries), stride: len(xs) + 1, xs: xIndex, ys: yIndex}
	inside := make([]bool, len(xs))
	crossings := make([]int, 0, len(points))
	for row, y := range ys {
		clear(inside)
		crossings = crossings[:0]
		for i, a := range points {
			b := points[(i+1)%len(points)]
			if a.x == b.x && y >= min(a.y, b.y) && y < max(a.y, b.y) {
				crossings = append(crossings, xIndex[a.x])
			}
		}
		slices.Sort(crossings)
		if len(crossings)%2 != 0 {
			return nil, errors.New("polygon has an unmatched scanline crossing")
		}
		for i := 0; i < len(crossings); i += 2 {
			for column := crossings[i]; column <= crossings[i+1]; column++ {
				inside[column] = true
			}
		}
		for i, a := range points {
			b := points[(i+1)%len(points)]
			if a.y == b.y && a.y == y {
				for column := min(xIndex[a.x], xIndex[b.x]); column <= max(xIndex[a.x], xIndex[b.x]); column++ {
					inside[column] = true
				}
			} else if a.x == b.x && y >= min(a.y, b.y) && y <= max(a.y, b.y) {
				inside[xIndex[a.x]] = true
			}
		}
		outside := 0
		for column, on := range inside {
			if !on {
				outside++
			}
			mask.prefix[(row+1)*mask.stride+column+1] = mask.prefix[row*mask.stride+column+1] + outside
		}
	}
	return mask, nil
}

func (m *theaterMask) contains(a, b theaterPoint) bool {
	x1, x2 := min(m.xs[a.x], m.xs[b.x]), max(m.xs[a.x], m.xs[b.x])+1
	y1, y2 := min(m.ys[a.y], m.ys[b.y]), max(m.ys[a.y], m.ys[b.y])+1
	right := m.prefix[y2*m.stride+x2] - m.prefix[y1*m.stride+x2]
	left := m.prefix[y2*m.stride+x1] - m.prefix[y1*m.stride+x1]
	return right == left
}

func (d Day9) Part1(input string) (string, error) {
	points, err := parseTheaterPoints(input)
	if err != nil {
		return "", err
	}
	var best int64
	for i, a := range points {
		for _, b := range points[i+1:] {
			area, err := theaterArea(a, b)
			if err != nil {
				return "", err
			}
			best = max(best, area)
		}
	}
	return strconv.FormatInt(best, 10), nil
}

func (d Day9) Part2(input string) (string, error) {
	points, err := parseTheaterPoints(input)
	if err != nil {
		return "", err
	}
	mask, err := buildTheaterMask(points)
	if err != nil {
		return "", err
	}
	var best int64
	for i, a := range points {
		for _, b := range points[i+1:] {
			area, err := theaterArea(a, b)
			if err != nil {
				return "", err
			}
			if area > best && mask.contains(a, b) {
				best = area
			}
		}
	}
	return strconv.FormatInt(best, 10), nil
}

func init() { solve.Register(Day9{}) }
