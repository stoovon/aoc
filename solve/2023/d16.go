package solve2023

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Direction int

const (
	Right Direction = iota
	Down
	Left
	Up
)

type beamState struct {
	row, col int
	dir      Direction
}

type day16Grid struct {
	device     [][]rune
	visited    map[beamState]struct{}
	energized  map[[2]int]struct{}
	rows, cols int
}

func newDay16Grid(input string) *day16Grid {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	device := make([][]rune, len(lines))
	for i, line := range lines {
		device[i] = []rune(line)
	}
	return &day16Grid{
		device:    device,
		visited:   make(map[beamState]struct{}),
		energized: make(map[[2]int]struct{}),
		rows:      len(device),
		cols:      len(device[0]),
	}
}

func (g *day16Grid) reset() {
	g.visited = make(map[beamState]struct{})
	g.energized = make(map[[2]int]struct{})
}

func (g *day16Grid) beam(row, col int, dir Direction) {
	stack := []beamState{{row, col, dir}}
	for len(stack) > 0 {
		b := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if b.row < 0 || b.row >= g.rows || b.col < 0 || b.col >= g.cols {
			continue
		}
		if _, ok := g.visited[b]; ok {
			continue
		}
		g.visited[b] = struct{}{}
		g.energized[[2]int{b.row, b.col}] = struct{}{}

		cell := g.device[b.row][b.col]
		switch cell {
		case '/':
			switch b.dir {
			case Up:
				b.dir = Right
			case Right:
				b.dir = Up
			case Down:
				b.dir = Left
			case Left:
				b.dir = Down
			}
			stack = append(stack, g.next(b))
		case '\\':
			switch b.dir {
			case Up:
				b.dir = Left
			case Left:
				b.dir = Up
			case Down:
				b.dir = Right
			case Right:
				b.dir = Down
			}
			stack = append(stack, g.next(b))
		case '-':
			if b.dir == Up || b.dir == Down {
				stack = append(stack, g.next(beamState{b.row, b.col, Left}))
				stack = append(stack, g.next(beamState{b.row, b.col, Right}))
			} else {
				stack = append(stack, g.next(b))
			}
		case '|':
			if b.dir == Left || b.dir == Right {
				stack = append(stack, g.next(beamState{b.row, b.col, Up}))
				stack = append(stack, g.next(beamState{b.row, b.col, Down}))
			} else {
				stack = append(stack, g.next(b))
			}
		default:
			stack = append(stack, g.next(b))
		}
	}
}

func (g *day16Grid) next(b beamState) beamState {
	switch b.dir {
	case Up:
		b.row--
	case Down:
		b.row++
	case Left:
		b.col--
	case Right:
		b.col++
	}
	return b
}

type Day16 struct{}

func (d Day16) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2023, Day: 16}
}

func (d Day16) Part1(input string) (string, error) {
	g := newDay16Grid(input)
	g.beam(0, 0, Right)
	return strconv.Itoa(len(g.energized)), nil
}

func (d Day16) Part2(input string) (string, error) {
	g := newDay16Grid(input)
	maxSeen := 0
	for col := 0; col < g.cols; col++ {
		g.reset()
		g.beam(0, col, Down)
		if len(g.energized) > maxSeen {
			maxSeen = len(g.energized)
		}
		g.reset()
		g.beam(g.rows-1, col, Up)
		if len(g.energized) > maxSeen {
			maxSeen = len(g.energized)
		}
	}
	for row := 0; row < g.rows; row++ {
		g.reset()
		g.beam(row, 0, Right)
		if len(g.energized) > maxSeen {
			maxSeen = len(g.energized)
		}
		g.reset()
		g.beam(row, g.cols-1, Left)
		if len(g.energized) > maxSeen {
			maxSeen = len(g.energized)
		}
	}
	return strconv.Itoa(maxSeen), nil
}

func init() {
	solve.Register(Day16{})
}
