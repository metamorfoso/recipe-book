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

type RecipePullResult struct {
	Url          string
	Ingredients  RecipeSectionResult
	Instructions RecipeSectionResult
}

func PullRecipe(url string) (RecipePullResult, error) {
	result := RecipePullResult{
		Url:         url,
		Ingredients: RecipeSectionResult{},
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
		SubResult RecipeSectionResult
		Error     error
	}

	wg := &sync.WaitGroup{}
	ch := make(chan SubresultChannelOuput)

	wg.Add(1)
	go func() {
		ingredientsResult, err := findRecipeSection(doc, "ul", []string{ingredientsKeyword})
		ch <- SubresultChannelOuput{SubResult: ingredientsResult, Error: err, Type: "ingredients"}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		instructionsResult, err := findRecipeSection(doc, "ol", instructionsKeywords)
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
