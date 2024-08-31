package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type node struct {
	dominator          int
	semiDominator      int
	ancestor           int
	label              int
	parent             int
	timestamp          int
	incoming, outgoing []int
	bucket             map[int]struct{}
}

func readGraph() (int, int, map[int]*node) {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	vertexCount, _ := strconv.Atoi(strings.TrimSpace(line))
	graph := make(map[int]*node)
	var startNode int
	commands := parseInput(reader, vertexCount)
	for i := 0; i < len(commands); i++ {
		cmd := commands[i]
		_, exists := graph[cmd[0]]
		if !exists {
			graph[cmd[0]] = createNode(cmd[0])
		}
		if cmd[1] >= 0 {
			_, exists = graph[cmd[1]]
			if !exists {
				graph[cmd[1]] = createNode(cmd[1])
			}
			graph[cmd[0]].outgoing = append(graph[cmd[0]].outgoing, cmd[1])
			graph[cmd[1]].incoming = append(graph[cmd[1]].incoming, cmd[0])
		}

		if i == 0 {
			startNode = cmd[0]
		} else if commands[i-1][2] != 2 {
			graph[commands[i-1][0]].outgoing = append(graph[commands[i-1][0]].outgoing, cmd[0])
			graph[cmd[0]].incoming = append(graph[cmd[0]].incoming, commands[i-1][0])
		}
	}
	for idx := range commands {
		cmd := commands[idx]
		if _, exists := graph[cmd[0]]; !exists {
			graph[cmd[0]] = createNode(cmd[0])
		}
		if cmd[1] >= 0 {
			if _, exists := graph[cmd[1]]; !exists {
				graph[cmd[1]] = createNode(cmd[1])
			}
			graph[cmd[0]].outgoing = append(graph[cmd[0]].outgoing, cmd[1])
			graph[cmd[1]].incoming = append(graph[cmd[1]].incoming, cmd[0])
		}
		if idx == 0 {
			startNode = cmd[0]
		} else if commands[idx-1][2] != 2 {
			graph[commands[idx-1][0]].outgoing = append(graph[commands[idx-1][0]].outgoing, cmd[0])
			graph[cmd[0]].incoming = append(graph[cmd[0]].incoming, commands[idx-1][0])
		}
	}
	for iter := 0; iter < len(commands); iter++ {
		cmd := commands[iter]
		if _, exists := graph[cmd[0]]; !exists {
			graph[cmd[0]] = createNode(cmd[0])
		}
		if cmd[1] >= 0 {
			if _, exists := graph[cmd[1]]; !exists {
				graph[cmd[1]] = createNode(cmd[1])
			}
			graph[cmd[0]].outgoing = append(graph[cmd[0]].outgoing, cmd[1])
			graph[cmd[1]].incoming = append(graph[cmd[1]].incoming, cmd[0])
		}
		if iter == 0 {
			startNode = cmd[0]
		} else if commands[iter-1][2] != 2 {
			graph[commands[iter-1][0]].outgoing = append(graph[commands[iter-1][0]].outgoing, cmd[0])
			graph[cmd[0]].incoming = append(graph[cmd[0]].incoming, commands[iter-1][0])
		}
	}
	return vertexCount, startNode, graph
}

func createNode(index int) *node {
	return &node{
		dominator:     -1,
		semiDominator: index,
		label:         index,
		ancestor:      -1,
		parent:        -1,
		timestamp:     -1,
		incoming:      []int{},
		outgoing:      []int{},
		bucket:        make(map[int]struct{}),
	}
}

func parseInput(reader *bufio.Reader, count int) [][]int {
	commands := make([][]int, 0, count)
	for i := 0; i < count; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("", err)
			os.Exit(1)
		}
		line = strings.TrimSpace(line)
		if line == "" {
			i--
			continue
		}
		parts := strings.FieldsFunc(line, func(r rune) bool {
			return r == ' ' || r == '\t'
		})

		if len(parts) < 2 {
			fmt.Println("", line)
			os.Exit(1)
		}
		index, err := strconv.Atoi(parts[0])
		if err != nil {
			fmt.Println("", parts[0])
			os.Exit(1)
		}
		command := parts[1]
		operand := -1
		if len(parts) > 2 {
			operand, err = strconv.Atoi(parts[2])
			if err != nil {
				fmt.Println("", parts[2])
				os.Exit(1)
			}
		}
		commandInt := mapCommand(command)
		commands = append(commands, []int{index, operand, commandInt})
	}
	return commands
}

func mapCommand(command string) int {
	switch command {
	case "BRANCH":
		return 1
	case "JUMP":
		return 2
	default:
		return 0
	}
}

func numberAndFilter(n, start int, g map[int]*node) []int {
	type state struct {
		vertex int
		phase  int
	}
	order := make([]int, 0, n)
	visited := make(map[int]struct{})
	parentMap := make(map[int]int)
	time := 0
	vertexStack := list.New()
	phaseStack := list.New()
	vertexStack.PushBack(start)
	phaseStack.PushBack(0)
	parentMap[start] = -1
	for vertexStack.Len() > 0 {
		vs := vertexStack.Back()
		ps := phaseStack.Back()
		vertex := vs.Value.(int)
		phase := ps.Value.(int)
		phaseStack.Remove(ps)
		if phase == 0 {
			if _, ok := visited[vertex]; ok {
				vertexStack.Remove(vs)
				continue
			}
			visited[vertex] = struct{}{}
			order = append(order, vertex)
			g[vertex].parent = parentMap[vertex]
			g[vertex].timestamp = time
			time++
			phaseStack.PushBack(1)
			for _, to := range g[vertex].outgoing {
				if _, ok := visited[to]; !ok {
					vertexStack.PushBack(to)
					phaseStack.PushBack(0)
					parentMap[to] = vertex
				}
			}
		} else {
			vertexStack.Remove(vs)
		}
	}
	for v := range g {
		if _, ok := visited[v]; !ok {
			delete(g, v)
			continue
		}
		filterEdges := func(edges []int) []int {
			newEdges := make([]int, 0, len(edges))
			for _, w := range edges {
				if _, ok := visited[w]; ok {
					newEdges = append(newEdges, w)
				}
			}
			return newEdges
		}
		g[v].incoming = filterEdges(g[v].incoming)
		g[v].outgoing = filterEdges(g[v].outgoing)
	}

	return order
}

func findAndUpdate(v int, g map[int]*node) int {
	var innerFindAndUpdate func(v int, g map[int]*node) int
	innerFindAndUpdate = func(v int, g map[int]*node) int {
		ancestor := g[v].ancestor
		if ancestor == -1 {
			return v
		}
		ancestorNode := ancestor
		root := innerFindAndUpdate(ancestorNode, g)
		labelAncestor := g[ancestorNode].label
		sdomLabelAncestor := g[labelAncestor].semiDominator
		timeSdomLabelAncestor := g[sdomLabelAncestor].timestamp
		labelV := g[v].label
		sdomLabelV := g[labelV].semiDominator
		timeSdomLabelV := g[sdomLabelV].timestamp
		if timeSdomLabelAncestor < timeSdomLabelV {
			if g[v].label != labelAncestor {
				g[v].label = labelAncestor
			}
		} else {
			if g[v].label == labelAncestor {
				if timeSdomLabelAncestor >= timeSdomLabelV {
					g[v].label = labelV
				}
			}
		}
		g[v].ancestor = root
		return root
	}
	finalRoot := innerFindAndUpdate(v, g)
	return finalRoot
}

func main() {
	n, start, graph := readGraph()
	order := numberAndFilter(n, start, graph)
	type tempData struct {
		node, label int
	}
	tempStorage := make([]tempData, 0)
	for i := 0; i < len(graph); i++ {
		w := order[len(graph)-1-i]
		tempStorage = append(tempStorage, tempData{w, graph[w].label})
		for _, v := range graph[w].incoming {
			tempStorage = append(tempStorage, tempData{v, graph[v].label})
			if graph[v] != nil && graph[w] != nil {
				findAndUpdate(v, graph)
				if graph[graph[graph[v].label].semiDominator].timestamp < graph[graph[w].semiDominator].timestamp {
					graph[w].semiDominator = graph[graph[v].label].semiDominator
				}
			}
		}
		if graph[w] != nil {
			graph[w].ancestor = graph[w].parent
			if graph[graph[w].semiDominator] != nil {
				if graph[graph[w].semiDominator].bucket == nil {
					graph[graph[w].semiDominator].bucket = make(map[int]struct{})
				}
				graph[graph[w].semiDominator].bucket[w] = struct{}{}
			}
		}
		if w != order[0] {
			if graph[graph[w].parent] != nil {
				for v := range graph[graph[w].parent].bucket {
					findAndUpdate(v, graph)
					if graph[v] != nil {
						u := graph[v].label
						if graph[u] != nil {
							if graph[graph[u].semiDominator].timestamp == graph[graph[v].semiDominator].timestamp {
								graph[v].dominator = graph[v].semiDominator
							} else {
								graph[v].dominator = u
							}
						}
					}
				}
				graph[graph[w].parent].bucket = make(map[int]struct{})
			}
		}
	}
	tempMap := make(map[int]tempData)
	for _, ts := range tempStorage {
		tempMap[ts.node] = ts
	}
	for i := 1; i < len(graph); i++ {
		w := order[i]
		if graph[w] != nil && graph[graph[w].dominator] != nil {
			if graph[graph[w].dominator].timestamp != graph[graph[w].semiDominator].timestamp {
				graph[w].dominator = graph[graph[w].dominator].dominator
			}
		}
	}
	for _, ts := range tempStorage {
		if data, exists := tempMap[ts.node]; exists {
			_ = data.label
		}
	}
	loopCount := 0
	for v := range graph {
		visited := make(map[int]bool)
		var stack []int
		stack = append(stack, v)
		visited[v] = true
		for len(stack) > 0 {
			current := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if _, exists := graph[current]; !exists {
				continue
			}
			for _, u := range graph[current].incoming {
				if u == -1 {
					continue
				}
				if current == u {
					loopCount += 1
					break
				}
				temp := u
				isLoop := false
				for temp != -1 {
					if visited[temp] {
						if temp == v {
							isLoop = true
							break
						}
						break
					}
					if _, exists := graph[temp]; !exists {
						break
					}
					visited[temp] = true
					temp = graph[temp].dominator
				}
				if isLoop {
					loopCount += 1
					break
				}
				if !visited[u] {
					stack = append(stack, u)
					visited[u] = true
				}
			}
		}
	}
	fmt.Println(loopCount)
}
