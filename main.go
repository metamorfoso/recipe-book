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
var keyword string = "ingredient"

func main() {
	res, err := http.Get(url)

	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}
	defer res.Body.Close()

	fmt.Printf("%v %v\n", res.StatusCode, url)

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		fmt.Printf("error parsing document: %s\n", err)
		os.Exit(1)
	}

	// look for ul elements that have a class or id that indicates ingredients
	var possibleIngredientUls []*goquery.Selection
	doc.Find("ul").Each(func(_ int, selection *goquery.Selection) {
		class := selection.AttrOr("class", "")
		id := selection.AttrOr("id", "")

		if strings.Contains(class, keyword) || strings.Contains(id, keyword) {
			possibleIngredientUls = append(possibleIngredientUls, selection)
		}
	})

	var ingredientCandidates [][]string

	for _, selection := range possibleIngredientUls {
		var ingredientSet []string
		selection.Find("li").Each(func(_ int, s *goquery.Selection) {
			ingredientSet = append(ingredientSet, s.Text())
		})

		ingredientCandidates = append(ingredientCandidates, ingredientSet)
	}

	fmt.Printf("%v possible sets of ingredients found\n", len(ingredientCandidates))
	for index, set := range ingredientCandidates {
		fmt.Printf("Set %v\n", index+1)
		for _, ingredient := range set {
			fmt.Printf("- %v\n", ingredient)
		}
	}

	// search through all elements in the document
	var elementsMentioningIngredients []*goquery.Selection
	doc.Find("*").Each(func(_ int, selection *goquery.Selection) {
		text := strings.ToLower(selection.Text())

		if strings.Contains(text, keyword) {
			elementsMentioningIngredients = append(elementsMentioningIngredients, selection)
		}
	})

	// for _, selection := range elementsMentioningIngredients {
	// 	fmt.Println(selection.Text())
	// 	for _, node := range selection.Nodes {
	// 		fmt.Println(node.Data)
	// 	}
	// }

	// fmt.Println(elementsMentioningIngredients)

	// c := colly.NewCollector()
	//
	// // Find and visit all links
	// c.OnHTML("ul", func(e *colly.HTMLElement) {
	// 	text := e.Text
	// 	fmt.Println(text)
	// })
	//
	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting", r.URL)
	// })
	//
	// c.Visit("https://anaffairfromtheheart.com/meat-ragout/")
}
