package solve2020

import (
	"aoc/solve"
	"errors"
	"strconv"
	"strings"
)

type Day12 struct{}

func (d Day12) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2020, Day: 12}
}

type Instruction struct {
	action rune
	value  int
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func manhattanDistance(x, y int) int {
	return abs(x) + abs(y)
}

func parseInstructions(input string) ([]Instruction, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	instructions := make([]Instruction, len(lines))

	for i, line := range lines {
		if len(line) < 2 {
			return nil, errors.New("invalid instruction format")
		}

		action := rune(line[0])
		value, err := strconv.Atoi(line[1:])
		if err != nil {
			return nil, err
		}

		instructions[i] = Instruction{action: action, value: value}
	}

	return instructions, nil
}

func (d Day12) Part1(input string) (string, error) {
	instructions, err := parseInstructions(input)
	if err != nil {
		return "", err
	}

	// Ship position and direction
	x, y := 0, 0
	direction := 1 // 0=North, 1=East, 2=South, 3=West (start facing East)

	for _, inst := range instructions {
		switch inst.action {
		case 'N':
			y += inst.value
		case 'S':
			y -= inst.value
		case 'E':
			x += inst.value
		case 'W':
			x -= inst.value
		case 'L':
			turns := inst.value / 90
			direction = (direction - turns + 4) % 4
		case 'R':
			turns := inst.value / 90
			direction = (direction + turns) % 4
		case 'F':
			switch direction {
			case 0: // North
				y += inst.value
			case 1: // East
				x += inst.value
			case 2: // South
				y -= inst.value
			case 3: // West
				x -= inst.value
			}
		default:
			return "", errors.New("unknown action: " + string(inst.action))
		}
	}

	return strconv.Itoa(manhattanDistance(x, y)), nil
}

func rotateWaypoint(wx, wy int, degrees int, clockwise bool) (int, int) {
	// Normalize degrees to 0, 90, 180, 270
	degrees = degrees % 360
	if degrees < 0 {
		degrees += 360
	}

	// Convert counter-clockwise to clockwise if needed
	if !clockwise {
		degrees = 360 - degrees
	}

	// Apply rotation
	switch degrees {
	case 0:
		return wx, wy
	case 90:
		return wy, -wx
	case 180:
		return -wx, -wy
	case 270:
		return -wy, wx
	default:
		// This shouldn't happen with valid input
		return wx, wy
	}
}

func (d Day12) Part2(input string) (string, error) {
	instructions, err := parseInstructions(input)
	if err != nil {
		return "", err
	}

	// Ship position
	shipX, shipY := 0, 0
	// Waypoint position relative to ship
	waypointX, waypointY := 10, 1

	for _, inst := range instructions {
		switch inst.action {
		case 'N':
			waypointY += inst.value
		case 'S':
			waypointY -= inst.value
		case 'E':
			waypointX += inst.value
		case 'W':
			waypointX -= inst.value
		case 'L':
			waypointX, waypointY = rotateWaypoint(waypointX, waypointY, inst.value, false)
		case 'R':
			waypointX, waypointY = rotateWaypoint(waypointX, waypointY, inst.value, true)
		case 'F':
			shipX += waypointX * inst.value
			shipY += waypointY * inst.value
		default:
			return "", errors.New("unknown action: " + string(inst.action))
		}
	}

	manhattanDist := manhattanDistance(shipX, shipY)
	return strconv.Itoa(manhattanDist), nil
}

func init() {
	solve.Register(Day12{})
}
