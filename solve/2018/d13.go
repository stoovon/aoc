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

// Direction represents the direction a cart is facing
type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

// TurnState represents the state of intersection turns for a cart
type TurnState int

const (
	TurnLeft TurnState = iota
	GoStraight
	TurnRight
)

// Cart represents a cart on the track
type Cart struct {
	x, y      int
	direction Direction
	turnState TurnState
}

// Move the cart one step forward
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

// Turn the cart based on the track piece it's on
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
		// Intersection - turn based on turn state
		switch c.turnState {
		case TurnLeft:
			c.direction = (c.direction + 3) % 4 // Turn left
		case GoStraight:
			// No change in direction
		case TurnRight:
			c.direction = (c.direction + 1) % 4 // Turn right
		}
		c.turnState = (c.turnState + 1) % 3 // Cycle through turn states
	}
}

func parseTrackAndCarts(input string) ([][]byte, []*Cart) {
	lines := strings.Split(strings.TrimRight(input, "\n"), "\n")

	// Find the maximum width
	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	// Create the track grid
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
					track[y][x] = '|' // Replace cart with underlying track
				case 'v':
					carts = append(carts, &Cart{x: x, y: y, direction: Down, turnState: TurnLeft})
					track[y][x] = '|' // Replace cart with underlying track
				case '<':
					carts = append(carts, &Cart{x: x, y: y, direction: Left, turnState: TurnLeft})
					track[y][x] = '-' // Replace cart with underlying track
				case '>':
					carts = append(carts, &Cart{x: x, y: y, direction: Right, turnState: TurnLeft})
					track[y][x] = '-' // Replace cart with underlying track
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
		// Sort carts by position (top to bottom, left to right)
		sort.Slice(carts, func(i, j int) bool {
			if carts[i].y != carts[j].y {
				return carts[i].y < carts[j].y
			}
			return carts[i].x < carts[j].x
		})

		// Track which carts have crashed in this tick
		crashed := make(map[*Cart]bool)

		// Move each cart
		for _, cart := range carts {
			if crashed[cart] {
				continue // Skip carts that have already crashed
			}

			cart.Move()

			// Check for collision with other non-crashed carts
			for _, other := range carts {
				if cart != other && !crashed[other] && cart.x == other.x && cart.y == other.y {
					if stopOnFirstCrash {
						return fmt.Sprintf("%d,%d", cart.x, cart.y), nil
					}
					// Mark both carts as crashed
					crashed[cart] = true
					crashed[other] = true
					break
				}
			}

			// Turn the cart based on the track piece (only if not crashed)
			if !crashed[cart] {
				cart.Turn(track[cart.y][cart.x])
			}
		}

		if !stopOnFirstCrash {
			// Remove crashed carts
			var survivingCarts []*Cart
			for _, cart := range carts {
				if !crashed[cart] {
					survivingCarts = append(survivingCarts, cart)
				}
			}
			carts = survivingCarts

			// Check if only one cart remains
			if len(carts) == 1 {
				return fmt.Sprintf("%d,%d", carts[0].x, carts[0].y), nil
			}

			// If no carts remain, something went wrong
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
