package solve2016

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day20 struct {
}

func (d Day20) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2016, Day: 20}
}

// parseRanges parses the input into a sorted list of integer pairs
func (d Day20) parseRanges(input string) ([][2]int, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	pairs := make([][2]int, len(lines))

	for i, line := range lines {
		matches := digitRe.FindAllString(line, -1)
		if len(matches) != 2 {
			return nil, fmt.Errorf("invalid range: %s", line)
		}
		low, _ := strconv.Atoi(matches[0])
		high, _ := strconv.Atoi(matches[1])
		pairs[i] = [2]int{low, high}
	}

	// Sort pairs by the lower bound, and then by the upper bound
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i][0] == pairs[j][0] {
			return pairs[i][1] < pairs[j][1]
		}
		return pairs[i][0] < pairs[j][0]
	})

	return pairs, nil
}

// unblocked generates all unblocked integers given sorted ranges
func (d Day20) unblocked(pairs [][2]int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		i := 0
		for _, pair := range pairs {
			low, high := pair[0], pair[1]
			for i < low {
				ch <- i
				i++
			}
			if i <= high {
				i = high + 1
			}
		}
	}()
	return ch
}

func (d Day20) Part1(input string) (string, error) {
	pairs, err := d.parseRanges(input)
	if err != nil {
		return "", err
	}

	for i := range d.unblocked(pairs) {
		return strconv.Itoa(i), nil
	}

	return "", fmt.Errorf("no unblocked integers found")
}

func (d Day20) Part2(input string) (string, error) {
	pairs, err := d.parseRanges(input)
	if err != nil {
		return "", err
	}

	count := 0
	for range d.unblocked(pairs) {
		count++
	}

	return strconv.Itoa(count), nil
}

func init() {
	solve.Register(Day20{})
}
