package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ServiceName string `yaml:"service_name"`
	Server      struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
}

func (c *Config) LoadData(path string) error {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	return yaml.Unmarshal(data, c)
}
