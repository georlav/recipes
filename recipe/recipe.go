package recipe

import "sync"

type Recipe struct {
	Title       string
	Ingredients string
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
