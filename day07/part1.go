package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func parseFile(path string) map[string]map[string]int {
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

	rules := make(map[string]map[string]int)
	scanner := bufio.NewScanner(buf)
	fullRe := regexp.MustCompile(`(?P<key>[\w ]+) bags contain (?P<rule>[\d\w ,]+)\.`)
	ruleRe := regexp.MustCompile(`(?P<count>\d+) (?P<type>[\w ]+) bags*`)
	noBags := "no other bags"
	for scanner.Scan() {
		line := scanner.Text()
		contents := make(map[string]int)
		matches := fullRe.FindStringSubmatch(line)
		key := matches[fullRe.SubexpIndex("key")]
		rule := matches[fullRe.SubexpIndex("rule")]
		if rule != noBags {
			for _, loc := range ruleRe.FindAllSubmatchIndex([]byte(rule), -1) {
				count, err := strconv.Atoi(rule[loc[2]:loc[3]])
				if err != nil {
					fmt.Println("count value is not an integer", err)
					return rules
				}
				contents[rule[loc[4]:loc[5]]] = count
			}
		}
		rules[key] = contents
	}

	return rules
}

func getContainingBags(bagRef string, rules map[string]map[string]int, outerBags map[string]string) map[string]string {
	var validOuterBags []string
	for bag, rule := range rules {
		_, present := rule[bagRef]
		if present {
			validOuterBags = append(validOuterBags, bag)
		}
	}
	newOuterBags := make(map[string]string)
	for _, bagType := range validOuterBags {
		_, present := outerBags[bagType]
		if present {
			continue
		}
		newOuterBags[bagType] = bagRef
		outerBags[bagType] = bagRef
	}
	for key := range newOuterBags {
		outerBags = getContainingBags(key, rules, outerBags)
	}
	return outerBags
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

	refBag := strings.Join(os.Args[2:], " ")

	bagRules := parseFile(fp)

	validOuterBags := make(map[string]string)
	validOuterBags = getContainingBags(refBag, bagRules, validOuterBags)
	fmt.Println(len(validOuterBags))
}
