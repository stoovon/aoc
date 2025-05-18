package solve2023

import (
	"errors"
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/maths"
)

type Day20 struct{}

func (d Day20) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2023, Day: 20}
}

type PulseType int

const (
	Low PulseType = iota
	High
)

type ModuleType int

const (
	Broadcaster ModuleType = iota
	FlipFlop
	Conjunction
)

type Pulse struct {
	To, From string
	Type     PulseType
}

func (d Day20) parseInput(input string) (map[string][]string, map[string]ModuleType, map[string]map[string]PulseType) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	routes := make(map[string][]string)
	types := make(map[string]ModuleType)
	froms := make(map[string]map[string]PulseType)
	for _, line := range lines {
		parts := strings.Split(line, "->")
		src := strings.TrimSpace(parts[0])
		dests := strings.Split(strings.TrimSpace(parts[1]), ", ")
		name := src
		typ := Broadcaster
		if src == "broadcaster" {
			typ = Broadcaster
		} else if src[0] == '%' {
			typ = FlipFlop
			name = src[1:]
		} else if src[0] == '&' {
			typ = Conjunction
			name = src[1:]
		}
		types[name] = typ
		routes[name] = dests
		if typ == Conjunction {
			froms[name] = make(map[string]PulseType)
		}
	}
	// Fill in conjunction inputs
	for from, tos := range routes {
		for _, to := range tos {
			if t, ok := types[to]; ok && t == Conjunction {
				if froms[to] == nil {
					froms[to] = make(map[string]PulseType)
				}
				froms[to][from] = Low
			}
		}
	}
	return routes, types, froms
}

func (d Day20) Part1(input string) (string, error) {
	routes, types, froms := d.parseInput(input)
	on := make(map[string]bool)
	lo, hi := 0, 0
	queue := make([]Pulse, 0, 32)
	for t := 0; t < 1000; t++ {
		queue = append(queue, Pulse{To: "broadcaster", From: "button", Type: Low})
		for len(queue) > 0 {
			p := queue[0]
			queue = queue[1:]
			if p.Type == Low {
				lo++
			} else {
				hi++
			}
			typ, ok := types[p.To]
			if !ok {
				continue
			}
			switch typ {
			case Broadcaster:
				for _, to := range routes[p.To] {
					queue = append(queue, Pulse{To: to, From: p.To, Type: p.Type})
				}
			case FlipFlop:
				if p.Type == High {
					continue
				}
				if !on[p.To] {
					on[p.To] = true
					for _, to := range routes[p.To] {
						queue = append(queue, Pulse{To: to, From: p.To, Type: High})
					}
				} else {
					on[p.To] = false
					for _, to := range routes[p.To] {
						queue = append(queue, Pulse{To: to, From: p.To, Type: Low})
					}
				}
			case Conjunction:
				froms[p.To][p.From] = p.Type
				out := Low
				for _, v := range froms[p.To] {
					if v == Low {
						out = High
						break
					}
				}
				for _, to := range routes[p.To] {
					queue = append(queue, Pulse{To: to, From: p.To, Type: out})
				}
			}
		}
	}
	return strconv.Itoa(lo * hi), nil
}

func (d Day20) Part2(input string) (string, error) {
	routes, types, froms := d.parseInput(input)
	// Find the conjunction feeding into rx
	var rxInput string
	for from, tos := range routes {
		for _, to := range tos {
			if to == "rx" {
				rxInput = from
			}
		}
	}
	// Find all modules feeding into rxInput
	watch := []string{}
	for from, tos := range routes {
		for _, to := range tos {
			if to == rxInput {
				watch = append(watch, from)
			}
		}
	}
	cycles := make(map[string]int)
	on := make(map[string]bool)
	prev := make(map[string]int)
	count := make(map[string]int)
	queue := make([]Pulse, 0, 32)
	for t := 1; t < 1e7; t++ {
		queue = append(queue, Pulse{To: "broadcaster", From: "button", Type: Low})
		for len(queue) > 0 {
			p := queue[0]
			queue = queue[1:]
			typ, ok := types[p.To]
			if !ok {
				continue
			}
			if p.Type == Low && contains(watch, p.To) {
				count[p.To]++
				if prev[p.To] != 0 && count[p.To] == 2 && cycles[p.To] == 0 {
					cycles[p.To] = t - prev[p.To]
				}
				prev[p.To] = t
				if len(cycles) == len(watch) {
					// Compute LCM
					res := 1
					for _, v := range cycles {
						res = maths.LCM(res, v)
					}
					return strconv.Itoa(res), nil
				}
			}
			switch typ {
			case Broadcaster:
				for _, to := range routes[p.To] {
					queue = append(queue, Pulse{To: to, From: p.To, Type: p.Type})
				}
			case FlipFlop:
				if p.Type == High {
					continue
				}
				if !on[p.To] {
					on[p.To] = true
					for _, to := range routes[p.To] {
						queue = append(queue, Pulse{To: to, From: p.To, Type: High})
					}
				} else {
					on[p.To] = false
					for _, to := range routes[p.To] {
						queue = append(queue, Pulse{To: to, From: p.To, Type: Low})
					}
				}
			case Conjunction:
				froms[p.To][p.From] = p.Type
				out := Low
				for _, v := range froms[p.To] {
					if v == Low {
						out = High
						break
					}
				}
				for _, to := range routes[p.To] {
					queue = append(queue, Pulse{To: to, From: p.To, Type: out})
				}
			}
		}
	}
	return "", errors.New("no solution found")
}

func contains(xs []string, x string) bool {
	for _, y := range xs {
		if y == x {
			return true
		}
	}
	return false
}

func init() {
	solve.Register(Day20{})
}
