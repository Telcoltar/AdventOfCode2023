package main

import (
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/mat"
)

type Vec2D struct {
	X, Y int
}

type Vec3D struct {
	X, Y, Z int
}

func ParseStringToVec3D(s string) Vec3D {
	vec := Vec3D{}
	split := strings.Split(strings.TrimSpace(s), ",")
	x, err := strconv.Atoi(strings.TrimSpace(split[0]))
	if err != nil {
		log.Fatal(err)
	}
	vec.X = x
	y, err := strconv.Atoi(strings.TrimSpace(split[1]))
	if err != nil {
		log.Fatal(err)
	}
	vec.Y = y
	z, err := strconv.Atoi(strings.TrimSpace(split[2]))
	if err != nil {
		log.Fatal(err)
	}
	vec.Z = z
	return vec
}

type Line3D struct {
	FixPoint  Vec3D
	Direction Vec3D
}

type Line2D struct {
	FixPoint  Vec2D
	Direction Vec2D
}

func Line2DFrom3D(line *Line3D) *Line2D {
	return &Line2D{
		Vec2D{line.FixPoint.X, line.FixPoint.Y},
		Vec2D{line.Direction.X, line.Direction.Y},
	}
}

func ParseStringToLine(s string) *Line3D {
	line := Line3D{}
	split := strings.Split(s, "@")
	line.FixPoint = ParseStringToVec3D(split[0])
	line.Direction = ParseStringToVec3D(split[1])
	return &line
}

func readData(filepath string) []*Line3D {
	lines := []*Line3D{}
	data, readErr := os.ReadFile(filepath)
	if readErr != nil {
		log.Fatal(readErr)
	}
	for _, line := range strings.Split(string(data), "\n") {
		if line == "" {
			continue
		}
		lines = append(lines, ParseStringToLine(line))
	}
	return lines
}

func hasIntersection(line1, line2 *Line2D) bool {
	return line2.Direction.Y*line1.Direction.X != line1.Direction.Y*line2.Direction.X
}

func intersectionPoint(line1, line2 *Line2D) (float64, float64, float64, float64) {
	det := float64(line2.Direction.Y*line1.Direction.X - line1.Direction.Y*line2.Direction.X)
	s := float64((line2.FixPoint.X-line1.FixPoint.X)*line2.Direction.Y - (line2.FixPoint.Y-line1.FixPoint.Y)*line2.Direction.X)
	t := float64((line2.FixPoint.X-line1.FixPoint.X)*line1.Direction.Y - (line2.FixPoint.Y-line1.FixPoint.Y)*line1.Direction.X)
	return float64(line1.FixPoint.X) + s/det*float64(line1.Direction.X), float64(line1.FixPoint.Y) + s/det*float64(line1.Direction.Y), s / det, t / det
}

func solutionPart1(lines []*Line3D, minValue float64, maxValue float64) int {
	count := 0
	for i, line1 := range lines {
		for j := i + 1; j < len(lines); j++ {
			line2 := lines[j]
			if hasIntersection(Line2DFrom3D(line1), Line2DFrom3D(line2)) {
				x, y, s, t := intersectionPoint(Line2DFrom3D(line1), Line2DFrom3D(line2))
				if s >= 0 && t >= 0 && x >= minValue && x <= maxValue && y >= minValue && y <= maxValue {
					count++
				}
			}
		}
	}
	return count
}

func solutionPart2(lines []*Line3D) int {
	indices := [2][2]int{{0, 1}, {0, 2}}
	data := []float64{}
	rhs := []float64{}
	for _, indices := range indices {
		// row 0
		data = append(data, float64(lines[indices[1]].FixPoint.Y-lines[indices[0]].FixPoint.Y))
		data = append(data, float64(lines[indices[0]].FixPoint.X-lines[indices[1]].FixPoint.X))
		data = append(data, 0)
		data = append(data, float64(lines[indices[0]].Direction.Y-lines[indices[1]].Direction.Y))
		data = append(data, float64(lines[indices[1]].Direction.X-lines[indices[0]].Direction.X))
		data = append(data, 0)

		rhs = append(rhs, float64(
			lines[indices[1]].FixPoint.X*lines[indices[1]].Direction.Y-lines[indices[1]].FixPoint.Y*lines[indices[1]].Direction.X-
				lines[indices[0]].FixPoint.X*lines[indices[0]].Direction.Y+lines[indices[0]].FixPoint.Y*lines[indices[0]].Direction.X))
		// row 1
		data = append(data, float64(lines[indices[1]].FixPoint.Z-lines[indices[0]].FixPoint.Z))
		data = append(data, 0)
		data = append(data, float64(lines[indices[0]].FixPoint.X-lines[indices[1]].FixPoint.X))
		data = append(data, float64(lines[indices[0]].Direction.Z-lines[indices[1]].Direction.Z))
		data = append(data, 0)
		data = append(data, float64(lines[indices[1]].Direction.X-lines[indices[0]].Direction.X))
		rhs = append(rhs, float64(
			lines[indices[1]].FixPoint.X*lines[indices[1]].Direction.Z-lines[indices[1]].FixPoint.Z*lines[indices[1]].Direction.X-
				lines[indices[0]].FixPoint.X*lines[indices[0]].Direction.Z+lines[indices[0]].FixPoint.Z*lines[indices[0]].Direction.X))
		// row 2
		data = append(data, 0)
		data = append(data, float64(lines[indices[1]].FixPoint.Z-lines[indices[0]].FixPoint.Z))
		data = append(data, float64(lines[indices[0]].FixPoint.Y-lines[indices[1]].FixPoint.Y))
		data = append(data, 0)
		data = append(data, float64(lines[indices[0]].Direction.Z-lines[indices[1]].Direction.Z))
		data = append(data, float64(lines[indices[1]].Direction.Y-lines[indices[0]].Direction.Y))
		rhs = append(rhs, float64(
			lines[indices[1]].FixPoint.Y*lines[indices[1]].Direction.Z-lines[indices[1]].FixPoint.Z*lines[indices[1]].Direction.Y-
				lines[indices[0]].FixPoint.Y*lines[indices[0]].Direction.Z+lines[indices[0]].FixPoint.Z*lines[indices[0]].Direction.Y))
	}
	matrix := mat.NewDense(6, 6, data)
	rhsVec := mat.NewVecDense(6, rhs)
	res := mat.NewVecDense(6, nil)
	err := res.SolveVec(matrix, rhsVec)
	if err != nil {
		log.Fatal(err)
	}
	sum := 0
	for i := 0; i < 3; i++ {
		sum += int(math.Round(res.At(i+3, 0)))
	}
	if sum < 0 {
		sum = -sum
	}
	return sum
}

func main() {
	lines := readData("example.txt")
	log.Println("Solution Part 1:", solutionPart1(lines, 7, 27))
	log.Println("Solution Part 2:", solutionPart2(lines))
	lines = readData("input.txt")
	log.Println("Solution Part 1:", solutionPart1(lines, 200000000000000, 400000000000000))
	log.Println("Solution Part 2:", solutionPart2(lines))
}
