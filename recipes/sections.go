package recipes

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type RecipeSectionResult struct {
	Candidates    Candidates
	DiscoveredVia string
}

func (subResult *RecipeSectionResult) appendCandidates(candidates [][]string) {
	subResult.Candidates = append(subResult.Candidates, candidates...)
}

func findRecipeSection(doc *goquery.Document, priorityElementTypes []string, keywords []string) (RecipeSectionResult, error) {
	result := RecipeSectionResult{}

	for _, keyword := range keywords {
		selections := findByClassOrIdContains(doc, priorityElementTypes, keyword)

		if len(selections) == 0 {
			continue
		}

		candidates := listElementSelectionsToCandidates(selections)
		result.DiscoveredVia = strings.Join(priorityElementTypes, " ")
		result.appendCandidates(candidates)
	}

	if len(result.Candidates) > 0 {
		return result, nil
	}

	for _, keyword := range keywords {
		elementTypesToSearch := []string{"*"}
		selections := findByClassOrIdContains(doc, elementTypesToSearch, keyword)

		if len(selections) == 0 {
			continue
		}

		candidates := elementSelectionsToCandidates(selections)
		result.DiscoveredVia = "*"
		result.appendCandidates(candidates)
	}

	return result, nil
}
