package solve2025

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestDay8ClosestEdgesMatchSortedPairs(t *testing.T) {
	rng := rand.New(rand.NewSource(8))
	for n := 0; n <= 30; n++ {
		points := make([]vec3, n)
		for i := range points {
			points[i] = vec3{int64(rng.Intn(9)), int64(rng.Intn(9)), int64(rng.Intn(9))}
		}
		want := []edge{}
		for i, p := range points {
			for j, q := range points {
				if i >= j {
					continue
				}
				dx, dy, dz := p.x-q.x, p.y-q.y, p.z-q.z
				want = append(want, edge{i, j, dx*dx + dy*dy + dz*dz})
			}
		}
		sort.Slice(want, func(i, j int) bool {
			a, b := want[i], want[j]
			if a.d2 != b.d2 {
				return a.d2 < b.d2
			}
			if a.i != b.i {
				return a.i < b.i
			}
			return a.j < b.j
		})
		for _, limit := range []int{0, 1, 2, 7, 25, 1000} {
			got, err := buildEdges(points, limit)
			if err != nil {
				t.Fatal(err)
			}
			sort.Slice(got, func(i, j int) bool { return edgeLess(got[i], got[j]) })
			count := len(want)
			if limit > 0 {
				count = min(count, limit)
			}
			if !reflect.DeepEqual(got, want[:count]) {
				t.Fatalf("n=%d limit=%d: selected edges differ", n, limit)
			}
		}
	}
}

func TestDay8ConnectionLimitIncludesRedundantPairs(t *testing.T) {
	edges := []edge{{0, 1, 1}, {1, 2, 2}, {0, 2, 3}, {2, 3, 4}}
	got := processEdges(4, edges, 3, false)
	sort.Slice(got.componentSizes, func(i, j int) bool { return got.componentSizes[i] > got.componentSizes[j] })
	if !reflect.DeepEqual(got.componentSizes, []int64{3, 1}) {
		t.Fatalf("sizes %v; want [3 1]", got.componentSizes)
	}
	if got.lastMergeI != 1 || got.lastMergeJ != 2 {
		t.Fatal("redundant pair counted as a merge")
	}
}

func TestDay8DistanceMatchesExactArithmetic(t *testing.T) {
	values := []int64{math.MinInt64, math.MinInt64 + 1, -3037000500, -3037000499, -1, 0, 1, 3037000499, 3037000500, math.MaxInt64 - 1, math.MaxInt64}
	check := func(a, b vec3) {
		t.Helper()
		want := new(big.Int)
		for _, pair := range [3][2]int64{{a.x, b.x}, {a.y, b.y}, {a.z, b.z}} {
			delta := new(big.Int).Sub(big.NewInt(pair[0]), big.NewInt(pair[1]))
			want.Add(want, delta.Mul(delta, delta))
		}
		got, err := squaredDistance(a, b)
		if want.IsInt64() {
			if err != nil || got != want.Int64() {
				t.Fatalf("distance %v %v: %d, %v; want %s", a, b, got, err, want)
			}
		} else if err == nil {
			t.Fatalf("accepted overflowing distance %v %v", a, b)
		}
	}
	for _, a := range values {
		for _, b := range values {
			check(vec3{x: a}, vec3{x: b})
		}
	}
	for _, a := range []vec3{{3037000499, 0, 0}, {3037000499, 100000, 0}, {2000000000, 2000000000, 2000000000}} {
		check(a, vec3{})
	}
	rng := rand.New(rand.NewSource(81))
	for i := 0; i < 500; i++ {
		a, b := vec3{}, vec3{}
		for _, v := range []*int64{&a.x, &a.y, &a.z, &b.x, &b.y, &b.z} {
			*v = rng.Int63n(6000000001) - 3000000000
		}
		check(a, b)
	}
}

func TestDay8AnswerProductMatchesExactArithmetic(t *testing.T) {
	values := []int64{math.MinInt64, math.MinInt64 + 1, -3037000500, -2, -1, 0, 1, 2, 3037000500, math.MaxInt64}
	for _, a := range values {
		for _, b := range values {
			want := new(big.Int).Mul(big.NewInt(a), big.NewInt(b))
			got, err := multiplyCircuitAnswer(a, b)
			if want.IsInt64() {
				if err != nil || got != want.Int64() {
					t.Fatalf("%d*%d: %d,%v; want %s", a, b, got, err, want)
				}
			} else if err == nil {
				t.Fatalf("accepted overflowing product %d*%d", a, b)
			}
		}
	}
}

func TestDay8PropagatesArithmeticErrors(t *testing.T) {
	for _, input := range []string{
		"-9223372036854775808,0,0\n9223372036854775807,0,0",
		"0,0,0\n3037000500,0,0",
		"0,0,0\n2000000000,2000000000,2000000000",
	} {
		for part, solve := range []func(string) (string, error){(Day8{}).Part1, (Day8{}).Part2} {
			if _, err := solve(input); err == nil {
				t.Errorf("part %d accepted overflowing distance", part+1)
			}
		}
	}
	if _, err := (Day8{}).Part2("4000000000,0,0\n4000000001,0,0"); err == nil {
		t.Error("accepted overflowing answer")
	}
}

func TestDay8Part1RequiresThreeActualCircuits(t *testing.T) {
	for _, input := range []string{"1,0,0", "1,0,0\n3,0,0", "1,0,0\n3,0,0\n8,0,0"} {
		if _, err := (Day8{}).Part1(input); err == nil {
			t.Errorf("accepted fewer than three resulting circuits for %q", input)
		}
	}
}

func TestDay8RequiresAtLeastTwoBoxes(t *testing.T) {
	for _, input := range []string{"", " \t\r\n", "4,5,6", " \t4,5,6\r\n"} {
		for part, solve := range []func(string) (string, error){(Day8{}).Part1, (Day8{}).Part2} {
			if _, err := solve(input); err == nil {
				t.Errorf("part %d accepted fewer than two boxes", part+1)
			}
		}
	}
	got, err := (Day8{}).Part2("2,3,4\n5,7,11")
	if err != nil || got != "10" {
		t.Fatalf("two-box connection: got %q, %v; want 10", got, err)
	}
}

func TestDay8Part1MultipliesActualCircuitSizes(t *testing.T) {
	var input strings.Builder
	for x := 0; x < 46; x++ {
		fmt.Fprintf(&input, "%d,0,0\n", x)
	}
	input.WriteString("1000,0,0\n2000,0,0")
	// All 1000 selected pairs are among the 46 nearby boxes.
	got, err := (Day8{}).Part1(input.String())
	if err != nil || got != "46" {
		t.Fatalf("got %q, %v; want 46", got, err)
	}
}

func TestDay8Part2FinalMergeAndTies(t *testing.T) {
	for _, input := range []string{
		"1,0,0\n4,0,0\n6,0,0\n11,0,0",
		" \t1,0,0\r\n4,0,0\r\n6,0,0\r\n11,0,0\r\n",
	} {
		// The 1-to-6 pair ties the final 6-to-11 pair but is already connected.
		got, err := (Day8{}).Part2(input)
		if err != nil || got != "66" {
			t.Fatalf("got %q, %v; want 66", got, err)
		}
	}
}
