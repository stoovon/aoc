package solve2023

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func referenceRockTilt(cells []byte, width, height, direction int) {
	dr, dc := []int{-1, 0, 1, 0}[direction], []int{0, -1, 0, 1}[direction]
	for {
		changed := false
		for row := 0; row < height; row++ {
			for col := 0; col < width; col++ {
				r, c := row+dr, col+dc
				if r < 0 || c < 0 || r >= height || c >= width {
					continue
				}
				a, b := row*width+col, r*width+c
				if cells[a] == 'O' && cells[b] == '.' {
					cells[a], cells[b] = '.', 'O'
					changed = true
				}
			}
		}
		if !changed {
			return
		}
	}
}
func referenceRockSpin(cells []byte, width, height int) {
	for direction := 0; direction < 4; direction++ {
		referenceRockTilt(cells, width, height, direction)
	}
}
func TestRockTiltsAgainstCellMovement(t *testing.T) {
	cases := 0
	for _, shape := range [][2]int{{1, 1}, {1, 3}, {3, 1}, {2, 3}, {3, 2}, {3, 3}} {
		width, height := shape[0], shape[1]
		count := 1
		for i := 0; i < width*height; i++ {
			count *= 3
		}
		for code := 0; code < count; code++ {
			cells := make([]byte, width*height)
			n := code
			for i := range cells {
				cells[i] = []byte(".#O")[n%3]
				n /= 3
			}
			for direction := 0; direction < 4; direction++ {
				p := rockPlatform{width, height, bytes.Clone(cells)}
				want := bytes.Clone(cells)
				referenceRockTilt(want, width, height, direction)
				p.tilt(direction)
				if !bytes.Equal(want, p.cells) {
					t.Fatalf("shape %v code %d direction %d", shape, code, direction)
				}
				cases++
			}
		}
	}
	if cases < 80000 {
		t.Fatal("incomplete tilt corpus")
	}
}
func TestRockCycleSkippingAgainstDirectSimulation(t *testing.T) {
	for code := 0; code < 729; code++ {
		cells := make([]byte, 6)
		n := code
		for i := range cells {
			cells[i] = []byte(".#O")[n%3]
			n /= 3
		}
		direct := bytes.Clone(cells)
		for cycles := 0; cycles <= 20; cycles++ {
			p := rockPlatform{3, 2, bytes.Clone(cells)}
			if err := p.cycles(int64(cycles), defaultRockLimits); err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(p.cells, direct) {
				t.Fatalf("code %d cycles %d", code, cycles)
			}
			referenceRockSpin(direct, 3, 2)
		}
	}
	p, _ := parseRockPlatform("O")
	if err := p.cycles(1_000_000_000, defaultRockLimits); err != nil || string(p.cells) != "O" {
		t.Fatal("fixed-point skipping failed")
	}
}
func TestRockParserAndLoad(t *testing.T) {
	for _, input := range []string{"", ".", "###\n...", "O.\n.", "\nO", "O\n\nO", "O\n\n", "Ox", "O ", "Oé"} {
		for _, part := range []func(string) (string, error){(Day14{}).Part1, (Day14{}).Part2} {
			if _, err := part(input); err == nil {
				t.Fatalf("accepted %q", input)
			}
		}
	}
	got, err := (Day14{}).Part1("O.\r\n.O\r\n")
	if err != nil || got != "4" {
		t.Fatalf("CRLF/load: %s %v", got, err)
	}
	if _, err := parseRockPlatform("O" + strings.Repeat(".", rockCellLimit)); !errors.Is(err, errRockLimit) {
		t.Fatal("cell limit ignored")
	}
}
func TestRockResourceLimits(t *testing.T) {
	for _, limits := range []rockLimits{{0, 100, 1000}, {10, 3, 1000}, {10, 100, 0}, {1, 100, 1000}, {10, 4, 1000}} {
		p, _ := parseRockPlatform("O.\n..")
		if err := p.cycles(10, limits); !errors.Is(err, errRockLimit) {
			t.Fatalf("missing explicit limit error: %v", err)
		}
	}
	p, _ := parseRockPlatform("O")
	if err := p.cycles(0, rockLimits{}); err != nil {
		t.Fatal("zero cycles requires no work")
	}
	if err := p.cycles(-1, defaultRockLimits); err == nil {
		t.Fatal("negative target accepted")
	}
}
