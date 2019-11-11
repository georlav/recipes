package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Config object
type Config struct {
	APP            APP            `json:"app"`
	RecipePuppyAPI RecipePuppyAPI `json:"recipePuppyAPI"`
}

// APP hold generic app configuration
type APP struct {
	NumOfWorkers int  `json:"numOfWorkers"`
	NumOfPages   int  `json:"numOfPages"`
	Debug        bool `json:"debug"`
}

// RecipePuppyAPI holds configuration for recipe puppy api
type RecipePuppyAPI struct {
	Host    string `json:"host"`
	Timeout int64  `json:"timeout"`
}

// Load loads a json config file and returns a config object
func Load(cfgFile string) (cfg *Config, err error) {
	b, err := ioutil.ReadFile("./" + cfgFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file %s, %w", cfgFile, err)
	}

	if err := json.Unmarshal(b, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file %s, %w", cfgFile, err)
	}

	return cfg, nil
}
