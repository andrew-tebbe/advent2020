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

func parseFile(path string) (int, []int) {
	buf, err := os.Open(path)
	if err != nil {
		fmt.Println("File opening error", err)
		return 0, nil
	}

	defer func() {
		if err = buf.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	scanner := bufio.NewScanner(buf)
	lineCount := 0
	timestamp := 0
	var ids []int
	for scanner.Scan() {
		line := scanner.Text()
		if lineCount == 0 {
			timestamp, err = strconv.Atoi(line)
			if err != nil {
				fmt.Println("arg value is not an integer", err)
				return 0, []int{}
			}
			lineCount++
			continue
		}
		for _, val := range strings.Split(line, ",") {
			if val == "x" {
				continue
			}
			id, err := strconv.Atoi(val)
			if err != nil {
				fmt.Println("arg value is not an integer", err)
				return 0, []int{}
			}
			ids = append(ids, id)
		}
	}
	return timestamp, ids
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

	timestamp, ids := parseFile(fp)

	minWait := 0
	minID := 0
	for _, id := range ids {
		wait := id - timestamp%id
		if minWait == 0 || wait < minWait {
			minWait = wait
			minID = id
		}
	}
	fmt.Println(minWait * minID)
}
