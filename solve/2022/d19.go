package solve2022

import (
	"aoc/solve"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Day19 struct{}

func (Day19) Coords() solve.SolutionCoords { return solve.SolutionCoords{Year: 2022, Day: 19} }

const robotSearchWork = 10_000_000

var errRobotSearchLimit = errors.New("robot search exceeded its state-visit budget")
var errRobotOverflow = errors.New("robot answer exceeds int64")

type robotBlueprint struct {
	id      int64
	recipes [4][3]uint16
	caps    [3]uint16
}

func takeRobotDecimal(rest *string, suffix string) (string, error) {
	value, following, ok := strings.Cut(*rest, suffix)
	if !ok || value == "" {
		return "", errors.New("expected a complete robot blueprint record")
	}
	for _, c := range []byte(value) {
		if c < '0' || c > '9' {
			return "", errors.New("expected an unsigned decimal integer")
		}
	}
	*rest = following
	return value, nil
}
func parseRobotBlueprint(line string) (robotBlueprint, error) {
	var bp robotBlueprint
	rest, ok := strings.CutPrefix(line, "Blueprint ")
	if !ok {
		return bp, errors.New("expected 'Blueprint '")
	}
	number, err := takeRobotDecimal(&rest, ": Each ore robot costs ")
	if err != nil {
		return bp, err
	}
	bp.id, err = strconv.ParseInt(number, 10, 64)
	if err != nil || bp.id <= 0 {
		return bp, errors.New("expected positive int64 blueprint ID")
	}
	suffixes := [6]string{" ore. Each clay robot costs ", " ore. Each obsidian robot costs ", " ore and ", " clay. Each geode robot costs ", " ore and ", " obsidian."}
	var costs [6]uint16
	for i, suffix := range suffixes {
		number, err = takeRobotDecimal(&rest, suffix)
		if err != nil {
			return bp, err
		}
		value, err := strconv.ParseUint(number, 10, 16)
		if err != nil {
			return bp, fmt.Errorf("robot cost: %w", err)
		}
		costs[i] = uint16(value)
	}
	if rest != "" {
		return bp, errors.New("unexpected data after blueprint")
	}
	bp.recipes = [4][3]uint16{{costs[0], 0, 0}, {costs[1], 0, 0}, {costs[2], costs[3], 0}, {costs[4], 0, costs[5]}}
	for _, recipe := range bp.recipes {
		for r, cost := range recipe {
			if cost > bp.caps[r] {
				bp.caps[r] = cost
			}
		}
	}
	return bp, nil
}
func parseRobotBlueprints(input string) ([]robotBlueprint, error) {
	if strings.TrimSpace(input) == "" {
		return nil, errors.New("blueprint input is empty")
	}
	var result []robotBlueprint
	ids := map[int64]bool{}
	for i, line := range strings.Split(strings.TrimSuffix(input, "\n"), "\n") {
		bp, err := parseRobotBlueprint(strings.TrimSuffix(line, "\r"))
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", i+1, err)
		}
		if ids[bp.id] {
			return nil, fmt.Errorf("line %d: duplicate blueprint ID", i+1)
		}
		ids[bp.id] = true
		result = append(result, bp)
	}
	return result, nil
}

type robotState struct {
	remaining uint8
	resources [3]uint16
	robots    [4]uint8
	geodes    uint16
}
type robotSearch struct {
	blueprint     *robotBlueprint
	best          uint16
	remainingWork uint32
}

func (s *robotSearch) visit(state robotState) error {
	if s.remainingWork == 0 {
		return errRobotSearchLimit
	}
	s.remainingWork--
	remaining := uint32(state.remaining)
	// Stopping construction is always legal. This is a real outcome, unlike the bound.
	finish := uint32(state.geodes) + remaining*uint32(state.robots[3])
	if uint16(finish) > s.best {
		s.best = uint16(finish)
	}
	if remaining <= 1 || finish+remaining*(remaining-1)/2 <= uint32(s.best) {
		return nil
	}
	// Explore every legal next construction; ordering affects work, not correctness.
	for kind := 0; kind < 4; kind++ {
		if kind < 3 && uint16(state.robots[kind]) >= s.blueprint.caps[kind] {
			continue
		}
		recipe := s.blueprint.recipes[kind]
		wait := uint32(0)
		possible := true
		for resource, cost := range recipe {
			if cost <= state.resources[resource] {
				continue
			}
			rate := uint32(state.robots[resource])
			if rate == 0 {
				possible = false
				break
			}
			needed := uint32(cost - state.resources[resource])
			if neededWait := (needed + rate - 1) / rate; neededWait > wait {
				wait = neededWait
			}
		}
		duration := wait + 1
		if !possible || duration >= remaining {
			continue
		}
		next := state
		next.remaining -= uint8(duration)
		// One build per minute and the 32-minute horizon bound balances at 528,
		// robots at 33, and geodes at 496. Wait arithmetic uses uint32 before narrowing.
		for resource, cost := range recipe {
			next.resources[resource] = uint16(uint32(state.resources[resource]) + duration*uint32(state.robots[resource]) - uint32(cost))
		}
		next.geodes += uint16(duration * uint32(state.robots[3]))
		next.robots[kind]++
		if err := s.visit(next); err != nil {
			return err
		}
	}
	return nil
}
func maxRobotGeodes(bp *robotBlueprint, minutes uint8, work uint32) (uint16, error) {
	if minutes > 32 {
		return 0, errors.New("robot search supports at most 32 minutes")
	}
	search := robotSearch{blueprint: bp, remainingWork: work}
	err := search.visit(robotState{remaining: minutes, robots: [4]uint8{1, 0, 0, 0}})
	if err != nil {
		return 0, err
	}
	return search.best, nil
}
func checkedRobotProduct(a, b int64) (int64, error) {
	if b != 0 && a > math.MaxInt64/b {
		return 0, errRobotOverflow
	}
	return a * b, nil
}
func (Day19) Part1(input string) (string, error) {
	blueprints, err := parseRobotBlueprints(input)
	if err != nil {
		return "", err
	}
	sum := int64(0)
	for i := range blueprints {
		bp := &blueprints[i]
		geodes, err := maxRobotGeodes(bp, 24, robotSearchWork)
		if err != nil {
			return "", err
		}
		score, err := checkedRobotProduct(bp.id, int64(geodes))
		if err != nil {
			return "", err
		}
		if sum > math.MaxInt64-score {
			return "", errRobotOverflow
		}
		sum += score
	}
	return strconv.FormatInt(sum, 10), nil
}
func (Day19) Part2(input string) (string, error) {
	blueprints, err := parseRobotBlueprints(input)
	if err != nil {
		return "", err
	}
	product := int64(1)
	for i := 0; i < min(3, len(blueprints)); i++ {
		geodes, err := maxRobotGeodes(&blueprints[i], 32, robotSearchWork)
		if err != nil {
			return "", err
		}
		product, err = checkedRobotProduct(product, int64(geodes))
		if err != nil {
			return "", err
		}
	}
	return strconv.FormatInt(product, 10), nil
}
func init() { solve.Register(Day19{}) }
