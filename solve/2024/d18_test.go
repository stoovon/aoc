package solve2024

import (
	"strings"
	"testing"
)

func TestRAMPart1UsesExactly1024Bytes(t *testing.T) {
	input := strings.Repeat("1,0\n", 1024) + "0,1\n"
	got, err := (Day18{}).Part1(input)
	if err != nil || got != "140" {
		t.Fatalf("first 1024 bytes: got %q %v, want 140", got, err)
	}
}

func ramRelaxationReference(width int, falls []int, prefix int) (int, bool) {
	size := width * width
	blocked := make([]bool, size)
	for _, id := range falls[:prefix] {
		blocked[id] = true
	}
	if blocked[0] || blocked[size-1] {
		return 0, false
	}
	distances := make([]int, size)
	for i := range distances {
		distances[i] = size + 1
	}
	distances[0] = 0
	for pass := 0; pass < size; pass++ {
		for id := 0; id < size; id++ {
			if blocked[id] {
				continue
			}
			for other := 0; other < size; other++ {
				dx, dy := id%width-other%width, id/width-other/width
				if dx < 0 {
					dx = -dx
				}
				if dy < 0 {
					dy = -dy
				}
				if !blocked[other] && dx+dy == 1 && distances[other]+1 < distances[id] {
					distances[id] = distances[other] + 1
				}
			}
		}
	}
	if distances[size-1] > size {
		return 0, false
	}
	return distances[size-1], true
}

func checkRAMPrefixes(t *testing.T, width int, falls []int) {
	t.Helper()
	memory := ramMemory{width: width, falls: falls, blockedAt: make([]int, width*width)}
	for i, id := range falls {
		if memory.blockedAt[id] == 0 {
			memory.blockedAt[id] = i + 1
		}
	}
	first := 0
	for prefix := 0; prefix <= len(falls); prefix++ {
		want, reachable := ramRelaxationReference(width, falls, prefix)
		got, ok := memory.shortestPath(prefix)
		if ok != reachable || ok && got != want {
			t.Fatalf("width %d falls %v prefix %d: got %d %t want %d %t", width, falls, prefix, got, ok, want, reachable)
		}
		if !reachable && first == 0 {
			first = prefix
		}
	}
	got, err := memory.firstBlockingPrefix()
	if first == 0 {
		if err == nil {
			t.Fatal("expected no-blocking-byte error")
		}
	} else if err != nil || got != first {
		t.Fatalf("falls %v: first blocking prefix %d %v want %d", falls, got, err, first)
	}
}

func TestRAMPrefixedSearchAgainstRelaxation(t *testing.T) {
	combinations := 1
	for length := 0; length <= 5; length++ {
		if length > 0 {
			combinations *= 4
		}
		for code := 0; code < combinations; code++ {
			falls := make([]int, length)
			rest := code
			for i := range falls {
				falls[i] = rest % 4
				rest /= 4
			}
			checkRAMPrefixes(t, 2, falls)
		}
	}
	for mask := 0; mask < 512; mask++ {
		falls := []int{}
		for id := 0; id < 9; id++ {
			if mask&(1<<id) != 0 {
				falls = append(falls, id)
			}
		}
		checkRAMPrefixes(t, 3, falls)
	}
	checkRAMPrefixes(t, 1, nil)
	checkRAMPrefixes(t, 1, []int{0, 0})
}

func TestRAMParserAndFailureCases(t *testing.T) {
	for _, input := range []string{"", "1,2,3", "1,2junk", "-1,2", "+1,2", "71,0", "0,71", "1,", "1 2", "0,0\n\n0,1", "999999999999999999999,0"} {
		for _, solver := range []func(string) (string, error){(Day18{}).Part1, (Day18{}).Part2} {
			if _, err := solver(input); err == nil {
				t.Fatalf("accepted invalid input %q", input)
			}
		}
	}
	for _, input := range []string{"0,0", "70,70"} {
		if got, err := (Day18{}).Part2(input); err != nil || got != input {
			t.Fatalf("blocked endpoint: %q %v", got, err)
		}
	}
	if _, err := (Day18{}).Part1("1,0"); err == nil {
		t.Fatal("expected short-input error")
	}
	if _, err := (Day18{}).Part1(strings.Repeat("0,0\n", 1024)); err == nil {
		t.Fatal("expected unreachable error")
	}
	if _, err := (Day18{}).Part2("1,0"); err == nil {
		t.Fatal("expected never-blocked error")
	}
	memory, err := parseRAM(" 1 , 0 \r\n1,0\r\n0,1\r\n", 2)
	if err != nil {
		t.Fatal(err)
	}
	if first, err := memory.firstBlockingPrefix(); err != nil || first != 3 {
		t.Fatalf("duplicate byte chronology: %d %v", first, err)
	}
	if _, err := parseRAM("0,0\n0,2", 2); err == nil || !strings.Contains(err.Error(), "line 2:") {
		t.Fatalf("line context: %v", err)
	}
}
