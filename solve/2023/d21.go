package solve2023

import (
	"aoc/solve"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Day21 struct{}

func (Day21) Coords() solve.SolutionCoords { return solve.SolutionCoords{Year: 2023, Day: 21} }

const gardenCellLimit = 1_048_576
const gardenPatchLimit = 4_194_304

var errGardenLimit = errors.New("garden resource limit exceeded")
var errGardenModel = errors.New("repeating garden has no certified distance model within the fixed budget")
var errGardenOverflow = errors.New("garden count exceeds int64")

type gardenPoint struct{ x, y int }

var gardenDirections = [4]gardenPoint{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

type gardenMap struct {
	width, height int
	cells         []byte
	start         gardenPoint
}

func parseGarden(input string) (gardenMap, error) {
	var g gardenMap
	found := false
	for row, line := range strings.Split(strings.TrimSuffix(input, "\n"), "\n") {
		line = strings.TrimSuffix(line, "\r")
		invalid := func(message string) error { return fmt.Errorf("line %d: %s", row+1, message) }
		if line == "" {
			return g, invalid("expected a nonempty rectangular garden")
		}
		if g.width == 0 {
			g.width = len(line)
		}
		if len(line) != g.width {
			return g, invalid("garden rows must have equal width")
		}
		if len(line) > gardenCellLimit-len(g.cells) {
			return g, errGardenLimit
		}
		for col, c := range []byte(line) {
			switch c {
			case '.', '#':
			case 'S':
				if found {
					return g, invalid("expected exactly one start")
				}
				found = true
				g.start = gardenPoint{col, row}
			default:
				return g, invalid("expected only '.', '#' or 'S'")
			}
		}
		g.cells = append(g.cells, line...)
		g.height++
	}
	if !found {
		return g, errors.New("expected exactly one start")
	}
	return g, nil
}
func gardenMod(value, width int) int {
	r := value % width
	if r < 0 {
		r += width
	}
	return r
}
func (g *gardenMap) open(p gardenPoint, repeating bool) bool {
	x, y := p.x, p.y
	if repeating {
		x = gardenMod(x, g.width)
		y = gardenMod(y, g.height)
	} else if x < 0 || y < 0 || x >= g.width || y >= g.height {
		return false
	}
	return g.cells[y*g.width+x] != '#'
}
func (g *gardenMap) hasEdge(repeating bool) bool {
	for _, d := range gardenDirections {
		if g.open(gardenPoint{g.start.x + d.x, g.start.y + d.y}, repeating) {
			return true
		}
	}
	return false
}
func (g *gardenMap) finiteCount(steps int64) (int64, error) {
	if steps < 0 {
		return 0, errors.New("negative step count")
	}
	if !g.hasEdge(false) {
		if steps == 0 {
			return 1, nil
		}
		return 0, nil
	}
	start := g.start.y*g.width + g.start.x
	visited := make([]bool, len(g.cells))
	visited[start] = true
	queue := make([]int, 1, len(g.cells))
	queue[0] = start
	parity := [2]int64{1, 0}
	head := 0
	// Shortest distance plus parity suffices: any spare pair of steps can traverse
	// an undirected edge and return. Isolated starts were handled above.
	for distance := int64(1); distance <= steps; distance++ {
		end := len(queue)
		added := int64(0)
		for head < end {
			index := queue[head]
			head++
			x, y := index%g.width, index/g.width
			for _, d := range gardenDirections {
				nx, ny := x+d.x, y+d.y
				if nx < 0 || ny < 0 || nx >= g.width || ny >= g.height {
					continue
				}
				next := ny*g.width + nx
				if !visited[next] && g.cells[next] != '#' {
					visited[next] = true
					queue = append(queue, next)
					added++
				}
			}
		}
		parity[distance%2] += added
		if added == 0 {
			break
		}
	}
	return parity[steps%2], nil
}

type gardenDistanceModel struct {
	radius, side int
	distances    []int32
}

func buildGardenModel(g *gardenMap, radius int) (gardenDistanceModel, error) {
	copies := 2*radius + 1
	side := g.width * copies
	if side > gardenPatchLimit/side {
		return gardenDistanceModel{}, errGardenLimit
	}
	model := gardenDistanceModel{radius: radius, side: side, distances: make([]int32, side*side)}
	for i := range model.distances {
		model.distances[i] = -1
	}
	shift := radius * g.width
	start := (g.start.y+shift)*side + g.start.x + shift
	model.distances[start] = 0
	queue := make([]int32, 1, side*side)
	queue[0] = int32(start)
	for head := 0; head < len(queue); head++ {
		index := int(queue[head])
		x, y := index%side, index/side
		for _, d := range gardenDirections {
			nx, ny := x+d.x, y+d.y
			if nx < 0 || ny < 0 || nx >= side || ny >= side {
				continue
			}
			next := ny*side + nx
			local := (ny%g.height)*g.width + nx%g.width
			if model.distances[next] < 0 && g.cells[local] != '#' {
				model.distances[next] = model.distances[index] + 1
				queue = append(queue, int32(next))
			}
		}
	}
	return model, nil
}
func (m *gardenDistanceModel) distance(g *gardenMap, p gardenPoint) int64 {
	localX, localY := gardenMod(p.x, g.width), gardenMod(p.y, g.width)
	tx, ty := (p.x-localX)/g.width, (p.y-localY)/g.width
	cx, cy := max(-m.radius, min(m.radius, tx)), max(-m.radius, min(m.radius, ty))
	col, row := (cx+m.radius)*g.width+localX, (cy+m.radius)*g.width+localY
	base := m.distances[row*m.side+col]
	if base < 0 {
		return -1
	}
	dx, dy := tx-cx, ty-cy
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	return int64(base) + int64(dx+dy)*int64(g.width)
}
func (m *gardenDistanceModel) validates(g *gardenMap) (bool, error) {
	ringSide := (2*m.radius + 3) * g.width
	if ringSide > gardenPatchLimit/ringSide {
		return false, errGardenLimit
	}
	low, high := -(m.radius+1)*g.width, (m.radius+2)*g.width
	// Outward translations add one tile width to both endpoints of an edge.
	// This ring covers their possible distance differences: edge inequalities bound
	// every path below, and descending neighbours supply a path to the unique zero.
	for y := low; y < high; y++ {
		for x := low; x < high; x++ {
			p := gardenPoint{x, y}
			distance := m.distance(g, p)
			if !g.open(p, true) {
				if distance >= 0 {
					return false, nil
				}
				continue
			}
			if distance == 0 && p != g.start || p == g.start && distance != 0 {
				return false, nil
			}
			descending := p == g.start
			for _, d := range gardenDirections {
				next := gardenPoint{x + d.x, y + d.y}
				if !g.open(next, true) {
					continue
				}
				n := m.distance(g, next)
				if distance < 0 && n < 0 {
					continue
				}
				if distance < 0 || n < 0 || distance-n > 1 || n-distance > 1 {
					return false, nil
				}
				descending = descending || n+1 == distance
			}
			if distance >= 0 && !descending {
				return false, nil
			}
		}
	}
	return true, nil
}

type certifiedGarden struct {
	garden *gardenMap
	model  gardenDistanceModel
}

func (g *gardenMap) certify() (certifiedGarden, error) {
	if g.width != g.height {
		return certifiedGarden{}, errors.New("large repeating counts require a square garden")
	}
	for _, radius := range []int{1, 2} {
		model, err := buildGardenModel(g, radius)
		if err != nil {
			return certifiedGarden{}, err
		}
		valid, err := model.validates(g)
		if err != nil {
			return certifiedGarden{}, err
		}
		if valid {
			return certifiedGarden{g, model}, nil
		}
	}
	return certifiedGarden{}, errGardenModel
}
func gardenProduct(a, b int64) (int64, error) {
	if a != 0 && b > math.MaxInt64/a {
		return 0, errGardenOverflow
	}
	return a * b, nil
}
func gardenTranslations(remaining, width int64, axes int) (int64, error) {
	if remaining < 0 {
		return 0, nil
	}
	copies, parity := remaining/width, remaining%2
	switch axes {
	case 0:
		if parity == 0 {
			return 1, nil
		}
		return 0, nil
	case 1:
		if width%2 == 0 {
			if parity == 0 {
				return copies, nil
			}
			return 0, nil
		}
		if parity == 0 {
			return copies / 2, nil
		}
		return copies/2 + copies%2, nil
	case 2:
		if width%2 == 0 {
			if parity != 0 || copies < 2 {
				return 0, nil
			}
			a, b := copies, copies-1
			if a%2 == 0 {
				a /= 2
			} else {
				b /= 2
			}
			return gardenProduct(a, b)
		}
		if parity == 0 {
			return gardenProduct(copies/2, copies/2)
		}
		pairs := max(copies-1, 0) / 2
		return gardenProduct(pairs, pairs+1)
	}
	return 0, errors.New("invalid translation dimension")
}
func (c *certifiedGarden) count(steps int64) (int64, error) {
	if steps < 0 {
		return 0, errors.New("negative step count")
	}
	g := c.garden
	if !g.hasEdge(true) {
		if steps == 0 {
			return 1, nil
		}
		return 0, nil
	}
	total := int64(0)
	width, radius := g.width, c.model.radius
	for index, distance := range c.model.distances {
		if distance < 0 {
			continue
		}
		tileX, tileY := index%c.model.side/width-radius, index/c.model.side/width-radius
		if tileX < 0 {
			tileX = -tileX
		}
		if tileY < 0 {
			tileY = -tileY
		}
		dimensions := [4]int{0, -1, -1, -1}
		if tileX == radius {
			dimensions[1] = 1
		}
		if tileY == radius {
			dimensions[2] = 1
		}
		if tileX == radius && tileY == radius {
			dimensions[3] = 2
		}
		for _, axes := range dimensions {
			if axes < 0 {
				continue
			}
			add, err := gardenTranslations(steps-int64(distance), int64(width), axes)
			if err != nil {
				return 0, err
			}
			if total > math.MaxInt64-add {
				return 0, errGardenOverflow
			}
			total += add
		}
	}
	return total, nil
}
func (Day21) Part1(input string) (string, error) {
	g, err := parseGarden(input)
	if err != nil {
		return "", err
	}
	count, err := g.finiteCount(64)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(count, 10), nil
}
func (Day21) Part2(input string) (string, error) {
	g, err := parseGarden(input)
	if err != nil {
		return "", err
	}
	certified, err := g.certify()
	if err != nil {
		return "", err
	}
	count, err := certified.count(26_501_365)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(count, 10), nil
}
func init() { solve.Register(Day21{}) }
