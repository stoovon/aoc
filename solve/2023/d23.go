package solve2023

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day23 struct {
}

func (d Day23) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2023, Day: 23}
}

type day23Grid struct {
	data   []byte
	width  int
	height int
}

var arrows = []byte{'^', '>', 'v', '<'}

func newDay23Grid(input string) *day23Grid {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	width := len(lines[0])
	height := len(lines)
	data := make([]byte, 0, width*height)
	for _, l := range lines {
		data = append(data, l...)
	}
	return &day23Grid{data, width, height}
}

func (g *day23Grid) nextPos(p int, dir int) (int, bool) {
	switch dir {
	case 0: // up
		if p >= g.width {
			return p - g.width, true
		}
	case 1: // right
		if (p+1)%g.width != 0 {
			return p + 1, true
		}
	case 2: // down
		if p < len(g.data)-g.width {
			return p + g.width, true
		}
	case 3: // left
		if p%g.width != 0 {
			return p - 1, true
		}
	}
	return 0, false
}

func (g *day23Grid) intersectionsAndDistances(slopes bool) [][][2]int {
	start, end := -1, -1
	for i := 0; i < g.width; i++ {
		if g.data[i] == '.' {
			start = i
			break
		}
	}
	for i := len(g.data) - g.width; i < len(g.data); i++ {
		if g.data[i] == '.' {
			end = i
			break
		}
	}
	vpts := []int{start, end}
	for i, c := range g.data {
		if c == '#' {
			continue
		}
		f := 0
		for d := 0; d < 4; d++ {
			if np, ok := g.nextPos(i, d); ok && g.data[np] != '#' {
				f++
			}
		}
		if f > 2 {
			vpts = append(vpts, i)
		}
	}
	// sort and unique
	m := map[int]struct{}{}
	for _, v := range vpts {
		m[v] = struct{}{}
	}
	vpts = vpts[:0]
	for k := range m {
		vpts = append(vpts, k)
	}
	// sort
	for i := 0; i < len(vpts); i++ {
		for j := i + 1; j < len(vpts); j++ {
			if vpts[j] < vpts[i] {
				vpts[i], vpts[j] = vpts[j], vpts[i]
			}
		}
	}
	rmap := make([][][2]int, len(vpts))
	for id, ii := range vpts {
		iv := &rmap[id]
		visited := make([]bool, len(g.data))
		visited[ii] = true
		type pair struct{ i, d int }
		q := []pair{{ii, 0}}
		for len(q) > 0 {
			cur := q[len(q)-1]
			q = q[:len(q)-1]
			ai := 4
			for k, a := range arrows {
				if g.data[cur.i] == a {
					ai = k
					break
				}
			}
			for di := 0; di < 4; di++ {
				if slopes && ai < 4 && ai != di {
					continue
				}
				np, ok := g.nextPos(cur.i, di)
				if !ok || g.data[np] == '#' {
					continue
				}
				if !visited[np] {
					visited[np] = true
					idx := -1
					for k, v := range vpts {
						if v == np {
							idx = k
							break
						}
					}
					if idx == -1 {
						q = append(q, pair{np, cur.d + 1})
					} else {
						*iv = append(*iv, [2]int{idx, cur.d + 1})
					}
				}
			}
		}
	}
	return rmap
}

func dfs(i, d int, rmap [][][2]int, visited []bool, end int, r *int) {
	visited[i] = true
	if i == end {
		if d > *r {
			*r = d
		}
	}
	for _, e := range rmap[i] {
		j, md := e[0], e[1]
		if !visited[j] {
			dfs(j, d+md, rmap, visited, end, r)
		}
	}
	visited[i] = false
}

func (d Day23) Part1(input string) (string, error) {
	grid := newDay23Grid(input)
	rmap := grid.intersectionsAndDistances(true)
	visited := make([]bool, len(rmap))
	r := 0
	dfs(0, 0, rmap, visited, len(rmap)-1, &r)
	return strconv.Itoa(r), nil
}

func (d Day23) Part2(input string) (string, error) {
	grid := newDay23Grid(input)
	rmap := grid.intersectionsAndDistances(false)
	visited := make([]bool, len(rmap))
	r := 0
	dfs(0, 0, rmap, visited, len(rmap)-1, &r)
	return strconv.Itoa(r), nil
}

func init() {
	solve.Register(Day23{})
}
