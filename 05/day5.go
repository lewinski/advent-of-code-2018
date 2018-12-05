package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
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

var reactions = []string{
	"aA", "Aa",
	"bB", "Bb",
	"cC", "Cc",
	"dD", "Dd",
	"eE", "Ee",
	"fF", "Ff",
	"gG", "Gg",
	"hH", "Hh",
	"iI", "Ii",
	"jJ", "Jj",
	"kK", "Kk",
	"lL", "Ll",
	"mM", "Mm",
	"nN", "Nn",
	"oO", "Oo",
	"pP", "Pp",
	"qQ", "Qq",
	"rR", "Rr",
	"sS", "Ss",
	"tT", "Tt",
	"uU", "Uu",
	"vV", "Vv",
	"wW", "Ww",
	"xX", "Xx",
	"yY", "Yy",
	"zZ", "Zz",
}

func main() {
	// parse stuff
	file := "input"
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	lines := fileLines(file)

	polymer := lines[0]

	// part1
	reduced := reducePolymer(polymer)
	fmt.Println(len(reduced))

	// part2
	minLen := len(polymer)
	for i := 0; i < len(reactions); i += 2 {
		thisLen := len(reducePolymer(regexp.MustCompile(fmt.Sprintf("[%s]", reactions[i])).ReplaceAllLiteralString(polymer, "")))
		if thisLen < minLen {
			minLen = thisLen
		}
	}
	fmt.Println(minLen)
}

func reducePolymer(polymer string) string {
	for {
		reacted := polymer
		for _, r := range reactions {
			reacted = strings.Replace(reacted, r, "", -1)
		}
		if reacted == polymer {
			return polymer
		}
		polymer = reacted
	}
}
