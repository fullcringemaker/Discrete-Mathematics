package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func hsort(n int,
	less func(i, j int) bool,
	swap func(i, j int)) {
	var sift func(start, end int)
	sift = func(start, end int) {
		root := start
		for {
			child := 2*root + 1
			if child > end {
				break
			}
			if child+1 <= end && less(child, child+1) {
				child++
			}
			if less(root, child) {
				swap(root, child)
				root = child
			} else {
				break
			}
		}
	}
	build := func(length int) {
		for start := (length - 2) / 2; start >= 0; start-- {
			sift(start, length-1)
		}
	}
	build(n)
	for end := n - 1; end > 0; end-- {
		swap(0, end)
		sift(0, end-1)
	}
}

func getInput() ([]int, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter number of elements: ")
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	line = strings.TrimSpace(line)
	num, err := strconv.Atoi(line)
	if err != nil {
		return nil, err
	}
	arr := make([]int, num)
	fmt.Printf("Enter %d elements separated by space: ", num)
	line, err = reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	line = strings.TrimSpace(line)
	elements := strings.Split(line, " ")
	for i := 0; i < num; i++ {
		arr[i], err = strconv.Atoi(elements[i])
		if err != nil {
			return nil, err
		}
	}
	return arr, nil
}

func main() {
	arr, err := getInput()
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	n := len(arr)
	less := func(i, j int) bool {
		return arr[i] < arr[j]
	}
	swap := func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	}
	hsort(n, less, swap)
	fmt.Println(arr)
}
