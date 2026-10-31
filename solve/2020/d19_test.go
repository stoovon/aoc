package solve2020

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"testing"
)

type monsterOracleRule struct {
	literal      string
	alternatives [][]int
}

func monsterOracleLanguage(rules map[int]monsterOracleRule, limit int) map[string]bool {
	words := make(map[int]map[string]bool)
	for id, rule := range rules {
		words[id] = make(map[string]bool)
		if len(rule.literal) > 0 && len(rule.literal) <= limit {
			words[id][rule.literal] = true
		}
	}
	changed := true
	for changed {
		changed = false
		for id, rule := range rules {
			for _, sequence := range rule.alternatives {
				prefixes := map[string]bool{"": true}
				for _, sub := range sequence {
					next := make(map[string]bool)
					for prefix := range prefixes {
						for suffix := range words[sub] {
							if len(prefix)+len(suffix) <= limit {
								next[prefix+suffix] = true
							}
						}
					}
					prefixes = next
				}
				for word := range prefixes {
					if !words[id][word] {
						words[id][word] = true
						changed = true
					}
				}
			}
		}
	}
	return words[0]
}
func monsterOracleInput(rules map[int]monsterOracleRule, messages []string) string {
	ids := make([]int, 0, len(rules))
	for id := range rules {
		ids = append(ids, id)
	}
	sort.Ints(ids)
	var source strings.Builder
	for _, id := range ids {
		fmt.Fprintf(&source, "%d: ", id)
		rule := rules[id]
		if rule.literal != "" {
			source.WriteString(strconv.Quote(rule.literal))
		} else {
			for i, sequence := range rule.alternatives {
				if i > 0 {
					source.WriteString(" | ")
				}
				for j, sub := range sequence {
					if j > 0 {
						source.WriteByte(' ')
					}
					source.WriteString(strconv.Itoa(sub))
				}
			}
		}
		source.WriteByte('\n')
	}
	source.WriteByte('\n')
	source.WriteString(strings.Join(messages, "\n"))
	return source.String()
}
func monsterBinaryWords(limit int) []string {
	words := []string{}
	level := []string{""}
	for length := 1; length <= limit; length++ {
		next := []string{}
		for _, prefix := range level {
			for _, letter := range []string{"x", "y"} {
				next = append(next, prefix+letter)
			}
		}
		words = append(words, next...)
		level = next
	}
	return words
}

func TestMonsterMessagesAgainstBoundedLanguageOracle(t *testing.T) {
	messages := monsterBinaryWords(5)
	options := [][]int{{3}, {1, 2}, {2, 1}, {1, 1}, {4}}
	for mask := 0; mask < 1<<len(options); mask++ {
		alternatives := [][]int{{2}}
		for i, option := range options {
			if mask&(1<<i) != 0 {
				alternatives = append(alternatives, option)
			}
		}
		rules := map[int]monsterOracleRule{0: {alternatives: [][]int{{1}}}, 1: {alternatives: alternatives}, 2: {literal: "x"}, 3: {literal: "y"}, 4: {alternatives: [][]int{{1}}}}
		input := monsterOracleInput(rules, messages)
		want := monsterOracleLanguage(rules, 5)
		grammar, _, err := parseMonsterMessages(input)
		if err != nil {
			t.Fatal(err)
		}
		widths := grammar.fixedWidths()
		matcher := fixedMessageMatcher{grammar, widths, make(map[fixedMessageKey]bool)}
		count := 0
		for _, message := range messages {
			expected := want[message]
			if expected {
				count++
			}
			if got := grammar.recognizes(message); got != expected {
				t.Fatalf("mask=%d message=%s chart=%t want=%t", mask, message, got, expected)
			}
			if widths[grammar.index[0]] > 0 && matcher.matches(grammar.index[0], message) != expected {
				t.Fatal("fixed-width recognizer mismatch")
			}
		}
		got, err := (Day19{}).Part1(input)
		if err != nil || got != strconv.Itoa(count) {
			t.Fatalf("mask=%d part1=%s err=%v want=%d", mask, got, err, count)
		}
	}
}

func TestMonsterMessagesRecursiveRulesAgainstLanguageOracle(t *testing.T) {
	messages := monsterBinaryWords(7)
	leftOptions := [][][]int{{{1}}, {{1, 1}}, {{1}, {1, 1}}}
	rightOptions := [][][]int{{{1}}, {{2}}, {{1, 2}}, {{2}, {2, 2}}}
	for _, left := range leftOptions {
		for _, right := range rightOptions {
			rules := map[int]monsterOracleRule{0: {alternatives: [][]int{{8, 11}}}, 8: {alternatives: [][]int{{42}}}, 11: {alternatives: [][]int{{42, 31}}}, 42: {alternatives: left}, 31: {alternatives: right}, 1: {literal: "x"}, 2: {literal: "y"}}
			input := monsterOracleInput(rules, messages)
			rules[8] = monsterOracleRule{alternatives: [][]int{{42}, {42, 8}}}
			rules[11] = monsterOracleRule{alternatives: [][]int{{42, 31}, {42, 11, 31}}}
			want := monsterOracleLanguage(rules, 7)
			count := 0
			for _, message := range messages {
				if want[message] {
					count++
				}
			}
			got, err := (Day19{}).Part2(input)
			if err != nil || got != strconv.Itoa(count) {
				t.Fatalf("left=%v right=%v part2=%s err=%v want=%d", left, right, got, err, count)
			}
			patched := monsterOracleInput(rules, messages)
			grammar, _, err := parseMonsterMessages(patched)
			if err != nil {
				t.Fatal(err)
			}
			widths := grammar.fixedWidths()
			matcher := fixedMessageMatcher{grammar, widths, make(map[fixedMessageKey]bool)}
			l, r := grammar.index[42], grammar.index[31]
			for _, message := range messages {
				expected := want[message]
				if grammar.recognizes(message) != expected {
					t.Fatal("recursive chart mismatch", left, right, message)
				}
				if widths[l] > 0 && widths[l] == widths[r] && matcher.matchesRecursiveChunks(message, widths[l], l, r) != expected {
					t.Fatal("verified chunk matcher mismatch", left, right, message)
				}
			}
		}
	}
}

func TestMonsterMessagesRegressionCases(t *testing.T) {
	cases := []struct{ input, want string }{
		{"0: \"xy\"\n\nxy\nx\nxyx", "1"},
		{"0: 0 1 | 1\n1: \"x\"\n\nx\nxx\nxxx\ny", "3"},
		{"0: 1\n1: 0 | 2\n2: \"x\"\n\nx\nxx", "1"},
		{"0: 1\n1: 0\n\nx", "0"},
		{"0: 1 1\n1: 1 1 | 2\n2: \"x\"\n\nxx\nxxx\nxxxx\nx", "3"},
		{"0: 1 | 1\n1: \"x\"\n\nx\nx\nxx", "2"},
		{"0: 1000000000\n1000000000: \"x\"\n\nx\ny", "1"},
		{"0: 1 | 1 0\n1: \"x\"\n\n" + strings.Repeat("x", 80), "1"},
	}
	for _, c := range cases {
		for _, input := range []string{c.input, strings.ReplaceAll(c.input, "\n", "\r\n") + "\r\n"} {
			got, err := (Day19{}).Part1(input)
			if err != nil || got != c.want {
				t.Fatal(got, err, c.want)
			}
		}
	}
	overlap := "0: 8 11\n8: 42\n11: 42 31\n42: 1\n31: 1\n1: \"x\"\n\nxxx\nxxxx\nxx\nx"
	if got, err := (Day19{}).Part2(overlap); err != nil || got != "2" {
		t.Fatal(got, err)
	}
	variable := "0: 8 11\n8: 42\n11: 42 31\n42: 1 | 1 1\n31: 2\n1: \"x\"\n2: \"y\"\n\nxxy\nxxxy\nxxxyy\nxxyy\nxy"
	if got, err := (Day19{}).Part2(variable); err != nil || got != "3" {
		t.Fatal(got, err)
	}
	noncanonical := "0: 42\n8: 42\n11: 42 31\n42: \"x\"\n31: \"y\"\n\nx\nxxy"
	if got, err := (Day19{}).Part2(noncanonical); err != nil || got != "1" {
		t.Fatal("assumed rule-zero shape", got, err)
	}
}

func TestMonsterMessagesRejectInvalidInput(t *testing.T) {
	invalid := []string{"", "0: \"x\"", "0: \"x\"\n\n", "0: \"\"\n\nx", "0: \"x\" | 1\n1: \"y\"\n\nx", "0: 1\n\nx", "0: 1 |\n1: \"x\"\n\nx", "0: \"x\"\n0: \"y\"\n\nx", "1: \"x\"\n\nx", "-1: \"x\"\n0: 1\n\nx", "+0: \"x\"\n\nx", "0: +1\n1: \"x\"\n\nx", "0: \"海\"\n\nx", "0: \"x\"\n\nx\n\ny", "0: \"x\"\n\nx y", "0: \"x\"\n\n1"}
	for _, input := range invalid {
		for _, part := range []func(string) (string, error){(Day19{}).Part1, (Day19{}).Part2} {
			if _, err := part(input); err == nil {
				t.Fatalf("accepted %q", input)
			}
		}
	}
	if _, err := (Day19{}).Part2("0: \"x\"\n\nx"); err == nil || !strings.Contains(err.Error(), "requires rule") {
		t.Fatal(err)
	}
	if _, err := (Day19{}).Part1("0: 1\n1: |\n\nx"); err == nil || !strings.Contains(err.Error(), "line 2:") {
		t.Fatal(err)
	}
}

func TestMonsterMessagesAvoidWholeLanguageGeneration(t *testing.T) {
	rules := map[int]monsterOracleRule{0: {alternatives: [][]int{{8, 11}}}, 8: {alternatives: [][]int{{42}}}, 11: {alternatives: [][]int{{42, 31}}}, 42: {alternatives: [][]int{{100}}}, 31: {alternatives: [][]int{{100}}}, 110: {alternatives: [][]int{{1}, {2}}}, 1: {literal: "x"}, 2: {literal: "y"}}
	for id := 100; id < 110; id++ {
		rules[id] = monsterOracleRule{alternatives: [][]int{{id + 1, id + 1}}}
	}
	input := monsterOracleInput(rules, []string{"x", "xy", "xxxy"})
	for _, part := range []func(string) (string, error){(Day19{}).Part1, (Day19{}).Part2} {
		if got, err := part(input); err != nil || got != "0" {
			t.Fatal(got, err)
		}
	}
	overflow := map[int]monsterOracleRule{70: {literal: "x"}}
	for id := 0; id < 70; id++ {
		overflow[id] = monsterOracleRule{alternatives: [][]int{{id + 1, id + 1}}}
	}
	if got, err := (Day19{}).Part1(monsterOracleInput(overflow, []string{"x"})); err != nil || got != "0" {
		t.Fatal("overflowing width must use chart fallback", got, err)
	}
}
