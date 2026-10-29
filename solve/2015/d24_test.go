package solve2015

import (
	"math/bits"
	"math/rand/v2"
	"slices"
	"strconv"
	"strings"
	"testing"
)

// Assign remaining packages to explicit bins, rather than using group-subset memoization.
func packageBalanceReference(weights []int, groups int) (int, bool) {
	total := 0
	for _, weight := range weights {
		total += weight
	}
	if total%groups != 0 || groups > len(weights) {
		return 0, false
	}
	target := total / groups
	bestCount, bestQE := len(weights)+1, 0
	for mask := 1; mask < 1<<len(weights); mask++ {
		count := bits.OnesCount(uint(mask))
		if count > bestCount {
			continue
		}
		sum, product := 0, 1
		var rest []int
		for index, weight := range weights {
			if mask&(1<<index) != 0 {
				sum += weight
				product *= weight
			} else {
				rest = append(rest, weight)
			}
		}
		if sum != target || count == bestCount && product >= bestQE {
			continue
		}
		bins := make([]int, groups-1)
		var place func(int) bool
		place = func(index int) bool {
			if index == len(rest) {
				for _, sum := range bins {
					if sum != target {
						return false
					}
				}
				return true
			}
			for bin := range bins {
				if bins[bin]+rest[index] > target {
					continue
				}
				bins[bin] += rest[index]
				if place(index + 1) {
					return true
				}
				bins[bin] -= rest[index]
			}
			return false
		}
		if place(0) {
			bestCount, bestQE = count, product
		}
	}
	return bestQE, bestQE != 0
}

func TestDay24PartitionsMatchExplicitBinAssignments(t *testing.T) {
	random := rand.New(rand.NewPCG(2015, 24))
	for range 300 {
		weights := make([]int, 2+random.IntN(9))
		for index := range weights {
			weights[index] = 1 + random.IntN(9)
		}
		groups := 2 + random.IntN(3)
		want, possible := packageBalanceReference(weights, groups)
		original := slices.Clone(weights)
		got, err := balancePackages(weights, groups)
		if (err == nil) != possible || err == nil && got.String() != strconv.Itoa(want) {
			t.Fatalf("weights %v, groups %d: got %v, %v; want %d, possible=%v", weights, groups, got, err, want, possible)
		}
		if !slices.Equal(weights, original) {
			t.Fatal("search changed caller's weights")
		}
	}
}

func TestDay24RemainingGroupsMustBeFeasible(t *testing.T) {
	weights := []int{1, 3, 9, 11, 18, 22, 24, 31, 39, 40}
	got, err := balancePackages(weights, 3)
	if err != nil || got.String() != "6318" {
		t.Fatalf("got %v, %v; want 6318", got, err)
	}
	// The old 2808 candidate [3,24,39] leaves no way to form the other groups.
	if _, err := balancePackages([]int{2, 2, 2, 3}, 3); err == nil {
		t.Fatal("one target-weight group does not imply a complete partition")
	}
}

func TestDay24CardinalityPrecedesProductAndProductsStayExact(t *testing.T) {
	for _, test := range []struct {
		weights []int
		groups  int
		want    string
	}{
		{[]int{1, 1, 1, 1, 2, 2, 3, 3, 4, 6}, 3, "12"},
		{[]int{1, 1, 1, 1, 1, 1}, 3, "1"},
		{[]int{10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000}, 3, "1000000000000000000000"},
	} {
		got, err := balancePackages(test.weights, test.groups)
		if err != nil || got.String() != test.want {
			t.Fatalf("weights %v, groups %d: got %v, %v; want %s", test.weights, test.groups, got, err, test.want)
		}
	}
}

func TestDay24HighestSupportedPackageBit(t *testing.T) {
	weights := make([]int, 63)
	for index := range weights {
		weights[index] = 1
	}
	for _, groups := range []int{3, 63} {
		got, err := balancePackages(weights, groups)
		if err != nil || got.String() != "1" {
			t.Fatalf("groups %d: got %v, %v; want 1", groups, got, err)
		}
	}
}

func TestDay24SyntheticPackages(t *testing.T) {
	input := "1\r\n2\r\n3\r\n4\r\n5\r\n6\r\n7\r\n8\r\n"
	for part, test := range []struct {
		solver func(string) (string, error)
		want   string
	}{{(Day24{}).Part1, "32"}, {(Day24{}).Part2, "8"}} {
		got, err := test.solver(input)
		if err != nil || got != test.want {
			t.Fatalf("part %d: got %q, %v; want %q", part+1, got, err, test.want)
		}
	}
}

func TestDay24InputAndSearchErrors(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	for _, test := range []struct {
		weights []int
		groups  int
	}{
		{nil, 3}, {[]int{1, 2}, 0}, {[]int{1, 2}, 1}, {[]int{1, 2}, 3},
		{[]int{0, 1, 2}, 3}, {[]int{-1, 1, 2}, 3}, {[]int{1, 2, 4}, 2},
		{[]int{1, 1, 4}, 3}, {[]int{maxInt, 1}, 2}, {make([]int, 64), 3},
	} {
		if _, err := balancePackages(test.weights, test.groups); err == nil {
			t.Fatalf("weights %v, groups %d should return an error", test.weights, test.groups)
		}
	}
	for _, solver := range []func(string) (string, error){(Day24{}).Part1, (Day24{}).Part2} {
		for _, input := range []string{"", "1\ninvalid", "0", "-1", "9999999999999999999999999999999999999"} {
			if _, err := solver(input); err == nil {
				t.Fatalf("input %q should return an error", input)
			}
		}
		if _, err := solver("1\ninvalid"); err == nil || !strings.HasPrefix(err.Error(), "line 2:") {
			t.Fatalf("missing line context: %v", err)
		}
	}
}
