package solve2016

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day17 struct {
}

func (d Day17) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2016, Day: 17}
}

var moves = map[byte]func(x, y int) (int, int){
	'U': func(x, y int) (int, int) { return x, y - 1 },
	'D': func(x, y int) (int, int) { return x, y + 1 },
	'L': func(x, y int) (int, int) { return x - 1, y },
	'R': func(x, y int) (int, int) { return x + 1, y },
}

// doors determines which doors are open based on the MD5 hash of the passcode and path.
func (d Day17) doors(path, passcode string) []bool {
	hash := md5.Sum([]byte(passcode + path))
	hexHash := hex.EncodeToString(hash[:])
	open := make([]bool, 4)
	for i, c := range hexHash[:4] {
		open[i] = c >= 'b' && c <= 'f'
	}
	return open
}

// bfs performs a breadth-first search to find paths to the goal.
func (d Day17) bfs(start, goal [2]int, passcode string) (string, int) {
	type state struct {
		pos  [2]int
		path string
	}
	queue := []state{{start, ""}}
	var shortest string
	var longest int

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		x, y := current.pos[0], current.pos[1]
		if current.pos == goal {
			if shortest == "" {
				shortest = current.path
			}
			if len(current.path) > longest {
				longest = len(current.path)
			}
			continue
		}

		openDoors := d.doors(current.path, passcode)
		directions := "UDLR"
		for i, open := range openDoors {
			if !open {
				continue
			}
			dir := directions[i]
			nx, ny := moves[dir](x, y)
			if nx >= 0 && nx < 4 && ny >= 0 && ny < 4 {
				queue = append(queue, state{[2]int{nx, ny}, current.path + string(dir)})
			}
		}
	}

	return shortest, longest
}

func (d Day17) Part1(input string) (string, error) {
	passcode := strings.TrimSpace(input)
	shortest, _ := d.bfs([2]int{0, 0}, [2]int{3, 3}, passcode)
	return shortest, nil
}

func (d Day17) Part2(input string) (string, error) {
	passcode := strings.TrimSpace(input)
	_, longest := d.bfs([2]int{0, 0}, [2]int{3, 3}, passcode)
	return strconv.Itoa(longest), nil
}

func init() {
	solve.Register(Day17{})
}
