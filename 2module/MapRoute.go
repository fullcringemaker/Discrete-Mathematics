package main

import "fmt"

type Node struct {
	Distance int
	Pos      int
	Connects []Adjacency
	Cost     int
}

type Adjacency struct {
	Row, Col int
}

type MinHeap struct {
	Elements []*Node
	Size     int
}

func (h *MinHeap) Push(node *Node) {
	position := h.Size
	h.Elements = append(h.Elements, node)
	h.Size++
	for position > 0 && h.Elements[(position-1)/2].Distance > h.Elements[position].Distance {
		h.swap(position, (position-1)/2)
		position = (position - 1) / 2
	}
	h.Elements[position].Pos = position
}

func (h *MinHeap) Pop() *Node {
	root := h.Elements[0]
	h.Size--
	h.Elements[0] = h.Elements[h.Size]
	h.Elements[0].Pos = 0
	h.Elements = h.Elements[:h.Size]
	h.siftDown(0)
	return root
}

func (h *MinHeap) siftDown(position int) {
	smallest := position
	leftChild := 2*position + 1
	rightChild := 2*position + 2

	if leftChild < h.Size && h.Elements[leftChild].Distance < h.Elements[smallest].Distance {
		smallest = leftChild
	}
	if rightChild < h.Size && h.Elements[rightChild].Distance < h.Elements[smallest].Distance {
		smallest = rightChild
	}
	if smallest != position {
		h.swap(position, smallest)
		h.siftDown(smallest)
	}
}

func (h *MinHeap) swap(i, j int) {
	h.Elements[i], h.Elements[j] = h.Elements[j], h.Elements[i]
	h.Elements[i].Pos, h.Elements[j].Pos = i, j
}

func (h *MinHeap) ReduceKey(position, newDistance int) {
	if position < 0 || position >= h.Size {
		return
	}
	h.Elements[position].Distance = newDistance
	for position > 0 && h.Elements[(position-1)/2].Distance > h.Elements[position].Distance {
		h.swap(position, (position-1)/2)
		position = (position - 1) / 2
	}
}

func (h *MinHeap) IsEmpty() bool {
	return h.Size == 0
}

func optimize(u, v *Node, weight int) bool {
	newDistance := u.Distance + weight
	if newDistance < v.Distance {
		v.Distance = newDistance
		return true
	}
	return false
}

func setupNodes(gridSize int) [][]Node {
	grid := make([][]Node, gridSize)
	for i := range grid {
		grid[i] = make([]Node, gridSize)
		for j := range grid[i] {
			grid[i][j].Connects = []Adjacency{}
		}
	}
	return grid
}

func isValidIndex(row, col, size int) bool {
	return row >= 0 && row < size && col >= 0 && col < size
}

func assignWeightsAndNeighbors(nodes [][]Node) {
	directions := []struct {
		dRow, dCol int
	}{
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
	}
	size := len(nodes)

	for row := range nodes {
		for col := range nodes[row] {
			fmt.Scan(&nodes[row][col].Cost)
			for _, dir := range directions {
				adjRow, adjCol := row+dir.dRow, col+dir.dCol
				if isValidIndex(adjRow, adjCol, size) {
					nodes[adjRow][adjCol].Connects = append(nodes[adjRow][adjCol].Connects, Adjacency{Row: row, Col: col})
				}
			}
		}
	}
}

func main() {
	var gridSize int
	fmt.Scan(&gridSize)
	if gridSize == 1 {
		fmt.Scan(&gridSize)
		fmt.Println(gridSize)
		return
	}

	nodes := setupNodes(gridSize)
	assignWeightsAndNeighbors(nodes)

	heap := MinHeap{Elements: make([]*Node, 0, gridSize*gridSize)}
	for row := range nodes {
		for col := range nodes[row] {
			node := &nodes[row][col]
			if row == 0 && col == 0 {
				node.Distance = node.Cost
			} else {
				node.Distance = 150000
			}
			heap.Push(node)
		}
	}

	for !heap.IsEmpty() {
		current := heap.Pop()
		current.Pos = -1
		for _, adj := range current.Connects {
			neighbor := &nodes[adj.Row][adj.Col]
			if neighbor.Pos != -1 && optimize(current, neighbor, neighbor.Cost) {
				heap.ReduceKey(neighbor.Pos, neighbor.Distance)
			}
		}
	}

	fmt.Println(nodes[gridSize-1][gridSize-1].Distance)
}
