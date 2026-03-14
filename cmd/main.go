package main

import (
	"fmt"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ServiceName string `yaml:"service_name"`
	Server      struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
}

func main() {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		os.Exit(1)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("Error unmarshaling YAML: %v\n", err)
		os.Exit(1)
	}

	port := fmt.Sprintf(":%s", config.Server.Port)

	fmt.Printf("Running on %s port %s\n", config.Server.Host, config.Server.Port)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
