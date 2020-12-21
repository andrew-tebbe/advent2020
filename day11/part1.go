package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	OCCUPIED int = iota
	FLOOR
	EMPTY
)

func parseFile(path string) [][]int {
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

	var grid [][]int
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		var line []int
		for _, space := range scanner.Text() {
			ferrySpace := FLOOR
			if space == '.' {
				ferrySpace = FLOOR
			}
			if space == 'L' {
				ferrySpace = EMPTY
			}
			if space == '#' {
				ferrySpace = OCCUPIED
			}
			line = append(line, ferrySpace)
		}
		grid = append(grid, line)
	}
	return grid
}

func takeSeats(curFerry [][]int) ([][]int, bool) {
	var seats [][]int
	rows := len(curFerry)
	columns := len(curFerry[0])
	changed := false
	for i, seatRow := range curFerry {
		var row []int
		for j, space := range seatRow {
			var adjacentSeats []int
			if space == FLOOR {
				row = append(row, space)
				continue
			}
			if i > 0 {
				adjacentSeats = append(adjacentSeats, curFerry[i-1][j])
				if j > 0 {
					adjacentSeats = append(adjacentSeats, curFerry[i-1][j-1])
				}
				if j < columns-1 {
					adjacentSeats = append(adjacentSeats, curFerry[i-1][j+1])
				}
			}
			if i < rows-1 {
				adjacentSeats = append(adjacentSeats, curFerry[i+1][j])
				if j > 0 {
					adjacentSeats = append(adjacentSeats, curFerry[i+1][j-1])
				}
				if j < columns-1 {
					adjacentSeats = append(adjacentSeats, curFerry[i+1][j+1])
				}
			}
			if j > 0 {
				adjacentSeats = append(adjacentSeats, curFerry[i][j-1])
			}
			if j < columns-1 {
				adjacentSeats = append(adjacentSeats, curFerry[i][j+1])
			}
			occupiedCount := 0
			for _, state := range adjacentSeats {
				if state == OCCUPIED {
					occupiedCount++
				}
			}
			if space == OCCUPIED && occupiedCount >= 4 {
				row = append(row, EMPTY)
				changed = true
				continue
			}
			if space == EMPTY && occupiedCount == 0 {
				row = append(row, OCCUPIED)
				changed = true
				continue
			}
			row = append(row, space)
		}
		seats = append(seats, row)
	}
	return seats, changed
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

	ferryMap := parseFile(fp)
	for changed := true; changed; {
		ferryMap, changed = takeSeats(ferryMap)
		// printMap(ferryMap)
		// fmt.Println()
		// ferryMap, _ = takeSeats(ferryMap)
		// changed = false
	}

	occupiedCount := 0
	for i := 0; i < len(ferryMap); i++ {
		for j := 0; j < len(ferryMap[i]); j++ {
			if ferryMap[i][j] == OCCUPIED {
				occupiedCount++
			}
		}
	}
	fmt.Println(occupiedCount)
}

func printMap(ferryMap [][]int) {
	for _, row := range ferryMap {
		rowString := ""
		for _, column := range row {
			if column == OCCUPIED {
				rowString += "#"
			}
			if column == EMPTY {
				rowString += "L"
			}
			if column == FLOOR {
				rowString += "."
			}
		}
		fmt.Println(rowString)
	}
}
