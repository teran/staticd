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

	if objectName == "" || strings.HasSuffix(objectName, "/") {
		GetDirectory(w, r)
	} else {
		GetFile(w, r)
	}
}

func GetDirectory(w http.ResponseWriter, r *http.Request) {
	objectName := r.URL.Path[1:]

	if !config.Cfg.AllowAutoindex {
		log.WithFields(log.Fields{
			"remote": r.RemoteAddr,
			"method": "GET",
			"path":   "/" + objectName,
			"return": "403 " + http.StatusText(403),
		}).Warn("Attempt to access directory listing")
		http.Error(w, http.StatusText(403), 403)
		return
	}

	var fileList []string

	doneCh := make(chan struct{})
	defer close(doneCh)
	objects := s3.Client.ListObjects(config.Cfg.S3BucketName, objectName, false, doneCh)
	fileList = append(fileList, `<html><head><title>Index of /`+objectName+`</title></head><body bgcolor="white"><h1>Index of /`+objectName+`</h1><hr><pre><a href="../">../</a><br>`)
	for object := range objects {
		objectName := object.Key
		objectSize := strconv.FormatInt(object.Size, 10)
		objectLastModified := object.LastModified.Format(time.RFC3339)

		log.WithFields(log.Fields{
			"remote": r.RemoteAddr,
			"method": "GET",
			"path":   "/" + objectName,
		}).Debugf("Listing objects from s3 backend: name=%v ; size=%v", objectName, objectSize)
		if object.Err != nil {
			log.WithFields(log.Fields{
				"remote": r.RemoteAddr,
				"method": "GET",
				"path":   "/" + objectName,
				"return": "503 " + http.StatusText(503),
			}).Warn(object.Err)
			http.Error(w, http.StatusText(503), 503)
			return
		}
		if object.Size == 0 && object.LastModified.IsZero() {
			fileList = append(fileList, helpers.PadLink(path.Base(objectName), "/"+objectName, 45)+helpers.PadText("-", 20)+`      -<br>`)
		} else {
			fileList = append(fileList, helpers.PadLink(path.Base(objectName), "/"+objectName, 45)+helpers.PadText(objectLastModified, 20)+objectSize+`<br>`)
		}
	}
	fileList = append(fileList, `</pre><hr><center>staticd</center></body></html>`)
	w.Write([]byte(strings.Join(fileList, "")))
}

func GetFile(w http.ResponseWriter, r *http.Request) {
	objectName := r.URL.Path[1:]

	objectStat, err := s3.Client.StatObject(config.Cfg.S3BucketName, objectName)
	if err != nil {
		log.WithFields(log.Fields{
			"remote": r.RemoteAddr,
			"method": "GET",
			"path":   "/" + objectName,
			"return": "404 " + http.StatusText(404),
		}).Warn(err.Error())
		http.Error(w, http.StatusText(404), 404)
		return
	}

	log.WithFields(log.Fields{
		"remote": r.RemoteAddr,
		"method": "GET",
		"path":   "/" + objectName,
	}).Debugf("Stat object in s3 backend: name=%v ; size=%v", objectName, strconv.FormatInt(objectStat.Size, 10))

	if config.Cfg.S3Mode == "redirect" {
		reqParams := make(url.Values)
		reqParams.Set("response-content-disposition", "attachment; filename=\""+objectName+"\"")

		presignedURL, err := s3.Client.PresignedGetObject(config.Cfg.S3BucketName, objectName, config.Cfg.S3RedirectUrlTTL, reqParams)
		if err != nil {
			log.WithFields(log.Fields{
				"remote": r.RemoteAddr,
				"method": "GET",
				"path":   "/" + objectName,
				"return": "503 " + http.StatusText(503),
			}).Warn(err.Error())
			http.Error(w, http.StatusText(503), 503)
			return
		}

		http.Redirect(w, r, presignedURL.String(), http.StatusFound)
		log.WithFields(log.Fields{
			"remote":   r.RemoteAddr,
			"method":   "GET",
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
				"method": "GET",
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
				"method": "GET",
				"path":   "/" + objectName,
				"return": "503 " + http.StatusText(503),
			}).Warn(err.Error())
			http.Error(w, http.StatusText(503), 503)
			return
		}
		w.Write([]byte(content))
		log.WithFields(log.Fields{
			"remote": r.RemoteAddr,
			"method": "GET",
			"path":   "/" + objectName,
		}).Info("Sent to client")
		return
	}

	log.WithFields(log.Fields{
		"remote": r.RemoteAddr,
		"method": "GET",
		"path":   "/" + objectName,
		"return": "503 " + http.StatusText(503),
	}).Warn("Somehing wrong happend on server side, probably it's configuration issue.")
	http.Error(w, http.StatusText(503), 503)
	return
}
