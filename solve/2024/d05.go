package solve2024

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day5 struct {
}

func (d Day5) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2024, Day: 5}
}

func (d Day5) createComparator(pageOrderingRules string) func(a, b string) int {
	return func(a, b string) int {
		for _, rule := range strings.Split(pageOrderingRules, "\n") {
			if ruleComponent := strings.Split(rule, "|"); ruleComponent[0] == a && ruleComponent[1] == b {
				return -1
			}
		}
		return 0
	}
}

func (d Day5) solve(input string, shouldSort bool) (r int, err error) {
	sections := strings.Split(strings.TrimSpace(input), "\n\n")
	if len(sections) < 2 {
		return 0, fmt.Errorf("invalid input format")
	}
	pageOrderingRules := sections[0]
	pagesToProduce := sections[1]

	comparator := Day5{}.createComparator(pageOrderingRules)
	total := 0

	for _, ptp := range strings.Split(pagesToProduce, "\n") {
		if page := strings.Split(ptp, ","); slices.IsSortedFunc(page, comparator) == shouldSort {
			slices.SortFunc(page, comparator)
			midpoint, err := strconv.Atoi(page[len(page)/2])
			if err != nil {
				return 0, err
			}
			total += midpoint
		}
	}
	return total, nil
}

func (d Day5) formatResult(result int, err error) string {
	if err != nil {
		return err.Error()
	}
	return strconv.Itoa(result)
}

func (d Day5) Part1(input string) (string, error) {
	return d.formatResult(d.solve(input, true)), nil
}

func (d Day5) Part2(input string) (string, error) {
	return d.formatResult(d.solve(input, false)), nil
}

func init() {
	solve.Register(Day5{})
}
