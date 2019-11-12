package recipe_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/georlav/recipes/internal/recipe"

	"github.com/georlav/recipes/internal/config"
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
			error:      fmt.Errorf("failed to retrive results, 500 %s", http.StatusText(http.StatusInternalServerError)),
		},
		{
			desc:       "Should fail to unmarshal result due to invalid json response",
			qp:         recipe.QueryParams{},
			respFile:   "testdata/invalid.json",
			statusCode: http.StatusOK,
			error:      fmt.Errorf("failed to unmarshal response, invalid character '}' looking for beginning of value"),
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.desc, func(t *testing.T) {
			// Setup test server
			ts := server(tc.respFile, tc.statusCode, t)
			defer ts.Close()

			// Setup recipe puppy service
			c := config.RecipePuppyAPI{Host: ts.URL, Timeout: 5}
			s := recipe.NewService(c)

			result, err := s.Get(tc.qp)
			if err != nil && tc.error == nil {
				t.Fatal(err)
			}
			if tc.error != nil && err.Error() != tc.error.Error() {
				t.Fatalf("Expected to have error: \n%s\ngot\n%s", tc.error, err)
			}

			if rlen := len(result.Results); rlen != tc.resLen {
				t.Fatalf("Invalid result length expected %d got %d", tc.resLen, rlen)
			}
		})
	}
}

func server(responseFile string, status int, t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)

		clientResp, err := ioutil.ReadFile(responseFile)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = w.Write(clientResp); err != nil {
			t.Fatal(err)
		}
	}))
}
