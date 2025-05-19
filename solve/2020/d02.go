package solve2020

import (
	"fmt"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day2 struct {
}

func (d Day2) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2020, Day: 2}
}

func (d Day2) Part1(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	valid := 0
	for _, line := range lines {
		// Example: "1-3 a: abcde"
		parts := strings.Fields(line)
		if len(parts) != 3 {
			return "", fmt.Errorf("invalid line: %s", line)
		}
		rangeParts := strings.Split(parts[0], "-")
		minCount, _ := strconv.Atoi(rangeParts[0])
		maxCount, _ := strconv.Atoi(rangeParts[1])
		letter := parts[1][0]
		password := parts[2]

		count := 0
		for i := 0; i < len(password); i++ {
			if password[i] == letter {
				count++
			}
		}
		if count >= minCount && count <= maxCount {
			valid++
		}
	}
	return strconv.Itoa(valid), nil
}

func (d Day2) Part2(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	valid := 0
	for _, line := range lines {
		// Example: "1-3 a: abcde"
		parts := strings.Fields(line)
		if len(parts) != 3 {
			return "", fmt.Errorf("invalid line: %s", line)
		}
		rangeParts := strings.Split(parts[0], "-")
		pos1, _ := strconv.Atoi(rangeParts[0])
		pos2, _ := strconv.Atoi(rangeParts[1])
		letter := parts[1][0]
		password := parts[2]

		// Convert to 0-based index
		match1 := pos1-1 < len(password) && password[pos1-1] == letter
		match2 := pos2-1 < len(password) && password[pos2-1] == letter

		if match1 != match2 { // XOR: exactly one is true
			valid++
		}
	}
	return strconv.Itoa(valid), nil
}

func init() {
	solve.Register(Day2{})
}
