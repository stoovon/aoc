package solve2023

import (
	"errors"
	"math"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day24 struct {
}

func (d Day24) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2023, Day: 24}
}

const epsilon = 0.0001

type Vec3 struct {
	x, y, z float64
}

type Hailstone struct {
	pos Vec3
	vel Vec3
}

func parseVec3(s string) (Vec3, error) {
	parts := strings.Split(strings.TrimSpace(s), ",")
	if len(parts) != 3 {
		return Vec3{}, errors.New("invalid Vec3")
	}
	x, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return Vec3{}, err
	}
	y, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return Vec3{}, err
	}
	z, err := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
	if err != nil {
		return Vec3{}, err
	}
	return Vec3{x, y, z}, nil
}

func parseHailstone(line string) (Hailstone, error) {
	parts := strings.Split(line, "@")
	if len(parts) != 2 {
		return Hailstone{}, errors.New("invalid hailstone")
	}
	pos, err := parseVec3(parts[0])
	if err != nil {
		return Hailstone{}, err
	}
	vel, err := parseVec3(parts[1])
	if err != nil {
		return Hailstone{}, err
	}
	return Hailstone{pos, vel}, nil
}

func (d Day24) parseInput(input string) ([]Hailstone, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	hailstones := make([]Hailstone, 0, len(lines))
	for _, line := range lines {
		hs, err := parseHailstone(line)
		if err != nil {
			return nil, err
		}
		hailstones = append(hailstones, hs)
	}
	return hailstones, nil
}

func intersect2D(a, b Hailstone) (t1, t2, x, y float64) {
	p1x, p1y := a.pos.x, a.pos.y
	v1x, v1y := a.vel.x, a.vel.y
	p2x, p2y := b.pos.x, b.pos.y
	v2x, v2y := b.vel.x, b.vel.y

	denom := v2x*v1y - v2y*v1x
	if math.Abs(denom) < epsilon {
		return math.NaN(), math.NaN(), math.NaN(), math.NaN()
	}
	t2 = ((p2y-p1y)*v1x - (p2x-p1x)*v1y) / denom
	t1 = (p2x - p1x + t2*v2x) / v1x
	x = p1x + t1*v1x
	y = p1y + t1*v1y
	return
}

func (d Day24) Part1(input string) (string, error) {
	hailstones, err := d.parseInput(input)
	if err != nil {
		return "", err
	}
	const minSeen, maxSeen = 200000000000000, 400000000000000
	count := 0
	for i := 0; i < len(hailstones); i++ {
		for j := i + 1; j < len(hailstones); j++ {
			t1, t2, x, y := intersect2D(hailstones[i], hailstones[j])
			if t1 > 0 && t2 > 0 && x >= minSeen && x <= maxSeen && y >= minSeen && y <= maxSeen {
				count++
			}
		}
	}
	return strconv.Itoa(count), nil
}

func isInt(f float64) bool {
	return math.Abs(math.Round(f)-f) < epsilon
}

func solvePart2(a, b Hailstone, vx, vy, vz float64) (t1, t2 float64, px, py, pz float64, ok bool) {
	pax, pay, paz := a.pos.x, a.pos.y, a.pos.z
	vax, vay, vaz := a.vel.x, a.vel.y, a.vel.z
	pbx, pby, pbz := b.pos.x, b.pos.y, b.pos.z
	vbx, vby, vbz := b.vel.x, b.vel.y, b.vel.z

	denom := vy - vby - ((vay-vy)*(vx-vbx))/(vax-vx)
	if math.Abs(denom) < epsilon {
		return
	}
	t2 = (pby - pay - ((vay-vy)*(pbx-pax))/(vax-vx)) / denom
	t1 = (pbx - pax - t2*(vx-vbx)) / (vax - vx)
	px = pax - t1*(vx-vax)
	py = pay - t1*(vy-vay)
	pz = paz - t1*(vz-vaz)
	if math.Abs(pz+t2*(vz-vbz)-pbz) > epsilon {
		return
	}
	ok = true
	return
}

func (d Day24) Part2(input string) (string, error) {
	hailstones, err := d.parseInput(input)
	if err != nil {
		return "", err
	}
	a, b := hailstones[0], hailstones[1]
	for vx := -500; vx < 500; vx++ {
		for vy := -500; vy < 500; vy++ {
			for vz := -500; vz < 500; vz++ {
				t1, t2, px, py, pz, ok := solvePart2(a, b, float64(vx), float64(vy), float64(vz))
				if !ok || !isInt(t1) || !isInt(t2) || !isInt(px) || !isInt(py) || !isInt(pz) || t1 < 0 || t2 < 0 {
					continue
				}
				valid := true
				for i := 2; i < len(hailstones); i++ {
					c := hailstones[i]
					vcx, vcy, vcz := c.vel.x, c.vel.y, c.vel.z
					pcx, pcy, pcz := c.pos.x, c.pos.y, c.pos.z
					t3 := (pcx - px) / (float64(vx) - vcx)
					if math.Abs(py+t3*float64(vy)-(pcy+t3*vcy)) > epsilon ||
						math.Abs(pz+t3*float64(vz)-(pcz+t3*vcz)) > epsilon {
						valid = false
						break
					}
				}
				if valid {
					return strconv.FormatInt(int64(px+py+pz), 10), nil
				}
			}
		}
	}
	return "", errors.New("no solution found")
}

func init() {
	solve.Register(Day24{})
}
