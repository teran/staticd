package config

import (
	"time"
)

type Config struct {
	AllowAutoindex   bool          `default:"true"`
	AllowDelete      bool          `default:"false"`
	AllowHead        bool          `default:"true"`
	AllowGet         bool          `default:"true"`
	AllowPut         bool          `default:"false"`
	Debug            bool          `required:"false"`
	Listen           string        `default:":8080"`
	ListenMonitoring string        `default:":8081"`
	MaxUploadSize    int           `default:"1024"` // MB
	S3Endpoint       string        `required:"true"`
	S3AccessKey      string        `required:"true"`
	S3SecretKey      string        `required:"true"`
	S3UseSSL         bool          `default:"true"`
	S3BucketName     string        `required:"true"`
	S3Region         string        `required:"true"`
	S3Mode           string        `default:"proxy"`
	S3RedirectUrlTTL time.Duration `default:"1800s"`
}

var Cfg Config
