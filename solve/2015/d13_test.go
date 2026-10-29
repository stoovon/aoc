package solve2015

import (
	"fmt"
	"maps"
	"math/rand/v2"
	"slices"
	"testing"
)

// Keep every rotation in this oracle to check that anchoring loses no optimum.
func seatingReference(mapping map[string]map[string]int, includeSelf bool) int {
	people := slices.Sorted(maps.Keys(mapping))
	n := len(people)
	if includeSelf {
		n++
	}
	used := make([]bool, n)
	var order []int
	best, found := 0, false
	var visit func()
	visit = func() {
		if len(order) == n {
			total := 0
			for index, a := range order {
				b := order[(index+1)%n]
				if a < len(people) && b < len(people) {
					total += mapping[people[a]][people[b]] + mapping[people[b]][people[a]]
				}
			}
			if !found || total > best {
				best, found = total, true
			}
			return
		}
		for guest := range n {
			if used[guest] {
				continue
			}
			used[guest] = true
			order = append(order, guest)
			visit()
			order = order[:len(order)-1]
			used[guest] = false
		}
	}
	visit()
	return best
}

func TestDay13SubsetSeatingsMatchEnumeration(t *testing.T) {
	random := rand.New(rand.NewPCG(2015, 13))
	for n := 1; n <= 7; n++ {
		for trial := range 10 {
			mapping := make(map[string]map[string]int)
			for a := range n {
				mapping[fmt.Sprint(a)] = make(map[string]int)
				for b := range n {
					if a != b {
						mapping[fmt.Sprint(a)][fmt.Sprint(b)] = random.IntN(41) - 20
					}
				}
			}
			for _, includeSelf := range []bool{false, true} {
				want := seatingReference(mapping, includeSelf)
				got, err := findMaxHappiness(mapping, includeSelf)
				if err != nil || got != want {
					t.Fatalf("n=%d trial=%d includeSelf=%v: got %d, %v; want %d", n, trial, includeSelf, got, err, want)
				}
			}
		}
	}
}

func TestDay13NegativeHappinessAndNeutralGuest(t *testing.T) {
	input := "Alice would lose 5 happiness units by sitting next to Bob.\nBob would lose 7 happiness units by sitting next to Alice."
	for part, test := range []struct {
		solver func(string) (string, error)
		want   string
	}{{(Day13{}).Part1, "-24"}, {(Day13{}).Part2, "-12"}} {
		got, err := test.solver(input)
		if err != nil || got != test.want {
			t.Fatalf("part %d: got %q, %v; want %q", part+1, got, err, test.want)
		}
	}
	mapping := map[string]map[string]int{"me": {"other": 4}, "other": {"me": 3}}
	original := map[string]map[string]int{"me": maps.Clone(mapping["me"]), "other": maps.Clone(mapping["other"])}
	got, err := findMaxHappiness(mapping, true)
	if err != nil || got != 7 || !maps.Equal(mapping["me"], original["me"]) || !maps.Equal(mapping["other"], original["other"]) || len(mapping) != 2 {
		t.Fatalf("neutral guest changed existing guests or mapping: score %d, error %v, mapping %v", got, err, mapping)
	}
}

func TestDay13SyntheticSeating(t *testing.T) {
	input := "Ada would gain 2 happiness units by sitting next to Bea.\nAda would gain 4 happiness units by sitting next to Cy.\nBea would gain 6 happiness units by sitting next to Ada.\nBea would gain 8 happiness units by sitting next to Cy.\nCy would gain 10 happiness units by sitting next to Ada.\nCy would gain 12 happiness units by sitting next to Bea."
	for part, test := range []struct {
		solver func(string) (string, error)
		want   string
	}{{(Day13{}).Part1, "42"}, {(Day13{}).Part2, "34"}} {
		got, err := test.solver(input)
		if err != nil || got != test.want {
			t.Fatalf("part %d: got %q, %v; want %q", part+1, got, err, test.want)
		}
	}
}

func TestDay13SeatingBudgetAndEmptyInput(t *testing.T) {
	mapping := make(map[string]map[string]int)
	for guest := range 64 {
		mapping[fmt.Sprint(guest)] = nil
	}
	for _, includeSelf := range []bool{false, true} {
		if _, err := findMaxHappiness(mapping, includeSelf); err == nil {
			t.Fatal("expected subset table budget error")
		}
		if _, err := findMaxHappiness(nil, includeSelf); err == nil {
			t.Fatal("expected empty seating error")
		}
	}
}
