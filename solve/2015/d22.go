package solve2015

import (
	"container/heap"
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/queues"
)

type Day22 struct {
}

func (d Day22) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2015, Day: 22}
}

type spell struct {
	name    string
	cost    int
	effect  bool
	turns   int
	damage  int
	healing int
	armor   int
	mana    int
}

var spells = []spell{
	{"Magic Missile", 53, false, 0, 4, 0, 0, 0},
	{"Drain", 73, false, 0, 2, 2, 0, 0},
	{"Shield", 113, true, 6, 0, 0, 7, 0},
	{"Poison", 173, true, 6, 3, 0, 0, 0},
	{"Recharge", 229, true, 5, 0, 0, 0, 101},
}

type state struct {
	hp         int
	mana       int
	bossHP     int
	bossDamage int
	manaSpent  int
	effects    []effect
	hardMode   bool
	parent     *state
	spellCast  string
}

type effect struct {
	turns int
	spell spell
}

func (s *state) processEffects() (int, int, int, []effect) {
	hp, mana, bossHP := s.hp, s.mana, s.bossHP
	armor := 0
	var remainingEffects []effect

	for _, currentEffect := range s.effects {
		hp += currentEffect.spell.healing
		mana += currentEffect.spell.mana
		bossHP -= currentEffect.spell.damage
		if currentEffect.spell.armor > armor {
			armor = currentEffect.spell.armor
		}
		if currentEffect.turns > 1 {
			remainingEffects = append(remainingEffects, effect{currentEffect.turns - 1, currentEffect.spell})
		}
	}

	return hp, mana, bossHP, remainingEffects
}

func (s *state) bossTurn() {
	hp, mana, bossHP, effects := s.processEffects()
	s.hp, s.mana, s.bossHP, s.effects = hp, mana, bossHP, effects
	if s.bossHP > 0 {
		damage := max(1, s.bossDamage-armorFromEffects(effects))
		s.hp -= damage
	}
}

func armorFromEffects(effects []effect) int {
	armor := 0
	for _, effect := range effects {
		if effect.spell.armor > armor {
			armor = effect.spell.armor
		}
	}
	return armor
}

func (s *state) transitions() []*state {
	hp, mana, bossHP, effects := s.processEffects()
	if s.hardMode {
		hp--
	}
	if hp <= 0 {
		return nil
	}

	var nextStates []*state
	for _, spell := range spells {
		if spell.cost > mana || spellsInEffect(spell, effects) {
			continue
		}

		newState := &state{
			hp:         hp,
			mana:       mana - spell.cost,
			bossHP:     bossHP,
			bossDamage: s.bossDamage,
			manaSpent:  s.manaSpent + spell.cost,
			effects:    append([]effect{}, effects...),
			hardMode:   s.hardMode,
			parent:     s,
			spellCast:  spell.name,
		}

		if spell.effect {
			newState.effects = append(newState.effects, effect{spell.turns, spell})
		} else {
			newState.hp += spell.healing
			newState.bossHP -= spell.damage
		}

		newState.bossTurn()
		if newState.hp > 0 {
			nextStates = append(nextStates, newState)
		}
	}

	return nextStates
}

func spellsInEffect(spell spell, effects []effect) bool {
	for _, effect := range effects {
		if effect.spell.name == spell.name {
			return true
		}
	}
	return false
}

type pqItem struct {
	value    *state
	priority int
}

func (p *pqItem) Score() int {
	return p.priority
}

func searchAStar(start *state) *state {
	openSet := &queues.PriorityQueue[*pqItem]{}
	heap.Init(openSet)
	heap.Push(openSet, &pqItem{value: start, priority: start.manaSpent})
	closedSet := make(map[*state]bool)

	for openSet.Len() > 0 {
		current := heap.Pop(openSet).(*pqItem).value
		if current.bossHP <= 0 {
			return current
		}

		closedSet[current] = true
		for _, next := range current.transitions() {
			if closedSet[next] {
				continue
			}
			heap.Push(openSet, &pqItem{value: next, priority: next.manaSpent})
		}
	}

	return nil
}

func (d Day22) parseInput(input string) (int, int) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	bossHP, _ := strconv.Atoi(strings.Split(lines[0], ": ")[1])
	bossDamage, _ := strconv.Atoi(strings.Split(lines[1], ": ")[1])
	return bossHP, bossDamage
}

func (d Day22) Part1(input string) (string, error) {
	bossHP, bossDamage := d.parseInput(input)
	start := &state{hp: 50, mana: 500, bossHP: bossHP, bossDamage: bossDamage}
	end := searchAStar(start)
	return strconv.Itoa(end.manaSpent), nil
}

func (d Day22) Part2(input string) (string, error) {
	bossHP, bossDamage := d.parseInput(input)
	start := &state{hp: 50, mana: 500, bossHP: bossHP, bossDamage: bossDamage, hardMode: true}
	end := searchAStar(start)
	return strconv.Itoa(end.manaSpent), nil
}

func init() {
	solve.Register(Day22{})
}
