package handlers

import (
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"staticd/config"
	"staticd/s3"
)

func Put(w http.ResponseWriter, r *http.Request) {
	objectName := r.URL.Path[1:]

	objectSize, err := strconv.Atoi(r.Header["Content-Length"][0])
	if err != nil {
		log.WithFields(log.Fields{
			"remote": r.RemoteAddr,
			"method": "PUT",
			"path":   "/" + objectName,
			"return": "400 " + http.StatusText(400),
		}).Warn(err.Error())
		http.Error(w, http.StatusText(400), 400)
		return
	}

	if objectSize > config.Cfg.MaxUploadSize*1024*1024 {
		log.WithFields(log.Fields{
			"remote": r.RemoteAddr,
			"method": "PUT",
			"path":   "/" + objectName,
			"return": "413 " + http.StatusText(413),
		}).Warn("Object is bigger than allowed by server configuration")
		http.Error(w, http.StatusText(413), 413)
		return
	}

	if config.Cfg.S3Mode == "redirect" {
		presignedURL, err := s3.Client.PresignedPutObject(config.Cfg.S3BucketName, objectName, config.Cfg.S3RedirectUrlTTL)
		if err != nil {
			log.WithFields(log.Fields{
				"remote": r.RemoteAddr,
				"method": "PUT",
				"path":   "/" + objectName,
				"return": "503 " + http.StatusText(503),
			}).Warn(err.Error())
			http.Error(w, http.StatusText(503), 503)
			return
		}

		http.Redirect(w, r, presignedURL.String(), http.StatusFound)
		log.WithFields(log.Fields{
			"remote":   r.RemoteAddr,
			"method":   "PUT",
			"path":     "/" + objectName,
			"redirect": presignedURL,
		}).Info("Sent to client")
		return
	}

	_, err = s3.Client.PutObject(config.Cfg.S3BucketName, objectName, r.Body, "application/octet-stream")
	if err != nil {
		log.WithFields(log.Fields{
			"remote": r.RemoteAddr,
			"method": "PUT",
			"path":   "/" + objectName,
			"return": "500 " + http.StatusText(500),
		}).Warn(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	log.WithFields(log.Fields{
		"remote": r.RemoteAddr,
		"method": "PUT",
		"path":   "/" + objectName,
	}).Info("Object successfully created")
	http.Error(w, http.StatusText(201), 201)
	return
}
