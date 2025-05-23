package server

import (
	"fmt"
	"io"
	"net/http"
)


func QueryHandler(w http.ResponseWriter, req *http.Request) {
	// Check that method is POST
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Make sure content type text/plain
	if req.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Unsupported Content-Type. Expected text/plain for SQL.", http.StatusUnsupportedMediaType)
		return
	}
	
	// Read request body
	sqlBytes, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
	}

	sqlQuery := string(sqlBytes)

	// Print out original SQL for now
	fmt.Printf("Received SQL query: %s\n", sqlQuery)
	fmt.Fprintf(w, "Received your SQL query: %s\n", sqlQuery)
}


