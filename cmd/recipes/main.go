package main

import (
	"log"

	"github.com/georlav/recipes/config"
	"github.com/georlav/recipes/recipe"
)

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatal(err)
	}

	rs := recipe.NewService(cfg.RecipePuppyAPI)

	rs.GetRecipes()

}
