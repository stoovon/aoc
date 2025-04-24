package solve2016

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day4 struct {
}

func (d Day4) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2016, Day: 4}
}

// parseLine converts a line into name, sector, and checksum
func (d Day4) parseLine(line string) (string, int, string) {
	re := regexp.MustCompile(`^(.+)-(\d+)\[([a-z]+)]$`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 4 {
		return "", 0, ""
	}
	sector, _ := strconv.Atoi(matches[2])
	return matches[1], sector, matches[3]
}

// isValid verifies the checksum of the room name
func (d Day4) isValid(name, checksum string) bool {
	counts := make(map[rune]int)
	for _, char := range name {
		if char != '-' {
			counts[char]++
		}
	}

	type pair struct {
		char  rune
		count int
	}
	var pairs []pair
	for char, count := range counts {
		pairs = append(pairs, pair{char, count})
	}

	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count == pairs[j].count {
			return pairs[i].char < pairs[j].char
		}
		return pairs[i].count > pairs[j].count
	})

	var calculatedChecksum strings.Builder
	for i := 0; i < 5 && i < len(pairs); i++ {
		calculatedChecksum.WriteRune(pairs[i].char)
	}

	return calculatedChecksum.String() == checksum
}

// sectorNumber returns the sector number of a valid line, else 0
func (d Day4) sectorNumber(line string) int {
	name, sector, checksum := d.parseLine(line)
	if d.isValid(name, checksum) {
		return sector
	}
	return 0
}

// rotN implements a Caesar cipher
func (d Day4) rotN(text string, shift int) string {
	shift = shift % 26
	var result strings.Builder
	for _, char := range text {
		if char == '-' {
			result.WriteRune(' ')
		} else {
			newChar := ((char-'a')+rune(shift))%26 + 'a'
			result.WriteRune(newChar)
		}
	}
	return result.String()
}

// decrypt decrypts the room name using the sector ID
func (d Day4) decrypt(line string) string {
	name, sector, _ := d.parseLine(line)
	return d.rotN(name, sector) + " " + strconv.Itoa(sector)
}

// grep Searches for a pattern in decrypted room names and extracts the sector ID
func (d Day4) grep(pattern string, lines []string) int {
	re := regexp.MustCompile(pattern)
	for _, line := range lines {
		if re.MatchString(line) {
			matches := regexp.MustCompile(`\d+`).FindStringSubmatch(line)
			if len(matches) > 0 {
				sector, _ := strconv.Atoi(matches[0])
				return sector
			}
		}
	}
	return 0
}

func (d Day4) Part1(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	sum := 0
	for _, line := range lines {
		sum += d.sectorNumber(line)
	}
	return strconv.Itoa(sum), nil
}

func (d Day4) Part2(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var decrypted []string
	for _, line := range lines {
		decrypted = append(decrypted, d.decrypt(line))
	}
	sector := d.grep("north", decrypted)
	return strconv.Itoa(sector), nil
}

func init() {
	solve.Register(Day4{})
}
