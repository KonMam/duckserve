package main

import (
	"fmt"
	"log"

	"duckserve/internal/server"
)

func main() {
	fmt.Println("Starting DuckServe on :8080.")

	err := server.StartListener()
	if err != nil {
		log.Fatalf("Server failed to start: %v!", err)
	}
}
