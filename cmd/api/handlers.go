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
	Log    LogEntry    `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogEntry struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
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
	case "logger":
		log.Debug().Msg("handle log case")
		log.Debug().Msgf("Log.Name:(%s) Log.Data:(%s)", requestPayload.Log.Name, requestPayload.Log.Data)
		// writeLog(w, requestPayload.Log, s.Conf.LogService)
		s.logEventViaRebbit(w, requestPayload.Log)
	case "mailer":
		log.Debug().Msg("handle mail case")
		sendMail(w, requestPayload.Mail, s.Conf.MailService)
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

func writeLog(w http.ResponseWriter, logs LogEntry, logService string) {
	log.Debug().Msg("post log into logger service")
	jsonData, _ := json.MarshalIndent(logs, "", "\t")
	logURL := fmt.Sprintf("%s/log", logService)
	log.Debug().Msgf("logURL: %s", logURL)
	request, err := http.NewRequest("POST", logURL, bytes.NewBuffer(jsonData))
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

	// if response.StatusCode == http.StatusUnauthorized {
	// 	errorJSON(w, errors.New("invalid credentials"))
	// 	return
	// }
	if response.StatusCode != http.StatusAccepted {
		errorJSON(w, errors.New("error calling logger service"))
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
		errorJSON(w, errors.New(jsonFromService.Massage))
		return
	}

	payload := &jsonResponse{
		Error:   false,
		Massage: "logged!",
		Data:    jsonFromService.Data,
	}

	writeJSON(w, http.StatusAccepted, payload)
}

func sendMail(w http.ResponseWriter, mail MailPayload, mailService string) {
	log.Debug().Msg("send mail using mail service")
	jsonData, _ := json.MarshalIndent(mail, "", "\t")
	mailURL := fmt.Sprintf("%s/send", mailService)
	log.Debug().Msgf("mailURL: %s", mailURL)
	request, err := http.NewRequest("POST", mailURL, bytes.NewBuffer(jsonData))
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

	// if response.StatusCode == http.StatusUnauthorized {
	// 	errorJSON(w, errors.New("invalid credentials"))
	// 	return
	// }
	if response.StatusCode != http.StatusAccepted {
		errorJSON(w, errors.New("error calling mail service"))
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
		errorJSON(w, errors.New(jsonFromService.Massage))
		return
	}

	payload := &jsonResponse{
		Error:   false,
		Massage: "message sent successfully!",
		Data:    jsonFromService.Data,
	}

	writeJSON(w, http.StatusAccepted, payload)
}

func (s *Server) logEventViaRebbit(w http.ResponseWriter, l LogEntry) {
	err := s.pushToQueue(l.Name, l.Data)
	if err != nil {
		errorJSON(w, err)
		return
	}

	payload := &jsonResponse{
		Error:   false,
		Massage: "logged to event to rabbitmq",
	}

	writeJSON(w, http.StatusAccepted, payload)
}

func (s *Server) pushToQueue(name, msg string) error {
	payload := LogEntry{
		Name: name,
		Data: msg,
	}

	j, err := json.MarshalIndent(&payload, "", "\t")
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal type LogEntry to json")
		return err
	}
	log.Debug().Msgf("pushToQueue LogEntry: %s", j)
	err = s.Producer.Push(string(j), "log.INFO")
	if err != nil {
		log.Error().Err(err).Msg("failed to push event to rabbitmq")
		return err
	}

	return nil
}
