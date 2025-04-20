package solve2024

import (
	"fmt"
	"image"
	"strings"

	"aoc/solve"
)

type Day4 struct {
}

func (d Day4) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2024, Day: 4}
}

func (d Day4) getGrid(input string) map[image.Point]rune {
	grid := map[image.Point]rune{}
	for y, s := range strings.Split(strings.TrimSpace(string(input)), "\n") {
		for x, r := range s {
			grid[image.Point{X: x, Y: y}] = r
		}
	}

	return grid
}

func (d Day4) Part1(input string) (string, error) {
	grid := d.getGrid(input)
	// Would also work with plainPermute (0, 1, 2, 3, 5, 6, 7, 8)
	permute := d.twistedPermute(grid)

	part1 := 0
	for p := range grid {
		part1 += strings.Count(strings.Join(permute(p, 4), " "), "XMAS")
	}

	return fmt.Sprintf("%d", part1), nil
}

func (d Day4) twistedPermute(grid map[image.Point]rune) func(p image.Point, length int) []string {
	return func(p image.Point, length int) []string {
		//      |  -1   |   0   |   1
		//  -1  |   0   |   3   |   6
		//  0   |   1   |   4   |   7
		//  1   |   2   |   5   |   8

		// Have to switch to 90s because we're looking at 2 MAS in shape of an X.
		// If we did 180s we'd also find MAS going vertically and horizontally.
		Chebyshev := []image.Point{
			{-1, -1}, // 0 90s
			{1, -1},  // 6
			{1, 1},   // 8
			{-1, 1},  // 2
			{0, -1},  // 3 90s
			{1, 0},   // 7
			{0, 1},   // 5
			{-1, 0},  // 1
		}

		combinations := make([]string, len(Chebyshev))
		for i, neighbourVector := range Chebyshev {
			for n := range length {

				// So for e.g. neighbourVector 2
				//      |  -2   |  -1   |   0  |  1  |  2  |
				//  -2  |   A   |   .   |   A  |  .  |  A  |
				//  -1  |   .   |   M   |   M  |  M  |  .  |
				//  0   |   A   |   M   |   X  |  M  |  A  |
				//  1   |   .   |   M   |   M  |  M  |  .  |
				//  2   |   A   |   .   |   A  |  .  |  A  |

				combinations[i] += string(grid[p.Add(neighbourVector.Mul(n))])
			}
		}
		return combinations
	}
}

func (d Day4) Part2(input string) (string, error) {
	grid := d.getGrid(input)
	permute := d.twistedPermute(grid)

	part2 := 0
	for p := range grid {
		// This one will need to be backwards, because otherwise things could get quite complicated
		// It's a bit like the sea monster

		// NeighbourVector 1
		//      |  -1   |   0  |  1  |
		//  -1  | (M|S) | (M|S) | (M|S) |
		//  0   | (M|S) |   A   | (M|S) |
		//  1   | (M|S) | (M|S) | (M|S) |

		// Combo A - The two 90s here cover both combos
		//      | -1  |  0  |  1  |
		//  -1  |  M  |  M  |  M  |
		//  0   |  M  |  A  |  S  |
		//  1   |  S  |  S  |  S  |

		// Combo B - The two 90s here cover both combos inverted
		//      | -1  |  0  |  1  |
		//  -1  |  S  |  S  |  S  |
		//  0   |  S  |  A  |  M  |
		//  1   |  M  |  M  |  M  |

		//      |  -1   |   0   |   1
		//  -1  |   0   |   3   |   6
		//  0   |   1   |   4   |   7
		//  1   |   2   |   5   |   8

		//         | (0,6)    | (8,2)    | (3,7)    | (5,1)    |
		// Combo A | (AM, AM) | (AS, AS) | (AM, AS) | (AS, AM) | 90s
		// Combo B | (AS, AS) | (AM, AM) | (AS, AM) | (AM, AS) | 90s

		// AMAMASAS       Combo A, First Pair  90
		//   AMASASAM     Combo A, Second Pair 90
		//     ASASAMAM   Combo B, First Pair  90
		//       ASAMAMAS Combo B, Second Pair 90
		// AMAMASASAMAMAS

		part2 += strings.Count("AMAMASASAMAMAS", strings.Join(permute(p, 2)[:4], ""))
	}

	return fmt.Sprintf("%d", part2), nil
}

func init() {
	solve.Register(Day4{})
}
