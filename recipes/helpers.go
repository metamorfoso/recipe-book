package recipes

import "strings"

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
