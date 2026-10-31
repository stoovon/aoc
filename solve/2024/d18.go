package solve2024

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day18 struct{}

func (d Day18) Coords() solve.SolutionCoords { return solve.SolutionCoords{Year: 2024, Day: 18} }

type ramMemory struct {
	width     int
	falls     []int
	blockedAt []int
}

func ramCoordinate(text string, width int) (int, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return 0, fmt.Errorf("empty coordinate")
	}
	for i := 0; i < len(text); i++ {
		if text[i] < '0' || text[i] > '9' {
			return 0, fmt.Errorf("expected nonnegative decimal coordinate")
		}
	}
	value, err := strconv.Atoi(text)
	if err != nil || value >= width {
		return 0, fmt.Errorf("coordinate outside grid")
	}
	return value, nil
}

func parseRAM(input string, width int) (ramMemory, error) {
	memory := ramMemory{width: width}
	if width < 1 || width > math.MaxInt/width {
		return memory, fmt.Errorf("invalid grid dimensions")
	}
	memory.blockedAt = make([]int, width*width)
	for line, record := range strings.Split(strings.TrimSuffix(input, "\n"), "\n") {
		fields := strings.Split(strings.TrimSpace(record), ",")
		if len(fields) != 2 {
			return memory, fmt.Errorf("line %d: expected x,y", line+1)
		}
		x, err := ramCoordinate(fields[0], width)
		if err != nil {
			return memory, fmt.Errorf("line %d: x: %w", line+1, err)
		}
		y, err := ramCoordinate(fields[1], width)
		if err != nil {
			return memory, fmt.Errorf("line %d: y: %w", line+1, err)
		}
		id := y*width + x
		memory.falls = append(memory.falls, id)
		if memory.blockedAt[id] == 0 {
			memory.blockedAt[id] = len(memory.falls)
		}
	}
	return memory, nil
}

func (m ramMemory) shortestPath(prefix int) (int, bool) {
	end := len(m.blockedAt) - 1
	blocked := func(id int) bool { return m.blockedAt[id] > 0 && m.blockedAt[id] <= prefix }
	if blocked(0) || blocked(end) {
		return 0, false
	}
	distances := make([]int, len(m.blockedAt))
	for i := range distances {
		distances[i] = -1
	}
	distances[0] = 0
	queue := make([]int, 1, len(distances))
	for head := 0; head < len(queue); head++ {
		current := queue[head]
		if current == end {
			return distances[current], true
		}
		x, y := current%m.width, current/m.width
		for _, step := range [...][2]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}} {
			nx, ny := x+step[0], y+step[1]
			if nx < 0 || nx >= m.width || ny < 0 || ny >= m.width {
				continue
			}
			next := ny*m.width + nx
			if distances[next] < 0 && !blocked(next) {
				distances[next] = distances[current] + 1
				queue = append(queue, next)
			}
		}
	}
	return 0, false
}

func (m ramMemory) firstBlockingPrefix() (int, error) {
	if _, reachable := m.shortestPath(len(m.falls)); reachable {
		return 0, fmt.Errorf("no falling byte blocks the route")
	}
	low, high := 0, len(m.falls)
	for high-low > 1 {
		mid := low + (high-low)/2
		if _, reachable := m.shortestPath(mid); reachable {
			low = mid
		} else {
			high = mid
		}
	}
	return high, nil
}

func (d Day18) Part1(input string) (string, error) {
	memory, err := parseRAM(input, 71)
	if err != nil {
		return "", err
	}
	if len(memory.falls) < 1024 {
		return "", fmt.Errorf("need at least 1024 falling bytes")
	}
	distance, reachable := memory.shortestPath(1024)
	if !reachable {
		return "", fmt.Errorf("destination unreachable after 1024 bytes")
	}
	return strconv.Itoa(distance), nil
}
func (d Day18) Part2(input string) (string, error) {
	memory, err := parseRAM(input, 71)
	if err != nil {
		return "", err
	}
	prefix, err := memory.firstBlockingPrefix()
	if err != nil {
		return "", err
	}
	id := memory.falls[prefix-1]
	return fmt.Sprintf("%d,%d", id%memory.width, id/memory.width), nil
}
func init() { solve.Register(Day18{}) }
