package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func atoi(x string) int {
	i, _ := strconv.Atoi(x)
	return i
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func fileLines(path string) []string {
	f, err := os.Open(path)
	check(err)

	lines := make([]string, 0)
	s := bufio.NewScanner(f)
	for s.Scan() {
		lines = append(lines, s.Text())
	}

	return lines
}

func iabs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func imax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type point struct {
	x, y int
}

func (p point) offset(o point) point {
	return point{p.x + o.x, p.y + o.y}
}

type game [][]rune

func newGame(width, height int) game {
	g := make([][]rune, height)
	for y := range g {
		g[y] = make([]rune, width)
		for x := range g[y] {
			g[y][x] = OPEN
		}
	}
	return g
}

func (g game) get(p point) rune {
	return g[p.y][p.x]
}

func (g game) set(p point, r rune) {
	g[p.y][p.x] = r
}

func (g game) contains(p point) bool {
	if p.x < 0 || p.x >= len(g) {
		return false
	}
	if p.y < 0 || p.y >= len(g[0]) {
		return false
	}
	return true
}

func (g game) adjacent(p point) []rune {
	options := []point{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}
	rs := []rune{}
	for _, o := range options {
		if g.contains(p.offset(o)) {
			rs = append(rs, g.get(p.offset(o)))
		}
	}
	return rs
}

func (g game) String() string {
	var sb strings.Builder
	for y := range g {
		sb.WriteString(string(g[y]))
		sb.WriteRune('\n')
	}
	return sb.String()
}

func (g game) score() int {
	trees := 0
	lumberyards := 0
	for y := range g {
		for x := range g[y] {
			switch g.get(point{x, y}) {
			case TREES:
				trees++
				break
			case LUMBERYARD:
				lumberyards++
				break
			}
		}
	}
	return trees * lumberyards
}

const (
	OPEN       = '.'
	TREES      = '|'
	LUMBERYARD = '#'
)

func main() {
	// read file
	file := "input"
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	lines := fileLines(file)

	// set up the game field
	height := len(lines)
	width := len(lines[0])
	field := newGame(width, height)
	for y := range lines {
		for x, r := range lines[y] {
			field.set(point{x, y}, r)
		}
	}

	// set up cycle detection for part 2
	maxIter := 1000000000
	seen := make(map[string]int)
	scores := make(map[int]int)
	for i := 0; i < maxIter; i++ {

		// for part 1 we just need to know the state of the world after 10 iterations
		if i == 10 {
			fmt.Println("part 1:")
			fmt.Print(field.String())
			fmt.Println("score =", field.score())
		}

		// check if we've previously seen this state
		st := field.String()
		if iter, ok := seen[st]; ok {
			fmt.Println("part 2:")
			sameIter := iter + ((maxIter - iter) % (i - iter))
			fmt.Print(field.String())
			fmt.Println("score =", scores[sameIter])
			break
		}
		seen[field.String()] = i
		scores[i] = field.score()

		// generate the next state
		next := newGame(width, height)
		for y := range field {
			for x := range field[y] {
				p := point{x, y}
				r := field.get(p)
				adj := field.adjacent(p)

				// count trees and lumberyards for rules below
				trees := 0
				lumberyards := 0
				for _, a := range adj {
					switch a {
					case TREES:
						trees++
						break
					case LUMBERYARD:
						lumberyards++
						break
					}
				}

				switch r {
				case OPEN:
					if trees >= 3 {
						next.set(p, TREES)
					} else {
						next.set(p, r)
					}
					break
				case TREES:
					if lumberyards >= 3 {
						next.set(p, LUMBERYARD)
					} else {
						next.set(p, r)
					}
					break
				case LUMBERYARD:
					if lumberyards > 0 && trees > 0 {
						next.set(p, LUMBERYARD)
					} else {
						next.set(p, OPEN)
					}
					break
				}
			}
		}

		field = next
	}

}
