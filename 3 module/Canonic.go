package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type States []*State

func (s *States) ManageStates(operation string, index1 int, index2 int) interface{} {
	switch operation {
	case "length":
		return len(*s)
	case "compareLess":
		return (*s)[index1].order < (*s)[index2].order
	case "exchange":
		(*s)[index1], (*s)[index2] = (*s)[index2], (*s)[index1]
		return nil
	default:
		return nil
	}
}

type State struct {
	NumParents  int
	Transitions []Transition
	order       int
	marked      bool
}

type Transition struct {
	Target *State
	Output string
}

var (
	numStates       int
	numTransitions  int
	startStateIndex int
	states          States
)

func CreateStates() {
	for i := 0; i < numStates; i++ {
		states = append(states, &State{NumParents: 0, Transitions: []Transition{}, order: i, marked: false})
	}
}

func ConnectStates(states []*State, transitionMatrix [][]int, outputsMatrix [][]string) {
	numStates := len(states)
	if numStates == 0 {
		return
	}
	numTransitions := len(transitionMatrix[0])
	i := 0
	for i < numStates {
		j := 0
		for j < numTransitions {
			nextState := states[transitionMatrix[i][j]]
			currentState := states[i]
			nextState.NumParents++
			newTransition := Transition{Target: nextState, Output: outputsMatrix[i][j]}
			currentState.Transitions = append(currentState.Transitions, newTransition)
			j++
		}
		i++
	}
}

type DFSStackEntry struct {
	state *State
}


func main() {
	stdinReader := bufio.NewReader(os.Stdin)
	stdoutWriter := bufio.NewWriter(os.Stdout)
	fmt.Fscan(stdinReader, &numStates, &numTransitions, &startStateIndex)
	var transitionMatrix [][]int
	var outputsMatrix [][]string
	transitionMatrix = make([][]int, numStates)
	outputsMatrix = make([][]string, numStates)
	for i := 0; i < numStates; i++ {
		transitionMatrix[i] = make([]int, numTransitions)
		outputsMatrix[i] = make([]string, numTransitions)
	}
	CreateStates()
	i := 0
	for i < numStates {
		j := 0
		for j < numTransitions {
			fmt.Fscan(stdinReader, &transitionMatrix[i][j])
			j++
		}
		i++
	}
	i = 0
	for i < numStates {
		j := 0
		for j < numTransitions {
			fmt.Fscan(stdinReader, &outputsMatrix[i][j])
			j++
		}
		i++
	}
	ConnectStates(states, transitionMatrix, outputsMatrix)
	stack := []DFSStackEntry{{state: states[startStateIndex]}}
	var order = 0
	for len(stack) > 0 {
		entry := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if entry.state.marked {
			continue
		}
		entry.state.marked = true
		entry.state.order = order
		order++
		for i := len(entry.state.Transitions) - 1; i >= 0; i-- {
			transition := entry.state.Transitions[i]
			if !transition.Target.marked {
				stack = append(stack, DFSStackEntry{state: transition.Target})
			}
		}
	}
	fmt.Fprint(stdoutWriter, numStates, "\n", numTransitions, "\n", 0, "\n")
	sort.Slice(states, func(i, j int) bool { return states[i].order < states[j].order })
	for _, state := range states {
		for _, t := range state.Transitions {
			fmt.Fprint(stdoutWriter, t.Target.order, " ")
		}
		fmt.Fprint(stdoutWriter, "\n")
	}
	for _, state := range states {
		for _, t := range state.Transitions {
			fmt.Fprint(stdoutWriter, t.Output, " ")
		}
		fmt.Fprint(stdoutWriter, "\n")
	}
	stdoutWriter.Flush()
}
