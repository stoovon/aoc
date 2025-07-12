package solve2018

import (
	"aoc/solve"
	"fmt"
	"sort"
	"strings"
)

type Day13 struct{}

func (d Day13) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2018, Day: 13}
}

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

type TurnState int

const (
	TurnLeft TurnState = iota
	GoStraight
	TurnRight
)

type Cart struct {
	x, y      int
	direction Direction
	turnState TurnState
}

func (c *Cart) Move() {
	switch c.direction {
	case Up:
		c.y--
	case Down:
		c.y++
	case Left:
		c.x--
	case Right:
		c.x++
	}
}

func (c *Cart) Turn(track byte) {
	switch track {
	case '/':
		switch c.direction {
		case Up:
			c.direction = Right
		case Right:
			c.direction = Up
		case Down:
			c.direction = Left
		case Left:
			c.direction = Down
		}
	case '\\':
		switch c.direction {
		case Up:
			c.direction = Left
		case Left:
			c.direction = Up
		case Down:
			c.direction = Right
		case Right:
			c.direction = Down
		}
	case '+':
		switch c.turnState {
		case TurnLeft:
			c.direction = (c.direction + 3) % 4 // Turn left
		case GoStraight:
		case TurnRight:
			c.direction = (c.direction + 1) % 4 // Turn right
		}
		c.turnState = (c.turnState + 1) % 3
	}
}

func parseTrackAndCarts(input string) ([][]byte, []*Cart) {
	lines := strings.Split(strings.TrimRight(input, "\n"), "\n")

	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	track := make([][]byte, len(lines))
	var carts []*Cart

	for y, line := range lines {
		track[y] = make([]byte, maxWidth)
		for x := 0; x < maxWidth; x++ {
			if x < len(line) {
				char := line[x]
				switch char {
				case '^':
					carts = append(carts, &Cart{x: x, y: y, direction: Up, turnState: TurnLeft})
					track[y][x] = '|'
				case 'v':
					carts = append(carts, &Cart{x: x, y: y, direction: Down, turnState: TurnLeft})
					track[y][x] = '|'
				case '<':
					carts = append(carts, &Cart{x: x, y: y, direction: Left, turnState: TurnLeft})
					track[y][x] = '-'
				case '>':
					carts = append(carts, &Cart{x: x, y: y, direction: Right, turnState: TurnLeft})
					track[y][x] = '-'
				default:
					track[y][x] = char
				}
			} else {
				track[y][x] = ' '
			}
		}
	}

	return track, carts
}

func simulateUntilCrashOrLastCart(carts []*Cart, track [][]byte, stopOnFirstCrash bool) (string, error) {
	for {
		sort.Slice(carts, func(i, j int) bool {
			if carts[i].y != carts[j].y {
				return carts[i].y < carts[j].y
			}
			return carts[i].x < carts[j].x
		})

		crashed := make(map[*Cart]bool)

		for _, cart := range carts {
			if crashed[cart] {
				continue
			}

			cart.Move()

			for _, other := range carts {
				if cart != other && !crashed[other] && cart.x == other.x && cart.y == other.y {
					if stopOnFirstCrash {
						return fmt.Sprintf("%d,%d", cart.x, cart.y), nil
					}
					crashed[cart] = true
					crashed[other] = true
					break
				}
			}

			if !crashed[cart] {
				cart.Turn(track[cart.y][cart.x])
			}
		}

		if !stopOnFirstCrash {
			var survivingCarts []*Cart
			for _, cart := range carts {
				if !crashed[cart] {
					survivingCarts = append(survivingCarts, cart)
				}
			}
			carts = survivingCarts

			if len(carts) == 1 {
				return fmt.Sprintf("%d,%d", carts[0].x, carts[0].y), nil
			}

			if len(carts) == 0 {
				return "", fmt.Errorf("all carts crashed")
			}
		}
	}
}

func (d Day13) Part1(input string) (string, error) {
	track, carts := parseTrackAndCarts(input)
	return simulateUntilCrashOrLastCart(carts, track, true)
}

func (d Day13) Part2(input string) (string, error) {
	track, carts := parseTrackAndCarts(input)
	return simulateUntilCrashOrLastCart(carts, track, false)
}

func init() {
	solve.Register(Day13{})
}
