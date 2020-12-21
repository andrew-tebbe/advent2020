package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Food struct {
	allergens   map[string]bool
	ingredients map[string]bool
}

func parseFile(path string) []Food {
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

	var foods []Food
	fullRe := regexp.MustCompile(`(?P<ingredients>[\w\s]+) \(contains (?P<allergens>[\w\s,]+)\)`)
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		matches := fullRe.FindStringSubmatch(scanner.Text())
		if matches != nil {
			ingredients := make(map[string]bool)
			allergens := make(map[string]bool)
			for _, ingredient := range strings.Split(matches[fullRe.SubexpIndex("ingredients")], " ") {
				ingredients[ingredient] = true
			}
			for _, allergen := range strings.Split(matches[fullRe.SubexpIndex("allergens")], ", ") {
				allergens[allergen] = true
			}

			foods = append(foods, Food{
				allergens,
				ingredients,
			})
		}
	}
	return foods
}

func copyMap(mapRef map[string]bool) map[string]bool {
	mapCopy := make(map[string]bool)
	for key, value := range mapRef {
		mapCopy[key] = value
	}
	return mapCopy
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

	foodDictionary := make(map[string]map[string]bool)
	foods := parseFile(fp)
	// fmt.Println(foods)
	for _, food := range foods {
		for allergen := range food.allergens {
			possibleTranslations, present := foodDictionary[allergen]
			if !present {
				foodDictionary[allergen] = copyMap(food.ingredients)
			} else {
				for posTrans := range possibleTranslations {
					_, present = food.ingredients[posTrans]
					if !present {
						delete(foodDictionary[allergen], posTrans)
					}
				}
			}
		}
	}

	possibleAllergens := make(map[string]bool)
	for _, possibleMap := range foodDictionary {
		for possible := range possibleMap {
			possibleAllergens[possible] = true
		}
	}
	notAllergenCount := 0
	for _, food := range foods {
		for foreignFood := range food.ingredients {
			_, present := possibleAllergens[foreignFood]
			if !present {
				notAllergenCount++
			}
		}
	}
	fmt.Println(notAllergenCount)
}
