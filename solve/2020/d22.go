package solve2020

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"aoc/solve"
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
	winner, err := playStandardCombat(p1, p2)
	if err != nil {
		return "", err
	}
	return scoreDeck(winner), nil
}

func (d Day22) Part2(input string) (string, error) {
	p1, p2, err := parseDecks(input)
	if err != nil {
		return "", err
	}
	_, winner := playRecursiveCombat(p1, p2, make(map[combatState]int))
	return scoreDeck(winner), nil
}

func parseDecks(input string) ([]int, []int, error) {
	sections := strings.Split(strings.TrimSpace(strings.ReplaceAll(input, "\r\n", "\n")), "\n\n")
	if len(sections) != 2 {
		return nil, nil, errors.New("input must have two player decks")
	}
	var decks [2][]int
	used := make(map[int]bool)
	for player, section := range sections {
		lines := strings.Split(section, "\n")
		if strings.TrimSpace(lines[0]) != fmt.Sprintf("Player %d:", player+1) {
			return nil, nil, fmt.Errorf("player %d: invalid deck header", player+1)
		}
		if len(lines) < 2 {
			return nil, nil, fmt.Errorf("player %d: empty deck", player+1)
		}
		decks[player] = make([]int, 0, len(lines)-1)
		for line, text := range lines[1:] {
			text = strings.TrimSpace(text)
			if text == "" {
				return nil, nil, fmt.Errorf("player %d card %d: empty value", player+1, line+1)
			}
			for i := 0; i < len(text); i++ {
				if text[i] < '0' || text[i] > '9' {
					return nil, nil, fmt.Errorf("player %d card %d: expected positive decimal int", player+1, line+1)
				}
			}
			value, err := strconv.Atoi(text)
			if err != nil || value == 0 {
				return nil, nil, fmt.Errorf("player %d card %d: expected positive int", player+1, line+1)
			}
			if used[value] {
				return nil, nil, errors.New("card values must be distinct across both decks")
			}
			used[value] = true
			decks[player] = append(decks[player], value)
		}
	}
	return decks[0], decks[1], nil
}

func scoreDeck(deck []int) string {
	var score, term, weight big.Int
	for i, value := range deck {
		term.SetInt64(int64(value))
		weight.SetInt64(int64(len(deck) - i))
		term.Mul(&term, &weight)
		score.Add(&score, &term)
	}
	return score.String()
}

type combatState string

func combatKeyInto(buffer []byte, a, b []int) (combatState, []byte) {
	buffer = buffer[:0]
	buffer = binary.AppendUvarint(buffer, uint64(len(a)))
	buffer = binary.AppendUvarint(buffer, uint64(len(b)))
	for _, deck := range [2][]int{a, b} {
		for _, value := range deck {
			buffer = binary.AppendUvarint(buffer, uint64(value))
		}
	}
	// Conversion copies the scratch bytes; retained keys must stay immutable.
	return combatState(string(buffer)), buffer
}

func combatKey(a, b []int) combatState {
	key, _ := combatKeyInto(nil, a, b)
	return key
}

func copyCombatDeck(deck []int, count int) []int {
	result := make([]int, count)
	copy(result, deck[:count])
	return result
}

func awardCombatRound(p1, p2 []int, c1, c2, winner int) ([]int, []int) {
	if winner == 1 {
		p1 = append(p1, c1, c2)
	} else {
		p2 = append(p2, c2, c1)
	}
	return p1, p2
}

func playStandardCombat(p1, p2 []int) ([]int, error) {
	seen := make(map[combatState]struct{})
	var buffer []byte
	for len(p1) > 0 && len(p2) > 0 {
		var state combatState
		state, buffer = combatKeyInto(buffer, p1, p2)
		if _, repeat := seen[state]; repeat {
			return nil, errors.New("ordinary combat does not terminate")
		}
		seen[state] = struct{}{}
		c1, c2 := p1[0], p2[0]
		p1, p2 = p1[1:], p2[1:]
		winner := 2
		if c1 > c2 {
			winner = 1
		}
		p1, p2 = awardCombatRound(p1, p2, c1, c2, winner)
	}
	if len(p1) > 0 {
		return p1, nil
	}
	return p2, nil
}

func playRecursiveCombat(p1, p2 []int, winners map[combatState]int) (int, []int) {
	// A repeated round state only counts within this particular game.
	seen := make(map[combatState]struct{})
	var buffer []byte
	for len(p1) > 0 && len(p2) > 0 {
		var state combatState
		state, buffer = combatKeyInto(buffer, p1, p2)
		if _, repeat := seen[state]; repeat {
			return 1, p1
		}
		seen[state] = struct{}{}
		c1, c2 := p1[0], p2[0]
		p1, p2 = p1[1:], p2[1:]
		winner := 2
		if len(p1) >= c1 && len(p2) >= c2 {
			var initial combatState
			initial, buffer = combatKeyInto(buffer, p1[:c1], p2[:c2])
			var cached bool
			winner, cached = winners[initial]
			if !cached {
				winner, _ = playRecursiveCombat(copyCombatDeck(p1, c1), copyCombatDeck(p2, c2), winners)
				winners[initial] = winner
			}
		} else if c1 > c2 {
			winner = 1
		}
		p1, p2 = awardCombatRound(p1, p2, c1, c2, winner)
	}
	if len(p1) > 0 {
		return 1, p1
	}
	return 2, p2
}

func init() {
	solve.Register(Day22{})
}
