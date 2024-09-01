package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func readInputGraph() (int, int, [][]int, [][]int, [][2]int) {
	reader := bufio.NewReader(os.Stdin)
	var nodeCount, edgeCount int
	fmt.Fscan(reader, &nodeCount, &edgeCount)
	graph := make([][]int, nodeCount)
	reversedGraph := make([][]int, nodeCount)
	edgeList := make([][2]int, edgeCount)
	for i := 0; i < edgeCount; i++ {
		var startNode, endNode int
		fmt.Fscan(reader, &startNode, &endNode)
		edgeList[i] = [2]int{startNode, endNode}
		graph[startNode] = append(graph[startNode], endNode)
		reversedGraph[endNode] = append(reversedGraph[endNode], startNode)
	}
	return nodeCount, edgeCount, graph, reversedGraph, edgeList
}

func calculateFillOrder(nodeCount int, graph [][]int, isVisited []bool) []int {
	type Frame struct {
		node          int
		phase         int
		neighborIndex int
	}
	orderStack := []int{}
	recursionStack := []Frame{}
	recursionStack = append(recursionStack, Frame{node: -1, phase: 0, neighborIndex: 0})
	for len(recursionStack) > 0 {
		currentFrame := &recursionStack[len(recursionStack)-1]
		switch currentFrame.phase {
		case 0:
			if currentFrame.node != -1 && isVisited[currentFrame.node] {
				recursionStack = recursionStack[:len(recursionStack)-1]
				continue
			}
			if currentFrame.node != -1 {
				isVisited[currentFrame.node] = true
			}
			currentFrame.phase = 1
		case 1:
			if currentFrame.node != -1 {
				for currentFrame.neighborIndex < len(graph[currentFrame.node]) {
					neighborNode := graph[currentFrame.node][currentFrame.neighborIndex]
					currentFrame.neighborIndex++
					if !isVisited[neighborNode] {
						recursionStack = append(recursionStack, Frame{node: neighborNode, phase: 0, neighborIndex: 0})
						break
					}
				}
				if currentFrame.neighborIndex == len(graph[currentFrame.node]) {
					currentFrame.phase = 2
				}
			} else {
				for i := 0; i < nodeCount; i++ {
					if !isVisited[i] {
						recursionStack = append(recursionStack, Frame{node: i, phase: 0, neighborIndex: 0})
						break
					}
				}
				if len(recursionStack) > 0 && recursionStack[len(recursionStack)-1].node == -1 {
					recursionStack = recursionStack[:len(recursionStack)-1]
				}
			}
		case 2:
			if currentFrame.node != -1 {
				orderStack = append(orderStack, currentFrame.node)
			}
			recursionStack = recursionStack[:len(recursionStack)-1]
		}
	}
	return orderStack
}

func findStrongConnectComponents(nodeCount int, orderStack []int, reversedGraph [][]int) ([]int, int) {
	componentLabels := make([]int, nodeCount)
	for i := range componentLabels {
		componentLabels[i] = -1
	}
	type NodeInformation struct {
		node      int
		component int
	}
	componentCount := 0
	assignComponentLabels := func(startNode, component int) {
		nodeQueue := []NodeInformation{{startNode, component}}
		nodeLevels := make([]int, nodeCount)
		for i := range nodeLevels {
			nodeLevels[i] = -1
		}
		level := 0
		nodeLevels[startNode] = level
		for len(nodeQueue) > 0 {
			level++
			currentLevelQueue := []NodeInformation{}
			nextLevelQueue := []NodeInformation{}
			for _, info := range nodeQueue {
				node := info.node
				if componentLabels[node] == -1 {
					componentLabels[node] = info.component
					for _, neighbor := range reversedGraph[node] {
						if componentLabels[neighbor] == -1 {
							if nodeLevels[neighbor] == -1 {
								nodeLevels[neighbor] = level
								nextLevelQueue = append(nextLevelQueue, NodeInformation{neighbor, component})
							} else if nodeLevels[neighbor] == level {
								currentLevelQueue = append(currentLevelQueue, NodeInformation{neighbor, component})
							}
						}
					}
				}
			}
			nodeQueue = append(currentLevelQueue, nextLevelQueue...)
		}
	}
	for i := len(orderStack) - 1; i >= 0; i-- {
		if componentLabels[orderStack[i]] == -1 {
			assignComponentLabels(orderStack[i], componentCount)
			componentCount++
		}
	}
	return componentLabels, componentCount
}

func main() {
	nodeCount, _, graph, reversedGraph, edgeList := readInputGraph()
	isVisited := make([]bool, nodeCount)
	orderStack := calculateFillOrder(nodeCount, graph, isVisited)
	componentLabels, componentCount :=  findStrongConnectComponents(nodeCount, orderStack, reversedGraph)
	inDegrees := make([]int, componentCount)
	for _, edge := range edgeList {
		startNode, endNode := edge[0], edge[1]
		if componentLabels[startNode] != componentLabels[endNode] {
			inDegrees[componentLabels[endNode]]++
		}
	}
	minimalVertices := []int{}
	for i := 0; i < componentCount; i++ {
		if inDegrees[i] == 0 {
			minVertex := -1
			for node := 0; node < nodeCount; node++ {
				if componentLabels[node] == i {
					if minVertex == -1 || node < minVertex {
						minVertex = node
					}
				}
			}
			minimalVertices = append(minimalVertices, minVertex)
		}
	}
	sort.Ints(minimalVertices)
	for _, vertex := range minimalVertices {
		fmt.Printf("%d ", vertex)
	}
}
