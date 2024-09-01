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
	id      int
	minCost int
	pos     int
}

type MinHeap []*Node

func (h MinHeap) Len() int { return len(h) }

func (h MinHeap) Less(i, j int) bool {
	return h[i].minCost < h[j].minCost
}

func (h MinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].pos = i
	h[j].pos = j
}

func (h *MinHeap) Push(x interface{}) {
	node := x.(*Node)
	node.pos = len(*h)
	*h = append(*h, node)
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	node := old[n-1]
	node.pos = -1
	*h = old[0 : n-1]
	return node
}

func updateCost(h *MinHeap, node *Node, cost int) {
	node.minCost = cost
	heap.Fix(h, node.pos)
}

func minimalTelephoneLines(buildings map[int][][2]int, n int) int {
	h := &MinHeap{}
	connected := make([]bool, n)
	nodeList := make([]*Node, n)

	for i := 0; i < n; i++ {
		node := &Node{id: i, minCost: int(^uint(0) >> 1)}
		nodeList[i] = node
		heap.Push(h, node)
	}
	nodeList[0].minCost = 0
	updateCost(h, nodeList[0], 0)

	totalCost := 0

	for h.Len() > 0 {
		node := heap.Pop(h).(*Node)
		u := node.id
		totalCost += node.minCost
		connected[u] = true

		for _, road := range buildings[u] {
			v := road[0]
			length := road[1]
			if !connected[v] && nodeList[v].minCost > length {
				updateCost(h, nodeList[v], length)
			}
		}
	}

	return totalCost
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())
	scanner.Scan()
	m, _ := strconv.Atoi(scanner.Text())

	buildings := make(map[int][][2]int)

	for i := 0; i < m; i++ {
		scanner.Scan()
		edges := strings.Split(scanner.Text(), " ")
		u, _ := strconv.Atoi(edges[0])
		v, _ := strconv.Atoi(edges[1])
		length, _ := strconv.Atoi(edges[2])

		buildings[u] = append(buildings[u], [2]int{v, length})
		buildings[v] = append(buildings[v], [2]int{u, length})
	}

	fmt.Println(minimalTelephoneLines(buildings, n))
}
