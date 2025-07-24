package solve2020

import (
	"aoc/solve"
	"regexp"
	"strconv"
	"strings"
)

type Day7 struct{}

func (d Day7) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2020, Day: 7}
}

type bagRules struct {
	contains map[string][]string // for Part1: contained -> containers
	contents map[string][]struct {
		color string
		count int
	} // for Part2: container -> list of (contained, count)
}

func parseBagRules(input string) bagRules {
	contains := make(map[string][]string)
	contents := make(map[string][]struct {
		color string
		count int
	})
	re := regexp.MustCompile(`^([a-z ]+) bags contain (.+)\.$`)
	reInner := regexp.MustCompile(`(\d+) ([a-z ]+) bag`)
	for _, line := range strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		m := re.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		outer := m[1]
		inners := m[2]
		for _, inner := range reInner.FindAllStringSubmatch(inners, -1) {
			count, _ := strconv.Atoi(inner[1])
			color := inner[2]
			contains[color] = append(contains[color], outer)
			contents[outer] = append(contents[outer], struct {
				color string
				count int
			}{color, count})
		}
	}
	return bagRules{contains, contents}
}

func (d Day7) Part1(input string) (string, error) {
	rules := parseBagRules(input)
	seen := make(map[string]bool)
	queue := []string{"shiny gold"}
	for len(queue) > 0 {
		color := queue[0]
		queue = queue[1:]
		for _, outer := range rules.contains[color] {
			if !seen[outer] {
				seen[outer] = true
				queue = append(queue, outer)
			}
		}
	}
	return strconv.Itoa(len(seen)), nil
}

func (d Day7) Part2(input string) (string, error) {
	rules := parseBagRules(input)
	var countBags func(string) int
	countBags = func(color string) int {
		total := 0
		for _, inner := range rules.contents[color] {
			total += inner.count * (1 + countBags(inner.color))
		}
		return total
	}
	return strconv.Itoa(countBags("shiny gold")), nil
}

func init() {
	solve.Register(Day7{})
}
