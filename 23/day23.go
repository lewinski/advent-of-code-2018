package main

import (
	"bufio"
	"container/heap"
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

func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func imax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

type point struct {
	x, y, z int
}

// manhatten distance between two points
func mdist(a, b point) int {
	return abs(a.x-b.x) + abs(a.y-b.y) + abs(a.z-b.z)
}

type nanobot struct {
	pos    point
	radius int
}

// does the nanobot's radius include the point?
func (n nanobot) reaches(p point) bool {
	return mdist(n.pos, p) <= n.radius
}

// does the nanobot's radius include any point in the cube?
func (n nanobot) reachesCube(c cube) bool {
	return c.distanceToPoint(n.pos) <= n.radius
}

// cube is defined by the minimum corner and the length of the side
type cube struct {
	corner point
	side   int
}

// helper for cube.distanceToPoint
func distanceToRange(x, low, high int) int {
	if x < low {
		return low - x
	} else if x > high {
		return x - high
	}
	return 0
}

// distance between the point and the edge of the cube, or zero if the point
// is in the cube. dimensions are inclusive on the min side and exclusive on the
// max side. i.e. [x, x+side)
func (c cube) distanceToPoint(p point) int {
	d := 0
	d += distanceToRange(p.x, c.corner.x, c.corner.x+c.side-1)
	d += distanceToRange(p.y, c.corner.y, c.corner.y+c.side-1)
	d += distanceToRange(p.z, c.corner.z, c.corner.z+c.side-1)
	return d
}

// divide the cube into 8 subcubes along the midpoint
func (c cube) subdivisions() []cube {
	subs := make([]cube, 8)
	side := c.side / 2
	subs[0] = cube{point{c.corner.x, c.corner.y, c.corner.z}, side}
	subs[1] = cube{point{c.corner.x, c.corner.y, c.corner.z + side}, side}
	subs[2] = cube{point{c.corner.x, c.corner.y + side, c.corner.z}, side}
	subs[3] = cube{point{c.corner.x, c.corner.y + side, c.corner.z + side}, side}
	subs[4] = cube{point{c.corner.x + side, c.corner.y, c.corner.z}, side}
	subs[5] = cube{point{c.corner.x + side, c.corner.y, c.corner.z + side}, side}
	subs[6] = cube{point{c.corner.x + side, c.corner.y + side, c.corner.z}, side}
	subs[7] = cube{point{c.corner.x + side, c.corner.y + side, c.corner.z + side}, side}
	return subs
}

// returns number of bots that have ranges that overlap with the cube
func botsReachingCube(bots []nanobot, c cube) int {
	cnt := 0
	for _, bot := range bots {
		if bot.reachesCube(c) {
			cnt++
		}
	}
	return cnt
}

// heap implementation for search
type searchState struct {
	bounds cube
	count  int
	dist   int
}

type searchHeap []searchState

func (h searchHeap) Len() int      { return len(h) }
func (h searchHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h searchHeap) Less(i, j int) bool {
	// higher counts come first
	if h[i].count != h[j].count {
		return h[i].count > h[j].count
	}
	// if counts are the same, closer ones to the origin come first
	if h[i].dist != h[j].dist {
		return h[i].dist < h[j].dist
	}
	// if same distance to the origin, then smaller items come first
	return h[i].bounds.side < h[j].bounds.side
}

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

	// parse the data file and find the bot with the largest radius
	bots := []nanobot{}
	maxBot := nanobot{radius: -1}
	for _, l := range lines {
		n := nanobot{}
		fmt.Sscanf(l, "pos=<%d,%d,%d>, r=%d", &n.pos.x, &n.pos.y, &n.pos.z, &n.radius)
		bots = append(bots, n)
		if maxBot.radius < n.radius {
			maxBot = n
		}
	}

	// count the bots in the largest bot's radius
	botsInRange := 0
	for _, bot := range bots {
		distance := mdist(bot.pos, maxBot.pos)
		if distance <= maxBot.radius {
			botsInRange++
		}
	}
	fmt.Println("part 1 - bots in range of largest radius:", botsInRange)

	// part 2 basic ideas:
	// if you have boxes B and C and C is completely inside B, then C cannot be
	// reachable by any more bots that B is. so we want to start with a box that
	// covers the range of every bot and then subdivide that box and search only
	// in the subbox that intersects with the most bots. then keep subdividing
	// until we reach a box that contains only one point and that is the answer.

	// in order to do this, we need to be able to quickly calculate if a bot's
	// radius includes any point in the cube. we can do this by calculating the
	// distance between the bot and the closest edge/corner of the cube, and then
	// comparing this distance with the bot's radius. if the radius is greater
	// than the distance, the the bot's range intersects with with the cube.
	// this is nanobot.reachesCube() and cube.distanceToPoint() above

	// figure out a bounding cube symmetric around origin with sides a power of
	// two that covers the entire range of every nanobot
	maxExtent := 0
	for _, bot := range bots {
		maxExtent = imax(maxExtent, abs(bot.pos.x)+bot.radius)
		maxExtent = imax(maxExtent, abs(bot.pos.y)+bot.radius)
		maxExtent = imax(maxExtent, abs(bot.pos.z)+bot.radius)
	}
	size := 1
	for size <= maxExtent {
		size *= 2
	}
	cover := cube{point{-size, -size, -size}, size * 2}

	// heap to keep track of the best cubes we've found so far
	h := &searchHeap{}
	heap.Init(h)

	// initialize the search heap with our covering cube
	origin := point{0, 0, 0}
	init := searchState{cover, botsReachingCube(bots, cover), cover.distanceToPoint(origin)}
	heap.Push(h, init)

	// while we have cubes in our heap
	for h.Len() > 0 {
		// grab the best item in the heap
		currentState := heap.Pop(h).(searchState)

		// if the item covers only one integer location, it is the answer
		if currentState.bounds.side == 1 {
			fmt.Printf("part 2 - distance: %d; count: %d\n", currentState.dist, currentState.count)
			break
		}

		// otherwise, divide it into 8 smaller boxes and push each of them onto
		// the search heap so we can continue the search
		for _, c := range currentState.bounds.subdivisions() {
			st := searchState{bounds: c}
			st.count = botsReachingCube(bots, c)
			st.dist = c.distanceToPoint(origin)
			heap.Push(h, st)
		}
	}
}
