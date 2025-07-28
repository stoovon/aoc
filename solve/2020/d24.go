package solve2020

import (
	"aoc/solve"
	"strings"
	"strconv"
)

type Day24 struct{}

func (d Day24) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2020, Day: 24}
}

func (d Day24) Part1(input string) (string, error) {
	black := initialBlackTiles(input)
	return strconv.Itoa(len(black)), nil
}

func (d Day24) Part2(input string) (string, error) {
	black := initialBlackTiles(input)
	dirs := hexDirs()
	for day := 0; day < 100; day++ {
		neighborCount := make(map[[2]int]int)
		for pos := range black {
			for _, d := range dirs {
				n := [2]int{pos[0]+d[0], pos[1]+d[1]}
				neighborCount[n]++
			}
		}
		next := make(map[[2]int]bool)
		for pos, cnt := range neighborCount {
			if black[pos] {
				if cnt == 1 || cnt == 2 {
					next[pos] = true
				}
			} else {
				if cnt == 2 {
					next[pos] = true
				}
			}
		}
		black = next
	}
	return strconv.Itoa(len(black)), nil
}
// initialBlackTiles parses the input and returns the set of black tiles after all flips.
func initialBlackTiles(input string) map[[2]int]bool {
	lines := strings.Fields(strings.TrimSpace(input))
	dirs := hexDirs()
	dirMap := map[string][2]int{
		"e":  dirs[0], "se": dirs[1], "sw": dirs[2],
		"w":  dirs[3], "nw": dirs[4], "ne": dirs[5],
	}
	black := make(map[[2]int]bool)
	for _, line := range lines {
		pos := [2]int{0, 0}
		for i := 0; i < len(line); {
			var d string
			if line[i] == 'e' || line[i] == 'w' {
				d = line[i : i+1]
				i++
			} else {
				d = line[i : i+2]
				i += 2
			}
			delta := dirMap[d]
			pos[0] += delta[0]
			pos[1] += delta[1]
		}
		black[pos] = !black[pos]
	}
	for k, v := range black {
		if !v {
			delete(black, k)
		}
	}
	return black
}

// hexDirs returns the 6 axial direction vectors for a hex grid.
func hexDirs() [][2]int {
	return [][2]int{
		{1, 0},  // e
		{0, 1},  // se
		{-1, 1}, // sw
		{-1, 0}, // w
		{0, -1}, // nw
		{1, -1}, // ne
	}
}

func init() {
	solve.Register(Day24{})
}
