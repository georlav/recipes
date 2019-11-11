package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/georlav/recipes/internal/config"
	"github.com/georlav/recipes/internal/recipe"
)

func main() {
	logger := log.Logger{}
	logger.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// Load application configuration
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatal(err)
	}

	// Enable debug output
	if !cfg.APP.Debug {
		logger.SetOutput(ioutil.Discard)
	}

	// Create a channel with page numbers
	pageNums := func() chan int {
		p := make(chan int)
		go func() {
			for i := 1; i <= cfg.APP.NumOfPages; i++ {
				p <- i
			}
			close(p)
		}()
		return p
	}()

	// Channel of received recipes
	recipeCH := make(chan recipe.Recipe)

	// slice of recipes we are going to print
	recipes := recipe.Recipes{}

	// start a go routine that Saves data to the slice we are going to print
	go func() {
		for r := range recipeCH {
			recipes.Append(r)
		}
	}()

	// Create a sync waitGroup
	wg := sync.WaitGroup{}
	wg.Add(cfg.APP.NumOfWorkers)

	// Initialize puppy recipe api
	rs := recipe.NewService(cfg.RecipePuppyAPI)

	// Start X goroutines
	for i := 1; i <= cfg.APP.NumOfWorkers; i++ {
		go func() {
			defer wg.Done()

			for p := range pageNums {
				logger.Println("Requesting recipes page", p)
				results, err := rs.Get(p)
				if err != nil {
					log.Fatalf("Failed to fetch page %d\n", p)
				}
				logger.Println("Retrieved recipes page", p)

				for i := range results.Results {
					recipeCH <- recipe.Recipe{
						Title:       results.Results[i].Title,
						Ingredients: results.Results[i].Ingredients,
						PageFound:   p,
					}
				}
			}
		}()
	}

	wg.Wait()

	// Print recipes to the standard output
	for _, v := range recipes.Values() {
		fmt.Printf("Title: %s\nIngredients: %s\nRecipe was found at page %d\n\n",
			v.Title, v.Ingredients, v.PageFound,
		)
	}

	fmt.Println("Total retrieved recipes: ", len(recipes.Values()))
}
