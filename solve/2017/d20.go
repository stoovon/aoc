package solve2017

import (
	"aoc/solve"
	"aoc/utils/maths"

	"errors"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type Day20 struct {
}

func (d Day20) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2017, Day: 20}
}

type vec3 struct{ x, y, z int }

func manhattan(v vec3) int {
	return maths.Abs(v.x) + maths.Abs(v.y) + maths.Abs(v.z)
}

type particle struct {
	idx     int
	p, v, a vec3
}

var particleRe = regexp.MustCompile(`p=<\s*(-?\d+),\s*(-?\d+),\s*(-?\d+)>, v=<\s*(-?\d+),\s*(-?\d+),\s*(-?\d+)>, a=<\s*(-?\d+),\s*(-?\d+),\s*(-?\d+)>`)

func parseParticles(input string) ([]*particle, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	particles := make([]*particle, 0, len(lines))
	for i, line := range lines {
		m := particleRe.FindStringSubmatch(line)
		if m == nil {
			return nil, errors.New("invalid input line: " + line)
		}
		nums := make([]int, 9)
		for j := 1; j <= 9; j++ {
			n, _ := strconv.Atoi(m[j])
			nums[j-1] = n
		}
		particles = append(particles, &particle{
			idx: i,
			p:   vec3{nums[0], nums[1], nums[2]},
			v:   vec3{nums[3], nums[4], nums[5]},
			a:   vec3{nums[6], nums[7], nums[8]},
		})
	}
	return particles, nil
}

func (d Day20) Part1(input string) (string, error) {
	particles, err := parseParticles(input)
	if err != nil {
		return "", err
	}
	minIdx := -1
	minA, minV, minP := math.MaxInt, math.MaxInt, math.MaxInt
	for _, part := range particles {
		ma := manhattan(part.a)
		mv := manhattan(part.v)
		mp := manhattan(part.p)
		if ma < minA ||
			(ma == minA && mv < minV) ||
			(ma == minA && mv == minV && mp < minP) {
			minA, minV, minP = ma, mv, mp
			minIdx = part.idx
		}
	}
	return strconv.Itoa(minIdx), nil
}

func (d Day20) Part2(input string) (string, error) {
	particles, err := parseParticles(input)
	if err != nil {
		return "", err
	}
	alive := make(map[int]*particle, len(particles))
	for _, p := range particles {
		alive[p.idx] = p
	}
	const maxTicks = 1000
	for tick := 0; tick < maxTicks; tick++ {
		posMap := map[[3]int][]int{}
		for _, part := range alive {
			part.v.x += part.a.x
			part.v.y += part.a.y
			part.v.z += part.a.z
			part.p.x += part.v.x
			part.p.y += part.v.y
			part.p.z += part.v.z
			pos := [3]int{part.p.x, part.p.y, part.p.z}
			posMap[pos] = append(posMap[pos], part.idx)
		}
		collided := map[int]bool{}
		for _, idxs := range posMap {
			if len(idxs) > 1 {
				for _, idx := range idxs {
					collided[idx] = true
				}
			}
		}
		for idx := range collided {
			delete(alive, idx)
		}
	}
	return strconv.Itoa(len(alive)), nil
}

func init() {
	solve.Register(Day20{})
}
