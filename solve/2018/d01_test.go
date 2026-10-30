package solve2018

import (
	"fmt"
	"math"
	"math/big"
	"strings"
	"testing"
)

func TestFrequencyFirstRepeat(t *testing.T) {
	for _, tc := range []struct{ input, want string }{{"+6\n-6", "0"}, {"+6\n+4\n-4\n-6", "6"}, {"+6\n-4", "6"}, {"-6\n+4", "-6"}, {"+8\n-3\n-2", "8"}, {"0", "0"}} {
		got, err := (Day1{}).Part2(tc.input)
		if err != nil || got != tc.want {
			t.Fatal("wrong chronological repeat", got, err, tc.want)
		}
	}
}
func TestFrequencyNoRepeatAndEnormousCycles(t *testing.T) {
	for _, input := range []string{"+2", "-2", "+5\n-2", "-5\n+2", "+2\n+4"} {
		if _, err := (Day1{}).Part2(input); err == nil || !strings.Contains(err.Error(), "never repeats") {
			t.Fatal("missing no-repeat error", err)
		}
	}
	for _, tc := range []struct{ input, want string }{{"1000000000000\n-999999999999", "1000000000000"}, {"-1000000000000\n999999999999", "-1000000000000"}} {
		got, err := (Day1{}).Part2(tc.input)
		if err != nil || got != tc.want {
			t.Fatal("failed to jump cycles", err)
		}
	}
}
func TestFrequencyParserAndWideTotals(t *testing.T) {
	for _, input := range []string{"", " \n\t", "+6\n-6\njunk", "2 4", "2\n\n4", "--2", "1.5", "9223372036854775808"} {
		for _, solve := range []func(string) (string, error){(Day1{}).Part1, (Day1{}).Part2} {
			if _, err := solve(input); err == nil {
				t.Fatal("accepted invalid input")
			}
		}
	}
	for _, solve := range []func(string) (string, error){(Day1{}).Part1, (Day1{}).Part2} {
		if _, err := solve("+6\n-6\njunk"); err == nil || !strings.HasPrefix(err.Error(), "line 3:") {
			t.Fatal("missing line context")
		}
	}
	if got, err := (Day1{}).Part1(" +6\r\n -4\r\n"); err != nil || got != "2" {
		t.Fatal("CRLF", err)
	}
	max := int64(math.MaxInt64)
	want := new(big.Int).Mul(big.NewInt(max), big.NewInt(2)).String()
	if got, err := (Day1{}).Part1(fmt.Sprintf("%d\n%d", max, max)); err != nil || got != want {
		t.Fatal("wide sum", err)
	}
	input := fmt.Sprintf("%d\n%d\n%d\n%d", max, max-2, -max+4, int64(math.MinInt64))
	if got, err := (Day1{}).Part1(input); err != nil || got != "1" {
		t.Fatal("wide prefix sum", err)
	}
	want = new(big.Int).Add(big.NewInt(max), big.NewInt(2)).String()
	if got, err := (Day1{}).Part2(input); err != nil || got != want {
		t.Fatal("wide repeated frequency", err)
	}
	if _, err := firstRepeatedFrequency(nil); err == nil {
		t.Fatal("empty list")
	}
}
func frequencySimulate(changes []int64) (int64, bool) {
	seen := map[int64]bool{0: true}
	frequency := int64(0)
	// For at most six changes of magnitude at most two, 200 steps exceed
	// the possible first-repeat time implied by the prefix span.
	for step := 0; step < 200; step++ {
		frequency += changes[step%len(changes)]
		if seen[frequency] {
			return frequency, true
		}
		seen[frequency] = true
	}
	return 0, false
}
func TestFrequencyExhaustiveSmallSimulation(t *testing.T) {
	cases := 1
	for length := 1; length <= 6; length++ {
		cases *= 5
		for code := 0; code < cases; code++ {
			changes := make([]int64, length)
			encoded := code
			for i := range changes {
				changes[i] = int64(encoded%5) - 2
				encoded /= 5
			}
			got, err := firstRepeatedFrequency(changes)
			want, exists := frequencySimulate(changes)
			if exists {
				if err != nil || got.Cmp(big.NewInt(want)) != 0 {
					t.Fatal("oracle mismatch", changes, err)
				}
			} else if err == nil || !strings.Contains(err.Error(), "never repeats") {
				t.Fatal("oracle no-repeat mismatch", changes, err)
			}
		}
	}
}
