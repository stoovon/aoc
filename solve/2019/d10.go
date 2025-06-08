package solve2019

import (
	"errors"
	"image"
	"math"
	"sort"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day10 struct {
}

func (d Day10) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2019, Day: 10}
}

func (d Day10) parseAsteroids(input string) []image.Point {
	var asteroids []image.Point
	lines := strings.Split(strings.TrimSpace(input), "\n")
	for y, line := range lines {
		for x, c := range line {
			if c == '#' {
				asteroids = append(asteroids, image.Point{X: x, Y: y})
			}
		}
	}
	return asteroids
}

func (d Day10) bestStation(asteroids []image.Point) (station image.Point, maxVisible int) {
	for _, a := range asteroids {
		angles := make(map[float64]struct{})
		for _, b := range asteroids {
			if a == b {
				continue
			}
			dx := float64(b.X - a.X)
			dy := float64(b.Y - a.Y)
			angle := math.Atan2(dy, dx)
			angles[angle] = struct{}{}
		}
		if len(angles) > maxVisible {
			maxVisible = len(angles)
			station = a
		}
	}
	return
}

type asteroid struct {
	x, y  int
	angle float64
	dist  float64
}

func (d Day10) Part1(input string) (string, error) {
	asteroids := d.parseAsteroids(input)
	_, maxVisible := d.bestStation(asteroids)
	return strconv.Itoa(maxVisible), nil
}

func (d Day10) Part2(input string) (string, error) {
	asteroids := d.parseAsteroids(input)
	station, _ := d.bestStation(asteroids)

	byAngle := make(map[float64][]asteroid)
	for _, a := range asteroids {
		if a == station {
			continue
		}
		dx := float64(a.X - station.X)
		dy := float64(a.Y - station.Y)
		// Angle: 0 is up, increases clockwise
		angle := math.Atan2(dx, -dy)
		if angle < 0 {
			angle += 2 * math.Pi
		}
		dist := dx*dx + dy*dy
		byAngle[angle] = append(byAngle[angle], asteroid{a.X, a.Y, angle, dist})
	}

	// Sort asteroids in each angle by distance (closest first)
	for angle := range byAngle {
		sort.Slice(byAngle[angle], func(i, j int) bool {
			return byAngle[angle][i].dist < byAngle[angle][j].dist
		})
	}

	// Sort angles in increasing order (clockwise from up)
	var angles []float64
	for angle := range byAngle {
		angles = append(angles, angle)
	}
	sort.Float64s(angles)

	count := 0
	for {
		for _, angle := range angles {
			if len(byAngle[angle]) == 0 {
				continue
			}
			count++
			a := byAngle[angle][0]
			byAngle[angle] = byAngle[angle][1:]
			if count == 200 {
				return strconv.Itoa(a.x*100 + a.y), nil
			}
		}

		allEmpty := true
		for _, v := range byAngle {
			if len(v) > 0 {
				allEmpty = false
				break
			}
		}
		if allEmpty {
			break
		}
	}
	return "", errors.New("fewer than 200 asteroids vaporized")
}

func init() {
	solve.Register(Day10{})
}
