package solve2019

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day22 struct {
}

func (d Day22) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2019, Day: 22}
}

func modinv(a, m *big.Int) *big.Int {
	// Modular inverse using Fermat's little theorem (m is prime)
	return new(big.Int).Exp(a, new(big.Int).Sub(m, big.NewInt(2)), m)
}

func parseBigInt(s string) *big.Int {
	n := new(big.Int)
	n.SetString(s, 10)
	return n
}

func ints(line string) int64 {
	// Extracts the first integer from a string
	fields := strings.Fields(line)
	for _, f := range fields {
		if n, err := strconv.ParseInt(f, 10, 64); err == nil {
			return n
		}
	}
	return 0
}

func (d Day22) Part1(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	cards := 10007
	deck := make([]int, cards)
	for i := range deck {
		deck[i] = i
	}

	for _, line := range lines {
		switch {
		case line == "deal into new stack":
			// Reverse the deck
			for i, j := 0, len(deck)-1; i < j; i, j = i+1, j-1 {
				deck[i], deck[j] = deck[j], deck[i]
			}
		case strings.HasPrefix(line, "cut"):
			parts := strings.Fields(line)
			q, _ := strconv.Atoi(parts[len(parts)-1])
			if q < 0 {
				q = cards + q
			}
			q = q % cards // handle q > cards
			deck = append(deck[q:], deck[:q]...)
		case strings.HasPrefix(line, "deal with increment"):
			parts := strings.Fields(line)
			q, _ := strconv.Atoi(parts[len(parts)-1])
			newDeck := make([]int, cards)
			pos := 0
			for _, card := range deck {
				newDeck[pos] = card
				pos = (pos + q) % cards
			}
			deck = newDeck
		}
	}

	for i, v := range deck {
		if v == 2019 {
			return fmt.Sprintf("%d", i), nil
		}
	}
	return "", fmt.Errorf("card 2019 not found")
}

func (d Day22) Part2(input string) (string, error) {
	cards := parseBigInt("119315717514047")
	repeats := parseBigInt("101741582076661")

	incrementMul := big.NewInt(1)
	offsetDiff := big.NewInt(0)

	lines := strings.Split(strings.TrimSpace(input), "\n")
	for _, line := range lines {
		switch {
		case line == "deal into new stack":
			// incrementMul *= -1; offsetDiff += incrementMul
			incrementMul.Neg(incrementMul)
			incrementMul.Mod(incrementMul, cards)
			offsetDiff.Add(offsetDiff, incrementMul)
			offsetDiff.Mod(offsetDiff, cards)
		case strings.HasPrefix(line, "cut"):
			q := ints(line)
			// offsetDiff += q * incrementMul
			tmp := new(big.Int).Mul(big.NewInt(q), incrementMul)
			offsetDiff.Add(offsetDiff, tmp)
			offsetDiff.Mod(offsetDiff, cards)
		case strings.HasPrefix(line, "deal with increment"):
			q := ints(line)
			// incrementMul *= inv(q)
			invQ := modinv(big.NewInt(q), cards)
			incrementMul.Mul(incrementMul, invQ)
			incrementMul.Mod(incrementMul, cards)
		}
	}

	// increment = incrementMul^repeats mod cards
	increment := new(big.Int).Exp(incrementMul, repeats, cards)

	// offset = offsetDiff * (1 - increment) * inv(1 - incrementMul) mod cards
	one := big.NewInt(1)
	oneMinusInc := new(big.Int).Sub(one, increment)
	oneMinusInc.Mod(oneMinusInc, cards)
	oneMinusIncMul := new(big.Int).Sub(one, incrementMul)
	oneMinusIncMul.Mod(oneMinusIncMul, cards)
	invOneMinusIncMul := modinv(oneMinusIncMul, cards)
	offset := new(big.Int).Mul(offsetDiff, oneMinusInc)
	offset.Mul(offset, invOneMinusIncMul)
	offset.Mod(offset, cards)

	// result = (offset + 2020 * increment) % cards
	pos := big.NewInt(2020)
	result := new(big.Int).Mul(pos, increment)
	result.Add(result, offset)
	result.Mod(result, cards)

	return result.String(), nil
}

func init() {
	solve.Register(Day22{})
}
