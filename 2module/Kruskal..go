package main

import (
	"fmt"
	"math"
	"sort"
)

type Edge struct {
	start, end int
	weight     float64
}

type UnionFind struct {
	parent, rank []int
}

func NewUnionFind(size int) *UnionFind {
	parent := make([]int, size)
	rank := make([]int, size)
	for i := range parent {
		parent[i] = i
	}
	return &UnionFind{parent, rank}
}

func (uf *UnionFind) find(i int) int {
	if uf.parent[i] != i {
		uf.parent[i] = uf.find(uf.parent[i])
	}
	return uf.parent[i]
}

func (uf *UnionFind) union(x, y int) {
	rootX := uf.find(x)
	rootY := uf.find(y)
	if rootX != rootY {
		if uf.rank[rootX] > uf.rank[rootY] {
			uf.parent[rootY] = rootX
		} else if uf.rank[rootX] < uf.rank[rootY] {
			uf.parent[rootX] = rootY
		} else {
			uf.parent[rootY] = rootX
			uf.rank[rootX]++
		}
	}
}

func kruskal(n int, edges []Edge) float64 {
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].weight < edges[j].weight
	})

	uf := NewUnionFind(n)
	result := 0.0
	for _, edge := range edges {
		if uf.find(edge.start) != uf.find(edge.end) {
			uf.union(edge.start, edge.end)
			result += edge.weight
		}
	}
	return result
}

func main() {
	var n int
	fmt.Scan(&n)
	positions := make([][2]int, n)
	for i := range positions {
		fmt.Scan(&positions[i][0], &positions[i][1])
	}

	var edges []Edge
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			dist := math.Sqrt(float64((positions[i][0]-positions[j][0])*(positions[i][0]-positions[j][0]) + (positions[i][1]-positions[j][1])*(positions[i][1]-positions[j][1])))
			edges = append(edges, Edge{i, j, dist})
		}
	}

	result := kruskal(n, edges)
	fmt.Printf("%.2f\n", result)
}
