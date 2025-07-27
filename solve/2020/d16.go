package solve2020

import (
	"aoc/solve"
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type Day16 struct{}

func (d Day16) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2020, Day: 16}
}

func (d Day16) Part1(input string) (string, error) {
	sections := strings.Split(strings.TrimSpace(input), "\n\n")
	if len(sections) < 3 {
		return "", errors.New("input does not have three sections")
	}

	rules := parseRules(sections[0])
	validSet := make(map[int]bool)
	for _, r := range rules {
		for _, rg := range r.ranges {
			for i := rg[0]; i <= rg[1]; i++ {
				validSet[i] = true
			}
		}
	}

	lines := strings.Split(sections[2], "\n")
	var errorRate int
	for _, line := range lines[1:] {
		if strings.TrimSpace(line) == "" {
			continue
		}
		nums := strings.Split(line, ",")
		for _, n := range nums {
			v, err := strconv.Atoi(n)
			if err != nil {
				continue
			}
			if !validSet[v] {
				errorRate += v
			}
		}
	}
	return strconv.Itoa(errorRate), nil
}

type rule struct {
	name   string
	ranges [][2]int
}

func parseRules(rulesSection string) []rule {
	ruleRe := regexp.MustCompile(`([a-z ]+): (\d+)-(\d+) or (\d+)-(\d+)`)
	var rules []rule
	for _, line := range strings.Split(rulesSection, "\n") {
		m := ruleRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		a1, _ := strconv.Atoi(m[2])
		a2, _ := strconv.Atoi(m[3])
		b1, _ := strconv.Atoi(m[4])
		b2, _ := strconv.Atoi(m[5])
		rules = append(rules, rule{
			name: m[1],
			ranges: [][2]int{{a1, a2}, {b1, b2}},
		})
	}
	return rules
}

func (d Day16) Part2(input string) (string, error) {
	sections := strings.Split(strings.TrimSpace(input), "\n\n")
	if len(sections) < 3 {
		return "", errors.New("input does not have three sections")
	}

	// Parse field rules
	ruleRe := regexp.MustCompile(`([a-z ]+): (\d+)-(\d+) or (\d+)-(\d+)`)
	var rules []rule
	for _, line := range strings.Split(sections[0], "\n") {
		m := ruleRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		a1, _ := strconv.Atoi(m[2])
		a2, _ := strconv.Atoi(m[3])
		b1, _ := strconv.Atoi(m[4])
		b2, _ := strconv.Atoi(m[5])
		rules = append(rules, rule{
			name: m[1],
			ranges: [][2]int{{a1, a2}, {b1, b2}},
		})
	}

	// Parse your ticket
	yourTicketLines := strings.Split(sections[1], "\n")
	yourTicketStr := yourTicketLines[1]
	yourTicket := parseTicket(yourTicketStr)

	// Parse nearby tickets
	var validTickets [][]int
	for _, line := range strings.Split(sections[2], "\n")[1:] {
		if strings.TrimSpace(line) == "" {
			continue
		}
		ticket := parseTicket(line)
		if isValidTicket(ticket, rules) {
			validTickets = append(validTickets, ticket)
		}
	}
	validTickets = append(validTickets, yourTicket) // include your ticket for field deduction

	// For each field position, determine possible rules
	numFields := len(yourTicket)
	possible := make([]map[string]bool, numFields)
	for i := 0; i < numFields; i++ {
		possible[i] = make(map[string]bool)
		for _, r := range rules {
			possible[i][r.name] = true
		}
	}
	for _, ticket := range validTickets {
		for i, v := range ticket {
			for _, r := range rules {
				if !inAnyRange(v, r.ranges) {
					possible[i][r.name] = false
				}
			}
		}
	}

	// Deduce field positions
	fieldOrder := make([]string, numFields)
	decided := make(map[string]bool)
	for assigned := 0; assigned < numFields; {
		for i := 0; i < numFields; i++ {
			count := 0
			var candidate string
			for name, ok := range possible[i] {
				if ok && !decided[name] {
					count++
					candidate = name
				}
			}
			if count == 1 {
				fieldOrder[i] = candidate
				decided[candidate] = true
				assigned++
			}
		}
	}

	// Multiply values for fields starting with "departure"
	prod := 1
	found := false
	for i, name := range fieldOrder {
		if strings.HasPrefix(name, "departure") {
			prod *= yourTicket[i]
			found = true
		}
	}
	if !found {
		return "0", nil // no departure fields
	}
	return strconv.Itoa(prod), nil
}

func parseTicket(line string) []int {
	nums := strings.Split(line, ",")
	var ticket []int
	for _, n := range nums {
		v, err := strconv.Atoi(strings.TrimSpace(n))
		if err == nil {
			ticket = append(ticket, v)
		}
	}
	return ticket
}

func isValidTicket(ticket []int, rules []rule) bool {
	for _, v := range ticket {
		valid := false
		for _, r := range rules {
			if inAnyRange(v, r.ranges) {
				valid = true
				break
			}
		}
		if !valid {
			return false
		}
	}
	return true
}

func inAnyRange(v int, ranges [][2]int) bool {
	for _, r := range ranges {
		if v >= r[0] && v <= r[1] {
			return true
		}
	}
	return false
}

func init() {
	solve.Register(Day16{})
}