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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
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

type coord struct {
	x int
	y int
}

func distance(a, b coord) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}

func closest(p coord, c []coord) int {
	indexes := make([]int, 0)
	bestDistance := math.MaxInt32
	for i, x := range c {
		d := distance(p, x)
		if d < bestDistance {
			bestDistance = d
			indexes = []int{i}
		} else if d == bestDistance {
			indexes = append(indexes, i)
		}
	}
	if len(indexes) == 1 {
		return indexes[0]
	}
	return -1
}

func main() {
	// parse stuff
	file := "input"
	maxDistance := 10000
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	if len(os.Args) > 2 {
		maxDistance = atoi(os.Args[2])
	}
	lines := fileLines(file)

	coords := []coord{}
	minCoord := coord{math.MaxInt32, math.MaxInt32}
	maxCoord := coord{-1, -1}

	// parse the data file and find the min and max extent for the coordinates given
	for _, s := range lines {
		re := regexp.MustCompile("^(\\d+), (\\d+)")
		matches := re.FindStringSubmatch(s)

		x := atoi(matches[1])
		y := atoi(matches[2])
		c := coord{x, y}

		coords = append(coords, c)

		minCoord.x = min(minCoord.x, c.x)
		maxCoord.x = max(maxCoord.x, c.x)
		minCoord.y = min(minCoord.y, c.y)
		maxCoord.y = max(maxCoord.y, c.y)
	}

	// basically all of the edges are going to be infinite, so lets figure out which those are
	infinites := make(map[int]bool)
	for x := minCoord.x; x <= maxCoord.x; x++ {
		infinites[closest(coord{x, minCoord.y}, coords)] = true
		infinites[closest(coord{x, maxCoord.y}, coords)] = true
	}
	for y := minCoord.y; y <= maxCoord.y; y++ {
		infinites[closest(coord{minCoord.x, y}, coords)] = true
		infinites[closest(coord{maxCoord.x, y}, coords)] = true
	}

	// now go over all the points in the region and find the closest one to each point
	closestCount := make(map[int]int)
	for x := minCoord.x; x <= maxCoord.x; x++ {
		for y := minCoord.y; y <= maxCoord.y; y++ {
			c := closest(coord{x, y}, coords)
			closestCount[c]++
		}
	}

	// now search for the largest non-infinite region
	bestCount := 0
	for i, cnt := range closestCount {
		if infinites[i] {
			continue
		}
		if cnt > bestCount {
			bestCount = cnt
		}
	}
	fmt.Println(bestCount)

	// part 2
	pointsInRegion := 0
	for x := minCoord.x - maxDistance; x <= maxCoord.x+maxDistance; x++ {
		for y := minCoord.y - maxDistance; y <= maxCoord.y+maxDistance; y++ {
			totalDistance := 0
			for _, p := range coords {
				totalDistance += distance(coord{x, y}, p)
				if totalDistance >= maxDistance {
					break
				}
			}
			if totalDistance < maxDistance {
				pointsInRegion++
			}
		}
	}
	fmt.Println(pointsInRegion)
}
