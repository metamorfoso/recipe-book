package recipes

import (
	"fmt"
	"sync"
	"testing"

	"github.com/metamorfoso/recipe-book/recipes"
)

type ChannelOutput struct {
	Result recipes.RecipePullResult
	Error  error
}

var testUrls []string = []string{
	"https://anaffairfromtheheart.com/meat-ragout/",
	"https://www.delish.com/cooking/recipe-ideas/recipes/a47922/lemon-butter-chicken-pasta-recipe/",
	"https://recipes.co.nz/recipes/the-best-smash-burgers/",
	"https://mykoreankitchen.com/kimchi-recipe/",
	"https://www.womensweeklyfood.com.au/recipe/baking/shepherds-pie-7402/",
	"https://khinskitchen.com/lamb-karahi/",
	"https://www.recipetineats.com/beef-barbacoa/",
}

func TestPullRecipe(t *testing.T) {
	wg := &sync.WaitGroup{}

	channel := make(chan ChannelOutput)

	asyncPullRecipe := func(url string) {
		pullResult, err := recipes.PullRecipe(url)
		channel <- ChannelOutput{Result: pullResult, Error: err}
		wg.Done()
	}

	for _, url := range testUrls {
		wg.Add(1)
		go asyncPullRecipe(url)
	}

	go func() {
		wg.Wait()
		close(channel)
	}()

	for output := range channel {
		fmt.Println("------")

		if output.Error != nil {
			fmt.Printf(">>>> ERROR processing %v: %v\n", output.Result.Url, output.Error)
			t.Fail()
		} else {
			fmt.Printf("Possible ingredients for %v:\n", output.Result.Url)
			fmt.Printf("(Discovered by looking through %v elements)\n", output.Result.Ingredients.DiscoveredVia)

			for index, ingredientGroup := range output.Result.Ingredients.Candidates {
				fmt.Printf("Set %v:\n", index+1)
				for _, ingredient := range ingredientGroup {
					fmt.Printf("- %v\n", ingredient)
				}
			}

			fmt.Printf("Possible instructions for %v:\n", output.Result.Url)
			fmt.Printf("(Discovered by looking through %v elements)\n", output.Result.Instructions.DiscoveredVia)

			for index, instructionsGroup := range output.Result.Instructions.Candidates {
				fmt.Printf("Set %v:\n", index+1)
				for _, instruction := range instructionsGroup {
					fmt.Printf("- %v\n", instruction)
				}
			}
		}
	}
}
