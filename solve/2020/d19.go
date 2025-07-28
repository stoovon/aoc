package solve2020

import (
	"aoc/solve"
	"errors"
	"strconv"
	"strings"
)

type Day19 struct{}

func (d Day19) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2020, Day: 19}
}

func (d Day19) Part1(input string) (string, error) {
	rules, messages, err := parseInput19(input)
	if err != nil {
		return "", err
	}
	count := 0
	for _, msg := range messages {
		if matchRule0(rules, msg) {
			count++
		}
	}
	return strconv.Itoa(count), nil
}

type rule19 struct {
	char string
	alts [][]int
}

func matchRule0(rules map[int]rule19, msg string) bool {
	ends := matchRule(rules, 0, msg)
	for _, e := range ends {
		if e == len(msg) {
			return true
		}
	}
	return false
}

// Returns all possible end positions after matching rule at pos 0
func matchRule(rules map[int]rule19, idx int, msg string) []int {
	rule := rules[idx]
	if rule.char != "" {
		if strings.HasPrefix(msg, rule.char) {
			return []int{1}
		}
		return nil
	}
	var ends []int
	for _, seq := range rule.alts {
		positions := []int{0}
		for _, sub := range seq {
			var newPos []int
			for _, p := range positions {
				res := matchRule(rules, sub, msg[p:])
				for _, r := range res {
					newPos = append(newPos, p+r)
				}
			}
			positions = newPos
			if len(positions) == 0 {
				break
			}
		}
		ends = append(ends, positions...)
	}
	return ends
}

func (d Day19) Part2(input string) (string, error) {
	rules, messages, err := parseInput19(input)
	if err != nil {
		return "", err
	}
	// Patch rules 8 and 11 for recursion
	rules[8] = rule19{alts: [][]int{{42}, {42, 8}}}
	rules[11] = rule19{alts: [][]int{{42, 31}, {42, 11, 31}}}
	// Precompute all possible matches for rule 42 and 31
	rule42matches := d.allMatches(rules, 42)
	rule31matches := d.allMatches(rules, 31)
	segLen := 0
	for m := range rule42matches {
		if segLen == 0 {
			segLen = len(m)
		}
		if len(m) != segLen {
			return "", errors.New("rule 42 matches have inconsistent lengths")
		}
	}
	count := 0
	for _, msg := range messages {
		if d.matchPart2(msg, rule42matches, rule31matches, segLen) {
			count++
		}
	}
	return strconv.Itoa(count), nil
}

// Shared input parsing for both parts
func parseInput19(input string) (map[int]rule19, []string, error) {
	rulesSection, messagesSection, found := strings.Cut(input, "\n\n")
	if !found {
		return nil, nil, errors.New("input missing rules/messages section")
	}
	rules := parseRules19(strings.Split(strings.TrimSpace(rulesSection), "\n"))
	messages := strings.Split(strings.TrimSpace(messagesSection), "\n")
	return rules, messages, nil
}

// Standalone for input parsing
func parseRules19(lines []string) map[int]rule19 {
	rules := make(map[int]rule19)
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, ": ", 2)
		idx, _ := strconv.Atoi(parts[0])
		if strings.Contains(parts[1], "\"") {
			rules[idx] = rule19{char: strings.Trim(parts[1], "\"")}
		} else {
			var alts [][]int
			for _, alt := range strings.Split(parts[1], "|") {
				var seq []int
				for _, n := range strings.Fields(alt) {
					v, _ := strconv.Atoi(n)
					seq = append(seq, v)
				}
				alts = append(alts, seq)
			}
			rules[idx] = rule19{alts: alts}
		}
	}
	return rules
}

// Generate all possible matches for a rule (should only be used for rules 42 and 31)
func (d Day19) allMatches(rules map[int]rule19, idx int) map[string]struct{} {
	var gen func(int) []string
	memo := make(map[int][]string)
	gen = func(i int) []string {
		if v, ok := memo[i]; ok {
			return v
		}
		rule := rules[i]
		if rule.char != "" {
			memo[i] = []string{rule.char}
			return memo[i]
		}
		var out []string
		for _, alt := range rule.alts {
			seq := []string{""}
			for _, sub := range alt {
				var next []string
				for _, s := range seq {
					for _, t := range gen(sub) {
						next = append(next, s+t)
					}
				}
				seq = next
			}
			out = append(out, seq...)
		}
		memo[i] = out
		return out
	}
	res := make(map[string]struct{})
	for _, s := range gen(idx) {
		res[s] = struct{}{}
	}
	return res
}

// Check if a message matches rule 0 with recursive rules 8 and 11
func (d Day19) matchPart2(msg string, rule42, rule31 map[string]struct{}, segLen int) bool {
	n := len(msg)
	if n%segLen != 0 {
		return false
	}
	var segs []string
	for i := 0; i < n; i += segLen {
		segs = append(segs, msg[i:i+segLen])
	}
	// Count how many segments from the start match rule 42
	i := 0
	for i < len(segs) && contains(rule42, segs[i]) {
		i++
	}
	// At least two rule 42 matches (one for rule 8, one for rule 11)
	n42 := i
	// Count how many segments from the end match rule 31
	j := len(segs) - 1
	for j >= 0 && contains(rule31, segs[j]) {
		j--
	}
	n31 := len(segs) - 1 - j
	// There must be at least one rule 31, more rule 42 than rule 31, and all segments used
	return n31 >= 1 && n42 > n31 && n42+n31 == len(segs)
}

func contains(set map[string]struct{}, s string) bool {
	_, ok := set[s]
	return ok
}

func init() {
	solve.Register(Day19{})
}
