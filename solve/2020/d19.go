package solve2020

import (
	"aoc/solve"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

type Day19 struct{}

func (d Day19) Coords() solve.SolutionCoords { return solve.SolutionCoords{Year: 2020, Day: 19} }

type messageRule struct {
	literal      string
	alternatives [][]int
}
type messageGrammar struct {
	rules []messageRule
	index map[int]int
}

func messageWord(text string) bool {
	if text == "" {
		return false
	}
	for _, c := range []byte(text) {
		if !(c >= 'a' && c <= 'z') && !(c >= 'A' && c <= 'Z') {
			return false
		}
	}
	return true
}

func messageRuleID(text string) (int, error) {
	if text == "" {
		return 0, fmt.Errorf("expected nonnegative rule ID")
	}
	for _, c := range []byte(text) {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("invalid rule ID %q", text)
		}
	}
	id, err := strconv.Atoi(text)
	if err != nil {
		return 0, fmt.Errorf("rule ID out of range: %q", text)
	}
	return id, nil
}

func parseMonsterMessages(input string) (messageGrammar, []string, error) {
	input = strings.TrimSuffix(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	source, messagesText, ok := strings.Cut(input, "\n\n")
	if !ok || source == "" || messagesText == "" {
		return messageGrammar{}, nil, fmt.Errorf("expected rules and messages separated by a blank line")
	}
	lines := strings.Split(source, "\n")
	raw := make(map[int]messageRule)
	for number, line := range lines {
		parseLine := func() error {
			idText, body, ok := strings.Cut(line, ":")
			if !ok {
				return fmt.Errorf("expected rule ID: definition")
			}
			id, err := messageRuleID(strings.TrimSpace(idText))
			if err != nil {
				return err
			}
			if _, exists := raw[id]; exists {
				return fmt.Errorf("duplicate rule %d", id)
			}
			body = strings.TrimSpace(body)
			rule := messageRule{}
			if strings.HasPrefix(body, "\"") {
				literal, err := strconv.Unquote(body)
				if err != nil || !messageWord(literal) {
					return fmt.Errorf("expected a quoted nonempty alphabetic literal")
				}
				rule.literal = literal
			} else {
				for _, alternative := range strings.Split(body, "|") {
					fields := strings.Fields(alternative)
					if len(fields) == 0 {
						return fmt.Errorf("expected nonempty rule alternative")
					}
					sequence := make([]int, len(fields))
					for i, field := range fields {
						sequence[i], err = messageRuleID(field)
						if err != nil {
							return err
						}
					}
					rule.alternatives = append(rule.alternatives, sequence)
				}
			}
			raw[id] = rule
			return nil
		}
		if err := parseLine(); err != nil {
			return messageGrammar{}, nil, fmt.Errorf("line %d: %w", number+1, err)
		}
	}
	if _, ok := raw[0]; !ok {
		return messageGrammar{}, nil, fmt.Errorf("missing rule 0")
	}
	ids := make([]int, 0, len(raw))
	for id := range raw {
		ids = append(ids, id)
	}
	sort.Ints(ids)
	grammar := messageGrammar{rules: make([]messageRule, len(ids)), index: make(map[int]int, len(ids))}
	for index, id := range ids {
		grammar.index[id] = index
	}
	for index, id := range ids {
		rule := raw[id]
		for _, sequence := range rule.alternatives {
			for i, reference := range sequence {
				mapped, ok := grammar.index[reference]
				if !ok {
					return messageGrammar{}, nil, fmt.Errorf("rule %d references missing rule %d", id, reference)
				}
				sequence[i] = mapped
			}
		}
		grammar.rules[index] = rule
	}
	messages := strings.Split(messagesText, "\n")
	for i, message := range messages {
		if !messageWord(message) {
			return messageGrammar{}, nil, fmt.Errorf("line %d: expected a nonempty alphabetic message", len(lines)+2+i)
		}
	}
	return grammar, messages, nil
}

func (g messageGrammar) fixedWidths() []int {
	widths := make([]int, len(g.rules))
	state := make([]uint8, len(g.rules))
	var width func(int) int
	width = func(id int) int {
		if state[id] == 1 {
			return 0
		}
		if state[id] == 2 {
			return widths[id]
		}
		state[id] = 1
		rule := g.rules[id]
		result := len(rule.literal)
		if rule.literal == "" {
			for i, sequence := range rule.alternatives {
				total := 0
				for _, sub := range sequence {
					w := width(sub)
					if w == 0 || w > math.MaxInt-total {
						total = 0
						break
					}
					total += w
				}
				if total == 0 || (i > 0 && total != result) {
					result = 0
					break
				}
				result = total
			}
		}
		state[id] = 2
		widths[id] = result
		return result
	}
	for id := range g.rules {
		width(id)
	}
	return widths
}

type fixedMessageKey struct {
	rule int
	text string
}
type fixedMessageMatcher struct {
	grammar messageGrammar
	widths  []int
	memo    map[fixedMessageKey]bool
}

func (m fixedMessageMatcher) matches(id int, text string) bool {
	if len(text) != m.widths[id] {
		return false
	}
	rule := m.grammar.rules[id]
	if rule.literal != "" {
		return text == rule.literal
	}
	key := fixedMessageKey{id, text}
	if value, ok := m.memo[key]; ok {
		return value
	}
	matched := false
	for _, sequence := range rule.alternatives {
		offset, valid := 0, true
		for _, sub := range sequence {
			end := offset + m.widths[sub]
			if !m.matches(sub, text[offset:end]) {
				valid = false
				break
			}
			offset = end
		}
		if valid {
			matched = true
			break
		}
	}
	m.memo[key] = matched
	return matched
}

func (m fixedMessageMatcher) matchesRecursiveChunks(text string, width, left, right int) bool {
	if len(text)%width != 0 {
		return false
	}
	count := len(text) / width
	if count < 3 {
		return false
	}
	prefix := 0
	for prefix < count-1 && m.matches(left, text[prefix*width:(prefix+1)*width]) {
		prefix++
	}
	if prefix <= count-prefix {
		return false
	}
	for i := prefix; i < count; i++ {
		if !m.matches(right, text[i*width:(i+1)*width]) {
			return false
		}
	}
	return true
}

type messageItem struct{ rule, alternative, dot, start int }
type messageChart struct {
	items   []messageItem
	seen    map[messageItem]struct{}
	waiting map[int][]messageItem
}

func (c *messageChart) add(item messageItem) {
	if c.seen == nil {
		c.seen = make(map[messageItem]struct{})
	}
	if _, ok := c.seen[item]; ok {
		return
	}
	c.seen[item] = struct{}{}
	c.items = append(c.items, item)
}
func (g messageGrammar) recognizes(text string) bool {
	root := g.index[0]
	if g.rules[root].literal != "" {
		return text == g.rules[root].literal
	}
	chart := make([]messageChart, len(text)+1)
	for alternative := range g.rules[root].alternatives {
		chart[0].add(messageItem{rule: root, alternative: alternative})
	}
	for position := range chart {
		for cursor := 0; cursor < len(chart[position].items); cursor++ {
			item := chart[position].items[cursor]
			sequence := g.rules[item.rule].alternatives[item.alternative]
			if item.dot == len(sequence) {
				if item.rule == root && item.start == 0 && position == len(text) {
					return true
				}
				for _, parent := range chart[item.start].waiting[item.rule] {
					parent.dot++
					chart[position].add(parent)
				}
				continue
			}
			next := sequence[item.dot]
			rule := g.rules[next]
			if rule.literal != "" {
				if strings.HasPrefix(text[position:], rule.literal) {
					item.dot++
					chart[position+len(rule.literal)].add(item)
				}
			} else {
				if chart[position].waiting == nil {
					chart[position].waiting = make(map[int][]messageItem)
				}
				chart[position].waiting[next] = append(chart[position].waiting[next], item)
				for alternative := range rule.alternatives {
					chart[position].add(messageItem{rule: next, alternative: alternative, start: position})
				}
			}
		}
	}
	return false
}

func solveMonsterMessages(input string, recursive bool) (string, error) {
	grammar, messages, err := parseMonsterMessages(input)
	if err != nil {
		return "", err
	}
	if recursive {
		for _, id := range []int{8, 11, 42, 31} {
			if _, ok := grammar.index[id]; !ok {
				return "", fmt.Errorf("part two requires rule %d", id)
			}
		}
		r8, r11, r42, r31 := grammar.index[8], grammar.index[11], grammar.index[42], grammar.index[31]
		grammar.rules[r8] = messageRule{alternatives: [][]int{{r42}, {r42, r8}}}
		grammar.rules[r11] = messageRule{alternatives: [][]int{{r42, r31}, {r42, r11, r31}}}
	}
	widths := grammar.fixedWidths()
	matcher := fixedMessageMatcher{grammar, widths, make(map[fixedMessageKey]bool)}
	root := grammar.index[0]
	chunkWidth, left, right := 0, 0, 0
	if recursive {
		r8, r11 := grammar.index[8], grammar.index[11]
		left, right = grammar.index[42], grammar.index[31]
		alternatives := grammar.rules[root].alternatives
		if len(alternatives) == 1 && len(alternatives[0]) == 2 && alternatives[0][0] == r8 && alternatives[0][1] == r11 && widths[left] > 0 && widths[left] == widths[right] {
			chunkWidth = widths[left]
		}
	}
	count := 0
	for _, message := range messages {
		var matched bool
		if widths[root] > 0 {
			matched = matcher.matches(root, message)
		} else if chunkWidth > 0 {
			matched = matcher.matchesRecursiveChunks(message, chunkWidth, left, right)
		} else {
			matched = grammar.recognizes(message)
		}
		if matched {
			count++
		}
	}
	return strconv.Itoa(count), nil
}
func (d Day19) Part1(input string) (string, error) { return solveMonsterMessages(input, false) }
func (d Day19) Part2(input string) (string, error) { return solveMonsterMessages(input, true) }
func init()                                        { solve.Register(Day19{}) }
