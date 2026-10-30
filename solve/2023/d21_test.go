package solve2023

import (
	"errors"
	"math"
	"strings"
	"testing"
)

func testSmallGarden(width, height, start, rocks int) gardenMap {
	cells := make([]byte, width*height)
	for i := range cells {
		if i == start {
			cells[i] = 'S'
		} else {
			cells[i] = '.'
			if rocks&1 != 0 {
				cells[i] = '#'
			}
			rocks >>= 1
		}
	}
	return gardenMap{width, height, cells, gardenPoint{start % width, start / width}}
}

// Keep exact-step frontiers and wrap coordinates independently of production BFS
// and distance certification. This oracle is deliberately limited to small cases.
func referenceGardenCounts(g gardenMap, repeating bool, steps int) []int64 {
	frontier := map[gardenPoint]bool{g.start: true}
	counts := []int64{1}
	for step := 0; step < steps; step++ {
		next := map[gardenPoint]bool{}
		for p := range frontier {
			for _, d := range [4]gardenPoint{{0, -1}, {0, 1}, {-1, 0}, {1, 0}} {
				point := gardenPoint{p.x + d.x, p.y + d.y}
				x, y := point.x, point.y
				if repeating {
					x = ((x % g.width) + g.width) % g.width
					y = ((y % g.height) + g.height) % g.height
				} else if x < 0 || y < 0 || x >= g.width || y >= g.height {
					continue
				}
				if g.cells[y*g.width+x] != '#' {
					next[point] = true
				}
			}
		}
		frontier = next
		counts = append(counts, int64(len(frontier)))
	}
	return counts
}
func TestGardenParserAndBoundaries(t *testing.T) {
	for _, input := range []string{"", ".\n", "SS", "S?", "Sé", "S\n\n", "S.\n.", " S"} {
		if _, err := parseGarden(input); err == nil {
			t.Fatalf("accepted malformed input %q", input)
		}
	}
	g, err := parseGarden(".S\r\n#.\r\n")
	if err != nil || g.width != 2 || g.height != 2 || g.start != (gardenPoint{1, 0}) {
		t.Fatal("CRLF parse failed")
	}
	g, _ = parseGarden("S.#\n.#.")
	for _, tc := range []struct {
		point           gardenPoint
		repeating, want bool
	}{{gardenPoint{-1, 0}, false, false}, {gardenPoint{-1, 0}, true, false}, {gardenPoint{-1, -1}, true, true}, {gardenPoint{1, -1}, true, false}} {
		if g.open(tc.point, tc.repeating) != tc.want {
			t.Fatal("axis wrapping failure")
		}
	}
	if _, err := parseGarden("S" + strings.Repeat(".", gardenCellLimit)); !errors.Is(err, errGardenLimit) {
		t.Fatal("map size not bounded")
	}
	if _, err := parseGarden("S.\n.x"); err == nil || !strings.HasPrefix(err.Error(), "line 2:") {
		t.Fatal("missing line context")
	}
}
func TestGardenFiniteDistanceParity(t *testing.T) {
	cases := 0
	for width := 1; width <= 3; width++ {
		for height := 1; height <= 3; height++ {
			for rocks := 0; rocks < 1<<(width*height-1); rocks++ {
				g := testSmallGarden(width, height, 0, rocks)
				want := referenceGardenCounts(g, false, 8)
				for steps, count := range want {
					got, err := g.finiteCount(int64(steps))
					if err != nil || got != count {
						t.Fatalf("shape %dx%d rocks %d step %d: %d %v vs %d", width, height, rocks, steps, got, err, count)
					}
				}
				cases++
			}
		}
	}
	if cases != 341 {
		t.Fatal("incomplete rectangle corpus")
	}
	isolated, _ := parseGarden("###\n#S#\n###")
	if got, _ := isolated.finiteCount(2); got != 0 {
		t.Fatal("isolated start cannot spend spare steps")
	}
	if got, err := (Day21{}).Part1("S"); err != nil || got != "0" {
		t.Fatal("part one must not wrap boundaries")
	}
	if _, err := isolated.finiteCount(-1); err == nil {
		t.Fatal("negative steps accepted")
	}
}
func TestGardenCertificatesAgainstExactSteps(t *testing.T) {
	cases, accepted := 0, 0
	for width := 1; width <= 3; width++ {
		for start := 0; start < width*width; start++ {
			for rocks := 0; rocks < 1<<(width*width-1); rocks++ {
				g := testSmallGarden(width, width, start, rocks)
				c, err := g.certify()
				if err == nil {
					want := referenceGardenCounts(g, true, 20)
					for steps, count := range want {
						got, err := c.count(int64(steps))
						if err != nil || got != count {
							t.Fatalf("width %d start %d rocks %d step %d: %d %v vs %d", width, start, rocks, steps, got, err, count)
						}
					}
					accepted++
				}
				cases++
			}
		}
	}
	if cases != 2337 || accepted < 100 {
		t.Fatalf("incomplete certificate corpus: %d/%d", accepted, cases)
	}
}
func TestGardenRejectsCorruptCertificates(t *testing.T) {
	g, _ := parseGarden("S.\n..")
	model, _ := buildGardenModel(&g, 1)
	if valid, err := model.validates(&g); err != nil || !valid {
		t.Fatal("open-map model invalid")
	}
	index := 0
	for model.distances[index] != 1 {
		index++
	}
	saved := model.distances[index]
	for _, value := range []int32{0, -1, 4} {
		model.distances[index] = value
		if valid, _ := model.validates(&g); valid {
			t.Fatal("corrupted certificate accepted")
		}
	}
	model.distances[index] = saved
	for index, d := range model.distances {
		if d == 0 {
			model.distances[index] = 2
			break
		}
	}
	if valid, _ := model.validates(&g); valid {
		t.Fatal("missing zero accepted")
	}
}
func TestGardenTranslationArithmetic(t *testing.T) {
	for width := int64(1); width <= 6; width++ {
		for remaining := int64(-1); remaining <= 60; remaining++ {
			for axes := 0; axes <= 2; axes++ {
				want := int64(0)
				switch axes {
				case 0:
					if remaining >= 0 && remaining%2 == 0 {
						want = 1
					}
				case 1:
					for x := int64(1); x <= 60; x++ {
						if remaining >= x*width && (remaining-x*width)%2 == 0 {
							want++
						}
					}
				case 2:
					for x := int64(1); x <= 60; x++ {
						for y := int64(1); y <= 60; y++ {
							if remaining >= (x+y)*width && (remaining-(x+y)*width)%2 == 0 {
								want++
							}
						}
					}
				}
				got, err := gardenTranslations(remaining, width, axes)
				if err != nil || got != want {
					t.Fatalf("translation %d %d %d: %d %v vs %d", remaining, width, axes, got, err, want)
				}
			}
		}
	}
	g, _ := parseGarden("S")
	c, err := g.certify()
	if err != nil {
		t.Fatal(err)
	}
	steps := int64(100_000_000)
	if got, err := c.count(steps); err != nil || got != (steps+1)*(steps+1) {
		t.Fatal("large integer count lost precision")
	}
	if _, err := c.count(math.MaxInt64); !errors.Is(err, errGardenOverflow) {
		t.Fatal("count overflow not reported")
	}
	if _, err := c.count(-1); err == nil {
		t.Fatal("negative steps accepted")
	}
	if _, err := gardenTranslations(math.MaxInt64-1, 2, 2); !errors.Is(err, errGardenOverflow) {
		t.Fatal("translation overflow not reported")
	}
}
func TestGardenLargerCoreAndUnsupportedModels(t *testing.T) {
	// Independently authored geometry retained from the Rust review: radius one
	// fails its distance equations; radius two admits a checked certificate.
	g, _ := parseGarden("#...#\n....#\n#.S.#\n.###.\n...##")
	first, _ := buildGardenModel(&g, 1)
	if valid, _ := first.validates(&g); valid {
		t.Fatal("fixture must require larger core")
	}
	c, err := g.certify()
	if err != nil || c.model.radius != 2 {
		t.Fatalf("larger core failed: %v", err)
	}
	for steps, want := range referenceGardenCounts(g, true, 60) {
		if got, err := c.count(int64(steps)); err != nil || got != want {
			t.Fatalf("larger core step %d", steps)
		}
	}
	g, _ = parseGarden(".#..#\n.#...\n..S##\n.#.#.\n.....")
	if _, err := g.certify(); !errors.Is(err, errGardenModel) {
		t.Fatalf("unsupported geometry: %v", err)
	}
	g, _ = parseGarden("S.")
	if _, err := g.certify(); err == nil {
		t.Fatal("non-square model accepted")
	}
	g = testSmallGarden(700, 700, 0, 0)
	if _, err := g.certify(); !errors.Is(err, errGardenLimit) {
		t.Fatalf("patch validation budget ignored: %v", err)
	}
}
