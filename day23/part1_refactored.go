package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func parseFile(path string) []int {
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

	var cups []int
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		for _, cupStr := range scanner.Text() {
			cup, err := strconv.Atoi(string(cupStr))
			if err != nil {
				fmt.Println("cup number is not an integer", err)
				return cups
			}
			cups = append(cups, cup)
		}
	}
	return cups
}

func main() {
	argc := len(os.Args[1:])

	if argc < 2 {
		fmt.Println("not enough args")
		return
	}

	input := os.Args[1]

	moveCount, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("move count is not an integer", err)
		return
	}

	fp, err := filepath.Abs(input)
	if err != nil {
		fmt.Println("Filepath invalid", input, err)
		return
	}

	cups := parseFile(fp)
	// cup -> nextCup
	cupMap := make(map[int]int)
	currentCup := cups[0]
	i := 1
	for ; i < len(cups); i++ {
		cupMap[cups[i-1]] = cups[i]
	}
	cupMap[cups[i-1]] = currentCup
	for i := 0; i < moveCount; i++ {
		performMove(cupMap, currentCup)
		currentCup = cupMap[currentCup]
	}

	cupLabel := 1
	for i := 0; i < len(cups)-1; i++ {
		cupLabel = cupMap[cupLabel]
		fmt.Printf("%d", cupLabel)
	}
	fmt.Printf("\n")
}

func performMove(cups map[int]int, currentCup int) {
	maxCupVal := len(cups)
	destinationLabel := cups[currentCup] - 1
	selectedCupStart := cups[currentCup]
	selectedCupEnd := currentCup
	for i := 0; i < 3; i++ {
		selectedCupEnd = cups[selectedCupEnd]
	}
	cups[currentCup] = cups[selectedCupEnd]
	cups[selectedCupEnd] = 0
	destinationLabel = currentCup - 1
	for destinationLabel < 1 || destinationLabel == selectedCupStart || destinationLabel == selectedCupEnd || destinationLabel == cups[selectedCupStart] {
		destinationLabel--
		if destinationLabel < 1 {
			destinationLabel = maxCupVal
		}
	}
	cups[selectedCupEnd] = cups[destinationLabel]
	cups[destinationLabel] = selectedCupStart
}
