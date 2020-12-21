package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type Coord struct {
	x int
	y int
	z int
}

func parseFile(path string) map[Coord]bool {
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

	initalState := make(map[Coord]bool)
	scanner := bufio.NewScanner(buf)

	y := 0
	for scanner.Scan() {
		for x, space := range scanner.Text() {
			coordinate := Coord{
				x,
				y,
				0}
			if space == '.' {
				initalState[coordinate] = false
			}
			if space == '#' {
				initalState[coordinate] = true
			}
		}
		y++
	}
	return initalState
}

func generateNeighbors(pos Coord) []Coord {
	var neighbors []Coord
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			for z := -1; z <= 1; z++ {
				if x == 0 && y == 0 && z == 0 {
					continue
				}
				coordinate := Coord{
					pos.x + x,
					pos.y + y,
					pos.z + z,
				}
				neighbors = append(neighbors, coordinate)
			}
		}
	}
	return neighbors
}

func runCycle(cubeState map[Coord]bool, lookUpTable map[Coord][]Coord) map[Coord]bool {
	newCubeState := copyMap(cubeState)
	for pos, state := range cubeState {
		neighbors, present := lookUpTable[pos]
		if !present {
			neighbors = generateNeighbors(pos)
			for _, neighbor := range neighbors {
				_, present := cubeState[neighbor]
				if !present {
					newCubeState[neighbor] = false
				}
			}
			lookUpTable[pos] = neighbors
		}
		activeCount := 0
		for _, neighbor := range neighbors {
			// if activeCount > 3 {
			// 	break
			// }
			neighborState, present := cubeState[neighbor]
			if present && neighborState {
				activeCount++
			}
		}

		if state && activeCount != 2 && activeCount != 3 {
			newCubeState[pos] = false
			// fmt.Println("Deactivating", pos)
		}
		if !state && activeCount == 3 {
			newCubeState[pos] = true
			// fmt.Println("Activating", pos)
		}
	}
	return newCubeState
}

func copyMap(aMap map[Coord]bool) map[Coord]bool {
	newMap := make(map[Coord]bool)
	for key, val := range aMap {
		newMap[key] = val
	}
	return newMap
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

	numCycles, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("numCycles value is not an integer", err)
		return
	}

	cubeState := parseFile(fp)
	lookUpTable := make(map[Coord][]Coord)
	newCubeState := copyMap(cubeState)
	for pos, state := range cubeState {
		neighbors := generateNeighbors(pos)
		for _, neighbor := range neighbors {
			_, present := newCubeState[neighbor]
			if !present {
				newCubeState[neighbor] = false
			}
		}
		newCubeState[pos] = state
		lookUpTable[pos] = neighbors
	}
	cubeState = newCubeState
	for i := 0; i < numCycles; i++ {
		cubeState = runCycle(cubeState, lookUpTable)
		// printCore(cubeState)
	}
	occupiedCount := 0
	for _, active := range cubeState {
		if active {
			occupiedCount++
		}
	}
	fmt.Println(occupiedCount)

}

func printCore(cubeState map[Coord]bool) {
	for z := -2; z <= 2; z++ {
		fmt.Println("z=", z)
		for y := -1; y <= 3; y++ {
			for x := -1; x <= 3; x++ {
				pos := Coord{
					x,
					y,
					z,
				}
				state, present := cubeState[pos]
				if present && state {
					fmt.Printf("#")
				} else {
					fmt.Printf(".")
				}
			}
			fmt.Printf("\n")
		}
		fmt.Printf("\n")
	}
}
