package recipe

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/georlav/recipes/internal/config"
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

func (s Service) Get(page int) (rr RecipeResponse, err error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(`%s/%s?p=%d`, s.cfg.Host, "api", page),
		nil,
	)
	if err != nil {
		return rr, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return rr, err
	}

	if resp.StatusCode != http.StatusOK {
		return rr, fmt.Errorf("failed to retrive results, %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&rr); err != nil {
		return rr, fmt.Errorf("failed to unmarshal response, %w", err)
	}

	//for i := range rr.Results {
	//	r = append(r, Recipe{
	//		Title:       rr.Results[i].Title,
	//		Ingredients: rr.Results[i].Ingredients,
	//		PageFound:   page,
	//	})
	//}

	return rr, nil
}
