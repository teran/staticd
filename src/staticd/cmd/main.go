package main

import (
	"net/http"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"

	"staticd/config"
	"staticd/handlers"
	"staticd/monitoring"
	"staticd/s3"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"remote": r.RemoteAddr,
		"method": r.Method,
		"path":   "/" + r.URL.Path[1:],
	}).Info("Incoming request")

	if r.Method == http.MethodGet && config.Cfg.AllowGet {
		handlers.Get(w, r)
	} else if r.Method == http.MethodHead && config.Cfg.AllowHead {
		handlers.Head(w, r)
	} else if r.Method == http.MethodPut && config.Cfg.AllowPut {
		handlers.Put(w, r)
	} else if r.Method == http.MethodDelete && config.Cfg.AllowDelete {
		handlers.Delete(w, r)
	} else {
		http.Error(w, http.StatusText(405), 405)
		log.WithFields(log.Fields{
			"remote": r.RemoteAddr,
			"method": r.Method,
			"path":   "/" + r.URL.Path[1:],
			"return": "405 Method not allowed",
		}).Warn("Method not allowed by configuration")
	}
	monitoring.Requests.With(prometheus.Labels{"method": r.Method}).Inc()
	monitoring.Requests.With(prometheus.Labels{"method": "total"}).Inc()

	return
}

func main() {
	log.SetFormatter(&log.TextFormatter{})
	err := envconfig.Process("staticd", &config.Cfg)
	if err != nil {
		log.Fatalln(err.Error())
	}

	if config.Cfg.Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	s3.Client = s3.Connect(config.Cfg)

	mainMux := http.NewServeMux()
	mainMux.HandleFunc("/", handler)
	log.Infof("Listening on %v", config.Cfg.Listen)

	go monitoring.RunPrometheusServer(config.Cfg)
	http.ListenAndServe(config.Cfg.Listen, mainMux)
}
