package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Condition = int

const (
	Operational = iota
	Damaged
	Unknown
)

type Row struct {
	Conditions []Condition
	Blocks     []int
}

type Process struct {
	Row
	Invocations int
	Cache       [][]int
}

func (p *Process) processRow(condIndex int, blockIndex int) int {
	if result := p.Cache[condIndex][blockIndex]; result != -1 {
		return result
	}
	p.Invocations += 1
	if condIndex >= len(p.Conditions) {
		if blockIndex >= len(p.Blocks) {
			return 1
		}
		return 0
	}
	condition := p.Conditions[condIndex]
	result := 0
	switch condition {
	case Operational:
		result = p.processRow(condIndex+1, blockIndex)
	case Damaged:
		if blockIndex < len(p.Blocks) && p.isSpaceForBlock(condIndex, p.Blocks[blockIndex]) {
			result = p.processRow(condIndex+p.Blocks[blockIndex]+1, blockIndex+1)
		}
	case Unknown:
		result = p.processRow(condIndex+1, blockIndex)
		if blockIndex < len(p.Blocks) && p.isSpaceForBlock(condIndex, p.Blocks[blockIndex]) {
			result += p.processRow(condIndex+p.Blocks[blockIndex]+1, blockIndex+1)
		}
	}
	p.Cache[condIndex][blockIndex] = result
	return result
}

func (p *Process) isSpaceForBlock(condIndex int, block int) bool {
	currentCondIndex := condIndex
	currentBlock := block
	for currentCondIndex < len(p.Conditions) && currentBlock > 0 &&
		(p.Conditions[currentCondIndex] == Damaged || p.Conditions[currentCondIndex] == Unknown) {
		currentBlock -= 1
		currentCondIndex += 1
	}
	if currentBlock == 0 && (currentCondIndex == len(p.Conditions) ||
		p.Conditions[currentCondIndex] == Operational || p.Conditions[currentCondIndex] == Unknown) {
		return true
	}
	return false
}

func loadData(filename string) ([]Row, error) {
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
	rows := make([]Row, 0)
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		conditions := make([]Condition, 0)
		for _, char := range line[0] {
			switch char {
			case '#':
				conditions = append(conditions, Damaged)
			case '.':
				conditions = append(conditions, Operational)
			case '?':
				conditions = append(conditions, Unknown)
			}
		}
		blocks := strings.Split(line[1], ",")
		row := Row{
			Conditions: conditions,
			Blocks:     make([]int, len(blocks)),
		}
		for i, block := range blocks {
			num, err := strconv.Atoi(block)
			if err != nil {
				return nil, err
			}
			row.Blocks[i] = num
		}
		rows = append(rows, row)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return rows, nil
}

func multiplyRow(row *Row, factor int) {
	originalConditions := row.Conditions
	originalBlocks := row.Blocks
	for i := 0; i < factor-1; i++ {
		row.Conditions = append(row.Conditions, Unknown)
		row.Conditions = append(row.Conditions, originalConditions...)
		row.Blocks = append(row.Blocks, originalBlocks...)
	}
}

func solution(factor int) {
	rows, err := loadData("input.txt")
	if err != nil {
		panic(err)
	}
	totalSum := 0
	for _, row := range rows {
		multiplyRow(&row, factor)
		cache := make([][]int, len(row.Conditions)+2)
		for i := 0; i < len(row.Conditions)+2; i++ {
			cacheRow := make([]int, len(row.Blocks)+2)
			for j := 0; j < len(row.Blocks)+2; j++ {
				cacheRow[j] = -1
			}
			cache[i] = cacheRow
		}
		p := Process{row, 0, cache}
		totalSum += p.processRow(0, 0)
	}
	fmt.Printf("Sum: %d\n", totalSum)
}

func solutionPart1() {
	solution(1)
}

func solutionPart2() {
	solution(5)
}
func main() {
	solutionPart1()
	solutionPart2()
}
