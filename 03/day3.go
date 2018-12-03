package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type claim struct {
	id int
	x  int
	y  int
	w  int
	h  int
}

func main() {
	// read whole file into box id list
	f, err := os.Open("input")
	check(err)

	lines := make([]string, 0)
	s := bufio.NewScanner(f)
	for s.Scan() {
		lines = append(lines, s.Text())
	}

	claims := make([]claim, 0)

	// parse stuff
	for _, s := range lines {
		re := regexp.MustCompile("^#(\\d+) @ (\\d+),(\\d+): (\\d+)x(\\d+)$")
		matches := re.FindStringSubmatch(s)

		id, _ := strconv.Atoi(matches[1])
		x, _ := strconv.Atoi(matches[2])
		y, _ := strconv.Atoi(matches[3])
		w, _ := strconv.Atoi(matches[4])
		h, _ := strconv.Atoi(matches[5])

		claims = append(claims, claim{id, x, y, w, h})
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
