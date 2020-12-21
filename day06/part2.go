package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func trimGroupAnswers(memberAnswers []string) string {
	var groupAnswers string
	for _, character := range memberAnswers[0] {
		present := true
		for j := 1; j < len(memberAnswers); j++ {
			if !strings.ContainsRune(memberAnswers[j], character) {
				present = false
				break
			}
		}
		if present {
			groupAnswers += string(character)
		}
	}
	return groupAnswers
}

func parseFile(path string) []string {
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

	var memberAnswers []string
	var totalAnswers []string
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			groupAnswers := trimGroupAnswers(memberAnswers)
			memberAnswers = memberAnswers[:0]
			totalAnswers = append(totalAnswers, groupAnswers)
			continue
		}
		memberAnswers = append(memberAnswers, line)
	}
	groupAnswers := trimGroupAnswers(memberAnswers)
	totalAnswers = append(totalAnswers, groupAnswers)
	return totalAnswers
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

	groupAnswers := parseFile(fp)

	uniqueAnswerCount := 0
	for _, group := range groupAnswers {
		uniqueAnswerCount += len(group)
	}
	fmt.Println(uniqueAnswerCount)
}
