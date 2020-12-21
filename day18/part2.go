package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type OpPair struct {
	operator string
	operand  Computable
}

type Int struct {
	val int64
}

func (i Int) compute() Int {
	return i
}

func (i Int) add(operand Computable) Int {
	// fmt.Println(i.val, "+", operand.compute().val)
	return Int{i.val + operand.compute().val}
}

func (i Int) mult(operand Computable) Int {
	// fmt.Println(i.val, "*", operand.compute().val)
	return Int{i.val * operand.compute().val}
}

type Expression struct {
	base     Computable
	operands []OpPair
}

func (exp Expression) add(operand Computable) Int {
	// fmt.Println(i.val, "+", operand.compute().val)
	return exp.base.compute().add(operand.compute())
}

func (exp Expression) mult(operand Computable) Int {
	// fmt.Println(i.val, "*", operand.compute().val)
	return exp.base.compute().mult(operand.compute())
}

func (exp Expression) compute() Int {
	var value Int

	if len(exp.operands) == 0 {
		return exp.compute()
	}

	value = exp.base.compute()
	var simplifiedExpression []Computable
	for _, rho := range exp.operands {
		if rho.operator == "+" {
			// fmt.Println(value, rho.operand)
			value = value.add(rho.operand)
			// fmt.Println(value)
		} else if rho.operator == "*" {
			simplifiedExpression = append(simplifiedExpression, value)
			value = rho.operand.compute()
		}
	}
	simplifiedExpression = append(simplifiedExpression, value)
	// fmt.Println(simplifiedExpression)

	value = Int{1}
	for _, operand := range simplifiedExpression {
		value = value.mult(operand)
	}
	// fmt.Println(value, simplifiedExpression, exp)
	return value

}

type Computable interface {
	add(Computable) Int
	mult(Computable) Int
	compute() Int
}

func parseFile(path string) []Computable {
	buf, err := os.Open(path)
	if err != nil {
		fmt.Println("File opening error", err)
		return nil
	}

	defer func() {
		if err = buf.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	var expressions []Computable
	scanner := bufio.NewScanner(buf)

	for scanner.Scan() {
		expressions = append(expressions, parseExpression(scanner.Text()))
	}
	return expressions
}

func main() {
	argc := len(os.Args[1:])

	if argc < 1 {
		fmt.Println("not enough args")
		return
	}

	input := os.Args[1]

	fp, err := filepath.Abs(input)
	if err != nil {
		fmt.Println("Filepath invalid", input, err)
		return
	}

	expressions := parseFile(fp)

	var sum int64 = 0
	for _, expression := range expressions {
		value := expression.compute().val
		fmt.Println(value)
		sum += value
	}
	fmt.Println(sum)
}

func parseExpression(str string) Computable {
	base, i := getNextOperand(str)
	if i == len(str) {
		return base
	}

	var exp Expression
	exp.base = base

	for j := 0; i < len(str); i += j {
		var operand Computable
		splitRes := strings.SplitN(str[i:], " ", 3)
		operator := splitRes[1]
		operand, j = getNextOperand(splitRes[2])
		j += len(splitRes[0]) + len(splitRes[1]) + 2
		exp.operands = append(exp.operands, OpPair{operator, operand})
	}

	return exp
}

func getNextOperand(expStr string) (Computable, int) {
	char := rune(expStr[0])
	if char == '(' {
		rightParenIdx := findRightParen(expStr[1:]) + 1

		return parseExpression(expStr[1:rightParenIdx]), rightParenIdx + 1
	}

	nextSpace := strings.Index(expStr, " ")
	retIdx := nextSpace
	var valStr string
	if nextSpace != -1 {
		valStr = expStr[:nextSpace]
	} else {
		retIdx = len(expStr)
		valStr = expStr
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		fmt.Println("operand value is not an integer", err, 0)
		return nil, 0
	}

	return Int{int64(val)}, retIdx
}

func findRightParen(str string) int {
	leftParenIdx := strings.Index(str, "(")
	rightParenIdx := strings.Index(str, ")")
	for leftParenIdx != -1 && leftParenIdx < rightParenIdx {
		leftParenIdx = strings.Index(str[rightParenIdx+1:], "(") + rightParenIdx + 1
		rightParenIdx = strings.Index(str[rightParenIdx+1:], ")") + rightParenIdx + 1
	}
	return rightParenIdx
}
