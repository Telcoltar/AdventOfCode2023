package main

import (
	"bytes"
	"container/list"
	"fmt"
	"maps"
	"os"
	"slices"
	"time"
)

type Point struct {
	X, Y int
}

type Grid struct {
	data          [][]rune
	height, width int
}

func (grid Grid) Get(p Point) rune {
	return grid.data[p.Y][p.X]
}

func (grid Grid) Set(p Point, value rune) {
	grid.data[p.Y][p.X] = value
}

var DirMap = map[rune]Point{
	'>': {1, 0},
	'<': {-1, 0},
	'^': {0, -1},
	'v': {0, 1},
}

var dirs = slices.Collect(maps.Keys(DirMap))

func (grid Grid) ValidNeighbours(p Point) []Point {
	neighbours := []Point{}
	if slices.Contains(dirs, grid.Get(p)) {
		neighbours = append(neighbours, Point{p.X + DirMap[grid.Get(p)].X, p.Y + DirMap[grid.Get(p)].Y})
		return neighbours
	}
	for _, neighbour := range []Point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} {
		newPoint := Point{p.X + neighbour.X, p.Y + neighbour.Y}
		if newPoint.X >= 0 && newPoint.X < grid.width && newPoint.Y >= 0 && newPoint.Y < grid.height && grid.Get(newPoint) != '#' {
			neighbours = append(neighbours, newPoint)
		}
	}
	return neighbours
}

func (grid Grid) ValidNeighboursP2(p Point) []Point {
	neighbours := []Point{}
	for _, neighbour := range []Point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} {
		newPoint := Point{p.X + neighbour.X, p.Y + neighbour.Y}
		if newPoint.X >= 0 && newPoint.X < grid.width && newPoint.Y >= 0 && newPoint.Y < grid.height && grid.Get(newPoint) != '#' {
			neighbours = append(neighbours, newPoint)
		}
	}
	return neighbours
}

func (grid Grid) String() string {
	var buffer bytes.Buffer
	for _, row := range grid.data {
		buffer.WriteString(string(row) + "\n")
	}
	return buffer.String()
}

func (grid Grid) PrintPath(path map[Point]struct{}) {
	for y, row := range grid.data {
		for x, cell := range row {
			if _, ok := path[Point{x, y}]; ok {
				fmt.Print("O")
			} else {
				fmt.Print(string(cell))
			}
		}
		fmt.Println()
	}
}

func readData(filepath string) Grid {
	fileContent, readErr := os.ReadFile(filepath)
	if readErr != nil {
		panic(readErr)
	}
	field := Grid{data: [][]rune{}}
	for _, line := range bytes.Split(fileContent, []byte("\n")) {
		field.data = append(field.data, []rune(string(line)))
	}
	field.height = len(field.data)
	field.width = len(field.data[0])
	return field
}

type QueueItem struct {
	Point
	path map[Point]struct{}
}

func (qi QueueItem) Copy() QueueItem {
	newPath := map[Point]struct{}{}
	for k, v := range qi.path {
		newPath[k] = v
	}
	return QueueItem{qi.Point, newPath}
}

func solutionPart1(data Grid) int {
	// fmt.Printf("%v\n", data)
	queue := []QueueItem{{Point{1, 0}, map[Point]struct{}{{1, 0}: {}}}}
	weightMap := map[Point]int{{1, 0}: 0}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		neighbours := data.ValidNeighbours(current.Point)
		neighbours = slices.DeleteFunc(neighbours, func(p Point) bool {
			_, ok := current.path[p]
			return ok
		})
		neighbours = slices.DeleteFunc(neighbours, func(p Point) bool {
			return weightMap[p] >= weightMap[current.Point]+1
		})
		if len(neighbours) == 1 {
			weightMap[neighbours[0]] = weightMap[current.Point] + 1
			current.Point = neighbours[0]
			current.path[neighbours[0]] = struct{}{}
			queue = append(queue, current)
		} else {
			for _, neighbour := range neighbours {
				newPath := current.Copy()
				newPath.Point = neighbour
				newPath.path[neighbour] = struct{}{}
				queue = append(queue, newPath)
				weightMap[neighbour] = weightMap[current.Point] + 1
			}
		}
	}
	return weightMap[Point{data.width - 2, data.height - 1}]
}

func solutionPart2(data Grid) int {
	// fmt.Printf("%v\n", data)
	queue := []QueueItem{{Point{1, 0}, map[Point]struct{}{{1, 0}: {}}}}
	weightMap := map[Point]int{{1, 0}: 0}
	end := Point{data.width - 2, data.height - 1}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		neighbours := data.ValidNeighboursP2(current.Point)
		neighbours = slices.DeleteFunc(neighbours, func(p Point) bool {
			_, ok := current.path[p]
			return ok
		})
		neighbours = slices.DeleteFunc(neighbours, func(p Point) bool {
			return weightMap[p] >= weightMap[current.Point]+1
		})
		if len(neighbours) == 1 {
			weightMap[neighbours[0]] = weightMap[current.Point] + 1
			current.Point = neighbours[0]
			current.path[neighbours[0]] = struct{}{}
			queue = append(queue, current)
		} else {
			for _, neighbour := range neighbours {
				newPath := current.Copy()
				newPath.Point = neighbour
				newPath.path[neighbour] = struct{}{}
				queue = append(queue, newPath)
				weightMap[neighbour] = weightMap[current.Point] + 1
			}
		}
	}

	return weightMap[end]
}

type PointDistance struct {
	Point
	Distance int
}

func buildGraph(grid Grid, start Point, end Point) map[Point][]PointDistance {
	fmt.Println(start, end)
	graph := map[Point][]PointDistance{}
	queue := []Point{start}
	visited := map[Point]bool{}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if visited[current] {
			continue
		}
		visited[current] = true
		neighbours := grid.ValidNeighboursP2(current)
		for _, neighbour := range neighbours {
			prevInner := current
			innerCurrent := neighbour
			deadEnd := false
			distance := 0
			for {
				distance++
				currentNeighbours := grid.ValidNeighboursP2(innerCurrent)
				currentNeighbours = slices.DeleteFunc(currentNeighbours, func(p Point) bool {
					return p == prevInner
				})
				if len(currentNeighbours) == 1 {
					prevInner = innerCurrent
					innerCurrent = currentNeighbours[0]
				} else if len(currentNeighbours) == 0 {
					if innerCurrent == end {
						break
					}
					deadEnd = true
					break
				} else {
					break
				}
			}
			if !deadEnd {
				graph[current] = append(graph[current], PointDistance{innerCurrent, distance})
				queue = append(queue, innerCurrent)
			}
		}
	}
	return graph
}

func dfsLongestPath(graph map[Point][]PointDistance, current, goal Point, visited map[Point]bool) int {
	if current == goal {
		return 0
	}
	visited[current] = true
	maxLen := -1
	for _, next := range graph[current] {
		if visited[next.Point] {
			continue
		}
		subPath := dfsLongestPath(graph, next.Point, goal, visited)
		if subPath >= 0 {
			total := next.Distance + subPath
			if total > maxLen {
				maxLen = total
			}
		}
	}
	visited[current] = false
	return maxLen
}

type setBackIndicator struct {
	point Point
}

func dfsStack(graph map[Point][]PointDistance, end Point) int {
	stack := list.New()
	visited := map[Point]bool{}
	start := PointDistance{Point{1, 0}, 0}
	stack.PushBack(start)
	maxLen := -1
	for stack.Len() > 0 {
		last := stack.Back()
		stack.Remove(last)
		var current PointDistance
		switch val := last.Value.(type) {
		case setBackIndicator:
			visited[val.point] = false
			continue
		case PointDistance:
			current = val
		}
		if current.Point == end {
			if current.Distance > maxLen {
				maxLen = current.Distance
			}
			continue
		}
		visited[current.Point] = true
		stack.PushBack(setBackIndicator{current.Point})
		for _, next := range graph[current.Point] {
			if visited[next.Point] {
				continue
			}
			stack.PushBack(PointDistance{next.Point, next.Distance + current.Distance})
		}
	}
	return maxLen
}

func main() {
	start := time.Now()
	data := readData("input.txt")
	println(solutionPart1(data))
	graph := buildGraph(data, Point{1, 0}, Point{data.width - 2, data.height - 1})
	visited := map[Point]bool{}
	println(dfsLongestPath(graph, Point{1, 0}, Point{data.width - 2, data.height - 1}, visited))
	fmt.Println("----------")
	println(dfsStack(graph, Point{data.width - 2, data.height - 1}))
	fmt.Printf("Execution time: %v\n", time.Since(start))
}
