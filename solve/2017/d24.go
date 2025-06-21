package solve2017

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day24 struct {
}

func (d Day24) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2017, Day: 24}
}

type component struct {
	a, b int
	used bool
}

func parseComponents(input string) []component {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var comps []component
	for _, line := range lines {
		parts := strings.Split(line, "/")
		a, _ := strconv.Atoi(parts[0])
		b, _ := strconv.Atoi(parts[1])
		comps = append(comps, component{a, b, false})
	}
	return comps
}

func maxBridgeStrength(comps []component, port int) int {
	maxStrength := 0
	for i := range comps {
		if comps[i].used {
			continue
		}
		if comps[i].a == port || comps[i].b == port {
			comps[i].used = true
			nextPort := comps[i].b
			if comps[i].a == port {
				nextPort = comps[i].b
			} else {
				nextPort = comps[i].a
			}
			strength := comps[i].a + comps[i].b + maxBridgeStrength(comps, nextPort)
			if strength > maxStrength {
				maxStrength = strength
			}
			comps[i].used = false
		}
	}
	return maxStrength
}

func (d Day24) Part1(input string) (string, error) {
	comps := parseComponents(input)
	result := maxBridgeStrength(comps, 0)
	return strconv.Itoa(result), nil
}

func maxLongestBridge(comps []component, port int) (length, strength int) {
	maxLen, maxStr := 0, 0
	for i := range comps {
		if comps[i].used {
			continue
		}
		if comps[i].a == port || comps[i].b == port {
			comps[i].used = true
			nextPort := comps[i].b
			if comps[i].a == port {
				nextPort = comps[i].b
			} else {
				nextPort = comps[i].a
			}
			l, s := maxLongestBridge(comps, nextPort)
			l++
			s += comps[i].a + comps[i].b
			if l > maxLen || (l == maxLen && s > maxStr) {
				maxLen, maxStr = l, s
			}
			comps[i].used = false
		}
	}
	return maxLen, maxStr
}

func (d Day24) Part2(input string) (string, error) {
	comps := parseComponents(input)
	_, strength := maxLongestBridge(comps, 0)
	return strconv.Itoa(strength), nil
}

func init() {
	solve.Register(Day24{})
}
