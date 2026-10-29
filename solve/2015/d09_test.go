package solve2015

import (
	"fmt"
	"math/rand/v2"
	"testing"
)

func TestDay9SyntheticRoutes(t *testing.T) {
	input := "North to East = 7\nNorth to West = 11\nEast to West = 5"
	for part, test := range []struct {
		solver func(string) (string, error)
		want   string
	}{{(Day9{}).Part1, "12"}, {(Day9{}).Part2, "18"}} {
		got, err := test.solver(input)
		if err != nil || got != test.want {
			t.Fatalf("part %d: got %q, %v; want %q", part+1, got, err, test.want)
		}
	}
}

// Enumerate full routes, retaining their order, independently of the subset table.
func routeExtremesReference(graph map[string]map[string]int) (int, int, bool) {
	seen := make(map[string]bool)
	shortest, longest, found := 0, 0, false
	var visit func(string, int)
	visit = func(last string, total int) {
		if len(seen) == len(graph) {
			if !found {
				shortest, longest, found = total, total, true
			} else {
				shortest, longest = min(shortest, total), max(longest, total)
			}
			return
		}
		for next, distance := range graph[last] {
			if seen[next] {
				continue
			}
			seen[next] = true
			visit(next, total+distance)
			delete(seen, next)
		}
	}
	for first := range graph {
		seen[first] = true
		visit(first, 0)
		delete(seen, first)
	}
	return shortest, longest, found
}

func TestDay9SubsetRoutesMatchEnumeration(t *testing.T) {
	random := rand.New(rand.NewPCG(2015, 9))
	for n := 1; n <= 7; n++ {
		for trial := range 30 {
			graph := make(map[string]map[string]int)
			for city := range n {
				graph[fmt.Sprint(city)] = make(map[string]int)
			}
			for a := range n {
				for b := a + 1; b < n; b++ {
					if trial%2 == 0 && random.IntN(3) == 0 {
						continue
					}
					distance := random.IntN(30)
					graph[fmt.Sprint(a)][fmt.Sprint(b)] = distance
					graph[fmt.Sprint(b)][fmt.Sprint(a)] = distance
				}
			}
			wantMin, wantMax, found := routeExtremesReference(graph)
			gotMin, gotMax, err := routeExtremes(graph)
			if (err == nil) != found || err == nil && (gotMin != wantMin || gotMax != wantMax) {
				t.Fatalf("n=%d trial=%d: got %d/%d, %v; want %d/%d, found=%v", n, trial, gotMin, gotMax, err, wantMin, wantMax, found)
			}
		}
	}
}

func TestDay9OpenRoutesAndMissingEdges(t *testing.T) {
	for _, test := range []struct {
		input string
		want  string
	}{
		{"A to B = 0", "0"},
		{"A to B = 3\nB to C = 4", "7"},
	} {
		for _, solver := range []func(string) (string, error){(Day9{}).Part1, (Day9{}).Part2} {
			got, err := solver(test.input)
			if err != nil || got != test.want {
				t.Fatalf("input %q: got %q, %v; want %q", test.input, got, err, test.want)
			}
		}
	}
	for _, input := range []string{"", "A to B = 2\nC to D = 3"} {
		for _, solver := range []func(string) (string, error){(Day9{}).Part1, (Day9{}).Part2} {
			if _, err := solver(input); err == nil {
				t.Fatalf("input %q should have no spanning route", input)
			}
		}
	}
}

func TestDay9RouteTableBudget(t *testing.T) {
	graph := make(map[string]map[string]int)
	for city := range 64 {
		graph[fmt.Sprint(city)] = nil
	}
	if _, _, err := routeExtremes(graph); err == nil {
		t.Fatal("expected table budget error before shifting or allocating")
	}
}
