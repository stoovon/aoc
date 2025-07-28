package solve2020

import (
	"aoc/solve"
	"strings"
	"strconv"
)

type Day23 struct{}

func (d Day23) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2020, Day: 23}
}

// parseCups parses the input string into a slice of cup labels and returns the max label.
func parseCups(input string) ([]int, int) {
	input = strings.TrimSpace(input)
	cups := make([]int, len(input))
	max := 0
	for i, ch := range input {
		cups[i] = int(ch - '0')
		if cups[i] > max { max = cups[i] }
	}
	return cups, max
}

// buildCupLinks builds the array-based linked list for the cup game.
// next[label] = next cup label after 'label'.
func buildCupLinks(cups []int, total int) []int {
	next := make([]int, total+1) // 1-based indexing
	last := cups[0]
	for i := 1; i < len(cups); i++ {
		next[last] = cups[i]
		last = cups[i]
	}
	for v := len(cups) + 1; v <= total; v++ {
		next[last] = v
		last = v
	}
	if total > len(cups) {
		next[last] = cups[0]
	} else {
		next[last] = cups[0]
	}
	return next
}

// playCupGame simulates the cup game. If part2 is true, returns the product of the two cups after 1 as a string.
// Otherwise, returns the cup order after 1 as a string.
func playCupGame(next []int, current, moves int, part2 bool) string {
	max := len(next) - 1
	for i := 0; i < moves; i++ {
		pick1 := next[current]
		pick2 := next[pick1]
		pick3 := next[pick2]
		afterPick := next[pick3]
		// Find destination
		dest := current - 1
		if dest == 0 { dest = max }
		for dest == pick1 || dest == pick2 || dest == pick3 {
			dest--
			if dest == 0 { dest = max }
		}
		// Remove picked up cups and insert after dest
		next[current] = afterPick
		next[pick3] = next[dest]
		next[dest] = pick1
		current = next[current]
	}
	if part2 {
		a := next[1]
		b := next[a]
		prod := int64(a) * int64(b)
		return strconv.FormatInt(prod, 10)
	}
	// Otherwise, collect order after 1
	res := make([]byte, 0, max-1)
	for x := next[1]; x != 1; x = next[x] {
		res = append(res, byte('0'+x))
	}
	return string(res)
}

func (d Day23) Part1(input string) (string, error) {
	cups, _ := parseCups(input)
	next := buildCupLinks(cups, len(cups))
	res := playCupGame(next, cups[0], 100, false)
	return res, nil
}

func (d Day23) Part2(input string) (string, error) {
	cups, _ := parseCups(input)
	total := 1000000
	next := buildCupLinks(cups, total)
	res := playCupGame(next, cups[0], 10000000, true)
	return res, nil
}

func init() {
	solve.Register(Day23{})
}
