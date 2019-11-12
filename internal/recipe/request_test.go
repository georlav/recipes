package recipe_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/georlav/recipes/internal/recipe"
)

func TestQueryParams_Encode(t *testing.T) {
	testCases := []struct {
		input  recipe.QueryParams
		output string
	}{
		{input: recipe.QueryParams{}, output: ""},
		{input: recipe.QueryParams{Page: 0}, output: ""},
		{input: recipe.QueryParams{Page: 1}, output: "p=1"},
		{input: recipe.QueryParams{Page: 2, Term: "test"}, output: "p=2&q=test"},
		{input: recipe.QueryParams{Page: 3, Ingredients: "in1"}, output: "i=in1&p=3"},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(fmt.Sprintf(`Encoding %+v`, tc.input), func(t *testing.T) {
			if tc.input.Encode() != tc.output {
				log.Fatalf("Invalid parameter encoding expected %s got %s", tc.input.Encode(), tc.output)
			}
		})
	}
}
