package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func parseAutomaton() (int, int, int, [][]int, [][]string, int) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	readInt := func() int {
		scanner.Scan()
		var num int
		fmt.Sscan(scanner.Text(), &num)
		return num
	}
	readString := func() string {
		scanner.Scan()
		return scanner.Text()
	}
	stateCount := readInt()
	transTable := make([][]int, stateCount)
	for state := 0; state < stateCount; state++ {
		transTable[state] = make([]int, 2)
		for input := 0; input < 2; input++ {
			transTable[state][input] = readInt()
		}
	}
	outputTable := make([][]string, stateCount)
	for state := 0; state < stateCount; state++ {
		outputTable[state] = make([]string, 2)
		for input := 0; input < 2; input++ {
			outputTable[state][input] = readString()
		}
	}
	initialState := readInt()
	maxWordLength := readInt()
	return stateCount, 2, initialState, transTable, outputTable, maxWordLength
}

func generateLanguage(stateCount, inputCount, initialState int, 
	transTable [][]int, outputTable [][]string, maxWordLength int) map[string]struct{} {
	wordSet := make(map[string]struct{})
	type stateWordPair struct {
		state int
		word  string
	}
	visited := make(map[stateWordPair]bool)
	queue := []stateWordPair{{initialState, ""}}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if visited[current] {
			continue
		}
		visited[current] = true
		wordSet[current.word] = struct{}{}
		if len(current.word) >= maxWordLength {
			continue
		}
		for i := 0; i < inputCount; i++ {
			nextState := transTable[current.state][i]
			outputSymbol := outputTable[current.state][i]
			if nextState == current.state && outputSymbol == "-" {
				continue
			}
			newWord := current.word
			if outputSymbol != "-" {
				newWord += outputSymbol
			}
			nextPair := stateWordPair{nextState, newWord}
			queue = append(queue, nextPair)
		}
	}
	return wordSet
}

func sortLanguage(wordSet map[string]struct{}) []string {
	wordList := make([]string, 0, len(wordSet))
	for word := range wordSet {
		wordList = append(wordList, word)
	}
	sort.Slice(wordList, func(i, j int) bool {
		if len(wordList[i]) == len(wordList[j]) {
			return wordList[i] < wordList[j]
		}
		return len(wordList[i]) < len(wordList[j])
	})
	return wordList
}

func main() {
	stateCount, inputCount, initialState, transTable, outputTable, maxWordLength := parseAutomaton()
	wordSet := generateLanguage(stateCount, inputCount, initialState, 
		transTable, outputTable, maxWordLength)
	sortedWords := sortLanguage(wordSet)
	previousLength := -1
	for _, word := range sortedWords {
		if len(word) != previousLength {
			if previousLength != -1 {
				fmt.Println()
			}
			previousLength = len(word)
		} else {
			fmt.Print(" ")
		}
		fmt.Print(word)
	}
	fmt.Println()
}
