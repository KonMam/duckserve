package main

import (
	"fmt"
	"log"

	"duckserve/internal/config"
	"duckserve/internal/server"
)

func main() {
	fmt.Println("Starting DuckServe on :8080.")

	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	server.InitializeServerSettings(cfg.MaxConcurrency, cfg.GetQueryTimeout())

	err = server.StartListener()
	if err != nil {
		log.Fatalf("Server failed to start: %v!", err)
	}
}
