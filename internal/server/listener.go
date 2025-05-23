package server

import (
	"io"
	"net/http"
)

func StartListener() error {
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello world!\n")
	}

	http.HandleFunc("/query", helloHandler)
	return http.ListenAndServe(":8080", nil)
}
