package main

import (
	"bytes"
	"fmt"
	"os"
	"slices"
	"strconv"
	"time"
)

func readData(filepath string) ([][]int, error) {
	fileContent, readErr := os.ReadFile(filepath)
	if readErr != nil {
		return nil, readErr
	}
	grid := make([][]int, 0)
	for _, line := range bytes.Split(fileContent, []byte("\n")) {
		if len(line) == 0 {
			continue
		}
		row := make([]int, 0)
		for _, char := range line {
			parsedNum, parseErr := strconv.Atoi(string(char))
			if parseErr != nil {
				return nil, parseErr
			}
			row = append(row, parsedNum)
		}
		grid = append(grid, row)
	}
	return grid, nil
}

var UP Vec = Vec{X: 0, Y: -1}
var RIGHT Vec = Vec{X: 1, Y: 0}
var DOWN Vec = Vec{X: 0, Y: 1}
var LEFT Vec = Vec{X: -1, Y: 0}

type Vec struct {
	X, Y int
}

func getPerpendicular(dir Vec) []Vec {
	if dir == UP || dir == DOWN {
		return []Vec{LEFT, RIGHT}
	}
	return []Vec{UP, DOWN}
}

type Status struct {
	Direction Vec
	Position  Vec
	Count     int
}

type grid [][]int

func (g grid) get(pos Vec) int {
	return g[pos.Y][pos.X]
}

func getPoint(pos Vec, dir Vec, height, width int) (Vec, bool) {
	if pos.X+dir.X < 0 || pos.X+dir.X >= width || pos.Y+dir.Y < 0 || pos.Y+dir.Y >= height {
		return Vec{}, false
	}
	return Vec{X: pos.X + dir.X, Y: pos.Y + dir.Y}, true
}

func solutionPart1(weights grid) int {
	end := []int{len(weights[0]) - 1, len(weights) - 1}
	weightGrid := map[Status]int{
		{Direction: RIGHT, Position: Vec{X: 1, Y: 0}, Count: 1}: weights.get(Vec{X: 1, Y: 0}),
		{Direction: DOWN, Position: Vec{X: 0, Y: 1}, Count: 1}:  weights.get(Vec{X: 0, Y: 1}),
	}
	queue := []Vec{{X: 1, Y: 0}, {X: 0, Y: 1}}
	for len(queue) > 0 {
		pos := queue[0]
		queue = queue[1:]
		for _, dir := range []Vec{UP, RIGHT, DOWN, LEFT} {
			if point, valid := getPoint(pos, dir, len(weights), len(weights[0])); valid {
				prevWeights := make([]int, 0)
				nextPoint := false
				for i := 0; i <= 2; i++ {
					if prevWeight, ok := weightGrid[Status{Direction: dir, Position: pos, Count: i}]; ok {
						prevWeights = append(prevWeights, prevWeight)
						if weight, ok := weightGrid[Status{Direction: dir, Position: point, Count: i + 1}]; ok && weight > prevWeight+weights.get(point) {
							weightGrid[Status{Direction: dir, Position: point, Count: i + 1}] = prevWeight + weights.get(point)
							nextPoint = true
						} else if !ok {
							weightGrid[Status{Direction: dir, Position: point, Count: i + 1}] = prevWeight + weights.get(point)
							nextPoint = true
						}
					}
				}
				if len(prevWeights) > 0 {
					minPrevWeights := slices.Min(prevWeights)
					for _, dir := range getPerpendicular(dir) {
						if weight, ok := weightGrid[Status{Direction: dir, Position: point, Count: 0}]; ok && weight > minPrevWeights+weights.get(point) {
							weightGrid[Status{Direction: dir, Position: point, Count: 0}] = minPrevWeights + weights.get(point)
							nextPoint = true
						} else if !ok {
							weightGrid[Status{Direction: dir, Position: point, Count: 0}] = minPrevWeights + weights.get(point)
							nextPoint = true
						}
					}
				}
				if nextPoint {
					queue = append(queue, point)
				}
			}
		}
	}
	endWeights := make([]int, 0)
	for i := 0; i <= 3; i++ {
		for _, dir := range []Vec{UP, RIGHT, DOWN, LEFT} {
			if weight, ok := weightGrid[Status{Direction: dir, Position: Vec{X: end[0], Y: end[1]}, Count: i}]; ok {
				endWeights = append(endWeights, weight)
			}
		}
	}

	return slices.Min(endWeights)
}

func solutionPart2(weights grid) int {
	end := []int{len(weights[0]) - 1, len(weights) - 1}
	weightGrid := map[Status]int{
		{Direction: RIGHT, Position: Vec{X: 1, Y: 0}, Count: 1}: weights.get(Vec{X: 1, Y: 0}),
		{Direction: DOWN, Position: Vec{X: 0, Y: 1}, Count: 1}:  weights.get(Vec{X: 0, Y: 1}),
	}
	queue := []Vec{{X: 1, Y: 0}, {X: 0, Y: 1}}
	for len(queue) > 0 {
		pos := queue[0]
		queue = queue[1:]
		for _, dir := range []Vec{UP, RIGHT, DOWN, LEFT} {
			if point, valid := getPoint(pos, dir, len(weights), len(weights[0])); valid {

				nextPoint := false
				for i := 1; i < 10; i++ {
					if prevWeight, ok := weightGrid[Status{Direction: dir, Position: pos, Count: i}]; ok {
						if weight, ok := weightGrid[Status{Direction: dir, Position: point, Count: i + 1}]; ok && weight > prevWeight+weights.get(point) {
							weightGrid[Status{Direction: dir, Position: point, Count: i + 1}] = prevWeight + weights.get(point)
							nextPoint = true
						} else if !ok {
							weightGrid[Status{Direction: dir, Position: point, Count: i + 1}] = prevWeight + weights.get(point)
							nextPoint = true
						}
					}
				}

				prevWeights := make([]int, 0)

				for _, dir := range getPerpendicular(dir) {
					for i := 4; i < 11; i++ {
						if prevWeight, ok := weightGrid[Status{Direction: dir, Position: pos, Count: i}]; ok {
							prevWeights = append(prevWeights, prevWeight)
						}
					}
				}
				if len(prevWeights) > 0 {
					minPrevWeights := slices.Min(prevWeights)

					if weight, ok := weightGrid[Status{Direction: dir, Position: point, Count: 1}]; ok && weight > minPrevWeights+weights.get(point) {
						weightGrid[Status{Direction: dir, Position: point, Count: 1}] = minPrevWeights + weights.get(point)
						nextPoint = true
					} else if !ok {
						weightGrid[Status{Direction: dir, Position: point, Count: 1}] = minPrevWeights + weights.get(point)
						nextPoint = true
					}

				}

				if nextPoint {
					queue = append(queue, point)
				}
			}
		}
	}
	endWeights := make([]int, 0)
	for i := 4; i < 11; i++ {
		for _, dir := range []Vec{UP, RIGHT, DOWN, LEFT} {
			if weight, ok := weightGrid[Status{Direction: dir, Position: Vec{X: end[0], Y: end[1]}, Count: i}]; ok {
				endWeights = append(endWeights, weight)
			}
		}
	}

	return slices.Min(endWeights)
}

func main() {
	grid, readErr := readData("input.txt")
	if readErr != nil {
		panic(readErr)
	}
	startTime := time.Now()
	println(solutionPart1(grid))
	fmt.Printf("Part1: %s\n", time.Since(startTime))
	startTime = time.Now()
	println(solutionPart2(grid))
	fmt.Printf("Part2: %s\n", time.Since(startTime))
}
