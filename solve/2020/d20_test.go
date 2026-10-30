package solve2020

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"

	"aoc/solve"
)

func testJigsawLabels(labels [][4]byte) *jigsaw {
	j := &jigsaw{size: 2}
	for index, edges := range labels {
		tile := jigsawTile{id: int64(index + 1)}
		for i := range tile.cells {
			tile.cells[i] = '.'
		}
		for side, label := range edges {
			for bit := 0; bit < 8; bit++ {
				if label&(1<<bit) != 0 {
					row, col := 0, bit+1
					switch side {
					case 1:
						row, col = bit+1, 9
					case 2:
						row, col = 9, bit+1
					case 3:
						row, col = bit+1, 0
					}
					tile.cells[row*10+col] = '#'
				}
			}
		}
		tile.cells[(index+1)*10+index+1] = '#'
		tile.readEdges()
		j.tiles = append(j.tiles, tile)
	}
	j.indexBorders()
	return j
}

func testJigsawInput(tiles []jigsawTile) string {
	var out strings.Builder
	for _, tile := range tiles {
		fmt.Fprintf(&out, "Tile %d:\n", tile.id)
		for row := 0; row < 10; row++ {
			out.Write(tile.cells[row*10 : (row+1)*10])
			out.WriteByte('\n')
		}
		out.WriteByte('\n')
	}
	return out.String()
}

func TestJigsawParser(t *testing.T) {
	input := "Tile 12345:\n" + strings.Repeat("..........\n", 10)
	for _, valid := range []string{input, strings.ReplaceAll(input, "\n", "\r\n"), "\n" + input + "\n\n"} {
		j, err := parseJigsaw(valid)
		if err != nil {
			t.Fatal(err)
		}
		if j.size != 1 || j.tiles[0].id != 12345 {
			t.Fatal("incorrect dimensions or truncated ID")
		}
	}
	for _, invalid := range []string{
		"", "\n\n", strings.ReplaceAll(input, "12345", "0"), strings.ReplaceAll(input, "12345", "-3"),
		strings.ReplaceAll(input, "12345", "2x"), strings.ReplaceAll(input, "12345", "9223372036854775808"),
		strings.ReplaceAll(input, "Tile ", "Tiles "), strings.ReplaceAll(input, ":", ":junk"), strings.ReplaceAll(input, ":", ""),
		strings.ReplaceAll(input, "..........", "........."), strings.ReplaceAll(input, "..........", "..........."),
		strings.Replace(input, ".", "é", 1), strings.Replace(input, ".", "x", 1), "Tile 1:\n..........",
		input + "\n" + input, input + "\n" + strings.ReplaceAll(input, "12345", "2"),
		input + "..........\n", input + strings.ReplaceAll(input, "12345", "2"),
	} {
		if _, err := parseJigsaw(invalid); err == nil {
			t.Errorf("accepted malformed input %q", invalid)
		}
	}
	var large strings.Builder
	for id := 1; id <= 1025; id++ {
		fmt.Fprintf(&large, "Tile %d:\n%s\n", id, strings.Repeat("..........\n", 10))
	}
	if _, err := parseJigsaw(large.String()); err == nil {
		t.Fatal("tile limit ignored")
	}
}

func TestJigsawPartComposition(t *testing.T) {
	j := testJigsawLabels([][4]byte{{1, 9, 10, 2}, {3, 4, 11, 9}, {10, 12, 5, 6}, {11, 7, 8, 12}})
	ids := []int64{2, 3, 5, 7}
	for i := range j.tiles {
		j.tiles[i].id = ids[i]
	}
	input := testJigsawInput(j.tiles)
	got, err := (Day20{}).Part1(input)
	if err != nil || got != "210" {
		t.Fatalf("corner product: %s %v", got, err)
	}
	j.tiles[0].id = math.MaxInt64
	if _, err := (Day20{}).Part1(testJigsawInput(j.tiles)); err == nil {
		t.Fatal("unchecked product overflow")
	}
	single := "Tile 37:\n" + strings.Repeat("..........\n", 10)
	if _, err := (Day20{}).Part1(single); err == nil {
		t.Fatal("missing four-corner contract")
	}
	if got, err := (Day20{}).Part2(single); err != nil || got != "0" {
		t.Fatalf("single tile: %s %v", got, err)
	}
	solver, ok := solve.GetAllSolvers()[solve.SolutionCoords{Year: 2020, Day: 20}]
	if !ok {
		t.Fatal("solver not registered")
	}
	for _, part := range []func(string) (string, error){solver.Part1, solver.Part2} {
		if _, err := part("bad"); err == nil {
			t.Fatal("registered solver swallowed parse error")
		}
	}
}

func TestJigsawCanonicalEdgesAndOwnership(t *testing.T) {
	for edge := 0; edge < 1024; edge++ {
		digits := fmt.Sprintf("%010b", edge)
		reverse := []byte(digits)
		for a, b := 0, 9; a < b; a, b = a+1, b-1 {
			reverse[a], reverse[b] = reverse[b], reverse[a]
		}
		reversed := 0
		for _, digit := range reverse {
			reversed = reversed*2 + int(digit-'0')
		}
		if got := canonicalJigsawEdge(uint16(edge)); got != uint16(min(edge, reversed)) {
			t.Fatal("canonical edge mismatch")
		}
	}
	// All four borders are palindromes owned by one tile, not eight neighbours.
	j, _ := parseJigsaw("Tile 17:\n" + strings.Repeat("..........\n", 10))
	if len(j.borders[0]) != 1 || j.tiles[0].matched != [4]bool{} {
		t.Fatal("own repeated/palindromic borders counted as neighbours")
	}
	tile := j.tiles[0]
	tile.id = 23
	j.tiles = append(j.tiles, tile)
	j.indexBorders()
	if len(j.borders[0]) != 2 || j.tiles[0].matched != [4]bool{true, true, true, true} {
		t.Fatal("distinct neighbour not identified")
	}
}

func TestJigsawTransformsAndCorners(t *testing.T) {
	for mask := 0; mask < 16; mask++ {
		j := testJigsawLabels([][4]byte{{1, 2, 3, 4}})
		tile := j.tiles[0]
		for side := 0; side < 4; side++ {
			tile.matched[side] = mask&(1<<side) != 0
		}
		want := mask == 3 || mask == 6 || mask == 12 || mask == 9
		original := tile
		for rotation := 0; rotation < 4; rotation++ {
			if tile.corner() != want {
				t.Fatal("incorrect corner adjacency")
			}
			tile = tile.rotate()
		}
		if tile != original || tile.flip().flip() != original || tile.flip().corner() != want {
			t.Fatal("transform failed to preserve tile state")
		}
	}
	j := testJigsawLabels([][4]byte{{1, 2, 3, 4}})
	if len(j.tiles[0].orientations()) != 8 {
		t.Fatal("missing distinct orientation")
	}
	flat, _ := parseJigsaw("Tile 1:\n" + strings.Repeat("..........\n", 10))
	if len(flat.tiles[0].orientations()) != 1 {
		t.Fatal("symmetric orientations not deduplicated")
	}
	img := renderJigsaw(j.tiles, 1)
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			if img.cells[row*8+col] != j.tiles[0].cells[(row+1)*10+col+1] {
				t.Fatal("incorrect interior extraction")
			}
		}
	}
}

func checkTestJigsawAssembly(t *testing.T, j *jigsaw, placed []jigsawTile) {
	t.Helper()
	ids := make(map[int64]bool)
	if len(placed) != len(j.tiles) {
		t.Fatal("incomplete arrangement")
	}
	for index, tile := range placed {
		if ids[tile.id] {
			t.Fatal("tile reused")
		}
		ids[tile.id] = true
		for offset := 0; offset < 10; offset++ {
			if index%j.size > 0 && placed[index-1].cells[offset*10+9] != tile.cells[offset*10] {
				t.Fatal("left border mismatch")
			}
			if index >= j.size && placed[index-j.size].cells[90+offset] != tile.cells[offset] {
				t.Fatal("top border mismatch")
			}
		}
	}
}

// The oracle transforms only pixel arrays and compares individual touching cells.
// It uses neither canonical edges nor the production candidate index and pruning.
func referenceJigsawVariants(cells [100]byte) [][100]byte {
	var result [][100]byte
	for flipped := 0; flipped < 2; flipped++ {
		current := cells
		if flipped == 1 {
			for row := 0; row < 10; row++ {
				for col := 0; col < 10; col++ {
					current[row*10+col] = cells[row*10+9-col]
				}
			}
		}
		for rotation := 0; rotation < 4; rotation++ {
			result = append(result, current)
			var next [100]byte
			for row := 0; row < 10; row++ {
				for col := 0; col < 10; col++ {
					next[row*10+col] = current[(9-col)*10+row]
				}
			}
			current = next
		}
	}
	return result
}

func referenceJigsawCompletion(variants [][][100]byte, placed [][100]byte, used uint8) bool {
	if len(placed) == 4 {
		return true
	}
	position := len(placed)
	for index, orientations := range variants {
		if used&(1<<index) != 0 {
			continue
		}
		for _, cells := range orientations {
			fits := true
			for offset := 0; offset < 10; offset++ {
				if position >= 2 && placed[position-2][90+offset] != cells[offset] || position%2 == 1 && placed[position-1][offset*10+9] != cells[offset*10] {
					fits = false
					break
				}
			}
			if fits && referenceJigsawCompletion(variants, append(placed, cells), used|1<<index) {
				return true
			}
		}
	}
	return false
}

func TestJigsawForcedBacktracking(t *testing.T) {
	j := testJigsawLabels([][4]byte{{4, 2, 1, 4}, {4, 1, 4, 1}, {4, 4, 1, 1}, {3, 2, 2, 1}})
	variants := make([][]jigsawTile, 4)
	reference := make([][][100]byte, 4)
	for i, tile := range j.tiles {
		variants[i] = tile.orientations()
		reference[i] = referenceJigsawVariants(tile.cells)
	}
	work := jigsawWorkLimit
	candidates, err := j.candidates(variants, nil, make([]bool, 4), []int{0, 1, 2, 3}, &work)
	if err != nil || candidates[len(candidates)-1] != (jigsawPlacement{0, 0}) {
		t.Fatal("fixture does not select intended dead end")
	}
	if referenceJigsawCompletion(reference, [][100]byte{j.tiles[0].cells}, 1) {
		t.Fatal("first local choice unexpectedly solvable")
	}
	placed, err := j.assemble(jigsawWorkLimit)
	if err != nil {
		t.Fatal(err)
	}
	if placed[0].cells == j.tiles[0].cells {
		t.Fatal("failed to undo dead-end root")
	}
	checkTestJigsawAssembly(t, j, placed)
}

func TestJigsawAnchoringAndLimits(t *testing.T) {
	j := testJigsawLabels([][4]byte{{1, 9, 10, 2}, {3, 4, 11, 9}, {10, 12, 5, 6}, {11, 7, 8, 12}})
	if len(j.corners()) != 4 {
		t.Fatal("fixture must have four forced corners")
	}
	placed, err := j.assemble(jigsawWorkLimit)
	if err != nil {
		t.Fatal(err)
	}
	if placed[0].id != j.tiles[0].id || placed[0].matched[0] || placed[0].matched[3] {
		t.Fatal("proven corner not anchored")
	}
	checkTestJigsawAssembly(t, j, placed)
	single, _ := parseJigsaw("Tile 1:\n" + strings.Repeat("..........\n", 10))
	for _, budget := range []int{0, 1} {
		if _, err := single.assemble(budget); !errors.Is(err, errJigsawSearchLimit) {
			t.Fatal("limit must be distinct from no solution")
		}
	}
	if _, err := single.assemble(2); err != nil {
		t.Fatal(err)
	}
	impossible := testJigsawLabels([][4]byte{{1, 1, 1, 1}, {2, 2, 2, 2}, {3, 3, 3, 3}, {4, 4, 4, 4}})
	if _, err := impossible.assemble(jigsawWorkLimit); err == nil || errors.Is(err, errJigsawSearchLimit) {
		t.Fatal("impossible layout must exhaust choices, not report a work limit")
	}
}

func TestJigsawSmallAssembliesAgainstOracle(t *testing.T) {
	state := uint32(17)
	for trial := 0; trial < 160; trial++ {
		labels := make([][4]byte, 4)
		for i := range labels {
			for side := range labels[i] {
				state = state*1664525 + 1013904223
				labels[i][side] = byte(1 + (state>>16)%4)
			}
		}
		j := testJigsawLabels(labels)
		variants := make([][][100]byte, 4)
		for i, tile := range j.tiles {
			variants[i] = referenceJigsawVariants(tile.cells)
		}
		want := referenceJigsawCompletion(variants, nil, 0)
		placed, err := j.assemble(jigsawWorkLimit)
		if errors.Is(err, errJigsawSearchLimit) {
			t.Fatal("small oracle comparison hit limit")
		}
		if (err == nil) != want {
			t.Fatalf("trial %d: existence mismatch", trial)
		}
		if err == nil {
			checkTestJigsawAssembly(t, j, placed)
		}
	}
}

func blankTestJigsawImage(height, width int) jigsawImage {
	return jigsawImage{width: width, height: height, cells: []byte(strings.Repeat(".", height*width))}
}
func stampTestMonster(img jigsawImage, m jigsawMonster, row, col int) {
	for _, p := range m.pixels {
		img.cells[(row+p.row)*img.width+col+p.col] = '#'
	}
}

func TestJigsawMonsterBoundaries(t *testing.T) {
	monster := monsterOrientations()[0]
	for _, shape := range [][2]int{{3, 20}, {20, 20}, {23, 23}} {
		img := blankTestJigsawImage(shape[0], shape[1])
		stampTestMonster(img, monster, shape[0]-3, shape[1]-20)
		if got := jigsawRoughness(img); got != 0 {
			t.Fatalf("final valid position missed: %d", got)
		}
	}
	for _, shape := range [][2]int{{0, 0}, {2, 20}, {3, 19}, {1, 1}} {
		img := blankTestJigsawImage(shape[0], shape[1])
		if len(img.cells) > 0 {
			img.cells[0] = '#'
		}
		want := min(1, len(img.cells))
		if got := jigsawRoughness(img); got != want {
			t.Fatal("small-image bounds failure")
		}
	}
}

func TestJigsawMonsterOverlapsAndOrientations(t *testing.T) {
	variants := monsterOrientations()
	for _, monster := range variants {
		img := blankTestJigsawImage(monster.height, monster.width)
		stampTestMonster(img, monster, 0, 0)
		if got := jigsawRoughness(img); got != 0 {
			t.Fatal("missing monster orientation")
		}
	}
	img := blankTestJigsawImage(21, 21)
	stampTestMonster(img, variants[0], 7, 0)
	stampTestMonster(img, variants[0], 7, 1)
	img.cells[0] = '#'
	if got := jigsawRoughness(img); got != 1 {
		t.Fatalf("overlap double counted: %d", got)
	}
	img = blankTestJigsawImage(64, 64)
	stampTestMonster(img, variants[0], 0, 0)
	stampTestMonster(img, variants[1], 10, 40)
	stampTestMonster(img, variants[1], 35, 40)
	if got := jigsawRoughness(img); got != 15 {
		t.Fatalf("orientations must use greatest separate coverage: %d", got)
	}
}

func referenceMonsterCoverage(img jigsawImage) int {
	offsets := [3]string{"                  # ", "#    ##    ##    ###", " #  #  #  #  #  #   "}
	covered := make(map[int]bool)
	for row := 0; row+3 <= img.height; row++ {
		for col := 0; col+20 <= img.width; col++ {
			fits := true
			for r, line := range offsets {
				for c := range line {
					if line[c] == '#' && img.cells[(row+r)*img.width+col+c] != '#' {
						fits = false
					}
				}
			}
			if fits {
				for r, line := range offsets {
					for c := range line {
						if line[c] == '#' {
							covered[(row+r)*img.width+col+c] = true
						}
					}
				}
			}
		}
	}
	return len(covered)
}
func referenceRotateJigsawImage(img jigsawImage) jigsawImage {
	result := blankTestJigsawImage(img.width, img.height)
	for row := 0; row < img.height; row++ {
		for col := 0; col < img.width; col++ {
			result.cells[col*result.width+img.height-1-row] = img.cells[row*img.width+col]
		}
	}
	return result
}
func referenceFlipJigsawImage(img jigsawImage) jigsawImage {
	result := blankTestJigsawImage(img.height, img.width)
	for row := 0; row < img.height; row++ {
		for col := 0; col < img.width; col++ {
			result.cells[row*img.width+col] = img.cells[row*img.width+img.width-1-col]
		}
	}
	return result
}

func TestJigsawMonsterSearchAgainstImageTransformOracle(t *testing.T) {
	state := uint32(37)
	for trial := 0; trial < 48; trial++ {
		img := blankTestJigsawImage(24+trial%5, 25+trial%7)
		for i := range img.cells {
			state = state*1664525 + 1013904223
			if state>>30 != 0 {
				img.cells[i] = '#'
			}
		}
		variants := monsterOrientations()
		m := variants[trial%8]
		stampTestMonster(img, m, trial%(img.height-m.height+1), trial%(img.width-m.width+1))
		before := append([]byte(nil), img.cells...)
		best, total := 0, strings.Count(string(img.cells), "#")
		for flip := 0; flip < 2; flip++ {
			current := img
			if flip == 1 {
				current = referenceFlipJigsawImage(img)
			}
			for rotation := 0; rotation < 4; rotation++ {
				best = max(best, referenceMonsterCoverage(current))
				current = referenceRotateJigsawImage(current)
			}
		}
		if got := jigsawRoughness(img); got != total-best {
			t.Fatalf("trial %d: coverage mismatch", trial)
		}
		if !reflect.DeepEqual(img.cells, before) {
			t.Fatal("monster search mutated image")
		}
	}
}
