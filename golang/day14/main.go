package main

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"math"
	"os"
)

func loadData(filename string) ([]byte, error) {
	data, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := data.Close()
		if err != nil {
			panic(err)
		}
	}()
	grid := make([]byte, 0)
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		grid = append(grid, []byte(scanner.Text())...)
	}
	err = scanner.Err()
	if err != nil {
		return nil, err
	}
	return grid, nil
}

func toFlatIndex(row int, col int, height int) int {
	return row*height + col
}

func tiltNorth(grid []byte, dim int) {
	for row := 0; row < dim; row++ {
		for col := 0; col < dim; col++ {
			if grid[toFlatIndex(row, col, dim)] == byte('O') {
				currentRow := row
				for currentRow > 0 && grid[toFlatIndex(currentRow-1, col, dim)] == '.' {
					currentRow--
				}
				if currentRow != row {
					grid[toFlatIndex(currentRow, col, dim)] = byte('O')
					grid[toFlatIndex(row, col, dim)] = byte('.')
				}
			}
		}
	}
}

func calculateLoad(grid []byte, dim int) int {
	load := 0
	for row := 0; row < dim; row++ {
		for col := 0; col < dim; col++ {
			if grid[toFlatIndex(row, col, dim)] == byte('O') {
				load += dim - row
			}
		}
	}
	return load
}

func solutionPart1() {
	grid, err := loadData("input.txt")
	if err != nil {
		panic(err)
	}
	dim := int(math.Sqrt(float64(len(grid))))
	tiltNorth(grid, dim)
	load := calculateLoad(grid, dim)
	fmt.Printf("Load: %d\n", load)
}

func rotateGrid(grid []byte, dim int) {
	for row := 0; row < dim; row++ {
		for col := row + 1; col < dim; col++ {
			grid[toFlatIndex(row, col, dim)], grid[toFlatIndex(col, row, dim)] =
				grid[toFlatIndex(col, row, dim)], grid[toFlatIndex(row, col, dim)]
		}
	}
	for col := 0; col < dim/2; col++ {
		for row := 0; row < dim; row++ {
			grid[toFlatIndex(row, col, dim)], grid[toFlatIndex(row, dim-col-1, dim)] =
				grid[toFlatIndex(row, dim-col-1, dim)], grid[toFlatIndex(row, col, dim)]
		}
	}
}

func cycle(grid []byte, dim int) {
	for i := 0; i < 4; i++ {
		tiltNorth(grid, dim)
		rotateGrid(grid, dim)
	}
}

func findCycle(grid []byte, dim int) (int, int) {
	cache := make(map[[20]byte]int)
	cache[sha1.Sum(grid)] = 0
	index := 0
	for {
		cycle(grid, dim)
		if lastIndex, ok := cache[sha1.Sum(grid)]; ok {
			return index - lastIndex, index + 1
		} else {
			cache[sha1.Sum(grid)] = index
		}
		index += 1
	}
}

func solutionPart2() {
	grid, err := loadData("input.txt")
	if err != nil {
		panic(err)
	}
	dim := int(math.Sqrt(float64(len(grid))))
	remainingCycles := 1000000000
	cycleLength, currentCycle := findCycle(grid, dim)
	remainingCycles -= currentCycle
	remainingCycles -= (remainingCycles / cycleLength) * cycleLength
	for i := 0; i < remainingCycles; i++ {
		cycle(grid, dim)
	}
	load := calculateLoad(grid, dim)
	fmt.Printf("Load: %d\n", load)
}

func main() {
	solutionPart1()
	solutionPart2()
}
