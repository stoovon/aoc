package solve2016

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day10 struct {
}

func (d Day10) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2016, Day: 10}
}

var (
	valueRegex = regexp.MustCompile(`value (\d+) goes to bot (\d+)`)
	botRegex   = regexp.MustCompile(`bot (\d+) gives low to (bot|output) (\d+) and high to (bot|output) (\d+)`)
)

func (d Day10) Solve(input string) (found, botProduct string) {
	input = strings.TrimSpace(input)

	bot := make(map[int][]int)
	output := make(map[int][]int)
	pipeline := make(map[int][2]struct {
		targetType string
		targetID   int
	})

	lines := strings.Split(input, "\n")

	for _, line := range lines {
		if valueRegex.MatchString(line) {
			matches := valueRegex.FindStringSubmatch(line)
			value, _ := strconv.Atoi(matches[1])
			botID, _ := strconv.Atoi(matches[2])
			bot[botID] = append(bot[botID], value)
		} else if botRegex.MatchString(line) {
			matches := botRegex.FindStringSubmatch(line)
			botID, _ := strconv.Atoi(matches[1])
			lowTargetType := matches[2]
			lowTargetID, _ := strconv.Atoi(matches[3])
			highTargetType := matches[4]
			highTargetID, _ := strconv.Atoi(matches[5])
			pipeline[botID] = [2]struct {
				targetType string
				targetID   int
			}{
				{lowTargetType, lowTargetID},
				{highTargetType, highTargetID},
			}
		}
	}

	var foundBot int

	for len(bot) > 0 {
		for botID, values := range bot {
			if len(values) == 2 {
				sort.Ints(values)
				v1, v2 := values[0], values[1]
				if v1 == 17 && v2 == 61 {
					foundBot = botID
				}

				lowTarget := pipeline[botID][0]
				highTarget := pipeline[botID][1]

				if lowTarget.targetType == "bot" {
					bot[lowTarget.targetID] = append(bot[lowTarget.targetID], v1)
				} else {
					output[lowTarget.targetID] = append(output[lowTarget.targetID], v1)
				}

				if highTarget.targetType == "bot" {
					bot[highTarget.targetID] = append(bot[highTarget.targetID], v2)
				} else {
					output[highTarget.targetID] = append(output[highTarget.targetID], v2)
				}

				delete(bot, botID)
			}
		}
	}

	product := output[0][0] * output[1][0] * output[2][0]
	return strconv.Itoa(foundBot), strconv.Itoa(product)
}

func (d Day10) Part1(input string) (string, error) {
	foundBot, _ := d.Solve(input)
	return foundBot, nil
}

func (d Day10) Part2(input string) (string, error) {
	_, botProduct := d.Solve(input)
	return botProduct, nil
}

func init() {
	solve.Register(Day10{})
}
