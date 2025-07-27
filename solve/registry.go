package solve

import (
	"strconv"
)

type SolutionCoords struct {
	Year      int
	Day       int
	NoCaching bool
}

type SolverFunc func(input string) (string, error)

type Solver interface {
	Coords() SolutionCoords
	Part1(input string) (string, error)
	Part2(input string) (string, error)
}

var registry = make(map[SolutionCoords]Solver)

func Register(solver Solver) {
	coords := solver.Coords()
	if _, exists := registry[coords]; exists {
		panic("duplicate solver registration for coords: " + strconv.Itoa(coords.Year) + ", day: " + strconv.Itoa(coords.Day))
	}
	registry[coords] = solver
}

func GetAllSolvers() map[SolutionCoords]Solver {
	return registry
}
