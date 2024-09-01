package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Node struct {
	id          int
	connections map[int]*Node
}

func CreateNode(id int) *Node {
	return &Node{
		id:          id,
		connections: make(map[int]*Node),
	}
}

type Network struct {
	nodes map[int]*Node
}

func CreateNetwork() *Network {
	return &Network{
		nodes: make(map[int]*Node),
	}
}

func (net *Network) ConnectNodes(a, b int) {
	if _, exists := net.nodes[a]; !exists {
		net.nodes[a] = CreateNode(a)
	}
	if _, exists := net.nodes[b]; !exists {
		net.nodes[b] = CreateNode(b)
	}
	net.nodes[a].connections[b] = net.nodes[b]
	net.nodes[b].connections[a] = net.nodes[a]
}

type SearchState struct {
	visited map[int]bool
	dTime   map[int]int
	low     map[int]int
	parent  map[int]int
	time    int
	bridges int
	network *Network
}

func InitializeSearchState(network *Network) *SearchState {
	return &SearchState{
		visited: make(map[int]bool),
		dTime:   make(map[int]int),
		low:     make(map[int]int),
		parent:  make(map[int]int),
		time:    0,
		bridges: 0,
		network: network,
	}
}

func (state *SearchState) minimum(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (state *SearchState) DepthFirstSearch(id int) {
	state.visited[id] = true
	state.dTime[id] = state.time
	state.low[id] = state.time
	state.time++

	for _, neighbor := range state.network.nodes[id].connections {
		if !state.visited[neighbor.id] {
			state.parent[neighbor.id] = id
			state.DepthFirstSearch(neighbor.id)
			state.low[id] = state.minimum(state.low[id], state.low[neighbor.id])

			if state.low[neighbor.id] > state.dTime[id] {
				state.bridges++
			}
		} else if neighbor.id != state.parent[id] {
			state.low[id] = state.minimum(state.low[id], state.dTime[neighbor.id])
		}
	}
}

func (state *SearchState) CountBridges() int {
	for id := range state.network.nodes {
		if !state.visited[id] {
			state.DepthFirstSearch(id)
		}
	}
	return state.bridges
}

func main() {
	inputScanner := bufio.NewScanner(os.Stdin)
	inputScanner.Scan()
	_, _ = strconv.Atoi(inputScanner.Text())
	inputScanner.Scan()
	totalEdges, _ := strconv.Atoi(inputScanner.Text())

	network := CreateNetwork()

	for i := 0; i < totalEdges; i++ {
		inputScanner.Scan()
		edgeData := strings.Split(inputScanner.Text(), " ")
		a, _ := strconv.Atoi(edgeData[0])
		b, _ := strconv.Atoi(edgeData[1])
		network.ConnectNodes(a, b)
	}

	searchState := InitializeSearchState(network)
	fmt.Println(searchState.CountBridges())
}
