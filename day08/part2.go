package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
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

func checkTermination(program []Instruction) (bool, map[int]int) {
	traceback := make(map[int]int)
	terminated := true
	pc = 0
	acc = 0
	for instrCount := 0; pc < len(program); instrCount++ {
		_, alreadyAccessed := traceback[pc]
		if alreadyAccessed {
			terminated = false
			break
		}
		traceback[pc] = instrCount
		runInstr(program[pc])
	}
	return terminated, traceback
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

	_, badTraceback := checkTermination(program)
	for instrPC := range badTraceback {
		origOpcode := program[instrPC].opcode
		if reflect.ValueOf(origOpcode).Pointer() == reflect.ValueOf(accumulate).Pointer() {
			continue
		}
		if reflect.ValueOf(origOpcode).Pointer() == reflect.ValueOf(nop).Pointer() {
			program[instrPC].opcode = jmp
		}
		if reflect.ValueOf(origOpcode).Pointer() == reflect.ValueOf(jmp).Pointer() {
			program[instrPC].opcode = nop
		}
		terminated, _ := checkTermination(program)
		if terminated {
			break
		}
		program[instrPC].opcode = origOpcode
	}

	fmt.Println(acc)
}
