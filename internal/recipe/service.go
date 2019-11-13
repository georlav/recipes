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

func (s Service) Get(qp QueryParams) (rr ResultsResponse, err error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(`%s/%s?%s`, s.cfg.Host, "api", qp.Encode()),
		nil,
	)
	if err != nil {
		return rr, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return rr, err
	}
	// nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return rr, fmt.Errorf("failed to retrieve results, %s", resp.Status)
	}
	if err := json.NewDecoder(resp.Body).Decode(&rr); err != nil {
		return rr, fmt.Errorf("failed to unmarshal response, %w", err)
	}

	return rr, nil
}
