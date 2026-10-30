package solve2020

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day20 struct{}

func (Day20) Coords() solve.SolutionCoords { return solve.SolutionCoords{Year: 2020, Day: 20} }

const jigsawTileLimit = 1024
const jigsawWorkLimit = 1_000_000

var errJigsawSearchLimit = errors.New("jigsaw assembly search limit exceeded")

type jigsawTile struct {
	id      int64
	cells   [100]byte
	edges   [4]uint16
	matched [4]bool
}

// Edges use top/right/bottom/left order; horizontal edges run left to right,
// vertical edges run top to bottom. Canonical keys are only candidate filters.
func (t *jigsawTile) readEdges() {
	t.edges = [4]uint16{}
	for i := 0; i < 10; i++ {
		for side, cell := range [4]byte{t.cells[i], t.cells[i*10+9], t.cells[90+i], t.cells[i*10]} {
			t.edges[side] <<= 1
			if cell == '#' {
				t.edges[side] |= 1
			}
		}
	}
}

func canonicalJigsawEdge(edge uint16) uint16 {
	forward := edge
	var reversed uint16
	for i := 0; i < 10; i++ {
		reversed = reversed<<1 | edge&1
		edge >>= 1
	}
	return min(forward, reversed)
}

func (t jigsawTile) rotate() jigsawTile {
	result := t
	for row := 0; row < 10; row++ {
		for col := 0; col < 10; col++ {
			result.cells[col*10+9-row] = t.cells[row*10+col]
		}
	}
	result.readEdges()
	result.matched = [4]bool{t.matched[3], t.matched[0], t.matched[1], t.matched[2]}
	return result
}

func (t jigsawTile) flip() jigsawTile {
	result := t
	for row := 0; row < 10; row++ {
		for col := 0; col < 10; col++ {
			result.cells[row*10+9-col] = t.cells[row*10+col]
		}
	}
	result.readEdges()
	result.matched = [4]bool{t.matched[0], t.matched[3], t.matched[2], t.matched[1]}
	return result
}

func (t jigsawTile) orientations() []jigsawTile {
	result := make([]jigsawTile, 0, 8)
	for flip := 0; flip < 2; flip++ {
		current := t
		if flip == 1 {
			current = current.flip()
		}
		for rotation := 0; rotation < 4; rotation++ {
			duplicate := false
			for _, other := range result {
				if other.cells == current.cells {
					duplicate = true
					break
				}
			}
			if !duplicate {
				result = append(result, current)
			}
			current = current.rotate()
		}
	}
	return result
}

func (t jigsawTile) corner() bool {
	unmatched := 0
	for _, matched := range t.matched {
		if !matched {
			unmatched++
		}
	}
	if unmatched != 2 {
		return false
	}
	for side := 0; side < 4; side++ {
		if !t.matched[side] && !t.matched[(side+1)%4] {
			return true
		}
	}
	return false
}

type jigsaw struct {
	tiles   []jigsawTile
	size    int
	borders map[uint16][]int
}

func parseJigsaw(input string) (*jigsaw, error) {
	lines := strings.Split(input, "\n")
	for i := range lines {
		lines[i] = strings.TrimSuffix(lines[i], "\r")
	}
	j := &jigsaw{}
	ids := make(map[int64]bool)
	for line := 0; line < len(lines); {
		if lines[line] == "" {
			line++
			continue
		}
		headerLine := line + 1
		invalid := func(message string) error { return fmt.Errorf("line %d: %s", headerLine, message) }
		header, ok := strings.CutPrefix(lines[line], "Tile ")
		if !ok {
			return nil, invalid("expected 'Tile <positive ID>:'")
		}
		number, ok := strings.CutSuffix(header, ":")
		if !ok || number == "" {
			return nil, invalid("expected 'Tile <positive ID>:'")
		}
		for _, c := range number {
			if c < '0' || c > '9' {
				return nil, invalid("tile ID must be positive decimal")
			}
		}
		id, err := strconv.ParseInt(number, 10, 64)
		if err != nil || id <= 0 {
			return nil, invalid("tile ID must be a positive int64")
		}
		if ids[id] {
			return nil, invalid("duplicate tile ID")
		}
		ids[id] = true
		if len(j.tiles) == jigsawTileLimit {
			return nil, invalid("jigsaw supports at most 1024 tiles")
		}
		line++
		t := jigsawTile{id: id}
		for row := 0; row < 10; row++ {
			if line == len(lines) {
				return nil, invalid("tile must contain ten rows")
			}
			if len(lines[line]) != 10 {
				return nil, fmt.Errorf("line %d: tile row must contain ten cells", line+1)
			}
			for col, c := range []byte(lines[line]) {
				if c != '.' && c != '#' {
					return nil, fmt.Errorf("line %d: expected only '.' or '#' cells", line+1)
				}
				t.cells[row*10+col] = c
			}
			line++
		}
		if line < len(lines) && lines[line] != "" {
			return nil, fmt.Errorf("line %d: expected blank line between tiles", line+1)
		}
		t.readEdges()
		j.tiles = append(j.tiles, t)
	}
	if len(j.tiles) == 0 {
		return nil, errors.New("jigsaw input is empty")
	}
	for j.size = 1; j.size*j.size < len(j.tiles); j.size++ {
	}
	if j.size*j.size != len(j.tiles) {
		return nil, errors.New("tile count must be a perfect square")
	}
	sort.Slice(j.tiles, func(a, b int) bool { return j.tiles[a].id < j.tiles[b].id })
	j.indexBorders()
	return j, nil
}

func (j *jigsaw) indexBorders() {
	j.borders = make(map[uint16][]int)
	for index, t := range j.tiles {
		for _, edge := range t.edges {
			key := canonicalJigsawEdge(edge)
			indices := j.borders[key]
			// Tile indices arrive in order; repeated sides must not duplicate ownership.
			if len(indices) == 0 || indices[len(indices)-1] != index {
				j.borders[key] = append(indices, index)
			}
		}
	}
	for index := range j.tiles {
		for side, edge := range j.tiles[index].edges {
			j.tiles[index].matched[side] = len(j.borders[canonicalJigsawEdge(edge)]) > 1
		}
	}
}

func (j *jigsaw) corners() []int {
	var corners []int
	for index, t := range j.tiles {
		if t.corner() {
			corners = append(corners, index)
		}
	}
	return corners
}

func (j *jigsaw) cornerProduct() (int64, error) {
	corners := j.corners()
	if len(corners) != 4 {
		return 0, errors.New("expected exactly four tiles with two adjacent unmatched borders")
	}
	product := int64(1)
	for _, index := range corners {
		id := j.tiles[index].id
		if product > math.MaxInt64/id {
			return 0, errors.New("corner product exceeds int64")
		}
		product *= id
	}
	return product, nil
}

type jigsawPlacement struct{ tile, orientation int }

func spendJigsawWork(remaining *int) error {
	if *remaining <= 0 {
		return errJigsawSearchLimit
	}
	*remaining -= 1
	return nil
}

func (j *jigsaw) candidates(variants [][]jigsawTile, placed []jigsawTile, used []bool, roots []int, work *int) ([]jigsawPlacement, error) {
	row, col := len(placed)/j.size, len(placed)%j.size
	indices := roots
	if col > 0 {
		indices = j.borders[canonicalJigsawEdge(placed[len(placed)-1].edges[1])]
	} else if row > 0 {
		indices = j.borders[canonicalJigsawEdge(placed[len(placed)-j.size].edges[2])]
	}
	var result []jigsawPlacement
	for _, index := range indices {
		if err := spendJigsawWork(work); err != nil {
			return nil, err
		}
		if used[index] {
			continue
		}
		for orientation, t := range variants[index] {
			if err := spendJigsawWork(work); err != nil {
				return nil, err
			}
			if row > 0 && t.edges[0] != placed[len(placed)-j.size].edges[2] || col > 0 && t.edges[3] != placed[len(placed)-1].edges[1] {
				continue
			}
			if row > 0 && !t.matched[0] || col+1 < j.size && !t.matched[1] || row+1 < j.size && !t.matched[2] || col > 0 && !t.matched[3] {
				continue
			}
			result = append(result, jigsawPlacement{index, orientation})
		}
	}
	// Pop from the end while preserving deterministic tile/orientation order.
	for a, b := 0, len(result)-1; a < b; a, b = a+1, b-1 {
		result[a], result[b] = result[b], result[a]
	}
	return result, nil
}

func (j *jigsaw) assemble(work int) ([]jigsawTile, error) {
	variants := make([][]jigsawTile, len(j.tiles))
	roots := make([]int, len(j.tiles))
	for i, t := range j.tiles {
		variants[i] = t.orientations()
		roots[i] = i
	}
	if corners := j.corners(); len(corners) == 4 {
		roots = corners[:1]
	}
	used := make([]bool, len(j.tiles))
	placed := make([]jigsawTile, 0, len(j.tiles))
	indices := make([]int, 0, len(j.tiles))
	first, err := j.candidates(variants, placed, used, roots, &work)
	if err != nil {
		return nil, err
	}
	choices := [][]jigsawPlacement{first}
	for len(choices) > 0 {
		depth := len(choices) - 1
		remaining := choices[depth]
		if len(remaining) == 0 {
			choices = choices[:depth]
			if len(indices) > 0 {
				index := indices[len(indices)-1]
				used[index] = false
				indices = indices[:len(indices)-1]
				placed = placed[:len(placed)-1]
			}
			continue
		}
		choice := remaining[len(remaining)-1]
		choices[depth] = remaining[:len(remaining)-1]
		used[choice.tile] = true
		placed = append(placed, variants[choice.tile][choice.orientation])
		indices = append(indices, choice.tile)
		if len(placed) == len(j.tiles) {
			return placed, nil
		}
		next, err := j.candidates(variants, placed, used, roots, &work)
		if err != nil {
			return nil, err
		}
		choices = append(choices, next)
	}
	return nil, errors.New("tiles cannot form a square image")
}

type jigsawImage struct {
	width, height int
	cells         []byte
}

func renderJigsaw(placed []jigsawTile, size int) jigsawImage {
	img := jigsawImage{width: size * 8, height: size * 8, cells: make([]byte, size*size*64)}
	for index, t := range placed {
		row, col := index/size, index%size
		for r := 0; r < 8; r++ {
			copy(img.cells[(row*8+r)*img.width+col*8:], t.cells[(r+1)*10+1:(r+1)*10+9])
		}
	}
	return img
}

type monsterPoint struct{ row, col int }
type jigsawMonster struct {
	height, width int
	pixels        [15]monsterPoint
}

var jigsawMonsterPixels = [15]monsterPoint{{0, 18}, {1, 0}, {1, 5}, {1, 6}, {1, 11}, {1, 12}, {1, 17}, {1, 18}, {1, 19}, {2, 1}, {2, 4}, {2, 7}, {2, 10}, {2, 13}, {2, 16}}

func monsterOrientations() [8]jigsawMonster {
	original := jigsawMonster{height: 3, width: 20, pixels: jigsawMonsterPixels}
	var result [8]jigsawMonster
	for flip := 0; flip < 2; flip++ {
		current := original
		if flip == 1 {
			for i, p := range current.pixels {
				current.pixels[i].col = current.width - 1 - p.col
			}
		}
		for rotation := 0; rotation < 4; rotation++ {
			result[flip*4+rotation] = current
			for i, p := range current.pixels {
				current.pixels[i] = monsterPoint{p.col, current.height - 1 - p.row}
			}
			current.height, current.width = current.width, current.height
		}
	}
	return result
}

func (m jigsawMonster) coverage(img jigsawImage, mask []bool) int {
	if img.height < m.height || img.width < m.width {
		return 0
	}
	clear(mask)
	covered := 0
	for row := 0; row <= img.height-m.height; row++ {
		for col := 0; col <= img.width-m.width; col++ {
			matches := true
			for _, p := range m.pixels {
				if img.cells[(row+p.row)*img.width+col+p.col] != '#' {
					matches = false
					break
				}
			}
			if matches {
				for _, p := range m.pixels {
					index := (row+p.row)*img.width + col + p.col
					if !mask[index] {
						mask[index] = true
						covered++
					}
				}
			}
		}
	}
	return covered
}

func jigsawRoughness(img jigsawImage) int {
	total := 0
	for _, cell := range img.cells {
		if cell == '#' {
			total++
		}
	}
	mask := make([]bool, len(img.cells))
	best := 0
	for _, monster := range monsterOrientations() {
		best = max(best, monster.coverage(img, mask))
	}
	return total - best
}

func (Day20) Part1(input string) (string, error) {
	j, err := parseJigsaw(input)
	if err != nil {
		return "", err
	}
	product, err := j.cornerProduct()
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(product, 10), nil
}

func (Day20) Part2(input string) (string, error) {
	j, err := parseJigsaw(input)
	if err != nil {
		return "", err
	}
	placed, err := j.assemble(jigsawWorkLimit)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(jigsawRoughness(renderJigsaw(placed, j.size))), nil
}

func init() { solve.Register(Day20{}) }
