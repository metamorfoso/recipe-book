package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

var testUrls []string = []string{
	"https://anaffairfromtheheart.com/meat-ragout/",
	"https://www.delish.com/cooking/recipe-ideas/recipes/a47922/lemon-butter-chicken-pasta-recipe/",
	"https://recipes.co.nz/recipes/the-best-smash-burgers/",
	"https://mykoreankitchen.com/kimchi-recipe/",
	"https://www.womensweeklyfood.com.au/recipe/baking/shepherds-pie-7402/",
	"https://khinskitchen.com/lamb-karahi/",
	"https://www.recipetineats.com/beef-barbacoa/",
}

var ingredientsKeyword string = "ingredient"

func getUrl(url string) (*http.Response, error) {
	res, err := http.Get(url)

	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		// os.Exit(1)
	}

	// fmt.Printf("%v %v\n", res.StatusCode, url)

	if res.StatusCode != http.StatusOK {
		fmt.Printf("Exiting... got response status code %v for %v\n", res.StatusCode, url)
		// os.Exit(1)
	}

	return res, err
}

func ulToCandidates(selections []*goquery.Selection) [][]string {
	var candidates [][]string
	for _, selection := range selections {
		var ingredientSet []string
		selection.Find("li").Each(func(_ int, s *goquery.Selection) {
			ingredientSet = append(ingredientSet, s.Text())
		})

		tidiedIngredients := tidyIngredients(ingredientSet)

		candidates = append(candidates, tidiedIngredients)
	}

	return candidates
}

func findByClassOrIdContains(doc *goquery.Document, elType string, keyword string) []*goquery.Selection {
	var matchingElements []*goquery.Selection

	doc.Find(elType).Each(func(_ int, selection *goquery.Selection) {
		class := selection.AttrOr("class", "")
		id := selection.AttrOr("id", "")

		if strings.Contains(class, keyword) || strings.Contains(id, keyword) {
			matchingElements = append(matchingElements, selection)
		}
	})

	return matchingElements
}

type IngredientCandidates = [][]string

func pullRecipe(url string) IngredientCandidates {
	fmt.Printf("Pulling recipe from %v\n", url)
	res, err := getUrl(url)

	if err != nil {
		return [][]string{}
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		fmt.Printf("error parsing document: %s\n", err)
		// os.Exit(1)
		return [][]string{}
	}

	var ingredientCandidates IngredientCandidates

	possibleIngredientUlsByClassOrId := findByClassOrIdContains(doc, "ul", ingredientsKeyword)

	if len(possibleIngredientUlsByClassOrId) > 0 {
		ingredientCandidates = append(ingredientCandidates, ulToCandidates(possibleIngredientUlsByClassOrId)...)
	} else {
		// fmt.Println("")
		fmt.Printf("No ul found in %v, checking other types of elements...\n", url)
		// fmt.Println("")

		possibleIngredientsElements := findByClassOrIdContains(doc, "*", ingredientsKeyword)

		if len(possibleIngredientsElements) == 0 {
			fmt.Printf("No elements in %v found whose class or id contains keyword '%v'\n", url, ingredientsKeyword)
		} else {
			// Note: for now it seems only the first match is relevant. This needs more exploration.
			firstMatching := possibleIngredientsElements[0]

			var textItems []string
			firstMatching.Find("*>*:last-of-type").Each(func(_ int, s *goquery.Selection) {
				t := s.Text()
				textItems = append(textItems, t)
			})

			tidiedTextItems := tidyIngredients(textItems)

			ingredientCandidates = append(ingredientCandidates, tidiedTextItems)
		}
	}

	// fmt.Println("")
	// fmt.Printf("%v possible sets of ingredients found\n", len(ingredientCandidates))
	// for index, set := range ingredientCandidates {
	// 	fmt.Printf("Set %v:\n", index+1)
	// 	for _, ingredient := range set {
	// 		fmt.Printf("- %v\n", ingredient)
	// 	}
	// }
	// fmt.Println("")
	// fmt.Println("")
	// fmt.Println("")

	return ingredientCandidates
}

func tidyIngredients(textItems []string) []string {
	var trimmed []string
	for _, item := range textItems {
		trimmed = append(trimmed, strings.TrimSpace(item))
	}

	reDeduped := unique(trimmed)

	return reDeduped
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

func main() {
	wg := &sync.WaitGroup{}

	channel := make(chan map[string]IngredientCandidates)

	asyncPullRecipe := func(url string) {
		candidates := pullRecipe(url)
		candidatesForUrl := map[string]IngredientCandidates{
			url: candidates,
		}
		channel <- candidatesForUrl
		wg.Done()
	}

	for _, url := range testUrls {
		wg.Add(1)
		go asyncPullRecipe(url)
		// pullRecipe(url)
	}

	go func() {
		wg.Wait()
		close(channel)
	}()

	for result := range channel {
		for k, v := range result {
			fmt.Println("------")
			fmt.Printf("Possible ingredients for %v:\n", k)

			for index, ingredientGroup := range v {
				fmt.Printf("Set %v:\n", index+1)
				for _, ingredient := range ingredientGroup {
					fmt.Printf("- %v\n", ingredient)
				}
			}
			fmt.Println("------")
		}
	}
}
