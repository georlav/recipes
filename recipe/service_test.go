package recipe_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/georlav/recipes/config"
	"github.com/georlav/recipes/recipe"
)

func TestService_GetRecipes(t *testing.T) {
	testCases := []struct {
		desc       string
		page       int
		respFile   string
		resLen     int
		statusCode int
		error      error
	}{
		{
			desc:       "Should successfully fetch 10 results",
			page:       1,
			respFile:   "testdata/success.json",
			resLen:     10,
			statusCode: http.StatusOK,
		},
		{
			desc:       "Should successfully fetch no results",
			page:       1000000000,
			respFile:   "testdata/empty.json",
			statusCode: http.StatusOK,
		},
		{
			desc:       "Should fail with error due to invalid params",
			respFile:   "testdata/empty.json",
			statusCode: http.StatusInternalServerError,
			error:      fmt.Errorf("failed to retrive results, %s", http.StatusText(http.StatusInternalServerError)),
		},
		{
			desc:       "Should fail to unmarshal result due to invalid json response",
			page:       1,
			respFile:   "testdata/invalid.json",
			statusCode: http.StatusOK,
			error:      fmt.Errorf("invalid character '}' looking for beginning of value"),
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

			result, err := s.Get(tc.page)
			if err != nil && err.Error() != tc.error.Error() {
				t.Fatal(err)
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
