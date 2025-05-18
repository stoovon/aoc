package solve2023

import (
	"sort"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day7 struct {
}

func (d Day7) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2023, Day: 7}
}

var strongJ = map[byte]int{'2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9, 'T': 10, 'J': 11, 'Q': 12, 'K': 13, 'A': 14}
var weakJ = map[byte]int{'J': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9, 'T': 10, 'Q': 12, 'K': 13, 'A': 14}

type hand struct {
	cards string
	bid   int
	typ   [5]int // sorted counts
	order []int  // card order for tie-break
}

// Compare two [5]int arrays lexicographically
func less(a, b [5]int) bool {
	for i := 0; i < 5; i++ {
		if a[i] < b[i] {
			return true
		}
		if a[i] > b[i] {
			return false
		}
	}
	return false
}

func greater(a, b [5]int) bool {
	for i := 0; i < 5; i++ {
		if a[i] > b[i] {
			return true
		}
		if a[i] < b[i] {
			return false
		}
	}
	return false
}

func handType(cards string, face map[byte]int, jokerWild bool) ([5]int, []int) {
	count := make(map[byte]int)
	for i := 0; i < 5; i++ {
		count[cards[i]]++
	}
	jokers := count['J']
	best := [5]int{}
	cardFaces := []byte{'2', '3', '4', '5', '6', '7', '8', '9', 'T', 'Q', 'K', 'A'}
	if jokerWild && jokers > 0 {
		for _, c := range cardFaces {
			count2 := make(map[byte]int)
			for k, v := range count {
				count2[k] = v
			}
			count2[c] += jokers
			count2['J'] = 0
			vals := make([]int, 0, 5)
			for _, v := range count2 {
				vals = append(vals, v)
			}
			sort.Sort(sort.Reverse(sort.IntSlice(vals)))
			var arr [5]int
			copy(arr[:], vals)
			if greater(arr, best) {
				best = arr
			}
		}
	} else {
		vals := make([]int, 0, 5)
		for _, v := range count {
			vals = append(vals, v)
		}
		sort.Sort(sort.Reverse(sort.IntSlice(vals)))
		copy(best[:], vals)
	}
	// Card order for tie-break
	order := make([]int, 5)
	for i := 0; i < 5; i++ {
		order[i] = face[cards[i]]
	}
	return best, order
}

func parseHands(input string, face map[byte]int, jokerWild bool) []hand {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	hands := make([]hand, 0, len(lines))
	for _, l := range lines {
		parts := strings.Fields(l)
		if len(parts) != 2 {
			continue
		}
		bid, _ := strconv.Atoi(parts[1])
		typ, order := handType(parts[0], face, jokerWild)
		hands = append(hands, hand{cards: parts[0], bid: bid, typ: typ, order: order})
	}
	return hands
}

func (d Day7) solve(input string, face map[byte]int, jokerWild bool) string {
	hands := parseHands(input, face, jokerWild)
	sort.Slice(hands, func(i, j int) bool {
		if hands[i].typ != hands[j].typ {
			return less(hands[i].typ, hands[j].typ)
		}
		for k := 0; k < 5; k++ {
			if hands[i].order[k] != hands[j].order[k] {
				return hands[i].order[k] < hands[j].order[k]
			}
		}
		return false
	})
	sum := 0
	for i, h := range hands {
		sum += (i + 1) * h.bid
	}
	return strconv.Itoa(sum)
}

func (d Day7) Part1(input string) (string, error) {
	return d.solve(input, strongJ, false), nil
}

func (d Day7) Part2(input string) (string, error) {
	return d.solve(input, weakJ, true), nil
}

func init() {
	solve.Register(Day7{})
}
