package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func AnswerReceiver(expr []string) (int, error) {
	for i := 0; i < len(expr); {
		if expr[i] == ")" {
			if i < 3 { // Проверка на достаточность элементов для операции
				return 0, fmt.Errorf("invalid expression format")
			}
			num1, err := strconv.Atoi(expr[i-2])
			if err != nil {
				return 0, err
			}
			num2, err := strconv.Atoi(expr[i-1])
			if err != nil {
				return 0, err
			}
			var outcome int
			switch expr[i-3] {
			case "+":
				outcome = num1 + num2
			case "-":
				outcome = num1 - num2
			case "*":
				outcome = num1 * num2
			default:
				return 0, fmt.Errorf("unknown operator %s", expr[i-3])
			}
			expr[i] = strconv.Itoa(outcome)
			expr = append(expr[:i-4], expr[i:]...)
			i = 0
		} else {
			i++
		}
	}
	if len(expr) != 1 {
		return 0, fmt.Errorf("")
	}
	return strconv.Atoi(expr[0])
}

func ElemCreator(str string) []string {
	var elem []string
	for _, runeValue := range str {
		if !unicode.IsSpace(runeValue) {
			elem = append(elem, string(runeValue))
		}
	}
	return elem
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	expr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	expr = strings.TrimSuffix(expr, "\n")
	elements := ElemCreator(expr)
	result, err := AnswerReceiver(elements)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Println(result)
}
