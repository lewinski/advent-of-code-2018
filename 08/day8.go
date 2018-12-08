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

type node struct {
	children []node
	metadata []int
}

// value as defined by part 2
func (n *node) value() (v int) {
	if len(n.children) == 0 {
		for _, m := range n.metadata {
			v += m
		}
		return v
	}

	for _, m := range n.metadata {
		m--
		if m >= 0 && m < len(n.children) {
			v += n.children[m].value()
		}
	}
	return
}

func (n *node) print() {
	fmt.Printf("node{[%d]node, %v}\n", len(n.children), n.metadata)
}

func parseNode(numbers []string) (node, int) {
	used := 0

	childrenCount := atoi(numbers[0])
	metadataCount := atoi(numbers[1])
	used += 2

	children := make([]node, childrenCount)
	for i := 0; i < childrenCount; i++ {
		n, u := parseNode(numbers[used:])
		children[i] = n
		used += u
	}

	metadata := make([]int, metadataCount)
	for i := 0; i < metadataCount; i++ {
		metadata[i] = atoi(numbers[used])
		used++
	}

	return node{children, metadata}, used
}

func main() {
	// parse stuff
	file := "input"
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	words := fileWords(file)

	tree, _ := parseNode(words)

	// part 1
	stack := []node{tree}
	license := 0
	for len(stack) != 0 {
		// pop a node
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		// append node's children
		stack = append(stack, n.children...)
		// sum metadata
		for _, m := range n.metadata {
			license += m
		}
	}
	fmt.Println(license)

	// part 2
	fmt.Println(tree.value())
}
