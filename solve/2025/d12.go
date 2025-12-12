package solve2025

import (
	"errors"
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"aoc/solve"
)

type Day12 struct {
}

func (d Day12) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2025, Day: 12}
}

func (d Day12) Part1(input string) (string, error) {
	shapes, regions, err := parseDay12(input)
	if err != nil {
		return "", err
	}

	numWorkers := runtime.NumCPU()
	var ok int64
	var processed int64

	type job struct {
		idx int
		r   day12Region
	}
	jobs := make(chan job, len(regions))
	var wg sync.WaitGroup

	start := time.Now()

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				result := canFitRegion(j.r.w, j.r.h, j.r.counts, shapes)
				if result {
					atomic.AddInt64(&ok, 1)
				}
				p := atomic.AddInt64(&processed, 1)
				if p%100 == 0 {
					fmt.Printf("Processed %d/%d regions (%.1fs elapsed)\n", p, len(regions), time.Since(start).Seconds())
				}
			}
		}()
	}

	for i, r := range regions {
		jobs <- job{idx: i, r: r}
	}
	close(jobs)
	wg.Wait()

	fmt.Printf("Total: %d/%d regions fit (%.1fs)\n", ok, len(regions), time.Since(start).Seconds())
	return strconv.FormatInt(ok, 10), nil
}

func (d Day12) Part2(input string) (string, error) {
	return "", errors.New("not implemented")
}

func init() {
	solve.Register(Day12{})
}

type day12Point struct {
	x int
	y int
}

type day12Shape struct {
	index        int
	cells        []day12Point // normalized
	cellCount    int
	orientations [][]day12Point // each normalized
}

type day12Region struct {
	w      int
	h      int
	counts []int // counts by shape index
}

func parseDay12(input string) ([]day12Shape, []day12Region, error) {
	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")

	shapesByIdx := map[int]day12Shape{}
	maxIdx := -1
	idx := 0
	for idx < len(lines) {
		line := strings.TrimSpace(lines[idx])
		if line == "" {
			idx++
			continue
		}

		if _, _, _, ok := parseDay12RegionLine(line); ok {
			break
		}

		shapeIndex, ok := parseDay12ShapeHeader(line)
		if !ok {
			return nil, nil, fmt.Errorf("day12: expected shape header or region line, got %q", line)
		}
		idx++

		var grid []string
		for idx < len(lines) {
			row := strings.TrimSpace(lines[idx])
			if row == "" {
				break
			}
			for _, ch := range row {
				if ch != '.' && ch != '#' {
					return nil, nil, fmt.Errorf("day12: invalid shape row %q", row)
				}
			}
			grid = append(grid, row)
			idx++
		}

		if len(grid) == 0 {
			return nil, nil, fmt.Errorf("day12: shape %d has no grid", shapeIndex)
		}

		cells := make([]day12Point, 0)
		for y, row := range grid {
			for x := 0; x < len(row); x++ {
				if row[x] == '#' {
					cells = append(cells, day12Point{x: x, y: y})
				}
			}
		}
		if len(cells) == 0 {
			return nil, nil, fmt.Errorf("day12: shape %d has zero # cells", shapeIndex)
		}

		cells = day12NormalizePoints(cells)
		shape := day12Shape{index: shapeIndex, cells: cells, cellCount: len(cells)}
		shape.orientations = day12GenerateOrientations(shape.cells)
		shapesByIdx[shapeIndex] = shape
		if shapeIndex > maxIdx {
			maxIdx = shapeIndex
		}

		for idx < len(lines) && strings.TrimSpace(lines[idx]) == "" {
			idx++
		}
	}

	if maxIdx < 0 {
		return nil, nil, fmt.Errorf("day12: no shapes found")
	}

	shapes := make([]day12Shape, maxIdx+1)
	for i := 0; i <= maxIdx; i++ {
		shape, ok := shapesByIdx[i]
		if !ok {
			return nil, nil, fmt.Errorf("day12: missing shape index %d", i)
		}
		shapes[i] = shape
	}

	regions := make([]day12Region, 0)
	for idx < len(lines) {
		line := strings.TrimSpace(lines[idx])
		idx++
		if line == "" {
			continue
		}
		w, h, counts, ok := parseDay12RegionLine(line)
		if !ok {
			return nil, nil, fmt.Errorf("day12: expected region line, got %q", line)
		}
		if w <= 0 || h <= 0 {
			return nil, nil, fmt.Errorf("day12: invalid region dimensions %dx%d", w, h)
		}
		if len(counts) > len(shapes) {
			return nil, nil, fmt.Errorf("day12: region references %d shapes but only %d defined", len(counts), len(shapes))
		}
		if len(counts) < len(shapes) {
			padded := make([]int, len(shapes))
			copy(padded, counts)
			counts = padded
		}
		regions = append(regions, day12Region{w: w, h: h, counts: counts})
	}

	return shapes, regions, nil
}

func parseDay12ShapeHeader(line string) (int, bool) {
	line = strings.TrimSpace(line)
	if !strings.HasSuffix(line, ":") {
		return 0, false
	}
	nStr := strings.TrimSpace(strings.TrimSuffix(line, ":"))
	if nStr == "" {
		return 0, false
	}
	n, err := strconv.Atoi(nStr)
	if err != nil || n < 0 {
		return 0, false
	}
	return n, true
}

func parseDay12RegionLine(line string) (w, h int, counts []int, ok bool) {
	line = strings.TrimSpace(line)
	colon := strings.IndexByte(line, ':')
	if colon < 0 {
		return 0, 0, nil, false
	}
	dims := strings.TrimSpace(line[:colon])
	rest := strings.TrimSpace(line[colon+1:])

	x := strings.IndexByte(dims, 'x')
	if x < 0 {
		return 0, 0, nil, false
	}
	wStr := strings.TrimSpace(dims[:x])
	hStr := strings.TrimSpace(dims[x+1:])
	if wStr == "" || hStr == "" {
		return 0, 0, nil, false
	}

	wParsed, err := strconv.Atoi(wStr)
	if err != nil {
		return 0, 0, nil, false
	}
	hParsed, err := strconv.Atoi(hStr)
	if err != nil {
		return 0, 0, nil, false
	}

	counts = []int{}
	if rest != "" {
		fields := strings.Fields(rest)
		counts = make([]int, 0, len(fields))
		for _, f := range fields {
			v, err := strconv.Atoi(f)
			if err != nil {
				return 0, 0, nil, false
			}
			if v < 0 {
				return 0, 0, nil, false
			}
			counts = append(counts, v)
		}
	}

	return wParsed, hParsed, counts, true
}

type day12PieceType struct {
	shapeIdx  int
	cellCount int
	cellIndex [][]int    // cellIndex[cell] = list of placement indices covering that cell
	plMinCell []int      // plMinCell[pi] = minimum cell index covered by placement pi
	plMulti   [][]uint64 // multi-word placements
}

func canFitRegion(w, h int, counts []int, shapes []day12Shape) bool {
	area := w * h
	total := 0
	for i, c := range counts {
		if c == 0 {
			continue
		}
		if i < 0 || i >= len(shapes) {
			return false
		}
		total += c * shapes[i].cellCount
		if total > area {
			return false
		}
	}

	pieceTypes := make([]day12PieceType, 0)
	countsCopy := make([]int, 0)
	for i, c := range counts {
		if c == 0 {
			continue
		}
		pt := day12BuildPieceType(w, h, area, shapes[i])
		if len(pt.plMulti) == 0 {
			return false
		}
		pt.shapeIdx = i
		pieceTypes = append(pieceTypes, pt)
		countsCopy = append(countsCopy, c)
	}

	// Sort by cell count descending for better pruning.
	for i := 0; i < len(pieceTypes)-1; i++ {
		for j := i + 1; j < len(pieceTypes); j++ {
			if pieceTypes[j].cellCount > pieceTypes[i].cellCount {
				pieceTypes[i], pieceTypes[j] = pieceTypes[j], pieceTypes[i]
				countsCopy[i], countsCopy[j] = countsCopy[j], countsCopy[i]
			}
		}
	}

	words := (area + 63) / 64
	occupied := make([]uint64, words)
	// minPlacement[ti] = minimum placement index that next piece of type ti can use (symmetry breaking)
	minPlacement := make([]int, len(pieceTypes))
	return day12SearchOpt(countsCopy, pieceTypes, occupied, minPlacement)
}

// day12SearchOpt: piece-centric search with MRV heuristic + symmetry breaking.
func day12SearchOpt(counts []int, pts []day12PieceType, occupied []uint64, minPlacement []int) bool {
	// Find piece type with fewest valid placements (MRV heuristic).
	bestType := -1
	bestCount := 1 << 30

	for ti, pt := range pts {
		if counts[ti] <= 0 {
			continue
		}
		cnt := 0
		threshold := minPlacement[ti]
		for pi := threshold; pi < len(pt.plMulti); pi++ {
			if !day12IntersectsMulti(occupied, pt.plMulti[pi]) {
				cnt++
				if cnt > bestCount {
					break // Early exit: this type has more options than current best
				}
			}
		}
		if cnt == 0 {
			return false // No valid placements for a piece we need to place
		}
		if cnt < bestCount {
			bestCount = cnt
			bestType = ti
			if bestCount == 1 {
				break // Can't do better than 1
			}
		}
	}

	if bestType < 0 {
		return true
	}

	pt := pts[bestType]
	threshold := minPlacement[bestType]
	counts[bestType]--

	for pi := threshold; pi < len(pt.plMulti); pi++ {
		bits := pt.plMulti[pi]
		if day12IntersectsMulti(occupied, bits) {
			continue
		}

		day12OrInMulti(occupied, bits)
		oldThreshold := minPlacement[bestType]
		minPlacement[bestType] = pi + 1 // Symmetry: next piece of same type must use higher index

		if day12SearchOpt(counts, pts, occupied, minPlacement) {
			return true
		}

		minPlacement[bestType] = oldThreshold
		day12AndNotInMulti(occupied, bits)
	}

	counts[bestType]++
	return false
}

func day12IntersectsMulti(a, b []uint64) bool {
	for i := range a {
		if (a[i] & b[i]) != 0 {
			return true
		}
	}
	return false
}

func day12OrInMulti(dst, src []uint64) {
	for i := range dst {
		dst[i] |= src[i]
	}
}

func day12AndNotInMulti(dst, src []uint64) {
	for i := range dst {
		dst[i] &^= src[i]
	}
}

func day12BuildPieceType(w, h, area int, shape day12Shape) day12PieceType {
	pt := day12PieceType{
		cellCount: shape.cellCount,
		cellIndex: make([][]int, area),
	}
	for c := 0; c < area; c++ {
		pt.cellIndex[c] = []int{}
	}

	words := (area + 63) / 64

	plIdx := 0
	for _, ori := range shape.orientations {
		maxX, maxY := 0, 0
		for _, p := range ori {
			if p.x > maxX {
				maxX = p.x
			}
			if p.y > maxY {
				maxY = p.y
			}
		}
		shapeW := maxX + 1
		shapeH := maxY + 1
		if shapeW > w || shapeH > h {
			continue
		}

		for oy := 0; oy <= h-shapeH; oy++ {
			for ox := 0; ox <= w-shapeW; ox++ {
				bitsMulti := make([]uint64, words)

				cells := make([]int, 0, len(ori))
				minCell := area // Track minimum cell for symmetry breaking
				for _, p := range ori {
					x := ox + p.x
					y := oy + p.y
					idx := y*w + x
					cells = append(cells, idx)
					bitsMulti[idx/64] |= 1 << uint(idx%64)
					if idx < minCell {
						minCell = idx
					}
				}

				pt.plMulti = append(pt.plMulti, bitsMulti)
				pt.plMinCell = append(pt.plMinCell, minCell)

				for _, c := range cells {
					pt.cellIndex[c] = append(pt.cellIndex[c], plIdx)
				}
				plIdx++
			}
		}
	}

	return pt
}

func day12GenerateOrientations(cells []day12Point) [][]day12Point {
	seen := map[string]struct{}{}
	oris := make([][]day12Point, 0, 8)

	for rot := 0; rot < 4; rot++ {
		for flip := 0; flip < 2; flip++ {
			pts := make([]day12Point, len(cells))
			for i, p := range cells {
				x, y := p.x, p.y
				for r := 0; r < rot; r++ {
					x, y = -y, x
				}
				if flip == 1 {
					x = -x
				}
				pts[i] = day12Point{x: x, y: y}
			}

			pts = day12NormalizePoints(pts)
			key := day12EncodePoints(pts)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			oris = append(oris, pts)
		}
	}

	return oris
}

func day12NormalizePoints(pts []day12Point) []day12Point {
	minX, minY := pts[0].x, pts[0].y
	for _, p := range pts {
		if p.x < minX {
			minX = p.x
		}
		if p.y < minY {
			minY = p.y
		}
	}
	for i := range pts {
		pts[i].x -= minX
		pts[i].y -= minY
	}
	sort.Slice(pts, func(i, j int) bool {
		if pts[i].y == pts[j].y {
			return pts[i].x < pts[j].x
		}
		return pts[i].y < pts[j].y
	})
	return pts
}

func day12EncodePoints(pts []day12Point) string {
	// Small N (polyomino cell counts), so string encoding is fine.
	var b strings.Builder
	for i, p := range pts {
		if i > 0 {
			b.WriteByte(';')
		}
		b.WriteString(strconv.Itoa(p.x))
		b.WriteByte(',')
		b.WriteString(strconv.Itoa(p.y))
	}
	return b.String()
}
