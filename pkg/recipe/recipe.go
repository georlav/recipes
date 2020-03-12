package recipe

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

type Recipe struct {
	Title       string   `json:"title"`
	URL         string   `json:"url"`
	Thumbnail   string   `json:"thumbnail"`
	Ingredients []string `json:"ingredients"`
	PageFound   int
}

type Recipes struct {
	values []Recipe
	m      sync.RWMutex
}

func (r *Recipes) Append(recipes ...Recipe) {
	r.m.Lock()
	defer r.m.Unlock()

	r.values = append(r.values, recipes...)
}

func (r *Recipes) Values(recipes ...Recipe) []Recipe {
	r.m.RLock()
	defer r.m.RUnlock()

	return r.values
}

// Save results to file
func (r *Recipes) Save() error {
	b, err := json.Marshal(r.values)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile("recipes.json", b, 0644); err != nil {
		return err
	}

	return nil
}
