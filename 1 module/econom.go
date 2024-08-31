package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

func OperatorChecker(ch byte) bool {
	operators := map[byte]bool{'#': true, '$': true, '@': true}
	_, exists := operators[ch]
	return exists
}

func ElemCreator(str string) []string {
	var elem []string
	if len(str) == 1 {
		return elem
	}
	ind1 := 1
	for ind1 < len(str) && str[ind1] != ')' {
		if str[ind1] == ' ' {
			ind1++
			continue
		}
		if unicode.IsLetter(rune(str[ind1])) || OperatorChecker(str[ind1]) {
			elem = append(elem, str[ind1:ind1+1])
			ind1++
		} else if str[ind1] == '(' {
			var res = 0
			var depth = 1
			i := ind1 + 1
			for i < len(str) {
				switch str[i] {
				case '(':
					depth = depth + 1
				case ')':
					depth = depth - 1
				}
				if depth == 0 {
					res = i + 1
					break
				}
				i++
			}
			ind2 := res
			if ind2 == -1 {
				ind2 = len(str)
			}
			elem = append(elem, str[ind1:ind2])
			ind1 = ind2
		}
	}
	return elem
}

func isEmpty(elements []string) bool {
	return len(elements) == 0
}

func parse(input string) int {
	elements := ElemCreator(input)
	if isEmpty(elements) {
		return 0
	}
	return StackCounter(elements)
}

var Expr = map[string]bool{}

func StackCounter(elements []string) int {
	var stack []int
	for _, element := range elements {
		if len(element) > 1 {
			number := parse(element)
			Expr[element] = true
			stack = append(stack, number)
		} else if unicode.IsLetter(rune(element[0])) {
			stack = append(stack, int(element[0]-'0'))
		} else {
			if len(stack) < 2 {
				fmt.Println("Incorrect number of operands.")
				os.Exit(0)
			}
			Oper1 := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			Oper2 := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			result := 0
			switch element[0] {
			case '#':
				result = Oper1 + Oper2
			case '$':
				result = Oper1 * Oper2
			case '@':
				result = Oper1 - Oper2
			}
			stack = append(stack, result)
		}
	}
	if len(stack) != 1 {
		fmt.Println("Evaluation error.")
		os.Exit(0)
	}
	return stack[0]
}

func main() {
	var str string
	myscanner := bufio.NewScanner(os.Stdin)
	myscanner.Scan()
	str = myscanner.Text()
	runes := []rune(str)
	size := len(runes)
	i := 0
	for i < size/2 {
		runes[i], runes[size-1-i] = runes[size-1-i], runes[i]
		i++
	}
	str = string(runes)
	rune := []rune(str)
	i = 0
	for i < len(str) {
		switch rune[i] {
		case '(':
			rune[i] = ')'
		case ')':
			rune[i] = '('
		}
		i++
	}
	str = string(rune)
	if len(str) < 1 || len(str) > 1 {
		Expr[str] = true
	}
	parse(str)
	final := len(Expr)
	fmt.Println(final)
}
