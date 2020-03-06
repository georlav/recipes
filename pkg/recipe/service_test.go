package recipe_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/georlav/recipes/pkg/recipe"

	"github.com/georlav/recipes/pkg/config"
)

func TestService_GetRecipes(t *testing.T) {
	testCases := []struct {
		desc       string
		qp         recipe.QueryParams
		respFile   string
		resLen     int
		statusCode int
		error      error
	}{
		{
			desc:       "Should successfully fetch 10 results",
			qp:         recipe.QueryParams{Page: 1},
			respFile:   "testdata/success.json",
			resLen:     10,
			statusCode: http.StatusOK,
		},
		{
			desc:       "Should successfully fetch no results",
			qp:         recipe.QueryParams{Page: 10000000},
			respFile:   "testdata/empty.json",
			statusCode: http.StatusOK,
		},
		{
			desc:       "Should fail due to invalid params",
			qp:         recipe.QueryParams{Page: -1},
			respFile:   "testdata/empty.json",
			statusCode: http.StatusInternalServerError,
			error:      recipe.ErrNoResults,
		},
		{
			desc:       "Should fail to unmarshal result due to invalid json response",
			qp:         recipe.QueryParams{},
			respFile:   "testdata/invalid.json",
			statusCode: http.StatusOK,
			error:      recipe.ErrUnmarshalResponse,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.desc, func(t *testing.T) {
			// Setup test server
			ts := server(tc.respFile, tc.statusCode, t)
			defer ts.Close()

			// Setup recipe puppy service
			c := config.RecipePuppyAPI{URL: ts.URL, Timeout: 5}
			s := recipe.NewService(c)

			result, err := s.Get(tc.qp)
			if err != nil && !errors.Is(err, tc.error) {
				t.Fatal(err)
			}

			if err != nil && !errors.Is(err, tc.error) {
				t.Fatalf("Expected to have error: \n%s\ngot\n%s", tc.error, err)
			}

			if err == nil && len(result.Results) != tc.resLen {
				t.Fatalf("Invalid result length expected %d got %d", tc.resLen, len(result.Results))
			}
		})
	}
}

func server(responseFile string, status int, t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)

		// nolint:gosec
		clientResp, err := ioutil.ReadFile(responseFile)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = w.Write(clientResp); err != nil {
			t.Fatal(err)
		}
	}))
}
