package solve2019

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/maths"
)

type Day12 struct {
}

type moon struct {
	pos [3]int
	vel [3]int
}

var moonRe = regexp.MustCompile(`<x=(-?\d+), y=(-?\d+), z=(-?\d+)>`)

func parseMoons(input string) []moon {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	moons := make([]moon, 0, len(lines))
	for _, line := range lines {
		m := moonRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		x, _ := strconv.Atoi(m[1])
		y, _ := strconv.Atoi(m[2])
		z, _ := strconv.Atoi(m[3])
		moons = append(moons, moon{pos: [3]int{x, y, z}})
	}
	return moons
}

func applyGravity(moons []moon) {
	for i := 0; i < len(moons); i++ {
		for j := i + 1; j < len(moons); j++ {
			for axis := 0; axis < 3; axis++ {
				if moons[i].pos[axis] < moons[j].pos[axis] {
					moons[i].vel[axis]++
					moons[j].vel[axis]--
				} else if moons[i].pos[axis] > moons[j].pos[axis] {
					moons[i].vel[axis]--
					moons[j].vel[axis]++
				}
			}
		}
	}
}

func applyVelocity(moons []moon) {
	for i := range moons {
		for axis := 0; axis < 3; axis++ {
			moons[i].pos[axis] += moons[i].vel[axis]
		}
	}
}

func totalEnergy(moons []moon) int {
	sum := 0
	for _, m := range moons {
		pot := maths.Abs(m.pos[0]) + maths.Abs(m.pos[1]) + maths.Abs(m.pos[2])
		kin := maths.Abs(m.vel[0]) + maths.Abs(m.vel[1]) + maths.Abs(m.vel[2])
		sum += pot * kin
	}
	return sum
}

func (d Day12) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2019, Day: 12}
}

func (d Day12) Part1(input string) (string, error) {
	moons := parseMoons(input)
	for step := 0; step < 1000; step++ {
		applyGravity(moons)
		applyVelocity(moons)
	}
	return fmt.Sprint(totalEnergy(moons)), nil
}

func axisCycleLength(moons []moon, axis int) int64 {
	type state struct {
		pos [4]int
		vel [4]int
	}
	var initial state
	for i := 0; i < 4; i++ {
		initial.pos[i] = moons[i].pos[axis]
		initial.vel[i] = moons[i].vel[axis]
	}
	cur := initial
	steps := int64(0)
	for {
		// Apply gravity
		for i := 0; i < 4; i++ {
			for j := i + 1; j < 4; j++ {
				if cur.pos[i] < cur.pos[j] {
					cur.vel[i]++
					cur.vel[j]--
				} else if cur.pos[i] > cur.pos[j] {
					cur.vel[i]--
					cur.vel[j]++
				}
			}
		}
		// Apply velocity
		for i := 0; i < 4; i++ {
			cur.pos[i] += cur.vel[i]
		}
		steps++
		if cur == initial {
			return steps
		}
	}
}

func (d Day12) Part2(input string) (string, error) {
	moons := parseMoons(input)
	xCycle := axisCycleLength(moons, 0)
	yCycle := axisCycleLength(moons, 1)
	zCycle := axisCycleLength(moons, 2)
	return fmt.Sprint(maths.LCM(int(xCycle), maths.LCM(int(yCycle), int(zCycle)))), nil
}

func init() {
	solve.Register(Day12{})
}
