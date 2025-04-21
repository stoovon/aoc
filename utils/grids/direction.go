package grids

import "image"

// URDL returns the four possible directions we can face/move in
// using these indices lets us do rotation with modulo
func URDL() []image.Point {
	return []image.Point{
		{0, -1}, // Up
		{1, 0},  // Right
		{0, 1},  // Down
		{-1, 0}, // Left
	}
}

func DURL() []image.Point {
	return []image.Point{
		{0, 1},  // Down
		{0, -1}, // Up
		{1, 0},  // Right
		{-1, 0}, // Left
	}
}

func DirectionsPoints(udrl string) map[rune]image.Point {
	return map[rune]image.Point{
		rune(udrl[0]): {-1, 0},
		rune(udrl[1]): {1, 0},
		rune(udrl[2]): {0, 1},
		rune(udrl[3]): {0, -1},
	}
}

func DirectionsComplex(udrl string) map[rune]complex128 {
	return map[rune]complex128{
		rune(udrl[0]): -1i,
		rune(udrl[1]): 1,
		rune(udrl[2]): 1i,
		rune(udrl[3]): -1,
	}
}
