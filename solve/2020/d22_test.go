package solve2020

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

type referenceCombatState struct{ a, b string }

func referenceRecursiveCombat(a, b []int) (int, []int) {
	a = append([]int(nil), a...)
	b = append([]int(nil), b...)
	seen := make(map[referenceCombatState]bool)
	for len(a) > 0 && len(b) > 0 {
		key := referenceCombatState{fmt.Sprint(a), fmt.Sprint(b)}
		if seen[key] {
			return 1, a
		}
		seen[key] = true
		x, y := a[0], b[0]
		a = append([]int(nil), a[1:]...)
		b = append([]int(nil), b[1:]...)
		winner := 2
		if x <= len(a) && y <= len(b) {
			winner, _ = referenceRecursiveCombat(a[:x], b[:y])
		} else if x > y {
			winner = 1
		}
		if winner == 1 {
			a = append(a, x, y)
		} else {
			b = append(b, y, x)
		}
	}
	if len(a) > 0 {
		return 1, a
	}
	return 2, b
}

func TestCombatMemoAgainstIndependentFreshGames(t *testing.T) {
	for count := 2; count <= 6; count++ {
		memo := make(map[combatState]int)
		cards := make([]int, count)
		for i := range cards {
			cards[i] = i + 1
		}
		var check func(int)
		check = func(index int) {
			if index < len(cards) {
				for i := index; i < len(cards); i++ {
					cards[index], cards[i] = cards[i], cards[index]
					check(index + 1)
					cards[index], cards[i] = cards[i], cards[index]
				}
				return
			}
			for split := 1; split < count; split++ {
				wantWinner, wantDeck := referenceRecursiveCombat(cards[:split], cards[split:])
				gotWinner, gotDeck := playRecursiveCombat(append([]int(nil), cards[:split]...), append([]int(nil), cards[split:]...), memo)
				if gotWinner != wantWinner || !reflect.DeepEqual(gotDeck, wantDeck) {
					t.Fatalf("cards%v split%d: %d %v want %d %v", cards, split, gotWinner, gotDeck, wantWinner, wantDeck)
				}
			}
		}
		check(0)
	}
}

func TestCombatMemoReusesOnlyFreshGameWinners(t *testing.T) {
	a, b := []int{2, 5, 1}, []int{1, 6, 4, 3}
	// The first draw starts a subgame from exactly these remaining prefixes.
	key := combatKey(a[1:3], b[1:2])
	winner, _ := referenceRecursiveCombat(a[1:3], b[1:2])
	memo := map[combatState]int{key: winner}
	wantWinner, wantDeck := referenceRecursiveCombat(a, b)
	for trial := 0; trial < 2; trial++ {
		gotWinner, gotDeck := playRecursiveCombat(append([]int(nil), a...), append([]int(nil), b...), memo)
		if gotWinner != wantWinner || !reflect.DeepEqual(gotDeck, wantDeck) {
			t.Fatalf("cached subgame: %d %v want%d %v", gotWinner, gotDeck, wantWinner, wantDeck)
		}
	}
	if memo[key] != winner {
		t.Fatal("fresh subgame winner changed")
	}
}

func TestCombatCopiesSubgameDecks(t *testing.T) {
	parent := []int{11, 23, 37, 41}
	child := copyCombatDeck(parent, 2)
	child[0] = 99
	child = append(child, 101)
	if !reflect.DeepEqual(parent, []int{11, 23, 37, 41}) {
		t.Fatal("child altered parent storage")
	}
}

func TestCombatEncodingPreservesBoundaryAndOrder(t *testing.T) {
	cases := [][2][]int{
		{{1}, {2, 3}}, {{1, 2}, {3}}, {{2, 1}, {3}}, {{1}, {3, 2}},
		{{1, 23}, {4}}, {{12, 3}, {4}}, {{1, 257}, {513}}, {{1, 1}, {1}},
		{nil, {1, 2, 3}}, {{1, 2, 3}, nil},
	}
	seen := make(map[combatState]bool)
	var buffer []byte
	for _, decks := range cases {
		key, next := combatKeyInto(buffer, decks[0], decks[1])
		buffer = next
		if seen[key] {
			t.Fatalf("encoding collision for%v", decks)
		}
		seen[key] = true
	}
	if len(seen) != len(cases) {
		t.Fatal("scratch reuse altered retained keys")
	}
	if !seen[combatKey([]int{1}, []int{2, 3})] {
		t.Fatal("first key mutated after scratch reuse")
	}
}

func TestCombatEncodingDoesNotTruncateLargeCards(t *testing.T) {
	for _, value := range []int{127, 128, 255, 256, 16383, 16384} {
		if combatKey([]int{value}, []int{1}) == combatKey([]int{value + 256}, []int{1}) {
			t.Fatalf("card%d truncated", value)
		}
	}
}

func referenceStandardCombat(a, b []int) ([]int, bool) {
	a = append([]int(nil), a...)
	b = append([]int(nil), b...)
	seen := make(map[referenceCombatState]bool)
	for len(a) > 0 && len(b) > 0 {
		state := referenceCombatState{fmt.Sprint(a), fmt.Sprint(b)}
		if seen[state] {
			return nil, false
		}
		seen[state] = true
		x, y := a[0], b[0]
		a = append([]int(nil), a[1:]...)
		b = append([]int(nil), b[1:]...)
		if x > y {
			a = append(a, x, y)
		} else {
			b = append(b, y, x)
		}
	}
	if len(a) > 0 {
		return a, true
	}
	return b, true
}

func TestCombatStandardRulesAgainstIndependentGames(t *testing.T) {
	for count := 2; count <= 6; count++ {
		cards := make([]int, count)
		for i := range cards {
			cards[i] = i + 1
		}
		var check func(int)
		check = func(index int) {
			if index < len(cards) {
				for i := index; i < len(cards); i++ {
					cards[index], cards[i] = cards[i], cards[index]
					check(index + 1)
					cards[index], cards[i] = cards[i], cards[index]
				}
				return
			}
			for split := 1; split < count; split++ {
				want, terminates := referenceStandardCombat(cards[:split], cards[split:])
				got, err := playStandardCombat(append([]int(nil), cards[:split]...), append([]int(nil), cards[split:]...))
				if (err == nil) != terminates || (err == nil && !reflect.DeepEqual(got, want)) {
					t.Fatalf("ordinary cards%v split%d: %v %v want%v terminates%t", cards, split, got, err, want, terminates)
				}
			}
		}
		check(0)
	}
}

func TestCombatCheckedParsingAndExactScore(t *testing.T) {
	valid := "Player 1:\n301\n\nPlayer 2:\n2"
	for _, input := range []string{valid, strings.ReplaceAll(valid, "\n", "\r\n") + "\r\n"} {
		for _, fn := range []func(string) (string, error){(Day22{}).Part1, (Day22{}).Part2} {
			if got, err := fn(input); err != nil || got != "604" {
				t.Fatalf("large card: %q %v", got, err)
			}
		}
	}
	for _, input := range []string{"", "Player 1:", "Player 1:\n\nPlayer 2:\n2", strings.Replace(valid, "Player 1:", "Player 3:", 1), strings.Replace(valid, "301", "0", 1), strings.Replace(valid, "301", "-1", 1), strings.Replace(valid, "301", "+1", 1), strings.Replace(valid, "301", "abc", 1), strings.Replace(valid, "301", "2", 1), strings.Replace(valid, "301", "301\n301", 1), strings.Replace(valid, "301", "99999999999999999999999999", 1)} {
		for _, fn := range []func(string) (string, error){(Day22{}).Part1, (Day22{}).Part2} {
			if _, err := fn(input); err == nil {
				t.Fatalf("accepted malformed decks%q", input)
			}
		}
	}
	input := "Player 1:\n" + strconv.Itoa(math.MaxInt) + "\n\nPlayer 2:\n1"
	want := strconv.FormatUint(uint64(math.MaxInt)*2+1, 10)
	for _, fn := range []func(string) (string, error){(Day22{}).Part1, (Day22{}).Part2} {
		if got, err := fn(input); err != nil || got != want {
			t.Fatalf("exact score beyond int: %q %v want%q", got, err, want)
		}
	}
}
