package server

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

var querySemaphore chan struct{}
var queryTimeout time.Duration

func InitializeServerSettings(maxConcurrency int, timeout time.Duration) {
	querySemaphore = make(chan struct{}, maxConcurrency)
	queryTimeout = timeout

	for range make([]struct{}, maxConcurrency) {
		querySemaphore <- struct{}{}
	}
	fmt.Printf("Server initialized with MaxConcurrency %d, queryTimeout %s\n", maxConcurrency, timeout)
}


func QueryHandler(w http.ResponseWriter, req *http.Request) {
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

	sqlQuery := string(sqlBytes)

	select {
	case <- querySemaphore:
		defer func() {
			querySemaphore <- struct{}{}
			fmt.Println("Semaphore slot released.")
		}()
		fmt.Println("Semaphore slot acquired. Executing query...")
		
		select {
		case <- time.After(2 * time.Second):
			fmt.Printf("Received SQL Query: %s\n", sqlQuery)
			fmt.Fprintf(w, "Query proccessed (simulated): %s\n", sqlQuery)
		case <- time.After(queryTimeout):
			fmt.Printf("Query timed out after %s: %s\n", queryTimeout, sqlQuery)
			http.Error(w, fmt.Sprintf("Query timed out after %s: %s\n", queryTimeout, sqlQuery), http.StatusRequestTimeout)
			return
		}


	case <- time.After(5 * time.Second):
		fmt.Println("Failed to acquire semaphore slot within 5 seconds. Server busy.")
		http.Error(w, "Server too busy. Try again later.", http.StatusServiceUnavailable)
		return
	}
}


