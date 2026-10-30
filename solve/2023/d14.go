package solve2023

import (
	"aoc/solve"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Day14 struct{}

func (Day14) Coords() solve.SolutionCoords { return solve.SolutionCoords{Year: 2023, Day: 14} }

type rockPlatform struct {
	width, height int
	cells         []byte
}

const rockCellLimit = 1_048_576

var errRockLimit = errors.New("platform simulation resource limit exceeded")

type rockLimits struct{ states, bytes, work int }

var defaultRockLimits = rockLimits{65536, 32 * 1024 * 1024, 100_000_000}

const (
	rockNorth = iota
	rockWest
	rockSouth
	rockEast
)

func parseRockPlatform(input string) (rockPlatform, error) {
	var p rockPlatform
	hasRock := false
	lines := strings.Split(strings.TrimSuffix(input, "\n"), "\n")
	for row, line := range lines {
		line = strings.TrimSuffix(line, "\r")
		invalid := func(message string) error { return fmt.Errorf("line %d: %s", row+1, message) }
		if line == "" {
			return p, invalid("platform row is empty")
		}
		if p.width == 0 {
			p.width = len(line)
		}
		if len(line) != p.width {
			return p, invalid("platform rows must have equal width")
		}
		if len(line) > rockCellLimit-len(p.cells) {
			return p, errRockLimit
		}
		for _, cell := range []byte(line) {
			if cell != '.' && cell != '#' && cell != 'O' {
				return p, invalid("expected only '.', '#' or 'O'")
			}
			hasRock = hasRock || cell == 'O'
		}
		p.cells = append(p.cells, line...)
		p.height++
	}
	if !hasRock {
		return p, errors.New("platform must contain at least one round rock")
	}
	return p, nil
}

// Traverse each line toward its tilt direction, compacting rocks after blockers.
func (p *rockPlatform) tilt(direction int) {
	lines, length, lineStep, cellStep, origin := p.width, p.height, 1, p.width, 0
	switch direction {
	case rockWest:
		lines, length, lineStep, cellStep = p.height, p.width, p.width, 1
	case rockSouth:
		origin = (p.height - 1) * p.width
		cellStep = -p.width
	case rockEast:
		lines, length, lineStep, cellStep, origin = p.height, p.width, p.width, -1, p.width-1
	}
	for line := 0; line < lines; line++ {
		base, destination := origin+line*lineStep, 0
		for offset := 0; offset < length; offset++ {
			index := base + offset*cellStep
			switch p.cells[index] {
			case '#':
				destination = offset + 1
			case 'O':
				if offset != destination {
					p.cells[base+destination*cellStep] = 'O'
					p.cells[index] = '.'
				}
				destination++
			}
		}
	}
}
func (p *rockPlatform) spin() {
	for direction := rockNorth; direction <= rockEast; direction++ {
		p.tilt(direction)
	}
}
func (p *rockPlatform) cycles(target int64, limits rockLimits) error {
	if target == 0 {
		return nil
	}
	if target < 0 {
		return errors.New("negative cycle count")
	}
	if limits.states < 1 || len(p.cells) > limits.bytes {
		return errRockLimit
	}
	seen := map[string]int64{string(p.cells): 0}
	retained, work := len(p.cells), limits.work
	spin := func() error {
		cost := 4 * len(p.cells)
		if cost > work {
			return errRockLimit
		}
		work -= cost
		p.spin()
		return nil
	}
	for completed := int64(1); completed <= target; completed++ {
		if err := spin(); err != nil {
			return err
		}
		key := string(p.cells)
		if previous, ok := seen[key]; ok {
			remaining := (target - completed) % (completed - previous)
			for i := int64(0); i < remaining; i++ {
				if err := spin(); err != nil {
					return err
				}
			}
			return nil
		}
		if completed == target {
			return nil
		}
		if len(seen) >= limits.states || len(key) > limits.bytes-retained {
			return errRockLimit
		}
		retained += len(key)
		seen[key] = completed
	}
	return nil
}
func (p rockPlatform) northLoad() (int64, error) {
	var total int64
	for index, cell := range p.cells {
		if cell == 'O' {
			weight := int64(p.height - index/p.width)
			if total > math.MaxInt64-weight {
				return 0, errors.New("platform load exceeds int64")
			}
			total += weight
		}
	}
	return total, nil
}
func (Day14) Part1(input string) (string, error) {
	p, err := parseRockPlatform(input)
	if err != nil {
		return "", err
	}
	p.tilt(rockNorth)
	load, err := p.northLoad()
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(load, 10), nil
}
func (Day14) Part2(input string) (string, error) {
	p, err := parseRockPlatform(input)
	if err != nil {
		return "", err
	}
	if err = p.cycles(1_000_000_000, defaultRockLimits); err != nil {
		return "", err
	}
	load, err := p.northLoad()
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(load, 10), nil
}
func init() { solve.Register(Day14{}) }
