package solve2018

import (
	"strconv"
	"strings"
	"testing"
)

func TestInventoryChecksumCountsEachIDOnce(t *testing.T) {
	for _, input := range []string{"aaabbb\nccde\n", "aaabb\n"} {
		got, err := (Day2{}).Part1(input)
		if err != nil || got != "1" {
			t.Fatalf("checksum: got %q %v want 1", got, err)
		}
	}
}

func TestInventoryExactOneDifference(t *testing.T) {
	for _, tc := range []struct {
		first, second, common string
		matches               bool
	}{
		{"wxyz", "wxqz", "wxz", true},
		{"wxyz", "qxyz", "xyz", true},
		{"wxyz", "wxyq", "wxy", true},
		{"a", "b", "", true},
		{"wxyz", "wxyz", "", false},
		{"wxyz", "qqyz", "", false},
		{"wxyz", "wx", "", false},
		{"wx", "wxyz", "", false},
		{"", "", "", false},
	} {
		common, matches := commonLetters(tc.first, tc.second)
		if common != tc.common || matches != tc.matches {
			t.Fatalf("%q vs %q: got %q %t want %q %t", tc.first, tc.second, common, matches, tc.common, tc.matches)
		}
	}
}

func TestInventoryLastIDAndEmptyCommonResult(t *testing.T) {
	for _, input := range []string{"wxyz\nwxqz", "wxyz\nwxqz\n"} {
		if got, err := (Day2{}).Part2(input); err != nil || got != "wxz" {
			t.Fatalf("final ID: %q %v", got, err)
		}
	}
	if got, err := (Day2{}).Part2("a\nb"); err != nil || got != "" {
		t.Fatalf("successful empty common letters: %q %v", got, err)
	}
}

func inventorySyntheticWords(length int) []string {
	count := 1
	for i := 0; i < length; i++ {
		count *= 3
	}
	words := make([]string, count)
	for code := 0; code < count; code++ {
		word := make([]byte, length)
		rest := code
		for i := range word {
			word[i] = 'a' + byte(rest%3)
			rest /= 3
		}
		words[code] = string(word)
	}
	return words
}

func TestInventorySmallIDsAgainstDirectSpecification(t *testing.T) {
	for length := 1; length <= 6; length++ {
		for _, id := range inventorySyntheticWords(length) {
			twice, thrice := false, false
			for _, letter := range []string{"a", "b", "c"} {
				count := strings.Count(id, letter)
				twice = twice || count == 2
				thrice = thrice || count == 3
			}
			pairs, triples := 1, 0
			if twice {
				pairs++
			}
			if thrice {
				triples++
			}
			got, err := (Day2{}).Part1(id + "\nzz")
			if err != nil || got != strconv.Itoa(pairs*triples) {
				t.Fatalf("checksum %q: %q %v", id, got, err)
			}
		}
	}
	for length := 1; length <= 5; length++ {
		words := inventorySyntheticWords(length)
		for _, first := range words {
			for _, second := range words {
				differences := 0
				common := []byte{}
				for i := range first {
					if first[i] != second[i] {
						differences++
					} else {
						common = append(common, first[i])
					}
				}
				wantMatches := differences == 1
				want := ""
				if wantMatches {
					want = string(common)
				}
				got, matches := commonLetters(first, second)
				if got != want || matches != wantMatches {
					t.Fatalf("comparison %q %q: %q %t want %q %t", first, second, got, matches, want, wantMatches)
				}
			}
		}
	}
}

func TestInventoryInputValidationAndNoMatch(t *testing.T) {
	for _, input := range []string{"", "abc\n\ndef", "abc\n\n", "abc\ndef ", "ABC", "ab1", "ébc", " abc"} {
		for _, solver := range []func(string) (string, error){(Day2{}).Part1, (Day2{}).Part2} {
			if _, err := solver(input); err == nil {
				t.Fatalf("accepted malformed input %q", input)
			}
		}
	}
	if got, err := (Day2{}).Part2("wxyz\r\nwxqz\r\n"); err != nil || got != "wxz" {
		t.Fatalf("CRLF: %q %v", got, err)
	}
	for _, input := range []string{"abc", "abc\nabc", "ab\nabcdef", "ab\ncd"} {
		if _, err := (Day2{}).Part2(input); err == nil {
			t.Fatalf("accepted no-match input %q", input)
		}
	}
	if _, err := (Day2{}).Part1("abc\nA"); err == nil || !strings.Contains(err.Error(), "line 2:") {
		t.Fatalf("missing line context: %v", err)
	}
}
