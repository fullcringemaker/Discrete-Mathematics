package main

import (
	"fmt"
	"sort"
)

func AreSlicesIdentical(sliceA, sliceB []int) bool {
	if len(sliceA) != len(sliceB) {
		return false
	}
	elementFreq := make(map[int]int)
	for _, num := range sliceA {
		elementFreq[num]++
	}
	for _, num := range sliceB {
		elementFreq[num]--
		if elementFreq[num] < 0 {
			return false
		}
	}
	return true
}

func DepthFirstTraversal(adjList [][][]string, startNode int,
	visitedNodes *[]int, processedNodes *[]bool, epsilonTransition string) {
	stack := []int{startNode}
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if !(*processedNodes)[node] {
			*visitedNodes = append(*visitedNodes, node)
			(*processedNodes)[node] = true
			for neighbor, transitions := range adjList[node] {
				for _, transition := range transitions {
					if transition == epsilonTransition {
						stack = append(stack, neighbor)
						break
					}
				}
			}
		}
	}
}

func CalculateEpsilonClosure(adjList [][][]string, initialNodes []int, epsilonTransition string) []int {
	closureSet := make(map[int]bool)
	queue := append([]int{}, initialNodes...)
	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:]
		if !closureSet[currentNode] {
			closureSet[currentNode] = true
			for neighbor, transitions := range adjList[currentNode] {
				for _, transition := range transitions {
					if transition == epsilonTransition && !closureSet[neighbor] {
						queue = append(queue, neighbor)
						break
					}
				}
			}
		}
	}
	result := []int{}
	for node := range closureSet {
		result = append(result, node)
	}
	sort.Ints(result)
	return result
}

func FindReachableStates(adjList [][][]string, currentStates []int, symbol string) []int {
	var reachableNodes []int
	for _, state := range currentStates {
		for target, transitions := range adjList[state] {
			for _, transition := range transitions {
				if transition == symbol {
					reachableNodes = append(reachableNodes, target)
					break
				}
			}
		}
	}
	return reachableNodes
}

func GetStateIndex(state []int, stateCollection [][]int) int {
	sort.Ints(state)
	for idx, existingState := range stateCollection {
		sort.Ints(existingState)
		if AreSlicesIdentical(existingState, state) {
			return idx
		}
	}
	return -1
}

func main() {
	var stateCount, transitionCount int
	fmt.Scan(&stateCount, &transitionCount)
	adjacencyList := make([][][]string, stateCount)
	for i := range adjacencyList {
		adjacencyList[i] = make([][]string, stateCount)
	}
	inputSymbols := []string{}
	epsilonSymbol := "lambda"
	for i := 0; i < transitionCount; i++ {
		var fromState, toState int
		var transitionLabel string
		fmt.Scan(&fromState, &toState, &transitionLabel)
		labelExists := false
		for _, symbol := range inputSymbols {
			if symbol == transitionLabel {
				labelExists = true
				break
			}
		}
		if transitionLabel != epsilonSymbol && !labelExists {
			inputSymbols = append(inputSymbols, transitionLabel)
		}
		adjacencyList[fromState][toState] = append(adjacencyList[fromState][toState], transitionLabel)
	}
	finalStateFlags := make([]int, stateCount)
	for i := range finalStateFlags {
		fmt.Scan(&finalStateFlags[i])
	}
	var startingState int
	fmt.Scan(&startingState)
	initialClosure := CalculateEpsilonClosure(adjacencyList, []int{startingState}, epsilonSymbol)
	stateQueue := [][]int{initialClosure}
	allDiscoveredStates := [][]int{initialClosure}
	finalStateSets := [][]int{}
	stateTransitions := make([]map[int][]string, 0)
	for len(stateQueue) > 0 {
		currentStateSet := stateQueue[0]
		stateQueue = stateQueue[1:]
		isFinal := false
		for _, state := range currentStateSet {
			if finalStateFlags[state] == 1 {
				isFinal = true
				break
			}
		}
		if isFinal {
			finalStateSets = append(finalStateSets, currentStateSet)
		}
		currentIdx := GetStateIndex(currentStateSet, allDiscoveredStates)
		if len(stateTransitions) <= currentIdx {
			stateTransitions = append(stateTransitions, make(map[int][]string))
		}
		for _, symbol := range inputSymbols {
			nextStateSet := FindReachableStates(adjacencyList, currentStateSet, symbol)
			closure := CalculateEpsilonClosure(adjacencyList, nextStateSet, epsilonSymbol)

			stateIdx := GetStateIndex(closure, allDiscoveredStates)
			if stateIdx == -1 {
				allDiscoveredStates = append(allDiscoveredStates, closure)
				stateQueue = append(stateQueue, closure)
				stateIdx = len(allDiscoveredStates) - 1
			}
			if _, exists := stateTransitions[currentIdx][stateIdx]; !exists {
				stateTransitions[currentIdx][stateIdx] = []string{}
			}
			stateTransitions[currentIdx][stateIdx] =
				append(stateTransitions[currentIdx][stateIdx], symbol)
		}
	}
	fmt.Println("digraph {")
	fmt.Println("\trankdir = LR")
	for i, state := range allDiscoveredStates {
		fmt.Printf("\t%d [label = \"%v\", shape = ", i, state)
		if GetStateIndex(state, finalStateSets) != -1 {
			fmt.Println("doublecircle]")
		} else {
			fmt.Println("circle]")
		}
	}
	for i := range allDiscoveredStates {
		for j, labels := range stateTransitions[i] {
			fmt.Printf("\t%d -> %d [label = \"", i, j)
			for k, label := range labels {
				fmt.Print(label)
				if k < len(labels)-1 {
					fmt.Print(", ")
				}
			}
			fmt.Println("\"]")
		}
	}
	fmt.Println("}")
}
