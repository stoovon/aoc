package solve2021

import (
	"aoc/solve"
	"errors"
	"strconv"
	"strings"
)

type Day6 struct{}

func simulateLanternfish(input string, days int) (int64, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return 0, errors.New("empty input")
	}
	parts := strings.Split(input, ",")

	timers := make([]int64, 9)
	for _, p := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return 0, err
		}
		timers[n]++
	}
	for range days {
		newFish := timers[0]

		for i := range 8 {
			timers[i] = timers[i+1]
		}
		timers[6] += newFish // Just-spawned fish
		timers[8] = newFish  // Baby fish
	}
	var total int64
	for _, count := range timers {
		total += count
	}
	return total, nil
}

func (d Day6) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2021, Day: 6}
}

func (d Day6) Part1(input string) (string, error) {
	total, err := simulateLanternfish(input, 80)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(total, 10), nil
}

func (d Day6) Part2(input string) (string, error) {
	total, err := simulateLanternfish(input, 256)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(total, 10), nil
}

func init() {
	solve.Register(Day6{})
}
