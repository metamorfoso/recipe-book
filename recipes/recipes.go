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

type IngredientCandidates = [][]string

type RecipePullResult struct {
	Url             string
	Ingredients     IngredientCandidates
	DiscoveryMethod string
}

func (r *RecipePullResult) appendIngredients(ingredients [][]string) {
	r.Ingredients = append(r.Ingredients, ingredients...)
}

func PullRecipe(url string) (RecipePullResult, error) {
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

	candidates := textFromDeepestLastOfType(possibleIngredientsElements)
	result.appendIngredients(candidates)
	return result, nil
}
