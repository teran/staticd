package handlers

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"staticd/config"
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

	doneCh := make(chan struct{})
	defer close(doneCh)
	objects := s3.Client.ListObjects(config.Cfg.S3BucketName, objectName, false, doneCh)
	w.Write([]byte(`<html><head><title>Index of /` + objectName + `</title></head><body bgcolor="white"><h1>Index of /` + objectName + `</h1><hr><pre><a href="../">../</a><br>`))
	for object := range objects {
		if object.Err != nil {
			log.Println(object.Err)
			return
		}
		if object.Size == 0 && object.LastModified.IsZero() {
			w.Write([]byte(`<a href="/` + object.Key + `">` + path.Base(object.Key) + `/</a>                                        -            -<br>`))
		} else {
			w.Write([]byte(`<a href="/` + object.Key + `">` + path.Base(object.Key) + `</a>                                        ` + object.LastModified.Format(time.RFC3339) + `            ` + strconv.FormatInt(object.Size, 10) + `<br>`))
		}
	}
	w.Write([]byte(`</pre><hr><center>staticd</center></body></html>`))
}

func GetFile(w http.ResponseWriter, r *http.Request) {
	objectName := r.URL.Path[1:]

	if config.Cfg.S3Mode == "redirect" {
		_, err := s3.Client.StatObject(config.Cfg.S3BucketName, objectName)
		if err != nil {
			log.Printf("GET %v: %v", objectName, err.Error())
			http.Error(w, http.StatusText(404), 404)
			return
		}

		reqParams := make(url.Values)
		reqParams.Set("response-content-disposition", "attachment; filename=\""+objectName+"\"")

		presignedURL, err := s3.Client.PresignedGetObject(config.Cfg.S3BucketName, objectName, config.Cfg.S3RedirectUrlTTL, reqParams)
		if err != nil {
			log.Printf("GET %v: %v", objectName, err.Error())
			http.Error(w, http.StatusText(503), 503)
			return
		}

		http.Redirect(w, r, presignedURL.String(), http.StatusFound)
		log.Printf("GET %v: redirected to %v", objectName, presignedURL)
		return
	} else if config.Cfg.S3Mode == "proxy" {
		objectStat, err := s3.Client.StatObject(config.Cfg.S3BucketName, objectName)
		if err != nil {
			log.Printf("GET %v: %v", objectName, err.Error())
			http.Error(w, http.StatusText(404), 404)
			return
		}

		w.Header().Set("Content-Type", objectStat.ContentType)
		w.Header().Set("Content-Length", strconv.FormatInt(objectStat.Size, 10))
		w.Header().Set("Last-Modified", objectStat.LastModified.Format(http.TimeFormat))
		w.Header().Set("Etag", objectStat.ETag)

		object, err := s3.Client.GetObject(config.Cfg.S3BucketName, objectName)
		if err != nil {
			log.Printf("GET %v: %v", objectName, err.Error())
			http.Error(w, http.StatusText(503), 503)
			return
		}
		content, err := ioutil.ReadAll(object)
		if err != nil {
			log.Printf("GET %v: %v", objectName, err.Error())
			http.Error(w, http.StatusText(503), 503)
			return
		}
		w.Write([]byte(content))
		log.Printf("GET %v: sent to client", objectName)
		return
	}

	http.Error(w, http.StatusText(503), 503)
	log.Printf("GET %v: unknown request", objectName)
	return
}
