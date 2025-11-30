package main

import (
	"aoc/solve"
	_ "aoc/solve/2021"
	_ "aoc/solve/2025"
)

func main() {
	driver := solve.NewDriver()
	driver.SolveAll(true)
}
