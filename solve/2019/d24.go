package solve2019

import (
	"bufio"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day24 struct {
}

func (d Day24) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2019, Day: 24}
}

// Converts grid to a 25-bit integer
func parseGrid(input string) uint32 {
	var state uint32
	scanner := bufio.NewScanner(strings.NewReader(input))
	i := 0
	for scanner.Scan() {
		for _, c := range scanner.Text() {
			if c == '#' {
				state |= 1 << i
			}
			i++
		}
	}
	return state
}

// Counts adjacent bugs for a cell in a flat grid (Part 1)
func countAdjLevelsFlat(levels map[int]uint32, depth, i int) int {
	state := levels[depth]
	x, y := i%5, i/5
	adj := 0
	for _, d := range [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}} {
		nx, ny := x+d[0], y+d[1]
		if nx >= 0 && nx < 5 && ny >= 0 && ny < 5 {
			ni := ny*5 + nx
			if (state & (1 << ni)) != 0 {
				adj++
			}
		}
	}
	return adj
}

func countAdjBugsRecursive(levels map[int]uint32, depth, i int) int {
	if i == 12 {
		return 0 // center cell is always empty
	}
	x, y := i%5, i/5
	adj := 0
	for _, dxy := range [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}} {
		nx, ny := x+dxy[0], y+dxy[1]
		if nx == 2 && ny == 2 {
			// Into inner level
			inner := levels[depth+1]
			switch {
			case x == 2 && y == 1: // down into top row
				for k := 0; k < 5; k++ {
					if (inner & (1 << k)) != 0 {
						adj++
					}
				}
			case x == 2 && y == 3: // up into bottom row
				for k := 20; k < 25; k++ {
					if (inner & (1 << k)) != 0 {
						adj++
					}
				}
			case x == 1 && y == 2: // right into left col
				for k := 0; k < 25; k += 5 {
					if (inner & (1 << k)) != 0 {
						adj++
					}
				}
			case x == 3 && y == 2: // left into right col
				for k := 4; k < 25; k += 5 {
					if (inner & (1 << k)) != 0 {
						adj++
					}
				}
			}
		} else if nx < 0 || nx > 4 || ny < 0 || ny > 4 {
			// Into outer level
			outer := levels[depth-1]
			switch {
			case nx < 0:
				if (outer & (1 << 11)) != 0 {
					adj++
				}
			case nx > 4:
				if (outer & (1 << 13)) != 0 {
					adj++
				}
			case ny < 0:
				if (outer & (1 << 7)) != 0 {
					adj++
				}
			default: // ny > 4
				if (outer & (1 << 17)) != 0 {
					adj++
				}
			}
		} else {
			// Normal neighbour
			if (levels[depth] & (1 << (ny*5 + nx))) != 0 {
				adj++
			}
		}
	}
	return adj
}

func stepLevels(
	levels map[int]uint32,
	countAdj func(map[int]uint32, int, int) int,
	skipCenter bool,
) map[int]uint32 {
	// Find min/max levels
	minLevel, maxLevel := 0, 0
	for k := range levels {
		if k < minLevel {
			minLevel = k
		}
		if k > maxLevel {
			maxLevel = k
		}
	}
	levels[minLevel-1] = 0
	levels[maxLevel+1] = 0

	next := map[int]uint32{}
	for depth := minLevel - 1; depth <= maxLevel+1; depth++ {
		var newState uint32
		for i := 0; i < 25; i++ {
			if skipCenter && i == 12 {
				continue
			}
			adj := countAdj(levels, depth, i)
			if (levels[depth] & (1 << i)) != 0 {
				if adj == 1 {
					newState |= 1 << i
				}
			} else {
				if adj == 1 || adj == 2 {
					newState |= 1 << i
				}
			}
		}
		next[depth] = newState
	}
	return next
}

func (d Day24) Part1(input string) (string, error) {
	levels := map[int]uint32{0: parseGrid(input)}
	seen := map[uint32]bool{}
	for {
		state := levels[0]
		if seen[state] {
			return strconv.Itoa(int(state)), nil
		}
		seen[state] = true
		levels = stepLevels(levels, countAdjLevelsFlat, false)
	}
}

func (d Day24) Part2(input string) (string, error) {
	const minutes = 200
	levels := map[int]uint32{0: parseGrid(input)}
	for t := 0; t < minutes; t++ {
		levels = stepLevels(levels, countAdjBugsRecursive, true)
	}
	// Count bugs
	count := 0
	for _, state := range levels {
		for i := 0; i < 25; i++ {
			if i == 12 {
				continue
			}
			if (state & (1 << i)) != 0 {
				count++
			}
		}
	}
	return strconv.Itoa(count), nil
}

func init() {
	solve.Register(Day24{})
}
