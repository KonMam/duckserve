package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"duckserve/internal/engine"
)

var (
	querySemaphore chan struct{}
	queryTimeout   time.Duration
)

func InitializeServerSettings(maxConcurrency int, timeout time.Duration) {
	querySemaphore = make(chan struct{}, maxConcurrency)
	queryTimeout = timeout

	for range make([]struct{}, maxConcurrency) {
		querySemaphore <- struct{}{}
	}
	fmt.Printf("Server initialized with MaxConcurrency %d, queryTimeout %s\n", maxConcurrency, timeout)
}

func QueryHandler(w http.ResponseWriter, req *http.Request, db *engine.DB) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if req.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Unsupported Content-Type. Expected text/plain for SQL.", http.StatusUnsupportedMediaType)
		return
	}

	sqlBytes, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
	}

	sqlQuery := strings.TrimSpace(string(sqlBytes))

	select {
	case <-querySemaphore:
		defer func() {
			querySemaphore <- struct{}{}
			fmt.Println("Semaphore slot released.")
		}()
		fmt.Println("Semaphore slot acquired. Executing query...")

		ctx, cancel := context.WithTimeout(req.Context(), queryTimeout)
		defer cancel()

		resultChan := make(chan struct {
			res string
			err error
		}, 1)

		go func() {
			duckDBResult, duckDBErr := db.ExecuteQuery(sqlQuery)
			resultChan <- struct {
				res string
				err error
			}{duckDBResult, duckDBErr}
		}()

		select {
		case <-ctx.Done():
			fmt.Printf("Query execution context timed out for query '%s': %v\n", sqlQuery, ctx.Err())
			http.Error(w, fmt.Sprintf("Query timed out after %s", queryTimeout), http.StatusRequestTimeout)
			return
		case res := <-resultChan:
			if res.err != nil {
				fmt.Printf("DuckDB query execution failed for '%s': %v\n", sqlQuery, res.err)
				http.Error(w, fmt.Sprintf("DuckDB query failed: %v", res.err), http.StatusInternalServerError)
				return
			}
			fmt.Printf("Finished processing SQL query: %s\n", sqlQuery)
			fmt.Fprintf(w, "%s", res.res)
		}

	case <-time.After(5 * time.Second):
		fmt.Println("Failed to acquire semaphore slot within 5 seconds. Server busy.")
		http.Error(w, "Server too busy. Try again later.", http.StatusServiceUnavailable)
		return
	}
}
