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

var eyeColors = []string{
	"amb",
	"blu",
	"brn",
	"gry",
	"grn",
	"hzl",
	"oth"}

var fields = map[string]func(string) bool{
	"byr": validBirthYear,
	"iyr": validIssueYear,
	"eyr": validExpireYear,
	"hgt": validHeight,
	"hcl": validHairColor,
	"ecl": validEyeColor,
	"pid": validPid,
}

func validBirthYear(val string) bool {
	year, err := strconv.Atoi(val)
	if err != nil {
		fmt.Println("year value is not an integer", err)
		return false
	}
	return year >= 1920 && year <= 2002
}

func validIssueYear(val string) bool {
	year, err := strconv.Atoi(val)
	if err != nil {
		fmt.Println("year value is not an integer", err)
		return false
	}
	return year >= 2010 && year <= 2020
}

func validExpireYear(val string) bool {
	year, err := strconv.Atoi(val)
	if err != nil {
		fmt.Println("year value is not an integer", err)
		return false
	}
	return year >= 2020 && year <= 2030
}

func validHeight(val string) bool {
	if strings.HasSuffix(val, "cm") {
		height, err := strconv.Atoi(strings.TrimSuffix(val, "cm"))
		if err != nil {
			fmt.Println("height value is not an integer", err)
			return false
		}
		return height >= 150 && height <= 193
	}
	if strings.HasSuffix(val, "in") {
		height, err := strconv.Atoi(strings.TrimSuffix(val, "in"))
		if err != nil {
			fmt.Println("height value is not an integer", err)
			return false
		}
		return height >= 59 && height <= 76
	}
	return false
}

func validHairColor(val string) bool {
	if !strings.HasPrefix(val, "#") {
		return false
	}
	colorStr := strings.TrimPrefix(val, "#")
	if len(colorStr) != 6 {
		return false
	}
	_, err := strconv.ParseInt(colorStr, 16, 0)
	return err == nil
}

func validEyeColor(val string) bool {
	for _, color := range eyeColors {
		if color == val {
			return true
		}
	}
	return false
}

func validPid(val string) bool {
	if len(val) != 9 {
		return false
	}
	_, err := strconv.Atoi(val)
	return err == nil
}

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

	passports := parseFile(fp)

	validCount := 0
	for _, passport := range passports {
		valid := true
		for field, validator := range fields {
			value, present := passport[field]
			if !present || !validator(value) {
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
