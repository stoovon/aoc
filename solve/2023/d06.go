package solve2023

import (
	"errors"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day6 struct {
}

func (d Day6) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2023, Day: 6}
}

func race(time, distance int64) int64 {
	start := int64(0)
	for hold := int64(1); hold <= time; hold++ {
		if (time-hold)*hold > distance {
			start = hold - 1
			break
		}
	}

	end := int64(0)
	for hold := time; hold >= 1; hold-- {
		if (time-hold)*hold > distance {
			end = hold
			break
		}
	}

	return end - start
}

func (d Day6) Part1(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) < 2 {
		return "", errors.New("invalid input format")
	}

	// Parse times
	timesStr := strings.Fields(strings.Split(lines[0], ":")[1])
	times := make([]int64, len(timesStr))
	for i, t := range timesStr {
		val, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return "", err
		}
		times[i] = val
	}

	// Parse distances
	distancesStr := strings.Fields(strings.Split(lines[1], ":")[1])
	distances := make([]int64, len(distancesStr))
	for i, d := range distancesStr {
		val, err := strconv.ParseInt(d, 10, 64)
		if err != nil {
			return "", err
		}
		distances[i] = val
	}

	// Calculate result
	result := int64(1)
	for i := range times {
		result *= race(times[i], distances[i])
	}

	return strconv.FormatInt(result, 10), nil
}

func (d Day6) Part2(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) < 2 {
		return "", errors.New("invalid input format")
	}

	// Parse time
	timeStr := strings.Split(lines[0], ":")[1]
	time, err := strconv.ParseInt(strings.TrimSpace(filterDigits(timeStr)), 10, 64)
	if err != nil {
		return "", err
	}

	// Parse distance
	distanceStr := strings.Split(lines[1], ":")[1]
	distance, err := strconv.ParseInt(strings.TrimSpace(filterDigits(distanceStr)), 10, 64)
	if err != nil {
		return "", err
	}

	// Calculate result
	result := race(time, distance)
	return strconv.FormatInt(result, 10), nil
}

// Helper function to filter digits from a string
func filterDigits(input string) string {
	var result strings.Builder
	for _, r := range input {
		if r >= '0' && r <= '9' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func init() {
	solve.Register(Day6{})
}
