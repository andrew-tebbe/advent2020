package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func parseFile(path string) []map[string]string {
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

	var passports []map[string]string
	scanner := bufio.NewScanner(buf)
	var passport = make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			var tmp = make(map[string]string)
			for k, v := range passport {
				tmp[k] = v
				delete(passport, k)
			}
			passports = append(passports, tmp)
		}
		for _, data := range strings.Fields(line) {
			pair := strings.Split(data, ":")
			passport[pair[0]] = pair[1]
		}
	}
	passports = append(passports, passport)
	return passports
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

	requiredFields := strings.Split(os.Args[2], ",")

	passports := parseFile(fp)

	validCount := 0
	for _, passport := range passports {
		valid := true
		for _, field := range requiredFields {
			_, present := passport[field]
			if !present {
				valid = false
				break
			}
		}
		if valid {
			validCount++
		}
	}
	fmt.Println(validCount)
}
