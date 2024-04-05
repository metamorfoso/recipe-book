package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	// "github.com/gocolly/colly/v2"
)

var url string = "https://anaffairfromtheheart.com/meat-ragout/"
var url2 string = "https://www.delish.com/cooking/recipe-ideas/recipes/a47922/lemon-butter-chicken-pasta-recipe/"
var url3 string = "https://recipes.co.nz/recipes/the-best-smash-burgers/"
var url4 string = "https://mykoreankitchen.com/kimchi-recipe/"
var url5 string = "https://www.womensweeklyfood.com.au/recipe/baking/shepherds-pie-7402/"

var keyword string = "ingredient"

func getUrl(url string) *http.Response {
	res, err := http.Get(url)

	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("%v %v\n", res.StatusCode, url)

	if res.StatusCode != http.StatusOK {
		fmt.Printf("Exiting... got response status code %v for %v\n", res.StatusCode, url)
		os.Exit(1)
	}

	return res
}

func findUlByClassOrId(doc *goquery.Document) []*goquery.Selection {
	// look for ul elements that have a class or id that indicates ingredients
	var possibleIngredientUls []*goquery.Selection
	doc.Find("ul").Each(func(_ int, selection *goquery.Selection) {
		class := selection.AttrOr("class", "")
		id := selection.AttrOr("id", "")

		if strings.Contains(class, keyword) || strings.Contains(id, keyword) {
			possibleIngredientUls = append(possibleIngredientUls, selection)
		}
	})

	return possibleIngredientUls
}

func ulToCandidates(selections []*goquery.Selection) [][]string {
	var candidates [][]string
	for _, selection := range selections {
		var ingredientSet []string
		selection.Find("li").Each(func(_ int, s *goquery.Selection) {
			ingredientSet = append(ingredientSet, s.Text())
		})

		candidates = append(candidates, ingredientSet)
	}

	return candidates
}

func pullRecipeGoquery(url string) {
	fmt.Println("")
	fmt.Printf("Pulling recipe from %v\n", url)
	res := getUrl(url)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		fmt.Printf("error parsing document: %s\n", err)
		os.Exit(1)
	}

	var ingredientCandidates [][]string

	possibleIngredientUlsByClassOrId := findUlByClassOrId(doc)

	if len(possibleIngredientUlsByClassOrId) > 0 {
		ingredientCandidates = append(ingredientCandidates, ulToCandidates(possibleIngredientUlsByClassOrId)...)
	} else {
		fmt.Println("")
		fmt.Println("No ul found, checking other types of elements...")
		fmt.Println("")

		var matchingElements []*goquery.Selection

		doc.Find("*").Each(func(_ int, selection *goquery.Selection) {
			if len(matchingElements) > 0 {
				return
			}

			class := selection.AttrOr("class", "")
			id := selection.AttrOr("id", "")

			if strings.Contains(class, keyword) || strings.Contains(id, keyword) {
				matchingElements = append(matchingElements, selection)
			}
		})

		firstMatching := matchingElements[0]

		if firstMatching != nil {
			var textItems []string
			firstMatching.Find("*>*:last-of-type").Each(func(_ int, s *goquery.Selection) {
				t := s.Text()
				textItems = append(textItems, t)
			})
			dedupedTextItems := unique(textItems)

			var trimmed []string
			for _, item := range dedupedTextItems {
				trimmed = append(trimmed, strings.TrimSpace(item))
			}

			reDeduped := unique(trimmed)

			ingredientCandidates = append(ingredientCandidates, reDeduped)
		}
	}

	fmt.Println("")
	fmt.Printf("%v possible sets of ingredients found\n", len(ingredientCandidates))
	for index, set := range ingredientCandidates {
		fmt.Printf("Set %v:\n", index+1)
		for _, ingredient := range set {
			fmt.Printf("- %v\n", ingredient)
		}
	}
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
}

// func pullRecipe(url string) {
// 	res := getUrl(url)
// 	defer res.Body.Close()
//
// 	body, err := io.ReadAll(res.Body)
//
// 	if err != nil {
// 		fmt.Printf("error parsing response body: %s\n", err)
// 		os.Exit(1)
// 	}
//
// 	fmt.Printf("parsed body: %v \n", string(body))
// }

func main() {
	// pullRecipeGoquery(url)
	// pullRecipeGoquery(url2)
	// pullRecipeGoquery(url3)
	// pullRecipeGoquery(url4)
	pullRecipeGoquery(url5)
}

func unique(s []string) []string {
	inResult := make(map[string]bool)
	var result []string
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return result
}
