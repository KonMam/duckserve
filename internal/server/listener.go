package server

import (
	"net/http"
)

func StartListener() error {
	http.HandleFunc("/query", QueryHandler)
	return http.ListenAndServe(":8080", nil)
}
