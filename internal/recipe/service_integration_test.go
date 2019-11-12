package recipe_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/georlav/recipes/internal/recipe"

	"github.com/georlav/recipes/internal/config"
)

func TestService_GetRecipes2(t *testing.T) {
	c := config.RecipePuppyAPI{
		Host:    "http://www.recipepuppy.com",
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
			recipe.QueryParams{Page: 10000000000},
			0,
			nil,
		},
		{
			"should fail due to invalid page number",
			recipe.QueryParams{Page: -1},
			0,
			fmt.Errorf("failed to retrive results, 500 %s", http.StatusText(http.StatusInternalServerError)),
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.desc, func(t *testing.T) {
			result, err := s.Get(tc.page)
			if err != nil && tc.error != nil {
				if tc.error.Error() != err.Error() {
					t.Fatal(err)
				}
			}

			if rlen := len(result.Results); rlen != tc.resultCount {
				t.Fatalf("Invalid number of results expected %d got %d", tc.resultCount, rlen)
			}
		})
	}

}
