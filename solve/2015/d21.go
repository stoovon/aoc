package solve2015

import (
	"math"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day21 struct {
}

func (d Day21) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2015, Day: 21}
}

type item struct {
	name   string
	cost   int
	damage int
	armor  int
}

type player struct {
	hp     int
	damage int
	armor  int
}

func simulateBattle(player, boss player) bool {
	for {
		// Player attacks
		boss.hp -= max(1, player.damage-boss.armor)
		if boss.hp <= 0 {
			return true
		}

		// Boss attacks
		player.hp -= max(1, boss.damage-player.armor)
		if player.hp <= 0 {
			return false
		}
	}
}

func (d Day21) Part1(input string) (string, error) {
	bossStats := parseBossStats(input)

	weapons := []item{
		{"Dagger", 8, 4, 0},
		{"Shortsword", 10, 5, 0},
		{"Warhammer", 25, 6, 0},
		{"Longsword", 40, 7, 0},
		{"Greataxe", 74, 8, 0},
	}

	armor := []item{
		{"Leather", 13, 0, 1},
		{"Chainmail", 31, 0, 2},
		{"Splintmail", 53, 0, 3},
		{"Bandedmail", 75, 0, 4},
		{"Platemail", 102, 0, 5},
		{"Unarmoured", 0, 0, 0},
	}

	rings := []item{
		{"Damage +1", 25, 1, 0},
		{"Damage +2", 50, 2, 0},
		{"Damage +3", 100, 3, 0},
		{"Defense +1", 20, 0, 1},
		{"Defense +2", 40, 0, 2},
		{"Defense +3", 80, 0, 3},
		{"Left Empty", 0, 0, 0},
		{"Right Empty", 0, 0, 0},
	}

	minCost := math.MaxInt
	for _, w := range weapons {
		for _, a := range armor {
			for i := 0; i < len(rings); i++ {
				for j := i + 1; j < len(rings); j++ {
					cost := w.cost + a.cost + rings[i].cost + rings[j].cost
					playerCharacter := player{
						hp:     100,
						damage: w.damage + a.damage + rings[i].damage + rings[j].damage,
						armor:  w.armor + a.armor + rings[i].armor + rings[j].armor,
					}
					boss := player{hp: bossStats[0], damage: bossStats[1], armor: bossStats[2]}
					if simulateBattle(playerCharacter, boss) {
						if cost < minCost {
							minCost = cost
						}
					}
				}
			}
		}
	}

	return strconv.Itoa(minCost), nil
}

func (d Day21) Part2(input string) (string, error) {
	bossStats := parseBossStats(input)

	weapons := []item{
		{"Dagger", 8, 4, 0},
		{"Shortsword", 10, 5, 0},
		{"Warhammer", 25, 6, 0},
		{"Longsword", 40, 7, 0},
		{"Greataxe", 74, 8, 0},
	}

	armor := []item{
		{"Leather", 13, 0, 1},
		{"Chainmail", 31, 0, 2},
		{"Splintmail", 53, 0, 3},
		{"Bandedmail", 75, 0, 4},
		{"Platemail", 102, 0, 5},
		{"Unarmoured", 0, 0, 0},
	}

	rings := []item{
		{"Damage +1", 25, 1, 0},
		{"Damage +2", 50, 2, 0},
		{"Damage +3", 100, 3, 0},
		{"Defense +1", 20, 0, 1},
		{"Defense +2", 40, 0, 2},
		{"Defense +3", 80, 0, 3},
		{"Left Empty", 0, 0, 0},
		{"Right Empty", 0, 0, 0},
	}

	maxCost := 0
	for _, w := range weapons {
		for _, a := range armor {
			for i := 0; i < len(rings); i++ {
				for j := i + 1; j < len(rings); j++ {
					cost := w.cost + a.cost + rings[i].cost + rings[j].cost
					playerCharacter := player{
						hp:     100,
						damage: w.damage + a.damage + rings[i].damage + rings[j].damage,
						armor:  w.armor + a.armor + rings[i].armor + rings[j].armor,
					}
					boss := player{hp: bossStats[0], damage: bossStats[1], armor: bossStats[2]}
					if !simulateBattle(playerCharacter, boss) {
						if cost > maxCost {
							maxCost = cost
						}
					}
				}
			}
		}
	}

	return strconv.Itoa(maxCost), nil
}

func parseBossStats(input string) []int {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	stats := make([]int, len(lines))
	for i, line := range lines {
		parts := strings.Split(line, ": ")
		stats[i], _ = strconv.Atoi(parts[1])
	}
	return stats
}

func init() {
	solve.Register(Day21{})
}
