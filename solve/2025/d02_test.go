package solve2025

import (
	"math/big"
	"strconv"
	"strings"
	"testing"
)

func giftShopDirectInvalid(value int, anyRepeat bool) bool {
	text := strconv.Itoa(value)
	if !anyRepeat {
		half := len(text) / 2
		return len(text)%2 == 0 && text[:half] == text[half:]
	}
	for width := 1; width <= len(text)/2; width++ {
		if len(text)%width == 0 && strings.Repeat(text[:width], len(text)/width) == text {
			return true
		}
	}
	return false
}

func TestGiftShopPart1AllowsRepetitiveHalves(t *testing.T) {
	for _, value := range []int64{73737373, 48484848} {
		input := strconv.FormatInt(value, 10) + "-" + strconv.FormatInt(value, 10)
		for part, solver := range []func(string) (string, error){(Day2{}).Part1, (Day2{}).Part2} {
			got, err := solver(input)
			if err != nil || got != strconv.FormatInt(value, 10) {
				t.Fatalf("part %d repetitive half: got %q %v want %d", part+1, got, err, value)
			}
		}
	}
}

func TestGiftShopRepeatRulesAgainstDirectStrings(t *testing.T) {
	var totalDoubles, totalRepeated int64
	for value := 0; value < 100000; value++ {
		var double, repeated int64
		if giftShopDirectInvalid(value, false) {
			double = int64(value)
		}
		if giftShopDirectInvalid(value, true) {
			repeated = int64(value)
		}
		if got := sumDoublesInRange(int64(value), int64(value)); got.Cmp(big.NewInt(double)) != 0 {
			t.Fatalf("double %d: got %s want %d", value, got.String(), double)
		}
		if got := sumRepeatedAtLeastTwiceInRange(int64(value), int64(value)); got.Cmp(big.NewInt(repeated)) != 0 {
			t.Fatalf("repeated %d: got %s want %d", value, got.String(), repeated)
		}
		totalDoubles += double
		totalRepeated += repeated
	}
	if got := sumDoublesInRange(0, 99999); got.Cmp(big.NewInt(totalDoubles)) != 0 {
		t.Fatalf("double range: %s want %d", got.String(), totalDoubles)
	}
	if got := sumRepeatedAtLeastTwiceInRange(0, 99999); got.Cmp(big.NewInt(totalRepeated)) != 0 {
		t.Fatalf("repeated range: %s want %d", got.String(), totalRepeated)
	}
}

func TestGiftShopPart1WideSingleton(t *testing.T) {
	const value int64 = 737373737737373737
	got := sumDoublesInRange(value, value)
	if got.Cmp(big.NewInt(value)) != 0 {
		t.Fatalf("18-digit double: %s", got)
	}
}

func TestGiftShopGeneratedPeriods(t *testing.T) {
	for width := 2; width <= 8; width++ {
		seen := make(map[int64]bool)
		for period := 1; period <= width/2; period++ {
			if width%period != 0 {
				continue
			}
			low := 1
			for i := 1; i < period; i++ {
				low *= 10
			}
			for base := low; base < low*10; base++ {
				value, err := strconv.ParseInt(strings.Repeat(strconv.Itoa(base), width/period), 10, 64)
				if err != nil {
					t.Fatal(err)
				}
				seen[value] = true
			}
		}
		var low int64 = 1
		for i := 1; i < width; i++ {
			low *= 10
		}
		for _, interval := range [][2]int64{{low, low*10 - 1}, {low + low/3, low*7 + 19}, {low * 4, low*4 + 17003}} {
			if interval[1] >= low*10 {
				interval[1] = low*10 - 1
			}
			var want int64
			for value := range seen {
				if value >= interval[0] && value <= interval[1] {
					want += value
				}
			}
			got := sumRepeatedAtLeastTwiceInRange(interval[0], interval[1])
			if got.Cmp(big.NewInt(want)) != 0 {
				t.Fatalf("width %d interval %v: %s want %d", width, interval, got, want)
			}
		}
	}
}

func TestGiftShopNumericLimits(t *testing.T) {
	// Independent prime-factor inclusion-exclusion gives these full-domain totals.
	for part, tc := range []struct {
		fn   func(int64, int64) *big.Int
		want string
	}{
		{sumDoublesInRange, "495495495540950040450040950"},
		{sumRepeatedAtLeastTwiceInRange, "495990091040401571498681796"},
	} {
		if got := tc.fn(0, 9223372036854775807).String(); got != tc.want {
			t.Fatalf("part %d full domain: %s", part+1, got)
		}
		if got := tc.fn(9223372036854775807, 9223372036854775807); got.Sign() != 0 {
			t.Fatalf("part %d maximum: %s", part+1, got)
		}
	}
	const repeated int64 = 8888888888888888888
	if got := sumRepeatedAtLeastTwiceInRange(repeated, repeated); got.Cmp(big.NewInt(repeated)) != 0 {
		t.Fatalf("19-digit repeat: %s", got)
	}
	if got := sumDoublesInRange(repeated, repeated); got.Sign() != 0 {
		t.Fatalf("odd-width double: %s", got)
	}
}

func TestGiftShopRangeParsing(t *testing.T) {
	for _, input := range []string{"", " ", "1", "1-2-3", "-1-2", "+1-2", "1-", "1-2,", ",1-2", "1-2,,3-4", "1-x", "0-9223372036854775808"} {
		for _, fn := range []func(string) (string, error){(Day2{}).Part1, (Day2{}).Part2} {
			if _, err := fn(input); err == nil {
				t.Fatalf("accepted malformed range %q", input)
			}
		}
	}
	for _, fn := range []func(string) (string, error){(Day2{}).Part1, (Day2{}).Part2} {
		got, err := fn(" 0044 - 0022 , 33-55, 22-22 \n")
		if err != nil || got != "154" {
			t.Fatalf("normalized overlapping ranges: %q %v", got, err)
		}
	}
}
