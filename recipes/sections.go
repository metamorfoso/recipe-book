package recipes

import "github.com/PuerkitoBio/goquery"

type RecipeSectionResult struct {
	Candidates    Candidates
	DiscoveredVia string
}

func (subResult *RecipeSectionResult) appendCandidates(candidates [][]string) {
	subResult.Candidates = append(subResult.Candidates, candidates...)
}

func findRecipeSection(doc *goquery.Document, priorityElementType string, keywords []string) (RecipeSectionResult, error) {
	result := RecipeSectionResult{}

	possiblePriorityElementSelections := findByClassOrIdContains(doc, priorityElementType, keywords[0])

	if len(possiblePriorityElementSelections) > 0 {
		result.DiscoveredVia = priorityElementType

		var candidates [][]string

		switch priorityElementType {
		case "ul":
			candidates = ulToCandidates(possiblePriorityElementSelections)
		case "ol":
			candidates = ulToCandidates(possiblePriorityElementSelections)
			// TODO: other types of priority elements?
		}

		result.appendCandidates(candidates)
		return result, nil
	}

	result.DiscoveredVia = "*"
	possibleRelevantElementSelections := findByClassOrIdContains(doc, "*", instructionsKeywords[0])

	if len(possibleRelevantElementSelections) == 0 {
		return result, nil
	}

	candidates := textFromDeepestLastOfType(possibleRelevantElementSelections)
	result.appendCandidates(candidates)

	return result, nil
}
