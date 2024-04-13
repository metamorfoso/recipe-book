package recipes

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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

func textFromDeepestLastOfType(possibleIngredientsElements []*goquery.Selection) [][]string {
	// Note: for now it seems only the first match is relevant. This needs more exploration.
	firstMatching := possibleIngredientsElements[0]

	var textItems []string
	firstMatching.Find("*>*:last-of-type").Each(func(_ int, s *goquery.Selection) {
		t := s.Text()
		textItems = append(textItems, t)
	})

	tidiedTextItems := tidyIngredients(textItems)

	return Candidates{tidiedTextItems}

}
