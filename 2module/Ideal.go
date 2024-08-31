package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Node struct {
	id    int
	dist  int
	path  []int
	links []Connection
	index int
}

type Connection struct {
	target *Node
	color  int
}

type MinHeap []*Node

func (h MinHeap) Len() int { return len(h) }
func (h MinHeap) Less(i, j int) bool {
	return h[i].dist < h[j].dist || (h[i].dist == h[j].dist && pathPriority(h[i].path, h[j].path) < 0)
}
func (h MinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}
func (h *MinHeap) Push(x interface{}) {
	n := len(*h)
	node := x.(*Node)
	node.index = n
	*h = append(*h, node)
}
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	node := old[n-1]
	old[n-1] = nil
	node.index = -1
	*h = old[:n-1]
	return node
}

func pathPriority(a, b []int) int {
	minLength := len(a)
	if len(b) < minLength {
		minLength = len(b)
	}
	for i := 0; i < minLength; i++ {
		if a[i] < b[i] {
			return -1
		} else if a[i] > b[i] {
			return 1
		}
	}
	return len(a) - len(b)
}

func findShortestPath(nodes []*Node, start *Node) {
	pq := &MinHeap{}
	heap.Init(pq)
	initializeNodes(nodes, start, pq)
	processNodes(pq, len(nodes))
}

func initializeNodes(nodes []*Node, start *Node, pq *MinHeap) {
	for _, node := range nodes {
		node.dist = 1<<31 - 1 
		node.path = nil
		if node == start {
			node.dist = 0
			node.path = []int{}
		}
		heap.Push(pq, node)
	}
}

func processNodes(pq *MinHeap, targetID int) {
	for pq.Len() > 0 {
		u := heap.Pop(pq).(*Node)
		if u.id == targetID {
			break
		}
		updateDistances(u, pq)
	}
}

func updateDistances(u *Node, pq *MinHeap) {
	for _, conn := range u.links {
		v := conn.target
		newPath := append([]int(nil), u.path...)
		newPath = append(newPath, conn.color)
		newDist := u.dist + 1
		if newDist < v.dist || (newDist == v.dist && pathPriority(newPath, v.path) < 0) {
			v.dist = newDist
			v.path = newPath
			heap.Fix(pq, v.index)
		}
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	nm := strings.Split(scanner.Text(), " ")
	n, _ := strconv.Atoi(nm[0])
	m, _ := strconv.Atoi(nm[1])

	nodes := make([]*Node, n)
	for i := range nodes {
		nodes[i] = &Node{id: i + 1, links: make([]Connection, 0)}
	}
	readGraph(scanner, m, nodes)

	findShortestPath(nodes, nodes[0])
	displayPath(nodes[n-1])
}

func readGraph(scanner *bufio.Scanner, m int, nodes []*Node) {
	for i := 0; i < m; i++ {
		scanner.Scan()
		edgeData := strings.Split(scanner.Text(), " ")
		u, _ := strconv.Atoi(edgeData[0])
		v, _ := strconv.Atoi(edgeData[1])
		c, _ := strconv.Atoi(edgeData[2])
		u--
		v--
		nodes[u].links = append(nodes[u].links, Connection{target: nodes[v], color: c})
		if u != v {
			nodes[v].links = append(nodes[v].links, Connection{target: nodes[u], color: c})
		}
	}
}

func displayPath(node *Node) {
	fmt.Println(len(node.path))
	if len(node.path) > 0 {
		fmt.Print(node.path[0])
		for _, color := range node.path[1:] {
			fmt.Print(" ", color)
		}
		fmt.Println()
	}
}
