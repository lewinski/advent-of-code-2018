package main

import (
	"bufio"
	"fmt"
	"os"
)

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
	boxes := fileLines("input")

	// counters for how many inputs had a two or three
	twos := 0
	threes := 0

	// check each box
	for i := 0; i < len(boxes); i++ {

		// create frequency count
		counts := make(map[rune]int)
		for _, c := range boxes[i] {
			counts[c]++
		}

		// check if this is a two count or three count
		two := 0
		three := 0
		for _, n := range counts {
			if n == 2 {
				two = 1
			}
			if n == 3 {
				three = 1
			}
		}

		// update final answer for this box id
		twos += two
		threes += three
	}

	fmt.Println(twos * threes)
}
