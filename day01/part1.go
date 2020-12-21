package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	argc := len(os.Args[1:])

	if argc < 2 {
		fmt.Println("not enough args")
		return
	}

	input := os.Args[1]
	expected, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Expected value is not an integer", err)
		return
	}

	fp, err := filepath.Abs(input)
	if err != nil {
		fmt.Println("Filepath invalid", input, err)
		return
	}

	buf, err := os.Open(fp)
	if err != nil {
		fmt.Println("File opening error", err)
		return
	}

	defer func() {
		if err = buf.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	var lines []int
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("Integer convert error", err)
			return
		}
		lines = append(lines, val)
	}

	for _, a := range lines {
		for _, b := range lines {
			if a+b == expected {
				fmt.Println(a * b)
				return
			}
		}
	}

}
