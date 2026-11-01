package solve2024

import (
	"aoc/solve"
	"fmt"
	"strconv"
	"strings"
)

type Day6 struct{}

func (d Day6) Coords() solve.SolutionCoords { return solve.SolutionCoords{Year: 2024, Day: 6} }

type patrolGrid struct {
	cells []byte
	start int
	ahead []int
}

func (d Day6) parseGrid(input string) (patrolGrid, error) {
	var grid patrolGrid
	text := strings.TrimSuffix(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	lines := strings.Split(text, "\n")
	width, height := len(lines[0]), len(lines)
	if width == 0 {
		return grid, fmt.Errorf("empty patrol grid")
	}
	maxInt := int(^uint(0) >> 1)
	if height > maxInt/4/width {
		return grid, fmt.Errorf("patrol state count exceeds int")
	}
	starts := 0
	for y, row := range lines {
		if len(row) != width {
			return grid, fmt.Errorf("row %d: inconsistent width", y+1)
		}
		for x := 0; x < len(row); x++ {
			switch row[x] {
			case '.', '#':
			case '^':
				grid.start = y*width + x
				starts++
			default:
				return grid, fmt.Errorf("row %d: expected '.', '#' or '^'", y+1)
			}
		}
	}
	if starts != 1 {
		return grid, fmt.Errorf("expected exactly one guard start")
	}
	grid.cells = make([]byte, 0, width*height)
	for _, row := range lines {
		grid.cells = append(grid.cells, row...)
	}
	grid.ahead = make([]int, 4*len(grid.cells))
	for cell := range grid.cells {
		x, y := cell%width, cell/width
		for direction := 0; direction < 4; direction++ {
			nx, ny := x, y
			switch direction {
			case 0:
				ny--
			case 1:
				nx++
			case 2:
				ny++
			case 3:
				nx--
			}
			next := -1
			if nx >= 0 && nx < width && ny >= 0 && ny < height {
				next = ny*width + nx
			}
			grid.ahead[cell*4+direction] = next
		}
	}
	return grid, nil
}

func (grid patrolGrid) next(state, obstacle int) int {
	cell, direction := state/4, state%4
	ahead := grid.ahead[state]
	if ahead < 0 {
		return -1
	}
	if grid.cells[ahead] == '#' || ahead == obstacle {
		return cell*4 + (direction+1)%4
	}
	return ahead*4 + direction
}

func (grid patrolGrid) trace() ([]uint8, bool) {
	visited := make([]uint8, len(grid.cells))
	for state := grid.start * 4; state >= 0; state = grid.next(state, -1) {
		cell, direction := state/4, state%4
		bit := uint8(1) << direction
		if visited[cell]&bit != 0 {
			return visited, true
		}
		visited[cell] |= bit
	}
	return visited, false
}

func (grid patrolGrid) loops(obstacle int, marks []int, trial int) bool {
	for state := grid.start * 4; state >= 0; state = grid.next(state, obstacle) {
		if marks[state] == trial {
			return true
		}
		marks[state] = trial
	}
	return false
}

func (d Day6) Part1(input string) (string, error) {
	grid, err := d.parseGrid(input)
	if err != nil {
		return "", err
	}
	visited, _ := grid.trace()
	count := 0
	for _, directions := range visited {
		if directions != 0 {
			count++
		}
	}
	return strconv.Itoa(count), nil
}

func (d Day6) Part2(input string) (string, error) {
	grid, err := d.parseGrid(input)
	if err != nil {
		return "", err
	}
	visited, baseLoops := grid.trace()
	marks := make([]int, len(grid.ahead))
	count, trial := 0, 0
	for cell, content := range grid.cells {
		if content != '.' {
			continue
		}
		if visited[cell] == 0 {
			if baseLoops {
				count++
			}
			continue
		}
		// At most one trial per cell; parseGrid bounds 4*cellCount within int.
		trial++
		if grid.loops(cell, marks, trial) {
			count++
		}
	}
	return strconv.Itoa(count), nil
}
func init() { solve.Register(Day6{}) }
