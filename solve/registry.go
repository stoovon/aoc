package solve

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
	registry[solver.Coords()] = solver
}

func GetAllSolvers() map[SolutionCoords]Solver {
	return registry
}
