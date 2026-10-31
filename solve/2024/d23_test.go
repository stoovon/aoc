package solve2024

import (
	"fmt"
	"strings"
	"testing"
)

func TestLANExactSearchAgainstEverySixVertexGraph(t *testing.T) {
	names := []string{"aa", "ab", "ac", "ad", "ta", "tb"}
	pairs := [][2]int{}
	for a := 0; a < 6; a++ {
		for b := a + 1; b < 6; b++ {
			pairs = append(pairs, [2]int{a, b})
		}
	}
	for code := 0; code < 1<<len(pairs); code++ {
		graph := lanGraph{names: names, adjacent: make([][]uint64, 6)}
		var matrix [6][6]bool
		for i := range graph.adjacent {
			graph.adjacent[i] = make([]uint64, 1)
		}
		for bit, pair := range pairs {
			if code&(1<<bit) != 0 {
				a, b := pair[0], pair[1]
				matrix[a][b] = true
				matrix[b][a] = true
				graph.adjacent[a][0] |= 1 << b
				graph.adjacent[b][0] |= 1 << a
			}
		}
		bestSize, bestPassword := 0, ""
		var triangles uint64
		for subset := 1; subset < 64; subset++ {
			clique := true
			members := []string{}
			hasT := false
			for a := 0; a < 6; a++ {
				if subset&(1<<a) == 0 {
					continue
				}
				members = append(members, names[a])
				hasT = hasT || names[a][0] == 't'
				for b := a + 1; b < 6; b++ {
					if subset&(1<<b) != 0 && !matrix[a][b] {
						clique = false
					}
				}
			}
			if !clique {
				continue
			}
			password := strings.Join(members, ",")
			if len(members) > bestSize || len(members) == bestSize && password < bestPassword {
				bestSize = len(members)
				bestPassword = password
			}
			if len(members) == 3 && hasT {
				triangles++
			}
		}
		ids := graph.largestClique()
		gotNames := []string{}
		for _, id := range ids {
			gotNames = append(gotNames, names[id])
		}
		if got := strings.Join(gotNames, ","); got != bestPassword {
			t.Fatalf("graph %d: got %s want %s", code, got, bestPassword)
		}
		if got := graph.triangles(); got != triangles {
			t.Fatalf("graph %d: triangles %d want %d", code, got, triangles)
		}
	}
}

func TestLANGreedyTrapAndInputOrder(t *testing.T) {
	edges := []string{"aa-ba", "ab-bb", "ac-bc", "ta-bd", "aa-ab", "aa-ac", "aa-ta", "ab-ac", "ab-ta", "ac-ta"}
	for ordering := 0; ordering < 3; ordering++ {
		input := strings.Join(edges, "\n")
		if ordering == 1 {
			input += "\naa-ab\nab-aa\n"
		}
		if ordering == 2 {
			input = strings.ReplaceAll(input, "\n", "\r\n") + "\r\n"
		}
		if got, err := (Day23{}).Part2(input); err != nil || got != "aa,ab,ac,ta" {
			t.Fatalf("greedy trap: %q %v", got, err)
		}
		if got, err := (Day23{}).Part1(input); err != nil || got != "3" {
			t.Fatalf("triangles: %q %v", got, err)
		}
		for a, b := 0, len(edges)-1; a < b; a, b = a+1, b-1 {
			edges[a], edges[b] = edges[b], edges[a]
		}
	}
	if got, err := (Day23{}).Part2("ac-ad\naa-ab"); err != nil || got != "aa,ab" {
		t.Fatalf("tied maximum: %q %v", got, err)
	}
}

func TestLANParserAndWideAdjacency(t *testing.T) {
	for _, input := range []string{"", "aa", "aa-bb-cc", "aa-aa", "a-bb", "AA-bb", "aa-b1", " aa-bb", "aa-bb\n\ncc-dd"} {
		for _, solver := range []func(string) (string, error){(Day23{}).Part1, (Day23{}).Part2} {
			if _, err := solver(input); err == nil {
				t.Fatalf("accepted %q", input)
			}
		}
	}
	names := []string{}
	for i := 0; i < 130; i++ {
		names = append(names, fmt.Sprintf("%c%c", 'a'+i/26, 'a'+i%26))
	}
	edges := []string{}
	for i := 1; i < len(names); i++ {
		edges = append(edges, names[0]+"-"+names[i])
	}
	edges = append(edges, names[64]+"-"+names[129])
	input := strings.Join(edges, "\n")
	want := strings.Join([]string{names[0], names[64], names[129]}, ",")
	if got, err := (Day23{}).Part2(input); err != nil || got != want {
		t.Fatalf("wide adjacency: %q %v want %s", got, err, want)
	}
	if got, err := (Day23{}).Part1(input); err != nil || got != "0" {
		t.Fatalf("wide triangles: %q %v", got, err)
	}
}
