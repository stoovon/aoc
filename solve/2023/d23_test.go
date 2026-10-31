package solve2023

import (
	"strconv"
	"strings"
	"testing"
)

func hikeTileReference(input string, slopes bool) (int, bool) {
	lines := strings.Split(input, "\n")
	height, width := len(lines), len(lines[0])
	start := strings.IndexAny(lines[0], ".^>v<")
	end := strings.IndexAny(lines[height-1], ".^>v<")
	seen := make([]bool, height*width)
	var walk func(int, int) (int, bool)
	walk = func(row, column int) (int, bool) {
		if row == height-1 && column == end {
			return 0, true
		}
		seen[row*width+column] = true
		arrow := -1
		switch lines[row][column] {
		case '^':
			arrow = 0
		case '>':
			arrow = 1
		case 'v':
			arrow = 2
		case '<':
			arrow = 3
		}
		best, exists := 0, false
		for direction, delta := range [][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}} {
			if slopes && arrow >= 0 && arrow != direction {
				continue
			}
			r, c := row+delta[0], column+delta[1]
			if r < 0 || r >= height || c < 0 || c >= width || lines[r][c] == '#' || seen[r*width+c] {
				continue
			}
			if distance, ok := walk(r, c); ok {
				best = max(best, distance+1)
				exists = true
			}
		}
		seen[row*width+column] = false
		return best, exists
	}
	return walk(0, start)
}
func TestHikeParallelCorridorsPerDirection(t *testing.T) {
	plain := "#.####\n#....#\n#..#.#\n####.#\n####.#"
	slopes := "#.####\n#....#\n#>^#.#\n####.#\n####.#"
	for _, input := range []string{plain, slopes} {
		for _, part := range []struct {
			slopes bool
			solve  func(string) (string, error)
		}{{true, (Day23{}).Part1}, {false, (Day23{}).Part2}} {
			want, ok := hikeTileReference(input, part.slopes)
			if !ok {
				t.Fatal("synthetic path missing")
			}
			got, err := part.solve(input)
			if err != nil || got != strconv.Itoa(want) {
				t.Fatal("lost longer corridor", got, want, err)
			}
		}
	}
	grid, err := newDay23Grid(slopes)
	if err != nil {
		t.Fatal(err)
	}
	directed := grid.intersectionsAndDistances(true)
	undirected := grid.intersectionsAndDistances(false)
	has := func(edges [][2]int, target, distance int) bool {
		for _, edge := range edges {
			if edge == [2]int{target, distance} {
				return true
			}
		}
		return false
	}
	if !has(directed[1], 2, 3) || !has(directed[2], 1, 1) || !has(undirected[2], 1, 3) {
		t.Fatal("directed corridor dominance")
	}
}
func TestHikeCompressionAgainstExhaustiveTileSearch(t *testing.T) {
	for mask := 0; mask < 1<<12; mask++ {
		for _, withArrows := range []bool{false, true} {
			cells := []byte(strings.Repeat("#", 30))
			cells[1] = '.'
			cells[28] = '.'
			for tile := 0; tile < 12; tile++ {
				if mask&(1<<tile) != 0 {
					value := byte('.')
					if withArrows {
						value = ".^>v<"[(tile+mask)%5]
					}
					cells[(tile/4+1)*6+tile%4+1] = value
				}
			}
			rows := []string{}
			for row := 0; row < 5; row++ {
				rows = append(rows, string(cells[row*6:(row+1)*6]))
			}
			input := strings.Join(rows, "\n")
			for _, part := range []struct {
				slopes bool
				solve  func(string) (string, error)
			}{{true, (Day23{}).Part1}, {false, (Day23{}).Part2}} {
				want, exists := hikeTileReference(input, part.slopes)

				got, err := part.solve(input)
				if (!exists && (err == nil || !strings.Contains(err.Error(), "no route"))) || (exists && (err != nil || got != strconv.Itoa(want))) {
					t.Fatalf("compression mismatch slopes=%t: got %s want %d err %v\n%s", part.slopes, got, want, err, input)
				}
			}
		}
	}
}

func TestHikeEndpointsAndValidation(t *testing.T) {
	for input, want := range map[string]string{"###.#\n#...#\n#.###": "4", "###.#\r\n#...#\r\n#.###\r\n": "4", ".": "0", ".\n.\n.": "2", "v\n.": "1"} {
		for _, part := range []func(string) (string, error){(Day23{}).Part1, (Day23{}).Part2} {
			got, err := part(input)
			if err != nil || got != want {
				t.Fatalf("%q: got %s, %v; want %s", input, got, err, want)
			}
		}
	}
	for _, input := range []string{"", "\n", "##\n##", "..\n#.", ".#\n..", ".#\n.", ".#\n\n.#", ".#\n?海", ".#\n #"} {
		if _, err := newDay23Grid(input); err == nil {
			t.Fatalf("accepted %q", input)
		}
	}
}

func TestHikeNoRouteAndZeroSteps(t *testing.T) {
	for _, part := range []func(string) (string, error){(Day23{}).Part1, (Day23{}).Part2} {
		if got, err := part("."); err != nil || got != "0" {
			t.Fatal(got, err)
		}
		if _, err := part(".#\n##\n#."); err == nil || !strings.Contains(err.Error(), "no route") {
			t.Fatal(err)
		}
	}
	if _, err := (Day23{}).Part1("^\n."); err == nil {
		t.Fatal("ignored entrance slope")
	}
	if got, err := (Day23{}).Part2("^\n."); err != nil || got != "1" {
		t.Fatal(got, err)
	}
}

func TestHikeForcedExitChecksIncomingCorridors(t *testing.T) {
	cases := []struct {
		edges         [][][2]int
		end, distance int
		exists        bool
	}{
		{[][][2]int{{{1, 2}, {2, 3}}, {{2, 11}, {3, 5}}, {{1, 7}}, {{0, 1}, {2, 1}}}, 3, 15, true},
		{[][][2]int{{{1, 2}, {2, 3}}, {{3, 5}}, {{3, 7}}, {}}, 3, 10, true},
		{[][][2]int{{{1, 2}}, {{0, 2}}, {{3, 5}}, {}}, 3, 0, false},
		{[][][2]int{{}}, 0, 0, true},
		{[][][2]int{{}, {}}, 1, 0, false},
		{[][][2]int{{{1, 5}}, {}}, 1, 5, true},
	}
	for _, c := range cases {
		got, exists := hikeLongestRoute(c.edges, 0, c.end)
		if exists != c.exists || (exists && got != c.distance) {
			t.Fatalf("got %d,%t want %d,%t", got, exists, c.distance, c.exists)
		}
		plain, ok := hikeLongestPath(0, c.edges, make([]bool, len(c.edges)), c.end)
		if exists != ok || (exists && got != plain) {
			t.Fatal("forced exit differs from exhaustive search")
		}
	}
}
