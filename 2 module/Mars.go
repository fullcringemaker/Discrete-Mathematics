package main

import (
	"container/list"
	"fmt"
	"sort"
)

type Network struct {
	nodes []*Node
}

func identifyClusters(n Network) []Network {
	var clusters []Network
	visitedNodes := make(map[*Node]bool)
	queue := []*Node{}
	for _, node := range n.nodes {
		if !visitedNodes[node] {
			cluster := Network{}
			queue = append(queue, node)
			for len(queue) > 0 {
				currentNode := queue[0]
				queue = queue[1:]
				if !visitedNodes[currentNode] {
					visitedNodes[currentNode] = true
					cluster.nodes = append(cluster.nodes, currentNode)
					for _, neighbor := range currentNode.connections {
						if !visitedNodes[neighbor] {
							queue = append(queue, neighbor)
						}
					}
				}
			}
			clusters = append(clusters, cluster)
		}
	}
	return clusters
}

type Node struct {
	connections []*Node
	index       int
	visited     bool
	group       int
}

func breadthFirstSearch(startNode *Node, group int) (int, bool, Network) {
	queue := list.New()
	startNode.group = group
	queue.PushBack(startNode)
	counter := 0
	isBipartite := true
	subNetwork := Network{}
	for queue.Len() > 0 {
		element := queue.Front()
		queue.Remove(element)
		currentNode := element.Value.(*Node)
		if currentNode.group == 1 {
			counter++
			subNetwork.nodes = append(subNetwork.nodes, currentNode)
		}
		for _, neighbor := range currentNode.connections {
			if neighbor.group == currentNode.group {
				isBipartite = false
			}
			if neighbor.group == -1 {
				neighbor.group = 3 - currentNode.group
				queue.PushBack(neighbor)
			}
		}
	}
	return counter, isBipartite, subNetwork
}

func resetGroups(network *Network) {
	for _, node := range network.nodes {
		node.group = -1
	}
}

func analyzeCluster(network Network) (int, int, Network, Network, bool) {
	if len(network.nodes) == 0 {
		return 0, 0, Network{}, Network{}, false
	}
	resetGroups(&network)
	count1, bipartite1, subNetwork1 := breadthFirstSearch(network.nodes[0], 1)
	if !bipartite1 {
		return 0, 0, Network{}, Network{}, false
	}
	resetGroups(&network)
	count2, bipartite2, subNetwork2 := breadthFirstSearch(network.nodes[0], 2)
	if !bipartite2 {
		return 0, 0, Network{}, Network{}, false
	}
	return count1, count2, subNetwork1, subNetwork2, true
}

func main() {
	var size int
	fmt.Scanf("%d", &size)
	network := createNetwork(size)
	parseConnections(&network, size)
	clusters := identifyClusters(network)
	dpTable, graphData := initializeDPAndGraphData(len(clusters), size)
	if !fillDPTable(&dpTable, &graphData, clusters, size) {
		fmt.Println("No solution")
		return
	}
	resultNodes := assembleResult(dpTable, graphData, clusters, size)
	displayNodes(resultNodes)
}

func createNetwork(size int) Network {
	network := Network{}
	network.nodes = make([]*Node, 0, size)
	for i := 1; i <= size; i++ {
		node := &Node{
			index: i,
			group: -1,
		}
		network.nodes = append(network.nodes, node)
	}
	return network
}

func parseConnections(network *Network, size int) {
	var char rune
	for i := 0; i < size; i++ {
		for j := 0; j < size; {
			fmt.Scanf("%c", &char)
			switch char {
			case '+':
				network.nodes[i].connections = append(network.nodes[i].connections, network.nodes[j])
				fallthrough
			case '-':
				j++
			}
		}
	}
}

func initializeDPAndGraphData(clusterCount, totalSize int) ([][]bool, [][]Network) {
	dpMatrix := make([][]bool, clusterCount+1)
	graphMatrix := make([][]Network, clusterCount+1)
	populateMatrices := func(rows, cols int) ([][]bool, [][]Network) {
		boolMatrix := make([][]bool, rows)
		networkMatrix := make([][]Network, rows)
		for i := 0; i < rows; i++ {
			boolMatrix[i] = make([]bool, cols)
			networkMatrix[i] = make([]Network, cols)
		}
		return boolMatrix, networkMatrix
	}
	dpMatrix, graphMatrix = populateMatrices(clusterCount+1, totalSize/2+1)
	if len(dpMatrix) > 0 && len(dpMatrix[0]) > 0 {
		dpMatrix[0][0] = true
	}

	return dpMatrix, graphMatrix
}

func fillDPTable(dpMatrix *[][]bool, graphMatrix *[][]Network, clusters []Network, size int) bool {
	type state struct {
		groupACount int
		groupBCount int
	}

	states := make([][]state, len(clusters)+1)
	for i := range states {
		states[i] = make([]state, size/2+1)
	}
	for i := 1; i <= len(clusters); i++ {
		groupA, groupB, subNetworkA, subNetworkB, valid := analyzeCluster(clusters[len(clusters)-i])
		if !valid {
			return false
		}
		for j := 0; j <= size/2; j++ {
			previousState := states[i-1][j]
			updated := false
			if j >= groupA && (*dpMatrix)[i-1][j-groupA] {
				states[i][j] = state{groupACount: previousState.groupACount + 1, groupBCount: previousState.groupBCount}
				(*dpMatrix)[i][j] = true
				(*graphMatrix)[i][j] = subNetworkA
				updated = true
			}
			if j >= groupB && (*dpMatrix)[i-1][j-groupB] {
				if !updated || previousState.groupACount+previousState.groupBCount > states[i][j].groupACount+states[i][j].groupBCount {
					states[i][j] = state{groupACount: previousState.groupACount, groupBCount: previousState.groupBCount + 1}
					(*dpMatrix)[i][j] = true
					(*graphMatrix)[i][j] = subNetworkB
				}
			}
		}
	}
	return true
}

func assembleResult(dpMatrix [][]bool, graphMatrix [][]Network, clusters []Network, size int) []*Node {
	var resultNodes []*Node
	clusterCount := len(clusters)
	target := size / 2
	stack := make([][2]int, 0)
	stack = append(stack, [2]int{clusterCount, target})
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		i, j := current[0], current[1]
		if i == 0 || j == 0 {
			continue
		}
		if dpMatrix[i][j] {
			resultNodes = append(resultNodes, graphMatrix[i][j].nodes...)
			stack = append(stack, [2]int{i - 1, j - len(graphMatrix[i][j].nodes)})
		} else {
			stack = append(stack, [2]int{i, j - 1})
		}
	}
	sort.Slice(resultNodes, func(i, j int) bool {
		return resultNodes[i].index < resultNodes[j].index
	})

	return resultNodes
}

func displayNodes(nodes []*Node) {
	for _, node := range nodes {
		fmt.Printf("%d ", node.index)
	}
	fmt.Println()
}
