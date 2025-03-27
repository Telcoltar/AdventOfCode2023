package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	X, Y int
}

type Direction string

var UP Direction = "U"
var DOWN Direction = "D"
var LEFT Direction = "L"
var RIGHT Direction = "R"

type Cmd struct {
	Direction Direction
	Steps     int
}

func readData(filepath string) ([]Cmd, error) {
	// Read data from file
	fileContent, readErr := os.ReadFile(filepath)
	if readErr != nil {
		return nil, readErr
	}
	cmds := []Cmd{}
	for _, line := range strings.Split(string(fileContent), "\n") {
		// Process each line
		splitLine := strings.Split(line, " ")
		direction := Direction(splitLine[0])
		steps, _ := strconv.Atoi(splitLine[1])
		cmds = append(cmds, Cmd{Direction: direction, Steps: steps})
	}
	return cmds, nil
}

var directionsMap map[int]Direction = map[int]Direction{
	0: RIGHT,
	1: DOWN,
	2: LEFT,
	3: UP,
}

func readDataPart2(filepath string) ([]Cmd, error) {
	// Read data from file
	fileContent, readErr := os.ReadFile(filepath)
	if readErr != nil {
		return nil, readErr
	}
	cmds := []Cmd{}
	for _, line := range strings.Split(string(fileContent), "\n") {
		// Process each line
		splitLine := strings.Split(line, " ")
		colourPart := splitLine[2]
		steps, parseErr := strconv.ParseInt(colourPart[2:len(colourPart)-2], 16, 64)
		if parseErr != nil {
			return nil, parseErr
		}
		directionCode, parseErr := strconv.Atoi(string(colourPart[len(colourPart)-2]))
		if parseErr != nil {
			return nil, parseErr
		}
		cmds = append(cmds, Cmd{Direction: directionsMap[directionCode], Steps: int(steps)})
	}
	return cmds, nil
}

func printGrid(minY, maxY, minX, maxX int, border map[Point]bool) string {
	sb := strings.Builder{}
	for y := minY - 1; y <= maxY+1; y++ {
		for x := minX - 1; x <= maxX+1; x++ {
			if border[Point{x, y}] {
				sb.WriteString("#")
			} else {
				sb.WriteString(".")
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func mod(a, b int) int {
	return (a%b + b) % b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func solutionPart1(cmds []Cmd) int {
	// Start from (0, 0)
	currentPoint := Point{0, 0}
	corners := []Point{currentPoint}
	borderLen := 0
	// Loop through all commands
	for _, cmd := range cmds {
		switch cmd.Direction {
		case UP:
			currentPoint.Y += cmd.Steps
		case DOWN:
			currentPoint.Y -= cmd.Steps
		case LEFT:
			currentPoint.X -= cmd.Steps
		case RIGHT:
			currentPoint.X += cmd.Steps
		default:
			fmt.Printf("Invalid direction: %s\n", cmd.Direction)
		}

		borderLen += cmd.Steps
		corners = append(corners, currentPoint)
	}

	fmt.Println("Outer Corners: ", 4+(len(cmds)-4)/2)
	fmt.Println("Inner Corners: ", (len(cmds)-4)/2)
	correctionCorners := ((4+(len(cmds)-4)/2)*3 + ((len(cmds) - 4) / 2)) / 4
	fmt.Println("Correction Corners: ", correctionCorners)
	remainingCorrection := (borderLen - len(cmds)) / 2
	fmt.Println("Remaining Correction: ", remainingCorrection)
	fmt.Println("Border Length: ", borderLen)

	fmt.Println("Corners: ", len(corners))
	trapezoid := 0
	for i := 0; i < len(corners)-1; i++ {
		trapezoid += corners[i].X*corners[i+1].Y - corners[i+1].X*corners[i].Y
	}
	trapezoid = abs(trapezoid / 2)
	fmt.Println("Trapezoid: ", trapezoid)

	return trapezoid + correctionCorners + remainingCorrection
}

func main() {
	cmds, readErr := readData("input.txt")
	// Cmd length
	fmt.Println("Cmds length: ", len(cmds))
	if readErr != nil {
		fmt.Println("Error reading file: ", readErr)
		return
	}
	fmt.Println("Part 1: ", solutionPart1(cmds))
	cmds, readErr = readDataPart2("input.txt")
	// Cmd length
	fmt.Println("Cmds length: ", len(cmds))
	if readErr != nil {
		fmt.Println("Error reading file: ", readErr)
		return
	}
	fmt.Println("Part 2: ", solutionPart1(cmds))
}
