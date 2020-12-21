package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
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

	var allergenList []string
	// possibleAllergens := make(map[string]map[string]bool)
	for englishName := range foodDictionary {
		allergenList = append(allergenList, englishName)
		// for possible := range possibleMap {
		// 	_, present := possibleAllergens[possible]
		// 	if !present {
		// 		possibleAllergens[possible] = make(map[string]bool)
		// 	}
		// 	possibleAllergens[possible][englishName] = true
		// }
	}
	// fmt.Println(possibleAllergens)

	incomplete := true
	finalDictionary := make(map[string]string)
	reverseDictionary := make(map[string]string)
	for incomplete {
		incomplete = false
		for allergen, possibleTranslations := range foodDictionary {
			_, present := finalDictionary[allergen]
			if len(possibleTranslations) == 1 {
				if !present {
					for translation := range possibleTranslations {
						finalDictionary[allergen] = translation
						reverseDictionary[translation] = allergen

					}
				}
				continue
			}
			if len(possibleTranslations) > 1 {
				for translation := range possibleTranslations {
					_, present := reverseDictionary[translation]
					if present {
						delete(possibleTranslations, translation)
					}
				}
				incomplete = true
			}
		}
	}

	sort.Sort(sort.StringSlice(allergenList))
	var translatedAllergentList []string
	for _, allergen := range allergenList {
		translatedAllergentList = append(translatedAllergentList, finalDictionary[allergen])
	}
	fmt.Println(strings.Join(translatedAllergentList, ","))
}
