package main

import (
	"net/http"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"

	"staticd/config"
	"staticd/handlers"
	"staticd/s3"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"method": r.Method,
		"path":   "/" + r.URL.Path[1:],
	}).Info("Incoming request")

	if r.Method == http.MethodGet && config.Cfg.AllowGet {
		handlers.Get(w, r)
	} else if r.Method == http.MethodPut && config.Cfg.AllowPut {
		handlers.Put(w, r)
	} else if r.Method == http.MethodDelete && config.Cfg.AllowDelete {
		handlers.Delete(w, r)
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
	log.Printf("Listening on %v", config.Cfg.Listen)
	http.ListenAndServe(config.Cfg.Listen, nil)
}
