package main

import (
	"log"
	"maps"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Point struct {
	X, Y, Z int
}

func (p Point) String() string {
	return "(" + strconv.Itoa(p.X) + ", " + strconv.Itoa(p.Y) + ", " + strconv.Itoa(p.Z) + ")"
}

func (p *Point) Copy() *Point {
	return &Point{p.X, p.Y, p.Z}
}

type Brick struct {
	Start, End  *Point
	Points      []*Point
	LowestPlane []*Point
}

func (b *Brick) FallDown() {
	b.Start.Z--
	b.End.Z--
	for _, point := range b.Points {
		point.Z--
	}
	for _, point := range b.LowestPlane {
		point.Z--
	}
}

func (b *Brick) Copy() *Brick {
	copyBrick := &Brick{b.Start, b.End, nil, nil}
	copyBrick.Points = make([]*Point, len(b.Points))
	copy(copyBrick.Points, b.Points)
	copyBrick.LowestPlane = make([]*Point, len(b.LowestPlane))
	copy(copyBrick.LowestPlane, b.LowestPlane)
	return copyBrick
}

func (b Brick) String() string {
	return "Brick{" + b.Start.String() + " - " + b.End.String() + "}"
}

func calcPoints(b Brick) []*Point {
	points := []*Point{}
	for x := b.Start.X; x <= b.End.X; x++ {
		for y := b.Start.Y; y <= b.End.Y; y++ {
			for z := b.Start.Z; z <= b.End.Z; z++ {
				points = append(points, &Point{x, y, z})
			}
		}
	}
	return points
}

func calcLowestPlane(b Brick) []*Point {
	points := b.Points
	lowestPlane := []*Point{}
	for _, point := range points {
		if point.Z == b.Start.Z {
			lowestPlane = append(lowestPlane, point.Copy())
		}
	}
	return lowestPlane
}

func parseStringToPoint(str string) Point {
	str = strings.TrimSpace(str)
	parts := strings.Split(str, ",")
	x, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Fatal(err)
	}
	y, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Fatal(err)
	}
	z, err := strconv.Atoi(parts[2])
	if err != nil {
		log.Fatal(err)
	}
	return Point{x, y, z}
}

func parseLineToBrick(line string) *Brick {
	line = strings.TrimSpace(line)
	parts := strings.Split(line, "~")
	start := parseStringToPoint(parts[0])
	end := parseStringToPoint(parts[1])
	if start.Z > end.Z {
		start, end = end, start
	}
	brick := &Brick{&start, &end, nil, nil}
	brick.Points = calcPoints(*brick)
	brick.LowestPlane = calcLowestPlane(*brick)
	return brick
}

func readData(filepath string) []*Brick {
	// Read the data
	fileContent, readErr := os.ReadFile(filepath)
	if readErr != nil {
		log.Fatal(readErr)
	}
	bricks := []*Brick{}
	for _, line := range strings.Split(string(fileContent), "\n") {
		if line == "" {
			continue
		}
		bricks = append(bricks, parseLineToBrick(line))
	}
	return bricks
}

func checkForIntersection(b1 *Brick, occupiedPoints map[Point]int) bool {
	for _, point := range b1.LowestPlane {
		if _, ok := occupiedPoints[*point]; ok {
			return true
		}
	}
	return false
}

func getSupportOfBrick(b1 *Brick, occupiedPoints map[Point]int) map[int]bool {
	support := make(map[int]bool)
	for _, point := range b1.LowestPlane {
		if supportBrickIdx, ok := occupiedPoints[*point]; ok {
			support[supportBrickIdx] = true
		}
	}
	return support
}

func buildSupportMap(bricks []*Brick, occupiedPoints map[Point]int) map[int][]int {
	supportMap := make(map[int][]int)
	for idx, brick := range bricks {
		support := getSupportOfBrick(brick, occupiedPoints)
		supportMap[idx] = slices.Collect(maps.Keys(support))
	}
	return supportMap
}

func invertSupportMap(supportMap map[int][]int) map[int][]int {
	invertedSupportMap := make(map[int][]int)
	for idx, support := range supportMap {
		for _, supportIdx := range support {
			invertedSupportMap[supportIdx] = append(invertedSupportMap[supportIdx], idx)
		}
	}
	return invertedSupportMap
}

func isSubSlice(a, b []int) bool {
	for _, aVal := range a {
		if !slices.Contains(b, aVal) {
			return false
		}
	}
	return true
}

func calculateFallingBricks(supportMap map[int][]int, invertedSupportMap map[int][]int, start int) int {
	fallingsBricks := []int{start}
	currentBricks := []int{start}
	for len(currentBricks) > 0 {
		nextBricks := []int{}
		for _, brick := range currentBricks {
			nextBricks = append(nextBricks, invertedSupportMap[brick]...)
		}

		slices.Sort(nextBricks)
		nextBricks = slices.Compact(nextBricks)

		nextBricks = slices.DeleteFunc(nextBricks, func(brick int) bool {
			return !isSubSlice(supportMap[brick], fallingsBricks)
		})
		fallingsBricks = append(fallingsBricks, nextBricks...)
		currentBricks = nextBricks
	}

	slices.Sort(fallingsBricks)
	fallingsBricks = slices.Compact(fallingsBricks)
	return len(fallingsBricks) - 1
}

func common(bricks []*Brick) map[int][]int {
	// sort bricks by z coordinate
	slices.SortFunc(bricks, func(a, b *Brick) int {
		return a.Start.Z - b.Start.Z
	})
	occupied := map[Point]int{}
	// let bricks fall down as far as possible
	for idx, brick := range bricks {
		for brick.Start.Z > 1 && !checkForIntersection(brick, occupied) {
			brick.FallDown()
		}
		for _, brickPoint := range brick.Points {
			p := Point{brickPoint.X, brickPoint.Y, brickPoint.Z + 1}
			if _, ok := occupied[p]; ok {
				log.Fatal("Brick", idx, "is intersecting with another brick")
			}
			occupied[p] = idx
		}
	}

	return buildSupportMap(bricks, occupied)
}

func solutionPart1(bricks []*Brick) int {
	supportMap := common(bricks)
	bricksNeccessary := make(map[int]bool)
	for idx := range bricks {
		if len(supportMap[idx]) == 1 {
			bricksNeccessary[supportMap[idx][0]] = true
		}
	}
	log.Println(len(bricksNeccessary))
	log.Println(len(bricks))
	return len(bricks) - len(bricksNeccessary)
}

func solutionPart2(bricks []*Brick) int {
	supportMap := common(bricks)
	invertedSupportMap := invertSupportMap(supportMap)

	count := 0
	for idx := range bricks {
		count += calculateFallingBricks(supportMap, invertedSupportMap, idx)
	}
	return count
}

func main() {
	bricks := readData("input.txt")
	log.Println(solutionPart1(bricks))
	bricks = readData("example.txt")
	log.Println(solutionPart2(bricks))
}
