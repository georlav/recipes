package recipe

import (
	"fmt"
	"net/url"
)

// QueryParams object used by service to filter results
// Ingredients is a comma delimited ingredients
// Term you can search by text search term
// Page results are paginated you can requests a specific page using this param
// No parameters are required
type QueryParams struct {
	_           struct{}
	Ingredients string
	Term        string
	Page        int
}

func (q QueryParams) Encode() string {
	u := url.Values{}

	if q.Ingredients != "" {
		u.Set("i", q.Ingredients)
	}
	if q.Term != "" {
		u.Set("q", q.Term)
	}
	if q.Page != 0 {
		u.Set("p", fmt.Sprint(q.Page))
	}

	return u.Encode()
}
