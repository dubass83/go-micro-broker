package main

import (
	"log"
	"net/http"

	"github.com/dubass83/go-micro-broker/cmd/api"
)

const servicePort = ":8080"

func main() {
	s := api.CreateNewServer()
	s.ConfigureCORS()
	s.AddMiddleware()
	s.MountHandlers()
	log.Printf("start listening on the port %s\n", servicePort)
	err := http.ListenAndServe(servicePort, s.Router)
	if err != nil {
		log.Fatal(err)
	}
}
