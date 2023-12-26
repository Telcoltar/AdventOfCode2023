package main

import (
	"bufio"
	"os"
	"sort"
)

type Point = map[string]int

func loadData(filename string) ([]Point, error) {
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
	var galaxies []Point
	scanner := bufio.NewScanner(data)
	y := 0
	for scanner.Scan() {
		for x, char := range scanner.Text() {
			if char == '#' {
				galaxies = append(galaxies, Point{"x": x, "y": y})
			}
		}
		y += 1
	}
	return galaxies, nil
}

func spreadGalaxies(galaxies []Point, spreadFactor int) {
	for _, coord := range []string{"x", "y"} {
		sort.Slice(galaxies, func(i, j int) bool {
			return galaxies[i][coord] < galaxies[j][coord]
		})
		currentSpread := 0
		currentPos := 0
		for _, galaxy := range galaxies {
			diff := galaxy[coord] - currentPos
			currentPos = galaxy[coord]
			if diff > 1 {
				currentSpread += (diff - 1) * (spreadFactor - 1)
			}
			galaxy[coord] += currentSpread
		}
	}
}

func abs(num int) int {
	if num < 0 {
		return -num
	}
	return num
}

func sumDistances(galaxies []Point) int {
	sum := 0
	for i := 0; i < len(galaxies); i++ {
		for j := i + 1; j < len(galaxies); j++ {
			sum += abs(galaxies[i]["x"]-galaxies[j]["x"]) + abs(galaxies[i]["y"]-galaxies[j]["y"])
		}
	}
	return sum
}

func solution(filename string, spreadFactor int) {
	galaxies, err := loadData(filename)
	if err != nil {
		panic(err)
	}
	spreadGalaxies(galaxies, spreadFactor)
	totalDistance := sumDistances(galaxies)
	println(totalDistance)
}

func solutionPart1() {
	solution("input.txt", 2)
}

func solutionPart2() {
	solution("input.txt", 1000000)
}

func main() {
	solutionPart1()
	solutionPart2()
}
