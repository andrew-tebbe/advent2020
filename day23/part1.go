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
	for i := 0; i < moveCount; i++ {
		currentCup := i % len(cups)
		cups = performMove(cups, currentCup)
	}
	foundStart := false
	for i := 0; i < len(cups); i = (i + 1) % len(cups) {
		if cups[i] == 1 {
			if foundStart {
				break
			}
			foundStart = true
			continue
		}
		if foundStart {
			fmt.Printf("%d", cups[i])
		}
	}
	fmt.Printf("\n")
}

func performMove(cups []int, currentCup int) []int {
	maxCupVal := len(cups)
	destinationLabel := cups[currentCup] - 1
	selectedCups, remainingCups, reorderCount := splitCups(cups, currentCup)
	for validDest := false; !validDest; {
		validDest = true
		if destinationLabel < 1 {
			destinationLabel = maxCupVal
		}
		for _, cupVal := range selectedCups {
			if cupVal == destinationLabel {
				destinationLabel--
				validDest = false
				break
			}
		}
	}
	var destinationIdx int
	insertBefore := 3
	for idx, label := range remainingCups {
		if label == cups[currentCup] {
			insertBefore = 0
		}
		if label == destinationLabel {
			destinationIdx = idx
			break
		}
	}
	reorderCount += insertBefore
	var resultingCups []int
	resultingCups = append(resultingCups, remainingCups[:destinationIdx+1]...)
	resultingCups = append(resultingCups, selectedCups...)

	if destinationIdx < len(remainingCups)-1 {
		resultingCups = append(resultingCups, remainingCups[destinationIdx+1:]...)
	}
	if reorderCount > 0 {
		resultingCups = append(resultingCups[reorderCount:], resultingCups[:reorderCount]...)
	}
	if reorderCount < 0 {
		resultingCups = append(resultingCups[len(cups)+reorderCount:], resultingCups[:len(cups)+reorderCount]...)
	}

	return resultingCups
}

func splitCups(cups []int, currentCup int) ([]int, []int, int) {
	var selectedCups []int
	var remainingCups []int
	reordersNeeded := 0
	endIdx := currentCup + 3
	if endIdx > len(cups)-1 {
		endIdx = endIdx % len(cups)
		reordersNeeded = (endIdx + 1) * -1
		selectedCups = append(cups[currentCup+1:], cups[:endIdx+1]...)
		remainingCups = append(remainingCups, cups[endIdx+1:currentCup+1]...)
	} else {
		selectedCups = cups[currentCup+1 : endIdx+1]
		remainingCups = append(remainingCups, cups[:currentCup+1]...)
		if endIdx < len(cups)-1 {
			remainingCups = append(remainingCups, cups[endIdx+1:]...)
		}
	}
	return selectedCups, remainingCups, reordersNeeded
}
