package solve2024

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day20 struct{}

func (d Day20) Coords() solve.SolutionCoords { return solve.SolutionCoords{Year: 2024, Day: 20} }

type raceTrack struct {
	width, height int
	cells         []byte
	start, end    int
}

func parseRaceTrack(input string) (raceTrack, error) {
	track := raceTrack{start: -1, end: -1}
	rows := strings.Split(strings.TrimSuffix(input, "\n"), "\n")
	track.height = len(rows)
	for y, row := range rows {
		row = strings.TrimSuffix(row, "\r")
		if y == 0 {
			track.width = len(row)
			if track.width == 0 || track.height > math.MaxInt/track.width {
				return track, fmt.Errorf("invalid track dimensions")
			}
			track.cells = make([]byte, 0, track.width*track.height)
		}
		if len(row) != track.width {
			return track, fmt.Errorf("line %d: inconsistent track width", y+1)
		}
		for x := 0; x < len(row); x++ {
			id := y*track.width + x
			switch row[x] {
			case '#', '.':
			case 'S':
				if track.start >= 0 {
					return track, fmt.Errorf("line %d: duplicate start", y+1)
				}
				track.start = id
			case 'E':
				if track.end >= 0 {
					return track, fmt.Errorf("line %d: duplicate end", y+1)
				}
				track.end = id
			default:
				return track, fmt.Errorf("line %d: invalid track cell", y+1)
			}
		}
		track.cells = append(track.cells, row...)
	}
	if track.start < 0 || track.end < 0 {
		return track, fmt.Errorf("track requires one S and one E")
	}
	return track, nil
}

func (t raceTrack) distances(start int) []int {
	distances := make([]int, len(t.cells))
	for i := range distances {
		distances[i] = -1
	}
	distances[start] = 0
	queue := make([]int, 1, len(t.cells))
	queue[0] = start
	for head := 0; head < len(queue); head++ {
		id := queue[head]
		x, y := id%t.width, id/t.width
		for _, step := range [...][2]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}} {
			nx, ny := x+step[0], y+step[1]
			if nx < 0 || nx >= t.width || ny < 0 || ny >= t.height {
				continue
			}
			next := ny*t.width + nx
			if t.cells[next] != '#' && distances[next] < 0 {
				distances[next] = distances[id] + 1
				queue = append(queue, next)
			}
		}
	}
	return distances
}

type raceOffset struct{ dx, dy, length int }

func (t raceTrack) countCheats(radius, minSaving int) (uint64, error) {
	if radius < 1 || minSaving < 0 {
		return 0, fmt.Errorf("invalid cheat rules")
	}
	fromStart := t.distances(t.start)
	normal := fromStart[t.end]
	if normal < 0 {
		return 0, fmt.Errorf("race finish is unreachable")
	}
	if minSaving > normal {
		return 0, nil
	}
	toEnd := t.distances(t.end)
	budget := normal - minSaving
	offsets := []raceOffset{}
	maxX := radius
	if maxX >= t.width {
		maxX = t.width - 1
	}
	for dx := -maxX; dx <= maxX; dx++ {
		absX := dx
		if absX < 0 {
			absX = -absX
		}
		maxY := radius - absX
		if maxY >= t.height {
			maxY = t.height - 1
		}
		for dy := -maxY; dy <= maxY; dy++ {
			absY := dy
			if absY < 0 {
				absY = -absY
			}
			length := absX + absY
			if length > 0 {
				offsets = append(offsets, raceOffset{dx, dy, length})
			}
		}
	}
	var count uint64
	for id, prefix := range fromStart {
		if prefix < 0 || prefix > budget {
			continue
		}
		x, y := id%t.width, id/t.width
		remaining := budget - prefix
		for _, offset := range offsets {
			if offset.length > remaining {
				continue
			}
			nx, ny := x+offset.dx, y+offset.dy
			if nx < 0 || nx >= t.width || ny < 0 || ny >= t.height {
				continue
			}
			suffix := toEnd[ny*t.width+nx]
			if suffix >= 0 && suffix <= remaining-offset.length {
				if count == math.MaxUint64 {
					return 0, fmt.Errorf("cheat count exceeds uint64")
				}
				count++
			}
		}
	}
	return count, nil
}

func raceCheats(input string, radius int) (string, error) {
	track, err := parseRaceTrack(input)
	if err != nil {
		return "", err
	}
	count, err := track.countCheats(radius, 100)
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(count, 10), nil
}
func (d Day20) Part1(input string) (string, error) { return raceCheats(input, 2) }
func (d Day20) Part2(input string) (string, error) { return raceCheats(input, 20) }
func init()                                        { solve.Register(Day20{}) }
