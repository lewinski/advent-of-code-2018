package main

import (
	"bufio"
	"container/list"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

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

func iabs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func imax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type point struct {
	x, y int
}

func (p point) distanceTo(o point) int {
	return iabs(p.x-o.x) + iabs(p.y-o.y)
}

func (p point) offset(o point) point {
	return point{p.x + o.x, p.y + o.y}
}

type unit struct {
	id       int
	elf      bool
	position point
	health   int
	attack   int
}

func (u unit) Symbol() rune {
	if u.elf {
		return 'E'
	}
	return 'G'
}

type bookOrder []point

func (bo bookOrder) Len() int {
	return len(bo)
}

func (bo bookOrder) Swap(i, j int) {
	bo[i], bo[j] = bo[j], bo[i]
}

func (bo bookOrder) Less(i, j int) bool {
	if bo[i].x == bo[j].x {
		return bo[i].y < bo[j].y
	}
	return bo[i].x < bo[j].x
}

type bookOrderPosition []unit

func (bo bookOrderPosition) Len() int {
	return len(bo)
}

func (bo bookOrderPosition) Swap(i, j int) {
	bo[i], bo[j] = bo[j], bo[i]
}

func (bo bookOrderPosition) Less(i, j int) bool {
	if bo[i].position.x == bo[j].position.x {
		return bo[i].position.y < bo[j].position.y
	}
	return bo[i].position.x < bo[j].position.x
}

func main() {
	// read file
	file := "input"
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	lines := fileLines(file)

	// parse input map
	width, height := 0, 0
	walls := map[point]bool{}
	units := []unit{}
	unitCounter := 1
	for x := range lines {
		for y, c := range lines[x] {
			if c == '#' {
				walls[point{x, y}] = true
			}
			if c == 'G' {
				units = append(units, unit{unitCounter, false, point{x, y}, 200, 3})
				unitCounter++
			} else if c == 'E' {
				units = append(units, unit{unitCounter, true, point{x, y}, 200, 3})
				unitCounter++
			}
			width = imax(y+1, width)
		}
		height = imax(x+1, height)
	}

	// part 1: run the battle and report the results
	fmt.Println("[Part 1]")
	unitsBefore := make([]unit, len(units))
	copy(unitsBefore, units)
	rounds, unitsAfter := runBattle(width, height, walls, unitsBefore, false)
	hp := 0
	for _, u := range unitsAfter {
		if u.health > 0 {
			hp += u.health
		}
	}
	fmt.Println("After combat:")
	printBattle(width, height, walls, units)
	fmt.Println("Combat ends after", rounds, "full rounds")
	fmt.Println("Winner has", hp, "total hit points left")
	fmt.Println("Outcome:", rounds*hp)

	// part 2: start elfs at 3 strength and increase it until the battle concludes with no deaths
	fmt.Println("")
	fmt.Println("[Part 2]")
	strength := 3
	for {
		unitsBefore := make([]unit, len(units))
		copy(unitsBefore, units)
		rounds, unitsAfter := runBattle(width, height, walls, unitsBefore, false)
		elfDeaths := 0
		hp := 0
		for _, u := range unitsAfter {
			if u.health <= 0 && u.elf {
				elfDeaths++
			}
			if u.health > 0 {
				hp += u.health
			}
		}
		fmt.Println("")
		fmt.Println("Elf attack strength is", strength)
		fmt.Println("After combat:")
		printBattle(width, height, walls, unitsAfter)
		fmt.Println("Combat ends after", rounds, "full rounds")
		fmt.Println("Winner has", hp, "total hit points left")
		fmt.Println("Elf deaths:", elfDeaths)
		fmt.Println("Outcome:", rounds*hp)

		if elfDeaths == 0 {
			break
		}

		strength++
		for i := range units {
			if units[i].elf {
				units[i].attack = strength
			}
		}
	}
}

func printBattle(width int, height int, walls map[point]bool, units []unit) {
	battle := make([][]rune, height)
	for x := 0; x < height; x++ {
		battle[x] = make([]rune, width)
		for y := 0; y < width; y++ {
			battle[x][y] = '.'
		}
	}

	for pos := range walls {
		battle[pos.x][pos.y] = '#'
	}

	for _, u := range units {
		if u.health > 0 {
			battle[u.position.x][u.position.y] = u.Symbol()
		}
	}

	for x, s := range battle {
		var sb strings.Builder
		sb.WriteString(string(s))
		sb.WriteString("  ")
		for y := range battle[x] {
			for _, u := range units {
				if u.health > 0 && u.position.x == x && u.position.y == y {
					sb.WriteString(fmt.Sprintf(" %c(%d)", u.Symbol(), u.health))
				}
			}
		}
		fmt.Println(sb.String())
	}
}

func runBattle(width int, height int, walls map[point]bool, units []unit, print bool) (int, []unit) {
	rounds := 0
nextRound:
	for {
		if print {
			fmt.Println("After", rounds, "rounds:")
			printBattle(width, height, walls, units)
			fmt.Println("")
		}

		sort.Sort(bookOrderPosition(units))
		for ui, u := range units {
			if u.health <= 0 {
				continue
			}

			blocked := map[point]bool{}
			for _, v := range units {
				if v.health <= 0 {
					continue
				}
				blocked[v.position] = true
			}
			for p := range walls {
				blocked[p] = true
			}

			// find valid targets on the other team
			targets := []unit{}
			for _, v := range units {
				if v.health <= 0 {
					continue
				}
				if u == v {
					continue
				}
				if u.elf != v.elf {
					targets = append(targets, v)
				}
			}
			if len(targets) == 0 {
				break nextRound
			}

			// check if we need to move in order to attack
			canattack := map[point]unit{}
			for _, v := range targets {
				if u.position.distanceTo(v.position) == 1 {
					canattack[v.position] = v
				}
			}
			if len(canattack) == 0 {
				// find positions in attack range of each target
				inrange := []point{}
				for _, t := range targets {
					options := []point{
						t.position.offset(point{-1, 0}),
						t.position.offset(point{0, -1}),
						t.position.offset(point{0, 1}),
						t.position.offset(point{1, 0}),
					}
					for _, o := range options {
						if _, ok := blocked[o]; ok {
							continue
						}
						inrange = append(inrange, o)
					}
				}

				// find all reachable positions
				reachable := bfsDistances(u.position, blocked)

				// find the closest points that we could move towards
				closest := []point{}
				distance := math.MaxInt32
				for _, p := range inrange {
					if d, ok := reachable[p]; ok {
						if d < distance {
							closest = []point{p}
							distance = d
						} else if d == distance {
							closest = append(closest, p)
						}
					}
				}

				// take a step towards the closest target space
				sort.Sort(bookOrder(closest))
				if len(closest) > 0 {
					goal := closest[0]
					// work paths backwards from the goal
					backwards := bfsDistances(goal, blocked)

					// positions we are next to
					candidates := []point{
						u.position.offset(point{-1, 0}),
						u.position.offset(point{0, -1}),
						u.position.offset(point{0, 1}),
						u.position.offset(point{1, 0}),
					}

					// check all the positions we are next to and find the one that is closest to the goal and not blocked
					candidateDistances := map[point]int{}
					bestDistance := math.MaxInt32
					for _, p := range candidates {
						if _, ok := blocked[p]; ok {
							continue
						}
						if d, ok := backwards[p]; ok {
							candidateDistances[p] = d
							bestDistance = imin(bestDistance, d)
						}
					}

					// if there are more than one candidate moves, we need to pick the book order earliest one
					moveChoice := []point{}
					for p, d := range candidateDistances {
						if d == bestDistance {
							moveChoice = append(moveChoice, p)
						}
					}
					sort.Sort(bookOrder(moveChoice))

					// finally, move
					u.position = moveChoice[0]
					units[ui] = u
				}
			}

			// attack if the unit is now able to
			canattack = map[point]unit{}
			for _, v := range targets {
				if u.position.distanceTo(v.position) == 1 {
					canattack[v.position] = v
				}
			}
			if len(canattack) > 0 {
				// find the units with the lowest health
				minhealth := math.MaxInt32
				for _, v := range canattack {
					minhealth = imin(minhealth, v.health)
				}
				attackoptions := []unit{}
				for _, v := range canattack {
					if v.health == minhealth {
						attackoptions = append(attackoptions, v)
					}
				}
				// attack the first unit with the lowest health
				sort.Sort(bookOrderPosition(attackoptions))
				for i, v := range units {
					if v.id == attackoptions[0].id {
						v.health -= u.attack
						units[i] = v
						break
					}
				}
			}
		}

		// all units have moved, so advance the round marker
		rounds++
	}

	return rounds, units
}

// flood fill distance calculation working outwards from given point
func bfsDistances(pos point, blocked map[point]bool) map[point]int {
	candidates := list.New()
	reachable := map[point]int{}

	reachable[pos] = 0
	candidates.PushBack(pos)

	for candidates.Len() > 0 {
		p := candidates.Front().Value.(point)
		candidates.Remove(candidates.Front())
		options := []point{
			p.offset(point{-1, 0}),
			p.offset(point{0, -1}),
			p.offset(point{0, 1}),
			p.offset(point{1, 0}),
		}
		for _, o := range options {
			// already been here
			if _, ok := reachable[o]; ok {
				continue
			}
			// can't go here
			if _, ok := blocked[o]; ok {
				continue
			}
			// we can reach here with an additional step from where we are now
			reachable[o] = reachable[p] + 1
			candidates.PushBack(o)
		}
	}

	return reachable
}
