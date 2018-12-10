package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
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

func fileWords(path string) []string {
	f, err := os.Open(path)
	check(err)

	words := make([]string, 0)
	s := bufio.NewScanner(f)
	s.Split(bufio.ScanWords)
	for s.Scan() {
		words = append(words, s.Text())
	}

	return words
}

func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func imax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type point struct {
	x int
	y int
}
type light struct {
	pos point
	vel point
}

// advance time, moving the points
func tick(l light) light {
	return light{point{l.pos.x + l.vel.x, l.pos.y + l.vel.y}, l.vel}
}

func mapLights(vl []light, f func(light) light) []light {
	vlm := make([]light, len(vl))
	for i, v := range vl {
		vlm[i] = f(v)
	}
	return vlm
}

// calculates rectangle that covers points
func extents(lights []light) (min point, max point) {
	min = lights[0].pos
	max = lights[0].pos
	for _, l := range lights {
		min.x = imin(min.x, l.pos.x)
		min.y = imin(min.y, l.pos.y)
		max.x = imax(max.x, l.pos.x)
		max.y = imax(max.y, l.pos.y)
	}
	return min, max
}

// area of the covering rectangle
func covering(lights []light) int {
	min, max := extents(lights)
	return (max.x - min.x) * (max.y - min.y)
}

// print out a light pattern
func printLights(lights []light) {
	min, max := extents(lights)

	lines := make([]string, max.y-min.y+1)
	for i := range lines {
		lines[i] = strings.Repeat(" ", max.x-min.x+1)
	}

	for _, l := range lights {
		r := []rune(lines[l.pos.y-min.y])
		r[l.pos.x-min.x] = '*'
		lines[l.pos.y-min.y] = string(r)
	}

	for _, l := range lines {
		fmt.Println(l)
	}
}

func main() {
	// read file
	file := "input"
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	lines := fileLines(file)

	// parse input data
	lights := []light{}
	re := regexp.MustCompile("position=<\\s*(-?\\d+),\\s*(-?\\d+)> velocity=<\\s*(-?\\d+),\\s*(-?\\d+)>")
	for _, l := range lines {
		matches := re.FindStringSubmatch(l)
		x := atoi(matches[1])
		y := atoi(matches[2])
		xv := atoi(matches[3])
		yv := atoi(matches[4])
		lights = append(lights, light{point{x, y}, point{xv, yv}})
	}

	// tick time forward until we seem to reach a minimum area that the points live in
	// hopefully that is the answer!
	seconds := 0
	for {
		nextLights := mapLights(lights, tick)
		if covering(nextLights) < covering(lights) {
			lights = nextLights
			seconds++
		} else {
			break
		}
	}

	// print out the pattern and how many ticks we took to get there
	printLights(lights)
	fmt.Println(seconds)
}
