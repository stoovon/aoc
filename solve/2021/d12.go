package solve2021

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day12 struct {
}

func (d Day12) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2021, Day: 12}
}

type state struct {
	Pos  string
	Seen map[string]int
}

func (d Day12) parse(input string) map[string]map[string]struct{} {
	caves := map[string]map[string]struct{}{}
	for _, s := range strings.Fields(string(input)) {
		s := strings.Split(s, "-")

		for a, b := range []int{1, 0} {
			if caves[s[a]] == nil {
				caves[s[a]] = map[string]struct{}{}
			}
			caves[s[a]][s[b]] = struct{}{}
		}
	}

	return caves
}

func (d Day12) run(caves map[string]map[string]struct{}, part1 bool) (count int) {
	queue := []state{{"start", map[string]int{"start": 1}}}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur.Pos == "end" {
			count++
			continue
		}

	out:
		for c := range caves[cur.Pos] {
			if c == "start" {
				continue
			}

			seen := map[string]int{}
			for k, v := range cur.Seen {
				seen[k] = v
				if (part1 || v == 2) && cur.Seen[c] > 0 {
					continue out
				}
			}

			if c == strings.ToLower(c) {
				seen[c]++
			}

			queue = append(queue, state{c, seen})
		}
	}
	return
}

func (d Day12) Part1(input string) (string, error) {
	caves := d.parse(input)
	count := d.run(caves, true)
	return strconv.Itoa(count), nil
}

func (d Day12) Part2(input string) (string, error) {
	caves := d.parse(input)
	count := d.run(caves, false)
	return strconv.Itoa(count), nil
}

func init() {
	solve.Register(Day12{})
}
