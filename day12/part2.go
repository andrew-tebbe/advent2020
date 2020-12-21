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

var nsWay = -1
var ewWay = 10
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
	}
	fmt.Println(calcManhattanDist())
}

func processMove(move Move) {
	switch move.direction {
	case "L":
		rotateLeft(move.amount)
	case "R":
		rotateRight(move.amount)
	case "E":
		ewWay += move.amount
	case "S":
		nsWay += move.amount
	case "W":
		ewWay -= move.amount
	case "N":
		nsWay -= move.amount
	case "F":
		fallthrough
	default:
		nsDist += move.amount * nsWay
		ewDist += move.amount * ewWay
	}
}

func rotateLeft(degrees int) {
	for i := 0; i < degrees/90; i++ {
		tmp := ewWay
		ewWay = nsWay
		nsWay = -tmp
	}
}

func rotateRight(degrees int) {
	for i := 0; i < degrees/90; i++ {
		tmp := ewWay
		ewWay = -nsWay
		nsWay = tmp
	}
}

func calcManhattanDist() int {
	return int(math.Abs(float64(ewDist)) + math.Abs(float64(nsDist)))
}
