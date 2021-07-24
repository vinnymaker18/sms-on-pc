package main

import (
	"encoding/json"
	"net/http"
)

const (
	contentTypeHeader = "Content-Type"

	jsonContentType = "application/json"

	// Accepts server control commands (shutdown, restart etc...) on this port.
	// Useful for building tools to manage servers.
	controlPort = 8001
)

func main() {
	http.HandleFunc("/sms", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set(contentTypeHeader, "text/plain")
		req.ParseForm()

		if req.Method == http.MethodGet {

			w.Write([]byte("As"))
		} else if req.Method == http.MethodPost {
			decoder := json.NewDecoder(req.Body)
			reqParams := make(map[string]string)
			decoder.Decode(&reqParams)

			w.Write([]byte("okay!"))
		}
	})

	http.ListenAndServe(":8000", nil)
}
