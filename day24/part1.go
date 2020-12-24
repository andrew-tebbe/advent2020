package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type GridSpace struct {
	northSouth int
	eastWest   float32
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

	var tiles []string
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		tiles = append(tiles, scanner.Text())
	}
	return tiles
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

	tiles := parseFile(fp)
	tileMap := make(map[GridSpace]bool)
	for _, tile := range tiles {
		var space GridSpace
		eastCount := float32(0)
		westCount := float32(0)
		j := 0
		for i := strings.Index(tile, "n"); i != -1; i = strings.Index(tile[j:], "n") {
			j += i + 1
			space.northSouth++
			if tile[j] == 'e' {
				eastCount -= 0.5
			} else if tile[j] == 'w' {
				westCount -= 0.5
			}
		}

		j = 0
		for i := strings.Index(tile, "s"); i != -1; i = strings.Index(tile[j:], "s") {
			j += i + 1
			space.northSouth--
			if tile[j] == 'e' {
				eastCount -= 0.5
			} else if tile[j] == 'w' {
				westCount -= 0.5
			}
		}

		eastCount += float32(strings.Count(tile, "e"))
		westCount += float32(strings.Count(tile, "w"))
		space.eastWest = eastCount - westCount
		isBlack, present := tileMap[space]
		if !present {
			tileMap[space] = true
		} else {
			tileMap[space] = !isBlack
		}
	}

	blackTileCount := 0
	for _, isBlack := range tileMap {
		if isBlack {
			blackTileCount++
		}
	}
	fmt.Println(blackTileCount)
}
