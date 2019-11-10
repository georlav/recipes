package recipe_test

import (
	"fmt"
	"testing"

	"github.com/georlav/recipes/config"
	"github.com/georlav/recipes/recipe"
)

func TestService_GetRecipes(t *testing.T) {
	c := config.RecipePuppyAPI{
		Host:    "http://www.recipepuppy.com",
		Timeout: 15,
	}
	s := recipe.NewService(c)

	testCases := []struct {
		desc        string
		page        int
		resultCount int
		error       error
	}{
		{
			"should fail due to invalid page number",
			0,
			0,
			fmt.Errorf("invalid response, %s, url: %s", "500 Internal Server Error", c.Host+"/api?p=0"),
		},
		{
			"Should successfully fetch 10 results",
			1,
			10,
			nil,
		},
		{
			"Should fetch no results",
			100000000,
			0,
			nil,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.desc, func(t *testing.T) {
			results, err := s.GetRecipes(tc.page)
			if err != nil && err.Error() != tc.error.Error() {
				t.Fatal(err)
			}

			if rlen := len(results); rlen != tc.resultCount {
				t.Fatalf("Invalid number of results expected %d got %d", tc.resultCount, rlen)
			}
		})
	}

}
