package solve2019

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

func naiveFFTMessage(input string) string {
	signal := make([]int, len(input))
	for i := range input {
		signal[i] = int(input[i] - '0')
	}
	pattern := [4]int{0, 1, 0, -1}
	for phase := 0; phase < 100; phase++ {
		next := make([]int, len(signal))
		for i := range signal {
			sum := 0
			for j, digit := range signal {
				sum += digit * pattern[((j+1)/(i+1))%4]
			}
			if sum < 0 {
				sum = -sum
			}
			next[i] = sum % 10
		}
		signal = next
	}
	var out strings.Builder
	for _, digit := range signal[:min(8, len(signal))] {
		out.WriteByte(byte('0' + digit))
	}
	return out.String()
}

func TestDay16Part1MatchesNaiveTransform(t *testing.T) {
	rng := rand.New(rand.NewSource(16))
	for n := 1; n <= 80; n++ {
		for trial := 0; trial < 4; trial++ {
			input := make([]byte, n)
			for i := range input {
				input[i] = byte('0' + rng.Intn(10))
			}
			got, err := (Day16{}).Part1(string(input))
			want := naiveFFTMessage(string(input))
			if err != nil || got != want {
				t.Fatalf("length %d trial %d: got %q, %v; want %q", n, trial, got, err, want)
			}
		}
	}
}

func TestDay16RejectsMalformedSignal(t *testing.T) {
	for _, input := range []string{"", " \r\n", "12x34", "12 34", "12\n34", "１２３４５６７８", "1234567\x00"} {
		for part, solve := range []func(string) (string, error){(Day16{}).Part1, (Day16{}).Part2} {
			if _, err := solve(input); err == nil {
				t.Errorf("part %d accepted %q", part+1, input)
			}
		}
	}
}

func TestDay16TrimsSurroundingWhitespace(t *testing.T) {
	for _, tc := range []struct {
		input string
		solve func(string) (string, error)
	}{
		{"31425768912", (Day16{}).Part1},
		{"00799923", (Day16{}).Part2},
	} {
		want, err := tc.solve(tc.input)
		if err != nil {
			t.Fatal(err)
		}
		got, err := tc.solve(" \t" + tc.input + "\r\n")
		if err != nil || got != want {
			t.Fatalf("trim: got %q, %v; want %q", got, err, want)
		}
	}
}

func TestDay16Part2OffsetBounds(t *testing.T) {
	for _, input := range []string{
		"123456",  // Too short to encode an offset.
		"0000000", // Before the second half.
		"0034999", // Immediately before the second half.
		"0069993", // Only seven digits remain.
		"0070000", // Offset is at the end.
		"9999999", // Offset is beyond the end.
	} {
		if _, err := (Day16{}).Part2(input); err == nil {
			t.Errorf("accepted invalid offset in %q", input)
		}
	}
	for _, input := range []string{"0035000", "0069992", "0069991"} {
		got, err := (Day16{}).Part2(input)
		if err != nil || len(got) != 8 {
			t.Errorf("valid boundary %q: got %q, %v", input, got, err)
		}
	}
}

func TestDay16Part2MatchesDirectSuffixSums(t *testing.T) {
	for _, length := range []int{8, 9, 17, 31, 64} {
		for _, tail := range []string{"3", "97", "24581"} {
			baseLen := 7 + len(tail)
			total := baseLen * 10000
			offset := total - length
			input := fmt.Sprintf("%07d%s", offset, tail)
			signal := make([]int, length)
			for i := range signal {
				signal[i] = int(input[(offset+i)%baseLen] - '0')
			}
			for phase := 0; phase < 100; phase++ {
				next := make([]int, length)
				for i := range signal {
					for _, digit := range signal[i:] {
						next[i] += digit
					}
					next[i] %= 10
				}
				signal = next
			}
			var want strings.Builder
			for _, digit := range signal[:8] {
				want.WriteByte(byte('0' + digit))
			}
			got, err := (Day16{}).Part2(input)
			if err != nil || got != want.String() {
				t.Fatalf("length %d tail %q: got %q, %v; want %q", length, tail, got, err, want.String())
			}
		}
	}
}
