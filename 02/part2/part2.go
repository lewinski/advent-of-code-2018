package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

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

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// read whole file into box id list
	f, err := os.Open("input")
	check(err)

	boxes := make([]string, 0)
	s := bufio.NewScanner(f)
	for s.Scan() {
		boxes = append(boxes, s.Text())
	}

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
