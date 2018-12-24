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

type team int

const (
	IMMUNE    team = 1
	INFECTION team = 2
)

func (t team) String() string {
	if t == IMMUNE {
		return "Immune System"
	}
	return "Infection"
}

type unit struct {
	id             int
	team           team
	initiative     int
	count          int
	hp             int
	attack         int
	attackType     string
	damageModifier map[string]int
	target         *unit
}

func (u unit) effectivePower() int {
	return u.attack * u.count
}

func (u unit) damageFor(o unit) int {
	dmg := u.effectivePower()
	if mod, ok := o.damageModifier[u.attackType]; ok {
		dmg *= mod
	}
	return dmg
}

func (u unit) applyDamage() {
	if u.target != nil {
		dmg := u.damageFor(*u.target)
		if dmg == 0 {
			return
		}
		killed := dmg / u.target.hp
		// fmt.Printf("%s group %d attacks defending group %d, killing %d units\n", u.team.String(), u.id, u.target.id, killed)
		if killed > u.target.count {
			u.target.count = 0
		} else {
			u.target.count -= killed
		}
	}
}

type byInitiative []*unit

func (b byInitiative) Len() int {
	return len(b)
}

func (b byInitiative) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byInitiative) Less(i, j int) bool {
	return b[j].initiative < b[i].initiative
}

type byEffectivePower []*unit

func (b byEffectivePower) Len() int {
	return len(b)
}

func (b byEffectivePower) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byEffectivePower) Less(i, j int) bool {
	if b[j].effectivePower() == b[i].effectivePower() {
		return b[i].initiative < b[j].initiative
	}
	return b[j].effectivePower() < b[i].effectivePower()
}

func fight(file string, boost int) []*unit {

	lines := fileLines(file)

	unitre1 := regexp.MustCompile("^(\\d+) units each with (\\d+) hit points")
	unitre2 := regexp.MustCompile("with an attack that does (\\d+) ([a-z]+) damage at initiative (\\d+)$")
	unitre3 := regexp.MustCompile("weak to ((?:[a-z]+(?:, )?)+)")
	unitre4 := regexp.MustCompile("immune to ((?:[a-z]+(?:, )?)+)")

	var currentTeam team
	var currentID int
	var currentBoost int
	units := []*unit{}
	for _, l := range lines {
		if l == "Immune System:" {
			currentTeam = IMMUNE
			currentID = 1
			currentBoost = boost
		} else if l == "Infection:" {
			currentTeam = INFECTION
			currentID = 1
			currentBoost = 0
		} else if unitre1.MatchString(l) {
			u := unit{team: currentTeam, id: currentID}
			u.damageModifier = make(map[string]int)
			currentID++

			matches := unitre1.FindStringSubmatch(l)
			u.count = atoi(matches[1])
			u.hp = atoi(matches[2])

			matches = unitre2.FindStringSubmatch(l)
			u.attack = atoi(matches[1]) + currentBoost
			u.attackType = matches[2]
			u.initiative = atoi(matches[3])

			matches = unitre3.FindStringSubmatch(l)
			if len(matches) > 1 {
				weaknesses := strings.Split(matches[1], ", ")
				for _, w := range weaknesses {
					u.damageModifier[w] = 2
				}
			}

			matches = unitre4.FindStringSubmatch(l)
			if len(matches) > 1 {
				immunities := strings.Split(matches[1], ", ")
				for _, i := range immunities {
					u.damageModifier[i] = 0
				}
			}

			units = append(units, &u)
		}
	}

	// if fight isn't resolved after a bunch of iterations, it might be a stalemate
	for c := 0; c < 10000; c++ {
		// attack selection
		sort.Sort(byEffectivePower(units))

		// target positions that have already been selected this fight iteration
		selected := map[int]bool{}

		for i := range units {
			if units[i].count == 0 {
				continue
			}

			maxDamage := 0
			targetPos := -1

			for j := range units {
				// only attack other team
				if units[i].team == units[j].team {
					continue
				}
				// don't attack dead units
				if units[j].count == 0 {
					continue
				}
				// don't attack already selected units
				if _, ok := selected[j]; ok {
					continue
				}

				// only attack units we can damage
				dmg := units[i].damageFor(*units[j])
				if dmg == 0 {
					continue
				}

				// fmt.Printf("%v group %d would deal defending group %d %d damage\n", units[i].team, units[i].id, units[j].id, dmg)

				if dmg == maxDamage {
					// if this is equal to the best we could do, follow the tie breaking rules
					if units[j].effectivePower() == units[targetPos].effectivePower() {
						if units[j].initiative > units[targetPos].initiative {
							targetPos = j
						}
					} else if units[j].effectivePower() > units[targetPos].effectivePower() {
						targetPos = j
					}
				} else if dmg > maxDamage {
					// if it is better than previous options, select it
					maxDamage = dmg
					targetPos = j
				}
			}

			// if we picked a target, select it
			if targetPos >= 0 {
				units[i].target = units[targetPos]
				selected[targetPos] = true
			} else {
				// oterwise, deselect an old target
				units[i].target = nil
			}
		}

		// attack phase
		sort.Sort(byInitiative(units))

		// fmt.Println("")
		for i := range units {
			// dead units can't attack
			if units[i].count == 0 {
				continue
			}
			units[i].applyDamage()
		}
		// fmt.Println("")

		// exit early if the fight has a victor
		alive := map[team]int{}
		for i := range units {
			if units[i].count > 0 {
				alive[units[i].team]++
			}
		}
		if len(alive) == 1 {
			break
		}
	}

	return units
}

func main() {
	// read file
	file := "input"
	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	units := fight(file, 0)
	alive := map[team]int{}
	for i := range units {
		if units[i].count > 0 {
			alive[units[i].team] += units[i].count
			// fmt.Printf("%v group %d contains %d units\n", units[i].team, units[i].id, units[i].count)
		}
	}
	for t, c := range alive {
		fmt.Printf("part 1: %v wins, alive = %d\n", t, c)
	}

BOOST:
	// turns out the boost needed for my actual input is pretty low, so we can just run the fight until we find it.
	for boost := 1; boost < 2000; boost++ {
		units := fight(file, boost)

		// total up each team
		alive := map[team]int{}
		for i := range units {
			if units[i].count > 0 {
				alive[units[i].team] += units[i].count
			}
		}
		if len(alive) == 1 {
			for t, c := range alive {
				if t == IMMUNE {
					fmt.Printf("part 2 (boost = %d): %v wins, alive = %d\n", boost, t, c)
					break BOOST
				}
			}
		} else {
			// fmt.Printf("part 2 (boost = %d): stalemate %v\n", boost, alive)
		}
	}
}
