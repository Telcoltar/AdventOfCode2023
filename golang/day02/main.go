package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

var MAX_CUBES = map[string]int{
	"red":   12,
	"green": 13,
	"blue":  14,
}

func load_game(filename string) (map[int]map[string]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	gameData := make(map[int]map[string]int)
	for scanner.Scan() {
		line := scanner.Text()
		lineSplit := strings.Split(line, ":")
		gameSplit := strings.Split(lineSplit[0], " ")
		gameID, err := strconv.Atoi(gameSplit[1])
		if err != nil {
			return nil, err
		}
		maxCubes := map[string]int{
			"red":   0,
			"green": 0,
			"blue":  0,
		}
		for _, draw := range strings.Split(lineSplit[1], ";") {
			for _, cube := range strings.Split(draw, ",") {
				cubeSplit := strings.Split(strings.TrimSpace(cube), " ")
				cubeColor := cubeSplit[1]
				cubeCount, err := strconv.Atoi(cubeSplit[0])
				if err != nil {
					return nil, err
				}
				maxCubes[cubeColor] = max(maxCubes[cubeColor], cubeCount)
			}
		}
		gameData[gameID] = maxCubes

	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return gameData, nil
}

func solutionPart1() {
	gameData, err := load_game("input.txt")
	if err != nil {
		panic(err)
	}
	possibleGamesSum := 0
	for gameIndex, game := range gameData {
		if game["red"] <= MAX_CUBES["red"] &&
			game["green"] <= MAX_CUBES["green"] &&
			game["blue"] <= MAX_CUBES["blue"] {
			possibleGamesSum += gameIndex
		}
	}
	println(possibleGamesSum)
}

func solutionPart2() {
	gameData, err := load_game("input.txt")
	if err != nil {
		panic(err)
	}
	gamePowerSum := 0
	for _, game := range gameData {
		gamePowerSum += game["red"] * game["blue"] * game["green"]
	}
	println(gamePowerSum)
}

func main() {
	solutionPart1()
	solutionPart2()
}
