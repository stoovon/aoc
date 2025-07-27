package solve2020

import (
	"aoc/solve"
	"strconv"
	"strings"
	"sort"
)

type Day21 struct{}

func (d Day21) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2020, Day: 21}
}

func (d Day21) Part1(input string) (string, error) {
	_, allIngredients, possible := parseFoodsAndAllergens(input)
	unsafe := make(map[string]struct{})
	for _, s := range possible {
		for k := range s {
			unsafe[k] = struct{}{}
		}
	}
	safeCount := 0
	for ingr, count := range allIngredients {
		if _, bad := unsafe[ingr]; !bad {
			safeCount += count
		}
	}
	return strconv.Itoa(safeCount), nil
}

func (d Day21) Part2(input string) (string, error) {
	_, _, possible := parseFoodsAndAllergens(input)
	// Deduce allergen to ingredient mapping
	allergenToIngredient := make(map[string]string)
	used := make(map[string]struct{})
	for len(allergenToIngredient) < len(possible) {
		progress := false
		for allergen, opts := range possible {
			var candidate string
			count := 0
			for ingr := range opts {
				if _, already := used[ingr]; !already {
					candidate = ingr
					count++
				}
			}
			if count == 1 {
				allergenToIngredient[allergen] = candidate
				used[candidate] = struct{}{}
				progress = true
			}
		}
		if !progress {
			break // avoid infinite loop if input is malformed
		}
	}
	// Sort allergens alphabetically and join their ingredients
	var allergens []string
	for a := range allergenToIngredient {
		allergens = append(allergens, a)
	}
	sort.Strings(allergens)
	var result []string
	for _, a := range allergens {
		result = append(result, allergenToIngredient[a])
	}
	return strings.Join(result, ","), nil
}
// Shared parsing for both parts
func parseFoodsAndAllergens(input string) ([]map[string]struct{}, map[string]int, map[string]map[string]struct{}) {
	var foods []map[string]struct{}
	allIngredients := make(map[string]int)
	allergenToFoods := make(map[string][]map[string]struct{})
	lines := strings.Split(strings.TrimSpace(input), "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, " (contains ", 2)
		ingr := make(map[string]struct{})
		for _, i := range strings.Fields(parts[0]) {
			ingr[i] = struct{}{}
			allIngredients[i]++
		}
		if len(parts) > 1 {
			als := strings.TrimSuffix(parts[1], ")")
			for _, a := range strings.Split(als, ", ") {
				allergenToFoods[a] = append(allergenToFoods[a], ingr)
			}
		}
		foods = append(foods, ingr)
	}
	// For each allergen, find possible ingredients (intersection)
	possible := make(map[string]map[string]struct{})
	for allergen, foodList := range allergenToFoods {
		var inter map[string]struct{}
		for _, ingr := range foodList {
			if inter == nil {
				inter = make(map[string]struct{})
				for k := range ingr {
					inter[k] = struct{}{}
				}
			} else {
				for k := range inter {
					if _, ok := ingr[k]; !ok {
						delete(inter, k)
					}
				}
			}
		}
		possible[allergen] = inter
	}
	return foods, allIngredients, possible
}

func init() {
	solve.Register(Day21{})
}
