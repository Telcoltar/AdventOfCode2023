package main

import (
	"bufio"
	"fmt"
	"os"
)

type Direction = int

const (
	UP = iota
	RIGHT
	DOWN
	LEFT
)

func oppositeDirection(direction Direction) Direction {
	switch direction {
	case UP:
		return DOWN
	case DOWN:
		return UP
	case LEFT:
		return RIGHT
	case RIGHT:
		return LEFT
	}
	panic("Unknown Direction")
}

func loadData(filename string) ([][]rune, error) {
	data, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := data.Close(); err != nil {
			panic(err)
		}
	}()
	scanner := bufio.NewScanner(data)
	grid := make([][]rune, 0)
	for scanner.Scan() {
		grid = append(grid, []rune(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return grid, nil
}

func padGrid(grid [][]rune) [][]rune {
	dim := len(grid) + 2
	newGrid := make([][]rune, dim)
	newGrid[0] = make([]rune, dim)
	for i := 0; i < dim; i++ {
		newGrid[0][i] = '#'
	}
	for i := 1; i < dim-1; i++ {
		newGrid[i] = append(newGrid[i], '#')
		newGrid[i] = append(newGrid[i], grid[i-1]...)
		newGrid[i] = append(newGrid[i], '#')
	}
	newGrid[dim-1] = make([]rune, dim)
	for i := 0; i < dim; i++ {
		newGrid[dim-1][i] = '#'
	}
	return newGrid
}

type Point struct {
	x int
	y int
}

func (p *Point) move(direction Direction) {
	switch direction {
	case UP:
		p.y -= 1
	case DOWN:
		p.y += 1
	case LEFT:
		p.x -= 1
	case RIGHT:
		p.x += 1
	}
}

type Status struct {
	Point
	direction Direction
}

func (s *Status) move() {
	s.Point.move(s.direction)
}

func (s *Status) moveBack() {
	s.Point.move(oppositeDirection(s.direction))
}

func processCorner(status *Status) {
	switch status.direction {
	case UP:
		status.direction = RIGHT
	case DOWN:
		status.direction = LEFT
	case LEFT:
		status.direction = DOWN
	case RIGHT:
		status.direction = UP
	}
}

func processBackwardCorner(status *Status) {
	switch status.direction {
	case UP:
		status.direction = LEFT
	case DOWN:
		status.direction = RIGHT
	case LEFT:
		status.direction = UP
	case RIGHT:
		status.direction = DOWN
	}
}

func processUpDown(status Status) []Status {
	if status.direction == LEFT || status.direction == RIGHT {
		statusOne := status
		statusOne.direction = UP
		statusTwo := status
		statusTwo.direction = DOWN
		return []Status{statusOne, statusTwo}
	}
	return []Status{status}
}

func processLeftRight(status Status) []Status {
	if status.direction == UP || status.direction == DOWN {
		statusOne := status
		statusOne.direction = LEFT
		statusTwo := status
		statusTwo.direction = RIGHT
		return []Status{statusOne, statusTwo}
	}
	return []Status{status}
}

func calculateEnergizedTiles(grid [][]rune, start Status) int {
	visited := make(map[Status]struct{})
	count := make(map[Point]struct{})
	queue := make([]Status, 0, 5)
	queue = append(queue, start)
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if _, ok := visited[current]; ok {
			continue
		}
		visited[current] = struct{}{}
		count[current.Point] = struct{}{}
		current.move()
		tile := grid[current.y][current.x]
		switch tile {
		case '/':
			processCorner(&current)
			queue = append(queue, current)
		case '\\':
			processBackwardCorner(&current)
			queue = append(queue, current)
		case '|':
			queue = append(queue, processUpDown(current)...)
		case '-':
			queue = append(queue, processLeftRight(current)...)
		case '.':
			if grid[current.y][current.x] == '.' {
				count[current.Point] = struct{}{}
				current.move()
			}
			current.moveBack()
			queue = append(queue, current)
		}
	}
	return len(count) - 1
}

func solutionPart1() {
	grid, err := loadData("input.txt")
	grid = padGrid(grid)
	if err != nil {
		panic(err)
	}
	count := calculateEnergizedTiles(grid, Status{Point{x: 0, y: 1}, RIGHT})
	fmt.Println("Count:", count)
}

func solutionPart2() {
	grid, err := loadData("input.txt")
	grid = padGrid(grid)
	if err != nil {
		panic(err)
	}
	energizedTiles := make([]int, 0)
	for i := 1; i < len(grid)-1; i++ {
		energizedTiles = append(energizedTiles, calculateEnergizedTiles(grid, Status{Point{0, i}, RIGHT}))
		energizedTiles = append(energizedTiles, calculateEnergizedTiles(grid, Status{Point{len(grid) - 1, i}, LEFT}))
		energizedTiles = append(energizedTiles, calculateEnergizedTiles(grid, Status{Point{i, 0}, DOWN}))
		energizedTiles = append(energizedTiles, calculateEnergizedTiles(grid, Status{Point{i, len(grid) - 1}, UP}))
	}
	maxTiles := 0
	for _, el := range energizedTiles {
		maxTiles = max(maxTiles, el)
	}
	fmt.Println("Count:", maxTiles)
}

func main() {
	solutionPart1()
	solutionPart2()
}
