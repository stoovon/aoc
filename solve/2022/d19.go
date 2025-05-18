package solve2022

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day19 struct{}

func (d Day19) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2022, Day: 19}
}

type OreAmount = uint16
type Recipe = [4]OreAmount

type BlueprintSim struct {
	Recipes [4]Recipe
}

type State struct {
	Ores   [4]OreAmount
	Robots [4]uint16
	Time   uint16
}

func blueprintToSim(bp Blueprint) BlueprintSim {
	return BlueprintSim{
		Recipes: [4]Recipe{
			{OreAmount(bp.oreRobotOre), 0, 0, 0},
			{OreAmount(bp.clayRobotOre), 0, 0, 0},
			{OreAmount(bp.obsRobotOre), OreAmount(bp.obsRobotClay), 0, 0},
			{OreAmount(bp.geodeRobotOre), 0, OreAmount(bp.geodeRobotObsidian), 0},
		},
	}
}

var blueprintRE = regexp.MustCompile(`Blueprint (\d+): Each ore robot costs (\d+) ore\. Each clay robot costs (\d+) ore\. Each obsidian robot costs (\d+) ore and (\d+) clay\. Each geode robot costs (\d+) ore and (\d+) obsidian\.`)

type Blueprint struct {
	id                 int
	oreRobotOre        int
	clayRobotOre       int
	obsRobotOre        int
	obsRobotClay       int
	geodeRobotOre      int
	geodeRobotObsidian int
}

func parseBlueprints(input string) []Blueprint {
	var blueprints []Blueprint
	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		m := blueprintRE.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		nums := make([]int, 0, 7)
		for i := 1; i <= 7; i++ {
			n, _ := strconv.Atoi(m[i])
			nums = append(nums, n)
		}
		blueprints = append(blueprints, Blueprint{
			id:                 nums[0],
			oreRobotOre:        nums[1],
			clayRobotOre:       nums[2],
			obsRobotOre:        nums[3],
			obsRobotClay:       nums[4],
			geodeRobotOre:      nums[5],
			geodeRobotObsidian: nums[6],
		})
	}
	return blueprints
}

func recurseSimulation(
	blueprint *BlueprintSim,
	state State,
	maxTime uint16,
	maxRobots *[4]uint16,
	maxGeodes *OreAmount,
) {
	hasRecursed := false
	for i := 0; i < 4; i++ {
		if state.Robots[i] == maxRobots[i] {
			continue
		}
		recipe := blueprint.Recipes[i]
		waitTime := uint16(0)
		for oreType := 0; oreType < 3; oreType++ {
			if recipe[oreType] == 0 {
				continue
			} else if recipe[oreType] <= state.Ores[oreType] {
				continue
			} else if state.Robots[oreType] == 0 {
				waitTime = maxTime + 1
				break
			} else {
				needed := recipe[oreType] - state.Ores[oreType]
				wt := (needed + state.Robots[oreType] - 1) / state.Robots[oreType]
				if wt > waitTime {
					waitTime = wt
				}
			}
		}
		timeFinished := state.Time + waitTime + 1
		if timeFinished >= maxTime {
			continue
		}
		var newOres [4]OreAmount
		var newRobots [4]uint16
		for o := 0; o < 4; o++ {
			newOres[o] = state.Ores[o] + OreAmount(state.Robots[o])*(waitTime+1) - recipe[o]
			newRobots[o] = state.Robots[o]
			if o == i {
				newRobots[o]++
			}
		}
		remainingTime := maxTime - timeFinished
		// Prune if we can't beat current max
		possible := ((remainingTime-1)*remainingTime)/2 + uint16(newOres[3]) + remainingTime*newRobots[3]
		if possible < uint16(*maxGeodes) {
			continue
		}
		hasRecursed = true
		recurseSimulation(blueprint, State{Ores: newOres, Robots: newRobots, Time: timeFinished}, maxTime, maxRobots, maxGeodes)
	}
	if !hasRecursed {
		geodes := state.Ores[3] + OreAmount(state.Robots[3])*(maxTime-state.Time)
		if geodes > *maxGeodes {
			*maxGeodes = geodes
		}
	}
}

func simulateBlueprint(blueprint *BlueprintSim, maxTime uint16) OreAmount {
	var maxRobots [4]uint16
	for i := 0; i < 3; i++ {
		maxSeen := OreAmount(0)
		for _, r := range blueprint.Recipes {
			if r[i] > maxSeen {
				maxSeen = r[i]
			}
		}
		maxRobots[i] = uint16(maxSeen)
	}
	maxRobots[3] = math.MaxUint16
	maxGeodes := OreAmount(0)
	recurseSimulation(blueprint, State{
		Ores:   [4]OreAmount{},
		Robots: [4]uint16{1, 0, 0, 0},
		Time:   0,
	}, maxTime, &maxRobots, &maxGeodes)
	return maxGeodes
}

func (d Day19) Part1(input string) (string, error) {
	blueprints := parseBlueprints(input)
	sum := 0
	for i, bp := range blueprints {
		sim := blueprintToSim(bp)
		geodes := simulateBlueprint(&sim, 24)
		sum += int(geodes) * (i + 1)
	}
	return fmt.Sprintf("%d", sum), nil
}

func (d Day19) Part2(input string) (string, error) {
	blueprints := parseBlueprints(input)
	prod := 1
	for i := 0; i < 3 && i < len(blueprints); i++ {
		sim := blueprintToSim(blueprints[i])
		geodes := simulateBlueprint(&sim, 32)
		prod *= int(geodes)
	}
	return fmt.Sprintf("%d", prod), nil
}

func init() {
	solve.Register(Day19{})
}
