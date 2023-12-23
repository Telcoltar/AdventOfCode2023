package main

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

type Direction int

const (
	R = iota
	L
)

type Node struct {
	Left  *Node
	Right *Node
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int) int {
	return a * b / gcd(a, b)
}

func loadData(filename string) ([]Direction, map[string]map[Direction]string, error) {
	data, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer func(data *os.File) {
		err := data.Close()
		if err != nil {
			panic(err)
		}
	}(data)
	scanner := bufio.NewScanner(data)
	scanner.Scan()
	directions := make([]Direction, 0)
	for _, d := range scanner.Text() {
		switch d {
		case 'R':
			directions = append(directions, R)
		case 'L':
			directions = append(directions, L)
		}
	}
	scanner.Scan()
	network := make(map[string]map[Direction]string)
	matcher, err := regexp.Compile(`\((\w{3}), (\w{3})\)`)
	if err != nil {
		return nil, nil, err
	}
	for scanner.Scan() {
		split := strings.Split(scanner.Text(), " = ")
		destinations := matcher.FindStringSubmatch(split[1])
		network[split[0]] = make(map[Direction]string)
		network[split[0]][L] = destinations[1]
		network[split[0]][R] = destinations[2]
	}
	return directions, network, nil
}

func followPath(start string, directions []Direction, network map[string]map[Direction]string, checkFn func(string) bool) int {
	currentNode := start
	steps := 0
	currentIndex := 0
	for checkFn(currentNode) {
		currentNode = network[currentNode][directions[currentIndex]]
		steps++
		currentIndex = (currentIndex + 1) % len(directions)
	}
	return steps
}

func solutionPart1() {
	directions, network, err := loadData("input.txt")
	if err != nil {
		panic(err)
	}
	steps := followPath("AAA", directions, network, func(node string) bool { return node != "ZZZ" })
	println(steps)
}

func solutionPart2() {
	directions, network, err := loadData("input.txt")
	if err != nil {
		panic(err)
	}
	cycles := make([]int, 0)
	for key := range network {
		if strings.HasSuffix(key, "A") {
			cycles = append(cycles, followPath(key, directions, network, func(node string) bool { return !strings.HasSuffix(node, "Z") }))
		}
	}
	currentLcm := 1
	for _, cycle := range cycles {
		currentLcm = lcm(currentLcm, cycle)
	}
	println(currentLcm)
}

func main() {
	solutionPart1()
	solutionPart2()
}
