package main

import (
	"fmt"
	"sort"
)

type Node struct {
	Adjacents []int
	Visited   bool
}

type Graph struct {
	Nodes []Node
}

type ComponentInfo struct {
	NodeCount  int
	LinkCount  int
	SmallestID int
	NodeIDs    map[int]bool
}

func NewGraph(size int) *Graph {
	g := &Graph{
		Nodes: make([]Node, size),
	}
	for i := range g.Nodes {
		g.Nodes[i].Adjacents = []int{}
	}
	return g
}

func (g *Graph) AddEdge(a, b int) {
	g.Nodes[a].Adjacents = append(g.Nodes[a].Adjacents, b)
	g.Nodes[b].Adjacents = append(g.Nodes[b].Adjacents, a)
}

func (g *Graph) ExploreComponent(start int, component *ComponentInfo) {
	stack := []int{start}
	g.Nodes[start].Visited = true
	component.NodeCount++
	component.SmallestID = min(component.SmallestID, start)
	component.NodeIDs[start] = true

	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		for _, adj := range g.Nodes[v].Adjacents {
			component.LinkCount++
			if !g.Nodes[adj].Visited {
				g.Nodes[adj].Visited = true
				stack = append(stack, adj)
				component.NodeCount++
				component.SmallestID = min(component.SmallestID, adj)
				component.NodeIDs[adj] = true
			}
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	var n, m, a, b int
	fmt.Scanf("%d\n", &n)
	fmt.Scanf("%d\n", &m)

	graph := NewGraph(n)

	for i := 0; i < m; i++ {
		fmt.Scanf("%d %d\n", &a, &b)
		graph.AddEdge(a, b)
	}

	var components []ComponentInfo
	for i := range graph.Nodes {
		if !graph.Nodes[i].Visited {
			comp := ComponentInfo{
				NodeCount:  0,
				LinkCount:  0,
				SmallestID: n,
				NodeIDs:    make(map[int]bool),
			}
			graph.ExploreComponent(i, &comp)
			components = append(components, comp)
		}
	}

	sort.Slice(components, func(i, j int) bool {
		if components[i].NodeCount != components[j].NodeCount {
			return components[i].NodeCount > components[j].NodeCount
		}
		if components[i].LinkCount != components[j].LinkCount {
			return components[i].LinkCount > components[j].LinkCount
		}
		return components[i].SmallestID < components[j].SmallestID
	})

	largestComponent := components[0]

	fmt.Printf("graph {\n")
	for i := 0; i < n; i++ {
		fmt.Printf("%d", i)
		if largestComponent.NodeIDs[i] {
			fmt.Printf(" [color = red]")
		}
		fmt.Printf("\n")
	}

	for i := 0; i < n; i++ {
		for _, adj := range graph.Nodes[i].Adjacents {
			if i < adj {
				fmt.Printf("%d -- %d", i, adj)
				if largestComponent.NodeIDs[i] {
					fmt.Printf(" [color = red]")
				}
				fmt.Printf("\n")
			}
		}
	}
	fmt.Printf("}\n")
}
