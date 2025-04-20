package solve2018

import (
	"aoc/solve"
	"fmt"
	"strings"
)

type Day2 struct {
}

func (Day2) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2018, Day: 2}
}

func (Day2) Part1(input string) (string, error) {
	var twos, threes int
	for _, boxID := range strings.Split(input, "\n") {
		charCounts := getCharCount(boxID)
		for _, v := range charCounts {
			if v == 2 {
				twos++
				break
			}
		}
		for _, v := range charCounts {
			if v == 3 {
				threes++
			}
		}
	}

	return fmt.Sprintf("%d", twos*threes), nil
}

func getCharCount(box string) map[rune]int {
	chars := make(map[rune]int)
	for _, c := range box {
		chars[c]++
	}
	return chars
}

func (Day2) Part2(input string) (string, error) {
	lines := strings.Split(input, "\n")
	for i := 0; i < len(lines); i++ {
		for j := i + 1; j < len(lines)-1; j++ {
			if sameChars := getSameCharacters(lines[i], lines[j]); sameChars != "" {
				return sameChars, nil
			}
		}
	}
	return "", nil
}
func getSameCharacters(str1, str2 string) string {
	var mismatchSeen bool
	var sameChars string
	for i := 0; i < len(str1); i++ {
		if str1[i] == str2[i] {
			sameChars += string(str1[i])
		} else if mismatchSeen {
			// if a mismatch has already been seen, then it's 2 characters off
			// return an empty string
			return ""
		} else {
			mismatchSeen = true
		}
	}
	return sameChars
}

func init() {
	solve.Register(Day2{})
}
