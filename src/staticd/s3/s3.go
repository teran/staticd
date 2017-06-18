package s3

import (
	"github.com/minio/minio-go"
	log "github.com/sirupsen/logrus"

	"staticd/config"
)

var Client *minio.Client

func Connect(cfg config.Config) *minio.Client {
	c, err := minio.New(cfg.S3Endpoint, cfg.S3AccessKey, cfg.S3SecretKey, cfg.S3UseSSL)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Printf("Successfully created S3 client")

	return c
}
