package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

const (
	East int = iota
	South
	West
	North
)

var nsDist = 0
var ewDist = 0
var direction = East

type Move struct {
	direction string
	amount    int
}

func parseFile(path string) []Move {
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

	var moves []Move
	scanner := bufio.NewScanner(buf)
	fullRe := regexp.MustCompile(`(?P<direction>\w)(?P<arg>\d+)`)
	for scanner.Scan() {
		line := scanner.Text()
		matches := fullRe.FindStringSubmatch(line)
		direction := matches[fullRe.SubexpIndex("direction")]
		argStr := matches[fullRe.SubexpIndex("arg")]
		arg, err := strconv.Atoi(argStr)
		if err != nil {
			fmt.Println("arg value is not an integer", err)
			return moves
		}
		move := Move{
			direction,
			arg,
		}
		moves = append(moves, move)
	}
	return moves
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

	moves := parseFile(fp)

	for _, move := range moves {
		processMove(move)
		// 	row := getRow(pass.rowPattern)
		// 	column := getColumn(pass.columnPattern)
		// 	id := row*8 + column
		// 	if id > maxID {
		// 		maxID = id
		// 	}
	}
	fmt.Println(calcManhattanDist())
}

func processMove(move Move) {
	switch move.direction {
	case "L":
		direction -= move.amount / 90
		if direction < 0 {
			direction = North + direction + 1
		}
	case "R":
		direction += move.amount / 90
		direction = direction % 4
	case "E":
		ewDist += move.amount
	case "S":
		nsDist += move.amount
	case "W":
		ewDist -= move.amount
	case "N":
		nsDist -= move.amount
	case "F":
		fallthrough
	default:
		switch direction {
		case South:
			nsDist += move.amount
		case West:
			ewDist -= move.amount
		case North:
			nsDist -= move.amount
		case East:
			fallthrough
		default:
			ewDist += move.amount
		}
	}
	// fmt.Println(direction, nsDist, ewDist)
}

func calcManhattanDist() int {
	return int(math.Abs(float64(ewDist)) + math.Abs(float64(nsDist)))
}
