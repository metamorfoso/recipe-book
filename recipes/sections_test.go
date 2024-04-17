package recipes

import (
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func TestFindRecipeSection(t *testing.T) {
	fileReader, err := os.Open("../test/fixtures/meat-ragout.html")
	if err != nil {
		t.Errorf("Error opening test file: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(fileReader)
	if err != nil {
		t.Errorf("Error parsing html document: %v", err)
	}

	cases := []struct {
		name                 string
		priorityElementTypes []string
		classAndIdKeywords   []string
		expected             RecipeSectionResult
	}{
		{
			name:                 "Find ingredient section via ul & ol",
			priorityElementTypes: []string{"ul", "ol"},
			classAndIdKeywords:   []string{"ingredient"},
			expected: RecipeSectionResult{
				Candidates: [][]string{
					{"2 pounds lean ground beef", "1 3 ounce jar of Real Bacon Pieces", "1 medium onion chopped", "1/4 cup Italian dressing", "3 stalks celery chopped", "2 cups carrot coins or carrot slices roughly chopped", "6 cloves garlic minced", "4 14.5 ounce cans diced tomatoes, Garlic and Olive Oil flavored", "2 beef bouillon cube dissolved in 2 cups warm water", "1 6 oz can tomato paste", "1 16 oz package Rigatoni pasta, cooked and drained (for serving)", "Freshly Grated Parmesan cheese for serving"},
				},
				DiscoveredVia: "ul ol",
			},
		},
		{
			name:                 "Find instructions sectio via ul & oln",
			priorityElementTypes: []string{"ul", "ol"},
			classAndIdKeywords:   []string{"instruction", "method"},
			expected: RecipeSectionResult{
				Candidates:    [][]string{{"In a dutch oven over medium heat, cook the ground beef, veggies, and bacon until the beef is cooked and the veggies are tender.", "Stir in the rest of the ingredients (except pasta) and simmer on medium-low heat for about 1 1/4 hours until thickened.", "Serve over pasta."}},
				DiscoveredVia: "ul ol",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result, err := findRecipeSection(doc, c.priorityElementTypes, c.classAndIdKeywords)
			if err != nil {
				t.Errorf("Function returned an error: %v", err)
			}
			assert.EqualValues(t, c.expected, result)
		})
	}

}
