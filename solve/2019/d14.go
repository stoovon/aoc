package solve2019

import (
	"fmt"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day14 struct {
}

func (d Day14) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2019, Day: 14}
}

type chemical struct {
	name  string
	count int64
}

type reaction struct {
	output chemical
	inputs []chemical
}

func parseReactions(input string) map[string]reaction {
	reactions := make(map[string]reaction)
	lines := strings.Split(strings.TrimSpace(input), "\n")
	for _, line := range lines {
		parts := strings.Split(line, " => ")
		inStrs := strings.Split(parts[0], ", ")
		inputs := make([]chemical, len(inStrs))
		for i, s := range inStrs {
			fields := strings.Fields(s)
			n, _ := strconv.ParseInt(fields[0], 10, 64)
			inputs[i] = chemical{fields[1], n}
		}
		outFields := strings.Fields(parts[1])
		outN, _ := strconv.ParseInt(outFields[0], 10, 64)
		outChem := chemical{outFields[1], outN}
		reactions[outChem.name] = reaction{outChem, inputs}
	}
	return reactions
}

func oreFor(chem string, amount int64, reactions map[string]reaction, leftovers map[string]int64) int64 {
	if chem == "ORE" {
		return amount
	}
	if leftovers[chem] > 0 {
		used := min(amount, leftovers[chem])
		amount -= used
		leftovers[chem] -= used
	}
	if amount == 0 {
		return 0
	}
	react := reactions[chem]
	times := (amount + react.output.count - 1) / react.output.count
	totalOre := int64(0)
	for _, in := range react.inputs {
		totalOre += oreFor(in.name, in.count*times, reactions, leftovers)
	}
	leftovers[chem] += times*react.output.count - amount
	return totalOre
}

func (d Day14) Part1(input string) (string, error) {
	reactions := parseReactions(input)
	ore := oreFor("FUEL", 1, reactions, make(map[string]int64))
	return fmt.Sprint(ore), nil
}

func (d Day14) Part2(input string) (string, error) {
	reactions := parseReactions(input)
	oreAvailable := int64(1_000_000_000_000)
	// Set reasonable search bounds
	low, high := int64(1), oreAvailable

	for low < high {
		mid := (low + high + 1) / 2
		oreNeeded := oreFor("FUEL", mid, reactions, make(map[string]int64))
		if oreNeeded > oreAvailable {
			high = mid - 1
		} else {
			low = mid
		}
	}
	return fmt.Sprint(low), nil
}

func init() {
	solve.Register(Day14{})
}
