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

type Range struct {
	start int
	stop  int
}

func parseFile(path string) (map[string][]Range, []int, [][]int) {
	buf, err := os.Open(path)
	if err != nil {
		fmt.Println("File opening error", err)
		return map[string][]Range{}, []int{}, [][]int{}
	}

	defer func() {
		if err = buf.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	var myTicket []int
	var otherTickets [][]int
	scanner := bufio.NewScanner(buf)
	rules := make(map[string][]Range)
	for scanner.Scan() {
		var ranges []Range
		line := scanner.Text()
		if len(line) == 0 {
			break
		}
		sides := strings.Split(line, ": ")
		for _, rangeStr := range strings.Split(sides[1], " or ") {
			rangeArr := strings.Split(rangeStr, "-")
			start, err1 := strconv.Atoi(rangeArr[0])
			stop, err2 := strconv.Atoi(rangeArr[1])
			if err1 != nil || err2 != nil {
				fmt.Println("val value is not an integer", err1, err2)
				return rules, []int{}, [][]int{}
			}
			ranges = append(ranges, Range{start, stop})
		}
		rules[sides[0]] = ranges
	}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "your ticket:" {
			continue
		}
		if len(line) == 0 {
			break
		}

		for _, valStr := range strings.Split(line, ",") {
			val, err := strconv.Atoi(valStr)
			if err != nil {
				fmt.Println("val value is not an integer", err)
				return rules, []int{}, [][]int{}
			}
			myTicket = append(myTicket, val)
		}
	}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "nearby tickets:" {
			continue
		}
		if len(line) == 0 {
			break
		}

		var otherTicket []int
		for _, valStr := range strings.Split(line, ",") {
			val, err := strconv.Atoi(valStr)
			if err != nil {
				fmt.Println("val value is not an integer", err)
				return rules, []int{}, [][]int{}
			}
			otherTicket = append(otherTicket, val)
		}
		otherTickets = append(otherTickets, otherTicket)

	}
	return rules, myTicket, otherTickets
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

	rules, myTicket, otherTickets := parseFile(fp)

	// lastVal := base[len(base)-1]
	// for i := len(base); i < 2020; i++ {
	// 	idx, present := valMap[lastVal]
	// 	age := 0
	// 	if present {
	// 		age = i - idx - 1
	// 	}
	// 	valMap[lastVal] = i - 1
	// 	lastVal = age
	// }
	scanErrorRate := 0
	for _, ticket := range otherTickets {
		for _, val := range ticket {
			valid := false
			for _, ranges := range rules {
				if valid {
					break
				}
				for _, rangeInst := range ranges {
					if val >= rangeInst.start && val < rangeInst.stop {
						valid = true
						break
					}
				}
			}
			if !valid {
				scanErrorRate += val
			}
		}
	}
	fmt.Println(rules, myTicket, otherTickets)
	fmt.Println(scanErrorRate)
}
