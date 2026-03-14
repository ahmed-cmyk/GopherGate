package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ahmed-cmyk/GopherGate/internal/config"
)

func main() {
	// Create a context that is cancelled when SIGINT or SIGNTERM is received
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var config config.Config

	err := config.LoadData("config.yaml")
	if err != nil {
		log.Fatalf("Error unmarshaling YAML: %v\n", err)
		os.Exit(1)
	}

	// Run server inside a goroutine so that it doesn't block
	go func() {
		port := fmt.Sprintf(":%s", config.Server.Port)
		err := http.ListenAndServe(port, nil)
		if err != nil {
			log.Fatalf("Failed to start server: %v\n", err)
		}
	}()

	log.Printf("Starting service: %s\n", config.ServiceName)
	log.Printf("Listening on port %s\n", config.Server.Port)

	// Wait for the interrupt signal
	<-ctx.Done()

	log.Println("Shutting down server gracefully...")
	log.Println("Server gracefully stopped")
}
