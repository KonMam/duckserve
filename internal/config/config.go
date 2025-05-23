package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Config struct {
	MaxConcurrency int `json:"MaxConcurrency"`
	QueryTimeoutSec	int `json:"QueryTimeoutSec"`
}

func DefaultConfig() *Config {
	return &Config{
		MaxConcurrency: 4,
		QueryTimeoutSec: 30,
	}
}

func LoadConfig(filePath string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Config file not found '%s', using default config.\n", filePath)
			return cfg, nil
		}
		return nil, fmt.Errorf("Failed to read config file '%s': '%w'", filePath, err)
	}
	
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse config file, '%s': '%w'", filePath, err)
	}

	if cfg.MaxConcurrency <= 0 {
		return nil, fmt.Errorf("MaxConcurrency must be above 0.", filePath, err)
	}

	if cfg.QueryTimeoutSec <= 0 {
		return nil, fmt.Errorf("Max QueryTimeoutSec must be above 0.", filePath, err)
	}
	
	return cfg, nil
}


func (c *Config) GetQueryTimeout() time.Duration {
	return time.Duration(c.QueryTimeoutSec) * time.Second
}
