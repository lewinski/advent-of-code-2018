package main

import (
	"fmt"
)

// simple fixed-size fifo
type last6 struct {
	a, b, c, d, e, f int
}

func (l *last6) push(x int) {
	l.a, l.b, l.c, l.d, l.e, l.f = l.b, l.c, l.d, l.e, l.f, x
}

func (l last6) equals(x last6) bool {
	return l.a == x.a && l.b == x.b && l.c == x.c && l.d == x.d && l.e == x.e && l.f == x.f
}

func main() {
	input := 554401
	scores := []int{3, 7}
	positions := []int{0, 1}

	// part 1
	for {
		if len(scores) > input+10 {
			fmt.Println(scores[input : input+10])
			break
		}

		recipe := scores[positions[0]] + scores[positions[1]]
		if recipe >= 10 {
			scores = append(scores, 1, recipe-10)
		} else {
			scores = append(scores, recipe)
		}

		positions[0] = (positions[0] + 1 + scores[positions[0]]) % len(scores)
		positions[1] = (positions[1] + 1 + scores[positions[1]]) % len(scores)
	}

	// part 2
	target := last6{5, 5, 4, 4, 0, 1}
	last := last6{0, 0, 0, 0, 0, 0}
	scores = []int{3, 7}
	positions = []int{0, 1}
	for {
		recipe := scores[positions[0]] + scores[positions[1]]
		if recipe >= 10 {
			scores = append(scores, 1, recipe-10)
			last.push(1)
			if last.equals(target) {
				fmt.Println(len(scores) - 7)
				break
			}
			last.push(recipe - 10)
			if last.equals(target) {
				fmt.Println(len(scores) - 6)
				break
			}
		} else {
			scores = append(scores, recipe)
			last.push(recipe)
			if last.equals(target) {
				fmt.Println(len(scores) - 6)
				break
			}
		}

		positions[0] = (positions[0] + 1 + scores[positions[0]]) % len(scores)
		positions[1] = (positions[1] + 1 + scores[positions[1]]) % len(scores)
	}
}
