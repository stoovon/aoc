package grids

import (
	"image"
	"testing"
)

func TestKeyPadParsing(t *testing.T) {
	expected := map[rune]image.Point{
		'7': {},
		'8': {Y: 1},
		'9': {Y: 2},
		'4': {X: 1},
		'5': {X: 1, Y: 1},
		'6': {X: 1, Y: 2},
		'1': {X: 2},
		'2': {X: 2, Y: 1},
		'3': {X: 2, Y: 2},
		'0': {X: 3, Y: 1},
		'A': {X: 3, Y: 2},
	}

	actual := NewGridOptions().Parse(`
789
456
123
#0A
`).PointsFromTopLeft()

	if len(actual) != len(expected) {
		t.Errorf("Expected %d keys, got %d", len(expected), len(actual))
	}

	for k, v := range expected {
		if actual[k] != v {
			t.Errorf("Key %c: expected %v, got %v", k, v, actual[k])
		}
	}
}

func TestDirPadParsing(t *testing.T) {
	expected := map[rune]image.Point{
		'^': {Y: 1},
		'A': {Y: 2},
		'<': {X: 1},
		'v': {X: 1, Y: 1},
		'>': {X: 1, Y: 2},
	}

	actual := NewGridOptions().Parse(`
#^A
<v>
`).PointsFromTopLeft()
	if len(actual) != len(expected) {
		t.Errorf("Expected %d keys, got %d", len(expected), len(actual))
	}

	for k, v := range expected {
		if actual[k] != v {
			t.Errorf("Key %c: expected %v, got %v", k, v, actual[k])
		}
	}
}
