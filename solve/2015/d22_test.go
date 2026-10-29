package solve2015

import (
	"container/heap"
	"slices"
	"testing"

	"aoc/utils/queues"
)

func TestPoisonWinAtPlayerTurnStartCostsNoAdditionalMana(t *testing.T) {
	for _, mana := range []int{0, 53} {
		start := &state{
			hp: 10, mana: mana, bossHP: 3, bossDamage: 8, manaSpent: 173,
			effects: []effect{{turns: 1, spell: spells[3]}},
		}
		winner := searchMinMana(start)
		if winner == nil {
			t.Fatalf("mana=%d: expected poison to win without a spell", mana)
		}
		if winner.manaSpent != start.manaSpent || winner.mana != mana || winner.bossHP > 0 {
			t.Fatalf("mana=%d: poison win charged mana or left boss alive: %+v", mana, winner)
		}
		if start.bossHP != 3 || start.effects[0].turns != 1 {
			t.Fatal("search mutated the starting state")
		}
	}
}

func TestHardModeDeathPrecedesPoisonWin(t *testing.T) {
	start := &state{
		hp: 1, mana: 500, bossHP: 3, bossDamage: 8, hardMode: true,
		effects: []effect{{turns: 1, spell: spells[3]}},
	}
	if winner := searchMinMana(start); winner != nil {
		t.Fatalf("expected hard-mode penalty to kill player before poison: %+v", winner)
	}
}

func TestBattleKeyIgnoresCostAndEffectOrder(t *testing.T) {
	first := &state{
		hp: 12, mana: 212, bossHP: 17, bossDamage: 5, manaSpent: 500,
		effects: []effect{{turns: 3, spell: spells[3]}, {turns: 2, spell: spells[2]}},
	}
	second := *first
	second.manaSpent = 600
	second.effects = slices.Clone(first.effects)
	slices.Reverse(second.effects)
	if first.key() != second.key() {
		t.Fatal("equivalent battles received different keys")
	}
	for _, change := range []func(*state){
		func(s *state) { s.hp++ },
		func(s *state) { s.mana++ },
		func(s *state) { s.bossHP++ },
		func(s *state) { s.bossDamage++ },
		func(s *state) { s.hardMode = !s.hardMode },
		func(s *state) { s.effects[0].turns++ },
		func(s *state) { s.effects[1].turns++ },
		func(s *state) { s.effects = append(s.effects, effect{turns: 2, spell: spells[4]}) },
	} {
		changed := *first
		changed.effects = slices.Clone(first.effects)
		change(&changed)
		if first.key() == changed.key() {
			t.Fatalf("different battle received the same key: %+v", changed)
		}
	}
}

// Search every spell history on these small fights, without value-state pruning.
func minimumManaReference(start *state) *state {
	queue := &queues.PriorityQueue[*pqItem]{}
	heap.Push(queue, &pqItem{value: start, priority: start.manaSpent})
	for queue.Len() > 0 {
		current := heap.Pop(queue).(*pqItem).value
		if current.bossHP <= 0 {
			return current
		}
		for _, next := range current.transitions() {
			heap.Push(queue, &pqItem{value: next, priority: next.manaSpent})
		}
	}
	return nil
}

func TestMinimumManaMatchesUnprunedSearch(t *testing.T) {
	cases := 0
	for _, hp := range []int{5, 10} {
		for _, mana := range []int{126, 250} {
			for _, bossHP := range []int{4, 8, 12, 16} {
				for _, damage := range []int{2, 4, 6, 8} {
					for _, hard := range []bool{false, true} {
						start := &state{hp: hp, mana: mana, bossHP: bossHP, bossDamage: damage, hardMode: hard}
						got, want := searchMinMana(start), minimumManaReference(start)
						if (got == nil) != (want == nil) || got != nil && got.manaSpent != want.manaSpent {
							t.Fatalf("start %+v: got %+v, want %+v", start, got, want)
						}
						cases++
					}
				}
			}
		}
	}
	if cases != 128 {
		t.Fatalf("checked %d fights, want 128", cases)
	}
}

func TestWizardSyntheticFights(t *testing.T) {
	for _, test := range []struct{ bossHP, want int }{{5, 106}, {9, 159}} {
		start := &state{hp: 12, mana: 212, bossHP: test.bossHP, bossDamage: 2}
		winner := searchMinMana(start)
		if winner == nil || winner.manaSpent != test.want {
			t.Fatalf("boss HP %d: got %+v, want cost %d", test.bossHP, winner, test.want)
		}
	}
}
