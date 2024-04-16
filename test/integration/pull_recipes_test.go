package recipes

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/metamorfoso/recipe-book/recipes"
	"github.com/stretchr/testify/assert"
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
		// defer jsonPrint(output.Result)

		if output.Error != nil {
			fmt.Printf(">>>> ERROR processing %v: %v\n", output.Result.Url, output.Error)
			t.Fail()
		}

		content, err := os.ReadFile("./expected.json")
		if err != nil {
			t.Fatal("Error when opening file: ", err)
		}

		var payload []recipes.RecipePullResult
		err = json.Unmarshal(content, &payload)
		if err != nil {
			t.Fatal("Error during Unmarshal(): ", err)
		}

		var expected recipes.RecipePullResult
		for _, result := range payload {
			if result.Url == output.Result.Url {
				expected = result
				break
			}
		}

		actual := output.Result

		assert.EqualValues(t, actual, expected)
	}
}

func jsonPrint(i interface{}) {
	bytes, _ := json.MarshalIndent(i, "", "  ")
	formatted := string(bytes)
	fmt.Println(formatted)
}
