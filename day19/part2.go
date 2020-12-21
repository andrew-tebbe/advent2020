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

type BaseRule struct {
	character string
}

func (rule BaseRule) valid(message string, rules map[int]Rule) (bool, int) {
	// fmt.Println(rule)
	return strings.HasPrefix(message, rule.character), 1
}

type RecursiveRule struct {
	subRules [][]int
}

func (rule RecursiveRule) valid(message string, rules map[int]Rule) (bool, int) {
	valid := true
	inc := 0
	for _, ruleChain := range rule.subRules {
		inc = 0
		valid = true
		msg := message
		for j, ruleInst := range ruleChain {
			ruleValid, i := rules[ruleInst].valid(msg, rules)
			valid = valid && ruleValid
			if !valid {
				break
			}
			msg = msg[i:]
			inc += i
			if inc == len(message) {
				if j == 0 {
					valid = false
				}
				break
			}
		}
		if valid {
			break
		}
	}
	return valid, inc
}

type Rule interface {
	valid(string, map[int]Rule) (bool, int)
}

func parseFile(path string) (map[int]Rule, []string) {
	buf, err := os.Open(path)
	if err != nil {
		fmt.Println("File opening error", err)
		return nil, nil
	}

	defer func() {
		if err = buf.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	rules := make(map[int]Rule)
	var messages []string
	rulesComplete := false
	fullRe := regexp.MustCompile(`(?P<id>\d+): (?P<rule>[\d |"\w]+)`)
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			rulesComplete = true
			continue
		}
		if !rulesComplete {
			matches := fullRe.FindStringSubmatch(line)
			if matches != nil {
				idStr := matches[fullRe.SubexpIndex("id")]
				id, err := strconv.Atoi(idStr)
				if err != nil {
					fmt.Println("id value is not an integer", err)
					return rules, messages
				}
				ruleStr := matches[fullRe.SubexpIndex("rule")]
				if strings.Contains(ruleStr, "\"") {
					rules[id] = BaseRule{
						strings.Split(ruleStr, "\"")[1],
					}
				} else {
					var subRules [][]int
					for _, subRuleStr := range strings.Split(ruleStr, " | ") {
						var subRuleNums []int
						for _, subRuleNumStr := range strings.Split(subRuleStr, " ") {
							subRule, err := strconv.Atoi(subRuleNumStr)
							if err != nil {
								fmt.Println("subRule value is not an integer", err)
								return rules, messages
							}
							subRuleNums = append(subRuleNums, subRule)
						}
						subRules = append(subRules, subRuleNums)
					}
					rules[id] = RecursiveRule{
						subRules,
					}
				}
			}
		} else {
			messages = append(messages, line)
		}
	}
	return rules, messages
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

	rules, messages := parseFile(fp)

	// Overwrite rules
	rules[8] = RecursiveRule{
		[][]int{
			[]int{
				42,
			},
			[]int{
				42, 8,
			},
		},
	}
	rules[11] = RecursiveRule{
		[][]int{
			[]int{
				42, 31,
			},
			[]int{
				42, 11, 31,
			},
		},
	}

	validMessages := 0
	for _, message := range messages {
		valid, i := rules[0].valid(message, rules)
		if valid && i == len(message) {
			validMessages++
		}
	}
	fmt.Println(validMessages)
}
