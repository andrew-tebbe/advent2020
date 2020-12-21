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

	lastJoltage := 0
	diff1 := 0
	diff3 := 0
	for _, joltage := range adapters {
		if joltage-lastJoltage == 1 {
			diff1++
		}
		if joltage-lastJoltage == 3 {
			diff3++
		}
		lastJoltage = joltage
	}
	diff3++
	fmt.Println(diff1 * diff3)
}
