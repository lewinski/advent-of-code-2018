package main

import (
	"fmt"
	"math"
)

// alternate solution
// after reading another using cumulative sums, i wanted to code it up because it seemed interesting

func main() {
	// generate the rack per the rules
	serial := 1718
	// grid starts at 1,1 so we can have an row/column of zeroes
	cells := [301][301]int{}
	for x := 1; x < len(cells); x++ {
		for y := 1; y < len(cells[x]); y++ {
			rack := x + 10
			power := rack * y
			power += serial
			power *= rack
			power = (power % 1000) / 100
			power -= 5
			cells[x][y] = power
		}
	}

	// calculate partial sums from 1,1 to x,y
	sums := [301][301]int{}
	for x := 1; x < len(cells); x++ {
		for y := 1; y < len(cells[x]); y++ {
			sums[x][y] = cells[x][y] + sums[x-1][y] + sums[x][y-1] - sums[x-1][y-1]
		}
	}

	// with our sums, we can find the sum of square from positions E to D as (D - C - B + A)
	// A . . B
	// . E . .
	// . . . .
	// C . . D

	// part 1: find the largest 3x3
	// O(n^2)
	maxPower := math.MinInt32
	maxX, maxY := 0, 0
	size := 3
	for x := 1; x < len(cells)-(size-1); x++ {
		for y := 1; y < len(cells[x])-(size-1); y++ {
			totalPower := sums[x+size-1][y+size-1] - sums[x+size-1][y-1] - sums[x-1][y+size-1] + sums[x-1][y-1]
			if totalPower > maxPower {
				maxPower = totalPower
				maxX, maxY = x, y
			}
		}
	}
	fmt.Printf("%d,%d [power = %d]\n", maxX-1, maxY-1, maxPower)

	// part 2: find the largest nxn
	// O(n^3)
	maxPower = math.MinInt32
	maxX, maxY = 0, 0
	maxSize := 0
	for size := 1; size <= 300; size++ {
		for x := 1; x < len(cells)-(size-1); x++ {
			for y := 1; y < len(cells[x])-(size-1); y++ {
				totalPower := sums[x+size-1][y+size-1] - sums[x+size-1][y-1] - sums[x-1][y+size-1] + sums[x-1][y-1]
				if totalPower > maxPower {
					maxPower = totalPower
					maxX, maxY = x, y
					maxSize = size
				}
			}
		}
	}
	fmt.Printf("%d,%d,%d [power = %d]\n", maxX-1, maxY-1, maxSize, maxPower)
}
