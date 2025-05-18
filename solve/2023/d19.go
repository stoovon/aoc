package solve2023

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day19 struct {
}

func (d Day19) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2023, Day: 19}
}

type Rating int

const (
	X Rating = iota
	M
	A
	S
)

func parseRating(r byte) (Rating, error) {
	switch r {
	case 'x':
		return X, nil
	case 'm':
		return M, nil
	case 'a':
		return A, nil
	case 's':
		return S, nil
	default:
		return 0, fmt.Errorf("invalid rating: %c", r)
	}
}

type Part struct {
	X, M, A, S int64
}

func (p Part) Get(r Rating) int64 {
	switch r {
	case X:
		return p.X
	case M:
		return p.M
	case A:
		return p.A
	case S:
		return p.S
	default:
		panic("invalid rating")
	}
}

func (p Part) Total() int64 {
	return p.X + p.M + p.A + p.S
}

var partRE = regexp.MustCompile(`\{x=(\d+),m=(\d+),a=(\d+),s=(\d+)}`)

func parsePart(line string) (Part, error) {
	m := partRE.FindStringSubmatch(line)
	if m == nil {
		return Part{}, errors.New("invalid part")
	}
	x, _ := strconv.ParseInt(m[1], 10, 64)
	mm, _ := strconv.ParseInt(m[2], 10, 64)
	a, _ := strconv.ParseInt(m[3], 10, 64)
	s, _ := strconv.ParseInt(m[4], 10, 64)
	return Part{X: x, M: mm, A: a, S: s}, nil
}

type WorkflowName string

const ACCEPTED WorkflowName = "A  "

type Comparison int

const (
	Greater Comparison = iota
	Less
)

type WorkflowStep struct {
	Rating     Rating
	Comparison Comparison
	Comparator int64
	Target     WorkflowName
}

func parseWorkflowStep(s string) (*WorkflowStep, error) {
	colon := strings.IndexByte(s, ':')
	if colon == -1 {
		return nil, errors.New("no colon in step")
	}
	r, err := parseRating(s[0])
	if err != nil {
		return nil, err
	}
	var cmp Comparison
	switch s[1] {
	case '>':
		cmp = Greater
	case '<':
		cmp = Less
	default:
		return nil, errors.New("invalid comparison")
	}
	val, err := strconv.ParseInt(s[2:colon], 10, 64)
	if err != nil {
		return nil, err
	}
	target := parseWorkflowName(s[colon+1:])
	return &WorkflowStep{
		Rating:     r,
		Comparison: cmp,
		Comparator: val,
		Target:     target,
	}, nil
}

type Workflow struct {
	Name    WorkflowName
	Steps   []WorkflowStep
	Default WorkflowName
}

func parseWorkflow(line string) (*Workflow, error) {
	brace := strings.IndexByte(line, '{')
	if brace == -1 {
		return nil, errors.New("no { in workflow")
	}
	name := parseWorkflowName(line[:brace])
	stepsStr := line[brace+1 : len(line)-1]
	parts := strings.Split(stepsStr, ",")
	var steps []WorkflowStep
	var defaultTarget WorkflowName
	for _, part := range parts {
		if step, err := parseWorkflowStep(part); err == nil {
			steps = append(steps, *step)
		} else {
			defaultTarget = parseWorkflowName(part)
		}
	}
	return &Workflow{
		Name:    name,
		Steps:   steps,
		Default: defaultTarget,
	}, nil
}

func (w *Workflow) Process(p Part) WorkflowName {
	for _, step := range w.Steps {
		val := p.Get(step.Rating)
		switch step.Comparison {
		case Greater:
			if val > step.Comparator {
				return step.Target
			}
		case Less:
			if val < step.Comparator {
				return step.Target
			}
		}
	}
	return w.Default
}

func (w *Workflow) ProcessTesseract(space StateSpace, queue *[]PossibilityState) {
	next := &space
	for _, step := range w.Steps {
		if next == nil {
			break
		}
		split, retain := next.Split(step)
		if split != nil {
			*queue = append(*queue, PossibilityState{
				StateSpace: *split,
				Workflow:   step.Target,
			})
		}
		next = retain
	}
	if next != nil {
		*queue = append(*queue, PossibilityState{
			StateSpace: *next,
			Workflow:   w.Default,
		})
	}
}

type Range struct{ Lo, Hi int64 }

func (r Range) Size() int64 {
	return 1 + r.Hi - r.Lo
}

func (r Range) Split(cmp Comparison, comp int64) (split, retain *Range) {
	switch cmp {
	case Greater:
		if comp >= r.Hi {
			retain = &r
		} else if comp < r.Lo {
			split = &r
		} else {
			retain = &Range{r.Lo, comp}
			split = &Range{comp + 1, r.Hi}
		}
	case Less:
		if comp > r.Hi {
			split = &r
		} else if comp <= r.Lo {
			retain = &r
		} else {
			retain = &Range{comp, r.Hi}
			split = &Range{r.Lo, comp - 1}
		}
	}
	return
}

type StateSpace struct {
	X, M, A, S Range
}

func initialStateSpace() StateSpace {
	return StateSpace{
		X: Range{1, 4000},
		M: Range{1, 4000},
		A: Range{1, 4000},
		S: Range{1, 4000},
	}
}

func (s StateSpace) Get(r Rating) Range {
	switch r {
	case X:
		return s.X
	case M:
		return s.M
	case A:
		return s.A
	case S:
		return s.S
	default:
		panic("invalid rating")
	}
}

func (s StateSpace) SplitSpace(r Rating, rg Range) StateSpace {
	switch r {
	case X:
		return StateSpace{X: rg, M: s.M, A: s.A, S: s.S}
	case M:
		return StateSpace{X: s.X, M: rg, A: s.A, S: s.S}
	case A:
		return StateSpace{X: s.X, M: s.M, A: rg, S: s.S}
	case S:
		return StateSpace{X: s.X, M: s.M, A: s.A, S: rg}
	default:
		panic("invalid rating")
	}
}

func (s StateSpace) Split(step WorkflowStep) (split, retain *StateSpace) {
	rg := s.Get(step.Rating)
	splitRg, retainRg := rg.Split(step.Comparison, step.Comparator)
	if splitRg != nil {
		ss := s.SplitSpace(step.Rating, *splitRg)
		split = &ss
	}
	if retainRg != nil {
		ss := s.SplitSpace(step.Rating, *retainRg)
		retain = &ss
	}
	return
}

func (s StateSpace) Volume() int64 {
	return s.X.Size() * s.M.Size() * s.A.Size() * s.S.Size()
}

type PossibilityState struct {
	StateSpace StateSpace
	Workflow   WorkflowName
}

type WorkflowSystem struct {
	Workflows map[WorkflowName]*Workflow
	Parts     []Part
}

func parseWorkflowName(s string) WorkflowName {
	// Pad to 3 chars for compatibility with Rust's tuple struct
	if len(s) == 1 {
		return WorkflowName(s + "  ")
	}
	if len(s) == 2 {
		return WorkflowName(s + " ")
	}
	if len(s) >= 3 {
		return WorkflowName(s[:3])
	}
	return WorkflowName("   ")
}

func parseWorkflowSystem(input string) (*WorkflowSystem, error) {
	sections := strings.SplitN(input, "\n\n", 2)
	if len(sections) != 2 {
		return nil, errors.New("invalid input")
	}
	wfLines := strings.Split(strings.TrimSpace(sections[0]), "\n")
	partLines := strings.Split(strings.TrimSpace(sections[1]), "\n")
	workflows := make(map[WorkflowName]*Workflow)
	for _, line := range wfLines {
		wf, err := parseWorkflow(line)
		if err != nil {
			return nil, err
		}
		workflows[wf.Name] = wf
	}
	var parts []Part
	for _, line := range partLines {
		p, err := parsePart(line)
		if err != nil {
			return nil, err
		}
		parts = append(parts, p)
	}
	return &WorkflowSystem{Workflows: workflows, Parts: parts}, nil
}

func (ws *WorkflowSystem) Process(p Part) WorkflowName {
	loc := parseWorkflowName("in")
	for {
		wf, ok := ws.Workflows[loc]
		if !ok {
			break
		}
		loc = wf.Process(p)
	}
	return loc
}

func (ws *WorkflowSystem) TotalOfAcceptedParts() int64 {
	var sum int64
	for _, p := range ws.Parts {
		if ws.Process(p) == ACCEPTED {
			sum += p.Total()
		}
	}
	return sum
}

func (ws *WorkflowSystem) AcceptedPossibilities() int64 {
	var total int64
	queue := []PossibilityState{{
		StateSpace: initialStateSpace(),
		Workflow:   parseWorkflowName("in"),
	}}
	for len(queue) > 0 {
		state := queue[0]
		queue = queue[1:]
		if state.Workflow == ACCEPTED {
			total += state.StateSpace.Volume()
			continue
		}
		wf, ok := ws.Workflows[state.Workflow]
		if !ok {
			continue
		}
		wf.ProcessTesseract(state.StateSpace, &queue)
	}
	return total
}

func (d Day19) Part1(input string) (string, error) {
	ws, err := parseWorkflowSystem(input)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(ws.TotalOfAcceptedParts(), 10), nil
}

func (d Day19) Part2(input string) (string, error) {
	ws, err := parseWorkflowSystem(input)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(ws.AcceptedPossibilities(), 10), nil
}

func init() {
	solve.Register(Day19{})
}
