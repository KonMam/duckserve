package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"duckserve/internal/config"
	"duckserve/internal/engine"
	"duckserve/internal/server"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	duckDB, err := engine.NewDB()
	if err != nil {
		log.Fatalf("Failed to initialize duckDB: %v", err)
	}
	defer duckDB.Close()

	server.InitializeServerSettings(cfg.MaxConcurrency, cfg.GetQueryTimeout())

	mux := http.NewServeMux()

	mux.HandleFunc("/query", func (w http.ResponseWriter, r *http.Request) {
		server.QueryHandler(w, r, duckDB)
	})

	port := "8080"
	srv := &http.Server{
		Addr: ":" + port,
		Handler: mux,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout: 30 * time.Second,
	}

	go func() {
		log.Printf("Server starting on port %s...", port)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited cleanly.")
}
