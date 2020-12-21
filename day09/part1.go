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

	scanner := bufio.NewScanner(buf)
	var code []int
	for scanner.Scan() {
		value, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("value is not an integer", err)
			return code
		}
		code = append(code, value)
	}

	return code
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

	code := parseFile(fp)

	preambleLen := 25
	for i := preambleLen; i < len(code); i++ {
		valid := false
		for j := i - preambleLen; j < i && !valid; j++ {
			for k := i - preambleLen; k < i; k++ {
				if code[i] == code[j]+code[k] && j != k {
					valid = true
					break
				}
			}
		}
		if !valid {
			fmt.Println(code[i])
			break
		}
	}
}
