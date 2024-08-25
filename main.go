package main

import (
	"log"
	"net/http"

	"github.com/dubass83/go-micro-broker/cmd/api"
)

func main() {
	s := api.CreateNewServer()
	s.ConfigureCORS()
	s.AddMiddleware()
	s.MountHandlers()
	err := http.ListenAndServe(":8080", s.Router)
	if err != nil {
		log.Fatal(err)
	}
}
