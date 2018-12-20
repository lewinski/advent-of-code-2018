package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
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

type point struct {
	x, y int
}

type pointlist []point

func (pl pointlist) dup() pointlist {
	d := make(pointlist, len(pl))
	copy(d, pl)
	return d
}

func (pl pointlist) unique() pointlist {
	keys := make(map[point]bool)
	list := pointlist{}
	for _, entry := range pl {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func (p point) offset(o point) point {
	return point{p.x + o.x, p.y + o.y}
}

// doors are at odd coordinates, rooms are at even, so walking from
// room to room we need to go 2 in some direction
func (p point) traverseDoor(door point) point {
	return p.offset(point{2 * (door.x - p.x), 2 * (door.y - p.y)})
}

type doors map[point]bool

func (d doors) addAndTraverseDoor(current point, direction rune) point {
	var move point
	switch direction {
	case 'N':
		move = point{0, -1}
	case 'E':
		move = point{1, 0}
	case 'S':
		move = point{0, 1}
	case 'W':
		move = point{-1, 0}
	default:
		log.Panicf("invalid direction: %c", direction)
	}

	door := current.offset(move)
	d[door] = true

	return current.traverseDoor(door)
}

func (d doors) adjacentDoors(current point) []point {
	adj := []point{}
	options := []point{
		current.offset(point{0, -1}),
		current.offset(point{1, 0}),
		current.offset(point{0, 1}),
		current.offset(point{-1, 0}),
	}
	for _, o := range options {
		if _, ok := d[o]; ok {
			adj = append(adj, o)
		}
	}
	return adj
}

func main() {
	// read file
	file := "input"
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	lines := fileLines(file)

	doors := make(doors, 0)
	pos := pointlist{point{0, 0}}
	heads := make([]pointlist, 0)
	tails := pointlist{}

	for _, r := range lines[0] {
		pos = pos.unique()
		switch r {
		case '^', '$':
			// nothing
		case 'N', 'E', 'S', 'W':
			// from each possible endpoint, move through the door
			for i := range pos {
				pos[i] = doors.addAndTraverseDoor(pos[i], r)
			}
		case '(':
			// save the current position so we can refer back to it
			heads = append(heads, pos.dup())
			// initialize a new list of alternates
			tails = pointlist{}
		case '|':
			// copy the current position(s) to the alternates list
			tails = append(tails, pos...)
			// reset positions back to the initial state of this alternation
			pos = heads[len(heads)-1].dup()
		case ')':
			// add all the alternates to the position list
			pos = append(pos, tails...)
			// forget the last set of starting positions
			heads = heads[:len(heads)-1]
		default:
			panic("ack")
		}
	}

	distances := bfs(point{0, 0}, doors)

	maxDistance := 0
	thousandPlus := 0
	for _, d := range distances {
		if d > maxDistance {
			maxDistance = d
		}
		if d >= 1000 {
			thousandPlus++
		}
	}

	// part 1
	fmt.Println("max distance:", maxDistance)
	if len(lines) > 1 {
		fmt.Println("test data max distance:", lines[1])
	}

	// part 2
	fmt.Println("rooms 1000+ doors away:", thousandPlus)
}

func bfs(pos point, d doors) map[point]int {
	reachable := map[point]int{}
	reachable[pos] = 0

	candidates := list.New()
	candidates.PushBack(pos)

	for candidates.Len() > 0 {
		p := candidates.Front().Value.(point)
		candidates.Remove(candidates.Front())

		for _, o := range d.adjacentDoors(p) {
			q := p.traverseDoor(o)
			// already been here
			if _, ok := reachable[q]; ok {
				continue
			}
			// we can reach here with an additional step from where we are now
			reachable[q] = reachable[p] + 1
			candidates.PushBack(q)
		}
	}

	return reachable
}
