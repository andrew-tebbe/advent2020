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
	index1    int
	index2    int
	character byte
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
		index1, err := strconv.Atoi(counts[0])
		if err != nil {
			fmt.Println("index1 integer convert error", err)
			return
		}
		index2, err := strconv.Atoi(counts[1])
		if err != nil {
			fmt.Println("index2 integer convert error", err)
			return
		}
		character := parts[1][0]
		var policy passwordPolicy
		policy.index1 = index1 - 1
		policy.index2 = index2 - 1
		policy.character = character
		var entry corruptPasswordEntry
		entry.policy = policy
		entry.password = parts[2]
		passwords = append(passwords, entry)
	}

	validCount := 0
	for _, entry := range passwords {
		password := entry.password
		if password[entry.policy.index1] == entry.policy.character && password[entry.policy.index2] != entry.policy.character {
			validCount++
			continue
		}
		if password[entry.policy.index1] != entry.policy.character && password[entry.policy.index2] == entry.policy.character {
			validCount++
			continue
		}

	}
	fmt.Println(validCount)

}
