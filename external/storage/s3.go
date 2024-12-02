package storage

import (
	"GoFiber-API/internal/config"

	"github.com/gofiber/storage/s3/v2"
)

var S3 *s3.Storage

func InitStorage() {
	S3 = s3.New(s3.Config{
		Bucket:   config.GetConfig.S3_BUCKET,
		Endpoint: config.GetConfig.S3_ENDPOINT,
		Region:   config.GetConfig.S3_REGION,
		Credentials: s3.Credentials{
			AccessKey:       config.GetConfig.S3_ACCESS_KEY,
			SecretAccessKey: config.GetConfig.S3_SECRET_KEY,
		},
	})
}
