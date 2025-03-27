package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"time"
)

func positiveModulo(a, b int) int {
	return (a%b + b) % b
}

type Point struct {
	X, Y int
}

type Grid struct {
	Cells         [][]rune
	Width, Height int
	Offset        int
}

func (g *Grid) Set(p Point, value rune) {
	g.Cells[p.Y][p.X] = value
}

func (g *Grid) Get(p Point) rune {
	return g.Cells[positiveModulo(p.Y+g.Offset, g.Height)][positiveModulo(p.X+g.Offset, g.Width)]
}

func (g *Grid) isInside(p Point) bool {
	return p.X >= -g.Offset && p.X <= g.Offset && p.Y >= -g.Offset && p.Y <= g.Offset
}

func (g *Grid) getValidNeighbours(p Point) []Point {
	neighbours := []Point{}
	for _, dir := range []Point{{0, -1}, {0, 1}, {-1, 0}, {1, 0}} {
		neighbour := Point{p.X + dir.X, p.Y + dir.Y}
		if g.isInside(neighbour) && g.Get(neighbour) == '.' {
			neighbours = append(neighbours, neighbour)
		}
	}
	return neighbours
}

func (g *Grid) getValidNeighboursPart2(p Point) []Point {
	neighbours := []Point{}
	for _, dir := range []Point{{0, -1}, {0, 1}, {-1, 0}, {1, 0}} {
		neighbour := Point{p.X + dir.X, p.Y + dir.Y}
		if g.Get(neighbour) != '#' {
			neighbours = append(neighbours, neighbour)
		}
	}
	return neighbours
}

func readData(filepath string) (Grid, Point) {
	fileContent, readErr := os.ReadFile(filepath)
	if readErr != nil {
		log.Fatal(readErr)
	}
	grid := Grid{}
	grid.Cells = [][]rune{}
	var start Point
	for y, line := range strings.Split(string(fileContent), "\n") {
		row := []rune(line)
		if idx := slices.Index(row, 'S'); idx != -1 {
			start = Point{idx, y}
		}
		grid.Cells = append(grid.Cells, row)
	}
	grid.Set(start, '.')
	grid.Width = len(grid.Cells[0])
	grid.Height = len(grid.Cells)
	grid.Offset = grid.Width / 2
	fmt.Println("Width: ", grid.Width)
	fmt.Println("Height: ", grid.Height)
	fmt.Println("Offset: ", grid.Offset)
	return grid, Point{0, 0}
}

func solutionPart1(grid Grid, start Point) {
	currentPoints := map[Point]bool{start: true}
	for i := 0; i < 64; i++ {
		nextPoints := map[Point]bool{}
		for point := range currentPoints {
			for _, neighbour := range grid.getValidNeighbours(point) {
				nextPoints[neighbour] = true
			}
		}
		currentPoints = nextPoints
	}
	fmt.Println(len(currentPoints))
}

func (g *Grid) countPointsInGridWithOffset(offset Point, points map[Point]int) int {
	count := 0
	for point := range points {
		if point.X >= g.Width*offset.X-g.Offset && point.X <= g.Width*offset.X+g.Offset &&
			point.Y >= g.Height*offset.Y-g.Offset && point.Y <= g.Height*offset.Y+g.Offset {
			count++
		}
	}
	return count
}

func PowInts(x, n int) int {
	if n == 0 {
		return 1
	}
	if n == 1 {
		return x
	}
	y := PowInts(x, n/2)
	if n%2 == 0 {
		return y * y
	}
	return x * y * y
}

func solutionPart2(grid Grid, start Point) {
	currentPoints := map[Point]int{start: 0}
	for i := 0; i < grid.Height*2+grid.Offset; i++ {
		nextPoints := map[Point]int{}
		for point, val := range currentPoints {
			for _, neighbour := range grid.getValidNeighboursPart2(point) {
				nextPoints[neighbour] = val + 1
			}
		}
		currentPoints = nextPoints

	}

	countOdd := 0
	fmt.Println("(0,0)", grid.countPointsInGridWithOffset(Point{0, 0}, currentPoints))
	fmt.Println("(1,0)", grid.countPointsInGridWithOffset(Point{1, 0}, currentPoints))
	fullOdd := grid.countPointsInGridWithOffset(Point{0, 0}, currentPoints)
	fullEven := grid.countPointsInGridWithOffset(Point{1, 0}, currentPoints)
	oddEdgePoints := []Point{
		{-1, 1},
		{-1, -1},
		{1, -1},
		{1, 1},
	}
	for _, point := range oddEdgePoints {
		// fmt.Println(point, (fullOdd - grid.countPointsInGridWithOffset(point, currentPoints)))
		countOdd += (fullOdd - grid.countPointsInGridWithOffset(point, currentPoints))
	}
	fmt.Println("Count Odd: ", countOdd)

	countEven := 0
	evenEdgePoints := []Point{
		{2, 1},
		{2, -1},
		{-2, -1},
		{-2, 1},
	}
	for _, point := range evenEdgePoints {
		// fmt.Println(point, grid.countPointsInGridWithOffset(point, currentPoints))
		countEven += grid.countPointsInGridWithOffset(point, currentPoints)
	}

	n := 202300
	//n := 2
	log.Println("Points: ", PowInts(n, 2)*fullEven+PowInts(n+1, 2)*fullOdd-(n+1)*countOdd+n*countEven)

}

func main() {
	grid, start := readData("input.txt")
	solutionPart1(grid, start)
	grid, start = readData("input.txt")
	startTime := time.Now()
	solutionPart2(grid, start)
	log.Println(time.Since(startTime))
}
