package main

import (
	"bufio"
	"os"
	"strconv"
	"unicode"
)

type Point struct {
	x int
	y int
}
type Symbol struct {
	value    rune
	location Point
}

type Number struct {
	Value   int
	Symbols []Symbol
}

func padLine(line []rune, padding rune) []rune {
	return append(append([]rune{padding}, line...), padding)
}

func filledLine(width int, fill rune) []rune {
	line := make([]rune, width)
	for i := 0; i < width; i++ {
		line[i] = fill
	}
	return line
}

func loadBoard(filename string) ([]Number, error) {
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
	numbers := make([]Number, 0)
	paddedBoard := make([][]rune, 0)
	// get first line to get width
	scanner.Scan()
	firstLine := scanner.Text()
	width := len(firstLine)
	paddedBoard = append(paddedBoard, filledLine(width, '.'))
	paddedBoard = append(paddedBoard, padLine([]rune(firstLine), '.'))
	for scanner.Scan() {
		line := scanner.Text()
		paddedBoard = append(paddedBoard, padLine([]rune(line), '.'))
	}
	paddedBoard = append(paddedBoard, filledLine(width, '.'))
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	// find numbers
	for y, line := range paddedBoard {
		numberInProgress := false
		var currentNumber []rune
		var start Point
		for x, char := range line {
			if unicode.IsDigit(char) {
				if numberInProgress {
					currentNumber = append(currentNumber, char)
				} else {
					numberInProgress = true
					currentNumber = make([]rune, 0)
					currentNumber = append(currentNumber, char)
					start = Point{x: x - 1, y: y - 1}
				}
			} else {
				if numberInProgress {
					parsedNumber, err := strconv.Atoi(string(currentNumber))
					if err != nil {
						return nil, err
					}
					numbers = append(numbers, Number{Value: parsedNumber, Symbols: getSymbols(paddedBoard, start, Point{x, y + 1})})
					numberInProgress = false
				}
			}
		}
	}
	return numbers, nil
}

func getSymbols(board [][]rune, start Point, end Point) []Symbol {
	symbols := make([]Symbol, 0)
	for i := start.x; i <= end.x; i++ {
		if board[start.y][i] != '.' && !unicode.IsDigit(board[start.y][i]) {
			symbols = append(symbols, Symbol{value: board[start.y][i], location: Point{x: i, y: start.y}})
		}
		if board[end.y][i] != '.' && !unicode.IsDigit(board[end.y][i]) {
			symbols = append(symbols, Symbol{value: board[end.y][i], location: Point{x: i, y: end.y}})
		}
	}
	if board[start.y+1][start.x] != '.' && !unicode.IsDigit(board[start.y+1][start.x]) {
		symbols = append(symbols, Symbol{value: board[start.y+1][start.x], location: Point{x: start.x, y: start.y + 1}})
	}
	if board[start.y+1][end.x] != '.' && !unicode.IsDigit(board[start.y+1][end.x]) {
		symbols = append(symbols, Symbol{value: board[start.y+1][end.x], location: Point{x: end.x, y: start.y + 1}})
	}
	return symbols
}

func solutionPart1() {
	board, err := loadBoard("input.txt")
	if err != nil {
		panic(err)
	}
	partNumbersSum := 0
	for _, number := range board {
		if len(number.Symbols) > 0 {
			partNumbersSum += number.Value
		}
	}
	println(partNumbersSum)
}

func solutionPart2() {
	board, err := loadBoard("input.txt")
	if err != nil {
		panic(err)
	}
	gearNumbers := make(map[Point][]int)
	for _, number := range board {
		for _, symbol := range number.Symbols {
			if symbol.value == '*' {
				gearNumbers[symbol.location] = append(gearNumbers[symbol.location], number.Value)
				break
			}
		}
	}
	gearRatioSum := 0
	for _, gearNumber := range gearNumbers {
		if len(gearNumber) == 2 {
			gearRatioSum += gearNumber[0] * gearNumber[1]
		}
	}
	println(gearRatioSum)
}

func main() {
	solutionPart1()
	solutionPart2()
}
