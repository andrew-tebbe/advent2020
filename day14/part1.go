package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	Zero int = iota
	One
	DontCase
)

type InitOperation struct {
	address  int
	setValue int
}

type InitSection struct {
	mask       []int
	operations []InitOperation
}

var memState = make(map[int]int)

func parseFile(path string) []InitSection {
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
	var initialization []InitSection
	var section InitSection
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "mask") {
			if len(section.operations) != 0 {
				initialization = append(initialization, section)
			}
			maskStr := strings.Split(line, " = ")[1]
			var mask []int
			for _, char := range maskStr {
				switch char {
				case 'X':
					mask = append(mask, DontCase)
				case '1':
					mask = append(mask, One)
				case '0':
					fallthrough
				default:
					mask = append(mask, Zero)
				}
			}
			section = InitSection{
				mask,
				[]InitOperation{},
			}
			continue
		}

		var operation InitOperation
		fullRe := regexp.MustCompile(`mem\[(?P<addr>\d+)\] = (?P<val>\d+)`)
		matches := fullRe.FindStringSubmatch(line)
		addr, err := strconv.Atoi(matches[fullRe.SubexpIndex("addr")])
		if err != nil {
			fmt.Println("arg value is not an integer", err)
			return initialization
		}
		operation.address = addr
		val, err := strconv.Atoi(matches[fullRe.SubexpIndex("val")])
		if err != nil {
			fmt.Println("arg value is not an integer", err)
			return initialization
		}
		operation.setValue = val
		section.operations = append(section.operations, operation)
	}
	initialization = append(initialization, section)
	return initialization
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

	initialization := parseFile(fp)

	for _, section := range initialization {
		mask := section.mask
		for _, operation := range section.operations {
			memState[operation.address] = applyMask(operation.setValue, mask)
		}
	}
	sum := 0
	for _, val := range memState {
		sum += val
	}
	fmt.Println(sum)
	// fmt.Println(memState)
	// fmt.Println(initialization)
}

func applyMask(memVal int, mask []int) int {
	val := 0
	for i := len(mask) - 1; i >= 0; i-- {
		bitPos := 1 << i
		idxVal := memVal & bitPos
		if mask[len(mask)-i-1] == Zero || mask[len(mask)-i-1] == One {
			idxVal = mask[len(mask)-i-1] << i
		}
		val |= idxVal
		bitPos = bitPos >> 1
	}
	return val
}
