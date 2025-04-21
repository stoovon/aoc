package solve2015

import (
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day19 struct {
}

func (d Day19) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2015, Day: 19}
}

var replacementRegex = regexp.MustCompile(`(\S+) => (\S+)`)

func part1(replacements [][2]string, target string) int {
	distinct := make(map[string]struct{})

	for _, repl := range replacements {
		from, to := repl[0], repl[1]
		for i := 0; i <= len(target)-len(from); i++ {
			if target[i:i+len(from)] == from {
				newMolecule := target[:i] + to + target[i+len(from):]
				distinct[newMolecule] = struct{}{}
			}
		}
	}

	return len(distinct)
}

func part2(replacements [][2]string, molecule string) int {
	steps := 0
	target := molecule

	for target != "e" {
		tmp := target
		for _, repl := range replacements {
			from, to := repl[1], repl[0]
			if strings.Contains(target, from) {
				target = strings.Replace(target, from, to, 1)
				steps++
			}
		}

		if tmp == target { // Restart if stuck
			target = molecule
			steps = 0
			rand.Shuffle(len(replacements), func(i, j int) {
				replacements[i], replacements[j] = replacements[j], replacements[i]
			})
		}
	}

	return steps
}

func (d Day19) Part1(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	replacements := parseReplacements(lines[:len(lines)-2])
	target := lines[len(lines)-1]

	return strconv.Itoa(part1(replacements, target)), nil
}

func (d Day19) Part2(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	replacements := parseReplacements(lines[:len(lines)-2])
	target := lines[len(lines)-1]

	return strconv.Itoa(part2(replacements, target)), nil
}

func parseReplacements(lines []string) [][2]string {
	replacements := make([][2]string, 0, len(lines))

	for _, line := range lines {
		matches := replacementRegex.FindStringSubmatch(line)
		if len(matches) == 3 {
			replacements = append(replacements, [2]string{matches[1], matches[2]})
		}
	}

	return replacements
}

func init() {
	solve.Register(Day19{})
}
