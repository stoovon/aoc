package solve2025

import (
	"fmt"
	"math"
	"math/bits"
	"slices"
	"strings"
	"testing"
)

func packingShape(cells []day12Point) day12Shape {
	cells = day12NormalizePoints(slices.Clone(cells))
	return day12Shape{cells, day12GenerateOrientations(cells)}
}

func packingOraclePlacements(cells []day12Point, w, h int) []uint64 {
	unique := make(map[uint64]bool)
	for flipped := 0; flipped < 2; flipped++ {
		for turn := 0; turn < 4; turn++ {
			transformed := make([]day12Point, len(cells))
			for i, p := range cells {
				if flipped == 1 {
					p.y = -p.y
				}
				for r := 0; r < turn; r++ {
					p.x, p.y = p.y, -p.x
				}
				transformed[i] = p
			}
			loX, loY := transformed[0].x, transformed[0].y
			hiX, hiY := loX, loY
			for _, p := range transformed {
				loX = min(loX, p.x)
				loY = min(loY, p.y)
				hiX = max(hiX, p.x)
				hiY = max(hiY, p.y)
			}
			for y := 0; y+hiY-loY < h; y++ {
				for x := 0; x+hiX-loX < w; x++ {
					var placement uint64
					for _, p := range transformed {
						placement |= uint64(1) << uint((y+p.y-loY)*w+x+p.x-loX)
					}
					unique[placement] = true
				}
			}
		}
	}
	var result []uint64
	for placement := range unique {
		result = append(result, placement)
	}
	return result
}

func packingOracle(w, h int, counts []int, shapes []day12Shape) bool {
	needed := 0
	var pieces [][]uint64
	for i, count := range counts {
		needed += count * len(shapes[i].cells)
		placements := packingOraclePlacements(shapes[i].cells, w, h)
		for range count {
			pieces = append(pieces, placements)
		}
	}
	if needed > w*h {
		return false
	}
	var search func(int, uint64) bool
	search = func(index int, occupied uint64) bool {
		if index == len(pieces) {
			return true
		}
		for _, placement := range pieces[index] {
			if occupied&placement == 0 && search(index+1, occupied|placement) {
				return true
			}
		}
		return false
	}
	return search(0, 0)
}

func TestDay12AllSmallShapesAgainstIndependentPacking(t *testing.T) {
	for mask := 1; mask < 64; mask++ {
		var cells []day12Point
		for index := 0; index < 6; index++ {
			if mask&(1<<index) != 0 {
				cells = append(cells, day12Point{index % 3, index / 3})
			}
		}
		shapes := []day12Shape{packingShape(cells)}
		for w := 1; w <= 4; w++ {
			for h := 1; h <= 3; h++ {
				for count := 0; count <= 3; count++ {
					expected := packingOracle(w, h, []int{count}, shapes)
					actual, err := canFitRegion(day12Region{w: w, h: h, counts: []int{count}}, shapes)
					if err != nil || actual != expected {
						t.Fatalf("mask%d %dx%d count%d got%t,%v want%t", mask, w, h, count, actual, err, expected)
					}
					exact, err := day12ExactFit(day12Region{w: w, h: h, counts: []int{count}}, shapes, 200000)
					if err != nil || exact != expected {
						t.Fatalf("exact mask%d %dx%d count%d got%t,%v want%t", mask, w, h, count, exact, err, expected)
					}
				}
			}
		}
	}
}

func TestDay12MixedPiecesAgainstIndependentPacking(t *testing.T) {
	shapes := []day12Shape{
		packingShape([]day12Point{{0, 0}, {1, 0}}),
		packingShape([]day12Point{{0, 0}, {0, 1}, {1, 1}}),
		packingShape([]day12Point{{0, 0}, {1, 0}, {0, 1}, {1, 1}}),
	}
	for code := 0; code < 27; code++ {
		value := code
		counts := make([]int, 3)
		for i := range counts {
			counts[i] = value % 3
			value /= 3
		}
		for w := 1; w <= 4; w++ {
			for h := 1; h <= 3; h++ {
				expected := packingOracle(w, h, counts, shapes)
				actual, err := canFitRegion(day12Region{w: w, h: h, counts: counts}, shapes)
				if err != nil || actual != expected {
					t.Fatalf("mixed%v %dx%d got%t,%v want%t", counts, w, h, actual, err, expected)
				}
			}
		}
	}
}

func TestDay12AreaIsNotPackingAndRotationWorks(t *testing.T) {
	cases := []struct{ input, want string }{
		{"0:\n###\n\n2x2: 1\n", "0"},
		{"0:\n###\n\n1x3: 1\n", "1"},
		{"0:\n#.\n##\n\n3x2: 2\n", "1"},
		{"0:\n##\n##\n\n1x4: 1\n", "0"},
		{"0:\n#######\n\n9x8: 10\n", "1"},
		{"0:\n#######\n\n8x8: 10\n", "0"},
		{"0:\n.#.\n...\n\n2x2: 4\n2x2: 5\n2x2: 0\n", "2"},
	}
	for _, test := range cases {
		actual, err := (Day12{}).Part1(test.input)
		if err != nil || actual != test.want {
			t.Fatalf("%q got%s,%v want%s", test.input, actual, err, test.want)
		}
	}
	input := strings.ReplaceAll(cases[2].input, "\n", "\r\n")
	if actual, err := (Day12{}).Part1(input); err != nil || actual != "1" {
		t.Fatal(actual, err)
	}
}

func TestDay12OrientationsAreDistinctAndPreserveCells(t *testing.T) {
	for mask := 1; mask < 64; mask++ {
		var cells []day12Point
		for i := 0; i < 6; i++ {
			if mask&(1<<i) != 0 {
				cells = append(cells, day12Point{i % 3, i / 3})
			}
		}
		shape := packingShape(cells)
		want := packingOraclePlacements(cells, 5, 5)
		actual := make(map[uint64]bool)
		for _, orientation := range shape.orientations {
			if len(orientation) != bits.OnesCount(uint(mask)) {
				t.Fatal("orientation changed cell count")
			}
			width, height := day12Bounds(orientation)
			for y := 0; y <= 5-height; y++ {
				for x := 0; x <= 5-width; x++ {
					var placement uint64
					for _, p := range orientation {
						placement |= uint64(1) << uint((p.y+y)*5+p.x+x)
					}
					actual[placement] = true
				}
			}
		}
		if len(actual) != len(want) {
			t.Fatal("orientation placement count", mask, len(actual), len(want))
		}
		for _, placement := range want {
			if !actual[placement] {
				t.Fatal("missing reflected or rotated placement", mask)
			}
		}
	}
}

func TestDay12ParsingAndArithmeticErrors(t *testing.T) {
	invalid := []string{
		"", "0:\n...\n\n2x2: 0\n", "0:\n##\n#\n\n2x2: 1\n",
		"0:\n#x\n\n2x2: 1\n", "0:\n#\n\n0:\n#\n\n2x2: 1\n",
		"99999999:\n#\n\n2x2: 1\n", "0:\n#\n\n2x2: -1\n", "0:\n#\n\n0x2: 1\n",
		"0:\n#\n\n2x2: 1 2\n", "0:\n#\n\n2x2:\n", "0:\n#\n", "2x2: 0\n",
		"0:\n#\n\n2x2: 1\n1:\n#\n", "0:\n#\n\n2xQ: 1\n",
		fmt.Sprintf("0:\n#\n\n%dx2: 1\n", math.MaxInt),
		fmt.Sprintf("0:\n##\n\n2x2: %d\n", math.MaxInt),
		fmt.Sprintf("0:\n#\n\n1:\n#\n\n2x2: %d %d\n", math.MaxInt, math.MaxInt),
	}
	for _, input := range invalid {
		if _, err := (Day12{}).Part1(input); err == nil {
			t.Fatalf("accepted invalid input%q", input)
		}
	}
	if _, err := (Day12{}).Part2("anything"); err == nil {
		t.Fatal("Part2 unexpectedly implemented")
	}
}

func TestDay12ProofsAndExactSearchBudgets(t *testing.T) {
	bars := []day12Shape{packingShape([]day12Point{{0, 0}, {1, 0}, {2, 0}})}
	if fits, err := canFitRegion(day12Region{w: 1, h: 6000, counts: []int{2000}}, bars); err != nil || !fits {
		t.Fatal("rotated box certificate failed", fits, err)
	}
	shapes := []day12Shape{packingShape([]day12Point{{0, 0}, {0, 1}, {1, 1}})}
	if fits, err := canFitRegion(day12Region{w: 100000, h: 100000, counts: []int{1000}}, shapes); err != nil || !fits {
		t.Fatal("box certificate failed", fits, err)
	}
	if fits, err := canFitRegion(day12Region{w: 1, h: 1, counts: []int{1}}, shapes); err != nil || fits {
		t.Fatal("area rejection failed", fits, err)
	}
	if _, err := day12ExactFit(day12Region{w: 3, h: 2, counts: []int{2}}, shapes, 0); err == nil {
		t.Fatal("node budget exhaustion must be an error")
	}
	if _, err := day12ExactFit(day12Region{w: 100, h: 100, counts: []int{2}}, shapes, 10); err == nil {
		t.Fatal("board budget exhaustion must be an error")
	}
	cells := make([]day12Point, 0, 1024)
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			cells = append(cells, day12Point{x, y})
		}
	}
	large := []day12Shape{packingShape(cells)}
	if _, err := day12ExactFit(day12Region{w: 64, h: 64, counts: []int{1}}, large, 10); err == nil {
		t.Fatal("placement storage exhaustion must be an error")
	}
}
