package handlers

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"staticd/config"
	"staticd/s3"
)

func Delete(w http.ResponseWriter, r *http.Request) {
	objectName := r.URL.Path[1:]

	err := s3.Client.RemoveObject(config.Cfg.S3BucketName, objectName)
	if err != nil {
		log.WithFields(log.Fields{
			"remote": r.RemoteAddr,
			"method": r.Method,
			"path":   "/" + r.URL.Path[1:],
			"return": "500 " + http.StatusText(500),
		}).Warn(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	log.WithFields(log.Fields{
		"remote": r.RemoteAddr,
		"method": "DELETE",
		"path":   "/" + objectName,
	}).Info("Object successfully deleted")
	http.Error(w, http.StatusText(204), 204)
	return
}
