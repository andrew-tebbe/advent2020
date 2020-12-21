package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

type Instruction struct {
	opcode func(int)
	arg    int
}

var pc int = 0
var acc = 0

func accumulate(arg int) {
	acc += arg
	pc++
	return
}

func nop(arg int) {
	pc++
	return
}

func jmp(arg int) {
	pc += arg
	return
}

func parseFile(path string) []Instruction {
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

	scanner := bufio.NewScanner(buf)
	fullRe := regexp.MustCompile(`(?P<opcode>\w{3}) (?P<arg>[-+]\d+)`)
	var program []Instruction
	for scanner.Scan() {
		var instr Instruction
		line := scanner.Text()
		matches := fullRe.FindStringSubmatch(line)
		opcode := matches[fullRe.SubexpIndex("opcode")]
		argStr := matches[fullRe.SubexpIndex("arg")]
		arg, err := strconv.Atoi(argStr)
		if err != nil {
			fmt.Println("arg value is not an integer", err)
			return program
		}
		instr.arg = arg
		switch opcode {
		case "acc":

			instr.opcode = accumulate
		case "jmp":
			instr.opcode = jmp
		case "nop":
			fallthrough
		default:
			instr.opcode = nop
		}
		program = append(program, instr)
	}

	return program
}

func runInstr(instr Instruction) {
	instr.opcode(instr.arg)
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

	program := parseFile(fp)

	traceback := make(map[int]int)
	for instrCount := 0; pc < len(program); instrCount++ {
		_, alreadyAccessed := traceback[pc]
		if alreadyAccessed {
			break
		}
		traceback[pc] = instrCount
		runInstr(program[pc])
	}

	fmt.Println(acc)
}
