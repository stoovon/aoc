package solve2017

import (
	"fmt"
	"strings"

	"aoc/solve"
)

type Day16 struct {
}

func (d Day16) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2017, Day: 16}
}

func (d Day16) Part1(input string) (string, error) {
	programs := []byte("abcdefghijklmnop")
	moves := strings.Split(strings.TrimSpace(input), ",")
	err := applyDance(programs, moves)
	if err != nil {
		return "", fmt.Errorf("error applying dance: %v", err)
	}
	return string(programs), nil
}

func (d Day16) Part2(input string) (string, error) {
	programs := []byte("abcdefghijklmnop")
	moves := strings.Split(strings.TrimSpace(input), ",")
	seen := make(map[string]int)
	var order string
	iterations := 1000000000

	for i := 0; i < iterations; i++ {
		order = string(programs)
		if prev, ok := seen[order]; ok {
			// Cycle detected
			cycleLen := i - prev
			remaining := (iterations - i) % cycleLen
			for j := 0; j < remaining; j++ {
				err := applyDance(programs, moves)
				if err != nil {
					return "", fmt.Errorf("error applying dance: %v", err)
				}
			}
			return string(programs), nil
		}
		seen[order] = i
		err := applyDance(programs, moves)
		if err != nil {
			return "", fmt.Errorf("error applying dance: %v", err)
		}
	}
	return string(programs), nil
}

func applyDance(programs []byte, moves []string) error {
	for _, move := range moves {
		switch move[0] {
		case 's':
			var x int
			_, err := fmt.Sscanf(move[1:], "%d", &x)
			if err != nil {
				return fmt.Errorf("invalid spin move: %v", err)
			}
			n := len(programs)
			copy(programs, append(programs[n-x:], programs[:n-x]...))
		case 'x':
			var a, b int
			_, err := fmt.Sscanf(move[1:], "%d/%d", &a, &b)
			if err != nil {
				return fmt.Errorf("invalid exchange move: %v", err)
			}
			programs[a], programs[b] = programs[b], programs[a]
		case 'p':
			a, b := move[1], move[3]
			var ia, ib int
			for i, c := range programs {
				if c == a {
					ia = i
				}
				if c == b {
					ib = i
				}
			}
			programs[ia], programs[ib] = programs[ib], programs[ia]
		}
	}

	return nil
}

func init() {
	solve.Register(Day16{})
}
