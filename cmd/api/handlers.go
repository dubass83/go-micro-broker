package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Broker api Handler
func Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Massage: "Hello from Broker!",
	}

	_ = writeJSON(w, http.StatusAccepted, payload)
}

func (s *Server) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := readJSON(w, r, &requestPayload)
	if err != nil {
		errorJSON(w, err, http.StatusBadRequest)
		return
	}
	switch requestPayload.Action {
	case "auth":
		log.Debug().Msg("handle auth case")
		authenticate(w, requestPayload.Auth, s.Conf.AuthService)
	default:
		errorJSON(w, errors.New("unknown action"))
	}
}

func authenticate(w http.ResponseWriter, auth AuthPayload, authService string) {
	log.Debug().Msg("run authenticate method")
	jsonData, _ := json.MarshalIndent(auth, "", "\t")
	authURL := fmt.Sprintf("%s/authenticate", authService)
	log.Debug().Msgf("authURL: %s", authURL)
	request, err := http.NewRequest("POST", authURL, bytes.NewBuffer(jsonData))
	if err != nil {
		errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		errorJSON(w, errors.New("invalid credentials"))
		return
	}
	if response.StatusCode != http.StatusAccepted {
		errorJSON(w, errors.New("error calling auth service"))
		return
	}

	var jsonFromService jsonResponse
	maxBytes := 1048576 // 1 Mb

	response.Body = http.MaxBytesReader(w, response.Body, int64(maxBytes))
	dec := json.NewDecoder(response.Body)
	err = dec.Decode(&jsonFromService)
	log.Debug().Msgf("jsonFromService: %+v", jsonFromService)
	if err != nil {
		errorJSON(w, errors.New(err.Error()))
		log.Error().Err(err)
		return
	}

	if jsonFromService.Error {
		errorJSON(w, errors.New(jsonFromService.Massage), http.StatusUnauthorized)
		return
	}

	payload := &jsonResponse{
		Error:   false,
		Massage: "Authenticated!",
		Data:    jsonFromService.Data,
	}

	writeJSON(w, http.StatusAccepted, payload)
}
