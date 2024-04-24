package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Suppliers map[string]SupplierConfig `yaml:"suppliers"`
}

type SupplierConfig struct {
	URL string `yaml:"url"`
}

// LoadConfig reads and parses the YAML configuration from a file.
func LoadConfig(filename string) (*Config, error) {
	var cfg Config
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
