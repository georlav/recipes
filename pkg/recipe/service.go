package recipe

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/georlav/recipes/pkg/config"
)

type Service struct {
	client *http.Client
	cfg    config.RecipePuppyAPI
}

func NewService(cfg config.RecipePuppyAPI) *Service {
	return &Service{
		client: &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		},
		cfg: cfg,
	}
}

func (s Service) Get(qp QueryParams) (*ResultsResponse, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(`%s/%s?%s`, s.cfg.URL, "api", qp.Encode()),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	// nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(`status: %s, %w`, resp.Status, ErrNoResults)
	}

	response := ResultsResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf(`status: %s, %w`, resp.Status, ErrUnmarshalResponse)
	}

	return &response, nil
}
