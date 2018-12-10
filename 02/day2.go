package main

import (
	"bufio"
	"errors"
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

// copy from my hamming exercise on exercism

// Distance calculates the hamming distance between two strings
func Distance(a, b string) (int, error) {
	r := []rune(a)
	s := []rune(b)

	if len(r) != len(s) {
		return -1, errors.New("string lengths do not match")
	}

	distance := 0
	for i := 0; i < len(r); i++ {
		if r[i] != s[i] {
			distance++
		}
	}

	return distance, nil
}

func main() {
	boxes := fileLines("input")

	// part 1
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

	// part 2
	// check each pair of box ids for a hamming distance of 1
	for i := 0; i < len(boxes); i++ {
		for j := i + 1; j < len(boxes); j++ {
			d, _ := Distance(boxes[i], boxes[j])
			if d == 1 {
				// create a new string with contains only the common chars
				answer := ""
				for k := 0; k < len(boxes[i]); k++ {
					if boxes[i][k] == boxes[j][k] {
						answer += string(boxes[i][k])
					}
				}
				fmt.Println(answer)
			}
		}
	}
}
