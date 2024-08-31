package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Trans struct {
	next int
	out  string
}

type FSM struct {
	start  int
	states []int
	inputs []string
	outs   []string
	trans  [][]Trans
}

func NewFSM(start int, states []int, inputs, outs []string, transMatrix [][]int, outMatrix [][]string) *FSM {
	trans := make([][]Trans, len(states))
	for i := range states {
		trans[i] = make([]Trans, len(inputs))
		for j := range inputs {
			trans[i][j] = Trans{next: transMatrix[i][j], out: outMatrix[i][j]}
		}
	}
	return &FSM{
		start:  start,
		states: states,
		inputs: inputs,
		outs:   outs,
		trans:  trans,
	}
}

func (f *FSM) UniqueOutArray() [][]string {
	transSet := make(map[struct {
		state int
		out   string
	}]bool)
	outArray := make([][]string, len(f.states))
	for i := range f.states {
		for j := range f.inputs {
			t := f.trans[i][j]
			key := struct {
				state int
				out   string
			}{state: t.next, out: t.out}
			if !transSet[key] {
				outArray[t.next] = append(outArray[t.next], t.out)
				transSet[key] = true
			}
		}
	}
	for i := range outArray {
		sort.Strings(outArray[i])
	}
	return outArray
}

func (f *FSM) ShowGraph() {
	outArray := f.UniqueOutArray()
	nodeMap := mapNodeIDs(outArray, f.outs)
	graphBuilder := strings.Builder{}
	graphBuilder.WriteString("digraph {\n")
	graphBuilder.WriteString("\trankdir = LR\n")
	for state, outs := range outArray {
		for _, out := range outs {
			outIdx, _ := strconv.Atoi(out)
			nodeID := nodeMap[fmt.Sprintf("(%d,%s)", state, out)]
			graphBuilder.WriteString(fmt.Sprintf("\t%d [label = \"(%d,%s)\"]\n", nodeID, state, f.outs[outIdx]))
		}
	}
	edges := collectEdges(f.trans, f.inputs, outArray, nodeMap)
	for _, edge := range edges {
		graphBuilder.WriteString(edge)
	}
	graphBuilder.WriteString("}\n")
	fmt.Print(graphBuilder.String())
}

func mapNodeIDs(outArray [][]string, outs []string) map[string]int {
	nodeID := 0
	nodeMap := make(map[string]int)
	for state, outs := range outArray {
		for _, out := range outs {
			key := fmt.Sprintf("(%d,%s)", state, out)
			nodeMap[key] = nodeID
			nodeID++
		}
	}
	return nodeMap
}

func collectEdges(trans [][]Trans, inputs []string, outArray [][]string, nodeMap map[string]int) []string {
	var edges []string
	for state, outs := range outArray {
		for _, out := range outs {
			for inputIdx, t := range trans[state] {
				fromNode := nodeMap[fmt.Sprintf("(%d,%s)", state, out)]
				toNode := nodeMap[fmt.Sprintf("(%d,%s)", t.next, t.out)]
				edges = append(edges, fmt.Sprintf("\t%d -> %d [label = \"%s\"]\n", fromNode, toNode, inputs[inputIdx]))
			}
		}
	}
	return edges
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	readLine := func() string {
		line, _ := reader.ReadString('\n')
		return strings.TrimSpace(line)
	}
	parseInts := func(line string) []int {
		parts := strings.Fields(line)
		nums := make([]int, len(parts))
		for i, part := range parts {
			nums[i], _ = strconv.Atoi(part)
		}
		return nums
	}
	parseStrings := func(line string) []string {
		return strings.Fields(line)
	}
	if _, err := strconv.Atoi(readLine()); err != nil {
		return
	}
	inputSet := parseStrings(readLine())

	if _, err := strconv.Atoi(readLine()); err != nil {
		return
	}
	outSet := parseStrings(readLine())
	stateCount, err := strconv.Atoi(readLine())
	if err != nil {
		return
	}
	states := make([]int, stateCount)
	for i := range states {
		states[i] = i
	}
	transMatrix := make([][]int, stateCount)
	for i := 0; i < stateCount; i++ {
		transMatrix[i] = parseInts(readLine())
	}
	outMatrix := make([][]string, stateCount)
	for i := 0; i < stateCount; i++ {
		outMatrix[i] = parseStrings(readLine())
	}
	fsm := NewFSM(0, states, inputSet, outSet, transMatrix, outMatrix)
	fsm.ShowGraph()
}
