package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Interval struct {
	start *int64
	end   *int64
}

func (i Interval) String() string {
	if i.start != nil && i.end != nil {
		return fmt.Sprintf("(Start: %d, End: %d)", *i.start, *i.end)
	}
	if i.start != nil {
		return fmt.Sprintf("(Start: %d, End: inf)", *i.start)
	}
	if i.end != nil {
		return fmt.Sprintf("(Start: -inf, End: %d)", *i.end)
	}
	return fmt.Sprintf("Invalid Interval")
}

func NewInterval(start int64, end int64) *Interval {
	return &Interval{&start, &end}
}

type OffsetInterval struct {
	Interval
	offset int64
}

func (oi OffsetInterval) String() string {
	if oi.start != nil && oi.end != nil {
		return fmt.Sprintf("(Start: %d, End: %d, Offset: %d)", *oi.start, *oi.end, oi.offset)
	}
	if oi.start != nil {
		return fmt.Sprintf("(Start: %d, End: inf, Offset: %d)", *oi.start, oi.offset)
	}
	if oi.end != nil {
		return fmt.Sprintf("(Start: -inf, End: %d, Offset: %d)", *oi.end, oi.offset)
	}
	return fmt.Sprintf("Invalid Interval")
}

func (oi OffsetInterval) Contains(num int64) bool {
	if oi.start != nil && oi.end != nil {
		return num >= *oi.start && num <= *oi.end
	}
	if oi.start != nil {
		return num >= *oi.start
	}
	if oi.end != nil {
		return num <= *oi.end
	}
	return true
}

func (oi OffsetInterval) IntersectWithInterval(interval *Interval) *Interval {
	var start int64
	if oi.start != nil {
		start = max(*oi.start, *interval.start)
	} else {
		start = *interval.start
	}
	var end int64
	if oi.end != nil {
		end = min(*oi.end, *interval.end)
	} else {
		end = *interval.end
	}
	start += oi.offset
	end += oi.offset
	return &Interval{&start, &end}
}

func NewRightOpenOffsetInterval(start int64, offset int64) *OffsetInterval {
	return &OffsetInterval{
		Interval: Interval{
			start: &start,
			end:   nil,
		},
		offset: offset,
	}
}

func NewOffsetInterval(start int64, end int64, offset int64) *OffsetInterval {
	return &OffsetInterval{
		Interval: Interval{
			start: &start,
			end:   &end,
		},
		offset: offset,
	}
}

func fillIntervalsGaps(intervals []*OffsetInterval) []*OffsetInterval {
	intervalsToAdd := make([]*OffsetInterval, 0)
	sort.Slice(intervals, func(i, j int) bool {
		return *intervals[i].start < *intervals[j].start
	})
	if *intervals[0].start != 0 {
		intervalsToAdd = append(intervalsToAdd,
			NewOffsetInterval(0, *intervals[0].start-1, 0),
		)
	}
	for i := 0; i < len(intervals)-1; i++ {
		if *intervals[i].end+1 != *intervals[i+1].start {
			intervalsToAdd = append(intervalsToAdd,
				NewOffsetInterval(*intervals[i].end+1, *intervals[i+1].start-1, 0),
			)
		}
	}
	intervalsToAdd = append(intervalsToAdd,
		NewRightOpenOffsetInterval(*intervals[len(intervals)-1].end+1, 0),
	)
	intervals = append(intervals, intervalsToAdd...)
	sort.Slice(intervals, func(i, j int) bool {
		return *intervals[i].start < *intervals[j].start
	})
	return intervals
}

func loadData(filename string) ([]int64, [][]*OffsetInterval, error) {
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
	intervals := make([][]*OffsetInterval, 0)
	scanner.Scan()
	seedString := strings.TrimSpace(strings.Split(scanner.Text(), ":")[1])
	seeds := make([]int64, 0)
	for _, seed := range strings.Split(seedString, " ") {
		parsedSeed, err := strconv.ParseInt(seed, 10, 64)
		if err != nil {
			return nil, nil, err
		}
		seeds = append(seeds, parsedSeed)
	}
	scanner.Scan()
	for scanner.Scan() {
		scanner.Scan()
		line := scanner.Text()
		currentMap := make([]*OffsetInterval, 0)
		for line != "" {
			numbers := make([]int64, 0)
			for _, num := range strings.Split(line, " ") {
				parsedNum, err := strconv.ParseInt(num, 10, 64)
				if err != nil {
					return nil, nil, err
				}
				numbers = append(numbers, parsedNum)
			}
			currentMap = append(currentMap,
				NewOffsetInterval(
					numbers[1],
					numbers[1]+numbers[2]-1,
					numbers[0]-numbers[1],
				),
			)
			scanner.Scan()
			line = scanner.Text()
		}
		currentMap = fillIntervalsGaps(currentMap)
		intervals = append(intervals, currentMap)
	}
	return seeds, intervals, nil
}

func solutionPart1() {
	seeds, maps, err := loadData("input.txt")
	if err != nil {
		fmt.Println("Error loading data:", err)
		return
	}
	sort.Slice(seeds, func(i, j int) bool {
		return seeds[i] < seeds[j]
	})
	currentSeeds := seeds
	for _, seedMap := range maps {
		currentMapIndex := 0
		nextSeeds := make([]int64, 0)
		for _, seed := range currentSeeds {
			for !seedMap[currentMapIndex].Contains(seed) {
				currentMapIndex += 1
			}
			nextSeeds = append(nextSeeds, seed+seedMap[currentMapIndex].offset)
		}
		sort.Slice(nextSeeds, func(i, j int) bool {
			return nextSeeds[i] < nextSeeds[j]
		})
		currentSeeds = nextSeeds
	}
	println(currentSeeds[0])
}

func solutionPart2() {
	seeds, maps, err := loadData("input.txt")
	if err != nil {
		fmt.Println("Error loading data:", err)
		return
	}
	currentSeedIntervals := make([]*Interval, 0)
	for i := 0; i < len(seeds)-1; i += 2 {
		currentSeedIntervals = append(currentSeedIntervals, NewInterval(seeds[i], seeds[i]+seeds[i+1]-1))
	}
	sort.Slice(currentSeedIntervals, func(i, j int) bool {
		return *currentSeedIntervals[i].start < *currentSeedIntervals[j].start
	})
	for _, seedMap := range maps {
		nextSeedIntervals := make([]*Interval, 0)
		currentSeedMapIndex := 0
		for _, interval := range currentSeedIntervals {
			for !seedMap[currentSeedMapIndex].Contains(*interval.start) {
				currentSeedMapIndex += 1
			}
			start := currentSeedMapIndex
			for !seedMap[currentSeedMapIndex].Contains(*interval.end) {
				currentSeedMapIndex += 1
			}
			for i := start; i <= currentSeedMapIndex; i++ {
				nextSeedIntervals = append(
					nextSeedIntervals,
					seedMap[i].IntersectWithInterval(interval),
				)
			}
		}
		currentSeedIntervals = nextSeedIntervals
		sort.Slice(currentSeedIntervals, func(i, j int) bool {
			return *currentSeedIntervals[i].start < *currentSeedIntervals[j].start
		})
	}
	println(*currentSeedIntervals[0].start)
}

func main() {
	solutionPart1()
	solutionPart2()
}
