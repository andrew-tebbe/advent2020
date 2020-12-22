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
	winner := playGame(players, 0)
	fmt.Println(calcScores(players)[winner[0]])
}

func copyMap(aMap map[int][]int) map[int][]int {
	newMap := make(map[int][]int)
	for key, val := range aMap {
		newMap[key] = []int{}
		for _, i := range val {
			newMap[key] = append(newMap[key], i)
		}
	}
	return newMap
}

func playGame(players map[int][]int, depth int) []int {
	// player1Deck score: [] player2Deck score
	duplicateMap := make(map[int][]int)
	for len(players[1]) > 0 && len(players[2]) > 0 {
		playerScores := calcScores(players)
		player2Scores, present := duplicateMap[playerScores[1]]
		if present {
			for _, score := range player2Scores {
				if score == playerScores[2] {
					return []int{1, 2}
				}
			}
			duplicateMap[playerScores[1]] = append(duplicateMap[playerScores[1]], playerScores[2])
		} else {
			duplicateMap[playerScores[1]] = []int{playerScores[2]}
		}
		playRound(players, depth)
	}
	if len(players[1]) > 0 {
		return []int{1, 2}
	}
	return []int{2, 1}
}

func playRound(players map[int][]int, depth int) {
	dealtCards := make(map[int]int)
	dealtCards[1] = players[1][0]
	dealtCards[2] = players[2][0]
	players[1] = players[1][1:]
	players[2] = players[2][1:]
	if dealtCards[1] <= len(players[1]) && dealtCards[2] <= len(players[2]) {
		subGameMap := copyMap(players)
		subGameMap[1] = subGameMap[1][:dealtCards[1]]
		subGameMap[2] = subGameMap[2][:dealtCards[2]]
		winners := playGame(subGameMap, depth+1)
		players[winners[0]] = append(players[winners[0]], dealtCards[winners[0]])
		players[winners[0]] = append(players[winners[0]], dealtCards[winners[1]])
		return
	}
	if dealtCards[1] > dealtCards[2] {
		players[1] = append(players[1], dealtCards[1], dealtCards[2])
	}
	if dealtCards[2] > dealtCards[1] {
		players[2] = append(players[2], dealtCards[2], dealtCards[1])
	}
}

func calcScores(players map[int][]int) map[int]int {
	scoreMap := make(map[int]int)
	for playerID, deck := range players {
		score := 0
		if len(deck) > 0 {
			for i, card := range deck {
				score += card * (len(deck) - i)
			}
		}
		scoreMap[playerID] = score
	}
	return scoreMap
}
