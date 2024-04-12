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
	}

	if res.StatusCode != http.StatusOK {
		fmt.Printf("Exiting... got response status code %v for %v\n", res.StatusCode, url)
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

type RecipePullResult struct {
	Url             string
	Ingredients     IngredientCandidates
	DiscoveryMethod string
}

func (r *RecipePullResult) appendIngredients(ingredients [][]string) {
	r.Ingredients = append(r.Ingredients, ingredients...)
}

func pullRecipe(url string) (RecipePullResult, error) {
	result := RecipePullResult{
		Url:         url,
		Ingredients: IngredientCandidates{},
	}

	fmt.Printf("Pulling recipe from %v\n", url)
	res, err := getUrl(url)

	if err != nil {
		return result, err
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		fmt.Printf("error parsing document: %s\n", err)
		return result, err
	}

	result.DiscoveryMethod = "ul"
	possibleIngredientUlsByClassOrId := findByClassOrIdContains(doc, "ul", ingredientsKeyword)

	if len(possibleIngredientUlsByClassOrId) > 0 {
		candidates := ulToCandidates(possibleIngredientUlsByClassOrId)
		result.appendIngredients(candidates)
		return result, nil
	}

	result.DiscoveryMethod = "*"
	possibleIngredientsElements := findByClassOrIdContains(doc, "*", ingredientsKeyword)

	if len(possibleIngredientsElements) == 0 {
		return result, nil
	}

	// Note: for now it seems only the first match is relevant. This needs more exploration.
	firstMatching := possibleIngredientsElements[0]

	var textItems []string
	firstMatching.Find("*>*:last-of-type").Each(func(_ int, s *goquery.Selection) {
		t := s.Text()
		textItems = append(textItems, t)
	})

	tidiedTextItems := tidyIngredients(textItems)

	result.appendIngredients(IngredientCandidates{tidiedTextItems})
	return result, nil
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

type ChannelOutput struct {
	Result RecipePullResult
	Error  error
}

func testPullRecipes() {
	wg := &sync.WaitGroup{}

	channel := make(chan ChannelOutput)

	asyncPullRecipe := func(url string) {
		pullResult, err := pullRecipe(url)
		channel <- ChannelOutput{Result: pullResult, Error: err}
		wg.Done()
	}

	for _, url := range testUrls {
		wg.Add(1)
		go asyncPullRecipe(url)
	}

	go func() {
		wg.Wait()
		close(channel)
	}()

	for output := range channel {
		fmt.Println("------")

		if output.Error != nil {
			fmt.Printf(">>>> ERROR processing %v: %v\n", output.Result.Url, output.Error)
		} else {
			fmt.Printf("Possible ingredients for %v:\n", output.Result.Url)
			fmt.Printf("(Discovered by looking through %v elements)\n", output.Result.DiscoveryMethod)

			for index, ingredientGroup := range output.Result.Ingredients {
				fmt.Printf("Set %v:\n", index+1)
				for _, ingredient := range ingredientGroup {
					fmt.Printf("- %v\n", ingredient)
				}
			}
		}
	}
}
