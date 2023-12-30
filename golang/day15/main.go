package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func loadData(filename string) ([]string, error) {
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
	scanner := bufio.NewScanner(data)
	scanner.Scan()
	ops := strings.Split(scanner.Text(), ",")
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return ops, nil
}

func hash(input string) int {
	hashSum := int32(0)
	for _, c := range input {
		hashSum += c
		hashSum *= 17
		hashSum %= 256
	}
	return int(hashSum)
}

func solutionPart1() {
	ops, err := loadData("input.txt")
	if err != nil {
		panic(err)
	}
	totalSum := 0
	for _, op := range ops {
		totalSum += hash(op)
	}
	fmt.Println(totalSum)
}

type Box struct {
	indexMap     map[string]int
	values       []int
	currentIndex int
}

func newBox() Box {
	b := Box{
		indexMap:     make(map[string]int),
		values:       make([]int, 0),
		currentIndex: 0,
	}
	return b
}

func (b *Box) insert(label string, value int) {
	if index, ok := b.indexMap[label]; ok {
		b.values[index] = value
	} else {
		b.indexMap[label] = b.currentIndex
		b.values = append(b.values, value)
		b.currentIndex += 1
	}
}

func (b *Box) remove(label string) {
	if index, ok := b.indexMap[label]; ok {
		b.values[index] = -1
		delete(b.indexMap, label)
	}
}

func solutionPart2() {
	ops, err := loadData("input.txt")
	if err != nil {
		panic(err)
	}
	opMatcher, err := regexp.Compile("([a-z]+)([=-])(\\d*)")
	if err != nil {
		panic(err)
	}
	boxes := make([]Box, 256)
	for i := range boxes {
		boxes[i] = newBox()
	}
	for _, op := range ops {
		match := opMatcher.FindStringSubmatch(op)
		label := match[1]
		opType := match[2]
		if opType == "=" {
			value, err := strconv.Atoi(match[3])
			if err != nil {
				panic(err)
			}
			boxes[hash(label)].insert(label, value)
		} else {
			boxes[hash(label)].remove(label)
		}
	}
	totalSum := 0
	for index, b := range boxes {
		if b.currentIndex > 0 {
			boxSum := 0
			currentIndex := 1
			for _, v := range b.values {
				if v != -1 {
					boxSum += v * currentIndex
					currentIndex += 1
				}
			}
			totalSum += boxSum * (index + 1)
		}
	}
	fmt.Println(totalSum)
}

func main() {
	solutionPart1()
	solutionPart2()
}
