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

type Edges struct {
	transform []int
	sides     map[int][]int
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

var transformFuncMap = map[int]func(ImageSquare) ImageSquare{
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
	tileMap := make(map[int]ImageSquare)
	for _, tile := range tiles {
		combos := make(map[int]map[int]int)
		for _, transformPair := range uniqueTransforms {
			comboID, _ := strconv.Atoi(strconv.Itoa(transformPair[0]) + strconv.Itoa(transformPair[1]))
			combos[comboID] = calcEdges(transformFuncMap[transformPair[0]](transformFuncMap[transformPair[1]](tile)))
		}
		squareEdges[tile.id] = combos
		tileMap[tile.id] = tile
	}

	grid := make(map[int]Edges)

	tileID := tiles[0].id
	comboID := 33
	grid[tileID] = Edges{
		[]int{FLIP_HORZ, FLIP_HORZ},
		map[int][]int{
			Top:    []int{0, squareEdges[tileID][comboID][Top]},
			Bottom: []int{0, squareEdges[tileID][comboID][Bottom]},
			Left:   []int{0, squareEdges[tileID][comboID][Left]},
			Right:  []int{0, squareEdges[tileID][comboID][Right]},
		},
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
						if grid[tileID].sides[i][1] == grid[posTile.id].sides[sideMap[i]][1] {
							grid[tileID].sides[i][0] = posTile.id
							grid[posTile.id].sides[sideMap[i]][0] = tileID
						}
						continue
					}
					for _, posOrientation := range uniqueTransforms {
						posComboID, _ := strconv.Atoi(strconv.Itoa(posOrientation[0]) + strconv.Itoa(posOrientation[1]))
						if grid[tileID].sides[i][1] == squareEdges[posTile.id][posComboID][sideMap[i]] {
							grid[tileID].sides[i][0] = posTile.id
							grid[posTile.id] = Edges{
								posOrientation,
								map[int][]int{
									Top:    []int{0, squareEdges[posTile.id][posComboID][Top]},
									Bottom: []int{0, squareEdges[posTile.id][posComboID][Bottom]},
									Left:   []int{0, squareEdges[posTile.id][posComboID][Left]},
									Right:  []int{0, squareEdges[posTile.id][posComboID][Right]},
								},
							}
							grid[posTile.id].sides[sideMap[i]][0] = tileID
							break
						}
					}
				}
			}
		}
	}

	var topLeft int
	var corners []int
	for tileID, borders := range grid {
		edgeCount := 0
		top := false
		left := false
		for sideID, borderTile := range borders.sides {
			if borderTile[0] == 0 {
				if sideID == Top {
					top = true
				}
				if sideID == Left {
					left = true
				}
				edgeCount++
			}
		}
		if edgeCount == 2 {
			corners = append(corners, tileID)
			if top && left {
				topLeft = tileID
			}
		}
	}
	var topRow []ImageSquare
	for nextID := topLeft; nextID != 0; nextID = grid[nextID].sides[Right][0] {
		topRow = append(topRow, transformFuncMap[grid[nextID].transform[0]](transformFuncMap[grid[nextID].transform[1]](tileMap[nextID])))
	}

	tileLen := len(tiles[0].contents)
	var stitchedImage ImageSquare
	stitchedImage.id = 0
	nextRow := topRow
	for i := 0; i < len(topRow); i++ {
		var newNextRow []ImageSquare
		for j := 1; j < tileLen-1; j++ {
			var row []bool
			for _, image := range nextRow {
				row = append(row, image.contents[j][1:tileLen-1]...)
				if j == 1 {
					// fmt.Printf("%d ", image.id)
					below := grid[image.id].sides[Bottom][0]
					if below != 0 {
						nextTile := transformFuncMap[grid[below].transform[0]](transformFuncMap[grid[below].transform[1]](tileMap[below]))
						newNextRow = append(newNextRow, nextTile)
					}
				}
			}
			stitchedImage.contents = append(stitchedImage.contents, row)
		}
		// fmt.Printf("\n")
		nextRow = newNextRow
	}

	for _, posOrientation := range uniqueTransforms {
		transformedImage := transformFuncMap[posOrientation[0]](transformFuncMap[posOrientation[1]](stitchedImage))
		totalCount, monsterCount := getMonsterCount(transformedImage)
		if monsterCount != 0 {
			fmt.Println(totalCount - monsterCount)
			break
		}
	}
}

func getMonsterCount(image ImageSquare) (int, int) {
	totalCount := 0
	monsterCount := 0
	for i := 0; i < len(image.contents); i++ {
		for j := 0; j < len(image.contents); j++ {
			if image.contents[i][j] {
				totalCount++
			}
		}
	}
	monsterTemplate := [][]bool{
		[]bool{false, false, false, false, false, false, false, false, false, false,
			false, false, false, false, false, false, false, false, true, false},
		[]bool{true, false, false, false, false, true, true, false, false, false,
			false, true, true, false, false, false, false, true, true, true},
		[]bool{false, true, false, false, true, false, false, true, false, false,
			true, false, false, true, false, false, true, false, false, false},
	}
	hitsInMonster := 15
	monsterHeight := len(monsterTemplate)
	monsterLen := len(monsterTemplate[0])
	for i := 0; i < len(image.contents)-monsterHeight; i++ {
		for j := 0; j < len(image.contents)-monsterLen; j++ {
			isMonster := true
			for k, monsterRow := range monsterTemplate {
				for l, monsterPiece := range monsterRow {
					if monsterPiece && !image.contents[i+k][j+l] {
						isMonster = false
						break
					}
				}
				if !isMonster {
					break
				}
			}
			if isMonster {
				monsterCount += hitsInMonster
			}

		}
	}

	return totalCount, monsterCount
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
