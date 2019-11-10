package recipe

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/georlav/recipes/config"
)

type Service struct {
	client http.Client
	cfg    config.RecipePuppyAPI
}

func NewService(cfg config.RecipePuppyAPI) *Service {
	return &Service{
		client: http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		},
		cfg: cfg,
	}
}

func (s Service) GetRecipes(page int) (r Recipes, err error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(`%s/%s?p=%d`, s.cfg.Host, "api", page),
		nil,
	)
	if err != nil {
		return r, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return r, err
	}

	if resp.StatusCode != http.StatusOK {
		return r, fmt.Errorf("invalid response, %s, url: %s", resp.Status, req.URL)
	}

	rr := RecipeResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&rr); err != nil {
		return r, err
	}

	for i := range rr.Results {
		r = append(r, Recipe{
			Title:       rr.Results[i].Title,
			Ingredients: rr.Results[i].Ingredients,
			PageFound:   page,
		})
	}

	return r, nil
}
