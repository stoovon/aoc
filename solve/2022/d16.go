package solve2022

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day16 struct{}

func (d Day16) Coords() solve.SolutionCoords { return solve.SolutionCoords{Year: 2022, Day: 16} }

var valveRecord = regexp.MustCompile(`^Valve ([A-Z]{2}) has flow rate=([0-9]+); (?:tunnels lead to valves ([A-Z]{2}(?:, [A-Z]{2})*)|tunnel leads to valve ([A-Z]{2}))$`)

type valveNetwork struct {
	rates []uint64
	edges [][]int
	start int
}

func parseValveNetwork(input string) (valveNetwork, error) {
	var g valveNetwork
	ids := map[string]int{}
	targets := [][]string{}
	for line, record := range strings.Split(strings.TrimSuffix(input, "\n"), "\n") {
		record = strings.TrimSuffix(record, "\r")
		match := valveRecord.FindStringSubmatch(record)
		if match == nil {
			return g, fmt.Errorf("line %d: invalid valve record", line+1)
		}
		if _, exists := ids[match[1]]; exists {
			return g, fmt.Errorf("line %d: duplicate valve %s", line+1, match[1])
		}
		rate, err := strconv.ParseUint(match[2], 10, 64)
		if err != nil {
			return g, fmt.Errorf("line %d: invalid flow rate: %w", line+1, err)
		}
		ids[match[1]] = len(g.rates)
		g.rates = append(g.rates, rate)
		names := match[3]
		if names == "" {
			names = match[4]
		}
		targets = append(targets, strings.Split(names, ", "))
	}
	start, ok := ids["AA"]
	if !ok {
		return g, fmt.Errorf("missing starting valve AA")
	}
	g.start = start
	g.edges = make([][]int, len(g.rates))
	for from, names := range targets {
		seen := map[int]bool{}
		for _, name := range names {
			to, ok := ids[name]
			if !ok {
				return g, fmt.Errorf("line %d: undefined valve %s", from+1, name)
			}
			if !seen[to] {
				seen[to] = true
				g.edges[from] = append(g.edges[from], to)
			}
		}
	}
	return g, nil
}

func (g valveNetwork) distances(from int) []int {
	dist := make([]int, len(g.rates))
	for i := range dist {
		dist[i] = -1
	}
	dist[from] = 0
	queue := []int{from}
	for head := 0; head < len(queue); head++ {
		for _, to := range g.edges[queue[head]] {
			if dist[to] < 0 {
				dist[to] = dist[queue[head]] + 1
				queue = append(queue, to)
			}
		}
	}
	return dist
}

type valveSearch struct {
	rates  []uint64
	travel [][]int
}

func (g valveNetwork) prepare(minutes int) valveSearch {
	initial := g.distances(g.start)
	useful := []int{}
	for id, rate := range g.rates {
		if rate > 0 && initial[id] >= 0 && initial[id]+1 < minutes {
			useful = append(useful, id)
		}
	}
	sort.Slice(useful, func(i, j int) bool { return g.rates[useful[i]] > g.rates[useful[j]] })
	var search valveSearch
	for _, id := range useful {
		search.rates = append(search.rates, g.rates[id])
	}
	for _, from := range append(useful, g.start) {
		distances := g.distances(from)
		row := make([]int, len(useful))
		for i, to := range useful {
			row[i] = distances[to]
		}
		search.travel = append(search.travel, row)
	}
	return search
}

type valveSearchState struct {
	current, timeLeft int
	opened            string
}

func (s valveSearch) scores(minutes int) (map[string]uint64, error) {
	best := map[string]uint64{}
	arrivals := map[valveSearchState]uint64{}
	var visit func(valveSearchState, uint64) error
	visit = func(state valveSearchState, score uint64) error {
		if previous, ok := arrivals[state]; ok && previous >= score {
			return nil
		}
		arrivals[state] = score
		if previous, ok := best[state.opened]; !ok || score > previous {
			best[state.opened] = score
		}
		for next, rate := range s.rates {
			if state.opened[next/8]&(1<<(next%8)) != 0 {
				continue
			}
			distance := s.travel[state.current][next]
			remaining := state.timeLeft - distance - 1
			if distance < 0 || remaining <= 0 {
				continue
			}
			if rate > math.MaxUint64/uint64(remaining) {
				return fmt.Errorf("pressure exceeds uint64")
			}
			gain := rate * uint64(remaining)
			if score > math.MaxUint64-gain {
				return fmt.Errorf("pressure exceeds uint64")
			}
			opened := []byte(state.opened)
			opened[next/8] |= 1 << (next % 8)
			if err := visit(valveSearchState{next, remaining, string(opened)}, score+gain); err != nil {
				return err
			}
		}
		return nil
	}
	err := visit(valveSearchState{len(s.rates), minutes, string(make([]byte, (len(s.rates)+7)/8))}, 0)
	return best, err
}

type valveSetScore struct {
	opened string
	score  uint64
}

func combineValveScores(scores map[string]uint64) (uint64, error) {
	sets := make([]valveSetScore, 0, len(scores))
	for opened, score := range scores {
		sets = append(sets, valveSetScore{opened, score})
	}
	sort.Slice(sets, func(i, j int) bool { return sets[i].score > sets[j].score })
	var best uint64
	for i, a := range sets {
		if a.score <= best/2 {
			break
		}
		for _, b := range sets[i:] {
			if a.score <= best && b.score <= best-a.score {
				break
			}
			disjoint := true
			for k := 0; k < len(a.opened); k++ {
				if a.opened[k]&b.opened[k] != 0 {
					disjoint = false
					break
				}
			}
			if !disjoint {
				continue
			}
			if a.score > math.MaxUint64-b.score {
				return 0, fmt.Errorf("combined pressure exceeds uint64")
			}
			best = a.score + b.score
			break
		}
	}
	return best, nil
}

func valvePressure(input string, minutes int, twoWorkers bool) (string, error) {
	network, err := parseValveNetwork(input)
	if err != nil {
		return "", err
	}
	scores, err := network.prepare(minutes).scores(minutes)
	if err != nil {
		return "", err
	}
	var best uint64
	if twoWorkers {
		best, err = combineValveScores(scores)
	} else {
		for _, score := range scores {
			if score > best {
				best = score
			}
		}
	}
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(best, 10), nil
}

func (d Day16) Part1(input string) (string, error) { return valvePressure(input, 30, false) }
func (d Day16) Part2(input string) (string, error) { return valvePressure(input, 26, true) }
func init()                                        { solve.Register(Day16{}) }
