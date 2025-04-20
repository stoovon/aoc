package solve2024

import (
	"fmt"
	"slices"
	"strings"

	"aoc/solve"
	"aoc/utils/maths"
)

type Day1 struct {
}

type day1DS struct {
	list1   []int
	list2   []int
	counts2 map[int]int
}

func (d Day1) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2024, Day: 1}
}

func (d Day1) getLists(input string) (day1DS, error) {
	var list1, list2 []int
	counts2 := map[int]int{}
	for _, s := range strings.Split(strings.TrimSpace(string(input)), "\n") {
		var n1, n2 int
		_, err := fmt.Sscanf(s, "%d   %d", &n1, &n2)
		if err != nil {
			return day1DS{}, err
		}
		list1, list2 = append(list1, n1), append(list2, n2)
		counts2[n2]++
	}

	slices.Sort(list1)
	slices.Sort(list2)

	return day1DS{
		list1:   list1,
		list2:   list2,
		counts2: counts2,
	}, nil
}

func (d Day1) Part1(input string) (string, error) {
	lists, err := d.getLists(input)
	if err != nil {
		return "", err
	}

	part1 := 0

	for i := range lists.list1 {
		part1 += maths.Abs(lists.list2[i] - lists.list1[i])
	}

	return fmt.Sprintf("%d", part1), nil
}

func (d Day1) Part2(input string) (string, error) {
	lists, err := d.getLists(input)
	if err != nil {
		return "", err
	}

	part2 := 0

	for i := range lists.list2 {
		part2 += lists.list1[i] * lists.counts2[lists.list1[i]]
	}

	return fmt.Sprintf("%d", part2), nil
}

func init() {
	solve.Register(Day1{})
}
