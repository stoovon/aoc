package solve2015

import "testing"

func TestPoisonWinAtPlayerTurnStartCostsNoAdditionalMana(t *testing.T) {
	for _, mana := range []int{0, 53} {
		start := &state{
			hp: 10, mana: mana, bossHP: 3, bossDamage: 8, manaSpent: 173,
			effects: []effect{{turns: 1, spell: spells[3]}},
		}
		winner := searchAStar(start)
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
	if winner := searchAStar(start); winner != nil {
		t.Fatalf("expected hard-mode penalty to kill player before poison: %+v", winner)
	}
}
