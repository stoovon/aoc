package solve2023

import (
	"sort"
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/maths"
)

type Day22 struct {
}

func (d Day22) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2023, Day: 22}
}

type Point3D struct {
	X, Y, Z int
}

type Brick struct {
	Start, End  Point3D
	Supports    []int
	SupportedBy []int
}

func parsePoint3D(s string) Point3D {
	parts := strings.Split(s, ",")
	x, _ := strconv.Atoi(parts[0])
	y, _ := strconv.Atoi(parts[1])
	z, _ := strconv.Atoi(parts[2])
	return Point3D{x, y, z}
}

func pointBetween(p int, rng [2]int) bool {
	return p >= rng[0] && p <= rng[1]
}

func hasCollision(a, b Brick) bool {
	for d := 0; d < 3; d++ {
		var a0, a1, b0, b1 int
		switch d {
		case 0:
			a0, a1 = min(a.Start.X, a.End.X), max(a.Start.X, a.End.X)
			b0, b1 = min(b.Start.X, b.End.X), max(b.Start.X, b.End.X)
		case 1:
			a0, a1 = min(a.Start.Y, a.End.Y), max(a.Start.Y, a.End.Y)
			b0, b1 = min(b.Start.Y, b.End.Y), max(b.Start.Y, b.End.Y)
		case 2:
			a0, a1 = min(a.Start.Z, a.End.Z), max(a.Start.Z, a.End.Z)
			b0, b1 = min(b.Start.Z, b.End.Z), max(b.Start.Z, b.End.Z)
		}
		if !(pointBetween(a0, [2]int{b0, b1}) || pointBetween(a1, [2]int{b0, b1}) ||
			pointBetween(b0, [2]int{a0, a1}) || pointBetween(b1, [2]int{a0, a1})) {
			return false
		}
	}
	return true
}

func drop(brick *Brick) {
	if brick.Start.Z > 1 {
		brick.Start.Z--
		brick.End.Z--
	}
}

func raise(brick *Brick) {
	brick.Start.Z++
	brick.End.Z++
}

func parseBricks(input string) []*Brick {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	bricks := make([]*Brick, 0, len(lines))
	for _, line := range lines {
		parts := strings.Split(line, "~")
		bricks = append(bricks, &Brick{
			Start: parsePoint3D(parts[0]),
			End:   parsePoint3D(parts[1]),
		})
	}
	sort.Slice(bricks, func(i, j int) bool {
		return min(bricks[i].Start.Z, bricks[i].End.Z) < min(bricks[j].Start.Z, bricks[j].End.Z)
	})
	result := make([]*Brick, 0, len(bricks))
	maxZReached := 1
	for _, orig := range bricks {
		brick := *orig // copy value
		diff := maths.Abs(brick.End.Z - brick.Start.Z)
		brick.Start.Z = maxZReached + 1
		brick.End.Z = brick.Start.Z + diff
		newIdx := len(result)
		for {
			canMoveDown := true
			drop(&brick)
			for i := len(result) - 1; i >= 0; i-- {
				if max(result[i].Start.Z, result[i].End.Z) < min(brick.Start.Z, brick.End.Z) {
					continue
				}
				if hasCollision(brick, *result[i]) {
					canMoveDown = false
					result[i].Supports = append(result[i].Supports, newIdx)
					brick.SupportedBy = append(brick.SupportedBy, i)
				}
			}
			if !canMoveDown {
				raise(&brick)
				break
			}
			if brick.Start.Z == 1 {
				break
			}
		}
		maxZReached = max(maxZReached, brick.End.Z)
		// append pointer to new brick
		result = append(result, &brick)
	}
	return result
}

func (d Day22) Part1(input string) (string, error) {
	bricks := parseBricks(input)
	count := 0
	for i := range bricks {
		ok := true
		for _, idx := range bricks[i].Supports {
			if len(bricks[idx].SupportedBy) <= 1 {
				ok = false
				break
			}
		}
		if ok {
			count++
		}
	}
	return strconv.Itoa(count), nil
}

func (d Day22) Part2(input string) (string, error) {
	bricks := parseBricks(input)
	sum := 0
	for i := len(bricks) - 1; i >= 0; i-- {
		if len(bricks[i].Supports) == 0 {
			continue
		}
		unsupported := make(map[int]bool)
		queue := []int{i}
		for len(queue) > 0 {
			idx := queue[0]
			queue = queue[1:]
			if unsupported[idx] {
				continue
			}
			unsupported[idx] = true
			for _, sidx := range bricks[idx].Supports {
				all := true
				for _, sup := range bricks[sidx].SupportedBy {
					if !unsupported[sup] {
						all = false
						break
					}
				}
				if all {
					queue = append(queue, sidx)
				}
			}
		}
		sum += len(unsupported) - 1
	}
	return strconv.Itoa(sum), nil
}

func init() {
	solve.Register(Day22{})
}
