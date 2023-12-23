package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type TypeOfHand int

const (
	HighCard TypeOfHand = iota + 1
	OnePair
	TwoPair
	ThreeOfAKind
	FullHouse
	FourOfAKind
	FiveOfAKind
)

type Hand struct {
	Hand     []rune
	HandType TypeOfHand
	bid      int
}

func loadData(filename string) ([]string, []int, error) {
	data, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	scanner := bufio.NewScanner(data)
	defer func(data *os.File) {
		err := data.Close()
		if err != nil {

		}
	}(data)
	var hands []string
	var bids []int
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		hand := line[0]
		bid, err := strconv.Atoi(line[1])
		if err != nil {
			return nil, nil, err
		}
		hands = append(hands, hand)
		bids = append(bids, bid)
	}
	return hands, bids, nil
}

func countCards(hand string) map[rune]int {
	cards := make(map[rune]int)
	for _, card := range hand {
		cards[card]++
	}
	return cards
}

func parseValues(values []int) TypeOfHand {
	switch {
	case values[0] == 5:
		return FiveOfAKind
	case values[0] == 4:
		return FourOfAKind
	case values[0] == 3 && values[1] == 2:
		return FullHouse
	case values[0] == 3:
		return ThreeOfAKind
	case values[0] == 2 && values[1] == 2:
		return TwoPair
	case values[0] == 2:
		return OnePair
	}
	return HighCard
}

func parseHandPart1(hand string) TypeOfHand {
	count := countCards(hand)
	values := make([]int, 0, len(count))
	for _, value := range count {
		values = append(values, value)
	}
	slices.Sort(values)
	slices.Reverse(values)
	return parseValues(values)
}

func parseHandPart2(hand string) TypeOfHand {
	count := countCards(hand)
	jCount := count['J']
	if jCount == 4 || jCount == 5 {
		return FiveOfAKind
	}
	delete(count, 'J')
	values := make([]int, 0, len(count))
	for _, value := range count {
		values = append(values, value)
	}
	slices.Sort(values)
	slices.Reverse(values)
	values[0] += jCount
	return parseValues(values)
}

func createCmpFunction(jValue int) func(hand1, hand2 Hand) int {
	return func(hand1, hand2 Hand) int {
		if hand1.HandType == hand2.HandType {
			for i := 0; i < len(hand1.Hand); i++ {
				if cardValue(hand1.Hand[i], jValue) == cardValue(hand2.Hand[i], jValue) {
					continue
				}
				return cardValue(hand1.Hand[i], jValue) - cardValue(hand2.Hand[i], jValue)
			}
			return 0
		}
		return int(hand1.HandType - hand2.HandType)
	}
}

func cardValue(card rune, jValue int) int {
	switch card {
	case 'A':
		return 14
	case 'K':
		return 13
	case 'Q':
		return 12
	case 'J':
		return jValue
	case 'T':
		return 10
	}
	return int(card - '0')
}

func solution(filename string, parseFunc func(string) TypeOfHand, jValue int) {
	handStrs, bids, err := loadData(filename)
	if err != nil {
		fmt.Println("Error loading data:", err)
		return
	}
	hands := make([]Hand, len(bids))
	for index, handStr := range handStrs {
		hands[index] = Hand{[]rune(handStr), parseFunc(handStr), bids[index]}
	}
	slices.SortFunc(hands, createCmpFunction(jValue))
	winnings := 0
	for index, hand := range hands {
		winnings += (index + 1) * hand.bid
	}
	fmt.Printf("Total: %+v\n", winnings)
}

func solutionPart1() {
	solution("input.txt", parseHandPart1, 11)
}

func solutionPart2() {
	solution("input.txt", parseHandPart2, 1)
}

func main() {
	solutionPart1()
	solutionPart2()
}
