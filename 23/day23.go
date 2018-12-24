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

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

type point struct {
	x, y, z int
}
type nanobot struct {
	pos    point
	radius int
}

func (n nanobot) reaches(p point) bool {
	return mdist(n.pos, p) <= n.radius
}

func mdist(a, b point) int {
	return abs(a.x-b.x) + abs(a.y-b.y) + abs(a.z-b.z)
}

func main() {
	// read file
	file := "input"
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	lines := fileLines(file)

	bots := []nanobot{}
	maxBot := nanobot{radius: -1}
	for _, l := range lines {
		n := nanobot{}
		fmt.Sscanf(l, "pos=<%d,%d,%d>, r=%d", &n.pos.x, &n.pos.y, &n.pos.z, &n.radius)
		bots = append(bots, n)
		if maxBot.radius < n.radius {
			maxBot = n
		}
	}

	botsInRange := 0
	for _, bot := range bots {
		distance := mdist(bot.pos, maxBot.pos)
		if distance <= maxBot.radius {
			botsInRange++
		}
	}
	fmt.Println("part 1 - bots in range of largest radius:", botsInRange)
}
