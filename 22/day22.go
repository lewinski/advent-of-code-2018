package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
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

// point helpers
type point struct {
	x, y int
}

func (p point) offset(o point) point {
	return point{p.x + o.x, p.y + o.y}
}

func (p point) neighbors() []point {
	return []point{
		p.offset(point{0, -1}),
		p.offset(point{0, 1}),
		p.offset(point{-1, 0}),
		p.offset(point{1, 0}),
	}
}

// gear types
type gear int

const (
	NEITHER  gear = 0
	TORCH    gear = 1
	CLIMBING gear = 2
)

// region types
type regionType rune

const (
	ROCKY  regionType = '.'
	WET    regionType = '='
	NARROW regionType = '|'
)

// calculated geology values
type geology struct {
	geologicIndex int
	erosionLevel  int
	regionType    regionType
	risk          int
}

func printRegion(geo [][]geology, t point) {
	for y := range geo {
		var sb strings.Builder
		for x := range geo[y] {
			if x == 0 && y == 0 {
				sb.WriteRune('M')
			} else if x == t.x && y == t.y {
				sb.WriteRune('T')
			} else {
				sb.WriteRune(rune(geo[y][x].regionType))
			}
		}
		fmt.Println(sb.String())
	}
}

// states for path finding algorithm
type searchPosition struct {
	loc  point
	gear gear
}

// all possible position transitions from current state
// does not consider whether the transition is valid because
// it doesn't know about the regions for each position
func (sp searchPosition) neighbors() []searchPosition {
	n := []searchPosition{}
	for _, p := range sp.loc.neighbors() {
		n = append(n, searchPosition{p, sp.gear})
	}
	if sp.gear != NEITHER {
		n = append(n, searchPosition{sp.loc, NEITHER})
	}
	if sp.gear != TORCH {
		n = append(n, searchPosition{sp.loc, TORCH})
	}
	if sp.gear != CLIMBING {
		n = append(n, searchPosition{sp.loc, CLIMBING})
	}
	return n
}

// given the current gear and the current and next region, is
// this state transition valid?
func validMove(g gear, cur regionType, next regionType) bool {
	if cur == next {
		return true
	}
	if cur == ROCKY {
		if next == WET {
			return g == CLIMBING
		} else if next == NARROW {
			return g == TORCH
		}
	} else if cur == WET {
		if next == ROCKY {
			return g == CLIMBING
		} else if next == NARROW {
			return g == NEITHER
		}
	} else if cur == NARROW {
		if next == ROCKY {
			return g == TORCH
		} else if next == WET {
			return g == NEITHER
		}
	}
	panic("validMove: some case wasn't covered")
}

// heap implementation for graph search
type searchState struct {
	pos  searchPosition
	cost int
}

type searchHeap []searchState

func (h searchHeap) Len() int           { return len(h) }
func (h searchHeap) Less(i, j int) bool { return h[i].cost < h[j].cost }
func (h searchHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *searchHeap) Push(x interface{}) {
	*h = append(*h, x.(searchState))
}

func (h *searchHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func main() {
	// read file
	file := "input"
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	lines := fileLines(file)

	var depth int
	fmt.Sscanf(lines[0], "depth: %d", &depth)
	var tgtx, tgty int
	fmt.Sscanf(lines[1], "target: %d,%d", &tgtx, &tgty)

	regionTypes := []regionType{ROCKY, WET, NARROW}

	// generate the geological map per the puzzle constraints
	margin := 50
	width, height := tgtx+margin, tgty+margin
	geo := make([][]geology, height)
	for y := range geo {
		geo[y] = make([]geology, width)
		for x := range geo[y] {
			if x == 0 && y == 0 {
				geo[y][x].geologicIndex = 0
			} else if x == tgtx && y == tgty {
				geo[y][x].geologicIndex = 0
			} else if x == 0 {
				geo[y][x].geologicIndex = y * 48271
			} else if y == 0 {
				geo[y][x].geologicIndex = x * 16807
			} else {
				geo[y][x].geologicIndex = geo[y-1][x].erosionLevel * geo[y][x-1].erosionLevel
			}
			geo[y][x].erosionLevel = (geo[y][x].geologicIndex + depth) % 20183
			geo[y][x].risk = geo[y][x].erosionLevel % 3
			geo[y][x].regionType = regionTypes[geo[y][x].risk]
		}
	}
	// printRegion(geo, point{tgtx, tgty})

	// calculate the "risk" over the region from the origin to the target
	risk := 0
	for y := range geo {
		for x := range geo[y] {
			if x <= tgtx && y <= tgty {
				risk += geo[y][x].risk
			}
		}
	}
	fmt.Println("part 1 risk =", risk)

	// use dijkstra's to find the minimum amount of time to reach the target
	// - nodes are a combination of x,y position and gear in use
	// - adjacent edges to walk down could be a move in x or y direction
	//   or change in gear, but not both

	// minimum time from origin to each positiion
	distances := map[searchPosition]int{}

	// have we actually been to this position
	visited := map[searchPosition]bool{}

	// min heap for quickly getting next lowest position
	unvisited := &searchHeap{}
	heap.Init(unvisited)
	heap.Push(unvisited, searchState{searchPosition{point{0, 0}, TORCH}, 0})

	for unvisited.Len() > 0 {
		next := heap.Pop(unvisited).(searchState)
		visited[next.pos] = true

		// consider neighboring states (moves or gear changes)
		for _, sp := range next.pos.neighbors() {
			// we've already visited this state
			if _, ok := visited[sp]; ok {
				continue
			}
			// this state is out of bounds
			if sp.loc.x < 0 || sp.loc.x >= width || sp.loc.y < 0 || sp.loc.y >= height {
				continue
			}
			// this state is not a valid move (because the gear doesn't match both regions)
			currentRegion := geo[next.pos.loc.y][next.pos.loc.x].regionType
			nextRegion := geo[sp.loc.y][sp.loc.x].regionType
			if !validMove(next.pos.gear, currentRegion, nextRegion) {
				continue
			}

			// cost for this transition
			cost := 1
			if next.pos.gear != sp.gear {
				cost = 7
			}

			// get the cost for the new position
			oldCost := math.MaxInt32
			if d, ok := distances[sp]; ok {
				oldCost = d
			}
			tmpCost := next.cost + cost

			// if this move is better than the last known one, add it and add
			// the position to the heap of positions to visit
			if tmpCost < oldCost {
				distances[sp] = tmpCost
				heap.Push(unvisited, searchState{sp, tmpCost})
			}
		}
	}

	// target position is at the given coordinates with the torch
	target := searchPosition{point{tgtx, tgty}, TORCH}
	fmt.Println("part 2 minutes =", distances[target])
}
