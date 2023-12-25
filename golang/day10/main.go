package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

type Direction int

const (
	up = iota
	down
	left
	right
)

func (d Direction) opposite() Direction {
	if d == up {
		return down
	}
	if d == down {
		return up
	}
	if d == left {
		return right
	}
	if d == right {
		return left
	}
	return -1
}

type Tile struct {
	symbol      rune
	connections map[Direction]bool
	path        map[Direction]Direction
}

func newTile(symbol rune) Tile {
	connections := make(map[Direction]bool)
	for i := 0; i < 4; i++ {
		connections[Direction(i)] = false
	}
	connectionsList := make([]Direction, 0)
	path := make(map[Direction]Direction, 2)
	switch symbol {
	case '|':
		connectionsList = append(connectionsList, up, down)
	case '-':
		connectionsList = append(connectionsList, left, right)
	case 'L':
		connectionsList = append(connectionsList, up, right)
	case 'J':
		connectionsList = append(connectionsList, up, left)
	case '7':
		connectionsList = append(connectionsList, left, down)
	case 'F':
		connectionsList = append(connectionsList, right, down)
	}
	if len(connectionsList) > 0 {
		for _, d := range connectionsList {
			connections[d] = true
		}
		path[connectionsList[0]] = connectionsList[1]
		path[connectionsList[1]] = connectionsList[0]
	}
	return Tile{
		symbol:      symbol,
		connections: connections,
		path:        path,
	}
}

type Point struct {
	x int
	y int
}

func (p *Point) getNeighbourInDirection(dir Direction) Point {
	if dir == up {
		return Point{p.x, p.y - 1}
	}
	if dir == down {
		return Point{p.x, p.y + 1}
	}
	if dir == left {
		return Point{p.x - 1, p.y}
	}
	if dir == right {
		return Point{p.x + 1, p.y}
	}
	return Point{p.x, p.y}
}

func loadData(filename string) ([][]Tile, error) {
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
	scanner := bufio.NewScanner(data)
	tileMap := make([][]Tile, 0)
	for scanner.Scan() {
		line := make([]Tile, 0)
		for _, tile := range scanner.Text() {
			line = append(line, newTile(tile))
		}
		tileMap = append(tileMap, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	tileMap = padMap(tileMap)
	return tileMap, nil
}

func padMap(tileMap [][]Tile) [][]Tile {
	paddedMap := make([][]Tile, len(tileMap)+2)
	for i := range paddedMap {
		paddedMap[i] = make([]Tile, len(tileMap[0])+2)
	}
	for i := 0; i < len(tileMap[0])+2; i++ {
		paddedMap[0][i] = newTile('.')
	}
	for i := 0; i < len(tileMap); i++ {
		for j := 0; j < len(tileMap[0]); j++ {
			paddedMap[i+1][j+1] = tileMap[i][j]
		}
	}
	for i := 0; i < len(tileMap[0])+2; i++ {
		paddedMap[len(tileMap)+1][i] = newTile('.')
	}
	return paddedMap
}

func findStart(tileMap [][]Tile) (Point, error) {
	for y, row := range tileMap {
		for x, tile := range row {
			if tile.symbol == 'S' {
				return Point{x, y}, nil
			}
		}
	}
	return Point{}, fmt.Errorf("no start found")
}

func replaceStart(tileMap [][]Tile, start Point) {
	connections := make([]int, 4)
	for i := 0; i < 4; i++ {
		neighbour := start.getNeighbourInDirection(Direction(i))
		if tileMap[neighbour.y][neighbour.x].connections[Direction(i).opposite()] {
			connections[i] = 1
		}
	}
	code := ""
	for _, conn := range connections {
		code += fmt.Sprintf("%d", conn)
	}
	var symbol rune
	switch code {
	case "1100":
		symbol = '|'
	case "0011":
		symbol = '_'
	case "1001":
		symbol = 'L'
	case "1010":
		symbol = 'J'
	case "0110":
		symbol = '7'
	case "0101":
		symbol = 'F'
	}
	tileMap[start.y][start.x] = newTile(symbol)
}

func followPath(tileMap [][]Tile, start Point) []Point {
	currentTile := start
	startTile := tileMap[start.y][start.x].connections
	pathTiles := make([]Point, 0)
	var currentDirection Direction
	for k, v := range startTile {
		if v {
			currentDirection = k
			break
		}
	}
	for {
		currentTile = currentTile.getNeighbourInDirection(currentDirection)
		currentDirection = tileMap[currentTile.y][currentTile.x].path[currentDirection.opposite()]
		pathTiles = append(pathTiles, currentTile)
		if currentTile == start {
			break
		}
	}
	return pathTiles
}

func scanRow(tileMap [][]Tile, path []Point) int {
	interiorSum := 0
	height := len(tileMap)
	width := len(tileMap[0])
	for y := 1; y < height-1; y++ {
		isInside := false
		x := 1
		for x < width-1 {
			currentPoint := Point{x: x, y: y}
			currentTile := tileMap[y][x]
			if !slices.Contains(path, currentPoint) {
				if isInside {
					interiorSum += 1
				}
			} else if currentTile.symbol == '|' {
				isInside = !isInside
			} else if currentTile.symbol == 'F' {
				x++
				for tileMap[y][x].symbol == '-' {
					x++
				}
				if tileMap[y][x].symbol == 'J' {
					isInside = !isInside
				}
			} else if currentTile.symbol == 'L' {
				x++
				for tileMap[y][x].symbol == '-' {
					x++
				}
				if tileMap[y][x].symbol == '7' {
					isInside = !isInside
				}
			}
			x++
		}
	}
	return interiorSum
}

func solutionPart1() {
	tileMap, err := loadData("input.txt")
	if err != nil {
		panic(err)
	}
	start, err := findStart(tileMap)
	if err != nil {
		panic(err)
	}
	replaceStart(tileMap, start)
	pathPoints := followPath(tileMap, start)
	fmt.Println(len(pathPoints) / 2)
}

func solutionPart2() {
	tileMap, err := loadData("input.txt")
	if err != nil {
		panic(err)
	}
	start, err := findStart(tileMap)
	if err != nil {
		panic(err)
	}
	replaceStart(tileMap, start)
	pathPoints := followPath(tileMap, start)
	interiorSum := scanRow(tileMap, pathPoints)
	println(interiorSum)
}

func main() {
	solutionPart1()
	solutionPart2()
}
