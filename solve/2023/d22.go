package solve2023

import (
	"aoc/solve"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

type Day22 struct{}

func (d Day22) Coords() solve.SolutionCoords { return solve.SolutionCoords{Year: 2023, Day: 22} }

type Point3D struct{ X, Y, Z int }
type Brick struct {
	Start, End            Point3D
	Supports, SupportedBy []int
}

func parsePoint3D(s string) (Point3D, error) {
	fields := strings.Split(s, ",")
	if len(fields) != 3 {
		return Point3D{}, fmt.Errorf("expected three coordinates")
	}
	var values [3]int
	for i, field := range fields {
		if field == "" {
			return Point3D{}, fmt.Errorf("empty coordinate")
		}
		for _, c := range field {
			if c < '0' || c > '9' {
				return Point3D{}, fmt.Errorf("coordinates must be nonnegative decimal integers")
			}
		}
		value, err := strconv.Atoi(field)
		if err != nil {
			return Point3D{}, err
		}
		values[i] = value
	}
	return Point3D{values[0], values[1], values[2]}, nil
}
func overlapsXY(a, b Brick) bool {
	return a.Start.X <= b.End.X && b.Start.X <= a.End.X && a.Start.Y <= b.End.Y && b.Start.Y <= a.End.Y
}
func parseBricks(input string) ([]Brick, error) {
	if strings.TrimSpace(input) == "" {
		return nil, fmt.Errorf("puzzle input is empty")
	}
	input = strings.ReplaceAll(input, "\r\n", "\n")
	lines := strings.Split(strings.TrimSuffix(input, "\n"), "\n")
	bricks := make([]Brick, 0, len(lines))
	for i, line := range lines {
		parse := func() (Brick, error) {
			fields := strings.Split(line, "~")
			if len(fields) != 2 {
				return Brick{}, fmt.Errorf("expected two endpoints separated by ~")
			}
			a, err := parsePoint3D(fields[0])
			if err != nil {
				return Brick{}, err
			}
			b, err := parsePoint3D(fields[1])
			if err != nil {
				return Brick{}, err
			}
			if a.Z == 0 || b.Z == 0 {
				return Brick{}, fmt.Errorf("brick must be above ground (z >= 1)")
			}
			axes := 0
			if a.X != b.X {
				axes++
			}
			if a.Y != b.Y {
				axes++
			}
			if a.Z != b.Z {
				axes++
			}
			if axes > 1 {
				return Brick{}, fmt.Errorf("brick must be a straight axis-aligned line or a single cube")
			}
			return Brick{Start: Point3D{min(a.X, b.X), min(a.Y, b.Y), min(a.Z, b.Z)}, End: Point3D{max(a.X, b.X), max(a.Y, b.Y), max(a.Z, b.Z)}}, nil
		}
		brick, err := parse()
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", i+1, err)
		}
		bricks = append(bricks, brick)
	}
	for i, brick := range bricks {
		for j, previous := range bricks[:i] {
			if brick.Start.Z <= previous.End.Z && previous.Start.Z <= brick.End.Z && overlapsXY(brick, previous) {
				return nil, fmt.Errorf("line %d: brick overlaps the brick on line %d before falling", i+1, j+1)
			}
		}
	}
	sort.SliceStable(bricks, func(i, j int) bool { return bricks[i].Start.Z < bricks[j].Start.Z })
	for i := range bricks {
		brick := &bricks[i]
		height := brick.End.Z - brick.Start.Z
		highest := 0
		for j := i - 1; j >= 0; j-- {
			support := &bricks[j]
			if support.End.Z < highest || !overlapsXY(*brick, *support) {
				continue
			}
			if support.End.Z > highest {
				highest = support.End.Z
				brick.SupportedBy = brick.SupportedBy[:0]
			}
			brick.SupportedBy = append(brick.SupportedBy, j)
		}
		if highest == math.MaxInt || height > math.MaxInt-(highest+1) {
			return nil, fmt.Errorf("landing height overflow")
		}
		brick.Start.Z = highest + 1
		brick.End.Z = brick.Start.Z + height
		for _, parent := range brick.SupportedBy {
			bricks[parent].Supports = append(bricks[parent].Supports, i)
		}
	}
	return bricks, nil
}
func (d Day22) Part1(input string) (string, error) {
	bricks, err := parseBricks(input)
	if err != nil {
		return "", err
	}
	count := 0
	for _, brick := range bricks {
		safe := true
		for _, child := range brick.Supports {
			if len(bricks[child].SupportedBy) <= 1 {
				safe = false
				break
			}
		}
		if safe {
			count++
		}
	}
	return strconv.Itoa(count), nil
}

// Each support edge is removed once; a child is queued when its last support falls.
func cascadeWithScratch(bricks []Brick, removed int, remaining, queue []int) int {
	for i, brick := range bricks {
		remaining[i] = len(brick.SupportedBy)
	}
	queue = append(queue[:0], removed)
	for cursor := 0; cursor < len(queue); cursor++ {
		for _, child := range bricks[queue[cursor]].Supports {
			remaining[child]--
			if remaining[child] == 0 {
				queue = append(queue, child)
			}
		}
	}
	return len(queue) - 1
}
func (d Day22) Part2(input string) (string, error) {
	bricks, err := parseBricks(input)
	if err != nil {
		return "", err
	}
	sum := 0
	remaining := make([]int, len(bricks))
	queue := make([]int, 0, len(bricks))
	for i, brick := range bricks {
		if len(brick.Supports) == 0 {
			continue
		}
		fallen := cascadeWithScratch(bricks, i, remaining, queue)
		if fallen > math.MaxInt-sum {
			return "", fmt.Errorf("cascade total overflow")
		}
		sum += fallen
	}
	return strconv.Itoa(sum), nil
}
func init() { solve.Register(Day22{}) }
