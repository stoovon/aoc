package solve2015

import (
	"fmt"
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/maths"
)

type Day15 struct {
}

func (d Day15) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2015, Day: 15}
}

type ingredient struct {
	Capacity   int
	Durability int
	Flavor     int
	Texture    int
	Calories   int
}

func parseIngredients(data string) map[string]ingredient {
	lines := strings.Split(strings.TrimSpace(data), "\n")
	ingredients := make(map[string]ingredient)

	for _, line := range lines {
		parts := strings.Split(line, ": ")
		name := parts[0]
		values := strings.Split(parts[1], ", ")

		currentIngredient := ingredient{}
		for _, value := range values {
			keyVal := strings.Split(value, " ")
			val, _ := strconv.Atoi(keyVal[1])
			switch keyVal[0] {
			case "capacity":
				currentIngredient.Capacity = val
			case "durability":
				currentIngredient.Durability = val
			case "flavor":
				currentIngredient.Flavor = val
			case "texture":
				currentIngredient.Texture = val
			case "calories":
				currentIngredient.Calories = val
			}
		}
		ingredients[name] = currentIngredient
	}

	return ingredients
}

func bestScore(ingredients map[string]ingredient, calorieTarget *int) int {
	best := 0
	keys := make([]string, 0, len(ingredients))
	for key := range ingredients {
		keys = append(keys, key)
	}

	for a := 0; a <= 100; a++ {
		for b := 0; b <= 100-a; b++ {
			for c := 0; c <= 100-a-b; c++ {
				d := 100 - a - b - c

				capacity := maths.Max(0, a*ingredients[keys[0]].Capacity+b*ingredients[keys[1]].Capacity+c*ingredients[keys[2]].Capacity+d*ingredients[keys[3]].Capacity)
				durability := maths.Max(0, a*ingredients[keys[0]].Durability+b*ingredients[keys[1]].Durability+c*ingredients[keys[2]].Durability+d*ingredients[keys[3]].Durability)
				flavor := maths.Max(0, a*ingredients[keys[0]].Flavor+b*ingredients[keys[1]].Flavor+c*ingredients[keys[2]].Flavor+d*ingredients[keys[3]].Flavor)
				texture := maths.Max(0, a*ingredients[keys[0]].Texture+b*ingredients[keys[1]].Texture+c*ingredients[keys[2]].Texture+d*ingredients[keys[3]].Texture)

				score := capacity * durability * flavor * texture

				if calorieTarget != nil {
					calories := maths.Max(0, a*ingredients[keys[0]].Calories+b*ingredients[keys[1]].Calories+c*ingredients[keys[2]].Calories+d*ingredients[keys[3]].Calories)
					if calories == *calorieTarget {
						best = maths.Max(best, score)
					}
				} else {
					best = maths.Max(best, score)
				}
			}
		}
	}

	return best
}

func (d Day15) Part1(input string) (string, error) {
	ingredients := parseIngredients(input)
	return fmt.Sprintf("%d", bestScore(ingredients, nil)), nil
}

func (d Day15) Part2(input string) (string, error) {
	ingredients := parseIngredients(input)
	calorieTarget := 500
	return fmt.Sprintf("%d", bestScore(ingredients, &calorieTarget)), nil
}

func init() {
	solve.Register(Day15{})
}
