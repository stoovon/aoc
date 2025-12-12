package solve2025

import (
	"errors"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day9 struct {
}

func (d Day9) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2025, Day: 9}
}

func (d Day9) Part1(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	type point struct {
		x, y int64
	}
	var pts []point

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			return "", errors.New("invalid line: " + line)
		}
		xStr := strings.TrimSpace(parts[0])
		yStr := strings.TrimSpace(parts[1])

		x, err := strconv.ParseInt(xStr, 10, 64)
		if err != nil {
			return "", err
		}
		y, err := strconv.ParseInt(yStr, 10, 64)
		if err != nil {
			return "", err
		}
		pts = append(pts, point{x: x, y: y})
	}

	var best int64
	for i := 0; i < len(pts); i++ {
		for j := i + 1; j < len(pts); j++ {
			dx := pts[i].x - pts[j].x
			if dx < 0 {
				dx = -dx
			}
			dy := pts[i].y - pts[j].y
			if dy < 0 {
				dy = -dy
			}
			area := (dx + 1) * (dy + 1)
			if area > best {
				best = area
			}
		}
	}

	return strconv.FormatInt(best, 10), nil
}

func (d Day9) Part2(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	type point struct {
		x, y int64
	}
	var verts []point

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			return "", errors.New("invalid line: " + line)
		}
		xStr := strings.TrimSpace(parts[0])
		yStr := strings.TrimSpace(parts[1])

		x, err := strconv.ParseInt(xStr, 10, 64)
		if err != nil {
			return "", err
		}
		y, err := strconv.ParseInt(yStr, 10, 64)
		if err != nil {
			return "", err
		}
		verts = append(verts, point{x: x, y: y})
	}

	if len(verts) < 2 {
		return "0", nil
	}

	// Build row spans using scanline over vertical edges.
	// For each row, store single [l,r] span (simple polygon assumption).
	type span struct {
		l, r int64
	}
	rowSpan := make(map[int64]span)

	n := len(verts)
	rowEdges := make(map[int64][]int64)
	for i := 0; i < n; i++ {
		p1 := verts[i]
		p2 := verts[(i+1)%n]
		if p1.x == p2.x {
			x := p1.x
			y1, y2 := p1.y, p2.y
			if y1 > y2 {
				y1, y2 = y2, y1
			}
			for y := y1; y < y2; y++ {
				rowEdges[y] = append(rowEdges[y], x)
			}
		}
	}

	for y, xs := range rowEdges {
		if len(xs) < 2 {
			continue
		}
		minX, maxX := xs[0], xs[0]
		for _, x := range xs[1:] {
			if x < minX {
				minX = x
			}
			if x > maxX {
				maxX = x
			}
		}
		rowSpan[y] = span{l: minX, r: maxX}
	}

	var yMin, yMax int64
	first := true
	for y := range rowSpan {
		if first {
			yMin, yMax = y, y
			first = false
		} else {
			if y < yMin {
				yMin = y
			}
			if y > yMax {
				yMax = y
			}
		}
	}

	if first {
		return "0", nil
	}

	// Build prefix structures for O(1) range queries:
	// For any y-range [yLo, yHi], we need:
	//   max of all rowSpan[y].l for y in [yLo, yHi]  (tightest left bound)
	//   min of all rowSpan[y].r for y in [yLo, yHi]  (tightest right bound)
	// If maxL <= xLo and xHi <= minR, rectangle is valid.

	height := yMax - yMin + 1
	// Arrays indexed by (y - yMin)
	leftBound := make([]int64, height)  // rowSpan[y].l or a sentinel
	rightBound := make([]int64, height) // rowSpan[y].r or a sentinel

	const bigVal = int64(1 << 60)
	for i := range leftBound {
		leftBound[i] = -bigVal // will take max, so start very negative
		rightBound[i] = bigVal // will take min, so start very positive
	}

	for y, sp := range rowSpan {
		idx := y - yMin
		leftBound[idx] = sp.l
		rightBound[idx] = sp.r
	}

	// Build sparse table for range-max of leftBound and range-min of rightBound
	// Using simple O(H) scan per query is still too slow for O(n^2) pairs.
	// Instead, use a sparse table for O(1) queries after O(H log H) preprocessing.

	logH := 0
	for (1 << logH) <= int(height) {
		logH++
	}

	// sparseMaxL[k][i] = max of leftBound[i..i+2^k-1]
	// sparseMinR[k][i] = min of rightBound[i..i+2^k-1]
	sparseMaxL := make([][]int64, logH)
	sparseMinR := make([][]int64, logH)
	for k := 0; k < logH; k++ {
		sparseMaxL[k] = make([]int64, height)
		sparseMinR[k] = make([]int64, height)
	}

	for i := range height {
		sparseMaxL[0][i] = leftBound[i]
		sparseMinR[0][i] = rightBound[i]
	}

	for k := 1; k < logH; k++ {
		span := int64(1) << (k - 1)
		for i := int64(0); i+span < height; i++ {
			sparseMaxL[k][i] = sparseMaxL[k-1][i]
			if sparseMaxL[k-1][i+span] > sparseMaxL[k][i] {
				sparseMaxL[k][i] = sparseMaxL[k-1][i+span]
			}
			sparseMinR[k][i] = sparseMinR[k-1][i]
			if sparseMinR[k-1][i+span] < sparseMinR[k][i] {
				sparseMinR[k][i] = sparseMinR[k-1][i+span]
			}
		}
	}

	queryMaxL := func(lo, hi int64) int64 {
		// max over [lo, hi] inclusive, indices relative to yMin
		if lo > hi {
			return -bigVal
		}
		length := hi - lo + 1
		k := 0
		for (int64(1) << (k + 1)) <= length {
			k++
		}
		span := int64(1) << k
		a := sparseMaxL[k][lo]
		b := sparseMaxL[k][hi-span+1]
		if b > a {
			return b
		}
		return a
	}

	queryMinR := func(lo, hi int64) int64 {
		// min over [lo, hi] inclusive
		if lo > hi {
			return bigVal
		}
		length := hi - lo + 1
		k := 0
		for (int64(1) << (k + 1)) <= length {
			k++
		}
		span := int64(1) << k
		a := sparseMinR[k][lo]
		b := sparseMinR[k][hi-span+1]
		if b < a {
			return b
		}
		return a
	}

	var best int64
	for i := 0; i < len(verts); i++ {
		for j := i + 1; j < len(verts); j++ {
			x1, y1 := verts[i].x, verts[i].y
			x2, y2 := verts[j].x, verts[j].y

			var xLo, xHi, yLo, yHi int64
			if x1 < x2 {
				xLo, xHi = x1, x2
			} else {
				xLo, xHi = x2, x1
			}
			if y1 < y2 {
				yLo, yHi = y1, y2
			} else {
				yLo, yHi = y2, y1
			}

			width := xHi - xLo + 1
			rectHeight := yHi - yLo + 1
			area := width * rectHeight
			if area <= best {
				continue
			}

			idxLo := yLo - yMin
			idxHi := yHi - yMin

			if idxLo < 0 || idxHi >= height {
				continue
			}

			maxL := queryMaxL(idxLo, idxHi)
			minR := queryMinR(idxLo, idxHi)

			// Rectangle [xLo, xHi] x [yLo, yHi] is valid iff
			// for every row y in [yLo, yHi]: rowSpan[y].l <= xLo AND xHi <= rowSpan[y].r
			// Equivalently: max(rowSpan[y].l) <= xLo AND xHi <= min(rowSpan[y].r)
			if maxL <= xLo && xHi <= minR {
				best = area
			}
		}
	}

	return strconv.FormatInt(best, 10), nil
}

func init() {
	solve.Register(Day9{})
}
