package solve2022

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day2 struct {
}

func (d Day2) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2022, Day: 2}
}

var moveScore = map[string]int{"X": 1, "Y": 2, "Z": 3}
var oppMove = map[string]int{"A": 0, "B": 1, "C": 2}
var myMove = map[string]int{"X": 0, "Y": 1, "Z": 2}
var outcomeScorePart1 = []int{3, 6, 0} // draw, win, lose
var outcomeScorePart2 = []int{0, 3, 6} // lose, draw, win

func (d Day2) Part1(input string) (string, error) {
	lines := strings.FieldsFunc(input, func(r rune) bool { return r == '\n' || r == '\r' })
	score := 0
	for _, line := range lines {
		if len(line) < 3 {
			continue
		}
		opp := oppMove[string(line[0])]
		me := myMove[string(line[2])]
		// 0: draw, 1: win, 2: lose
		outcome := (me - opp + 3) % 3
		score += moveScore[string(line[2])] + outcomeScorePart1[outcome]
	}
	return strconv.Itoa(score), nil
}

func (d Day2) Part2(input string) (string, error) {
	lines := strings.FieldsFunc(input, func(r rune) bool { return r == '\n' || r == '\r' })
	score := 0
	for _, line := range lines {
		if len(line) < 3 {
			continue
		}
		opp := oppMove[string(line[0])]
		// X: lose, Y: draw, Z: win
		outcome := myMove[string(line[2])]
		// 0: lose, 1: draw, 2: win
		me := (opp + outcome - 1 + 3) % 3
		score += (me + 1) + outcomeScorePart2[outcome]
	}
	return strconv.Itoa(score), nil
}

func init() {
	solve.Register(Day2{})
}
