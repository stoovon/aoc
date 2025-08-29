package solve2021

import (
	"aoc/solve"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Day19 struct{}

type Vec3 struct {
	X, Y, Z int
}

// Returns all 24 possible orientations of a Vec3
func (v Vec3) Orientations() []Vec3 {
	// 6 axis directions (x, y, z, -x, -y, -z), 4 rotations each
	o := make([]Vec3, 0, 24)
	x, y, z := v.X, v.Y, v.Z
	axes := [][3]int{
		{x, y, z}, {x, -y, -z}, {x, z, -y}, {x, -z, y},
		{-x, y, -z}, {-x, -y, z}, {-x, z, y}, {-x, -z, -y},
		{y, z, x}, {y, -z, -x}, {y, x, -z}, {y, -x, z},
		{-y, z, -x}, {-y, -z, x}, {-y, x, z}, {-y, -x, -z},
		{z, x, y}, {z, -x, -y}, {z, y, -x}, {z, -y, x},
		{-z, x, -y}, {-z, -x, y}, {-z, y, x}, {-z, -y, -x},
	}
	for _, a := range axes {
		o = append(o, Vec3{a[0], a[1], a[2]})
	}
	return o
}

// Returns all 24 orientations for a slice of Vec3
func allOrientations(beacons []Vec3) [][]Vec3 {
	out := make([][]Vec3, 24)
	for i := 0; i < 24; i++ {
		out[i] = make([]Vec3, len(beacons))
		for j, b := range beacons {
			out[i][j] = b.Orientations()[i]
		}
	}
	return out
}

type Scanner struct {
	Beacons []Vec3
}

func (d Day19) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2021, Day: 19}
}

func parseInput(input string) []Scanner {
	blocks := strings.Split(strings.TrimSpace(input), "\n\n")
	scanners := make([]Scanner, 0, len(blocks))
	for _, block := range blocks {
		lines := strings.Split(block, "\n")
		beacons := make([]Vec3, 0)
		for _, line := range lines[1:] { // skip header
			if line == "" {
				continue
			}
			parts := strings.Split(line, ",")
			if len(parts) != 3 {
				continue
			}
			x, _ := strconv.Atoi(parts[0])
			y, _ := strconv.Atoi(parts[1])
			z, _ := strconv.Atoi(parts[2])
			beacons = append(beacons, Vec3{x, y, z})
		}
		scanners = append(scanners, Scanner{Beacons: beacons})
	}
	return scanners
}

func (a Vec3) Sub(b Vec3) Vec3 {
	return Vec3{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

func (a Vec3) Add(b Vec3) Vec3 {
	return Vec3{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

func beaconKey(v Vec3) string {
	return fmt.Sprintf("%d,%d,%d", v.X, v.Y, v.Z)
}

const numerOfOverlappingBeaconsRequiredForAlignment = 12

// Aligns all scanners and returns the set of global beacons and scanner positions
func alignScanners(scanners []Scanner) (map[string]Vec3, map[int]Vec3, error) {
	aligned := map[int]struct{}{0: {}} // scanner 0 is reference
	scannerPos := map[int]Vec3{0: {0, 0, 0}}
	beaconSet := map[string]Vec3{}
	for _, b := range scanners[0].Beacons {
		beaconSet[beaconKey(b)] = b
	}
	unaligned := map[int]struct{}{}
	for i := 1; i < len(scanners); i++ {
		unaligned[i] = struct{}{}
	}

	for len(unaligned) > 0 {
		progress := false
		for ua := range unaligned {
			found := false
			// Use the union of all currently aligned beacons as reference
			refBeacons := make([]Vec3, 0, len(beaconSet))
			for _, b := range beaconSet {
				refBeacons = append(refBeacons, b)
			}
			orientations := allOrientations(scanners[ua].Beacons)
			for _, obs := range orientations {
				// Count translation vectors - greatly reduces unnecessary comparisons
				transCount := map[string]int{}
				transVec := map[string]Vec3{}
				for _, rb := range refBeacons {
					for _, ob := range obs {
						offset := rb.Sub(ob)
						key := beaconKey(offset)
						transCount[key]++
						transVec[key] = offset
					}
				}
				for key, count := range transCount {
					if count >= numerOfOverlappingBeaconsRequiredForAlignment {
						offset := transVec[key]
						// Verify actual overlap
						actualOverlap := 0
						refSet := map[string]struct{}{}
						for _, b := range refBeacons {
							refSet[beaconKey(b)] = struct{}{}
						}
						for _, b := range obs {
							global := b.Add(offset)
							if _, ok := refSet[beaconKey(global)]; ok {
								actualOverlap++
							}
						}
						if actualOverlap >= numerOfOverlappingBeaconsRequiredForAlignment {
							aligned[ua] = struct{}{}
							scannerPos[ua] = offset
							for _, b := range obs {
								global := b.Add(offset)
								beaconSet[beaconKey(global)] = global
							}
							delete(unaligned, ua)
							progress = true
							found = true
							break
						}
					}
				}
				if found {
					break
				}
			}
			if found {
				break
			}
		}
		if !progress {
			return nil, nil, errors.New("Failed to align all scanners")
		}
	}
	return beaconSet, scannerPos, nil
}

func (d Day19) Part1(input string) (string, error) {
	scanners := parseInput(input)
	beaconSet, _, err := alignScanners(scanners)
	if err != nil {
		return err.Error(), nil
	}
	return strconv.Itoa(len(beaconSet)), nil
}

func manhattan(a, b Vec3) int {
	return abs(a.X-b.X) + abs(a.Y-b.Y) + abs(a.Z-b.Z)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (d Day19) Part2(input string) (string, error) {
	scanners := parseInput(input)
	_, scannerPos, err := alignScanners(scanners)
	if err != nil {
		return err.Error(), nil
	}
	maxDist := 0
	for i := range scannerPos {
		for j := range scannerPos {
			if i == j {
				continue
			}
			dist := manhattan(scannerPos[i], scannerPos[j])
			if dist > maxDist {
				maxDist = dist
			}
		}
	}
	return strconv.Itoa(maxDist), nil
}

func init() {
	solve.Register(Day19{})
}
