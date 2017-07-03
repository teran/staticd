package monitoring

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	"staticd/config"
)

var (
	Requests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "staticd",
		Subsystem: "http",
		Name:      "requests",
		Help:      "Number of requests received",
	}, []string{"method"})
)

func RunPrometheusServer(cfg config.Config) {
	monitoringMux := http.NewServeMux()
	monitoringMux.Handle("/metrics", prometheus.Handler())
	monitoringMux.HandleFunc("/ping", PingHandler)

	prometheus.MustRegister(
		Requests,
	)

	log.Infof("Listening monitoring server on %v", config.Cfg.ListenMonitoring)
	log.Fatal(http.ListenAndServe(cfg.ListenMonitoring, monitoringMux))
}
