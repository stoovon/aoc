package solve2021

import (
	"aoc/solve"
	"fmt"
	"strconv"
)

type Day21 struct{}

func (d Day21) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2021, Day: 21}
}

func (d Day21) Part1(input string) (string, error) {
	// Parse the input to get the starting positions of the players
	var p1Start, p2Start int
	_, err := fmt.Sscanf(input, "Player 1 starting position: %d\nPlayer 2 starting position: %d", &p1Start, &p2Start)
	if err != nil {
		return "", err
	}

	// Initialize game state
	p1Score, p2Score := 0, 0
	die := 1
	dieRolls := 0

	rollDie := func() int {
		result := die
		die = (die % 100) + 1
		dieRolls++
		return result
	}

	// Play the game
	for {
		// Player 1's turn
		p1Move := rollDie() + rollDie() + rollDie()
		p1Start = (p1Start+p1Move-1)%10 + 1
		p1Score += p1Start
		if p1Score >= 1000 {
			return fmt.Sprintf("%d", p2Score*dieRolls), nil
		}

		// Player 2's turn
		p2Move := rollDie() + rollDie() + rollDie()
		p2Start = (p2Start+p2Move-1)%10 + 1
		p2Score += p2Start
		if p2Score >= 1000 {
			return strconv.Itoa(p1Score * dieRolls), nil
		}
	}
}

func (d Day21) Part2(input string) (string, error) {
	// Parse the input to get the starting positions of the players
	var p1Start, p2Start int
	_, err := fmt.Sscanf(input, "Player 1 starting position: %d\nPlayer 2 starting position: %d", &p1Start, &p2Start)
	if err != nil {
		return "", err
	}

	// Memoization for quantum universes
	type state struct {
		p1Pos, p2Pos, p1Score, p2Score int
	}
	memo := make(map[state][2]int64)

	// Quantum die rolls and their frequencies
	rollFrequencies := map[int]int{
		3: 1, 4: 3, 5: 6, 6: 7, 7: 6, 8: 3, 9: 1,
	}

	var play func(p1Pos, p2Pos, p1Score, p2Score int) [2]int64
	play = func(p1Pos, p2Pos, p1Score, p2Score int) [2]int64 {
		if p1Score >= 21 {
			return [2]int64{1, 0}
		}
		if p2Score >= 21 {
			return [2]int64{0, 1}
		}

		key := state{p1Pos, p2Pos, p1Score, p2Score}
		if result, exists := memo[key]; exists {
			return result
		}

		result := [2]int64{0, 0}
		for roll, freq := range rollFrequencies {
			newP1Pos := (p1Pos+roll-1)%10 + 1
			newP1Score := p1Score + newP1Pos
			outcomes := play(p2Pos, newP1Pos, p2Score, newP1Score)
			result[0] += int64(freq) * outcomes[1]
			result[1] += int64(freq) * outcomes[0]
		}

		memo[key] = result
		return result
	}

	wins := play(p1Start, p2Start, 0, 0)
	return fmt.Sprintf("%d", max(wins[0], wins[1])), nil
}

func init() {
	solve.Register(Day21{})
}
