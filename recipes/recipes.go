package recipes

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

var ingredientsKeyword string = "ingredient"

var instructionsKeywords []string = []string{"instruction", "method"}

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

func (subResult *SubResult) appendCandidates(candidates [][]string) {
	subResult.Candidates = append(subResult.Candidates, candidates...)
}

type RecipePullResult struct {
	Url          string
	Ingredients  SubResult
	Instructions SubResult
}

func findIngredients(doc *goquery.Document) (SubResult, error) {
	result := SubResult{}

	result.DiscoveredVia = "ul"
	possibleIngredientUlsByClassOrId := findByClassOrIdContains(doc, "ul", ingredientsKeyword)

	if len(possibleIngredientUlsByClassOrId) > 0 {
		candidates := ulToCandidates(possibleIngredientUlsByClassOrId)
		result.appendCandidates(candidates)
		return result, nil
	}

	result.DiscoveredVia = "*"
	possibleIngredientsElements := findByClassOrIdContains(doc, "*", ingredientsKeyword)

	if len(possibleIngredientsElements) == 0 {
		return result, nil
	}

	candidates := textFromDeepestLastOfType(possibleIngredientsElements)
	result.appendCandidates(candidates)

	return result, nil
}

func findInstructions(doc *goquery.Document) (SubResult, error) {
	result := SubResult{}

	result.DiscoveredVia = "ol"

	possibleInstructionsOls := findByClassOrIdContains(doc, "ol", instructionsKeywords[0])

	if len(possibleInstructionsOls) > 0 {
		candidates := ulToCandidates(possibleInstructionsOls)
		result.appendCandidates(candidates)
		return result, nil
	}

	result.DiscoveredVia = "*"
	possibleInstructionsElements := findByClassOrIdContains(doc, "*", instructionsKeywords[0])

	if len(possibleInstructionsElements) == 0 {
		return result, nil
	}

	candidates := textFromDeepestLastOfType(possibleInstructionsElements)
	result.appendCandidates(candidates)

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

	type SubresultChannelOuput struct {
		Type      string
		SubResult SubResult
		Error     error
	}

	wg := &sync.WaitGroup{}
	ch := make(chan SubresultChannelOuput)

	wg.Add(1)
	go func() {
		ingredientsResult, err := findIngredients(doc)
		ch <- SubresultChannelOuput{SubResult: ingredientsResult, Error: err, Type: "ingredients"}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		instructionsResult, err := findInstructions(doc)
		ch <- SubresultChannelOuput{SubResult: instructionsResult, Error: err, Type: "instructions"}
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(ch)
	}()

	for output := range ch {
		switch output.Type {
		case "ingredients":
			result.Ingredients = output.SubResult
		case "instructions":
			result.Instructions = output.SubResult
		}
	}

	return result, nil
}
