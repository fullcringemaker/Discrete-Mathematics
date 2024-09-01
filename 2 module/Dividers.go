package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

type FactorList []int

func calculateDivisors(n int) FactorList {
	limit := int(math.Sqrt(float64(n)))
	divisors := FactorList{}
	var i = 1
	for i <= limit {
		if n%i == 0 {
			divisors = append(divisors, i)
			if i != n/i {
				divisors = append(divisors, n/i)
			}
		}
		i++
	}
	sort.Slice(divisors, func(i, j int) bool {
		return divisors[i] > divisors[j]
	})
	return divisors
}
func drawGraph(divisors FactorList) {
	fmt.Println("graph {")
	for _, divisor := range divisors {
		fmt.Printf("    %d\n", divisor)
	}
	linkDivisors(divisors)
	fmt.Println("}")
}
func linkDivisors(divisors FactorList) {
	size := len(divisors)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if i != j && divisors[i]%divisors[j] == 0 {
				if !isTransitive(divisors, i, j) {
					fmt.Printf("    %d -- %d\n", divisors[i], divisors[j])
				}
			}
		}
	}
}
func isTransitive(divisors FactorList, i, j int) bool {
	var k = 0
	for k < len(divisors) {
		if k != i && k != j && divisors[i]%divisors[k] == 0 && divisors[k]%divisors[j] == 0 {
			return true
		}
		k++
	}
	return false
}
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	n, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Ошибка ввода:", err)
		return
	}
	divisors := calculateDivisors(n)
	drawGraph(divisors)
}
