package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
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

	scanner := bufio.NewScanner(buf)
	var adapters []int
	for scanner.Scan() {
		value, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("value is not an integer", err)
			return adapters
		}
		adapters = append(adapters, value)
	}
	sort.Ints(adapters)

	return adapters
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

	adapters := parseFile(fp)
	adapters = append([]int{0}, adapters...)

	pathCount := make(map[int]int)
	for i, val := range adapters {
		if val == 0 {
			pathCount[i] = 1
			continue
		}

		sum := 0
		for j := i - 1; j >= i-3 && j >= 0; j-- {
			if validNext(adapters[j], val) {
				sum += pathCount[j]
				continue
			}
			break
		}
		pathCount[i] = sum
	}
	fmt.Println(pathCount[len(pathCount)-1])
}

func validNext(base int, next int) bool {
	if next-base > 3 {
		return false
	}
	return true
}
