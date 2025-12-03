package solve2025

import (
	"errors"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day2 struct {
}

func (d Day2) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2025, Day: 2}
}

func (d Day2) Part1(input string) (string, error) {
	total, err := sumRangesBy(input, sumDoublesInRange)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(total, 10), nil
}

func (d Day2) Part2(input string) (string, error) {
	total, err := sumRangesBy(input, sumRepeatedAtLeastTwiceInRange)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(total, 10), nil
}

func init() {
	solve.Register(Day2{})
}

// sumRangesBy parses the single-line comma-separated ranges and returns the summed
// result of applying f to each numeric range [a,b]. Leading zeros are ignored by parsing.
func sumRangesBy(input string, f func(lo, hi int64) int64) (int64, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return 0, errors.New("empty input")
	}
	parts := strings.Split(input, ",")
	var total int64
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		dash := strings.IndexByte(p, '-')
		if dash == -1 {
			return 0, errors.New("invalid range: missing dash")
		}
		aStr := strings.TrimSpace(p[:dash])
		bStr := strings.TrimSpace(p[dash+1:])
		a, err := strconv.ParseInt(aStr, 10, 64)
		if err != nil {
			return 0, err
		}
		b, err := strconv.ParseInt(bStr, 10, 64)
		if err != nil {
			return 0, err
		}
		if a > b {
			a, b = b, a
		}
		total += f(a, b)
	}
	return total, nil
}

var pow10 = [...]int64{
	1,
	10,
	100,
	1000,
	10000,
	100000,
	1000000,
	10000000,
	100000000,
	1000000000,
	10000000000,
	100000000000,
	1000000000000,
	10000000000000,
	100000000000000,
	1000000000000000,
	10000000000000000,
	100000000000000000,
}

func digits(n int64) int {
	if n == 0 {
		return 1
	}
	d := 0
	for d < len(pow10) && n >= pow10[d] {
		d++
	}
	return d
}

// sumDoublesInRange returns the sum of all numbers in [lo, hi]
// that are formed by repeating a sequence of digits twice (e.g., 55, 6464, 123123).
func sumDoublesInRange(lo, hi int64) int64 {
	return sumRepeatedInRange(lo, hi, 2, 2)
}

// sumRepeatedAtLeastTwiceInRange sums numbers in [lo,hi] that are composed of a base block
// of k digits repeated m times for any m >= 2. Each number is counted once by requiring
// the base block to be primitive (i.e., not itself a repetition).
func sumRepeatedAtLeastTwiceInRange(lo, hi int64) int64 {
	return sumRepeatedInRange(lo, hi, 2, 0)
}

func divCeil(a, b int64) int64 {
	if b == 0 {
		return 0
	}
	if a >= 0 {
		return (a + b - 1) / b
	}
	// For negative a, integer division in Go rounds toward zero.
	// Ceil for negative a can be just a / b.
	return a / b
}

// sumRepeatedInRange sums numbers in [lo,hi] that are composed of a base block
// of k digits repeated m times, with m constrained by [minRepeats, maxRepeats].
// When maxRepeats==0, m ranges from minRepeats up to floor(digits(hi)/k), so
// total digits D = k*m never exceed digits(hi). Values are further clipped to
// [lo,hi] by intersecting the implied base range.
func sumRepeatedInRange(lo, hi int64, minRepeats, maxRepeats int) int64 {
	if hi < 11 {
		return 0
	}
	if lo < 0 {
		lo = 0
	}

	var sum int64
	maxDigits := digits(hi)
	// k is base block length
	for k := 1; k <= maxDigits/2; k++ {
		kLo := pow10[k-1]
		kHi := pow10[k] - 1

		maxM := maxDigits / k
		if maxRepeats > 0 && maxRepeats < maxM {
			maxM = maxRepeats
		}
		if maxM < minRepeats {
			continue
		}

		// m is repeat count in [minRepeats, maxM]
		for m := minRepeats; m <= maxM; m++ {
			// Compute factor = 1 + 10^k + 10^{2k} + ... + 10^{(m-1)k}
			km := k * m
			if km >= len(pow10) {
				// Out of precomputed range; safely stop as numbers would exceed our pow10 table
				break
			}
			denom := pow10[k] - 1
			if denom == 0 {
				continue
			}
			numer := pow10[km] - 1
			factor := numer / denom

			hLo := kLo
			hHi := kHi
			boundLo := divCeil(lo, factor)
			if boundLo > hLo {
				hLo = boundLo
			}
			boundHi := hi / factor
			if boundHi < hHi {
				hHi = boundHi
			}
			if hLo > hHi {
				continue
			}

			// Enumerate candidate base blocks and require primitive to avoid duplicates
			for h := hLo; h <= hHi; h++ {
				if !isPrimitiveBlock(h, k) {
					continue
				}
				sum += h * factor
			}
		}
	}

	return sum
}

// isPrimitiveBlock returns true if the k-digit number h is not composed of
// repeating a smaller digit block.
func isPrimitiveBlock(h int64, k int) bool {
	s := strconv.FormatInt(h, 10)
	if len(s) != k {
		if len(s) > k {
			return false
		}
		s = strings.Repeat("0", k-len(s)) + s
	}
	for d := 1; d <= k/2; d++ {
		if k%d != 0 {
			continue
		}
		base := s[:d]
		// If base starts with '0', then repeating would start with '0', not a valid k-digit block
		if base[0] == '0' {
			continue
		}
		repeat := k / d
		ok := true
		for i := 1; i < repeat; i++ {
			if s[i*d:(i+1)*d] != base {
				ok = false
				break
			}
		}
		if ok {
			return false
		}
	}
	return true
}
