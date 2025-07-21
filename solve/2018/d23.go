package solve2018

import (
	"aoc/solve"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Nanobot struct {
	x, y, z, r int
}

type Day23 struct{}

func (d Day23) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2018, Day: 23}
}

func parseNanobots(input string) []Nanobot {
	lines := strings.Split(input, "\n")
	nanobots := []Nanobot{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		var x, y, z, r int
		fmt.Sscanf(line, "pos=<%d,%d,%d>, r=%d", &x, &y, &z, &r)
		nanobots = append(nanobots, Nanobot{x, y, z, r})
	}
	return nanobots
}

func manhattanDistance(a, b Nanobot) int {
	return int(math.Abs(float64(a.x-b.x)) + math.Abs(float64(a.y-b.y)) + math.Abs(float64(a.z-b.z)))
}

func (d Day23) Part1(input string) (string, error) {
	nanobots := parseNanobots(input)

	// Find the strongest nanobot
	strongest := nanobots[0]
	for _, bot := range nanobots {
		if bot.r > strongest.r {
			strongest = bot
		}
	}

	// Count nanobots in range of the strongest nanobot
	count := 0
	for _, bot := range nanobots {
		if manhattanDistance(strongest, bot) <= strongest.r {
			count++
		}
	}

	return strconv.Itoa(count), nil
}

func (d Day23) Part2(input string) (string, error) {
	nanobots := parseNanobots(input)

	// Define the search space
	minX, minY, minZ := math.MaxInt, math.MaxInt, math.MaxInt
	maxX, maxY, maxZ := math.MinInt, math.MinInt, math.MinInt
	for _, bot := range nanobots {
		minX = int(math.Min(float64(minX), float64(bot.x-bot.r)))
		minY = int(math.Min(float64(minY), float64(bot.y-bot.r)))
		minZ = int(math.Min(float64(minZ), float64(bot.z-bot.r)))
		maxX = int(math.Max(float64(maxX), float64(bot.x+bot.r)))
		maxY = int(math.Max(float64(maxY), float64(bot.y+bot.r)))
		maxZ = int(math.Max(float64(maxZ), float64(bot.z+bot.r)))
	}

	// Start with a large step size and refine
	step := 1
	for step < maxX-minX || step < maxY-minY || step < maxZ-minZ {
		step *= 2
	}

	bestCount := 0
	bestDistance := math.MaxInt
	bestCoord := Nanobot{}

	for step > 0 {
		for x := minX; x <= maxX; x += step {
			for y := minY; y <= maxY; y += step {
				for z := minZ; z <= maxZ; z += step {
					count := 0
					for _, bot := range nanobots {
						distance := int(math.Abs(float64(x-bot.x)) + math.Abs(float64(y-bot.y)) + math.Abs(float64(z-bot.z)))
						if distance <= bot.r {
							count++
						}
					}

					distanceToOrigin := int(math.Abs(float64(x)) + math.Abs(float64(y)) + math.Abs(float64(z)))
					if count > bestCount || (count == bestCount && distanceToOrigin < bestDistance) {
						bestCount = count
						bestDistance = distanceToOrigin
						bestCoord = Nanobot{x, y, z, 0}
					}
				}
			}
		}
		minX, minY, minZ = bestCoord.x-step, bestCoord.y-step, bestCoord.z-step
		maxX, maxY, maxZ = bestCoord.x+step, bestCoord.y+step, bestCoord.z+step
		step /= 2
	}

	return strconv.Itoa(bestDistance), nil
}

func init() {
	solve.Register(Day23{})
}
