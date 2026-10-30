package solve2023

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
)

func cascadeSize(bricks []Brick, removed int) int {
	return cascadeWithScratch(bricks, removed, make([]int, len(bricks)), make([]int, 0, len(bricks)))
}

func TestBrickParser(t *testing.T) {
	invalid := []string{"", " \n\t", "0,0,1", "0,0,1~0,0,1~0,0,2", "0,0~0,0,1", "0,0,1,2~0,0,1", "-1,0,1~0,0,1", "+0,0,1~0,0,1", "0, 0,1~0,0,1", "0,0,1~0,0,1\n\n", strconv.Itoa(math.MaxInt) + "0,0,1~0,0,1", "0,0,0~0,0,3", "0,0,3~0,0,0", "0,0,1~1,1,1", "0,0,1~1,0,2", "0,0,1~0,1,2", "0,0,1~1,1,2"}
	for _, input := range invalid {
		for _, solve := range []func(string) (string, error){(Day22{}).Part1, (Day22{}).Part2} {
			if _, err := solve(input); err == nil {
				t.Fatalf("accepted invalid %q", input)
			}
		}
	}
	for _, line := range []string{"0,0,2~2,0,2\n1,0,1~1,0,3", "0,0,2~2,0,2\n2,0,2~2,0,2", "0,0,2~0,0,2\n0,0,2~0,0,2", "0,0,1~0,0,5\n0,0,4~0,0,2"} {
		_, err := parseBricks(line)
		if err == nil || !strings.Contains(err.Error(), "line 2:") || !strings.Contains(err.Error(), "line 1") {
			t.Fatal("missing intersection context", err)
		}
	}
	_, err := parseBricks("0,0,1~0,0,1\nbad")
	if err == nil || !strings.Contains(err.Error(), "line 2:") {
		t.Fatal("missing parse line context")
	}
	for _, input := range []string{"0,0,1~0,0,1\r\n", "0,0,1~0,0,1\n1,0,1~1,0,1", "0,0,1~0,0,1\n0,0,2~0,0,2", "0,0,1~0,0,1\n0,1,1~0,1,1", fmt.Sprintf("%d,%d,1~%d,%d,%d", math.MaxInt, math.MaxInt, math.MaxInt, math.MaxInt, math.MaxInt), fmt.Sprintf("0,0,%d~0,0,%d", math.MaxInt, math.MaxInt)} {
		if _, err := parseBricks(input); err != nil {
			t.Fatal("valid rejected", err)
		}
	}
}

// The reference occupies individual cubes and lowers each brick one unit at a time.
func voxelSettle(bricks []Brick) ([]Brick, int) {
	bricks = append([]Brick(nil), bricks...)
	sort.SliceStable(bricks, func(i, j int) bool { return bricks[i].Start.Z < bricks[j].Start.Z })
	occupied := map[Point3D]int{}
	fallen := 0
	for i := range bricks {
		b := &bricks[i]
		b.Supports = nil
		b.SupportedBy = nil
		original := b.Start.Z
		for b.Start.Z > 1 {
			blocked := false
			for x := b.Start.X; x <= b.End.X; x++ {
				for y := b.Start.Y; y <= b.End.Y; y++ {
					if _, ok := occupied[Point3D{x, y, b.Start.Z - 1}]; ok {
						blocked = true
					}
				}
			}
			if blocked {
				break
			}
			b.Start.Z--
			b.End.Z--
		}
		if b.Start.Z != original {
			fallen++
		}
		parents := map[int]bool{}
		for x := b.Start.X; x <= b.End.X; x++ {
			for y := b.Start.Y; y <= b.End.Y; y++ {
				if parent, ok := occupied[Point3D{x, y, b.Start.Z - 1}]; ok {
					parents[parent] = true
				}
				for z := b.Start.Z; z <= b.End.Z; z++ {
					cell := Point3D{x, y, z}
					if _, ok := occupied[cell]; ok {
						panic("invalid test scene")
					}
					occupied[cell] = i
				}
			}
		}
		for parent := range parents {
			b.SupportedBy = append(b.SupportedBy, parent)
			bricks[parent].Supports = append(bricks[parent].Supports, i)
		}
		sort.Ints(b.SupportedBy)
	}
	return bricks, fallen
}
func physicalRemoval(bricks []Brick, removed int) int {
	survivors := make([]Brick, 0, len(bricks)-1)
	for i, b := range bricks {
		if i != removed {
			survivors = append(survivors, b)
		}
	}
	_, fallen := voxelSettle(survivors)
	return fallen
}
func TestBrickGeneratedOracles(t *testing.T) {
	for seed := uint64(0); seed < 256; seed++ {
		state := seed + 1
		var input strings.Builder
		raw := []Brick{}
		for i := 0; i < 12; i++ {
			state = state*1664525 + 1013904223
			x, y, z := int(state%5), int((state/5)%5), 1+i*4
			a, b := Point3D{x, y, z}, Point3D{x, y, z}
			switch state % 3 {
			case 0:
				b.X += 2
			case 1:
				b.Y += 2
			case 2:
				b.Z += 2
			}
			raw = append(raw, Brick{Start: a, End: b})
			fmt.Fprintf(&input, "%d,%d,%d~%d,%d,%d\n", a.X, a.Y, a.Z, b.X, b.Y, b.Z)
		}
		actual, err := parseBricks(input.String())
		if err != nil {
			t.Fatal(err)
		}
		expected, _ := voxelSettle(raw)
		total, safe := 0, 0
		for i, a := range actual {
			b := expected[i]
			parents := append([]int(nil), a.SupportedBy...)
			sort.Ints(parents)
			if a.Start != b.Start || a.End != b.End || !reflect.DeepEqual(a.Supports, b.Supports) || !reflect.DeepEqual(parents, b.SupportedBy) {
				t.Fatalf("settling mismatch seed %d brick %d", seed, i)
			}
			want := physicalRemoval(expected, i)
			got := cascadeSize(actual, i)
			if got != want {
				t.Fatalf("cascade mismatch seed %d removed %d: %d vs %d", seed, i, got, want)
			}
			total += want
			if want == 0 {
				safe++
			}
		}
		p1, e1 := (Day22{}).Part1(input.String())
		p2, e2 := (Day22{}).Part2(input.String())
		if e1 != nil || e2 != nil || p1 != strconv.Itoa(safe) || p2 != strconv.Itoa(total) {
			t.Fatal("parts differ from oracle")
		}
	}
}
func TestBrickReversalsAndSharedSupports(t *testing.T) {
	input := "0,0,1~0,0,3\n0,0,8~2,0,8\n2,0,12~2,2,12\n2,2,16~2,2,16"
	expected, _ := parseBricks(input)
	for mask := 0; mask < 16; mask++ {
		lines := strings.Split(input, "\n")
		for i, line := range lines {
			if mask&(1<<i) != 0 {
				a, b, _ := strings.Cut(line, "~")
				lines[i] = b + "~" + a
			}
		}
		actual, err := parseBricks(strings.Join(lines, "\n"))
		if err != nil || !reflect.DeepEqual(actual, expected) {
			t.Fatal("reversal changed brick", mask, err)
		}
	}
	bricks, err := parseBricks("0,0,1~2,0,1\n0,0,2~0,0,2\n2,0,2~2,0,2\n0,0,3~2,0,3\n1,0,4~1,0,4")
	if err != nil {
		t.Fatal(err)
	}
	for i, want := range []int{4, 0, 0, 1, 0} {
		if got := cascadeSize(bricks, i); got != want || got != physicalRemoval(bricks, i) {
			t.Fatal("fork/rejoin cascade", i, got)
		}
	}
	bridge, err := parseBricks("0,0,1~0,0,3\n1,0,1~1,0,1\n2,0,1~2,0,3\n0,0,9~2,0,9")
	if err != nil {
		t.Fatal(err)
	}
	sort.Ints(bridge[3].SupportedBy)
	if bridge[3].Start.Z != 4 || !reflect.DeepEqual(bridge[3].SupportedBy, []int{0, 2}) || len(bridge[1].Supports) != 0 {
		t.Fatal("incorrect bridge supporters")
	}
}
