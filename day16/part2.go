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

	scanErrorRate := 0
	var validTickets [][]int
	for _, ticket := range otherTickets {
		ticketValid := true
		for _, val := range ticket {
			valid := false
			for _, ranges := range rules {
				if valid {
					break
				}
				for _, rangeInst := range ranges {
					if val >= rangeInst.start && val <= rangeInst.stop {
						valid = true
						break
					}
				}
			}
			if !valid {
				scanErrorRate += val
				ticketValid = false
				break
			}
		}
		if ticketValid {
			validTickets = append(validTickets, ticket)
		}
	}

	fmt.Println(len(otherTickets), len(validTickets))
	possibleIdx := make(map[int]bool)
	for i := 0; i < len(validTickets[0]); i++ {
		possibleIdx[i] = true
	}

	ruleMap := make(map[string]map[int]bool)
	for ruleKey, rule := range rules {
		validIdx := make(map[int]bool)
		for idx := range possibleIdx {
			valid := true
			for _, ticket := range validTickets {
				inRange := false
				for _, rangeInst := range rule {
					if ticket[idx] >= rangeInst.start && ticket[idx] <= rangeInst.stop {
						inRange = true
						// fmt.Println(ticket[idx], "<", rangeInst.start, "||", ticket[idx], ">", rangeInst.stop, valid)
						break
					}
					// fmt.Println(idx, valid, ticket[idx], ruleKey)
				}
				if !inRange {
					valid = false
					break
				}
			}
			if valid {
				validIdx[idx] = true
				// possibleIdx = remove(possibleIdx, i)
				// break
			}
		}
		ruleMap[ruleKey] = validIdx
	}

	for allSimplified := false; !allSimplified; {
		allSimplified = true
		keepName := ""
		eliminateIdx := -1
		for name, validIdx := range ruleMap {
			if len(validIdx) == 1 {
				for idx := range validIdx {
					if possibleIdx[idx] {
						eliminateIdx = idx
						keepName = name
						possibleIdx[idx] = false
						allSimplified = false
					}
				}
				if !allSimplified {
					break
				}
			}
		}
		if !allSimplified {
			for ruleKey := range ruleMap {
				if ruleKey == keepName {
					continue
				}
				_, present := ruleMap[ruleKey][eliminateIdx]
				if present {
					delete(ruleMap[ruleKey], eliminateIdx)
				}
			}
		}
	}

	idxMapping := make(map[string]int)
	for name, idxs := range ruleMap {
		for i := range idxs {
			idxMapping[name] = i
		}
	}
	fmt.Println(idxMapping)
	ticketTotal := 1
	for name, idx := range idxMapping {
		if strings.HasPrefix(name, "departure") {
			ticketTotal *= myTicket[idx]
		}
	}
	fmt.Println(ticketTotal)
}

func remove(arr []int, idx int) []int {
	arr[idx] = arr[len(arr)-1]
	return arr[:len(arr)-1]
}
