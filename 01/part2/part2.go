package main

import (
	"bufio"
	"container/ring"
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

func main() {
	lines := fileLines("input")

	changes := make([]int, 0)
	for _, l := range lines {
		changes = append(changes, atoi(l))
	}

	// load change list into ring
	r := ring.New(len(changes))
	for _, i := range changes {
		r.Value = i
		r = r.Next()
	}

	// make note of frequencies we've seen
	freq := 0
	seen := make(map[int]bool)
	seen[freq] = true

	// walk the ring and apply changes until we get a duplicate
	for {
		value, _ := r.Value.(int)
		freq += value
		if seen[freq] {
			fmt.Println(freq)
			break
		} else {
			seen[freq] = true
		}
		r = r.Next()
	}
}
