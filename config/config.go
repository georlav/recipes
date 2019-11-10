package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	RecipePuppyAPI RecipePuppyAPI
}

type RecipePuppyAPI struct {
	Host         string
	Timeout      int64
	NumOfWorkers int
	NumOfPages   int
}

func Load(cfgFile string) (cfg *Config, err error) {
	b, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file %s, %w", cfgFile, err)
	}

	if err := json.Unmarshal(b, &Config{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file %s, %w", cfgFile, err)
	}

	return cfg, nil
}
