package solve2020

import (
	"aoc/solve"
	"errors"
	"strconv"
	"strings"
)

type Day22 struct{}

func (d Day22) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2020, Day: 22}
}

func (d Day22) Part1(input string) (string, error) {
	p1, p2, err := parseDecks(input)
	if err != nil {
		return "", err
	}
	_, winner := playCombat(p1, p2, false)
	return strconv.Itoa(scoreDeck(winner)), nil
}

func (d Day22) Part2(input string) (string, error) {
	p1, p2, err := parseDecks(input)
	if err != nil {
		return "", err
	}
	_, winner := playCombat(p1, p2, true)
	return strconv.Itoa(scoreDeck(winner)), nil
}

func parseDecks(input string) ([]int, []int, error) {
	sections := strings.Split(strings.TrimSpace(input), "\n\n")
	if len(sections) != 2 {
		return nil, nil, errors.New("input must have two player decks")
	}
	parseDeck := func(s string) []int {
		var deck []int
		for _, line := range strings.Split(s, "\n")[1:] {
			if line == "" { continue }
			n, _ := strconv.Atoi(line)
			deck = append(deck, n)
		}
		return deck
	}
	return parseDeck(sections[0]), parseDeck(sections[1]), nil
}

func scoreDeck(deck []int) int {
	score := 0
	for i, v := range deck {
		score += v * (len(deck) - i)
	}
	return score
}

func playCombat(p1, p2 []int, recursive bool) (int, []int) {
	type state struct{ a, b string }
	seen := make(map[state]struct{})
	copyDeck := func(d []int, n int) []int {
		out := make([]int, n)
		copy(out, d[:n])
		return out
	}
	joinInts := func(a []int) string {
		if len(a) == 0 {
			return ""
		}
		b := make([]byte, 0, len(a)*3)
		for i, v := range a {
			if i > 0 {
				b = append(b, ',')
			}
			b = strconv.AppendInt(b, int64(v), 10)
		}
		return string(b)
	}
	toKey := func(a, b []int) state {
		return state{joinInts(a), joinInts(b)}
	}
	for len(p1) > 0 && len(p2) > 0 {
		if recursive {
			k := toKey(p1, p2)
			if _, ok := seen[k]; ok {
				return 1, p1
			}
			seen[k] = struct{}{}
		}
		c1, c2 := p1[0], p2[0]
		p1, p2 = p1[1:], p2[1:]
		var winner int
		if recursive && len(p1) >= c1 && len(p2) >= c2 {
			winner, _ = playCombat(copyDeck(p1, c1), copyDeck(p2, c2), true)
		} else if c1 > c2 {
			winner = 1
		} else {
			winner = 2
		}
		if winner == 1 {
			p1 = append(p1, c1, c2)
		} else {
			p2 = append(p2, c2, c1)
		}
	}
	if len(p1) > 0 {
		return 1, p1
	}
	return 2, p2
}

func init() {
	solve.Register(Day22{})
}
