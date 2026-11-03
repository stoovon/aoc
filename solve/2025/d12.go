package solve2025

import (
	"cmp"
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day12 struct{}
type day12Point struct{ x, y int }
type day12Shape struct {
	cells        []day12Point
	orientations [][]day12Point
}
type day12Region struct {
	w, h   int
	counts []int
	line   int
}

func (d Day12) Coords() solve.SolutionCoords { return solve.SolutionCoords{Year: 2025, Day: 12} }

func day12Number(text string) (int, error) {
	if text == "" {
		return 0, errors.New("expected a nonnegative decimal integer")
	}
	for _, b := range []byte(text) {
		if b < '0' || b > '9' {
			return 0, errors.New("expected a nonnegative decimal integer")
		}
	}
	return strconv.Atoi(text)
}

func parseDay12(input string) ([]day12Shape, []day12Region, error) {
	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	shapesByIndex := make(map[int]day12Shape)
	var regions []day12Region
	regionPhase := false
	for index := 0; index < len(lines); {
		line := strings.TrimSpace(lines[index])
		lineNumber := index + 1
		index++
		if line == "" {
			continue
		}
		label, rest, ok := strings.Cut(line, ":")
		if !ok {
			return nil, nil, fmt.Errorf("line %d: expected a shape header or region", lineNumber)
		}
		label = strings.TrimSpace(label)
		if strings.Contains(label, "x") || regionPhase {
			regionPhase = true
			wText, hText, ok := strings.Cut(label, "x")
			if !ok {
				return nil, nil, fmt.Errorf("line %d: expected widthxheight", lineNumber)
			}
			w, err := day12Number(strings.TrimSpace(wText))
			if err != nil {
				return nil, nil, fmt.Errorf("line %d: width: %w", lineNumber, err)
			}
			h, err := day12Number(strings.TrimSpace(hText))
			if err != nil {
				return nil, nil, fmt.Errorf("line %d: height: %w", lineNumber, err)
			}
			if w == 0 || h == 0 {
				return nil, nil, fmt.Errorf("line %d: region dimensions must be positive", lineNumber)
			}
			fields := strings.Fields(rest)
			if len(fields) != len(shapesByIndex) {
				return nil, nil, fmt.Errorf("line %d: expected one count per shape", lineNumber)
			}
			counts := make([]int, len(fields))
			for i, field := range fields {
				value, err := day12Number(field)
				if err != nil {
					return nil, nil, fmt.Errorf("line %d: count: %w", lineNumber, err)
				}
				counts[i] = value
			}
			regions = append(regions, day12Region{w, h, counts, lineNumber})
			continue
		}
		shapeIndex, err := day12Number(label)
		if err != nil || strings.TrimSpace(rest) != "" {
			return nil, nil, fmt.Errorf("line %d: expected a numbered shape header", lineNumber)
		}
		if _, exists := shapesByIndex[shapeIndex]; exists {
			return nil, nil, fmt.Errorf("line %d: duplicate shape index", lineNumber)
		}
		var cells []day12Point
		width, height := 0, 0
		for index < len(lines) && strings.TrimSpace(lines[index]) != "" {
			row := strings.TrimSpace(lines[index])
			index++
			if width == 0 {
				width = len(row)
			}
			if len(row) != width {
				return nil, nil, fmt.Errorf("line %d: shape rows must have equal width", index)
			}
			for x, b := range []byte(row) {
				switch b {
				case '#':
					cells = append(cells, day12Point{x, height})
				case '.':
				default:
					return nil, nil, fmt.Errorf("line %d: shape rows may contain only '.' and '#'", index)
				}
			}
			height++
		}
		if len(cells) == 0 {
			return nil, nil, fmt.Errorf("line %d: shape must contain an occupied cell", lineNumber)
		}
		cells = day12NormalizePoints(cells)
		shapesByIndex[shapeIndex] = day12Shape{cells, day12GenerateOrientations(cells)}
	}
	if len(shapesByIndex) == 0 || len(regions) == 0 {
		return nil, nil, errors.New("expected shapes and at least one region")
	}
	shapes := make([]day12Shape, len(shapesByIndex))
	for i := range shapes {
		shape, ok := shapesByIndex[i]
		if !ok {
			return nil, nil, fmt.Errorf("missing shape index %d", i)
		}
		shapes[i] = shape
	}
	return shapes, regions, nil
}

func day12NormalizePoints(points []day12Point) []day12Point {
	minX, minY := points[0].x, points[0].y
	for _, p := range points {
		minX = min(minX, p.x)
		minY = min(minY, p.y)
	}
	for i := range points {
		points[i].x -= minX
		points[i].y -= minY
	}
	slices.SortFunc(points, func(a, b day12Point) int {
		if a.y != b.y {
			return cmp.Compare(a.y, b.y)
		}
		return cmp.Compare(a.x, b.x)
	})
	return points
}

func day12GenerateOrientations(cells []day12Point) [][]day12Point {
	seen := make(map[string]bool)
	var orientations [][]day12Point
	for rotation := 0; rotation < 4; rotation++ {
		for reflection := 0; reflection < 2; reflection++ {
			transformed := make([]day12Point, len(cells))
			for i, p := range cells {
				for turn := 0; turn < rotation; turn++ {
					p.x, p.y = -p.y, p.x
				}
				if reflection == 1 {
					p.x = -p.x
				}
				transformed[i] = p
			}
			transformed = day12NormalizePoints(transformed)
			var key strings.Builder
			for _, p := range transformed {
				fmt.Fprintf(&key, "%d,%d;", p.x, p.y)
			}
			encoded := key.String()
			if !seen[encoded] {
				seen[encoded] = true
				orientations = append(orientations, transformed)
			}
		}
	}
	return orientations
}

func day12Product(a, b int) (int, error) {
	if a < 0 || b < 0 || (b != 0 && a > math.MaxInt/b) {
		return 0, errors.New("packing arithmetic exceeds the integer range")
	}
	return a * b, nil
}

func day12Bounds(cells []day12Point) (int, int) {
	width, height := 0, 0
	for _, p := range cells {
		width = max(width, p.x+1)
		height = max(height, p.y+1)
	}
	return width, height
}

func canFitRegion(region day12Region, shapes []day12Shape) (bool, error) {
	area, err := day12Product(region.w, region.h)
	if err != nil {
		return false, err
	}
	needed, pieces, boxWidth, boxHeight := 0, 0, 1, 1
	for i, count := range region.counts {
		cells, err := day12Product(count, len(shapes[i].cells))
		if err != nil {
			return false, err
		}
		if needed > math.MaxInt-cells {
			return false, errors.New("total occupied area exceeds the integer range")
		}
		needed += cells
		if pieces > math.MaxInt-count {
			return false, errors.New("piece count exceeds the integer range")
		}
		pieces += count
		if count > 0 {
			w, h := day12Bounds(shapes[i].cells)
			boxWidth = max(boxWidth, w)
			boxHeight = max(boxHeight, h)
		}
	}
	if needed > area {
		return false, nil
	}
	if pieces <= (region.w/boxWidth)*(region.h/boxHeight) || pieces <= (region.w/boxHeight)*(region.h/boxWidth) {
		return true, nil
	}
	return day12ExactFit(region, shapes, 2000000)
}

type day12PackingType struct {
	count      int
	placements [][]int
	minimum    int
}

func day12ExactFit(region day12Region, shapes []day12Shape, nodeBudget int) (bool, error) {
	area, err := day12Product(region.w, region.h)
	if err != nil {
		return false, err
	}
	if area > 4096 {
		return false, errors.New("undecided region exceeds the exact packing board budget")
	}
	var types []day12PackingType
	placementCells := 0
	for i, count := range region.counts {
		if count == 0 {
			continue
		}
		piece := day12PackingType{count: count}
		for _, orientation := range shapes[i].orientations {
			width, height := day12Bounds(orientation)
			for y := 0; y <= region.h-height; y++ {
				for x := 0; x <= region.w-width; x++ {
					if placementCells > 1000000-len(orientation) {
						return false, errors.New("exact packing placements exceed the storage budget")
					}
					placement := make([]int, len(orientation))
					for j, p := range orientation {
						placement[j] = (y+p.y)*region.w + x + p.x
					}
					piece.placements = append(piece.placements, placement)
					placementCells += len(orientation)
				}
			}
		}
		if len(piece.placements) < count {
			return false, nil
		}
		types = append(types, piece)
	}
	occupied := make([]bool, area)
	fits := func(placement []int) bool {
		for _, cell := range placement {
			if occupied[cell] {
				return false
			}
		}
		return true
	}
	var search func() (bool, error)
	search = func() (bool, error) {
		chosen, options := -1, math.MaxInt
		for i, piece := range types {
			if piece.count == 0 {
				continue
			}
			available := 0
			for _, placement := range piece.placements[piece.minimum:] {
				if fits(placement) {
					available++
				}
			}
			if available < piece.count {
				return false, nil
			}
			if available < options {
				chosen, options = i, available
			}
		}
		if chosen < 0 {
			return true, nil
		}
		if nodeBudget <= 0 {
			return false, errors.New("exact packing search exceeds the node budget")
		}
		nodeBudget--
		piece := &types[chosen]
		minimum := piece.minimum
		piece.count--
		for index := minimum; index < len(piece.placements); index++ {
			placement := piece.placements[index]
			if !fits(placement) {
				continue
			}
			for _, cell := range placement {
				occupied[cell] = true
			}
			piece.minimum = index + 1
			result, err := search()
			for _, cell := range placement {
				occupied[cell] = false
			}
			piece.minimum = minimum
			if err != nil || result {
				piece.count++
				return result, err
			}
		}
		piece.count++
		return false, nil
	}
	return search()
}

func (d Day12) Part1(input string) (string, error) {
	shapes, regions, err := parseDay12(input)
	if err != nil {
		return "", err
	}
	count := 0
	for _, region := range regions {
		fits, err := canFitRegion(region, shapes)
		if err != nil {
			return "", fmt.Errorf("line %d: %w", region.line, err)
		}
		if fits {
			count++
		}
	}
	return strconv.Itoa(count), nil
}
func (d Day12) Part2(input string) (string, error) { return "", errors.New("not implemented") }

func init() { solve.Register(Day12{}) }
