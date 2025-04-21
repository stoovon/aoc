package solve2015

import (
	"fmt"
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/maths"
)

type Day14 struct {
}

func (d Day14) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2015, Day: 14}
}

type reindeer struct {
	Name     string
	Speed    int
	Duration int
	Rest     int
}

func parseReindeer(data string) []reindeer {
	lines := strings.Split(strings.TrimSpace(data), "\n")
	var trackedReindeer []reindeer

	for _, line := range lines {
		parts := strings.Fields(line)
		speed, _ := strconv.Atoi(parts[3])
		duration, _ := strconv.Atoi(parts[6])
		rest, _ := strconv.Atoi(parts[13])

		trackedReindeer = append(trackedReindeer, reindeer{
			Name:     parts[0],
			Speed:    speed,
			Duration: duration,
			Rest:     rest,
		})
	}

	return trackedReindeer
}

func distances(trackedReindeer []reindeer, time int) []int {
	distanceList := make([]int, len(trackedReindeer))

	for i, r := range trackedReindeer {
		cycleTime := r.Duration + r.Rest
		fullCycles := time / cycleTime
		remainingTime := time % cycleTime

		distanceList[i] = r.Speed * (fullCycles*r.Duration + min(remainingTime, r.Duration))
	}

	return distanceList
}

func points(trackedReindeer []reindeer, time int) []int {
	points := make([]int, len(trackedReindeer))

	for t := 1; t <= time; t++ {
		distanceList := distances(trackedReindeer, t)
		maxDistance := maths.Max(distanceList...)

		for i, d := range distanceList {
			if d == maxDistance {
				points[i]++
			}
		}
	}

	return points
}

func (d Day14) Part1(input string) (string, error) {
	reindeer := parseReindeer(input)
	distanceList := distances(reindeer, 2503)
	return fmt.Sprintf("%d", maths.Max(distanceList...)), nil
}

func (d Day14) Part2(input string) (string, error) {
	reindeer := parseReindeer(input)
	pointsList := points(reindeer, 2503)
	return fmt.Sprintf("%d", maths.Max(pointsList...)), nil
}

func init() {
	solve.Register(Day14{})
}
