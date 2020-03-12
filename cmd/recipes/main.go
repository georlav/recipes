package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/georlav/recipes/pkg/config"
	"github.com/georlav/recipes/pkg/recipe"
)

func main() {
	// initialize logger
	logger := log.New(
		os.Stdout, "", log.LstdFlags|log.Lmicroseconds,
	)

	// Load application configuration
	cfg, err := config.Load("config.json")
	if err != nil {
		logger.Fatal(err)
	}

	// Disable debug output
	if !cfg.APP.Debug {
		logger.SetOutput(ioutil.Discard)
	}

	// User can also change those from command line, using as defaults the cfg values
	flag.IntVar(&cfg.APP.NumOfPages, "pages", cfg.APP.NumOfPages, "number of pages to retrieve --page=1")
	flag.IntVar(&cfg.APP.NumOfWorkers, "workers", cfg.APP.NumOfWorkers, "number of workers to start --workers=1")
	flag.Parse()

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

	// Initialize puppy recipe client
	rs := recipe.NewService(cfg.RecipePuppyAPI)

	// Start X goroutines
	for i := 1; i <= cfg.APP.NumOfWorkers; i++ {
		go func() {
			defer wg.Done()

			for p := range pageNums {
				logger.Println("Requesting recipes page", p)
				results, err := rs.Get(recipe.QueryParams{Page: p})
				if err != nil {
					logger.Fatalf("Failed to fetch page %d\n", p)
				}
				logger.Println("Retrieved recipes page", p)

				for i := range results.Results {
					ingredients := strings.Split(results.Results[i].Ingredients, ",")
					for i := range ingredients {
						ingredients[i] = strings.TrimSpace(ingredients[i])
					}

					recipeCH <- recipe.Recipe{
						Title:       results.Results[i].Title,
						URL:         results.Results[i].Href,
						Thumbnail:   results.Results[i].Thumbnail,
						Ingredients: ingredients,
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
	fmt.Println("Total pages: ", cfg.APP.NumOfPages)
	fmt.Println("Results per page: 10")
	fmt.Println("Total retrieved recipes: ", len(recipes.Values()))

	// Save result to file for later use
	if err := recipes.Save(); err != nil {
		logger.Fatal(err)
	}
}
