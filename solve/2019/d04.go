package solve2019

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day4 struct {
}

func (d Day4) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2019, Day: 4}
}

func isValidPassword(n int) bool {
	s := strconv.Itoa(n)
	hasDouble := false
	for i := 1; i < len(s); i++ {
		if s[i] == s[i-1] {
			hasDouble = true
		}
		if s[i] < s[i-1] {
			return false
		}
	}
	return hasDouble
}

func isValidPasswordExactDouble(n int) bool {
	s := strconv.Itoa(n)
	counts := make(map[byte]int)
	for i := 0; i < len(s); {
		j := i + 1
		for j < len(s) && s[j] == s[i] {
			j++
		}
		counts[s[i]] = j - i
		if j < len(s) && s[j] < s[j-1] {
			return false
		}
		i = j
	}
	hasExactDouble := false
	for _, v := range counts {
		if v == 2 {
			hasExactDouble = true
			break
		}
	}
	// Check for non-decreasing
	for i := 1; i < len(s); i++ {
		if s[i] < s[i-1] {
			return false
		}
	}
	return hasExactDouble
}

func (d Day4) Part1(input string) (string, error) {
	parts := strings.Split(strings.TrimSpace(input), "-")
	if len(parts) != 2 {
		return "", errors.New("invalid input")
	}
	start, err1 := strconv.Atoi(parts[0])
	end, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return "", errors.New("invalid range")
	}
	count := 0
	for n := start; n <= end; n++ {
		if isValidPassword(n) {
			count++
		}
	}
	return fmt.Sprintf("%d", count), nil
}

func (d Day4) Part2(input string) (string, error) {
	parts := strings.Split(strings.TrimSpace(input), "-")
	if len(parts) != 2 {
		return "", errors.New("invalid input")
	}
	start, err1 := strconv.Atoi(parts[0])
	end, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return "", errors.New("invalid range")
	}
	count := 0
	for n := start; n <= end; n++ {
		if isValidPasswordExactDouble(n) {
			count++
		}
	}
	return fmt.Sprintf("%d", count), nil
}

func init() {
	solve.Register(Day4{})
}
