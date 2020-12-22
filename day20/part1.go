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

type Edge struct {
	side      int
	transform int
}

const (
	ROTATE_LEFT int = iota
	ROTATE_RIGHT
	FLIP_VERT
	FLIP_HORZ
	UNTRANSFORMED
)

const (
	Top int = iota
	Right
	Bottom
	Left
)

var pixToBool = map[rune]bool{
	'#': true,
	'.': false,
}

var boolToPix = map[bool]rune{
	true:  '#',
	false: '.',
}

var tranformFuncMap = map[int]func(ImageSquare) ImageSquare{
	ROTATE_LEFT:   rotateLeft,
	ROTATE_RIGHT:  rotateRight,
	FLIP_VERT:     flipVert,
	FLIP_HORZ:     flipHorz,
	UNTRANSFORMED: nullTransform,
}

var uniqueTransforms = [][]int{
	[]int{ROTATE_LEFT, UNTRANSFORMED},
	[]int{ROTATE_RIGHT, UNTRANSFORMED},
	[]int{FLIP_HORZ, UNTRANSFORMED},
	[]int{FLIP_VERT, UNTRANSFORMED},
	[]int{FLIP_HORZ, FLIP_HORZ},
	[]int{ROTATE_RIGHT, FLIP_HORZ},
	[]int{ROTATE_LEFT, FLIP_HORZ},
	[]int{FLIP_VERT, FLIP_HORZ},
}

var sideMap = map[int]int{
	Top:    Bottom,
	Bottom: Top,
	Left:   Right,
	Right:  Left,
}

type ImageSquare struct {
	id       int
	contents [][]bool
}

func nullTransform(square ImageSquare) ImageSquare {
	return square
}

func rotateRight(square ImageSquare) ImageSquare {
	squareLen := len(square.contents)
	var rotated ImageSquare
	rotated.id = square.id
	for i := 0; i < squareLen; i++ {
		var row []bool
		for j := 0; j < squareLen; j++ {
			row = append(row, square.contents[j][i])
		}
		rotated.contents = append(rotated.contents, row)
	}
	return rotated
}

func rotateLeft(square ImageSquare) ImageSquare {
	squareLen := len(square.contents)
	var rotated ImageSquare
	rotated.id = square.id
	for i := 0; i < squareLen; i++ {
		var row []bool
		for j := 0; j < squareLen; j++ {
			row = append(row, square.contents[squareLen-j-1][i])
		}
		rotated.contents = append(rotated.contents, row)
	}
	return rotated
}

func flipVert(square ImageSquare) ImageSquare {
	squareLen := len(square.contents)
	var flipped ImageSquare
	flipped.id = square.id
	for i := 0; i < squareLen; i++ {
		flipped.contents = append(flipped.contents, square.contents[squareLen-i-1])
	}
	return flipped
}

func flipHorz(square ImageSquare) ImageSquare {
	squareLen := len(square.contents)
	var flipped ImageSquare
	flipped.id = square.id
	for _, row := range square.contents {
		var flippedRow []bool
		for i := 0; i < squareLen; i++ {
			flippedRow = append(flippedRow, row[squareLen-i-1])
		}
		flipped.contents = append(flipped.contents, flippedRow)
	}
	return flipped
}

func parseFile(path string) []ImageSquare {
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

	var squares []ImageSquare
	var tile ImageSquare
	tileRe := regexp.MustCompile(`Tile (?P<title>\d+):`)
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			squares = append(squares, tile)
			tile = ImageSquare{}
			continue
		}

		matches := tileRe.FindStringSubmatch(line)
		if matches != nil {
			tileID := matches[tileRe.SubexpIndex("title")]
			id, err := strconv.Atoi(tileID)
			if err != nil {
				fmt.Println("id value is not an integer", err)
				return squares
			}
			tile.id = id
			continue
		}
		var lineList []bool
		for _, char := range line {
			lineList = append(lineList, pixToBool[char])
		}
		tile.contents = append(tile.contents, lineList)
	}
	squares = append(squares, tile)
	return squares
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

	squareEdges := make(map[int]map[int]map[int]int)
	for _, tile := range tiles {
		combos := make(map[int]map[int]int)
		for _, transformPair := range uniqueTransforms {
			comboID, _ := strconv.Atoi(strconv.Itoa(transformPair[0]) + strconv.Itoa(transformPair[1]))
			combos[comboID] = calcEdges(tranformFuncMap[transformPair[0]](tranformFuncMap[transformPair[1]](tile)))
		}
		squareEdges[tile.id] = combos
	}

	grid := make(map[int]map[int][]int)

	tileID := tiles[0].id
	comboID := 33
	grid[tileID] = map[int][]int{
		Top:    []int{0, squareEdges[tileID][comboID][Top]},
		Bottom: []int{0, squareEdges[tileID][comboID][Bottom]},
		Left:   []int{0, squareEdges[tileID][comboID][Left]},
		Right:  []int{0, squareEdges[tileID][comboID][Right]},
	}
	complete := false
	for !complete {
		complete = true
		for _, tile := range tiles {
			tileID := tile.id
			_, present := grid[tileID]
			if !present {
				complete = false
				continue
			}
			for i := Top; i <= Left; i++ {
				for _, posTile := range tiles {
					if tileID == posTile.id {
						continue
					}

					_, present := grid[posTile.id]
					if present {
						if grid[tileID][i][1] == grid[posTile.id][sideMap[i]][1] {
							grid[tileID][i][0] = posTile.id
							grid[posTile.id][sideMap[i]][0] = tileID
						}
						continue
					}
					for _, posOrientation := range uniqueTransforms {
						posComboID, _ := strconv.Atoi(strconv.Itoa(posOrientation[0]) + strconv.Itoa(posOrientation[1]))
						if grid[tileID][i][1] == squareEdges[posTile.id][posComboID][sideMap[i]] {
							grid[tileID][i][0] = posTile.id
							grid[posTile.id] = map[int][]int{
								Top:    []int{0, squareEdges[posTile.id][posComboID][Top]},
								Bottom: []int{0, squareEdges[posTile.id][posComboID][Bottom]},
								Left:   []int{0, squareEdges[posTile.id][posComboID][Left]},
								Right:  []int{0, squareEdges[posTile.id][posComboID][Right]},
							}
							grid[posTile.id][sideMap[i]][0] = tileID
							break
						}
					}
				}
			}
		}
	}

	var corners []int
	for tileID, borders := range grid {
		edgeCount := 0
		for _, borderTile := range borders {
			if borderTile[0] == 0 {
				edgeCount++
			}
		}
		if edgeCount == 2 {
			corners = append(corners, tileID)
		}
	}
	cornerMult := 1
	for _, id := range corners {
		cornerMult *= id
	}
	fmt.Println(cornerMult)
}

func calcEdges(tile ImageSquare) map[int]int {
	squareLen := len(tile.contents)
	edges := make(map[int]int)
	edges[Top] = 0
	edges[Bottom] = 0
	edges[Left] = 0
	edges[Right] = 0
	for i := 0; i < squareLen; i++ {
		if tile.contents[0][i] {
			edges[Top] += 1 << i
		}
		if tile.contents[squareLen-1][i] {
			edges[Bottom] += 1 << i
		}
		if tile.contents[i][0] {
			edges[Left] += 1 << i
		}
		if tile.contents[i][squareLen-1] {
			edges[Right] += 1 << i
		}
	}
	return edges
}

func printTile(tile ImageSquare) {
	fmt.Println("Tile", tile.id)
	for _, row := range tile.contents {
		for _, pixel := range row {
			fmt.Printf("%c", boolToPix[pixel])
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}
