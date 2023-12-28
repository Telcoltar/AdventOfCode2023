package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

func loadData(filename string) ([][]string, error) {
	data, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(data *os.File) {
		err := data.Close()
		if err != nil {
			panic(err)
		}
	}(data)
	var grids [][]string
	scanner := bufio.NewScanner(data)
	var grid []string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			grids = append(grids, grid)
			grid = make([]string, 0)
			continue
		}
		grid = append(grid, line)
	}
	grids = append(grids, grid)
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return grids, nil
}

func isStrictEqual(a, b string) bool {
	return a == b
}

func isSimilar(a, b string) bool {
	diff := 0
	aRune := []rune(a)
	bRune := []rune(b)
	for i := 0; i < len(aRune); i++ {
		if aRune[i] != bRune[i] {
			diff++
		}
	}
	return diff == 0 || diff == 1
}

func getStartingPoint(grid []string, equalFn func(string, string) bool) []int {
	points := make([]int, 0)
	for i := 0; i < len(grid)-1; i++ {
		if equalFn(grid[i], grid[i+1]) {
			points = append(points, i)
		}
	}
	return points
}

func checkIfStartingPointWorks(grid []string, startingPoint int, equalFn func(string, string) bool) bool {
	lower := startingPoint
	upper := startingPoint + 1
	for lower >= 0 && upper < len(grid) && equalFn(grid[lower], grid[upper]) {
		upper++
		lower--
	}
	return lower == -1 || upper == len(grid)
}

func scanForHorizontalReflectionLine(grid []string, equalFn func(string, string) bool) []int {
	startingPoints := getStartingPoint(grid, equalFn)
	validPoints := make([]int, 0)
	for _, startingPoint := range startingPoints {
		if checkIfStartingPointWorks(grid, startingPoint, equalFn) {
			validPoints = append(validPoints, startingPoint)
		}
	}
	return validPoints
}

func transposeGrid(grid []string) []string {
	transposedGrid := make([]string, len(grid[0]))
	for i := 0; i < len(grid[0]); i++ {
		for j := 0; j < len(grid); j++ {
			transposedGrid[i] += string(grid[j][i])
		}
	}
	return transposedGrid
}

func solutionPart1() {
	grids, err := loadData("input.txt")
	if err != nil {
		panic(err)
	}
	horizontalSum := 0
	verticalSum := 0
	for _, grid := range grids {
		horizontalStartingPoints := scanForHorizontalReflectionLine(grid, isStrictEqual)
		if len(horizontalStartingPoints) == 1 {
			horizontalSum += horizontalStartingPoints[0] + 1
		} else {
			transposedGrid := transposeGrid(grid)
			verticalStartingPoints := scanForHorizontalReflectionLine(transposedGrid, isStrictEqual)
			if len(verticalStartingPoints) == 1 {
				verticalSum += verticalStartingPoints[0] + 1
			} else {
				panic("Does not found exactly one solution")
			}
		}
	}
	totalSum := verticalSum + horizontalSum*100
	fmt.Println(totalSum)
}

func getCorrectedStartingPoint(grid []string) []int {
	startingPoints := scanForHorizontalReflectionLine(grid, isStrictEqual)
	diff := make([]int, 0)
	for _, startingPoint := range scanForHorizontalReflectionLine(grid, isSimilar) {
		if !slices.Contains(startingPoints, startingPoint) {
			diff = append(diff, startingPoint)
		}
	}
	return diff
}

func solutionPart2() {
	grids, err := loadData("input.txt")
	if err != nil {
		panic(err)
	}
	horizontalSum := 0
	verticalSum := 0
	for _, grid := range grids {
		horizontalStartingPoints := getCorrectedStartingPoint(grid)
		if len(horizontalStartingPoints) == 1 {
			horizontalSum += horizontalStartingPoints[0] + 1
		} else {
			transposedGrid := transposeGrid(grid)
			verticalStartingPoints := getCorrectedStartingPoint(transposedGrid)
			if len(verticalStartingPoints) == 1 {
				verticalSum += verticalStartingPoints[0] + 1
			} else {
				panic("Did not find unique starting point")
			}
		}
	}
	totalSum := verticalSum + horizontalSum*100
	fmt.Println(totalSum)
}

func main() {
	solutionPart1()
	solutionPart2()
}
