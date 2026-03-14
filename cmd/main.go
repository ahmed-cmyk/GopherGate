package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ahmed-cmyk/GopherGate/internal/config"
)

func main() {
	var config config.Config

	err := config.LoadData("config.yaml")
	if err != nil {
		log.Fatalf("Error unmarshaling YAML: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Starting service: %s\n", config.ServiceName)
	fmt.Printf("Listening on port %s\n", config.Server.Port)

	port := fmt.Sprintf(":%s", config.Server.Port)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
