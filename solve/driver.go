package solve

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"

	"aoc/httpclient"
)

type Driver struct {
	client *httpclient.Client
}

func NewDriver() *Driver {
	return &Driver{
		client: httpclient.NewClient(),
	}
}

func (d *Driver) SolveAll(fullCache bool) {
	allSolvers := GetAllSolvers()

	coordsList := make([]SolutionCoords, 0, len(GetAllSolvers()))
	for coords := range allSolvers {
		coordsList = append(coordsList, coords)
	}

	// Sort the coordinates by year and day
	sort.Slice(coordsList, func(i, j int) bool {
		if coordsList[i].Year == coordsList[j].Year {
			return coordsList[i].Day < coordsList[j].Day
		}
		return coordsList[i].Year < coordsList[j].Year
	})

	for _, coords := range coordsList {
		solver := allSolvers[coords]
		if coords.Day == 1 {
			fmt.Printf("\nSolve %d: ", coords.Year)
		}

		err := d.Solve(solver, fullCache)
		if err != nil {
			fmt.Printf("Error solving %d-%02d: %v\n", coords.Year, coords.Day, err)
			continue
		}
	}
}

func (d *Driver) Solve(solver Solver, fullCache bool) error {
	coords := solver.Coords()

	input, err := d.client.GetInput(coords.Year, coords.Day)
	if err != nil {
		return fmt.Errorf("could not get input for day %d-%d: %v", coords.Year, coords.Day, err)
	}

	err = d.submitPart(input, solver, 1, fullCache)
	if err != nil {
		return fmt.Errorf("error submitting part 1: %v", err)
	}

	fmt.Printf(".")

	if coords.Day == 25 {
		fmt.Printf("🎄")
		return nil
	}

	err = d.submitPart(input, solver, 2, fullCache)
	if err != nil {
		return fmt.Errorf("error submitting part 2: %v", err)
	}

	fmt.Printf(":")

	return nil
}

func (d *Driver) submitPart(input string, solver Solver, part int, fullCache bool) error {
	coords := solver.Coords()

	cache := fullCache
	if coords.NoCaching {
		cache = false
	}

	if cache && d.client.HasSolution(coords.Year, coords.Day, part) {
		return nil
	}

	solverFunc := solver.Part1
	if part == 2 {
		solverFunc = solver.Part2
	}

	solution, err := solverFunc(input)
	if err != nil {
		return fmt.Errorf("error solving part %d: %v", part, err)
	}

	response, err := d.client.SubmitAnswer(coords.Year, coords.Day, part, solution)
	if err != nil {
		return err
	}

	if strings.Contains(response, "That's the right answer!") {
		return nil
	}

	fmt.Printf("Part %d: %s\n", part, solution)

	if strings.HasPrefix(response, "That's not the right answer.") {
		err := d.client.KnownBadAnswer(coords.Year, coords.Day, part, response)
		if err != nil {
			return fmt.Errorf("error marking answer as incorrect: %v", err)
		}

		return fmt.Errorf("Incorrect answer for part %d.\n", part)
	} else if strings.HasPrefix(response, "Please don't repeatedly") {
		// Given enhanced caching, we don't expect to see this one much.
		return fmt.Errorf("Rate limit exceeded for part %d.\n", part)
	} else if response == httpclient.NotKnownGood {
		return errors.New(response)
	} else {
		// To avoid spamming the AoC servers, we should exit if anything is unexpected.
		log.Fatalf("Unexpected response for part %d: %s\n", part, response)
	}

	return nil
}
