package solve2025

import (
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day2 struct{}

func (d Day2) Coords() solve.SolutionCoords { return solve.SolutionCoords{Year: 2025, Day: 2} }

type giftShopRange struct{ low, high int64 }

func giftShopEndpoint(text string) (int64, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return 0, fmt.Errorf("empty endpoint")
	}
	for i := 0; i < len(text); i++ {
		if text[i] < '0' || text[i] > '9' {
			return 0, fmt.Errorf("expected nonnegative decimal endpoint")
		}
	}
	value, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("endpoint exceeds int64")
	}
	return value, nil
}

func parseGiftShopRanges(input string) ([]giftShopRange, error) {
	records := strings.Split(strings.TrimSpace(input), ",")
	ranges := make([]giftShopRange, 0, len(records))
	for index, record := range records {
		pair := strings.Split(record, "-")
		if len(pair) != 2 {
			return nil, fmt.Errorf("range %d: expected low-high", index+1)
		}
		low, err := giftShopEndpoint(pair[0])
		if err != nil {
			return nil, fmt.Errorf("range %d: %w", index+1, err)
		}
		high, err := giftShopEndpoint(pair[1])
		if err != nil {
			return nil, fmt.Errorf("range %d: %w", index+1, err)
		}
		if low > high {
			low, high = high, low
		}
		ranges = append(ranges, giftShopRange{low, high})
	}
	sort.Slice(ranges, func(i, j int) bool {
		if ranges[i].low != ranges[j].low {
			return ranges[i].low < ranges[j].low
		}
		return ranges[i].high < ranges[j].high
	})
	merged := ranges[:0]
	for _, next := range ranges {
		if len(merged) == 0 || next.low > merged[len(merged)-1].high {
			merged = append(merged, next)
		} else if next.high > merged[len(merged)-1].high {
			merged[len(merged)-1].high = next.high
		}
	}
	return merged, nil
}

func sumPeriodInRange(low, high int64, width, period int) *big.Int {
	power := int64(1)
	for i := 0; i < period; i++ {
		power *= 10
	}
	factor := int64(0)
	for i := 0; i < width/period; i++ {
		factor = factor*power + 1
	}
	first := low / factor
	if low%factor != 0 {
		first++
	}
	if first < power/10 {
		first = power / 10
	}
	last := high / factor
	if last >= power {
		last = power - 1
	}
	if first > last {
		return new(big.Int)
	}
	sum := big.NewInt(first + last)
	sum.Mul(sum, big.NewInt(last-first+1))
	sum.Quo(sum, big.NewInt(2))
	sum.Mul(sum, big.NewInt(factor))
	return sum
}

func sumDoublesInRange(low, high int64) *big.Int {
	total := new(big.Int)
	width := len(strconv.FormatInt(low, 10))
	if width%2 != 0 {
		width++
	}
	for ; width <= len(strconv.FormatInt(high, 10)); width += 2 {
		total.Add(total, sumPeriodInRange(low, high, width, width/2))
	}
	return total
}

func sumRepeatedAtLeastTwiceInRange(low, high int64) *big.Int {
	total := new(big.Int)
	for width := len(strconv.FormatInt(low, 10)); width <= len(strconv.FormatInt(high, 10)); width++ {
		primitive := make([]*big.Int, width/2+1)
		for period := 1; period <= width/2; period++ {
			if width%period != 0 {
				continue
			}
			sum := sumPeriodInRange(low, high, width, period)
			for smaller := 1; smaller < period; smaller++ {
				if period%smaller == 0 {
					sum.Sub(sum, primitive[smaller])
				}
			}
			primitive[period] = sum
			total.Add(total, sum)
		}
	}
	return total
}

func sumGiftShop(input string, anyRepeat bool) (string, error) {
	ranges, err := parseGiftShopRanges(input)
	if err != nil {
		return "", err
	}
	total := new(big.Int)
	for _, interval := range ranges {
		if anyRepeat {
			total.Add(total, sumRepeatedAtLeastTwiceInRange(interval.low, interval.high))
		} else {
			total.Add(total, sumDoublesInRange(interval.low, interval.high))
		}
	}
	return total.String(), nil
}
func (d Day2) Part1(input string) (string, error) { return sumGiftShop(input, false) }
func (d Day2) Part2(input string) (string, error) { return sumGiftShop(input, true) }
func init()                                       { solve.Register(Day2{}) }
