package solve2022

import (
	"aoc/solve"
	"errors"
	"strings"
)

type Day25 struct{}

func (d Day25) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2022, Day: 25}
}

func (d Day25) Part1(input string) (string, error) {
	// Helper function to convert SNAFU to decimal
	snafuToDecimal := func(snafu string) int {
		value := 0
		place := 1
		for i := len(snafu) - 1; i >= 0; i-- {
			switch snafu[i] {
			case '2':
				value += 2 * place
			case '1':
				value += 1 * place
			case '0':
				// No change needed for 0
			case '-':
				value -= 1 * place
			case '=':
				value -= 2 * place
			}
			place *= 5
		}
		return value
	}

	// Helper function to convert decimal to SNAFU
	decimalToSnafu := func(decimal int) string {
		if decimal == 0 {
			return "0"
		}
		snafu := ""
		for decimal != 0 {
			remainder := decimal % 5
			decimal /= 5
			switch remainder {
			case 0:
				snafu = "0" + snafu
			case 1:
				snafu = "1" + snafu
			case 2:
				snafu = "2" + snafu
			case 3:
				snafu = "=" + snafu
				decimal++
			case 4:
				snafu = "-" + snafu
				decimal++
			}
		}
		return snafu
	}

	// Parse input and calculate the sum in decimal
	sum := 0
	for _, line := range strings.Fields(input) {
		sum += snafuToDecimal(line)
	}

	// Convert the sum back to SNAFU
	result := decimalToSnafu(sum)
	return result, nil
}

func (d Day25) Part2(input string) (string, error) {
	return "", errors.New("not implemented")
}

func init() {
	solve.Register(Day25{})
}
