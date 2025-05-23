package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	serverURL   = flag.String("url", "http://localhost:8080/query", "URL of the DuckServe /query endpoint")
	numRequests = flag.Int("n", 10, "Total number of requests to send")
	concurrency = flag.Int("c", 5, "Number of concurrent client requests")
	clientTimeout = flag.Int("t", 15, "Client-side HTTP request timeout in seconds")
)

func main() {
	flag.Parse()

	if *numRequests <= 0 {
		log.Fatalf("Number of requests (-n) must be greater than 0.")
	}

	if *concurrency <= 0 {
		log.Fatalf("Concurrency (-c) must be greater than 0.")
	}

	fmt.Printf("Starting simple DuckServe Load Tester\n")
	fmt.Printf("-----------------------------------\n")
	fmt.Printf("Target URL: %s\n", *serverURL)
	fmt.Printf("Total Requests: %d\n", *numRequests)
	fmt.Printf("Client Concurrency: %d\n", *concurrency)
	fmt.Println("-----------------------------------")

	var wg sync.WaitGroup

	requestLimiter := make(chan struct{}, *concurrency)

	httpClient := &http.Client{
		Timeout: time.Duration(*clientTimeout) * time.Second,
	}

	for i := 1; i <= *numRequests; i++ {
		requestLimiter <- struct{}{}
		wg.Add(1)

		go func(requestID int) {
			defer func() {
				<- requestLimiter
				wg.Done()
			}()

			sqlQuery := fmt.Sprintf("SELECT %d;", requestID)
			requstStartTime := time.Now()

			fmt.Printf("[%d] Sending query: %s\n", requestID, sqlQuery)

			resp, err := sendQuery(httpClient, sqlQuery)
			requestDuration := time.Since(requstStartTime)

			if err != nil {
				log.Printf("[%d] ERROR: %v (took %s)\n", requestID, err, requestDuration)
				return
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("[%d] Response Status: %d, Body: \"%s\" (took %s)\n",
				requestID, resp.StatusCode, string(body), requestDuration)
		}(i)
	}

	wg.Wait()

	fmt.Println("\n-----------------------------------")
	fmt.Println("Load test finished.")

}



func sendQuery(client *http.Client, query string) (*http.Response, error ) {
	req, err := http.NewRequest(http.MethodPost, *serverURL, bytes.NewBufferString(query))
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "text/plain")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Client request failed: %w", err)
	}
	return resp, nil
}
