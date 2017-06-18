package handlers

import (
	"log"
	"net/http"
	"strconv"

	"staticd/config"
	"staticd/s3"
)

func Put(w http.ResponseWriter, r *http.Request) {
	objectName := r.URL.Path[1:]
	objectSize, err := strconv.Atoi(r.Header["Content-Length"][0])

	if err != nil {
		log.Printf("PUT %v: %v", objectName, err.Error())
		http.Error(w, http.StatusText(400), 400)
		return
	}

	if objectSize > config.Cfg.MaxUploadSize*1024*1024 {
		log.Printf("PUT %v: object size is higher than allowed limit. HTTP 413 code returned", objectName)
		http.Error(w, http.StatusText(413), 413)
		return
	}

	_, err = s3.Client.PutObject(config.Cfg.S3BucketName, objectName, r.Body, "application/octet-stream")
	if err != nil {
		log.Printf("PUT %v: %v", objectName, err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
}
