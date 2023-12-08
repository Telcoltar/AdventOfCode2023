package main

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"unicode"
)

var NUMBERS = map[string]string{
	"one":   "1",
	"two":   "2",
	"three": "3",
	"four":  "4",
	"five":  "5",
	"six":   "6",
	"seven": "7",
	"eight": "8",
	"nine":  "9",
}

func solutionPart1() error {
	file, err := os.Open("input.txt")
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(file)
	lineNumbers := make([]int, 0)
	for scanner.Scan() {
		line := []rune(scanner.Text())
		lineNumber := make([]rune, 2)
		for _, c := range line {
			if unicode.IsDigit(c) {
				lineNumber[0] = c
				break
			}
		}
		for i := len(line) - 1; i >= 0; i-- {
			if unicode.IsDigit(line[i]) {
				lineNumber[1] = line[i]
				break
			}
		}
		parsedNumber, err := strconv.Atoi(string(lineNumber))
		if err != nil {
			return err
		}
		lineNumbers = append(lineNumbers, parsedNumber)
	}
	sum := 0
	for _, lineNumber := range lineNumbers {
		sum += lineNumber
	}
	println(sum)
	return nil
}

func solutionPart2() error {
	file, err := os.Open("input.txt")
	if err != nil {
		return err
	}
	lineNumbers := make([]int, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineNumber := ""
		reDigits, err := regexp.Compile(`(one|two|three|four|five|six|seven|eight|nine|\d)`)
		if err != nil {
			return err
		}
		reReversedDigits, err := regexp.Compile(`(eno|owt|eerht|ruof|evif|xis|neves|thgie|enin|\d)`)
		if err != nil {
			return err
		}
		line := scanner.Text()
		firstDigit := reDigits.FindString(line)
		if len(firstDigit) > 1 {
			lineNumber += NUMBERS[firstDigit]
		} else {
			lineNumber += firstDigit
		}
		reversedLine := []rune(line)
		for i, j := 0, len(reversedLine)-1; i < j; i, j = i+1, j-1 {
			reversedLine[i], reversedLine[j] = reversedLine[j], reversedLine[i]
		}
		lastDigitReversed := reReversedDigits.FindString(string(reversedLine))
		lastDigit := []rune(lastDigitReversed)
		for i, j := 0, len(lastDigit)-1; i < j; i, j = i+1, j-1 {
			lastDigit[i], lastDigit[j] = lastDigit[j], lastDigit[i]
		}
		if len(lastDigit) > 1 {
			lineNumber += NUMBERS[string(lastDigit)]
		} else {
			lineNumber += lastDigitReversed
		}
		parsedNumber, err := strconv.Atoi(lineNumber)
		if err != nil {
			return err
		}
		lineNumbers = append(lineNumbers, parsedNumber)
	}
	sum := 0
	for _, lineNumber := range lineNumbers {
		sum += lineNumber
	}
	println(sum)
	return nil
}

func main() {
	err := solutionPart2()
	if err != nil {
		panic(err)
	}
}
