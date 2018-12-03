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

func main() {
	lines := fileLines("input")

	changes := make([]int, 0)
	for _, l := range lines {
		changes = append(changes, atoi(l))
	}

	freq := 0
	for _, i := range changes {
		freq += i
	}

	fmt.Println(freq)
}
