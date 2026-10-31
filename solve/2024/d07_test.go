package solve2024

import (
	"math"
	"math/big"
	"strconv"
	"strings"
	"testing"
)

func calibrationForwardReference(numbers []uint64, concat bool) map[string]bool {
	result := make(map[string]bool)
	var visit func(int, *big.Int)
	visit = func(index int, value *big.Int) {
		if index == len(numbers) {
			result[value.String()] = true
			return
		}
		next := new(big.Int).SetUint64(numbers[index])
		visit(index+1, new(big.Int).Add(value, next))
		visit(index+1, new(big.Int).Mul(value, next))
		if concat {
			combined, ok := new(big.Int).SetString(value.String()+next.String(), 10)
			if !ok {
				panic("invalid synthetic concatenation")
			}
			visit(index+1, combined)
		}
	}
	visit(1, new(big.Int).SetUint64(numbers[0]))
	return result
}

func TestCalibrationMatchesExactForwardEvaluation(t *testing.T) {
	combinations := 1
	for length := 1; length <= 5; length++ {
		combinations *= 4
		for code := 0; code < combinations; code++ {
			numbers := make([]uint64, length)
			remaining := code
			for i := range numbers {
				numbers[i] = uint64(remaining % 4)
				remaining /= 4
			}
			for _, concat := range []bool{false, true} {
				reachable := calibrationForwardReference(numbers, concat)
				for target := uint64(0); target < 64; target++ {
					want := reachable[strconv.FormatUint(target, 10)]
					got := calibrationMatches(calibrationEquation{target: target, numbers: numbers}, concat)
					if got != want {
						t.Fatalf("numbers=%v target=%d concat=%t: got %t want %t", numbers, target, concat, got, want)
					}
				}
			}
		}
	}
}

func TestCalibrationExactBoundariesAndZero(t *testing.T) {
	cases := []calibrationEquation{
		{1, []uint64{math.MaxInt64, 2, 3}},
		{math.MaxInt64, []uint64{922337203685477580, 8}},
		{7, []uint64{math.MaxUint64, 2, 0, 7}},
		{math.MaxUint64, []uint64{0, math.MaxUint64}},
		{math.MaxUint64, []uint64{2, math.MaxUint64}},
		{10000000000000000000, []uint64{0, 10000000000000000000}},
		{0, []uint64{math.MaxUint64, math.MaxUint64, 0}},
		{0, []uint64{1, 2}},
		{0, []uint64{0}},
		{17, []uint64{17}},
		{123, []uint64{1, 2, 3}},
		{100, []uint64{10, 0}},
		{20, []uint64{2, 3, 4}},
		{14, []uint64{2, 3, 4}},
	}
	for _, equation := range cases {
		for _, concat := range []bool{false, true} {
			want := calibrationForwardReference(equation.numbers, concat)[strconv.FormatUint(equation.target, 10)]
			if got := calibrationMatches(equation, concat); got != want {
				t.Fatalf("%v concat=%t: got %t want %t", equation, concat, got, want)
			}
		}
	}
	for _, part := range []func(string) (string, error){(Day7{}).Part1, (Day7{}).Part2} {
		if got, err := part("1: 9223372036854775807 2 3"); err != nil || got != "0" {
			t.Fatal(got, err)
		}
		if got, err := part("7: 18446744073709551615 2 0 7"); err != nil || got != "7" {
			t.Fatal(got, err)
		}
		if _, err := part("18446744073709551615: 18446744073709551615\n1: 1"); err == nil || !strings.Contains(err.Error(), "total exceeds") {
			t.Fatal(err)
		}
	}
	if got, err := (Day7{}).Part2("9223372036854775807: 922337203685477580 8"); err != nil || got != "0" {
		t.Fatal(got, err)
	}
}

func TestCalibrationParsingAndPartComposition(t *testing.T) {
	input := "20: 2 3 4\n123: 1 2 3\n19: 2 3 4"
	for _, lineEndings := range []string{input, input + "\n", strings.ReplaceAll(input, "\n", "\r\n") + "\r\n"} {
		if got, err := (Day7{}).Part1(lineEndings); err != nil || got != "20" {
			t.Fatal(got, err)
		}
		if got, err := (Day7{}).Part2(lineEndings); err != nil || got != "143" {
			t.Fatal(got, err)
		}
	}
	if got, err := (Day7{}).Part2(" 0012 : 01   02 "); err != nil || got != "12" {
		t.Fatal(got, err)
	}
	for _, input := range []string{"", " ", "1", "1:", ": 1", "-1: 1", "+1: 1", "1: -1", "1: +1", "1: 1.0", "1: １", "1:: 1", "1: 1\n\n2: 2", "18446744073709551616: 1", "1: 18446744073709551616"} {
		for _, part := range []func(string) (string, error){(Day7{}).Part1, (Day7{}).Part2} {
			if _, err := part(input); err == nil {
				t.Fatalf("accepted %q", input)
			}
		}
	}
	if _, err := (Day7{}).Part1("1: 1\nno colon"); err == nil || !strings.Contains(err.Error(), "line 2:") {
		t.Fatal(err)
	}
}

func TestCalibrationPrefixUpperBounds(t *testing.T) {
	ones := make([]uint64, 80)
	for i := range ones {
		ones[i] = 1
	}
	if calibrationMatches(calibrationEquation{81, ones}, false) {
		t.Fatal("exceeded maximum prefix sum")
	}
	zeros := make([]uint64, 80)
	for _, concat := range []bool{false, true} {
		if calibrationMatches(calibrationEquation{1, zeros}, concat) {
			t.Fatal("zero prefix produced one")
		}
		if !calibrationMatches(calibrationEquation{2, ones}, concat) {
			t.Fatal("lost long reachable equation")
		}
		equation := calibrationEquation{7, []uint64{math.MaxUint64, math.MaxUint64, 0, 7}}
		if !calibrationMatches(equation, concat) {
			t.Fatal("saturated bound discarded a valid zero reset")
		}
	}
}
