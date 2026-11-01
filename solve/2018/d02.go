package solve2018

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day2 struct{}

func (Day2) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2018, Day: 2}
}

func (Day2) Part1(input string) (string, error) {
	ids, err := inventoryIDs(input)
	if err != nil {
		return "", err
	}
	var twos, threes uint64
	for _, boxID := range ids {
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
				break
			}
		}
	}

	if threes != 0 && twos > math.MaxUint64/threes {
		return "", fmt.Errorf("checksum exceeds uint64")
	}
	return strconv.FormatUint(twos*threes, 10), nil
}

func getCharCount(box string) map[rune]int {
	chars := make(map[rune]int)
	for _, c := range box {
		chars[c]++
	}
	return chars
}

func (Day2) Part2(input string) (string, error) {
	lines, err := inventoryIDs(input)
	if err != nil {
		return "", err
	}
	for i := 0; i < len(lines); i++ {
		for j := i + 1; j < len(lines); j++ {
			if sameChars, matches := commonLetters(lines[i], lines[j]); matches {
				return sameChars, nil
			}
		}
	}
	return "", fmt.Errorf("no IDs differ at exactly one position")
}
func inventoryIDs(input string) ([]string, error) {
	ids := strings.Split(strings.TrimSuffix(input, "\n"), "\n")
	for line, id := range ids {
		id = strings.TrimSuffix(id, "\r")
		if id == "" {
			return nil, fmt.Errorf("line %d: empty box ID", line+1)
		}
		for i := 0; i < len(id); i++ {
			if id[i] < 'a' || id[i] > 'z' {
				return nil, fmt.Errorf("line %d: box IDs must contain lowercase ASCII letters", line+1)
			}
		}
		ids[line] = id
	}
	return ids, nil
}

func commonLetters(first, second string) (string, bool) {
	if len(first) != len(second) {
		return "", false
	}
	mismatch := -1
	for i := 0; i < len(first); i++ {
		if first[i] == second[i] {
			continue
		}
		if mismatch >= 0 {
			return "", false
		}
		mismatch = i
	}
	if mismatch < 0 {
		return "", false
	}
	return first[:mismatch] + first[mismatch+1:], true
}

func init() {
	solve.Register(Day2{})
}
