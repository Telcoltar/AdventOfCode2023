package main

import (
	"bufio"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

func loadData(filename string) ([][]int, error) {
	result := make([][]int, 0)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		row := make([]int, 0)
		for _, numStr := range strings.Split(line, " ") {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				return nil, err
			}
			row = append(row, num)
		}
		result = append(result, row)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return result, nil
}

func processLine(line []int) [][]int {
	result := make([][]int, 0)
	result = append(result, line)
	for {
		currentLine := make([]int, 0)
		lastLine := result[len(result)-1]
		for i := 0; i < len(lastLine)-1; i++ {
			currentLine = append(currentLine, lastLine[i+1]-lastLine[i])
		}
		result = append(result, currentLine)
		// check if currentLine only contains 0
		allZero := true
		for _, num := range currentLine {
			if num != 0 {
				allZero = false
				break
			}
		}
		if allZero {
			break
		}
	}
	return result
}

func predictEnd(lines [][]int) int {
	sum := 0
	for _, line := range lines {
		sum += line[len(line)-1]
	}
	return sum
}

func predictBeginning(lines [][]int) int {
	acc := 0
	slices.Reverse(lines)
	for _, line := range lines {
		acc = line[0] - acc
	}
	return acc
}

func solution(predictFn func([][]int) int) {
	data, err := loadData("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	sum := 0
	for _, line := range data {
		processedLines := processLine(line)
		sum += predictFn(processedLines)
	}
	println(sum)
}

func solutionPart1() {
	solution(predictEnd)
}

func solutionPart2() {
	solution(predictBeginning)
}

func main() {
	solutionPart1()
	solutionPart2()
}
