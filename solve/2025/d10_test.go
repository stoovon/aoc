package solve2025

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func buttonCounterOracle(buttons []int, target [3]int) (int, bool) {
	type state struct {
		counters [3]int
		presses  int
	}
	queue := []state{{}}
	seen := map[[3]int]bool{{}: true}
	for head := 0; head < len(queue); head++ {
		cur := queue[head]
		if cur.counters == target {
			return cur.presses, true
		}
		for _, button := range buttons {
			next := cur.counters
			valid := true
			for i := range next {
				if button&(1<<i) != 0 {
					next[i]++
				}
				if next[i] > target[i] {
					valid = false
				}
			}
			if valid && !seen[next] {
				seen[next] = true
				queue = append(queue, state{next, cur.presses + 1})
			}
		}
	}
	return 0, false
}

func TestButtonJoltageBoundsAgainstCounterSearch(t *testing.T) {
	for subset := 0; subset < 128; subset++ {
		buttons := []int{}
		records := []string{"[...]"}
		for mask := 1; mask < 8; mask++ {
			if subset&(1<<(mask-1)) != 0 {
				buttons = append(buttons, mask)
				indices := []string{}
				for i := 0; i < 3; i++ {
					if mask&(1<<i) != 0 {
						indices = append(indices, strconv.Itoa(i))
					}
				}
				records = append(records, "("+strings.Join(indices, ",")+")")
			}
		}
		for code := 0; code < 27; code++ {
			target := [3]int{code % 3, (code / 3) % 3, code / 9}
			line := strings.Join(records, " ") + fmt.Sprintf(" {%d,%d,%d}", target[0], target[1], target[2])
			want, exists := buttonCounterOracle(buttons, target)
			got, err := minPressesForMachineJolts(line)
			if exists && (err != nil || got != int64(want)) {
				t.Fatalf("subset=%d target=%v: got %d %v want %d", subset, target, got, err, want)
			}
			if !exists && err == nil {
				t.Fatalf("subset=%d target=%v: accepted impossible target", subset, target)
			}
		}
	}
}

func TestButtonJoltageUnequalBounds(t *testing.T) {
	for _, tc := range []struct {
		line string
		want int64
	}{
		{"[..] (0) (1) (0,1) {12,20}", 20},
		{"[...] (0) (1) (2) (0,1) (0,2) (1,2) (0,1,2) {0,0,16}", 16},
		{"[...] (0,1) (0,2) (1,2) {3,3,4}", 5},
		{"[..] (0) (1) (0,1) {0,0}", 0},
	} {
		got, err := minPressesForMachineJolts(tc.line)
		if err != nil || got != tc.want {
			t.Fatalf("%s: got %d %v want %d", tc.line, got, err, tc.want)
		}
	}
}

func TestButtonJoltageExactLargeFractions(t *testing.T) {
	for _, tc := range []struct {
		line string
		want int64
	}{
		{"[...] (0,1) (0,2) (1,2) {6000000000000000000,6000000000000000000,6000000000000000000}", 9000000000000000000},
		{"[...] (0,1) (0,2) (1,2) {9223372036854775806,9223372036854775806,0}", 9223372036854775806},
		{"[.] (0) {9223372036854775807}", 9223372036854775807},
	} {
		got, err := minPressesForMachineJolts(tc.line)
		if err != nil || got != tc.want {
			t.Fatalf("large exact arithmetic: got %d %v want %d", got, err, tc.want)
		}
	}
	for _, line := range []string{
		"[...] (0,1) (0,2) (1,2) {6000000000000000001,6000000000000000001,6000000000000000001}",
		"[...] (0,1) (0,2) (1,2) {9223372036854775806,9223372036854775806,9223372036854775806}",
	} {
		if _, err := minPressesForMachineJolts(line); err == nil {
			t.Fatal("expected fractional or oversized minimum error")
		}
	}
	line := "[.] (0) {9223372036854775807}"
	if _, err := (Day10{}).Part2(line + "\n" + line); err == nil {
		t.Fatal("expected machine-total overflow error")
	}
}

func TestButtonJoltageSignaturesAcrossByteBoundaries(t *testing.T) {
	for _, width := range []int{17, 65, 129} {
		targets := make([]string, width)
		for i := range targets {
			targets[i] = "0"
		}
		targets[0] = "2"
		targets[width-1] = "1"
		line := fmt.Sprintf("[%s] (0) (0,%d) (%d,0) {%s}", strings.Repeat(".", width), width-1, width-1, strings.Join(targets, ","))
		got, err := minPressesForMachineJolts(line)
		if err != nil || got != 2 {
			t.Fatalf("width %d distinct and reordered signatures: got %d %v want 2", width, got, err)
		}
		targets[0] = "0"
		line = fmt.Sprintf("[%s] (%d) {%s}", strings.Repeat(".", width), width-1, strings.Join(targets, ","))
		got, err = minPressesForMachineJolts(line)
		if err != nil || got != 1 {
			t.Fatalf("width %d high bit alone: got %d %v want 1", width, got, err)
		}
	}
}

func TestButtonRepeatedIndicesAreRejected(t *testing.T) {
	for part, solver := range []func(string) (string, error){(Day10{}).Part1, (Day10{}).Part2} {
		for _, line := range []string{"[.] (0,0) {1}", "[##] (1,0,1) {1,1}"} {
			if _, err := solver(line); err == nil || !strings.Contains(err.Error(), "repeated button index") {
				t.Fatalf("part %d: expected repeated index error, got %v", part+1, err)
			}
		}
		got, err := solver("[#] (0) (0) {1}")
		if err != nil || got != "1" {
			t.Fatalf("part %d identical separate buttons: %s %v", part+1, got, err)
		}
	}
}

func TestButtonLightsIgnoreDuplicateWiring(t *testing.T) {
	line := "[#] " + strings.Repeat("(0) ", 32) + "{0}"
	if got, err := (Day10{}).Part1(line); err != nil || got != "1" {
		t.Fatalf("32 identical buttons: %q %v", got, err)
	}
	for subset := 0; subset < 128; subset++ {
		records := []string{}
		buttons := []int{}
		for mask := 1; mask < 8; mask++ {
			if subset&(1<<(mask-1)) != 0 {
				buttons = append(buttons, mask)
				forward, backward := []string{}, []string{}
				for bit := 0; bit < 3; bit++ {
					if mask&(1<<bit) != 0 {
						forward = append(forward, strconv.Itoa(bit))
					}
				}
				for bit := 2; bit >= 0; bit-- {
					if mask&(1<<bit) != 0 {
						backward = append(backward, strconv.Itoa(bit))
					}
				}
				records = append(records, "("+strings.Join(forward, ",")+")", "("+strings.Join(backward, ",")+")")
			}
		}
		distances := [8]int{0, -1, -1, -1, -1, -1, -1, -1}
		queue := []int{0}
		for head := 0; head < len(queue); head++ {
			for _, button := range buttons {
				next := queue[head] ^ button
				if distances[next] < 0 {
					distances[next] = distances[queue[head]] + 1
					queue = append(queue, next)
				}
			}
		}
		exhaustive := [8]int{0, -1, -1, -1, -1, -1, -1, -1}
		for choice := 0; choice < 1<<len(buttons); choice++ {
			state, presses := 0, 0
			for i, button := range buttons {
				if choice&(1<<i) != 0 {
					state ^= button
					presses++
				}
			}
			if exhaustive[state] < 0 || presses < exhaustive[state] {
				exhaustive[state] = presses
			}
		}
		for target := 0; target < 8; target++ {
			pattern := []byte("...")
			for bit := 0; bit < 3; bit++ {
				if target&(1<<bit) != 0 {
					pattern[bit] = '#'
				}
			}
			line := "[" + string(pattern) + "] " + strings.Join(records, " ") + " {0,0,0}"
			got, err := minPressesForMachine(line)
			want := exhaustive[target]
			if distances[target] != want {
				t.Fatal("independent references disagree")
			}
			if want < 0 {
				if err == nil {
					t.Fatalf("subset %d target %d: accepted unreachable lights", subset, target)
				}
			} else if err != nil || got != want {
				t.Fatalf("subset %d target %d: got %d %v want %d", subset, target, got, err, want)
			}
		}
	}
}

func TestButtonLightsAdaptiveSearch(t *testing.T) {
	for _, includeTarget := range []bool{true, false} {
		records := []string{}
		limit := 63
		if !includeTarget {
			limit = 62
		}
		for mask := 1; mask <= limit; mask++ {
			indices := []string{}
			for bit := 0; bit < 6; bit++ {
				if mask&(1<<bit) != 0 {
					indices = append(indices, strconv.Itoa(bit))
				}
			}
			records = append(records, "("+strings.Join(indices, ",")+")")
		}
		line := "[######] " + strings.Join(records, " ") + " {0,0,0,0,0,0}"
		want := 1
		if !includeTarget {
			want = 2
		}
		got, err := minPressesForMachine(line)
		if err != nil || got != want {
			t.Fatalf("dense six lights: %d %v want %d", got, err, want)
		}
	}
	records := []string{}
	for bit := 0; bit < 129; bit++ {
		records = append(records, fmt.Sprintf("(%d)", bit))
	}
	line := "[" + strings.Repeat("#", 129) + "] " + strings.Join(records, " ") + " {" + strings.Repeat("0,", 128) + "0}"
	if got, err := minPressesForMachine(line); err != nil || got != 129 {
		t.Fatalf("wide independent lights: %d %v", got, err)
	}
	target := make([]uint8, 129)
	target[128] = 1
	if got, err := minimumLightPatternPresses(target, [][]int{{0}, {64}, {128}}); err != nil || got != 1 {
		t.Fatalf("wide byte BFS: %d %v", got, err)
	}
}

func TestMachineParserValidation(t *testing.T) {
	invalid := []string{
		"garbage[.] (0) {1}", "[] {0}", "[x] (0) {1}", "[é] (0) {1}",
		"{0}[.] (0)", "[.] (0", "[.] () {0}", "[.] (0,) {0}", "[.] (,0) {0}",
		"[.] (0,,0) {0}", "[.] (-1) {0}", "[.] (+0) {0}", "[.] (1) {0}",
		"[.] (999999999999999999999999) {0}", "[.] garbage (0) {0}",
		"[.] (0) {}", "[.] (0) {0,}", "[.] (0) {,0}", "[.] (0) {0,0}",
		"[.] (0) {-1}", "[.] (0) {+1}", "[.] (0) {9223372036854775808}",
		"[.] (0) {0} trailing", "[.] (0) {0} {0}", "[.] {0} (0)", "[.] (0) {0",
	}
	for _, line := range invalid {
		for part, solver := range []func(string) (string, error){(Day10{}).Part1, (Day10{}).Part2} {
			if _, err := solver(line); err == nil {
				t.Fatalf("part %d accepted %q", part+1, line)
			}
		}
	}
	for part, solver := range []func(string) (string, error){(Day10{}).Part1, (Day10{}).Part2} {
		valid := "  [#.] ( 0 ) { 1, 0 } \r\n"
		if got, err := solver(valid); err != nil || got != "1" {
			t.Fatalf("part %d whitespace: %q %v", part+1, got, err)
		}
		if got, err := solver(""); err != nil || got != "0" {
			t.Fatalf("part %d empty: %q %v", part+1, got, err)
		}
		if got, err := solver("[..] {0,0}"); err != nil || got != "0" {
			t.Fatalf("part %d no buttons: %q %v", part+1, got, err)
		}
		if _, err := solver("[.] {0}\n\n[.] {0}"); err == nil || !strings.Contains(err.Error(), "line 2:") {
			t.Fatalf("part %d blank line: %v", part+1, err)
		}
	}
	if got, err := (Day10{}).Part1("[#] (0)"); err != nil || got != "1" {
		t.Fatalf("part1 optional joltage: %q %v", got, err)
	}
	if _, err := (Day10{}).Part2("[#] (0)"); err == nil {
		t.Fatal("part2 requires joltage")
	}
}

func FuzzMachineParser(f *testing.F) {
	for _, seed := range []string{"[.] (0) {1}", "{0}[.]", "[.] (0,0) {0}", "[##] (0,1) {2,2}", "", "[.] () {}"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, line string) {
		pattern, buttons, targets, err := parseMachineLine(line)
		if err != nil {
			return
		}
		if pattern == "" {
			t.Fatal("accepted empty pattern")
		}
		for _, button := range buttons {
			for _, index := range button {
				if index < 0 || index >= len(pattern) {
					t.Fatal("accepted unsafe button index")
				}
			}
		}
		if targets != nil && len(targets) != len(pattern) {
			t.Fatal("accepted mismatched targets")
		}
	})
}
