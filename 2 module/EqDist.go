package main

import (
	"fmt"
	"sort"
)

type Vertex struct {
	neighbors []int
	dists     []int
}

func bfs(graph []Vertex, root, id int) {
	queue := make([]int, 0)
	queue = append(queue, root)
	graph[root].dists[id] = 0

	begin := 0
	for begin < len(queue) {
		v := queue[begin]
		begin++
		for _, to := range graph[v].neighbors {
			if graph[to].dists[id] == -1 {
				graph[to].dists[id] = graph[v].dists[id] + 1
				queue = append(queue, to)
			}
		}
	}
}

func main() {
	var n, m int
	fmt.Scan(&n, &m)
	graph := make([]Vertex, n)
	for i := range graph {
		graph[i].dists = make([]int, 0)
	}

	for i := 0; i < m; i++ {
		var from, to int
		fmt.Scan(&from, &to)
		graph[from].neighbors = append(graph[from].neighbors, to)
		graph[to].neighbors = append(graph[to].neighbors, from)
	}

	var k int
	fmt.Scan(&k)
	roots := make([]int, k)
	for i := range graph {
		graph[i].dists = make([]int, k) 
		for j := range graph[i].dists {
			graph[i].dists[j] = -1
		}
	}
	for i := 0; i < k; i++ {
		fmt.Scan(&roots[i])
		bfs(graph, roots[i], i)
	}

	result := make([]int, 0)
	for i := 0; i < n; i++ {
		allEqual := true
		for j := 1; j < k && allEqual; j++ {
			if graph[i].dists[j] != graph[i].dists[0] {
				allEqual = false
			}
		}
		if allEqual && graph[i].dists[0] != -1 { 
			result = append(result, i)
		}
	}

	if len(result) == 0 {
		fmt.Println("-")
	} else {
		sort.Ints(result) 
		for _, v := range result {
			fmt.Print(v, " ")
		}
		fmt.Println()
	}
}
