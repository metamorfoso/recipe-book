package recipes

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func listElementSelectionsToCandidates(selections []*goquery.Selection) [][]string {
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

func findByClassOrIdContains(doc *goquery.Document, elTypes []string, keyword string) []*goquery.Selection {
	var matchingElements []*goquery.Selection

	for _, elType := range elTypes {
		doc.Find(elType).Each(func(_ int, selection *goquery.Selection) {
			class := selection.AttrOr("class", "")
			id := selection.AttrOr("id", "")

			if strings.Contains(class, keyword) || strings.Contains(id, keyword) {
				matchingElements = append(matchingElements, selection)
			}
		})
	}

	return matchingElements
}

func elementSelectionsToCandidates(selections []*goquery.Selection) [][]string {
	var candidates [][]string

	for _, selection := range selections {
		text := selection.Text()
		textItems := strings.Split(text, "\n")
		tidiedTextItems := tidyIngredients(textItems)
		candidates = append(candidates, tidiedTextItems)
	}

	return candidates
}
