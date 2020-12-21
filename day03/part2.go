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

type Slope struct {
	rise int
	run  int
}

func parseFile(path string) [][]bool {
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

	var grid [][]bool
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		var line []bool
		for _, space := range scanner.Text() {
			tree := true
			if space == '.' {
				tree = false
			}
			if space == '#' {
				tree = true
			}
			line = append(line, tree)
		}
		grid = append(grid, line)
	}
	return grid
}

func main() {
	argc := len(os.Args[1:])

	if argc < 2 {
		fmt.Println("not enough args")
		return
	}

	input := os.Args[1]

	fp, err := filepath.Abs(input)
	if err != nil {
		fmt.Println("Filepath invalid", input, err)
		return
	}

	var slopes []Slope
	for _, arg := range os.Args[2:] {
		var slope Slope
		invSlopeStr := strings.Split(arg, ",")
		rise, err := strconv.Atoi(invSlopeStr[0])
		if err != nil {
			fmt.Println("slope value is not an integer", err)
			return
		}
		slope.rise = rise
		run, err := strconv.Atoi(invSlopeStr[1])
		if err != nil {
			fmt.Println("slope value is not an integer", err)
			return
		}
		slope.run = run
		slopes = append(slopes, slope)
	}

	treeMap := parseFile(fp)
	treeMultiple := 1
	for _, slope := range slopes {
		treeCount := 0
		index := slope.run
		for i := slope.rise; i < len(treeMap); i += slope.rise {
			if treeMap[i][index] {
				treeCount++
			}
			index = (index + slope.run) % len(treeMap[i])
		}
		treeMultiple = treeMultiple * treeCount
	}
	fmt.Println(treeMultiple)
}
