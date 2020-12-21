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

var memState = make(map[int]int)

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

	scanner := bufio.NewScanner(buf)
	var base []int
	for scanner.Scan() {
		line := scanner.Text()
		nums := strings.Split(line, ",")
		for _, valStr := range nums {
			val, err := strconv.Atoi(valStr)
			if err != nil {
				fmt.Println("val value is not an integer", err)
				return base
			}
			base = append(base, val)
		}
	}
	return base
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

	base := parseFile(fp)

	valMap := make(map[int]int)
	for i, val := range base {
		valMap[val] = i
	}

	lastVal := base[len(base)-1]
	for i := len(base); i < 2020; i++ {
		idx, present := valMap[lastVal]
		age := 0
		if present {
			age = i - idx - 1
		}
		valMap[lastVal] = i - 1
		lastVal = age
	}
	fmt.Println(lastVal)
}
