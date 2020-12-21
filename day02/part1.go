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

type passwordPolicy struct {
	minCount  int
	maxCount  int
	character string
}

type corruptPasswordEntry struct {
	policy   passwordPolicy
	password string
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

	var passwords []corruptPasswordEntry
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		counts := strings.Split(parts[0], "-")
		minCount, err := strconv.Atoi(counts[0])
		if err != nil {
			fmt.Println("minCount integer convert error", err)
			return
		}
		maxCount, err := strconv.Atoi(counts[1])
		if err != nil {
			fmt.Println("maxCount integer convert error", err)
			return
		}
		character := parts[1][:1]
		var policy passwordPolicy
		policy.minCount = minCount
		policy.maxCount = maxCount
		policy.character = character
		var entry corruptPasswordEntry
		entry.policy = policy
		entry.password = parts[2]
		passwords = append(passwords, entry)
	}

	validCount := 0
	for _, entry := range passwords {
		charCount := strings.Count(entry.password, entry.policy.character)
		if charCount < entry.policy.minCount || charCount > entry.policy.maxCount {
			continue
		}
		validCount++
	}
	fmt.Println(validCount)

}
