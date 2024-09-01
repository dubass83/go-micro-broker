package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Massage string `json:"massage"`
	Data    any    `json:"data,omitempty"`
}

// readJSON reads from request body 1 Mb size and save to the data
func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // 1 Mb

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only one JSON value")
	}
	return nil
}

// writeJSON write data with status and optional headers into http.ResponseWriter
func writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "aplication/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	payload := jsonResponse{
		Error:   true,
		Massage: err.Error(),
	}

	return writeJSON(w, statusCode, payload)
}
