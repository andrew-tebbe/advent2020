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

type GridSpace struct {
	northSouth int
	eastWest   float32
}

var neighborOffsets = []GridSpace{
	GridSpace{1, -0.5},
	GridSpace{1, 0.5},
	GridSpace{0, -1},
	GridSpace{-1, -0.5},
	GridSpace{-1, 0.5},
	GridSpace{0, 1},
}

func (space GridSpace) add(other GridSpace) GridSpace {
	return GridSpace{space.northSouth + other.northSouth,
		space.eastWest + other.eastWest}
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

	numDays, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("day count is not an integer", err)
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

	for i := 0; i < numDays; i++ {
		generateNeighbors(tileMap)
		tileMap = cycleDay(tileMap)

	}

	blackTileCount := 0
	for _, isBlack := range tileMap {
		if isBlack {
			blackTileCount++
		}
	}
	fmt.Println(blackTileCount)
	// fmt.Println(len(tileMap))
}

func generateNeighbors(tiles map[GridSpace]bool) {
	for space := range tiles {
		for _, offset := range neighborOffsets {
			neighbor := space.add(offset)
			_, present := tiles[neighbor]
			if !present {
				tiles[neighbor] = false
			}
		}
	}
}

func cycleDay(tiles map[GridSpace]bool) map[GridSpace]bool {
	newMap := make(map[GridSpace]bool)
	for space, isBlack := range tiles {
		adjBlackTiles := 0
		for _, offset := range neighborOffsets {
			if tiles[space.add(offset)] {
				adjBlackTiles++
			}
		}

		if isBlack && adjBlackTiles > 0 && adjBlackTiles < 3 {
			newMap[space] = true
		}
		if !isBlack && adjBlackTiles == 2 {
			newMap[space] = true
		}
	}
	return newMap
}
