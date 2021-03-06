package recipe_test

import (
	"errors"
	"testing"

	"github.com/georlav/recipes/pkg/recipe"

	"github.com/georlav/recipes/pkg/config"
)

func TestService_GetRecipes2(t *testing.T) {
	c := config.RecipePuppyAPI{
		URL:     "http://www.recipepuppy.com",
		Timeout: 15,
	}
	s := recipe.NewService(c)

	testCases := []struct {
		desc        string
		page        recipe.QueryParams
		resultCount int
		error       error
	}{
		{
			"Should successfully fetch 10 results",
			recipe.QueryParams{Page: 1},
			10,
			nil,
		},
		{
			"Should successfully fetch 10 results",
			recipe.QueryParams{},
			10,
			nil,
		},
		{
			"Should successfully fetch  0 results",
			recipe.QueryParams{Page: 100000},
			0,
			nil,
		},
		{
			"should fail due to invalid page number",
			recipe.QueryParams{Page: -1},
			0,
			recipe.ErrNoResults,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.desc, func(t *testing.T) {
			result, err := s.Get(tc.page)
			if err != nil && !errors.Is(err, tc.error) {
				t.Fatalf("Unexpected error, expected \n%s got \n%s", tc.error, err)
			}

			if err == nil && len(result.Results) != tc.resultCount {
				t.Fatalf("Invalid number of results expected %d got %d", tc.resultCount, len(result.Results))
			}
		})
	}
}
