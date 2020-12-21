package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type BoardingPass struct {
	rowPattern    string
	columnPattern string
}

func getRow(searchPattern string) int {
	if len(searchPattern) != 7 {
		fmt.Println("Invalid row pattern", len(searchPattern))
	}
	min := 0
	max := 127
	for i := 0; i < len(searchPattern); i++ {
		diff := max - min + 1
		if searchPattern[i] == 'F' {
			max -= diff / 2
			continue
		}
		if searchPattern[i] == 'B' {
			min += diff / 2
			continue
		}
	}
	return min
}

func getColumn(searchPattern string) int {
	if len(searchPattern) != 3 {
		fmt.Println("Invalid column pattern", len(searchPattern))
	}
	min := 0
	max := 7
	for i := 0; i < len(searchPattern); i++ {
		diff := max - min + 1
		if searchPattern[i] == 'L' {
			max -= diff / 2
			continue
		}
		if searchPattern[i] == 'R' {
			min += diff / 2
			continue
		}
	}
	return min
}
func parseFile(path string) []BoardingPass {
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

	var boardingPasses []BoardingPass
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		line := scanner.Text()
		pass := BoardingPass{
			line[:7],
			line[7:],
		}
		boardingPasses = append(boardingPasses, pass)
	}
	return boardingPasses
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

	passes := parseFile(fp)

	var ids = make(map[int]bool)
	for _, pass := range passes {
		row := getRow(pass.rowPattern)
		column := getColumn(pass.columnPattern)
		id := row*8 + column
		ids[id] = true
	}
	maxID := 127*8 + 8
	for i := 1; i < maxID; i++ {
		_, left := ids[i-1]
		if !left {
			continue
		}
		_, right := ids[i+1]
		if !right {
			continue
		}
		_, here := ids[i]
		if !here {
			fmt.Println(i)
			break
		}
	}
}
