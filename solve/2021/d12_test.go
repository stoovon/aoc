package solve2021

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestCavePathRulesAndCycles(t *testing.T) {
	for _, tc := range []struct{ input, p1, p2 string }{
		{"start-a\na-end", "1", "1"},
		{"start-a\na-b\na-end\nb-end", "2", "3"},
		{"start-end\nstart-A\nA-B", "1", "1"},
		{"start-a\na-A\nA-B\na-end", "1", "infinite"},
		{"start-A\nA-B\nB-end", "infinite", "infinite"},
		{"start-A\nA-A\nA-end", "infinite", "infinite"},
		{"start-A\nA-B\nx-end", "0", "0"},
		{"start-a\na-a\na-end", "1", "2"},
		{"start-end\nend-A\nA-B", "1", "1"},
		{"start-a\na-start\na-end\na-end\r\n", "1", "1"},
	} {
		for part, solver := range []func(string) (string, error){(Day12{}).Part1, (Day12{}).Part2} {
			want := []string{tc.p1, tc.p2}[part]
			got, err := solver(tc.input)
			if want == "infinite" {
				if err == nil || !strings.Contains(err.Error(), "infinitely") {
					t.Fatalf("part %d: expected infinite error, got %q %v", part+1, got, err)
				}
				continue
			}
			if err != nil || got != want {
				t.Fatalf("%q part %d: got %q %v want %s", tc.input, part+1, got, err, want)
			}
		}
	}
}

func TestCavePathsAgainstPathEnumeration(t *testing.T) {
	names := []string{"start", "a", "b", "c", "A", "end"}
	pairs := [][2]int{}
	for a := range names {
		for b := a + 1; b < len(names); b++ {
			pairs = append(pairs, [2]int{a, b})
		}
	}
	for mask := 0; mask < 1<<len(pairs); mask++ {
		adj := make([][]int, len(names))
		records := []string{}
		for i, p := range pairs {
			if mask&(1<<i) != 0 {
				a, b := p[0], p[1]
				adj[a] = append(adj[a], b)
				adj[b] = append(adj[b], a)
				records = append(records, names[a]+"-"+names[b])
			}
		}
		if len(adj[0]) == 0 || len(adj[5]) == 0 {
			continue
		}
		input := strings.Join(records, "\n")
		for _, used := range []bool{true, false} {
			visits := make([]int, len(names))
			visits[0] = 1
			var walk func(int, bool) uint64
			walk = func(at int, used bool) uint64 {
				if at == 5 {
					return 1
				}
				var count uint64
				for _, to := range adj[at] {
					if to == 0 {
						continue
					}
					nextUsed := used
					if to != 4 {
						if visits[to] > 0 {
							if used {
								continue
							}
							nextUsed = true
						}
						visits[to]++
					}
					count += walk(to, nextUsed)
					if to != 4 {
						visits[to]--
					}
				}
				return count
			}
			want := walk(0, used)
			got, err := cavePaths(input, used)
			if err != nil || got != strconv.FormatUint(want, 10) {
				t.Fatalf("mask %d used %t: got %q %v want %d", mask, used, got, err, want)
			}
		}
	}
}

func TestCaveParserAndWideVisitSet(t *testing.T) {
	for _, input := range []string{"", "start-end\n\n", "start--end", "start-a1\na1-end", "start-aB\naB-end", "start-\n-end", "a-end", "start-a", " start-end"} {
		if _, err := (Day12{}).Part1(input); err == nil {
			t.Fatalf("accepted %q", input)
		}
	}
	names := []string{"start"}
	for i := 0; i < 80; i++ {
		names = append(names, fmt.Sprintf("c%c%c", 'a'+i/26, 'a'+i%26))
	}
	names = append(names, "end")
	lines := []string{}
	for i := 1; i < len(names); i++ {
		lines = append(lines, names[i-1]+"-"+names[i])
	}
	for _, solver := range []func(string) (string, error){(Day12{}).Part1, (Day12{}).Part2} {
		got, err := solver(strings.Join(lines, "\n"))
		if err != nil || got != "1" {
			t.Fatalf("wide chain: %q %v", got, err)
		}
	}
}

func TestCavePathCountOverflow(t *testing.T) {
	g := caveGraph{edges: make([][]int, 1), small: []int{0}, start: 0, bytes: 1}
	previous := 0
	for layer := 0; layer < 64; layer++ {
		a := len(g.edges)
		g.edges = append(g.edges, []int{a + 2}, []int{a + 2}, nil)
		g.small = append(g.small, -1, -1, -1)
		g.edges[previous] = []int{a, a + 1}
		previous = a + 2
	}
	g.end = previous
	if _, err := g.count(true); err == nil || !strings.Contains(err.Error(), "exceeds") {
		t.Fatalf("expected overflow error, got %v", err)
	}
}
