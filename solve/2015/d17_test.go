package solve2015

import (
	"math/bits"
	"math/rand/v2"
	"strconv"
	"strings"
	"testing"
)

func containerArrangementsReference(sizes []int, target int) containerCounts {
	var counts containerCounts
	for mask := 0; mask < 1<<len(sizes); mask++ {
		volume := 0
		for index, size := range sizes {
			if mask&(1<<index) != 0 {
				volume += size
			}
		}
		if volume != target {
			continue
		}
		counts.totalWays++
		used := bits.OnesCount(uint(mask))
		if counts.minWays == 0 || used < counts.minContainers {
			counts.minContainers, counts.minWays = used, 1
		} else if used == counts.minContainers {
			counts.minWays++
		}
	}
	return counts
}

func TestDay17VolumeCountsMatchSubsets(t *testing.T) {
	random := rand.New(rand.NewPCG(2015, 17))
	for trial := range 300 {
		sizes := make([]int, random.IntN(13))
		for index := range sizes {
			sizes[index] = random.IntN(11)
		}
		target := random.IntN(31)
		want := containerArrangementsReference(sizes, target)
		got, err := countContainerArrangements(sizes, target)
		if err != nil || got != want {
			t.Fatalf("trial %d, sizes %v, target %d: got %+v, %v; want %+v", trial, sizes, target, got, err, want)
		}
	}
}

func TestDay17ContainersAreDistinctAndUsedOnce(t *testing.T) {
	for _, test := range []struct {
		sizes  []int
		target int
		want   containerCounts
	}{
		{[]int{9, 7, 4, 2, 2}, 11, containerCounts{4, 2, 3}},
		{[]int{5}, 10, containerCounts{}},
		{[]int{5, 5}, 5, containerCounts{2, 1, 2}},
		{[]int{5, 5}, 10, containerCounts{1, 2, 1}},
		{[]int{0, 0, 5}, 5, containerCounts{4, 1, 1}},
		{[]int{0, 0}, 0, containerCounts{4, 0, 1}},
		{nil, 0, containerCounts{1, 0, 1}},
		{nil, 5, containerCounts{}},
	} {
		got, err := countContainerArrangements(test.sizes, test.target)
		if err != nil || got != test.want {
			t.Fatalf("sizes %v, target %d: got %+v, %v; want %+v", test.sizes, test.target, got, err, test.want)
		}
	}
}

func TestDay17InputAndCountErrors(t *testing.T) {
	for _, solver := range []func(string) (string, error){(Day17{}).Part1, (Day17{}).Part2} {
		for _, input := range []string{"", "-1", "10\ninvalid", "10\n\n20", "999999999999999999999999999999999"} {
			if _, err := solver(input); err == nil {
				t.Fatalf("input %q should return an error", input)
			}
		}
		if _, err := solver("10\ninvalid"); err == nil || !strings.HasPrefix(err.Error(), "line 2:") {
			t.Fatalf("missing line context: %v", err)
		}
	}
	for _, test := range []struct {
		sizes  []int
		target int
	}{
		{[]int{-1}, 5}, {nil, -1}, {nil, 1 << 20}, {make([]int, strconv.IntSize), 0},
	} {
		if _, err := countContainerArrangements(test.sizes, test.target); err == nil {
			t.Fatalf("sizes %v, target %d should return an error", test.sizes, test.target)
		}
	}
}

func TestDay17PartsUseCorrectSummary(t *testing.T) {
	input := "100\r\n50\r\n50\r\n0\r\n"
	for part, test := range []struct {
		solver func(string) (string, error)
		want   string
	}{{(Day17{}).Part1, "4"}, {(Day17{}).Part2, "2"}} {
		got, err := test.solver(input)
		if err != nil || got != test.want {
			t.Fatalf("part %d: got %q, %v; want %q", part+1, got, err, test.want)
		}
	}
}
