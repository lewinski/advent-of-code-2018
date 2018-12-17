package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
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

func (p point) below() point {
	return point{p.x, p.y + 1}
}

func (p point) left() point {
	return point{p.x - 1, p.y}
}

func (p point) right() point {
	return point{p.x + 1, p.y}
}

type field struct {
	startX int
	data   [][]rune
}

func newField(minX, minY, maxX, maxY int) field {
	// expand the bounds to ensure the border is free for water overflows
	minX--
	maxX++

	f := field{startX: minX}
	f.data = make([][]rune, maxY-minY+1)
	for y := range f.data {
		f.data[y] = make([]rune, maxX-minX+1)
		for x := range f.data[y] {
			f.data[y][x] = SAND
		}
	}
	return f
}

func (f field) contains(p point) bool {
	if p.y >= len(f.data) {
		return false
	}
	if (p.x - f.startX) >= len(f.data[p.y]) {
		return false
	}
	return true
}

func (f field) get(p point) rune {
	return f.data[p.y][p.x-f.startX]
}

func (f *field) set(p point, r rune) {
	f.data[p.y][p.x-f.startX] = r
}

func (f field) print() {
	for y := range f.data {
		fmt.Println(string(f.data[y]))
	}
}

const (
	CLAY   rune = '#'
	SAND   rune = '.'
	SPRING rune = '+'
	WATER  rune = '~'
	FLOW   rune = '|'
)

func main() {
	// read file
	file := "input"
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	lines := fileLines(file)

	// from the puzzle description
	spring := point{500, 0}

	// parse all of the input scan data
	re := regexp.MustCompile("^([xy])=(\\d+), ([xy])=(\\d+)..(\\d+)")
	clay := []point{}
	for _, l := range lines {
		matches := re.FindStringSubmatch(l)
		for i := atoi(matches[4]); i <= atoi(matches[5]); i++ {
			if matches[1] == "x" {
				clay = append(clay, point{atoi(matches[2]), i})
			} else {
				clay = append(clay, point{i, atoi(matches[2])})
			}
		}
	}

	// find the bounds of the clay we surveyed
	minX, minY := math.MaxInt32, math.MaxInt32
	maxX, maxY := math.MinInt32, math.MinInt32
	for _, p := range clay {
		minX = imin(minX, p.x)
		maxX = imax(maxX, p.x)
		minY = imin(minY, p.y)
		maxY = imax(maxY, p.y)
	}

	// set up the initial play state
	playField := newField(minX, 0, maxX, maxY)
	for _, p := range clay {
		playField.set(p, CLAY)
	}
	playField.set(spring, SPRING)

	// start with the spring as the only water source we need to resolve
	sources := []point{spring}
	for len(sources) > 0 {
		// pop the last source off the stack
		start := sources[len(sources)-1]
		sources = sources[:len(sources)-1]

		// if this source has been consumed by static water, ignore it
		if playField.get(start) == WATER {
			continue
		}

		// fill down from the source
		below := start
		for playField.contains(below.below()) && (playField.get(below.below()) == SAND || playField.get(below.below()) == FLOW) {
			below = below.below()
			playField.set(below, FLOW)
		}

		// if we've left the bottom of the field, we're done with that source
		if !playField.contains(below.below()) {
			continue
		}

		// if we hit clay or static water, fill across from where the down flow stopped
		if playField.get(below.below()) == CLAY || playField.get(below.below()) == WATER {
			left := below.left()
			for playField.get(left) != CLAY && (playField.get(left.below()) == CLAY || playField.get(left.below()) == WATER) {
				playField.set(left, FLOW)
				left = left.left()
			}
			right := below.right()
			for playField.get(right) != CLAY && (playField.get(right.below()) == CLAY || playField.get(right.below()) == WATER) {
				playField.set(right, FLOW)
				right = right.right()
			}
			// if this level is between two clay spots, turn the flowing water into static water and redo the same source again
			if playField.get(left) == CLAY && playField.get(right) == CLAY {
				left = left.right()
				for left != right {
					playField.set(left, WATER)
					left = left.right()
				}
				sources = append(sources, start)
			}
			// if either side is not a clay spot, add flowing water and a new source to evaluate
			if playField.get(left) == SAND {
				playField.set(left, FLOW)
				sources = append(sources, start, left)
			}
			if playField.get(right) == SAND {
				playField.set(right, FLOW)
				sources = append(sources, start, right)
			}
		}
	}

	// count flowing and static water
	cntFlow, cntWater := 0, 0
	for y := minY; y <= maxY; y++ {
		for x := range playField.data[y] {
			if playField.data[y][x] == FLOW {
				cntFlow++
			}
			if playField.data[y][x] == WATER {
				cntWater++
			}
		}
	}

	fmt.Println("Final state:")
	playField.print()
	fmt.Printf("flow = %d, water = %d, flow + water = %d\n", cntFlow, cntWater, cntFlow+cntWater)
}
