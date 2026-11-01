package solve2016

import (
	"fmt"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day11 struct{}

func (d Day11) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2016, Day: 11}
}

type Component uint

type Building struct {
	Floors   [4]Floor
	Elevator int
}

func (b Building) Valid() bool {
	for _, f := range b.Floors {
		if !f.Valid() {
			return false
		}
	}
	return true
}

type buildingKey struct {
	elevator uint8
	count    uint8
	pairs    [8]uint8
}

// Renaming complete chip/generator pairs leaves the remaining search unchanged.
func (b Building) Norm() buildingKey {
	key := buildingKey{elevator: uint8(b.Elevator)}
	for element := 0; element < 8; element++ {
		chip, generator := uint8(4), uint8(4)
		for floor, items := range b.Floors {
			if items&(1<<element) != 0 {
				chip = uint8(floor)
			}
			if items&(1<<(element+8)) != 0 {
				generator = uint8(floor)
			}
		}
		if chip == 4 && generator == 4 {
			continue
		}
		pair := chip*5 + generator
		index := int(key.count)
		for index > 0 && key.pairs[index-1] > pair {
			key.pairs[index] = key.pairs[index-1]
			index--
		}
		key.pairs[index] = pair
		key.count++
	}
	return key
}

// The building is safe on entry; only the two changed floors need rechecking.
func (b Building) AllCombinations() []Building {
	var buildings []Building
	var items [16]Floor
	count := 0
	for component := 0; component < 16; component++ {
		bit := Floor(1) << component
		if b.Floors[b.Elevator]&bit != 0 {
			items[count] = bit
			count++
		}
	}
	appendLoad := func(load Floor) {
		departure := b.Floors[b.Elevator] &^ load
		if !departure.Valid() {
			return
		}
		for _, direction := range [...]int{1, -1} {
			target := b.Elevator + direction
			if target < 0 || target >= len(b.Floors) {
				continue
			}
			arrival := b.Floors[target] | load
			if !arrival.Valid() {
				continue
			}
			next := b
			next.Floors[b.Elevator] = departure
			next.Floors[target] = arrival
			next.Elevator = target
			buildings = append(buildings, next)
		}
	}
	for i := 0; i < count; i++ {
		appendLoad(items[i])
	}
	for i := 0; i < count; i++ {
		for j := i + 1; j < count; j++ {
			appendLoad(items[i] | items[j])
		}
	}
	return buildings
}

type Floor uint16

func (f Floor) Valid() bool {
	rtgs := f >> 8
	chips := f & 0xFF
	unpairedChips := chips &^ rtgs

	if unpairedChips == 0 {
		return true
	}

	return rtgs == 0
}

func (f *Floor) Add(components ...Component) {
	for _, c := range components {
		*f |= 1 << c
	}
}

func (d Day11) bfs(start, end Building) (int, error) {
	if start.Elevator < 0 || start.Elevator >= len(start.Floors) {
		return 0, fmt.Errorf("invalid elevator floor")
	}
	if !start.Valid() {
		return 0, fmt.Errorf("unsafe initial building")
	}

	frontier := []Building{start}
	var next []Building
	visited := map[buildingKey]bool{start.Norm(): true}
	for distance := 0; len(frontier) > 0; distance++ {
		next = next[:0]
		for _, state := range frontier {
			if state.Floors == end.Floors {
				return distance, nil
			}
			for _, candidate := range state.AllCombinations() {
				key := candidate.Norm()
				if visited[key] {
					continue
				}
				visited[key] = true
				next = append(next, candidate)
			}
		}
		frontier, next = next, frontier
	}
	return 0, fmt.Errorf("no safe elevator route")
}

func (d Day11) parseInput(input string) (Building, error) {
	var building Building
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) != 4 {
		return building, fmt.Errorf("expected four floor descriptions")
	}
	elements := make(map[string]uint)
	var seen Floor
	for floor, ordinal := range [...]string{"first", "second", "third", "fourth"} {
		line := strings.ToLower(strings.TrimSpace(lines[floor]))
		body, ok := strings.CutPrefix(line, "the "+ordinal+" floor contains ")
		if !ok {
			return building, fmt.Errorf("floor %d: invalid description", floor+1)
		}
		body = strings.TrimSuffix(body, ".")
		if body == "nothing relevant" {
			continue
		}
		body = strings.ReplaceAll(body, ", and ", ", ")
		body = strings.ReplaceAll(body, " and ", ", ")
		for _, item := range strings.Split(body, ", ") {
			name, article := strings.CutPrefix(item, "a ")
			if !article {
				name, article = strings.CutPrefix(item, "an ")
			}
			if !article {
				return building, fmt.Errorf("floor %d: invalid component description", floor+1)
			}
			element, generator := strings.CutSuffix(name, " generator")
			if !generator {
				var chip bool
				element, chip = strings.CutSuffix(name, "-compatible microchip")
				if !chip {
					return building, fmt.Errorf("floor %d: expected generator or compatible microchip", floor+1)
				}
			}
			if element == "" {
				return building, fmt.Errorf("floor %d: empty element name", floor+1)
			}
			for i := range element {
				if element[i] < 'a' || element[i] > 'z' {
					return building, fmt.Errorf("floor %d: expected alphabetic element name", floor+1)
				}
			}
			id, found := elements[element]
			if !found {
				if len(elements) == 8 {
					return building, fmt.Errorf("at most eight element pairs supported")
				}
				id = uint(len(elements))
				elements[element] = id
			}
			if generator {
				id += 8
			}
			bit := Floor(1) << id
			if seen&bit != 0 {
				return building, fmt.Errorf("floor %d: duplicate component", floor+1)
			}
			seen |= bit
			building.Floors[floor] |= bit
		}
	}
	if seen&0xff != seen>>8 {
		return building, fmt.Errorf("each element needs a chip and generator")
	}
	if !building.Valid() {
		return building, fmt.Errorf("unsafe initial building")
	}
	return building, nil
}

func (d Day11) solve(input string, extraPairs int) (string, error) {
	start, err := d.parseInput(input)
	if err != nil {
		return "", err
	}
	count := int(start.Norm().count)
	if count+extraPairs > 8 {
		return "", fmt.Errorf("at most eight element pairs supported including additions")
	}
	for element := count; element < count+extraPairs; element++ {
		start.Floors[0].Add(Component(element), Component(element+8))
	}
	var end Building
	for _, floor := range start.Floors {
		end.Floors[3] |= floor
	}
	steps, err := d.bfs(start, end)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(steps), nil
}
func (d Day11) Part1(input string) (string, error) { return d.solve(input, 0) }
func (d Day11) Part2(input string) (string, error) { return d.solve(input, 2) }
func init()                                        { solve.Register(Day11{}) }
