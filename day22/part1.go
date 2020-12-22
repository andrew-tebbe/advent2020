package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

type Player struct {
	deck []int
}

func parseFile(path string) map[int][]int {
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

	players := make(map[int][]int)
	playerRe := regexp.MustCompile(`Player (?P<player>\d+):`)
	playerID := 3
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		matches := playerRe.FindStringSubmatch(line)
		if matches != nil {
			playerID, err = strconv.Atoi(matches[playerRe.SubexpIndex("player")])
			if err != nil {
				fmt.Println("player id is not an integer", err)
				return players
			}
			players[playerID] = []int{}
			continue
		}
		cardVal, err := strconv.Atoi(line)
		if err != nil {
			fmt.Println("card value is not an integer", err)
			return players
		}
		players[playerID] = append(players[playerID], cardVal)
	}
	return players
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

	players := parseFile(fp)
	for len(players[1]) > 0 && len(players[2]) > 0 {
		player1Card := players[1][0]
		player2Card := players[2][0]
		players[1] = players[1][1:]
		players[2] = players[2][1:]
		if player1Card > player2Card {
			players[1] = append(players[1], player1Card, player2Card)
		}
		if player2Card > player1Card {
			players[2] = append(players[2], player2Card, player1Card)
		}
	}
	fmt.Println(calcWinnerScore(players))
}

func calcWinnerScore(players map[int][]int) int {
	for _, deck := range players {
		if len(deck) > 0 {
			score := 0
			for i, card := range deck {
				score += card * (len(deck) - i)
			}
			return score
		}
	}
	return 0
}
