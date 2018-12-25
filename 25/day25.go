package main

import (
	"bufio"
	"fmt"
	"os"
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

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

type point struct {
	x, y, z, t int
}

func (p point) mdist(o point) int {
	return abs(p.x-o.x) + abs(p.y-o.y) + abs(p.z-o.z) + abs(p.t-o.t)
}

func main() {
	// read file
	file := "input"
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	lines := fileLines(file)

	// parse data
	fixedPoints := []point{}
	for _, l := range lines {
		var p point
		fmt.Sscanf(l, "%d,%d,%d,%d", &p.x, &p.y, &p.z, &p.t)
		fixedPoints = append(fixedPoints, p)
	}

	// has this point been used?
	used := map[point]bool{}

	// all the constellations
	constellations := [][]point{}

	// until we've assigned every point to a constellation
	for len(used) != len(fixedPoints) {
		// find the first unassigned point
		var seed point
		for _, p := range fixedPoints {
			if _, ok := used[p]; !ok {
				seed = p
				break
			}
		}

		// assign it to a new constellation
		constellation := []point{seed}
		used[seed] = true

		// find other points that are close, and then assign those points
		// and redo the search for those new ones as well
		stack := []point{seed}
		for len(stack) > 0 {
			p := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			for _, o := range fixedPoints {
				if _, ok := used[o]; ok {
					continue
				}
				if p.mdist(o) <= 3 {
					constellation = append(constellation, o)
					used[o] = true
					stack = append(stack, o)
				}
			}

		}

		// constellation has grown all it can so add it to the list
		constellations = append(constellations, constellation)
	}

	fmt.Println("number of constellations:", len(constellations))
}
