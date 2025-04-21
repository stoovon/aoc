package main

import (
	"aoc/solve"

	_ "aoc/solve/2015"
	_ "aoc/solve/2024"
)

func main() {
	driver := solve.NewDriver()
	driver.SolveAll(false)
}
