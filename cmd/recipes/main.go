package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/georlav/recipes/config"
	"github.com/georlav/recipes/recipe"
)

func main() {
	// Load application configuration
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatal(err)
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
	recipeCH := make(chan recipe.Recipe, cfg.APP.NumOfWorkers)

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
				results, err := rs.Get(p)
				if err != nil {
					log.Fatalf("Failed to fetch page %d\n", p)
				}

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

	// Print retrieved recipes to screen
	for _, v := range recipes.Values() {
		fmt.Println(v)
	}

	fmt.Println("Total retrieved recipes: ", len(recipes.Values()))
}
