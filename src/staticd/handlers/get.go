package handlers

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"staticd/config"
	"staticd/helpers"
	"staticd/s3"
)

func Get(w http.ResponseWriter, r *http.Request) {
	objectName := r.URL.Path[1:]

	log.WithFields(log.Fields{
		"method": "GET",
		"path":   objectName,
	}).Info("Incoming request")

	if objectName == "" || strings.HasSuffix(objectName, "/") {
		GetDirectory(w, r)
	} else {
		GetFile(w, r)
	}
}

func GetDirectory(w http.ResponseWriter, r *http.Request) {
	objectName := r.URL.Path[1:]

	doneCh := make(chan struct{})
	defer close(doneCh)
	objects := s3.Client.ListObjects(config.Cfg.S3BucketName, objectName, false, doneCh)
	w.Write([]byte(`<html><head><title>Index of /` + objectName + `</title></head><body bgcolor="white"><h1>Index of /` + objectName + `</h1><hr><pre><a href="../">../</a><br>`))
	for object := range objects {
		if object.Err != nil {
			log.WithFields(log.Fields{
				"method": "GET",
				"path":   objectName,
			}).Warn(object.Err)
			return
		}
		if object.Size == 0 && object.LastModified.IsZero() {
			w.Write([]byte(helpers.PadLink(path.Base(object.Key), "/"+object.Key, 45) + helpers.PadText("-", 20) + `      -<br>`))
		} else {
			w.Write([]byte(helpers.PadLink(path.Base(object.Key), "/"+object.Key, 45) + helpers.PadText(object.LastModified.Format(time.RFC3339), 20) + strconv.FormatInt(object.Size, 10) + `<br>`))
		}
	}
	w.Write([]byte(`</pre><hr><center>staticd</center></body></html>`))
}

func GetFile(w http.ResponseWriter, r *http.Request) {
	objectName := r.URL.Path[1:]

	if config.Cfg.S3Mode == "redirect" {
		_, err := s3.Client.StatObject(config.Cfg.S3BucketName, objectName)
		if err != nil {
			log.WithFields(log.Fields{
				"method": "GET",
				"path":   objectName,
			}).Warn(err.Error())
			http.Error(w, http.StatusText(404), 404)
			return
		}

		reqParams := make(url.Values)
		reqParams.Set("response-content-disposition", "attachment; filename=\""+objectName+"\"")

		presignedURL, err := s3.Client.PresignedGetObject(config.Cfg.S3BucketName, objectName, config.Cfg.S3RedirectUrlTTL, reqParams)
		if err != nil {
			log.WithFields(log.Fields{
				"method": "GET",
				"path":   objectName,
			}).Warn(err.Error())
			http.Error(w, http.StatusText(503), 503)
			return
		}

		http.Redirect(w, r, presignedURL.String(), http.StatusFound)
		log.WithFields(log.Fields{
			"method":   "GET",
			"path":     objectName,
			"redirect": presignedURL,
		}).Info("Sent to client")
		return
	} else if config.Cfg.S3Mode == "proxy" {
		objectStat, err := s3.Client.StatObject(config.Cfg.S3BucketName, objectName)
		if err != nil {
			log.WithFields(log.Fields{
				"method": "GET",
				"path":   objectName,
			}).Warn(err.Error())
			http.Error(w, http.StatusText(404), 404)
			return
		}

		w.Header().Set("Content-Type", objectStat.ContentType)
		w.Header().Set("Content-Length", strconv.FormatInt(objectStat.Size, 10))
		w.Header().Set("Last-Modified", objectStat.LastModified.Format(http.TimeFormat))
		w.Header().Set("Etag", objectStat.ETag)

		object, err := s3.Client.GetObject(config.Cfg.S3BucketName, objectName)
		if err != nil {
			log.WithFields(log.Fields{
				"method": "GET",
				"path":   objectName,
			}).Warn(err.Error())
			http.Error(w, http.StatusText(503), 503)
			return
		}
		content, err := ioutil.ReadAll(object)
		if err != nil {
			log.WithFields(log.Fields{
				"method": "GET",
				"path":   objectName,
			}).Warn(err.Error())
			http.Error(w, http.StatusText(503), 503)
			return
		}
		w.Write([]byte(content))
		log.WithFields(log.Fields{
			"method": "GET",
			"path":   objectName,
		}).Info("Sent to client")
		return
	}

	http.Error(w, http.StatusText(503), 503)
	log.WithFields(log.Fields{
		"method": "GET",
		"path":   objectName,
	}).Warn("Somehing wrong happend on server side, probably it's configuration issue.")
	return
}
