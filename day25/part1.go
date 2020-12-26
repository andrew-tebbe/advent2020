package main

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
)

func parseFile(path string) []int64 {
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

	var pubKeys []int64
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("Public key is not an integer", err)
			return pubKeys
		}
		pubKeys = append(pubKeys, int64(val))
	}
	return pubKeys
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

	subjectVal, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("subject value is not an integer", err)
		return
	}

	pubKeys := parseFile(fp)
	encKey := int64(0)
	loopVal := []int64{0, 0}
	val := int64(1)
	for i := int64(1); encKey == 0; i++ {
		val = calcCycle(int64(subjectVal), i)
		for j, key := range pubKeys {
			if val == key {
				loopVal[j] = i
				if j == 0 {
					encKey = calcCycle(pubKeys[1], i)
				} else {
					encKey = calcCycle(pubKeys[0], i)
				}
				break
			}
		}
	}

	fmt.Println(encKey)
}

func calcCycle(sub int64, loop int64) int64 {
	var val big.Int
	var mod big.Int
	var loopInt big.Int
	val.SetInt64(sub)
	mod.SetInt64(20201227)
	loopInt.SetInt64(loop)
	val.Exp(&val, &loopInt, &mod)
	return val.Int64()
}
