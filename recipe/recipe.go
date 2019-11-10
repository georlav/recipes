package recipe

type Recipes []Recipe

type Recipe struct {
	Title       string
	Ingredients string
	PageFound   int
}

//type Recipes struct {
//	values []Recipe
//	m      sync.RWMutex
//}
//
//func (r *Recipes) Append(recipe ...Recipe) {
//	r.m.Lock()
//	defer r.m.Unlock()
//
//	r.values = append(r.values, recipe...)
//}
//
//func (r *Recipes) Print(recipe ...Recipe) {
//	r.m.RLock()
//	defer r.m.RUnlock()
//
//	for _, recipe := range r.values {
//		fmt.Printf(`Recipe: %s, Inghedients: %s, found at page %d`,
//			recipe.Title, recipe.Ingredients, recipe.PageFound,
//		)
//	}
//}
