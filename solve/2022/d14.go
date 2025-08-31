package solve2022

import (
	"aoc/solve"
	"strconv"
	"strings"
)

type Day14 struct{}

func (d Day14) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2022, Day: 14}
}

func (d Day14) parseInput(input string) [][][2]int {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var paths [][][2]int
	for _, line := range lines {
		var path [][2]int
		segments := strings.Split(line, "->")
		for _, segment := range segments {
			coords := strings.Split(strings.TrimSpace(segment), ",")
			x, _ := strconv.Atoi(coords[0])
			y, _ := strconv.Atoi(coords[1])
			path = append(path, [2]int{x, y})
		}
		paths = append(paths, path)
	}
	return paths
}

func putPaths(grid map[[2]int]rune, paths [][][2]int) {
	for _, path := range paths {
		for i := 0; i < len(path)-1; i++ {
			x0, y0 := path[i][0], path[i][1]
			x1, y1 := path[i+1][0], path[i+1][1]
			for x := min(x0, x1); x <= max(x0, x1); x++ {
				for y := min(y0, y1); y <= max(y0, y1); y++ {
					grid[[2]int{x, y}] = '#'
				}
			}
		}
	}
}

func simulateSand(grid map[[2]int]rune, paths [][][2]int, entry [2]int, floor int) int {
	putPaths(grid, paths)
	bottom := 0
	for coord := range grid {
		if coord[1] > bottom {
			bottom = coord[1]
		}
	}
	falling := [][2]int{{0, 1}, {-1, 1}, {1, 1}}
	if floor > 0 {
		dx := bottom + 2
		bottom += floor
		putPaths(grid, [][][2]int{{{entry[0] - dx, bottom}, {entry[0] + dx, bottom}}})
	}
	particles := 0
	for {
		particles++
		loc := entry
		for {
			stopped := true
			for _, dir := range falling {
				next := [2]int{loc[0] + dir[0], loc[1] + dir[1]}
				if _, exists := grid[next]; !exists {
					loc = next
					stopped = false
					break
				}
			}
			if loc == entry {
				return particles
			}
			if stopped {
				grid[loc] = 'o'
				break
			}
			if loc[1] > bottom {
				grid[loc] = 'o'
				return particles - 1
			}
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (d Day14) Part1(input string) (string, error) {
	paths := d.parseInput(input)
	grid := make(map[[2]int]rune)
	result := simulateSand(grid, paths, [2]int{500, 0}, 0)
	return strconv.Itoa(result), nil
}

func (d Day14) Part2(input string) (string, error) {
	paths := d.parseInput(input)
	grid := make(map[[2]int]rune)
	result := simulateSand(grid, paths, [2]int{500, 0}, 2)
	return strconv.Itoa(result), nil
}

func init() {
	solve.Register(Day14{})
}
