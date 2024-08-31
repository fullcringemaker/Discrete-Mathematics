package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type StateList []*StateEntity

func (sl *StateList) StateOperations(opType string, idx1 int, idx2 int) interface{} {
	switch opType {
	case "size":
		return len(*sl)
	case "isSmaller":
		return (*sl)[idx1].index < (*sl)[idx2].index
	case "swap":
		(*sl)[idx1], (*sl)[idx2] = (*sl)[idx2], (*sl)[idx1]
		return nil
	default:
		return nil
	}
}

type StateEntity struct {
	index    int
	parent   *StateEntity
	rank     int
	visited  bool
	transits StateList
	outputs  []string
}

func CombineStates(stateOne, stateTwo *StateEntity) {
	rootOne := stateOne
	for rootOne.parent != rootOne {
		rootOne = rootOne.parent
	}
	rootTwo := stateTwo
	for rootTwo.parent != rootTwo {
		rootTwo = rootTwo.parent
	}
	if rootOne.rank < rootTwo.rank {
		rootOne.parent = rootTwo
	} else {
		rootTwo.parent = rootOne
		if rootOne != rootTwo && rootOne.rank == rootTwo.rank {
			rootOne.rank++
		}
	}
}

var totalStates, numSymbols int

func InitialSplit(stateList StateList, partitionCount *int, parentList *StateList) {
	index := 0
	for index < len(stateList) {
		state := stateList[index]
		state.parent = state
		state.rank = 0
		index++
	}
	*partitionCount = totalStates
	CheckEquivalence(stateList, partitionCount)
	for _, state := range stateList {
		root := state
		for root.parent != root {
			root = root.parent
		}
		(*parentList)[state.index] = root
	}
}

func CheckEquivalence(stateList StateList, partitionCount *int) {
	i := 0
	for i < totalStates {
		j := i + 1
		for j < totalStates {
			rootI := stateList[i]
			for rootI.parent != rootI {
				rootI = rootI.parent
			}
			rootJ := stateList[j]
			for rootJ.parent != rootJ {
				rootJ = rootJ.parent
			}
			if rootI != rootJ {
				equal := true
				c := 0
				for c < numSymbols {
					if rootI.outputs[c] != rootJ.outputs[c] {
						equal = false
						break
					}
					c++
				}
				if equal {
					rootI.parent = rootJ
					*partitionCount = *partitionCount - 1
				}
			}
			j++
		}
		i++
	}
}

func FineTunePartitions(stateList StateList, partitionCount *int, parentList *StateList) {
	for _, state := range stateList {
		state.parent = state
		state.rank = 0
	}
	RefinePartition(stateList, partitionCount, parentList, totalStates, numSymbols)
	for _, state := range stateList {
		root := state
		for root.parent != root {
			root = root.parent
		}
		(*parentList)[state.index] = root
	}
}

func RefinePartition(stateList StateList, partitionCount *int, parentList *StateList, stateCount int, symbolCount int) {
	*partitionCount = stateCount
	i := 0
	for i < stateCount-1 {
		j := i + 1
		for j < stateCount {
			rootI := stateList[i]
			for rootI.parent != rootI {
				rootI = rootI.parent
			}
			rootJ := stateList[j]
			for rootJ.parent != rootJ {
				rootJ = rootJ.parent
			}
			if (*parentList)[stateList[i].index] == (*parentList)[stateList[j].index] && rootI != rootJ {
				equal := true
				c := 0
				for c < symbolCount {
					w1 := stateList[i].transits[c].index
					w2 := stateList[j].transits[c].index
					if (*parentList)[w1] != (*parentList)[w2] {
						equal = false
						break
					}
					c++
				}
				if equal {
					rootI.parent = rootJ
					*partitionCount = *partitionCount - 1
				}
			}
			j++
		}
		i++
	}
}

func AutomataReduction(automata StateList) StateList {
	var count1 int
	parentRefs := make(StateList, totalStates)
	InitialSplit(automata, &count1, &parentRefs)
	for {
		var count2 int
		FineTunePartitions(automata, &count2, &parentRefs)
		if count1 != count2 {
			count1 = count2
		} else {
			break
		}
		count1 = count2
	}
	reducedStates := StateList{}
	for _, state := range automata {
		reducedRoot := parentRefs[state.index]
		found := false
		for _, s := range reducedStates {
			if s == reducedRoot {
				found = true
				break
			}
		}
		if !found {
			reducedStates = append(reducedStates, reducedRoot)
			var i = 0
			for i < numSymbols {
				reducedRoot.transits[i] = parentRefs[state.transits[i].index]
				reducedRoot.outputs[i] = state.outputs[i]
				i++
			}
		}
	}
	return reducedStates
}

var num, initialState int

func main() {
	num = 0
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Fscan(inputReader, &totalStates, &numSymbols, &initialState)
	automata := make(StateList, totalStates)
	var i = 0
	for i < totalStates {
		automata[i] = &StateEntity{i, nil, 0, false, make(StateList, numSymbols), make([]string, numSymbols)}
		i++
	}
	i = 0
	for i < totalStates {
		j := 0
		for j < numSymbols {
			var t int
			fmt.Fscan(inputReader, &t)
			automata[i].transits[j] = automata[t]
			j++
		}
		i++
	}
	i = 0
	for i < totalStates {
		j := 0
		for j < numSymbols {
			var s string
			fmt.Fscan(inputReader, &s)
			automata[i].outputs[j] = s
			j++
		}
		i++
	}

	minAutomata := AutomataReduction(automata)
	OutputGraph(minAutomata, automata, initialState, numSymbols)
}

func OutputGraph(minAutomata StateList, automata StateList, startState int, symbolCount int) {
	var counter int
	for _, q := range minAutomata {
		match := true
		for i := 0; i < symbolCount; i++ {
			if q.outputs[i] != automata[startState].outputs[i] {
				match = false
				break
			}
		}
		if match {
			stack := []*StateEntity{q}
			for len(stack) > 0 {
				current := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if !current.visited {
					current.index = counter
					counter++
					current.visited = true
					for i := symbolCount - 1; i >= 0; i-- {
						if !current.transits[i].visited {
							stack = append(stack, current.transits[i])
						}
					}
				}
			}
			break
		}
	}
	sort.Slice(minAutomata, func(i, j int) bool {
		return minAutomata.StateOperations("isSmaller", i, j).(bool)
	})
	fmt.Printf("digraph {\n")
	fmt.Printf("    rankdir = LR\n")
	for _, q := range minAutomata {
		for i := 0; i < symbolCount; i++ {
			symbol := fmt.Sprintf("%c", i+97)
			if i >= 26 {
				symbol = fmt.Sprintf("\\x%02x", i+97)
			}
			fmt.Printf("\t%d -> %d [label = \"%s(%s)\"]\n", q.index, q.transits[i].index, symbol, q.outputs[i])
		}
	}
	fmt.Printf("}\n")
}
