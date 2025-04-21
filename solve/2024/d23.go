package solve2024

import (
	"sort"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day23 struct {
}

func (d Day23) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2024, Day: 23}
}

func run(input string) (int, string) {
	var m = make(connectionsMap)

	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "-")
		start := parts[0]
		end := parts[1]

		m.get(start).add(end)
		m.get(end).add(start)
	}

	var threes = make(map[string]bool)
	var largest = make([]string, 0)

	for name, connections := range m {
		var lan = make([]string, 1)
		lan[0] = name

		for _, v1 := range connections.arr {
			a := m[v1]

			var include = true
			for _, n := range lan {
				if !a.has(n) {
					include = false
					break
				}
			}

			if include {
				lan = append(lan, v1)
			}

			for _, v2 := range connections.arr {
				if v1 == v2 {
					continue
				}

				b := m[v2]

				if a.has(v2) && b.has(v1) {
					var names = make([]string, 3)
					names[0] = name
					names[1] = v1
					names[2] = v2
					sort.Strings(names)

					threes[strings.Join(names, "-")] = true
				}
			}
		}

		sort.Strings(lan)

		if len(lan) > len(largest) {
			largest = lan
		}
	}

	var partOne int
	for v := range threes {
		if v[0] == 't' || v[3] == 't' || v[6] == 't' {
			partOne++
		}
	}

	return partOne, strings.Join(largest, ",")
}

type connectionsMap map[string]*connections

func (m connectionsMap) get(name string) *connections {
	if v, ok := m[name]; ok {
		return v
	}
	m[name] = &connections{
		arr: []string{},
	}
	return m[name]
}

type connections struct {
	arr []string
}

func (c *connections) add(name string) {
	c.arr = append(c.arr, name)
}

func (c *connections) has(name string) bool {
	for _, v := range c.arr {
		if v == name {
			return true
		}
	}

	return false
}

func (d Day23) Part1(input string) (string, error) {
	totalPartOne, _ := run(input)

	return strconv.Itoa(totalPartOne), nil
}

func (d Day23) Part2(input string) (string, error) {
	_, totalPartTwo := run(input)

	return totalPartTwo, nil
}

func init() {
	solve.Register(Day23{})
}
