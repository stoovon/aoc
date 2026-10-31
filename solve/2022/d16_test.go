package solve2022

import (
	"fmt"
	"math"
	"strconv"
	"testing"
)

type valveMinuteState struct {
	first, second int
	opened        uint64
}
type valveMinuteAction struct {
	position int
	opening  uint64
}

func valveMinuteOracle(g valveNetwork, minutes int, two bool) uint64 {
	states := map[valveMinuteState]uint64{{g.start, g.start, 0}: 0}
	actions := func(at int, opened uint64) []valveMinuteAction {
		result := []valveMinuteAction{{at, 0}}
		for _, to := range g.edges[at] {
			result = append(result, valveMinuteAction{to, 0})
		}
		if g.rates[at] > 0 && opened&(1<<at) == 0 {
			result = append(result, valveMinuteAction{at, 1 << at})
		}
		return result
	}
	for minute := 0; minute < minutes; minute++ {
		next := map[valveMinuteState]uint64{}
		for state, score := range states {
			var rate uint64
			for id, r := range g.rates {
				if state.opened&(1<<id) != 0 {
					rate += r
				}
			}
			second := []valveMinuteAction{{g.start, 0}}
			if two {
				second = actions(state.second, state.opened)
			}
			for _, a := range actions(state.first, state.opened) {
				for _, b := range second {
					if a.opening&b.opening != 0 {
						continue
					}
					key := valveMinuteState{a.position, b.position, state.opened | a.opening | b.opening}
					value := score + rate
					if old, ok := next[key]; !ok || value > old {
						next[key] = value
					}
				}
			}
		}
		states = next
	}
	var best uint64
	for _, score := range states {
		if score > best {
			best = score
		}
	}
	return best
}

func TestValveSearchAgainstMinuteActions(t *testing.T) {
	for code := 0; code < 384; code++ {
		size := 4
		if code >= 256 {
			size = 5
		}
		g := valveNetwork{rates: make([]uint64, size), edges: make([][]int, size)}
		for i := range g.rates {
			g.rates[i] = uint64((code*(i+3) + i*5) % 9)
		}
		if code < 256 {
			bit := 0
			for a := 0; a < size; a++ {
				for b := a + 1; b < size; b++ {
					if (code%64)&(1<<bit) != 0 {
						g.edges[a] = append(g.edges[a], b)
						g.edges[b] = append(g.edges[b], a)
					}
					bit++
				}
			}
		} else {
			seed := uint32(code)
			for a := 0; a < size; a++ {
				for b := 0; b < size; b++ {
					seed = seed*1664525 + 1013904223
					if seed>>30 == 0 {
						g.edges[a] = append(g.edges[a], b)
					}
				}
			}
		}
		for _, minutes := range []int{1, 4, 7} {
			scores, err := g.prepare(minutes).scores(minutes)
			if err != nil {
				t.Fatal(err)
			}
			var single uint64
			for _, score := range scores {
				if score > single {
					single = score
				}
			}
			pair, err := combineValveScores(scores)
			if err != nil {
				t.Fatal(err)
			}
			for _, two := range []bool{false, true} {
				got := single
				if two {
					got = pair
				}
				want := valveMinuteOracle(g, minutes, two)
				if got != want {
					t.Fatalf("case %d minutes %d two=%t: got %d want %d", code, minutes, two, got, want)
				}
			}
		}
	}
}

func TestValveGreedySplitRegression(t *testing.T) {
	input := "Valve AA has flow rate=0; tunnels lead to valves AB, AC, AD\n" +
		"Valve AB has flow rate=10; tunnel leads to valve AA\n" +
		"Valve AC has flow rate=9; tunnel leads to valve AA\n" +
		"Valve AD has flow rate=8; tunnel leads to valve AA"
	for i, solve := range []func(string) (string, error){(Day16{}).Part1, (Day16{}).Part2} {
		got, err := solve(input)
		want := []string{"681", "624"}[i]
		if err != nil || got != want {
			t.Fatalf("part %d: got %q %v want %s", i+1, got, err, want)
		}
	}
}

func TestValveParserAndBoundaries(t *testing.T) {
	valid := "Valve AA has flow rate=4; tunnel leads to valve AA"
	for _, input := range []string{valid, valid + "\r\n"} {
		for part, solve := range []func(string) (string, error){(Day16{}).Part1, (Day16{}).Part2} {
			got, err := solve(input)
			want := []string{"116", "100"}[part]
			if err != nil || got != want {
				t.Fatalf("positive AA part %d: %q %v", part+1, got, err)
			}
		}
	}
	for _, input := range []string{"", valid + "\n\n", valid + "\n" + valid,
		"Valve AB has flow rate=0; tunnel leads to valve AB",
		"Valve AA has flow rate=0; tunnel leads to valve ZZ",
		"Valve AA has flow rate=-1; tunnel leads to valve AA",
		"Valve AA has flow rate=18446744073709551616; tunnel leads to valve AA",
		"Valve AA has flow rate=1x; tunnel leads to valve AA",
		"Valve AA has flow rate=0; tunnels lead to valves ",
		"Valve A has flow rate=0; tunnel leads to valve AA"} {
		if _, err := (Day16{}).Part1(input); err == nil {
			t.Fatalf("accepted %q", input)
		}
	}
	zero := "Valve AA has flow rate=0; tunnel leads to valve AA\nValve ZZ has flow rate=99; tunnel leads to valve ZZ"
	for _, solve := range []func(string) (string, error){(Day16{}).Part1, (Day16{}).Part2} {
		if got, err := solve(zero); err != nil || got != "0" {
			t.Fatalf("unreachable valves: %q %v", got, err)
		}
	}
	g := valveNetwork{rates: make([]uint64, 71), edges: make([][]int, 71)}
	for i := 1; i < 71; i++ {
		g.rates[i] = uint64(i)
		g.edges[0] = append(g.edges[0], i)
	}
	scores, err := g.prepare(3).scores(3)
	if err != nil {
		t.Fatal(err)
	}
	if got, err := combineValveScores(scores); err != nil || got != 139 {
		t.Fatalf("wide masks: %d %v", got, err)
	}
}

func TestValvePressureOverflow(t *testing.T) {
	for _, rate := range []uint64{math.MaxUint64, math.MaxUint64/29 + 1} {
		input := fmt.Sprintf("Valve AA has flow rate=%d; tunnel leads to valve AA", rate)
		if _, err := (Day16{}).Part1(input); err == nil {
			t.Fatal("expected multiplication overflow")
		}
	}
	s := valveSearch{rates: []uint64{math.MaxUint64 / 3, math.MaxUint64 / 3}, travel: [][]int{{0, 0}, {0, 0}, {0, 0}}}
	if _, err := s.scores(4); err == nil {
		t.Fatal("expected accumulated overflow")
	}
	if _, err := combineValveScores(map[string]uint64{string([]byte{1}): math.MaxUint64, string([]byte{2}): 1}); err == nil {
		t.Fatal("expected combined overflow")
	}
	got, err := valvePressure("Valve AA has flow rate="+strconv.FormatUint(math.MaxUint64/29, 10)+"; tunnel leads to valve AA", 30, false)
	if err != nil || got != strconv.FormatUint((math.MaxUint64/29)*29, 10) {
		t.Fatalf("finite boundary: %s %v", got, err)
	}
}

func TestValvePairingAgainstAllPairs(t *testing.T) {
	for code := 0; code < 80; code++ {
		scores := map[string]uint64{}
		for mask := 0; mask < 64; mask++ {
			scores[string([]byte{byte(mask)})] = uint64((mask*37 + code*23 + mask*code) % 997)
		}
		var want uint64
		for a, sa := range scores {
			for b, sb := range scores {
				if a[0]&b[0] == 0 && sa+sb > want {
					want = sa + sb
				}
			}
		}
		got, err := combineValveScores(scores)
		if err != nil || got != want {
			t.Fatalf("pairing %d: %d %v want %d", code, got, err, want)
		}
	}
}

func TestValvePairingTreatsMasksAsBytes(t *testing.T) {
	scores := map[string]uint64{
		string([]byte{0, 0}):       0,
		string([]byte{0xc2, 0x80}): 100,
		string([]byte{0, 0x80}):    90,
	}
	if got, err := combineValveScores(scores); err != nil || got != 100 {
		t.Fatalf("byte masks: %d %v", got, err)
	}
}
