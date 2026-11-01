package solve2017

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day24 struct{}

func (d Day24) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2017, Day: 24}
}

type component struct {
	a, b int
	used bool
}

func bridgePort(text string) (int, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return 0, fmt.Errorf("empty port")
	}
	for i := 0; i < len(text); i++ {
		if text[i] < '0' || text[i] > '9' {
			return 0, fmt.Errorf("expected nonnegative decimal port")
		}
	}
	value, err := strconv.Atoi(text)
	if err != nil {
		return 0, fmt.Errorf("port exceeds int")
	}
	return value, nil
}

func parseComponents(input string) ([]component, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	comps := make([]component, 0, len(lines))
	total := 0
	for line, text := range lines {
		left, right, ok := strings.Cut(text, "/")
		if !ok {
			return nil, fmt.Errorf("line %d: expected port/port", line+1)
		}
		a, err := bridgePort(left)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", line+1, err)
		}
		b, err := bridgePort(right)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", line+1, err)
		}
		for _, value := range [...]int{a, b} {
			if value > math.MaxInt-total {
				return nil, fmt.Errorf("total component strength exceeds int")
			}
			total += value
		}
		comps = append(comps, component{a: a, b: b})
	}
	return comps, nil
}

// The port index guarantees that one endpoint matches the exposed port.
func (piece component) otherPort(port int) int {
	if piece.a == port {
		return piece.b
	}
	return piece.a
}

func bridgeCandidates(comps []component) map[int][]int {
	byPort := make(map[int][]int)
	for i, piece := range comps {
		byPort[piece.a] = append(byPort[piece.a], i)
		if piece.a != piece.b {
			byPort[piece.b] = append(byPort[piece.b], i)
		}
	}
	return byPort
}

func maxBridgeStrength(comps []component, byPort map[int][]int, port int) int {
	strongest := 0
	for _, i := range byPort[port] {
		if comps[i].used {
			continue
		}
		comps[i].used = true
		strength := comps[i].a + comps[i].b + maxBridgeStrength(comps, byPort, comps[i].otherPort(port))
		comps[i].used = false
		if strength > strongest {
			strongest = strength
		}
	}
	return strongest
}

func (d Day24) Part1(input string) (string, error) {
	comps, err := parseComponents(input)
	if err != nil {
		return "", err
	}
	result := maxBridgeStrength(comps, bridgeCandidates(comps), 0)
	return strconv.Itoa(result), nil
}

func maxLongestBridge(comps []component, byPort map[int][]int, port int) (length, strength int) {
	for _, i := range byPort[port] {
		if comps[i].used {
			continue
		}
		comps[i].used = true
		nextLength, nextStrength := maxLongestBridge(comps, byPort, comps[i].otherPort(port))
		comps[i].used = false
		nextLength++
		nextStrength += comps[i].a + comps[i].b
		if nextLength > length || (nextLength == length && nextStrength > strength) {
			length, strength = nextLength, nextStrength
		}
	}
	return length, strength
}

func (d Day24) Part2(input string) (string, error) {
	comps, err := parseComponents(input)
	if err != nil {
		return "", err
	}
	_, strength := maxLongestBridge(comps, bridgeCandidates(comps), 0)
	return strconv.Itoa(strength), nil
}

func init() {
	solve.Register(Day24{})
}
