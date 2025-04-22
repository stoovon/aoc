package solve2023

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day2 struct {
}

func (d Day2) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2023, Day: 2}
}

type cubeDraw struct {
	Red   int
	Green int
	Blue  int
}

func newCubeDraw(value string) cubeDraw {
	counts := make(map[string]int)
	parts := strings.Split(value, ", ")
	for _, part := range parts {
		split := strings.Fields(part)
		if len(split) != 2 {
			continue
		}
		count, _ := strconv.Atoi(split[0])
		color := split[1]
		counts[color] += count
	}

	return cubeDraw{
		Red:   counts["red"],
		Green: counts["green"],
		Blue:  counts["blue"],
	}
}

func (d Day2) parseGame(games []string) []cubeDraw {
	var result []cubeDraw
	for _, game := range games {
		result = append(result, newCubeDraw(game))
	}
	return result
}

func (d Day2) parseGames(input string) map[int][]cubeDraw {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	result := make(map[int][]cubeDraw)

	for _, line := range lines {
		parts := strings.Split(line, ": ")
		if len(parts) != 2 {
			continue
		}

		gameKey, _ := strconv.Atoi(strings.TrimPrefix(parts[0], "Game "))
		gameStrs := strings.Split(parts[1], "; ")
		result[gameKey] = d.parseGame(gameStrs)
	}

	return result
}

func (d Day2) Part1(input string) (string, error) {
	allGames := d.parseGames(input)
	sum := 0

	for key, draws := range allGames {
		isValid := true
		for _, draw := range draws {
			if draw.Blue > 14 || draw.Green > 13 || draw.Red > 12 {
				isValid = false
				break
			}
		}
		if isValid {
			sum += key
		}
	}

	return strconv.Itoa(sum), nil
}

func (d Day2) Part2(input string) (string, error) {
	allGames := d.parseGames(input)
	total := 0

	for _, draws := range allGames {
		biggestDraw := cubeDraw{}
		for _, draw := range draws {
			if draw.Red > biggestDraw.Red {
				biggestDraw.Red = draw.Red
			}
			if draw.Green > biggestDraw.Green {
				biggestDraw.Green = draw.Green
			}
			if draw.Blue > biggestDraw.Blue {
				biggestDraw.Blue = draw.Blue
			}
		}
		total += biggestDraw.Red * biggestDraw.Green * biggestDraw.Blue
	}

	return strconv.Itoa(total), nil
}

func init() {
	solve.Register(Day2{})
}
