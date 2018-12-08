package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
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

func findAvailableTask(tasks []string, requirements map[string][]string, done map[string]bool, assigned map[string]bool) string {
NextStep:
	for _, k := range tasks {
		if done[k] {
			continue
		}
		if a, ok := assigned[k]; ok && a {
			continue
		}
		req, ok := requirements[k]
		if !ok {
			return k
		}
		for _, x := range req {
			if !done[x] {
				continue NextStep
			}
		}
		return k
	}
	return "."
}

func main() {
	// parse stuff
	file := "input"
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	lines := fileLines(file)

	// parse data file
	steps := map[string]bool{}
	requirements := map[string][]string{}
	re := regexp.MustCompile("Step (.) must be finished before step (.) can begin.")
	for _, l := range lines {
		matches := re.FindStringSubmatch(l)
		requirements[matches[2]] = append(requirements[matches[2]], matches[1])
		steps[matches[1]] = true
		steps[matches[2]] = true
	}

	// get all the step names
	keys := make([]string, len(steps))
	i := 0
	for k := range steps {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	done := map[string]bool{}
	for _, k := range keys {
		done[k] = false
	}

	// part 1
	order := make([]string, len(keys))
	for i := 0; i < len(order); i++ {
		order[i] = findAvailableTask(keys, requirements, done, map[string]bool{})
		done[order[i]] = true
	}
	fmt.Println(strings.Join(order, ""))

	// part 2
	workers := 5
	minTask := 60

	for _, k := range keys {
		done[k] = false
	}
	order2 := []string{}
	available := make([]int, workers)
	doneAt := map[string]int{}
	current := 0

	// while all the tasks aren't completed
	for len(order2) != len(keys) {

		// find an available worker
		for w := 0; w < workers; w++ {
			if available[w] <= current {

				// find all the currently assigned tasks
				assigned := map[string]bool{}
				for task, t := range doneAt {
					if t >= current {
						assigned[task] = true
					}
				}

				// find a task to assign to the worker
				task := findAvailableTask(keys, requirements, done, assigned)
				if task != "." {
					// assign the task and update when the worker will be available again
					available[w] = current + minTask + int(task[0]-'A'+1)
					doneAt[task] = available[w]
				}
			}
		}

		// advance time
		current++
		for task, t := range doneAt {
			// complete any tasks
			if t == current {
				done[task] = true
				order2 = append(order2, task)
			}
		}
	}
	fmt.Println(current)
}
