package solve2018

import (
	"fmt"
	"image"
	"sort"
	"strings"

	"aoc/solve"
)

type Day15 struct {
}

func (d Day15) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2018, Day: 15}
}

type Team int

const (
	ELF Team = iota
	GOBLIN
)

type Unit struct {
	Team  Team
	Pos   image.Point
	HP    int
	Alive bool
	Power int
}

type Grid struct {
	Walls map[image.Point]bool
	Units []*Unit
}

func parseGrid(input string, elfPower int) *Grid {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	walls := make(map[image.Point]bool)
	var units []*Unit
	for y, line := range lines {
		for x, ch := range line {
			pt := image.Pt(x, y)
			switch ch {
			case '#':
				walls[pt] = true
			case 'E', 'G':
				team := ELF
				power := 3
				if ch == 'E' {
					power = elfPower
				} else {
					team = GOBLIN
				}
				units = append(units, &Unit{Team: team, Pos: pt, HP: 200, Alive: true, Power: power})
			}
		}
	}
	return &Grid{Walls: walls, Units: units}
}

func (g *Grid) occupied(exclude *Unit) map[image.Point]bool {
	occ := make(map[image.Point]bool)
	for _, u := range g.Units {
		if u.Alive && u != exclude {
			occ[u.Pos] = true
		}
	}
	return occ
}

func (g *Grid) play(elfDeath bool) (int, bool) {
	rounds := 0
	for {
		sort.Slice(g.Units, func(i, j int) bool {
			a, b := g.Units[i], g.Units[j]
			return a.Pos.Y < b.Pos.Y || (a.Pos.Y == b.Pos.Y && a.Pos.X < b.Pos.X)
		})
		for _, u := range g.Units {
			if !u.Alive {
				continue
			}
			if g.move(u, elfDeath) {
				sum := 0
				for _, u2 := range g.Units {
					if u2.Alive {
						sum += u2.HP
					}
				}
				return rounds * sum, true
			}
		}
		rounds++
	}
}

func (g *Grid) move(u *Unit, elfDeath bool) bool {
	var targets []*Unit
	for _, t := range g.Units {
		if t.Alive && t.Team != u.Team {
			targets = append(targets, t)
		}
	}
	if len(targets) == 0 {
		return true
	}
	occ := g.occupied(u)
	inRange := map[image.Point]bool{}
	for _, t := range targets {
		for _, nb := range nb4(t.Pos) {
			if !g.Walls[nb] && !occ[nb] {
				inRange[nb] = true
			}
		}
	}
	if !inRange[u.Pos] {
		move := g.findMove(u.Pos, inRange, occ)
		if move != nil {
			u.Pos = *move
		}
	}
	// Attack
	var opponents []*Unit
	for _, t := range targets {
		for _, nb := range nb4(u.Pos) {
			if t.Pos == nb {
				opponents = append(opponents, t)
			}
		}
	}
	if len(opponents) > 0 {
		sort.Slice(opponents, func(i, j int) bool {
			if opponents[i].HP == opponents[j].HP {
				a, b := opponents[i].Pos, opponents[j].Pos
				return a.Y < b.Y || (a.Y == b.Y && a.X < b.X)
			}
			return opponents[i].HP < opponents[j].HP
		})
		target := opponents[0]
		target.HP -= u.Power
		if target.HP <= 0 {
			target.Alive = false
			if elfDeath && target.Team == ELF {
				panic("elf died")
			}
		}
	}
	return false
}

func (g *Grid) findMove(start image.Point, targets map[image.Point]bool, occ map[image.Point]bool) *image.Point {
	type node struct {
		Pt   image.Point
		Dist int
		Prev *image.Point
	}
	queue := []node{{start, 0, nil}}
	visited := map[image.Point]int{start: 0}
	parents := map[image.Point]*image.Point{}
	var found []node
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if targets[cur.Pt] {
			found = append(found, cur)
		}
		for _, nb := range nb4(cur.Pt) {
			if g.Walls[nb] || occ[nb] {
				continue
			}
			if d, ok := visited[nb]; ok && d <= cur.Dist+1 {
				continue
			}
			visited[nb] = cur.Dist + 1
			parents[nb] = &cur.Pt
			queue = append(queue, node{nb, cur.Dist + 1, &cur.Pt})
		}
	}
	if len(found) == 0 {
		return nil
	}
	sort.Slice(found, func(i, j int) bool {
		a, b := found[i].Pt, found[j].Pt
		if found[i].Dist == found[j].Dist {
			return a.Y < b.Y || (a.Y == b.Y && a.X < b.X)
		}
		return found[i].Dist < found[j].Dist
	})
	// Backtrack to first step
	dest := found[0].Pt
	for parents[dest] != nil && *parents[dest] != start {
		dest = *parents[dest]
	}
	return &dest
}

func nb4(p image.Point) []image.Point {
	return []image.Point{
		{p.X, p.Y - 1},
		{p.X - 1, p.Y},
		{p.X + 1, p.Y},
		{p.X, p.Y + 1},
	}
}

func (Day15) Part1(input string) (string, error) {
	grid := parseGrid(input, 3)
	outcome, _ := grid.play(false)
	return fmt.Sprintf("%d", outcome), nil
}

func (Day15) Part2(input string) (string, error) {
	for power := 4; ; power++ {
		var outcome int
		var allElves bool
		func() {
			defer func() {
				_ = recover()
			}()
			grid := parseGrid(input, power)
			outcome, _ = grid.play(true)
			allElves = true
			for _, u := range grid.Units {
				if u.Team == ELF && !u.Alive {
					allElves = false
					break
				}
			}
		}()
		if allElves {
			return fmt.Sprintf("%d", outcome), nil
		}
	}
}

func init() {
	solve.Register(Day15{})
}
