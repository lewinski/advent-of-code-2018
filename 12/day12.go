package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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

func score(state string, start int) (s int) {
	for i, c := range state {
		if c == '#' {
			s += i - start
		}
	}
	return
}

func main() {
	// read file
	file := "input"
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	lines := fileLines(file)

	// initial state is in the first line
	initial := lines[0][15:]

	// rules start in the third line
	rules := map[string]string{}
	for _, l := range lines[2:] {
		rules[l[0:5]] = string(l[9])
	}

	// part 1
	// just run 20 generations and score it at the end
	start := 0
	state := initial
	for gen := 0; gen < 20; gen++ {
		// add padding around state and adjust the offset where zero is in our string
		// subtract 2 because when we map the next state, we lose two positions
		state = "....." + state + "....."
		start += 5 - 2

		// map next generation through rules
		var next strings.Builder
		for i := 2; i < len(state)-2; i++ {
			if res, ok := rules[state[i-2:i+3]]; ok {
				next.WriteString(res)
			} else {
				next.WriteString(".")
			}
		}

		// trim up the string and do any adjustments on the starting position
		state = strings.TrimRight(next.String(), ".")
		trim := strings.TrimLeft(state, ".")
		start -= len(state) - len(trim)
		state = trim
	}
	fmt.Println(score(state, start))

	// part 2
	// tons of generations, so we see if we can exit the loop early
	start = 0
	state = initial
	gens := 50000000000
	for gen := 0; gen < gens; gen++ {
		istate, istart := state, start

		// add padding around state and adjust the offset where zero is in our string
		// subtract 2 because when we map the next state, we lose two positions
		state = "....." + state + "....."
		start += 5 - 2

		// map next generation through rules
		var next strings.Builder
		for i := 2; i < len(state)-2; i++ {
			if res, ok := rules[state[i-2:i+3]]; ok {
				next.WriteString(res)
			} else {
				next.WriteString(".")
			}
		}

		// trim up the string and do any adjustments on the starting position
		state = strings.TrimRight(next.String(), ".")
		trim := strings.TrimLeft(state, ".")
		start -= len(state) - len(trim)
		state = trim

		if state == istate {
			// if we have the same state as we started with, then we just moved
			// left or right some number of positions and will do the same next
			// generation. as a shortcut, we can adjust our position by the
			// delta times the number of remaining generations and exit the loop early
			start += (start - istart) * (gens - gen - 1)
			break
		}
	}

	fmt.Println(score(state, start))
}
