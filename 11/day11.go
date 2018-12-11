package main

import (
	"fmt"
	"math"
)

func main() {
	// generate the rack per the rules
	serial := 1718
	cells := [300][300]int{}
	for x := range cells {
		for y := range cells[x] {
			rack := x + 10
			power := rack * y
			power += serial
			power *= rack
			power = (power % 1000) / 100
			power -= 5
			cells[x][y] = power
		}
	}

	// part 1: find the largest 3x3
	// O(n^2)
	maxPower := math.MinInt32
	maxX, maxY := 0, 0
	for x := 0; x < len(cells)-2; x++ {
		for y := 0; y < len(cells[x])-2; y++ {
			totalPower := cells[x][y] + cells[x][y+1] + cells[x][y+2] +
				cells[x+1][y] + cells[x+1][y+1] + cells[x+1][y+2] +
				cells[x+2][y] + cells[x+2][y+1] + cells[x+2][y+2]
			if totalPower > maxPower {
				maxPower = totalPower
				maxX, maxY = x, y
			}
		}
	}
	fmt.Printf("%d,%d [power = %d]\n", maxX-1, maxY-1, maxPower)

	// part2: find the largest nxn
	// O(n^5)
	maxPower = math.MinInt32
	maxX, maxY = 0, 0
	maxSize := 0
	for size := 1; size <= 300; size++ {
		for x := 0; x < len(cells)-size-1; x++ {
			for y := 0; y < len(cells[x])-size-1; y++ {
				totalPower := 0
				for i := 0; i < size; i++ {
					for j := 0; j < size; j++ {
						totalPower += cells[x+i][y+j]
					}
				}
				if totalPower > maxPower {
					maxPower = totalPower
					maxX, maxY, maxSize = x, y, size
				}
			}
		}
	}
	fmt.Printf("%d,%d,%d [power = %d]\n", maxX-1, maxY-1, maxSize, maxPower)
}
