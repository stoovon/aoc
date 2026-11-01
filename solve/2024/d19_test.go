package solve2024

import (
	"math/big"
	"strings"
	"testing"
)

func referenceTowelCount(patterns []string, design string) int64 {
	if design == "" {
		return 1
	}
	var count int64
	for _, pattern := range patterns {
		if len(pattern) <= len(design) && design[:len(pattern)] == pattern {
			count += referenceTowelCount(patterns, design[len(pattern):])
		}
	}
	return count
}

func TestTowelSeparateSearchesAgainstExhaustiveSplits(t *testing.T) {
	options := []string{"a", "b", "aa", "ab", "ba", "bb"}
	for subset := 1; subset < 1<<len(options); subset++ {
		var raw []string
		var patterns towelPatterns
		for i, text := range options {
			if subset&(1<<i) != 0 {
				raw = append(raw, text)
				patterns[text[0]] = append(patterns[text[0]], text)
			}
		}
		possible := make(map[string]bool)
		counts := make(map[string]*big.Int)
		for length := 0; length <= 7; length++ {
			for encoding := 0; encoding < 1<<length; encoding++ {
				var design strings.Builder
				for i := 0; i < length; i++ {
					if encoding&(1<<i) == 0 {
						design.WriteByte('a')
					} else {
						design.WriteByte('b')
					}
				}
				text := design.String()
				want := referenceTowelCount(raw, text)
				if got := patterns.canMake(text, possible); got != (want > 0) {
					t.Fatalf("patterns%v design%q possible%t want%d", raw, text, got, want)
				}
				if got := patterns.countWays(text, counts); got.Cmp(big.NewInt(want)) != 0 {
					t.Fatalf("patterns%v design%q ways%s want%d", raw, text, got, want)
				}
			}
		}
	}
}

func TestTowelPart1StopsAfterFirstArrangement(t *testing.T) {
	patterns := towelPatterns{}
	patterns['a'] = []string{"aa", "a"}
	cache := make(map[string]bool)
	if !patterns.canMake("aa", cache) {
		t.Fatal("possible design rejected")
	}
	if _, visited := cache["a"]; visited {
		t.Fatal("continued searching after successful first branch")
	}
}

func TestTowelExactCountsAndCacheOwnership(t *testing.T) {
	design := strings.Repeat("a", 100)
	input := "a, aa\n\n" + design + "\n" + design
	// Independent length recurrence for one- and two-character tiles.
	first, second := big.NewInt(1), big.NewInt(1)
	for length := 2; length <= 100; length++ {
		first, second = second, new(big.Int).Add(first, second)
	}
	want := new(big.Int).Lsh(new(big.Int).Set(second), 1)
	if got, err := (Day19{}).Part1(input); err != nil || got != "2" {
		t.Fatalf("overflow-independent possibility: %q %v", got, err)
	}
	if got, err := (Day19{}).Part2(input); err != nil || got != want.String() {
		t.Fatalf("exact repeated-design counts: %q %v want%s", got, err, want)
	}
	towels, err := parseTowels(input)
	if err != nil {
		t.Fatal(err)
	}
	cache := make(map[string]*big.Int)
	saved := towels.patterns.countWays("aa", cache)
	towels.patterns.countWays(design, cache)
	if saved.String() != "2" {
		t.Fatal("larger computation mutated cached child count")
	}
}

func TestTowelDuplicatePatternsAndDesignRecords(t *testing.T) {
	for part, tc := range []struct {
		fn   func(string) (string, error)
		want string
	}{{(Day19{}).Part1, "2"}, {(Day19{}).Part2, "4"}} {
		got, err := tc.fn(" a, aa, a \r\n\r\naa\r\naa\r\nb\r\n")
		if err != nil || got != tc.want {
			t.Fatalf("part%d duplicate rules: %q %v want%q", part+1, got, err, tc.want)
		}
	}
}

func TestTowelRejectsMalformedInput(t *testing.T) {
	for _, input := range []string{"", "a", "a\n\n", "a,\n\na", ",a\n\na", "a,,aa\n\na", "a b\n\nab", "A\n\na", "é\n\na", "a\naa\n\na", "a\n\na\n\nb", "a\n\na b"} {
		for _, fn := range []func(string) (string, error){(Day19{}).Part1, (Day19{}).Part2} {
			if _, err := fn(input); err == nil {
				t.Fatalf("accepted malformed input%q", input)
			}
		}
	}
	// A zero count must be cached just as reliably as a positive count.
	towels, err := parseTowels("a\n\nb")
	if err != nil {
		t.Fatal(err)
	}
	cache := make(map[string]*big.Int)
	if got := towels.patterns.countWays("b", cache); got.Sign() != 0 {
		t.Fatal("impossible design has nonzero count")
	}
	if _, found := cache["b"]; !found {
		t.Fatal("zero count not cached")
	}
}
