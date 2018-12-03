package main

import (
	"bufio"
	"fmt"
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

type claim struct {
	id int
	x  int
	y  int
	w  int
	h  int
}

func main() {
	// parse stuff
	lines := fileLines("input")
	claims := make([]claim, 0)
	for _, s := range lines {
		re := regexp.MustCompile("^#(\\d+) @ (\\d+),(\\d+): (\\d+)x(\\d+)$")
		matches := re.FindStringSubmatch(s)

		var c claim
		c.id = atoi(matches[1])
		c.x = atoi(matches[2])
		c.y = atoi(matches[3])
		c.w = atoi(matches[4])
		c.h = atoi(matches[5])

		claims = append(claims, c)
	}

	// make a map of claims
	var fabric [1000][1000]int
	for _, c := range claims {
		for i := 0; i < c.w; i++ {
			for j := 0; j < c.h; j++ {
				fabric[c.x+i][c.y+j]++
			}
		}
	}

	// part 1
	c := 0
	for i := range fabric {
		for j := range fabric[i] {
			if fabric[i][j] > 1 {
				c++
			}
		}
	}
	fmt.Println(c)

	// part 2
NextClaim:
	for _, c := range claims {
		for i := 0; i < c.w; i++ {
			for j := 0; j < c.h; j++ {
				if fabric[c.x+i][c.y+j] > 1 {
					continue NextClaim
				}
			}
		}
		fmt.Println(c.id)
	}
}
