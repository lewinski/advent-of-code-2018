package main

import (
	"bufio"
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
	freq := 0

	f, err := os.Open("input")
	check(err)

	s := bufio.NewScanner(f)
	for s.Scan() {
		t := s.Text()
		i, err := strconv.Atoi(t)
		check(err)
		freq += i
	}

	fmt.Println(freq)
}
