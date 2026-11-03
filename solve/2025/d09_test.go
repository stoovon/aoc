package solve2025

import (
	"fmt"
	"math"
	"math/big"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func theaterInput(points []theaterPoint) string {
	var input strings.Builder
	for _, p := range points {
		fmt.Fprintf(&input, "%d,%d\n", p.x, p.y)
	}
	return input.String()
}

func theaterTileOracle(points []theaterPoint, x, y int64) bool {
	inside := false
	for i, a := range points {
		b := points[(i+1)%len(points)]
		if a.x == b.x && x == a.x && y >= min(a.y, b.y) && y <= max(a.y, b.y) {
			return true
		}
		if a.y == b.y && y == a.y && x >= min(a.x, b.x) && x <= max(a.x, b.x) {
			return true
		}
		if a.y == b.y && x >= min(a.x, b.x) && x < max(a.x, b.x) && a.y > y {
			inside = !inside
		}
	}
	return inside
}

func theaterRectangleOracle(points []theaterPoint, a, b theaterPoint) bool {
	for x := min(a.x, b.x); x <= max(a.x, b.x); x++ {
		for y := min(a.y, b.y); y <= max(a.y, b.y); y++ {
			if !theaterTileOracle(points, x, y) {
				return false
			}
		}
	}
	return true
}

func checkTheaterOracle(t *testing.T, points []theaterPoint) {
	t.Helper()
	mask, err := buildTheaterMask(points)
	if err != nil {
		t.Fatal(err, points)
	}
	var best1, best2 int64
	for i, a := range points {
		for _, b := range points[i+1:] {
			area := (max(a.x, b.x) - min(a.x, b.x) + 1) * (max(a.y, b.y) - min(a.y, b.y) + 1)
			best1 = max(best1, area)
			valid := theaterRectangleOracle(points, a, b)
			if mask.contains(a, b) != valid {
				t.Fatalf("containment disagreement: %v %v in %v", a, b, points)
			}
			if valid {
				best2 = max(best2, area)
			}
		}
	}
	input := theaterInput(points)
	for i, solver := range []func(string) (string, error){(Day9{}).Part1, (Day9{}).Part2} {
		expected := []int64{best1, best2}[i]
		actual, err := solver(input)
		if err != nil || actual != strconv.FormatInt(expected, 10) {
			t.Fatalf("part%d got %s,%v want%d polygon%v", i+1, actual, err, expected, points)
		}
	}
	xvalues, yvalues := make([]int64, len(points)), make([]int64, len(points))
	for i, p := range points {
		xvalues[i], yvalues[i] = p.x, p.y
	}
	xs, _ := theaterAxis(xvalues)
	ys, _ := theaterAxis(yvalues)
	for row, y := range ys {
		for col, x := range xs {
			index := (row+1)*mask.stride + col + 1
			outside := mask.prefix[index] - mask.prefix[index-1] - mask.prefix[index-mask.stride] + mask.prefix[index-mask.stride-1]
			if (outside == 0) != theaterTileOracle(points, x, y) {
				t.Fatalf("tile disagreement at%d,%d in%v", x, y, points)
			}
		}
	}
}

func TestDay9BoundaryAndConcavity(t *testing.T) {
	shapes := [][]theaterPoint{
		{{0, 0}, {4, 0}, {4, 2}, {0, 2}},
		{{0, 0}, {8, 0}, {8, 8}, {6, 8}, {6, 2}, {2, 2}, {2, 8}, {0, 8}},
		{{0, 0}, {7, 0}, {7, 7}, {5, 7}, {5, 1}, {2, 1}, {2, 7}, {0, 7}},
		{{0, 0}, {7, 0}, {7, 6}, {4, 6}, {4, 3}, {3, 3}, {3, 6}, {0, 6}},
		{{0, 0}, {5, 0}, {5, 1}, {0, 1}},
		{{0, 0}, {2, 0}, {5, 0}, {5, 3}, {0, 3}},
	}
	for _, shape := range shapes {
		checkTheaterOracle(t, shape)
		reversed := slices.Clone(shape)
		slices.Reverse(reversed)
		checkTheaterOracle(t, reversed)
	}
	mask, err := buildTheaterMask(shapes[1])
	if err != nil {
		t.Fatal(err)
	}
	if mask.contains(theaterPoint{0, 8}, theaterPoint{8, 8}) {
		t.Fatal("rectangle crossed the exterior gap between two row intervals")
	}
	actual, err := (Day9{}).Part2(theaterInput(shapes[0]))
	if err != nil || actual != "15" {
		t.Fatalf("boundary regression: %s %v", actual, err)
	}
}

func TestDay9GeneratedHistograms(t *testing.T) {
	for encoded := 0; encoded < 81; encoded++ {
		heights := [4]int64{}
		value := encoded
		for i := range heights {
			heights[i] = int64(value%3 + 1)
			value /= 3
		}
		points := []theaterPoint{{0, 0}, {8, 0}, {8, heights[3]}}
		for i := 3; i > 0; i-- {
			if heights[i] != heights[i-1] {
				points = append(points, theaterPoint{int64(i * 2), heights[i]}, theaterPoint{int64(i * 2), heights[i-1]})
			}
		}
		points = append(points, theaterPoint{0, heights[0]})
		for rotation := 0; rotation < 4; rotation++ {
			transformed := make([]theaterPoint, len(points))
			for i, p := range points {
				for k := 0; k < rotation; k++ {
					p.x, p.y = -p.y, p.x
				}
				transformed[i] = theaterPoint{p.x - 5, p.y - 4}
			}
			checkTheaterOracle(t, transformed)
		}
	}
}

func TestDay9CoordinatesAndPolygonValidation(t *testing.T) {
	day := Day9{}
	for _, solver := range []func(string) (string, error){day.Part1, day.Part2} {
		for _, input := range []string{"", " \r\n", "1,2", "1,2,3\n2,3", "x,2\n2,3", "9223372036854775808,0\n0,0"} {
			if _, err := solver(input); err == nil {
				t.Fatalf("accepted malformed input%q", input)
			}
		}
		_, err := solver("0,0\nnot-a-point\n1,1")
		if err == nil || !strings.Contains(err.Error(), "line 2:") {
			t.Fatalf("missing line context:%v", err)
		}
	}
	input := " \r\n -4, -3\r\n+1, -3\r\n 1, 0\r\n-4,0\r\n\r\n"
	for _, solver := range []func(string) (string, error){day.Part1, day.Part2} {
		actual, err := solver(input)
		if err != nil || actual != "24" {
			t.Fatalf("signed/whitespace: %s %v", actual, err)
		}
	}
	invalid := [][]theaterPoint{
		{{0, 0}, {3, 0}, {3, 2}},
		{{0, 0}, {3, 0}, {4, 2}, {0, 2}},
		{{0, 0}, {3, 0}, {3, 2}, {0, 2}, {0, 0}},
		{{0, 0}, {4, 0}, {2, 0}, {2, 3}, {0, 3}},
		{{0, 0}, {4, 0}, {4, 4}, {2, 4}, {2, -1}, {0, -1}},
	}
	for _, points := range invalid {
		if _, err := day.Part2(theaterInput(points)); err == nil {
			t.Fatal("accepted invalid polygon", points)
		}
	}
	if actual, err := day.Part1("0,0\n3,4"); err != nil || actual != "20" {
		t.Fatal("Part1 incorrectly requires a polygon", actual, err)
	}
}

func TestDay9AreaMatchesBigIntegerOracle(t *testing.T) {
	values := []int64{math.MinInt64, math.MinInt64 + 1, -4, -1, 0, 1, 7, 3037000499, math.MaxInt64 - 1, math.MaxInt64}
	for _, x1 := range values {
		for _, x2 := range values {
			for _, y1 := range values {
				for _, y2 := range values {
					a, b := theaterPoint{x1, y1}, theaterPoint{x2, y2}
					width := new(big.Int).Sub(big.NewInt(max(x1, x2)), big.NewInt(min(x1, x2)))
					height := new(big.Int).Sub(big.NewInt(max(y1, y2)), big.NewInt(min(y1, y2)))
					width.Add(width, big.NewInt(1))
					height.Add(height, big.NewInt(1))
					expected := new(big.Int).Mul(width, height)
					actual, err := theaterArea(a, b)
					if !expected.IsInt64() {
						if err == nil {
							t.Fatal("overflow accepted", a, b)
						}
					} else if err != nil || actual != expected.Int64() {
						t.Fatal("area disagreement", a, b, actual, err, expected)
					}
				}
			}
		}
	}
	for _, solver := range []func(string) (string, error){(Day9{}).Part1, (Day9{}).Part2} {
		points := []theaterPoint{{0, 0}, {math.MaxInt64, 0}, {math.MaxInt64, 2}, {0, 2}}
		if _, err := solver(theaterInput(points)); err == nil {
			t.Fatal("entrypoint failed to propagate area overflow")
		}
	}
}

func TestDay9CompressionIgnoresCoordinateMagnitude(t *testing.T) {
	for _, origin := range []int64{math.MinInt64, 0, math.MaxInt64 - 8} {
		points := []theaterPoint{{origin, origin}, {origin + 8, origin}, {origin + 8, origin + 2}, {origin, origin + 2}}
		actual, err := (Day9{}).Part2(theaterInput(points))
		if err != nil || actual != "27" {
			t.Fatal("large translation", actual, err)
		}
	}
	points := []theaterPoint{{-1000000000, -1000000000}, {1000000000, -1000000000}, {1000000000, 1000000000}, {-1000000000, 1000000000}}
	mask, err := buildTheaterMask(points)
	if err != nil {
		t.Fatal(err)
	}
	if len(mask.prefix) != 16 {
		t.Fatal("grid grew with coordinate range", len(mask.prefix))
	}
	actual, err := (Day9{}).Part2(theaterInput(points))
	if err != nil || actual != "4000000004000000001" {
		t.Fatal(actual, err)
	}
	for _, size := range [][2]int{{math.MaxInt, 1}, {1, math.MaxInt}, {4096, 4096}, {0, 1}} {
		if _, err := theaterMaskSize(size[0], size[1]); err == nil {
			t.Fatal("unbounded grid accepted", size)
		}
	}
}
