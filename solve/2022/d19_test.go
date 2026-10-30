package solve2022

import (
	"aoc/solve"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"
)

func robotTestRecord(id int64, costs [6]uint16) string {
	return fmt.Sprintf("Blueprint %d: Each ore robot costs %d ore. Each clay robot costs %d ore. Each obsidian robot costs %d ore and %d clay. Each geode robot costs %d ore and %d obsidian.", id, costs[0], costs[1], costs[2], costs[3], costs[4], costs[5])
}
func TestRobotParser(t *testing.T) {
	record := robotTestRecord(401, [6]uint16{2, 1, 1, 1, 3, 1})
	for _, valid := range []string{record, record + "\n", record + "\r\n", robotTestRecord(73, [6]uint16{}), robotTestRecord(79, [6]uint16{65535, 65535, 65535, 65535, 65535, 65535})} {
		bps, err := parseRobotBlueprints(valid)
		if err != nil || len(bps) != 1 {
			t.Fatal("valid blueprint rejected", err)
		}
	}
	for _, invalid := range []string{"", " \n\t", robotTestRecord(0, [6]uint16{}), robotTestRecord(-1, [6]uint16{}), strings.Replace(record, "401", "9223372036854775808", 1), strings.Replace(record, "2 ore.", "65536 ore.", 1), strings.Replace(record, "2 ore.", "-2 ore.", 1), strings.Replace(record, "2 ore.", "+2 ore.", 1), strings.Replace(record, "Each ore", "Every ore", 1), strings.TrimSuffix(record, "."), "1 2 3 4 5 6 7", record + "junk", " " + record, record + "\n\n", record + "\n" + record} {
		if _, err := parseRobotBlueprints(invalid); err == nil {
			t.Fatalf("accepted invalid record %q", invalid)
		}
	}
	for _, part := range []func(string) (string, error){(Day19{}).Part1, (Day19{}).Part2} {
		if _, err := part(record + "\nbad"); err == nil || !strings.HasPrefix(err.Error(), "line 2:") {
			t.Fatal("missing line context")
		}
	}
}
func TestRobotPartsAndArithmetic(t *testing.T) {
	first, second, third, fourth := robotTestRecord(401, [6]uint16{}), robotTestRecord(902, [6]uint16{}), robotTestRecord(1203, [6]uint16{}), robotTestRecord(1604, [6]uint16{65535, 65535, 65535, 65535, 65535, 65535})
	geodes24, geodes32 := int64(24*23/2), int64(32*31/2)
	for _, tc := range []struct {
		input, want string
		part        int
	}{{first, strconv.FormatInt(401*geodes24, 10), 1}, {first + "\n" + second, strconv.FormatInt((401+902)*geodes24, 10), 1}, {first, strconv.FormatInt(geodes32, 10), 2}, {first + "\n" + second, strconv.FormatInt(geodes32*geodes32, 10), 2}, {strings.Join([]string{first, second, third, fourth}, "\n"), strconv.FormatInt(geodes32*geodes32*geodes32, 10), 2}} {
		part := (Day19{}).Part1
		if tc.part == 2 {
			part = (Day19{}).Part2
		}
		got, err := part(tc.input)
		if err != nil || got != tc.want {
			t.Fatalf("part%d: %s %v vs %s", tc.part, got, err, tc.want)
		}
	}
	if _, err := (Day19{}).Part2(strings.Join([]string{first, second, third, "bad"}, "\n")); err == nil {
		t.Fatal("unselected records still require validation")
	}
	if _, err := (Day19{}).Part1(robotTestRecord(math.MaxInt64, [6]uint16{})); !errors.Is(err, errRobotOverflow) {
		t.Fatal("unchecked product")
	}
	id := int64(math.MaxInt64 / geodes24)
	a, b := robotTestRecord(id, [6]uint16{}), robotTestRecord(id-1, [6]uint16{})
	for _, input := range []string{a, b} {
		if _, err := (Day19{}).Part1(input); err != nil {
			t.Fatal("single weighted result should fit")
		}
	}
	if _, err := (Day19{}).Part1(a + "\n" + b); !errors.Is(err, errRobotOverflow) {
		t.Fatal("unchecked sum")
	}
	solver, ok := solve.GetAllSolvers()[solve.SolutionCoords{Year: 2022, Day: 19}]
	if !ok {
		t.Fatal("solver not registered")
	}
	if _, err := solver.Part1("bad"); err == nil {
		t.Fatal("registered solver swallowed error")
	}
}
func TestRobotInvestingCanBeatImmediateGeodes(t *testing.T) {
	bp, err := parseRobotBlueprint(robotTestRecord(71, [6]uint16{2, 1, 1, 1, 3, 1}))
	if err != nil {
		t.Fatal(err)
	}
	got, err := maxRobotGeodes(&bp, 14, robotSearchWork)
	if err != nil || got != 22 {
		t.Fatalf("investment regression: %d %v", got, err)
	}
}
func TestRobotBoundsAndLimits(t *testing.T) {
	bp, _ := parseRobotBlueprint(robotTestRecord(73, [6]uint16{}))
	got, err := maxRobotGeodes(&bp, 32, robotSearchWork)
	if err != nil || got != 32*31/2 {
		t.Fatal("geode count above 255 was truncated")
	}
	for _, budget := range []uint32{0, 2} {
		if partial, err := maxRobotGeodes(&bp, 32, budget); !errors.Is(err, errRobotSearchLimit) || partial != 0 {
			t.Fatal("budget exhaustion returned a partial result")
		}
	}
	if got, err := maxRobotGeodes(&bp, 0, 1); err != nil || got != 0 {
		t.Fatal("zero horizon")
	}
	if _, err := maxRobotGeodes(&bp, 33, robotSearchWork); err == nil {
		t.Fatal("unsupported horizon")
	}
	// Large costs are accepted without narrowing the wait arithmetic to a byte.
	costly, _ := parseRobotBlueprint(robotTestRecord(79, [6]uint16{65535, 65535, 65535, 65535, 65535, 65535}))
	if got, err := maxRobotGeodes(&costly, 32, robotSearchWork); err != nil || got != 0 {
		t.Fatal("large recipe overflow")
	}
	// Free clay production can exceed 255 before the unaffordable obsidian recipe.
	growing, _ := parseRobotBlueprint(robotTestRecord(83, [6]uint16{0, 0, 0, 500, 0, 1}))
	if got, err := maxRobotGeodes(&growing, 32, robotSearchWork); err != nil || got != 0 {
		t.Fatal("resource balance overflow")
	}
}

type robotReferenceState struct {
	left      uint8
	resources [3]uint16
	robots    [3]uint8
}

// This oracle advances exactly one minute and considers all five actions. It
// shares no production caps, wait jumps, optimistic bounds, or branch ordering.
func robotExhaustive(recipes [4][3]uint16, state robotReferenceState, memo map[robotReferenceState]uint16) uint16 {
	if state.left <= 1 {
		return 0
	}
	if value, ok := memo[state]; ok {
		return value
	}
	var best uint16
	for action := 0; action <= 4; action++ {
		affordable := true
		if action < 4 {
			for r, cost := range recipes[action] {
				if cost > state.resources[r] {
					affordable = false
				}
			}
		}
		if !affordable {
			continue
		}
		next := state
		next.left--
		for r := 0; r < 3; r++ {
			next.resources[r] += uint16(state.robots[r])
			if action < 4 {
				next.resources[r] -= recipes[action][r]
			}
		}
		if action < 3 {
			next.robots[action]++
		}
		gain := uint16(0)
		if action == 3 {
			gain = uint16(state.left - 1)
		}
		value := gain + robotExhaustive(recipes, next, memo)
		if value > best {
			best = value
		}
	}
	memo[state] = best
	return best
}
func TestRobotSearchAgainstExhaustiveBlueprints(t *testing.T) {
	cases := 0
	for encoding := 0; encoding < 729; encoding++ {
		n := encoding
		var costs [6]uint16
		for i := range costs {
			costs[i] = uint16(n % 3)
			n /= 3
		}
		bp, err := parseRobotBlueprint(robotTestRecord(79, costs))
		if err != nil {
			t.Fatal(err)
		}
		memo := map[robotReferenceState]uint16{}
		for minutes := uint8(0); minutes <= 8; minutes++ {
			want := robotExhaustive(bp.recipes, robotReferenceState{left: minutes, robots: [3]uint8{1, 0, 0}}, memo)
			got, err := maxRobotGeodes(&bp, minutes, robotSearchWork)
			if err != nil || got != want {
				t.Fatalf("costs %v minutes %d: %d %v vs %d", costs, minutes, got, err, want)
			}
			cases++
		}
	}
	if cases != 6561 {
		t.Fatal("incomplete oracle corpus")
	}
}
