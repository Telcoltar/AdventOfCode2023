package main

import (
	"bytes"
	"log"
	"maps"
	"os"
	"reflect"
	"slices"
	"strings"
)

type ModuleType string

const (
	broadcast   ModuleType = "broadcaster"
	flipflop    ModuleType = "%"
	conjunction ModuleType = "&"
)

type ModuleDescription struct {
	Name    string
	Type    ModuleType
	Outputs []string
}

type Module interface {
	Process(input string, high bool) int
	GetOutputs() []string
}

type Broadcaster struct {
	Outputs []string
}

func (b *Broadcaster) Process(input string, high bool) int {
	// Convert the boolean to int (0 for false, 1 for true)
	if high {
		return 1
	}
	return 0
}

func (b *Broadcaster) GetOutputs() []string {
	return b.Outputs
}

type FlipFlop struct {
	Outputs []string
	Status  bool
}

func (f *FlipFlop) Process(input string, high bool) int {
	// If the input is low, set the status to low
	if !high {
		f.Status = !f.Status
		if f.Status {
			return 1
		}
		return 0
	}
	return -1
}

func (f *FlipFlop) GetOutputs() []string {
	return f.Outputs
}

type Conjunction struct {
	Status  map[string]bool
	Outputs []string
}

func (c *Conjunction) Process(input string, high bool) int {
	// Set the status of the input
	c.Status[input] = high
	// Check if all inputs are high
	for _, v := range c.Status {
		if !v {
			return 1
		}
	}
	return 0
}

func (c *Conjunction) GetOutputs() []string {
	return c.Outputs
}

func readData(filename string) map[string]Module {
	// Read the data
	fileContent, readErr := os.ReadFile(filename)
	if readErr != nil {
		log.Fatal(readErr)
	}
	modules := map[string]Module{}
	conjunctions := []string{}
	outputInputMap := map[string][]string{}
	for _, line := range strings.Split(string(fileContent), "\n") {
		if line == "" {
			continue
		}
		// Split the line
		parts := strings.Split(line, "->")
		// Check the type of the module
		rawOutputs := strings.Split(parts[1], ",")

		// Trim spaces for each output
		outputs := make([]string, len(rawOutputs))
		for i, output := range rawOutputs {
			outputs[i] = strings.TrimSpace(output)
		}

		moduleName := strings.TrimSpace(parts[0])
		if moduleName == "broadcaster" {
			modules[moduleName] = &Broadcaster{Outputs: outputs}
			for _, output := range outputs {
				outputInputMap[output] = append(outputInputMap[output], moduleName)
			}
		} else if strings.Contains(moduleName, "%") {
			modules[moduleName[1:]] = &FlipFlop{Outputs: outputs}
			for _, output := range outputs {
				outputInputMap[output] = append(outputInputMap[output], moduleName[1:])
			}
		} else if strings.Contains(moduleName, "&") {
			conjunctions = append(conjunctions, moduleName[1:])
			for _, output := range outputs {
				outputInputMap[output] = append(outputInputMap[output], moduleName[1:])
			}
			modules[moduleName[1:]] = &Conjunction{Status: map[string]bool{}, Outputs: outputs}
		}
	}

	// Set the outputs of the conjunctions
	for _, conjunction := range conjunctions {
		conjunctionModule := modules[conjunction].(*Conjunction)
		for _, input := range outputInputMap[conjunction] {
			conjunctionModule.Status[input] = false
		}
	}

	return modules
}

type Pulse struct {
	high     bool
	Receiver string
	Sender   string
}

func doCycle(modules map[string]Module, cycle int, mem map[string][]int) (int, int) {
	start := Pulse{high: false, Receiver: "broadcaster", Sender: "button"}
	queue := []Pulse{start}
	lowPulses := 0
	highPulses := 0
	seen := map[string]bool{}
	for len(queue) > 0 {
		// Get the first element of the queue
		pulse := queue[0]
		queue = queue[1:]
		if pulse.high {
			highPulses++
		} else {
			lowPulses++
		}
		// Get the module
		module := modules[pulse.Receiver]
		// Process the module
		if module == nil {
			continue
		}
		processOutput := module.Process(pulse.Sender, pulse.high)
		if pulse.Receiver == "hp" && (pulse.high || seen[pulse.Sender]) {
			seen[pulse.Sender] = true
			log.Println("Cycle:", cycle, "Sender:", pulse.Sender, "High:", pulse.high)
			mem[pulse.Sender] = append(mem[pulse.Sender], cycle)
		}
		// If the output is -1, continue
		if processOutput == -1 {
			continue
		}

		for _, output := range module.GetOutputs() {
			queue = append(queue, Pulse{high: processOutput == 1, Receiver: output, Sender: pulse.Receiver})
		}
	}
	return highPulses, lowPulses
}

func boolsToBytes(bools []bool) []byte {
	bytes := make([]byte, len(bools))
	for i, b := range bools {
		if b {
			bytes[i] = 1
		}
	}
	return bytes
}

func getFlipFlopStatuses(modules map[string]Module) []byte {
	statuses := map[string]bool{}
	for name, module := range modules {
		if flipFlop, ok := module.(*FlipFlop); ok {
			statuses[name] = flipFlop.Status
		}
	}
	statusSlice := make([]bool, len(statuses))
	keys := slices.Collect(maps.Keys(statuses))
	slices.Sort(keys)
	for i, key := range keys {
		statusSlice[i] = statuses[key]
	}
	// Convert the statuses to bytes
	return boolsToBytes(statusSlice)
}

func sliceContains(slices [][]byte, slice []byte) bool {
	for _, s := range slices {
		if bytes.Equal(s, slice) {
			return true
		}
	}
	return false
}

/*
	func solutionPart1(modules map[string]Module) int {
		highPulses := 0
		lowPulses := 0
		cycle := 0
		for {
			high, low := doCycle(modules, cycle)
			highPulses += high
			lowPulses += low
			cycle++
			if cycle == 1000 {
				break
			}
		}
		pulses := highPulses * lowPulses
		log.Println("High pulses:", highPulses)
		log.Println("Low pulses:", lowPulses)
		log.Printf("Cycle: %d", cycle)
		return pulses
	}
*/
func getPointerValue(i interface{}) uintptr {
	return reflect.ValueOf(i).Pointer()
}

func solutionPart2(modules map[string]Module) int {
	highPulses := 0
	lowPulses := 0
	cycle := 1
	mem := map[string][]int{}
	for {
		high, low := doCycle(modules, cycle, mem)
		highPulses += high
		lowPulses += low
		cycle++
		if cycle == 1000000 {
			break
		}
	}
	diffMap := map[string]int{}
	for k, v := range mem {
		diffs := []int{}
		for i := 1; i < len(v); i++ {
			diffs = append(diffs, v[i]-v[i-1])
		}
		diffMap[k] = diffs[0]
	}
	mult := 1
	for _, v := range diffMap {
		mult *= v
	}
	return mult
}

func main() {
	modules := readData("input.txt")
	//log.Printf("Solution part 1: %d", solutionPart1(modules))
	log.Printf("Solution part 2: %d", solutionPart2(modules))
}
