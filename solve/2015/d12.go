package solve2015

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"

	"aoc/solve"
)

type Day12 struct {
}

var (
	numberRegex = regexp.MustCompile(`-?\d+`)
)

func (d Day12) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2015, Day: 12}
}

func noRed(obj interface{}) interface{} {
	switch v := obj.(type) {
	case map[string]interface{}:
		for _, value := range v {
			if str, ok := value.(string); ok && str == "red" {
				return nil
			}
		}
		filtered := make(map[string]interface{})
		for key, value := range v {
			filtered[key] = noRed(value)
		}
		return filtered
	case []interface{}:
		for i, value := range v {
			v[i] = noRed(value)
		}
		return v
	default:
		return v
	}
}

func sumNumbers(obj interface{}) int {
	switch v := obj.(type) {
	case float64:
		return int(v)
	case map[string]interface{}:
		sum := 0
		for _, value := range v {
			sum += sumNumbers(value)
		}
		return sum
	case []interface{}:
		sum := 0
		for _, value := range v {
			sum += sumNumbers(value)
		}
		return sum
	default:
		return 0
	}
}

func (d Day12) sumOfAllNumbers(input string) (int, int, error) {
	matches := numberRegex.FindAllString(input, -1)
	sumA := 0
	for _, match := range matches {
		num, _ := strconv.Atoi(match)
		sumA += num
	}

	var jsonData interface{}
	if err := json.Unmarshal([]byte(input), &jsonData); err != nil {
		return 0, 0, errors.New("error parsing JSON")
	}
	filteredData := noRed(jsonData)
	sumB := sumNumbers(filteredData)

	return sumA, sumB, nil
}

func (d Day12) Part1(input string) (string, error) {
	result, _, err := d.sumOfAllNumbers(input)

	if err != nil {
		return "", err
	}

	return strconv.Itoa(result), nil
}

func (d Day12) Part2(input string) (string, error) {
	_, result, err := d.sumOfAllNumbers(input)

	if err != nil {
		return "", err
	}

	return strconv.Itoa(result), nil
}

func init() {
	solve.Register(Day12{})
}
