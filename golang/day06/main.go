package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func loadDataPart1(filename string) ([]int, []int, error) {
	data, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	scanner := bufio.NewScanner(data)
	scanner.Scan()
	times, err := parseNumbersString(scanner.Text())
	if err != nil {
		return nil, nil, err
	}
	scanner.Scan()
	distances, err := parseNumbersString(scanner.Text())
	if err != nil {
		return nil, nil, err
	}
	return distances, times, nil
}

func loadDataPart2(filename string) (int, int, error) {
	data, err := os.Open(filename)
	if err != nil {
		return 0, 0, err
	}
	scanner := bufio.NewScanner(data)
	scanner.Scan()
	time, err := strconv.Atoi(strings.Replace(strings.Split(scanner.Text(), ":")[1], " ", "", -1))
	if err != nil {
		return 0, 0, err
	}
	scanner.Scan()
	distance, err := strconv.Atoi(strings.Replace(strings.Split(scanner.Text(), ":")[1], " ", "", -1))
	if err != nil {
		return 0, 0, err
	}
	return time, distance, nil
}

func parseNumbersString(listString string) ([]int, error) {
	numberList := make([]int, 0)
	for _, numberStr := range strings.Split(strings.Split(listString, ":")[1], " ") {
		if numberStr == "" {
			continue
		}
		number, err := strconv.Atoi(numberStr)
		if err != nil {
			return nil, err
		}
		numberList = append(numberList, number)
	}
	return numberList, nil
}

func calculateRange(time float64, distance float64) int {
	p_2 := time / 2.0
	d := math.Sqrt(math.Pow(time/2.0, 2.0) - distance)
	zero_1 := p_2 - d
	zero_2 := p_2 + d
	high := math.Floor(zero_2)
	low := math.Ceil(zero_1)
	if high == zero_2 {
		high -= 1
	}
	if low == zero_1 {
		low += 1
	}
	return int(high) - int(low) + 1
}

func solutionPart1() {
	distances, times, err := loadDataPart1("input.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	ranges := make([]int, 0)
	for i := 0; i < len(times); i++ {
		ranges = append(ranges, calculateRange(float64(times[i]), float64(distances[i])))
	}
	product := 1
	for _, r := range ranges {
		product *= r
	}
	fmt.Println("Product:", product)
}

func solutionPart2() {
	time, distance, err := loadDataPart2("input.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	pressTime := calculateRange(float64(time), float64(distance))
	fmt.Println("Range:", pressTime)
}

func main() {
	solutionPart1()
	solutionPart2()
}
