package api

import (
	"net/http"
)

// Broker api Handler
func Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Massage: "Hello from Broker!",
	}

	_ = writeJSON(w, http.StatusAccepted, payload)
}
