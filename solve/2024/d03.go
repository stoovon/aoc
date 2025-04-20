package solve2024

import (
	"fmt"
	"regexp"
	"strconv"

	"aoc/solve"
)

type Day3 struct {
}

func (d Day3) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2024, Day: 3}
}

type registers struct {
	mulTotal         int64
	mulTotalExecuted int64
}

func (d Day3) eval(input string) (registers, error) {
	mulRe := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)|(do)\(\)|(don't)\(\)`)
	mulMatches := mulRe.FindAllStringSubmatch(string(input), -1)

	var mulTotal, mulTotalState int64

	execute := true

	for _, match := range mulMatches {
		if match[3] == "do" {
			execute = true
		} else if match[4] == "don't" {
			execute = false
		} else {
			n1, err := strconv.Atoi(match[1])
			if err != nil {
				return registers{}, err
			}

			n2, err := strconv.Atoi(match[2])
			if err != nil {
				return registers{}, err
			}

			mul := int64(n1 * n2)
			mulTotal += mul

			if execute {
				mulTotalState += mul
			}
		}
	}

	return registers{
		mulTotal:         mulTotal,
		mulTotalExecuted: mulTotalState,
	}, nil
}

func (d Day3) Part1(input string) (string, error) {
	res, err := d.eval(input)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", res.mulTotal), nil
}

func (d Day3) Part2(input string) (string, error) {
	res, err := d.eval(input)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", res.mulTotalExecuted), nil
}

func init() {
	solve.Register(Day3{})
}
