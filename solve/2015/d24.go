package solve2015

import (
	"cmp"
	"fmt"
	"math/big"
	"math/bits"
	"slices"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day24 struct{}

func (d Day24) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2015, Day: 24}
}

type packagePartitionKey struct {
	remaining uint64
	groups    int
}

func balancePackages(weights []int, groups int) (*big.Int, error) {
	if groups < 2 || groups > len(weights) {
		return nil, fmt.Errorf("expected at least two nonempty groups")
	}
	if len(weights) > 63 {
		return nil, fmt.Errorf("package search supports at most 63 packages")
	}
	weights = slices.Clone(weights)
	slices.SortFunc(weights, func(a, b int) int { return cmp.Compare(b, a) })
	prefix := make([]int, len(weights)+1)
	maxInt := int(^uint(0) >> 1)
	for index, weight := range weights {
		if weight <= 0 {
			return nil, fmt.Errorf("package weights must be positive")
		}
		if weight > maxInt-prefix[index] {
			return nil, fmt.Errorf("total package weight overflows int")
		}
		prefix[index+1] = prefix[index] + weight
	}
	total := prefix[len(weights)]
	if total%groups != 0 {
		return nil, fmt.Errorf("total package weight is not divisible by group count")
	}
	target := total / groups
	if weights[0] > target {
		return nil, fmt.Errorf("a package exceeds the group weight")
	}
	all := uint64(1)<<len(weights) - 1
	memo := make(map[packagePartitionKey]bool)
	var canPartition func(uint64, int) (bool, error)
	canPartition = func(remaining uint64, groupsLeft int) (bool, error) {
		// Every completed group removes target weight, so the final group is forced.
		if groupsLeft == 1 {
			return true, nil
		}
		key := packagePartitionKey{remaining: remaining, groups: groupsLeft}
		if possible, seen := memo[key]; seen {
			return possible, nil
		}
		// Groups are interchangeable: fix one remaining package in the next group.
		anchor := bits.TrailingZeros64(remaining)
		anchorBit := uint64(1) << anchor
		var fill func(int, int, uint64) (bool, error)
		fill = func(start, needed int, chosen uint64) (bool, error) {
			if needed == 0 {
				return canPartition(remaining&^chosen, groupsLeft-1)
			}
			previousWeight := -1
			for index := start; index < len(weights); index++ {
				bit := uint64(1) << index
				weight := weights[index]
				if remaining&bit == 0 || weight > needed || weight == previousWeight {
					continue
				}
				previousWeight = weight
				possible, err := fill(index+1, needed-weight, chosen|bit)
				if err != nil || possible {
					return possible, err
				}
			}
			return false, nil
		}
		possible, err := fill(anchor+1, target-weights[anchor], anchorBit)
		if err != nil {
			return false, err
		}
		if len(memo) >= 1<<20 {
			return false, fmt.Errorf("package partition exceeds the memo table budget")
		}
		memo[key] = possible
		return possible, nil
	}

	var best *big.Int
	var chooseFirst func(int, int, int, uint64, *big.Int) error
	chooseFirst = func(start, slots, needed int, chosen uint64, product *big.Int) error {
		if slots == 0 {
			if needed != 0 || best != nil && product.Cmp(best) >= 0 {
				return nil
			}
			possible, err := canPartition(all&^chosen, groups-1)
			if err != nil {
				return err
			}
			if possible {
				best = new(big.Int).Set(product)
			}
			return nil
		}
		if len(weights)-start < slots || needed <= 0 {
			return nil
		}
		if needed > prefix[start+slots]-prefix[start] || needed < prefix[len(weights)]-prefix[len(weights)-slots] {
			return nil
		}
		for index := start; index <= len(weights)-slots; index++ {
			weight := weights[index]
			if weight > needed || index > start && weight == weights[index-1] {
				continue
			}
			nextProduct := new(big.Int).Mul(product, big.NewInt(int64(weight)))
			if best != nil && nextProduct.Cmp(best) >= 0 {
				continue
			}
			if err := chooseFirst(index+1, slots-1, needed-weight, chosen|uint64(1)<<index, nextProduct); err != nil {
				return err
			}
		}
		return nil
	}
	// Some group in any partition has at most the average number of packages.
	for count := 1; count <= len(weights)/groups; count++ {
		if err := chooseFirst(0, count, target, 0, big.NewInt(1)); err != nil {
			return nil, err
		}
		if best != nil {
			return best, nil
		}
	}
	return nil, fmt.Errorf("packages cannot form equal-weight groups")
}

func (d Day24) solve(input string, groups int) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", fmt.Errorf("package list is empty")
	}
	var weights []int
	for index, line := range strings.Split(input, "\n") {
		weight, err := strconv.Atoi(strings.TrimSpace(line))
		if err != nil {
			return "", fmt.Errorf("line %d: invalid package weight: %w", index+1, err)
		}
		if weight <= 0 {
			return "", fmt.Errorf("line %d: package weight must be positive", index+1)
		}
		weights = append(weights, weight)
	}
	result, err := balancePackages(weights, groups)
	if err != nil {
		return "", err
	}
	return result.String(), nil
}

func (d Day24) Part1(input string) (string, error) {
	return d.solve(input, 3)
}

func (d Day24) Part2(input string) (string, error) {
	return d.solve(input, 4)
}

func init() {
	solve.Register(Day24{})
}
