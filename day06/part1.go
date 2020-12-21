package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type BoardingPass struct {
	rowPattern    string
	columnPattern string
}

func getRow(searchPattern string) int {
	if len(searchPattern) != 7 {
		fmt.Println("Invalid row pattern", len(searchPattern))
	}
	min := 0
	max := 127
	for i := 0; i < len(searchPattern); i++ {
		diff := max - min + 1
		if searchPattern[i] == 'F' {
			max -= diff / 2
			continue
		}
		if searchPattern[i] == 'B' {
			min += diff / 2
			continue
		}
	}
	return min
}

func getColumn(searchPattern string) int {
	if len(searchPattern) != 3 {
		fmt.Println("Invalid column pattern", len(searchPattern))
	}
	min := 0
	max := 7
	for i := 0; i < len(searchPattern); i++ {
		diff := max - min + 1
		if searchPattern[i] == 'L' {
			max -= diff / 2
			continue
		}
		if searchPattern[i] == 'R' {
			min += diff / 2
			continue
		}
	}
	return min
}
func parseFile(path string) []map[rune]bool {
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

	var uniqueAnswers []map[rune]bool
	scanner := bufio.NewScanner(buf)
	var answers = make(map[rune]bool)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			var tmp = make(map[rune]bool)
			for k, v := range answers {
				tmp[k] = v
				delete(answers, k)
			}
			uniqueAnswers = append(uniqueAnswers, tmp)
		}
		for _, answerID := range line {
			_, present := answers[answerID]
			if !present {
				answers[answerID] = true
			}
		}
	}
	uniqueAnswers = append(uniqueAnswers, answers)
	return uniqueAnswers
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
