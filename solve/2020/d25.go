package solve2020

import (
	"aoc/solve"
	"errors"
	"strings"
	"strconv"
)

type Day25 struct{}

func (d Day25) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2020, Day: 25}
}

func (d Day25) Part1(input string) (string, error) {
	lines := strings.Fields(strings.TrimSpace(input))
	if len(lines) < 2 {
		return "", errors.New("need two public keys")
	}
	cardPub, err1 := strconv.Atoi(lines[0])
	doorPub, err2 := strconv.Atoi(lines[1])
	if err1 != nil || err2 != nil {
		return "", errors.New("invalid public keys")
	}
	subject := 7
	mod := 20201227
	loop := babyStepGiantStep(subject, cardPub, mod)
	if loop == -1 {
		return "", errors.New("no discrete log found")
	}
	enc := powMod(doorPub, loop, mod)
	return strconv.Itoa(enc), nil
}

// babyStepGiantStep solves for x in base^x ≡ target (mod mod), returns -1 if not found.
func babyStepGiantStep(base, target, mod int) int {
	m := intSqrt(mod) + 1
	table := make(map[int]int, m)
	e := 1
	for j := 0; j < m; j++ {
		table[e] = j
		e = (e * base) % mod
	}
	inv := powMod(base, mod-2, mod) // Fermat's little theorem for inverse
	factor := powMod(inv, m, mod)
	e = target
	for i := 0; i < m; i++ {
		if j, ok := table[e]; ok {
			return i*m + j
		}
		e = (e * factor) % mod
	}
	return -1
}

func powMod(a, b, mod int) int {
	res := 1
	for b > 0 {
		if b&1 == 1 {
			res = (res * a) % mod
		}
		a = (a * a) % mod
		b >>= 1
	}
	return res
}

func intSqrt(n int) int {
	x := n
	y := (x + 1) / 2
	for y < x {
		x = y
		y = (x + n/x) / 2
	}
	return x
}

func (d Day25) Part2(input string) (string, error) {
	return "", errors.New("Mulligan")
}

func init() {
	solve.Register(Day25{})
}
