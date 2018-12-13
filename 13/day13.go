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

func fileWords(path string) []string {
	f, err := os.Open(path)
	check(err)

	words := make([]string, 0)
	s := bufio.NewScanner(f)
	s.Split(bufio.ScanWords)
	for s.Scan() {
		words = append(words, s.Text())
	}

	return words
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

type facing int

const (
	invalid facing = iota
	up
	right
	down
	left
)

func (f facing) String() string {
	switch f {
	case up:
		return "^"
	case right:
		return ">"
	case down:
		return "v"
	case left:
		return "<"
	}
	return "*"
}

func toFacing(r rune) (facing, bool) {
	switch r {
	case '^':
		return up, true
	case '>':
		return right, true
	case 'v':
		return down, true
	case '<':
		return left, true
	}
	return invalid, false
}

type cart struct {
	turns     int
	direction facing
}

func (c *cart) handleTrack(r rune) {
	switch r {
	case '+':
		turns := []map[facing]facing{
			{up: left, right: up, down: right, left: down},
			{up: up, right: right, down: down, left: left},
			{up: right, right: down, down: left, left: up},
		}
		c.direction = turns[c.turns%3][c.direction]
		c.turns++
		break
	case '/':
		turns := map[facing]facing{up: right, right: up, down: left, left: down}
		c.direction = turns[c.direction]
		break
	case '\\':
		turns := map[facing]facing{up: left, right: down, down: right, left: up}
		c.direction = turns[c.direction]
	}
}

func nextPosition(x int, y int, direction facing) (int, int) {
	switch direction {
	case up:
		return x - 1, y
	case right:
		return x, y + 1
	case down:
		return x + 1, y
	case left:
		return x, y - 1
	}
	return x, y
}

func printCarts(tracks []string, carts [][]cart) {
	output := make([]string, len(tracks))
	copy(output, tracks)
	for i := range carts {
		for j, c := range carts[i] {
			if c.direction != invalid {
				output[i] = output[i][:j] + c.direction.String() + output[i][j+1:]
			}
		}
	}
	for _, s := range output {
		fmt.Println(s)
	}
}

func main() {
	// read file
	file := "input"
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	lines := fileLines(file)

	tracks := make([]string, len(lines))
	copy(tracks, lines)

	carts := make([][]cart, len(tracks))
	for i := range carts {
		carts[i] = make([]cart, len(tracks[i]))
	}

	for i := range tracks {
		for j := range tracks[i] {
			if f, ok := toFacing([]rune(tracks[i])[j]); ok {
				carts[i][j] = cart{0, f}
				if f == up || f == down {
					tracks[i] = tracks[i][:j] + "|" + tracks[i][j+1:]
				} else {
					tracks[i] = tracks[i][:j] + "-" + tracks[i][j+1:]
				}
			}
		}
	}

	for {
		// where we will put the next generation
		newCarts := make([][]cart, len(carts))
		for i := range newCarts {
			newCarts[i] = make([]cart, len(carts[i]))
		}

		for i := range tracks {
			for j := range tracks[i] {
				c := carts[i][j]
				// if there is a cart in this position
				if c.direction != invalid {
					// move the cart out of this position
					carts[i][j].direction = invalid

					// figure out the next position
					a, b := nextPosition(i, j, c.direction)

					// if there is a cart in this position already, we got a collision
					if (newCarts[a][b].direction != invalid) || (carts[a][b].direction != invalid) {
						fmt.Println("collision at", b, a)
						// delete both carts
						carts[a][b].direction = invalid
						newCarts[a][b].direction = invalid
					} else {
						// move the cart into the new position
						newCarts[a][b] = c
						newCarts[a][b].handleTrack([]rune(tracks[a])[b])
					}
				}
			}
		}
		carts = newCarts

		// printCarts(tracks, carts)

		// see if there is only one cart left
		totalCarts, x, y := 0, 0, 0
		for i := range carts {
			for j := range carts[i] {
				if carts[i][j].direction != invalid {
					totalCarts++
					x, y = j, i
				}
			}
		}
		if totalCarts <= 1 {
			fmt.Println("last cart is at", x, y)
			break
		}
	}
}
