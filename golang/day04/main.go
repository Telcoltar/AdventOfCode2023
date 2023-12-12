package main

import (
	"bufio"
	"github.com/hashicorp/go-set"
	"os"
	"strconv"
	"strings"
)

type Card struct {
	WinningNumbers *set.Set[int]
	MyNumbers      *set.Set[int]
}

func parseNumberList(list string) (*set.Set[int], error) {
	numbers := set.New[int](0)
	for _, number := range strings.Split(list, " ") {
		if number == "" {
			continue
		}
		numberInt, err := strconv.Atoi(number)
		if err != nil {
			return nil, err
		}
		numbers.Insert(numberInt)
	}
	return numbers, nil
}

func loadData(filename string) ([]Card, error) {
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
	cards := []Card{}
	for scanner.Scan() {
		line := scanner.Text()
		lineSplit := strings.Split(line, ":")
		numbersSplit := strings.Split(lineSplit[1], "|")
		winningNumbers, err := parseNumberList(numbersSplit[0])
		if err != nil {
			return nil, err
		}
		myNumbers, err := parseNumberList(numbersSplit[1])
		if err != nil {
			return nil, err
		}
		cards = append(cards, Card{MyNumbers: myNumbers, WinningNumbers: winningNumbers})
	}
	return cards, nil
}

// implement int power function because golang has no std function for it
func pow(a int, b int) int {
	result := 1
	for i := 0; i < b; i++ {
		result *= a
	}
	return result
}

func solutionPart1() {
	cards, err := loadData("input.txt")
	if err != nil {
		panic(err)
	}
	score := 0
	for _, card := range cards {
		matches := card.MyNumbers.Intersect(card.WinningNumbers).Size()
		if matches > 0 {
			score += pow(2, matches-1)
		}
	}
	println(score)
}

func solutionPart2() {
	cards, err := loadData("input.txt")
	if err != nil {
		panic(err)
	}
	pile := make([]int, len(cards))
	for i := range pile {
		pile[i] = 1
	}
	for index, card := range cards {
		matches := card.MyNumbers.Intersect(card.WinningNumbers).Size()
		for i := 0; i < matches; i++ {
			pile[index+i+1] += pile[index]
		}
	}
	sum := 0
	for _, count := range pile {
		sum += count
	}
	println(sum)
}

func main() {
	solutionPart1()
	solutionPart2()
}
