package monitoring

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"staticd/config"
	"staticd/s3"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	s3.Client = s3.Connect(config.Cfg)

	exists, err := s3.Client.BucketExists(config.Cfg.S3BucketName)
	if err != nil || !exists {
		log.WithFields(log.Fields{
			"remote": r.RemoteAddr,
			"method": "GET",
			"server": "monitoring",
			"path":   "/ping",
			"return": "503 " + http.StatusText(503),
		}).Debug(err.Error())
		http.Error(w, http.StatusText(503), 503)
		return
	}

	log.WithFields(log.Fields{
		"remote": r.RemoteAddr,
		"method": "GET",
		"server": "monitoring",
		"path":   "/ping",
		"return": "200 " + http.StatusText(200),
	}).Debug("OK")
	http.Error(w, http.StatusText(200), 200)
	return
}
