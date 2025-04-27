package solve2016

import (
	"fmt"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day11 struct {
}

func (d Day11) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2016, Day: 11}
}

type Component uint

const (
	PoloniumChip Component = iota
	ThuliumChip
	PromethiumChip
	RutheniumChip
	CobaltChip
	EleriumChip
	DilithiumChip
)

const (
	PoloniumRTG Component = 8 + iota
	ThuliumRTG
	PromethiumRTG
	RutheniumRTG
	CobaltRTG
	EleriumRTG
	DilithiumRTG
)

type Building struct {
	Floors   [4]Floor
	Elevator int
	Parent   *Building
}

func (b Building) Valid() bool {
	for _, f := range b.Floors {
		if !f.Valid() {
			return false
		}
	}
	return true
}

func (b *Building) String() string {
	s := fmt.Sprintf("Elevator:%d Floors:", b.Elevator)
	for i := 0; i < 4; i++ {
		components := b.Floors[i].Components()
		s += fmt.Sprintf("%+v\t", components)
	}
	return s + "\n"
}

func (b Building) Norm() Building {
	b.Parent = nil
	return b
}

func (b Building) AllCombinations() []Building {
	buildings := make([]Building, 0)

	pairs := b.pairUp(b.Floors[b.Elevator].Components())
	for _, p := range pairs {
		if b.Elevator != 3 {
			newB := b.Move(1, p)
			if newB.Valid() {
				buildings = append(buildings, newB)
			}
		}

		if b.Elevator != 0 {
			newB := b.Move(-1, p)
			if newB.Valid() {
				buildings = append(buildings, newB)
			}
		}
	}

	return buildings
}

func (b Building) pairUp(components []Component) [][]Component {
	pairs := make([][]Component, 0)

	for _, c := range components {
		pairs = append(pairs, []Component{c})
	}

	for i := 0; i < len(components)-1; i++ {
		for j := i + 1; j < len(components); j++ {
			pairs = append(pairs, []Component{components[i], components[j]})
		}
	}
	return pairs
}

func (b Building) Move(direction int, components []Component) Building {
	targetFloor := b.Elevator + direction
	for _, c := range components {
		(&b.Floors[b.Elevator]).Remove(c)
		(&b.Floors[targetFloor]).Add(c)
	}
	b.Elevator = targetFloor
	return b
}

type Floor uint16

func (f Floor) Valid() bool {
	rtgs := f >> 8
	chips := f & 0xFF
	unpairedChips := chips &^ rtgs

	if unpairedChips == 0 {
		return true
	}

	// Boom
	return rtgs == 0
}

func (f *Floor) Add(components ...Component) {
	for _, c := range components {
		*f |= 1 << c
	}
}

func (f *Floor) Remove(components ...Component) {
	for _, c := range components {
		*f &^= 1 << c
	}
}

func (f Floor) Components() []Component {
	var components []Component

	var i Component
	for i = 0; i < 16; i++ {
		if f&(1<<i) != 0 {
			components = append(components, i)
		}
	}
	return components
}

func (d Day11) bfs(start, end Building) *Building {
	var todo = make([]Building, 0)
	var visited = make(map[Building]bool)

	todo = append(todo, start)

	for len(todo) > 0 {
		var v Building
		v, todo = todo[0], todo[1:]

		if v.Floors == end.Floors {
			return &v
		}

		for _, b := range v.AllCombinations() {
			if !visited[b.Norm()] {
				visited[b.Norm()] = true
				b.Parent = &v
				todo = append(todo, b)
			}
		}
	}

	return &Building{}
}

func (d Day11) countPath(b *Building) int {
	count := 0
	for b != nil {
		count++
		b = b.Parent
	}

	// Skip first state
	count--

	return count
}

func (d Day11) initializeBuilding(floorComponents [4][]Component) Building {
	var b Building
	for i, components := range floorComponents {
		(&b.Floors[i]).Add(components...)
	}
	return b
}

func (d Day11) calculateEnd(start Building) Building {
	var end Building
	for i, floor := range start.Floors {
		(&end.Floors[3]).Add(floor.Components()...) // Move all components to the last floor
		if i == 3 {
			break
		}
	}
	return end
}

func (d Day11) parseInput(input string) [4][]Component {
	input = strings.TrimSpace(input)

	componentMap := map[string]Component{
		"polonium generator":              PoloniumRTG,
		"thulium generator":               ThuliumRTG,
		"promethium generator":            PromethiumRTG,
		"ruthenium generator":             RutheniumRTG,
		"cobalt generator":                CobaltRTG,
		"elerium generator":               EleriumRTG,
		"dilithium generator":             DilithiumRTG,
		"polonium-compatible microchip":   PoloniumChip,
		"thulium-compatible microchip":    ThuliumChip,
		"promethium-compatible microchip": PromethiumChip,
		"ruthenium-compatible microchip":  RutheniumChip,
		"cobalt-compatible microchip":     CobaltChip,
		"elerium-compatible microchip":    EleriumChip,
		"dilithium-compatible microchip":  DilithiumChip,
	}

	var floors [4][]Component
	lines := strings.Split(input, "\n")

	for i, line := range lines {
		line = strings.ToLower(line)
		if strings.Contains(line, "nothing relevant") {
			continue
		}

		line = strings.ReplaceAll(line, ", and a", ", a")
		line = strings.ReplaceAll(line, " and a", ", a")

		// Extract components from the line
		parts := strings.Split(line, "contains ")[1]
		components := strings.Split(parts, ", ")

		for _, component := range components {
			component = strings.TrimSuffix(component, ".")
			component = strings.TrimPrefix(component, "a ")
			component = strings.TrimSpace(component)

			if c, ok := componentMap[component]; ok {
				floors[i] = append(floors[i], c)
			}
		}
	}

	return floors
}

func (d Day11) Part1(input string) (string, error) {
	start := d.initializeBuilding(d.parseInput(input))
	end := d.calculateEnd(start)

	final := d.bfs(start, end)

	steps := d.countPath(final)
	return strconv.Itoa(steps), nil
}

func (d Day11) Part2(input string) (string, error) {
	start := d.initializeBuilding(d.parseInput(input))
	start.Floors[0].Add(
		EleriumChip,
		EleriumRTG,
		DilithiumChip,
		DilithiumRTG,
	)

	end := d.calculateEnd(start)

	final := d.bfs(start, end)

	steps := d.countPath(final)

	return strconv.Itoa(steps), nil
}

func init() {
	solve.Register(Day11{})
}
