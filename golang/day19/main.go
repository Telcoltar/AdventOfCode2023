package main

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/tiendc/go-deepcopy"
)

type Part struct {
	X, M, A, S int
}

func (p *Part) AddUp() int {
	return p.X + p.M + p.A + p.S
}

type Operation rune

const gt Operation = '>'
const lt Operation = '<'

type AltRule struct {
	Field        string
	Operation    Operation
	Num          int
	NextWorkflow string
}

type Rule struct {
	Condition    func(*Part) bool
	NextWorkflow string
}

type Workflow struct {
	Rule       []Rule
	LastAction string
}

type AltWorkflow struct {
	Rule       []AltRule
	LastAction string
}

func (w *Workflow) nextAction(part *Part) string {
	for _, rule := range w.Rule {
		if rule.Condition(part) {
			return rule.NextWorkflow
		}
	}
	return w.LastAction
}

var ruleRegex = regexp.MustCompile(`([axsm])([<>])(\d+)`)

func parseAltRule(rule string) AltRule {
	split := strings.Split(rule, ":")
	nextWorkflow := split[1]
	firstPart := ruleRegex.FindStringSubmatch(split[0])
	num, parseErr := strconv.Atoi(firstPart[3])
	if parseErr != nil {
		log.Fatal(parseErr)
	}
	var operation Operation
	if firstPart[2] == ">" {
		operation = gt
	} else {
		operation = lt
	}
	return AltRule{
		Field:        firstPart[1],
		Operation:    operation,
		Num:          num,
		NextWorkflow: nextWorkflow,
	}
}

func parseRule(rule string) Rule {
	var condition func(*Part) bool
	split := strings.Split(rule, ":")
	nextWorkflow := split[1]
	firstPart := ruleRegex.FindStringSubmatch(split[0])
	num, parseErr := strconv.Atoi(firstPart[3])
	if parseErr != nil {
		log.Fatal(parseErr)
	}
	if firstPart[2] == ">" {
		switch firstPart[1] {
		case "a":
			condition = func(part *Part) bool {
				return part.A > num
			}
		case "x":
			condition = func(part *Part) bool {
				return part.X > num
			}
		case "s":
			condition = func(part *Part) bool {
				return part.S > num
			}
		case "m":
			condition = func(part *Part) bool {
				return part.M > num
			}
		default:
			log.Fatal("Invalid condition")
		}
	} else {
		switch firstPart[1] {
		case "a":
			condition = func(part *Part) bool {
				return part.A < num
			}
		case "x":
			condition = func(part *Part) bool {
				return part.X < num
			}
		case "s":
			condition = func(part *Part) bool {
				return part.S < num
			}
		case "m":
			condition = func(part *Part) bool {
				return part.M < num
			}
		default:
			log.Fatal("Invalid condition")
		}
	}

	return Rule{
		Condition:    condition,
		NextWorkflow: nextWorkflow,
	}
}

func readData(filepath string) (map[string]*Workflow, []*Part) {
	fileContent, readErr := os.ReadFile(filepath)
	if readErr != nil {
		log.Fatal(readErr)
	}
	splitContent := strings.Split(string(fileContent), "\n\n")
	workflowRegex := regexp.MustCompile(`^([a-z]{2,3})\{(\S+),([a-zRA]+)\}$`)
	workflows := make(map[string]*Workflow)
	for _, workflow := range strings.Split(splitContent[0], "\n") {
		matches := workflowRegex.FindStringSubmatch(workflow)
		workflowName := matches[1]
		defaultAction := matches[3]
		rules := make([]Rule, 0)
		for _, rule := range strings.Split(matches[2], ",") {
			rules = append(rules, parseRule(rule))
		}
		workflows[workflowName] = &Workflow{
			Rule:       rules,
			LastAction: defaultAction,
		}
	}

	parts := make([]*Part, 0)
	for _, part := range strings.Split(splitContent[1], "\n") {
		partRegex := regexp.MustCompile(`([axsm])=(\d+)`)
		matches := partRegex.FindAllStringSubmatch(part, -1)
		x, xErr := strconv.Atoi(matches[0][2])
		m, mErr := strconv.Atoi(matches[1][2])
		a, aErr := strconv.Atoi(matches[2][2])
		s, sErr := strconv.Atoi(matches[3][2])
		if xErr != nil || mErr != nil || aErr != nil || sErr != nil {
			log.Fatal("Invalid part")
		}
		parts = append(parts, &Part{
			X: x,
			M: m,
			A: a,
			S: s,
		})
	}
	return workflows, parts
}

func readPart2(filepath string) map[string]*AltWorkflow {
	fileContent, readErr := os.ReadFile(filepath)
	if readErr != nil {
		log.Fatal(readErr)
	}
	splitContent := strings.Split(string(fileContent), "\n\n")
	workflowRegex := regexp.MustCompile(`^([a-z]{2,3})\{(\S+),([a-zRA]+)\}$`)
	workflows := make(map[string]*AltWorkflow)
	for _, workflow := range strings.Split(splitContent[0], "\n") {
		matches := workflowRegex.FindStringSubmatch(workflow)
		workflowName := matches[1]
		defaultAction := matches[3]
		rules := make([]AltRule, 0)
		for _, rule := range strings.Split(matches[2], ",") {
			rules = append(rules, parseAltRule(rule))
		}
		workflows[workflowName] = &AltWorkflow{
			Rule:       rules,
			LastAction: defaultAction,
		}
	}
	return workflows
}

func solutionPart1(workflows map[string]*Workflow, parts []*Part) int {
	accepted := make([]*Part, 0)
	rejected := make([]*Part, 0)
	for _, part := range parts {
		nextAction := "in"
		for nextAction != "A" && nextAction != "R" {
			nextAction = workflows[nextAction].nextAction(part)
		}
		if nextAction == "A" {
			accepted = append(accepted, part)
		} else {
			rejected = append(rejected, part)
		}
	}

	log.Printf("Accepted: %d, Rejected: %d\n", len(accepted), len(rejected))

	sumAccepted := 0
	for _, part := range accepted {
		sumAccepted += part.AddUp()
	}

	return sumAccepted
}

type Status struct {
	X, M, A, S   *Intervall
	NextWorkflow string
}

func (s *Status) GetIntervall(field string) *Intervall {
	switch field {
	case "x":
		return s.X
	case "m":
		return s.M
	case "a":
		return s.A
	case "s":
		return s.S
	}
	return nil
}

func (s *Status) SetIntervall(field string, intervall *Intervall) {
	switch field {
	case "x":
		s.X = intervall
	case "m":
		s.M = intervall
	case "a":
		s.A = intervall
	case "s":
		s.S = intervall
	}
}

func (s *Status) AnyEmpty() bool {
	return s.X.Size() <= 0 || s.M.Size() <= 0 || s.A.Size() <= 0 || s.S.Size() <= 0
}

func (s *Status) Size() int {
	return s.X.Size() * s.M.Size() * s.A.Size() * s.S.Size()
}

type Intervall struct {
	Min, Max int
}

func (i *Intervall) contains(num int) bool {
	return i.Min <= num && i.Max >= num
}

func (i *Intervall) String() string {
	return strconv.Itoa(i.Min) + "-" + strconv.Itoa(i.Max)
}

func (i *Intervall) Size() int {
	return i.Max - i.Min + 1
}

func (i *Intervall) Intersect(other *Intervall) *Intervall {
	if i.Max < other.Min || i.Min > other.Max {
		return nil
	}
	return &Intervall{
		Min: max(i.Min, other.Min),
		Max: min(i.Max, other.Max),
	}
}

func (st *Status) Intersect(other *Status) *Status {
	x := st.X.Intersect(other.X)
	m := st.M.Intersect(other.M)
	a := st.A.Intersect(other.A)
	s := st.S.Intersect(other.S)
	if x == nil || m == nil || a == nil || s == nil {
		return nil
	}
	return &Status{
		X:            x,
		M:            m,
		A:            a,
		S:            s,
		NextWorkflow: st.NextWorkflow,
	}
}

func ApplyRule(status *Status, rule AltRule) (*Status, *Status) {
	intervall := status.GetIntervall(rule.Field)
	if intervall == nil {
		log.Fatal("Invalid field")
	}
	if intervall.contains(rule.Num) {
		higherInterval := &Intervall{}
		if err := deepcopy.Copy(higherInterval, intervall); err != nil {
			log.Fatal(err)
		}

		if rule.Operation == gt {
			intervall.Max = rule.Num
			higherInterval.Min = rule.Num + 1
			lowerStatus := &Status{}
			higherStatus := status
			deepcopy.Copy(lowerStatus, higherStatus)
			lowerStatus.SetIntervall(rule.Field, intervall)
			higherStatus.SetIntervall(rule.Field, higherInterval)
			higherStatus.NextWorkflow = rule.NextWorkflow
			return lowerStatus, higherStatus
		} else {
			intervall.Max = rule.Num - 1
			higherInterval.Min = rule.Num
			lowerStatus := &Status{}
			higherStatus := status
			deepcopy.Copy(lowerStatus, higherStatus)
			lowerStatus.SetIntervall(rule.Field, intervall)
			higherStatus.SetIntervall(rule.Field, higherInterval)
			lowerStatus.NextWorkflow = rule.NextWorkflow
			return higherStatus, lowerStatus
		}
	}
	return status, nil
}

func NextIntervall(status []Status, workflows map[string]*AltWorkflow) []Status {
	nextStatus := make([]Status, 0)
	for _, oneStatus := range status {
		workflow := workflows[oneStatus.NextWorkflow]
		currentStatus := &oneStatus
		for _, rule := range workflow.Rule {
			remain, updated := ApplyRule(currentStatus, rule)
			if updated != nil && !updated.AnyEmpty() {
				nextStatus = append(nextStatus, *updated)
			}
			if remain.AnyEmpty() {
				currentStatus = nil
				break
			}
			currentStatus = remain
		}
		if currentStatus != nil && !currentStatus.AnyEmpty() {
			currentStatus.NextWorkflow = workflow.LastAction
			nextStatus = append(nextStatus, *currentStatus)
		}
	}
	return nextStatus
}

func calculateIntervalsForField(workflows map[string]*AltWorkflow) []Status {
	s := []Status{
		{
			X:            &Intervall{Min: 1, Max: 4000},
			M:            &Intervall{Min: 1, Max: 4000},
			A:            &Intervall{Min: 1, Max: 4000},
			S:            &Intervall{Min: 1, Max: 4000},
			NextWorkflow: "in",
		},
	}
	finishedStatus := make([]Status, 0)
	currentStatus := s
	for len(currentStatus) > 0 {
		currentStatus = NextIntervall(currentStatus, workflows)
		nextStatus := make([]Status, 0)
		for _, status := range currentStatus {
			if !status.AnyEmpty() {
				if status.NextWorkflow == "A" {
					finishedStatus = append(finishedStatus, status)
				} else if status.NextWorkflow != "R" {
					nextStatus = append(nextStatus, status)
				}
			}
		}
		currentStatus = nextStatus
	}
	return finishedStatus
}

func solutionPart2(workflows map[string]*AltWorkflow) int {
	finidshedIntevals := calculateIntervalsForField(workflows)
	for _, status := range finidshedIntevals {
		for _, status2 := range finidshedIntevals {
			if status != status2 {
				intersect := status.Intersect(&status2)
				if intersect != nil && !intersect.AnyEmpty() {
					log.Println("Intersect", intersect)
				}
			}
		}
	}
	sum := 0
	for _, status := range finidshedIntevals {
		sum += status.Size()
	}
	return sum
}

func main() {
	workflows, parts := readData("input.txt")
	log.Println(solutionPart1(workflows, parts))
	log.Println(solutionPart2(readPart2("input.txt")))
}
