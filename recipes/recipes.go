package recipes

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

var ingredientsKeyword string = "ingredient"

func getUrl(url string) (*http.Response, error) {
	res, err := http.Get(url)

	if err != nil {
		fmt.Printf("Error making http request: %s\n", err)
	}

	if res.StatusCode != http.StatusOK {
		fmt.Printf("Got response status code %v for %v\n", res.StatusCode, url)
	}

	return res, err
}

type Candidates = [][]string

type SubResult struct {
	Candidates    Candidates
	DiscoveredVia string
}

func (subResult *SubResult) appendIngredients(ingredients [][]string) {
	subResult.Candidates = append(subResult.Candidates, ingredients...)
}

type RecipePullResult struct {
	Url           string
	Ingredients   SubResult
	CookingMethod SubResult
}

func findIngredients(doc *goquery.Document) (SubResult, error) {
	result := SubResult{}

	result.DiscoveredVia = "ul"
	possibleIngredientUlsByClassOrId := findByClassOrIdContains(doc, "ul", ingredientsKeyword)

	if len(possibleIngredientUlsByClassOrId) > 0 {
		candidates := ulToCandidates(possibleIngredientUlsByClassOrId)
		result.appendIngredients(candidates)
		return result, nil
	}

	result.DiscoveredVia = "*"
	possibleIngredientsElements := findByClassOrIdContains(doc, "*", ingredientsKeyword)

	if len(possibleIngredientsElements) == 0 {
		return result, nil
	}

	candidates := textFromDeepestLastOfType(possibleIngredientsElements)
	result.appendIngredients(candidates)

	return result, nil
}

func PullRecipe(url string) (RecipePullResult, error) {
	result := RecipePullResult{
		Url:         url,
		Ingredients: SubResult{},
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

	ingredientsResult, err := findIngredients(doc)
	if err != nil {
		return result, err
	}

	result.Ingredients = ingredientsResult

	return result, nil
}
