package main

import (
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"

	"staticd/config"
	"staticd/handlers"
	"staticd/s3"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		handlers.Get(w, r)
	} else {
		http.Error(w, http.StatusText(405), 405)
	}
}

func main() {
	err := envconfig.Process("staticd", &config.Cfg)
	if err != nil {
		log.Fatalln(err.Error())
	}

	s3.Client = s3.Connect(config.Cfg)

	http.HandleFunc("/", handler)
	log.Printf("Listening on %v\n", config.Cfg.Listen)
	http.ListenAndServe(config.Cfg.Listen, nil)
}
