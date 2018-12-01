package main

import (
	"bufio"
	"container/ring"
	"fmt"
	"os"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// read whole file into change list
	f, err := os.Open("input")
	check(err)

	changes := make([]int, 0)
	s := bufio.NewScanner(f)
	for s.Scan() {
		t := s.Text()
		i, err := strconv.Atoi(t)
		check(err)
		changes = append(changes, i)
	}

	// load change list into ring
	r := ring.New(len(changes))
	for i := 0; i < len(changes); i++ {
		r.Value = changes[i]
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
