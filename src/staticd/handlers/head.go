package handlers

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"staticd/config"
	"staticd/s3"
)

func Head(w http.ResponseWriter, r *http.Request) {
	objectName := r.URL.Path[1:]

	if objectName == "" || strings.HasSuffix(objectName, "/") {
		HeadDirectory(w, r)
	} else {
		HeadFile(w, r)
	}
}

func HeadDirectory(w http.ResponseWriter, r *http.Request) {
	if !config.Cfg.AllowAutoindex {
		http.Error(w, http.StatusText(403), 403)
		return
	}

	objectName := r.URL.Path[1:]

	doneCh := make(chan struct{})
	defer close(doneCh)
	objects := s3.Client.ListObjects(config.Cfg.S3BucketName, objectName, false, doneCh)
	for object := range objects {
		objectName := object.Key
		objectSize := strconv.FormatInt(object.Size, 10)

		log.WithFields(log.Fields{
			"remote": r.RemoteAddr,
			"method": "GET",
			"path":   "/" + objectName,
		}).Debugf("Listing objects from s3 backend: name=%v ; size=%v", objectName, objectSize)
		if object.Err != nil {
			log.WithFields(log.Fields{
				"remote": r.RemoteAddr,
				"method": "HEAD",
				"path":   "/" + objectName,
				"return": "503 " + http.StatusText(503),
			}).Warn(object.Err)
			http.Error(w, http.StatusText(503), 503)
			return
		}
	}
}

func HeadFile(w http.ResponseWriter, r *http.Request) {
	objectName := r.URL.Path[1:]

	objectStat, err := s3.Client.StatObject(config.Cfg.S3BucketName, objectName)
	if err != nil {
		log.WithFields(log.Fields{
			"remote": r.RemoteAddr,
			"method": "HEAD",
			"path":   "/" + objectName,
			"return": "404 " + http.StatusText(404),
		}).Warn(err.Error())
		http.Error(w, http.StatusText(404), 404)
		return
	}

	if config.Cfg.S3Mode == "redirect" {
		reqParams := make(url.Values)
		reqParams.Set("response-content-disposition", "attachment; filename=\""+objectName+"\"")

		presignedURL, err := s3.Client.PresignedGetObject(config.Cfg.S3BucketName, objectName, config.Cfg.S3RedirectUrlTTL, reqParams)
		if err != nil {
			log.WithFields(log.Fields{
				"remote": r.RemoteAddr,
				"method": "HEAD",
				"path":   "/" + objectName,
				"return": "503 " + http.StatusText(503),
			}).Warn(err.Error())
			http.Error(w, http.StatusText(503), 503)
			return
		}

		http.Redirect(w, r, presignedURL.String(), http.StatusFound)
		log.WithFields(log.Fields{
			"remote":   r.RemoteAddr,
			"method":   "HEAD",
			"path":     "/" + objectName,
			"redirect": presignedURL,
		}).Info("Sent to client")
		return
	} else if config.Cfg.S3Mode == "proxy" {
		w.Header().Set("Content-Type", objectStat.ContentType)
		w.Header().Set("Content-Length", strconv.FormatInt(objectStat.Size, 10))
		w.Header().Set("Last-Modified", objectStat.LastModified.Format(http.TimeFormat))
		w.Header().Set("Etag", objectStat.ETag)

		object, err := s3.Client.GetObject(config.Cfg.S3BucketName, objectName)
		if err != nil {
			log.WithFields(log.Fields{
				"remote": r.RemoteAddr,
				"method": "HEAD",
				"path":   "/" + objectName,
				"return": "503 " + http.StatusText(503),
			}).Warn(err.Error())
			http.Error(w, http.StatusText(503), 503)
			return
		}
		content, err := ioutil.ReadAll(object)
		if err != nil {
			log.WithFields(log.Fields{
				"remote": r.RemoteAddr,
				"method": "HEAD",
				"path":   "/" + objectName,
				"return": "503 " + http.StatusText(503),
			}).Warn(err.Error())
			http.Error(w, http.StatusText(503), 503)
			return
		}
		w.Write([]byte(content))
		log.WithFields(log.Fields{
			"remote": r.RemoteAddr,
			"method": "HEAD",
			"path":   "/" + objectName,
		}).Info("Sent to client")
		return
	}

	log.WithFields(log.Fields{
		"remote": r.RemoteAddr,
		"method": "HEAD",
		"path":   "/" + objectName,
		"return": "503 " + http.StatusText(503),
	}).Warn("Somehing wrong happend on server side, probably it's configuration issue.")
	http.Error(w, http.StatusText(503), 503)
	return
}
