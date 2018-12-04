package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
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
	// parse stuff
	file := "input"
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	lines := fileLines(file)

	startShift := regexp.MustCompile("^\\[\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}\\] Guard #(\\d+) begins shift$")
	fallsAsleep := regexp.MustCompile("^\\[\\d{4}-\\d{2}-\\d{2} \\d{2}:(\\d{2})\\] falls asleep$")
	wakesUp := regexp.MustCompile("^\\[\\d{4}-\\d{2}-\\d{2} \\d{2}:(\\d{2})\\] wakes up$")

	// map guard -> number of minutes asleep
	guardTotalSleep := map[int]int{}

	// map guard -> minute -> number of times asleep
	guardSleeping := map[int]map[int]int{}

	// current guard and when they fell asleep
	var guard int
	var asleep int
	for _, l := range lines {
		if startShift.MatchString(l) {
			matches := startShift.FindStringSubmatch(l)
			guard = atoi(matches[1])
			asleep = 0
			if _, ok := guardSleeping[guard]; !ok {
				guardSleeping[guard] = make(map[int]int)
			}
		} else if fallsAsleep.MatchString(l) {
			matches := fallsAsleep.FindStringSubmatch(l)
			asleep = atoi(matches[1])
		} else if wakesUp.MatchString(l) {
			matches := wakesUp.FindStringSubmatch(l)
			awake := atoi(matches[1])

			guardTotalSleep[guard] += awake - asleep - 1

			for i := asleep; i < awake; i++ {
				guardSleeping[guard][i]++
			}
		}
	}

	// part 1
	// find guard that has most sleep minutes
	maxSleep := 0
	maxGuard := 0
	for guard, sleep := range guardTotalSleep {
		if sleep > maxSleep {
			maxSleep = sleep
			maxGuard = guard
		}
	}
	// find minute that guard was asleep most frequently
	maxMinute := 0
	maxAmount := 0
	for minute, amount := range guardSleeping[maxGuard] {
		if amount > maxAmount {
			maxAmount = amount
			maxMinute = minute
		}
	}
	// print out answer
	fmt.Println(maxGuard * maxMinute)

	// part 2
	// find largest minute across all guard/minute combos
	maxMinute2 := 0
	maxGuard2 := 0
	maxAmount2 := 0
	for guard := range guardSleeping {
		for minute, amount := range guardSleeping[guard] {
			if amount > maxAmount2 {
				maxMinute2 = minute
				maxGuard2 = guard
				maxAmount2 = amount
			}
		}
	}
	// print out answer
	fmt.Println(maxGuard2 * maxMinute2)
}
